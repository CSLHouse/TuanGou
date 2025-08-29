package system

import (
	"cooller/server/middleware"
	"cooller/server/model/business"
	businessRes "cooller/server/model/business/response"
	wechatReq "cooller/server/model/wechat/request"
	wechatRes "cooller/server/model/wechat/response"
	"fmt"
	"strconv"
	"time"

	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	"cooller/server/model/system"
	systemReq "cooller/server/model/system/request"
	systemRes "cooller/server/model/system/response"
	"cooller/server/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Login
// @Tags     Base
// @Summary  用户登录
// @Produce   application/json
// @Param    data  body      systemReq.Login                                             true  "用户名, 密码, 验证码"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var l systemReq.Login
	err := c.ShouldBindJSON(&l)
	key := c.ClientIP()

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.LoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 判断验证码是否开启
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool = openCaptcha == 0 || openCaptcha < interfaceToInt(v)

	if !oc || store.Verify(l.CaptchaId, l.Captcha, true) {
		u := &system.SysUser{UserName: l.UserName, Password: l.Password}
		user, err := userService.Login(u)
		if err != nil {
			global.GVA_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
			// 验证码次数+1
			global.BlackCache.Increment(key, 1)
			response.FailWithMessage("用户名不存在或者密码错误", c)
			return
		}
		if user.Enable != 1 {
			global.GVA_LOG.Error("登陆失败! 用户被禁止登录!")
			// 验证码次数+1
			global.BlackCache.Increment(key, 1)
			response.FailWithMessage("用户被禁止登录", c)
			return
		}
		b.TokenNext(c, *user)
		return
	}
	// 验证码次数+1
	global.BlackCache.Increment(key, 1)
	response.FailWithMessage("验证码错误", c)
}

// TokenNext 登录以后签发jwt
func (b *BaseApi) TokenNext(c *gin.Context, user system.SysUser) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		UserName:    user.UserName,
		AuthorityId: user.AuthorityId,
		Telephone:   user.Phone,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.GVA_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.UUID.String()); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, user.UUID.String()); err != nil {
			global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
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
		if err := jwtService.SetRedisJWT(token, user.UUID.String()); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	}
}

// Register
// @Tags     SysUser
// @Summary  用户注册账号
// @Produce   application/json
// @Param    data  body      systemReq.Register                                            true  "用户名, 昵称, 密码, 角色ID"
// @Success  200   {object}  response.Response{data=systemRes.SysUserResponse,msg=string}  "用户注册账号,返回包括用户信息"
// @Router   /user/admin_register [post]
func (b *BaseApi) Register(c *gin.Context) {
	var r systemReq.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var authorities []system.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &system.SysUser{UserName: r.UserName, NickName: r.NickName, Password: r.Password, AvatarUrl: r.AvatarUrl, AuthorityId: r.AuthorityId, Authorities: authorities, Enable: r.Enable, Phone: r.Phone, Email: r.Email}
	userReturn, err := userService.Register(*user)
	if err != nil {
		global.GVA_LOG.Error("注册失败!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册成功", c)
}

// ChangePassword
// @Tags      SysUser
// @Summary   用户修改密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      systemReq.ChangePasswordReq    true  "用户名, 原密码, 新密码"
// @Success   200   {object}  response.Response{msg=string}  "用户修改密码"
// @Router    /user/changePassword [post]
func (b *BaseApi) ChangePassword(c *gin.Context) {
	var req systemReq.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.ChangePasswordVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	u := &system.SysUser{GVA_MODEL: global.GVA_MODEL{ID: uid}, Password: req.Password}
	_, err = userService.ChangePassword(u, req.NewPassword)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败，原密码与当前账户不符", c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// GetUserList
// @Tags      SysUser
// @Summary   分页获取用户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /user/getUserList [post]
func (b *BaseApi) GetUserList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userService.GetUserInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// SetUserAuthority
// @Tags      SysUser
// @Summary   更改用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuth          true  "用户UUID, 角色ID"
// @Success   200   {object}  response.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthority [post]
func (b *BaseApi) SetUserAuthority(c *gin.Context) {
	var sua systemReq.SetUserAuth
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if UserVerifyErr := utils.Verify(sua, utils.SetUserAuthorityVerify); UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	err = userService.SetUserAuthority(userID, sua.AuthorityId)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims := utils.GetUserInfo(c)
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims.AuthorityId = sua.AuthorityId
	if token, err := j.CreateToken(*claims); err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		c.Header("new-token", token)
		c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
		response.OkWithMessage("修改成功", c)
	}
}

// SetUserAuthorities
// @Tags      SysUser
// @Summary   设置用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuthorities   true  "用户UUID, 角色ID"
// @Success   200   {object}  response.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthorities [post]
func (b *BaseApi) SetUserAuthorities(c *gin.Context) {
	var sua systemReq.SetUserAuthorities
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.SetUserAuthorities(sua.ID, sua.AuthorityIds)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// DeleteUser
// @Tags      SysUser
// @Summary   删除用户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "用户ID"
// @Success   200   {object}  response.Response{msg=string}  "删除用户"
// @Router    /user/deleteUser [delete]
func (b *BaseApi) DeleteUser(c *gin.Context) {
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
	jwtId := utils.GetUserID(c)
	if jwtId == reqId.ID {
		response.FailWithMessage("删除失败, 自杀失败", c)
		return
	}
	err = userService.DeleteUser(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// SetUserInfo
// @Tags      SysUser
// @Summary   设置用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /user/setUserInfo [put]
func (b *BaseApi) SetUserInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(user, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if len(user.AuthorityIds) != 0 {
		err = userService.SetUserAuthorities(user.ID, user.AuthorityIds)
		if err != nil {
			global.GVA_LOG.Error("设置失败!", zap.Error(err))
			response.FailWithMessage("设置失败", c)
			return
		}
	}
	err = userService.SetUserInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:     user.NickName,
		AvatarUrl:    user.AvatarUrl,
		Phone:        user.Phone,
		Email:        user.Email,
		SideMode:     user.SideMode,
		Enable:       user.Enable,
		IsMembership: user.IsMembership,
		PayOnline:    user.PayOnline,
	})
	if err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败", c)
		return
	}
	response.OkWithMessage("设置成功", c)
}

// SetSelfInfo
// @Tags      SysUser
// @Summary   设置用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /user/SetSelfInfo [put]
func (b *BaseApi) SetSelfInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user.ID = utils.GetUserID(c)
	err = userService.SetSelfInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:     user.NickName,
		AvatarUrl:    user.AvatarUrl,
		Phone:        user.Phone,
		Email:        user.Email,
		SideMode:     user.SideMode,
		Enable:       user.Enable,
		IsMembership: user.IsMembership,
		PayOnline:    user.PayOnline,
	})
	if err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败", c)
		return
	}
	response.OkWithMessage("设置成功", c)
}

