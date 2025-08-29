package wechat

import (
	"cooller/server/global"
	"cooller/server/model/wechat"
	"gorm.io/gorm"
)

type AccountService struct{}

func (exa *AccountService) CreateWXAccount(e wechat.WXUser) (err error) {
	db := global.GVA_DB.Model(&wechat.WXUser{})
	var wxUser wechat.WXUser
	result := db.Where("open_id = ?", e.OpenId).First(&wxUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			err = global.GVA_DB.Create(&e).Error
			return err
		}
		err = result.Error
	} else {
		err = db.Debug().Where("open_id = ?", e.OpenId).Updates(map[string]interface{}{"session_key": e.SessionKey}).Error
		return err
	}

	return err
}

func (exa *AccountService) UpdateWXAccountInfo(e wechat.WXUser) (err error) {
	db := global.GVA_DB.Model(&wechat.WXUser{})

	err = db.Debug().Where("open_id = ?", e.OpenId).Updates(map[string]interface{}{"nick_name": e.NickName,
		"gender": e.Gender, "avatar_url": e.AvatarUrl, "token": e.Token}).Error
	return err
}

func (exa *AccountService) GetWXAccountByOpenID(openId string) (user wechat.WXUser, err error) {
	err = global.GVA_DB.Where("open_id = ?", openId).First(&user).Error
	return user, err
}

func (exa *AccountService) UpdateWXAccountPhone(openId string, phoneNum string) (err error) {
	db := global.GVA_DB.Model(&wechat.WXUser{})
	err = db.Debug().Where("open_id = ?", openId).UpdateColumn("phone_number", phoneNum).Error
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
