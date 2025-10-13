package wechat

import (
	"cooller/server/global"
	"cooller/server/model/wechat"
	"gorm.io/gorm"
)

type AccountService struct{}

func (exa *AccountService) CreateWXAccount(e *wechat.WXUser) (user wechat.WXUser, err error) {
	db := global.GVA_DB.Model(&wechat.WXUser{})
	var wxUser wechat.WXUser
	result := db.Where("open_id = ?", e.OpenId).First(&wxUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			err = global.GVA_DB.Create(&e).Error
			wxUser = *e
			return wxUser, err
		}
		err = result.Error
	} else {
		err = db.Debug().Where("open_id = ?", e.OpenId).Updates(map[string]interface{}{
			"session_key": e.SessionKey, "avatar_url": e.AvatarUrl,
			"user_name": e.UserName, "city": e.City, "telephone": e.Telephone,
		}).Error
		wxUser.SessionKey = e.SessionKey
		wxUser.AvatarUrl = e.AvatarUrl
		wxUser.UserName = e.UserName
		wxUser.City = e.City
		wxUser.Telephone = e.Telephone
		return wxUser, err
	}

	return wxUser, err
}

func (exa *AccountService) UpdateWXAccountInfo(e wechat.WXUser) (err error) {
	db := global.GVA_DB.Model(&wechat.WXUser{})

	err = db.Debug().Where("open_id = ?", e.OpenId).Updates(map[string]interface{}{"user_name": e.UserName,
		"gender": e.Gender, "avatar_url": e.AvatarUrl, "token": e.Token}).Error
	return err
}

func (exa *AccountService) UpdateWXAccount(e wechat.WXUser) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *AccountService) GetWXAccountByOpenID(openId string) (user wechat.WXUser, err error) {
	err = global.GVA_DB.Where("open_id = ?", openId).First(&user).Error
	return user, err
}

func (exa *AccountService) GetWXAccountByID(userId int) (user wechat.WXUser, err error) {
	err = global.GVA_DB.Where("id = ?", userId).First(&user).Error
	return user, err
}

func (exa *AccountService) GetWXAccountByInviteCode(inviteCode string) (user wechat.WXUser, err error) {
	err = global.GVA_DB.Where("invite_code = ?", inviteCode).First(&user).Error
	return user, err
}

func (exa *AccountService) UpdateWXAccountPhone(openId string, phoneNum string) (err error) {
	db := global.GVA_DB.Model(&wechat.WXUser{})
	err = db.Debug().Where("open_id = ?", openId).UpdateColumn("phone_number", phoneNum).Error
	return err
}

func (exa *AccountService) ResetWXNickName(e *wechat.WXUser) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *AccountService) CheckWXAccountPhone(openId string) (wxUser wechat.WXUser, err error) {
	db := global.GVA_DB.Model(&wechat.WXUser{})
	err = db.Where("open_id = ?", openId).First(&wxUser).Error
	return wxUser, err
}

func (exa *AccountService) CreateMemberReceiveAddress(e *wechat.MemberReceiveAddress) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *AccountService) GetMemberReceiveAddressById(id int) (e wechat.MemberReceiveAddress, err error) {
	var address wechat.MemberReceiveAddress
	err = global.GVA_DB.Where("id = ?", id).Find(&address).Error
	return address, err
}

func (exa *AccountService) GetMemberReceiveAddressList(userId int) (addressList []wechat.MemberReceiveAddress, err error) {
	err = global.GVA_DB.Where("user_id = ?", userId).Find(&addressList).Error
	return addressList, err
}

func (exa *AccountService) UpdateMemberReceiveAddress(e *wechat.MemberReceiveAddress) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *AccountService) DeleteMemberReceiveAddress(id int) (err error) {
	var address wechat.MemberReceiveAddress
	err = global.GVA_DB.Where("id = ?", id).Delete(&address).Error
	return err
}

func (exa *AccountService) RecordShareScanAccount(openId *string) (err error) {
	err = global.GVA_DB.Debug().Where("open_id = ?", openId).UpdateColumn("share_count", gorm.Expr("share_count+?", 1)).Error
	return err
}