// GetUserInfo
// @Tags      SysUser
// @Summary   获取用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "获取用户信息"
// @Router    /user/getUserInfo [get]
func (b *BaseApi) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	ReqUser, err := userService.GetUserInfo(uuid)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "获取成功", c)
}

// ResetPassword
// @Tags      SysUser
// @Summary   重置用户密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      system.SysUser                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "重置用户密码"
// @Router    /user/resetPassword [post]
func (b *BaseApi) ResetPassword(c *gin.Context) {
	var user system.SysUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.ResetPassword(user.ID)
	if err != nil {
		global.GVA_LOG.Error("重置失败!", zap.Error(err))
		response.FailWithMessage("重置失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("重置成功", c)
}

func (b *BaseApi) WXLogin(c *gin.Context) {
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
func (b *BaseApi) GetWXUserInfo(c *gin.Context) {
	var user wechatReq.UserTag
	err := c.ShouldBindQuery(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	wxUser, err := userService.GetWXAccountByOpenID(user.OpenID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(wxUser, c)
}

func (b *BaseApi) CreateWXUserInfo(c *gin.Context) {
	var userInfo wechatReq.WXUserInfo
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println("---userInfo:", userInfo)

	var wxUser system.SysUser
	//wxUser.OpenId = userInfo.OpenID
	wxUser.NickName = userInfo.NickName
	//wxUser.Gender = userInfo.Gender
	wxUser.AvatarUrl = userInfo.AvatarUrl
	wxUser.AuthorityId = 9528

	err = userService.CreateWXAccount(&wxUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println("---wxUser:", wxUser)

	response.OkWithDetailed(wxUser, "更新成功", c)
}

// WXTokenNext 登录以后签发jwt
func (b *BaseApi) WXTokenNext(c *gin.Context, customer business.Customer) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        customer.UUID,
		ID:          customer.ID,
		NickName:    customer.NickName,
		UserName:    customer.UserName,
		AuthorityId: customer.AuthorityId,
		Telephone:   customer.Telephone,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.GVA_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(businessRes.WXLoginResponse{
			Customer:  customer,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(customer.UUID.String()); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, customer.UUID.String()); err != nil {
			global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(businessRes.WXLoginResponse{
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
		if err := jwtService.SetRedisJWT(token, customer.UUID.String()); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(businessRes.WXLoginResponse{
			Customer:  customer,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	}
}

func (b *BaseApi) WXRefreshLogin(c *gin.Context) {
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
	wxUser, err := memberService.GetVIPMemberByOpenIdWithoutCardList(login.OpenID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	b.WXTokenNext(c, wxUser)
}

// TODO: 优化请求
func (b *BaseApi) ParsePhoneNumber(c *gin.Context) {
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
	var wxUser business.Customer
	wxUser.Telephone = phoneNumber
	wxUser.OpenId = loginInfo.OpenID
	wxUser.AuthorityId = 9528

	err = memberService.CreateCustomerFormWechat(&wxUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(wechatRes.WXPhoneNum{
		PhoneNumber: phoneNumber,
	}, c)
}

// CheckPhoneNumber 查询是否有手机号 true 有 false 无
func (b *BaseApi) CheckPhoneNumber(c *gin.Context) {
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
	if len(userInfo.PhoneNumber) < 11 {
		response.OkWithData(false, c)
		return
	}
	response.OkWithData(true, c)
}

func (b *BaseApi) ResetWXNickName(c *gin.Context) {
	var user business.Customer
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.ResetWXNickName(&user)
	if err != nil {
		global.GVA_LOG.Error("设置昵称失败!", zap.Error(err))
		response.FailWithMessage("设置昵称失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("设置昵称成功", c)
}

// RecordShareScanAccount 记录分享被读取次数
func (b *BaseApi) RecordShareScanAccount(c *gin.Context) {
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

	err = userService.RecordShareScanAccount(&openIdInfo.OpenId)
	if err != nil {
		global.GVA_LOG.Error("记录分享次数失败!", zap.Error(err))
		response.FailWithMessage("记录分享次数失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("记录分享次数成功", c)
}
