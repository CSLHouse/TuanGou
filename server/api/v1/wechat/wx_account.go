package wechat

import (
	"cooller/server/global"
	"cooller/server/middleware"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	"cooller/server/model/system"
	systemReq "cooller/server/model/system/request"
	"cooller/server/model/wechat"
	wechatReq "cooller/server/model/wechat/request"
	wechatRes "cooller/server/model/wechat/response"
	"cooller/server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type WXAccountApi struct{}

func (b *WXAccountApi) WXLogin(c *gin.Context) {
	var login wechatReq.WXLogin
	err := c.ShouldBindJSON(&login)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(login.Code) < 1 {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//发送jscode，获得用户的open_id
	wechatClient := middleware.NewWechatClient(nil)
	wxMap, err := wechatClient.WXLogin(login.Code)
	if err != nil {
		global.GVA_LOG.Error("登录失败!", zap.Error(err))
		response.FailWithMessage("登录失败", c)
		return
	}
	var session wechatRes.WXLoginRes
	session.OpenID = wxMap["openid"]
	if len(session.OpenID) < 1 {
		global.GVA_LOG.Error("登录失败!", zap.Error(err))
		response.FailWithMessage("登录失败", c)
		return
	}

	//var wxUser wechat.WXUser
	//wxUser.OpenId = wxMap["openid"]
	//wxUser.SessionKey = wxMap["session_key"]
	////wxUser.Token = b.CreateToken(wxMap["openid"], userInfo.NickName)
	//wxUser.Count = 1
	//fmt.Println("---wxUser:", wxUser)
	//err = userService.CreateWXAccount(wxUser)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	response.OkWithDetailed(session, "获取成功", c)
}

// GetWXUserInfo 获取用户详情
func (b *WXAccountApi) GetWXUserInfo(c *gin.Context) {
	var user wechatReq.UserTag
	err := c.ShouldBindQuery(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	wxUser, err := wxAccountServer.GetWXAccountByOpenID(user.OpenID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(wxUser, c)
}

func (b *WXAccountApi) CreateWXUserInfo(c *gin.Context) {
	var userInfo wechatReq.WXUserInfo
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println("---userInfo:", userInfo)
	var upperUserId uint16
	if len(userInfo.InviteCode) > 0 {
		upperUserInfo, err := accountService.GetWXAccountByInviteCode(userInfo.InviteCode)
		if err != nil {
			global.GVA_LOG.Error(err.Error(), zap.Error(err))
			response.FailWithMessage("请检查邀请码是否错误", c)
			return
		}
		//通过邀请码反推ID, 需保证邀请码生成器的CHARSET和长度一致
		g := utils.NewGenerator[uint16](utils.CHARSET, 6)
		upperUserId = g.Decode(userInfo.InviteCode)
		if upperUserId < 1 || upperUserInfo.ID != int(upperUserId) {
			global.GVA_LOG.Error("邀请码错误!", zap.Error(err))
			response.FailWithMessage("邀请码错误", c)
			return
		}
	}
	var wxUser wechat.WXUser
	wxUser.OpenId = userInfo.OpenID
	wxUser.UserName = userInfo.NickName
	wxUser.City = userInfo.City
	wxUser.AvatarUrl = userInfo.AvatarUrl
	wxUser.Telephone = userInfo.Telephone
	wxUser.CaptainId = int(upperUserId)
	wxUser.AuthorityId = 9528

	wxUser, err = accountService.CreateWXAccount(&wxUser)
	if err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println("---wxUser:", wxUser)
	// 生成邀请码
	if wxUser.ID > 0 {
		g := utils.NewGenerator[uint16](utils.CHARSET, 6)
		// 通过一个现有的非负整数ID生成对应的邀请码
		curUID := uint16(wxUser.ID)
		code, err := g.Encode(curUID)
		if err != nil {
			global.GVA_LOG.Error("生成邀请码失败!", zap.Error(err))
			response.FailWithMessage("生成邀请码失败", c)
			return
		}
		wxUser.InviteCode = code

		// 记录
		var inviteCodeData wechat.TeamRecord
		inviteCodeData.UserId = wxUser.ID
		inviteCodeData.CaptainId = int(upperUserId) // 队长id
		inviteCodeData.IsActivated = 0
		inviteCodeData.InviteCode = userInfo.InviteCode
		err = teamService.CreateInviteCodeRecode(&inviteCodeData)
		if err != nil {
			global.GVA_LOG.Error("存储邀请码记录失败!", zap.Error(err))
			response.FailWithMessage("存储邀请码记录失败", c)
			return
		}

		err = accountService.UpdateWXAccount(wxUser)
		if err != nil {
			global.GVA_LOG.Error("更新邀请码失败!", zap.Error(err))
			response.FailWithMessage("更新邀请码失败", c)
			return
		}
	}

	b.WXTokenNext(c, wxUser)
}

// WXTokenNext 登录以后签发jwt
func (b *WXAccountApi) WXTokenNext(c *gin.Context, customer wechat.WXUser) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		ID:          customer.ID,
		UserName:    customer.UserName,
		AuthorityId: customer.AuthorityId,
		Telephone:   customer.Telephone,
		CaptainId:   customer.CaptainId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.GVA_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(wechatRes.WXLoginResponse{
			Customer:  customer,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(customer.OpenId); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, customer.OpenId); err != nil {
			global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(wechatRes.WXLoginResponse{
			Customer:  customer,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := jwtService.SetRedisJWT(token, customer.OpenId); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(wechatRes.WXLoginResponse{
			Customer:  customer,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	}
}

func (b *WXAccountApi) WXRefreshLogin(c *gin.Context) {
	var login wechatReq.UserTag
	err := c.ShouldBindJSON(&login)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(login.OpenID) < 1 {
		response.FailWithMessage(err.Error(), c)
		return
	}
	wxUser, err := wxAccountServer.GetWXAccountByOpenID(login.OpenID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	b.WXTokenNext(c, wxUser)
}

// TODO: 优化请求
func (b *WXAccountApi) ParsePhoneNumber(c *gin.Context) {
	var loginInfo wechatReq.WXPhoneNumber
	err := c.ShouldBindJSON(&loginInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(loginInfo, utils.WxRegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	accessToken, ok := global.BlackCache.Get("access_token")
	if !ok || accessToken == nil {
		wechatClient := middleware.NewWechatClient(nil)
		wxMap, err := wechatClient.GetWXAccessToken()
		if err != nil {
			global.GVA_LOG.Error("登录失败!", zap.Error(err))
			response.FailWithMessage("登录失败", c)
			return
		}
		accessToken = wxMap["access_token"]
		openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
		global.BlackCache.Set("access_token", accessToken, time.Second*time.Duration(openCaptchaTimeOut))
	}

	//httpClient := http.Client{}
	wechatClient := middleware.NewWechatClient(nil)
	wxMap, err := wechatClient.GetWXTelephone(accessToken.(string), loginInfo.Code)

	if err != nil {
		global.GVA_LOG.Error("登录失败!", zap.Error(err))
		response.FailWithMessage("登录失败", c)
		return
	}

	var phoneNumber = wxMap.PhoneInfo.PurePhoneNumber
	if len(phoneNumber) < 1 {
		global.GVA_LOG.Error("登录失败!", zap.Error(err))
		response.FailWithMessage("登录失败", c)
		return
	}
	//var wxUser wechat.WXUser
	//wxUser.Telephone = phoneNumber
	//wxUser.OpenId = loginInfo.OpenID
	//wxUser.AuthorityId = 9528
	//
	//err = accountService.CreateWXAccount(&wxUser)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}

	response.OkWithData(wechatRes.WXPhoneNum{
		PhoneNumber: phoneNumber,
	}, c)
}

// CheckPhoneNumber 查询是否有手机号 true 有 false 无
func (b *WXAccountApi) CheckPhoneNumber(c *gin.Context) {
	var loginInfo wechatReq.UserTag
	err := c.ShouldBindQuery(&loginInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(loginInfo.OpenID) < 1 {
		fmt.Println("----login:", loginInfo)
		response.FailWithMessage(err.Error(), c)
		return
	}

	userInfo, err := accountService.CheckWXAccountPhone(loginInfo.OpenID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(userInfo.Telephone) < 11 {
		response.OkWithData(false, c)
		return
	}
	response.OkWithData(true, c)
}

func (b *WXAccountApi) ResetWXNickName(c *gin.Context) {
	var user wechat.WXUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = accountService.ResetWXNickName(&user)
	if err != nil {
		global.GVA_LOG.Error("设置昵称失败!", zap.Error(err))
		response.FailWithMessage("设置昵称失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("设置昵称成功", c)
}

// RecordShareScanAccount 记录分享被读取次数
func (b *WXAccountApi) RecordShareScanAccount(c *gin.Context) {
	var openIdInfo request.OpenIdInfo
	err := c.ShouldBindJSON(&openIdInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(openIdInfo.OpenId) < 1 {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = accountService.RecordShareScanAccount(&openIdInfo.OpenId)
	if err != nil {
		global.GVA_LOG.Error("记录分享次数失败!", zap.Error(err))
		response.FailWithMessage("记录分享次数失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("记录分享次数成功", c)
}

func (b *WXAccountApi) CreateMemberReceiveAddress(c *gin.Context) {
	var address wechat.MemberReceiveAddress
	err := c.ShouldBindJSON(&address)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	if userId < 1 {
		response.FailWithMessage("获取UserID失败", c)
		return
	}
	address.UserId = userId
	err = accountService.CreateMemberReceiveAddress(&address)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetMemberReceiveAddressList 获取会员地址列表
func (b *WXAccountApi) GetMemberReceiveAddressList(c *gin.Context) {
	var user wechatReq.UserTag
	err := c.ShouldBindQuery(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	if userId < 1 {
		response.FailWithMessage("获取UserID失败", c)
		return
	}
	addressList, err := accountService.GetMemberReceiveAddressList(userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(addressList, c)
}

// GetMemberReceiveAddressById 获取会员地址列表
func (b *WXAccountApi) GetMemberReceiveAddressById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	address, err := accountService.GetMemberReceiveAddressById(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(address, c)
}

// UpdateMemberReceiveAddress 更新会员地址列表
func (b *WXAccountApi) UpdateMemberReceiveAddress(c *gin.Context) {
	var address wechat.MemberReceiveAddress
	err := c.ShouldBindJSON(&address)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = accountService.UpdateMemberReceiveAddress(&address)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteMemberReceiveAddress 删除会员地址列表
func (b *WXAccountApi) DeleteMemberReceiveAddress(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = accountService.DeleteMemberReceiveAddress(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
