package response

import (
	"cooller/server/model/system"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	User      system.SysUser `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expiresAt"`
}

type PhoneInfo struct {
	CountryCode     string `json:"countryCode"`
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
}

type PhoneModel struct {
	Errcode   int               `json:"errcode"`
	Errmsg    string            `json:"errmsg"`
	PhoneInfo PhoneInfo         `json:"phone_info"`
	watermark map[string]string `json:"watermark"`
}
