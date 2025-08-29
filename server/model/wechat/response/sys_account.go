package response

type WXLoginRes struct {
	ExpireIn int    `json:"expires_in"`
	OpenID   string `json:"openid"`
	UnionId  string `json:"unionid"`
}

type WXPhoneNum struct {
	PhoneNumber string `json:"phoneNumber"`
}

//type WXUserInfo struct {
//	OpenID      string `json:"openid"`
//	NickName    string `json:"nickName" gorm:"not null;comment:昵称;size:100"`
//	AvatarUrl   string `json:"avatarUrl" form:"avatarUrl" gorm:"comment:头像"`
//	PhoneNumber string `json:"phoneNumber" form:"phoneNumber" gorm:"comment:手机号"`
//	Gender      uint   `json:"gender" form:"gender" gorm:"column:性别"`
//	Token       string `json:"token" form:"token" gorm:"comment:token;size:500"`
//	Permissions uint   `json:"permissions" form:"permissions" gorm:"column:权限"`
//}
