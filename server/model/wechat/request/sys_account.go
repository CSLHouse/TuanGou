package request

type WXLogin struct {
	Code string `json:"code"`
}

type WXUserInfo struct {
	OpenID    string `json:"openId" binding:"required"`
	NickName  string `json:"nickName" binding:"required"`
	AvatarUrl string `json:"avatarUrl" binding:"required"`
	Gender    int    `json:"gender"`
}

type WXPhoneNumber struct {
	OpenID string `json:"openId" binding:"required"`
	Code   string `json:"code" binding:"required"`
	//EncryptedData string `json:"encryptedData" binding:"required"`
	//Iv            string `json:"iv" binding:"required"`
	//SessionKey    string `json:"sessionKey" binding:"required"`
}

// UserTag ShouldBindQuery解析url参数 必须加上form
type UserTag struct {
	OpenID string `form:"openId"  json:"openId" binding:"required"`
}
