package wechat

import (
	"cooller/server/global"
)

// WXUser 微信用户
type WXUser struct {
	global.GVA_MODEL
	//UUID        uuid.UUID `json:"uuid" gorm:"index;comment:用户UUID"`    // 用户UUID
	UserName   string `json:"userName" gorm:"index;comment:用户登录名"` // 用户登录名
	OpenId     string `json:"openId" form:"openId" gorm:"primaryKey"`
	SessionKey string `json:"sessionKey" form:"sessionKey"`
	//NickName        string `json:"nickName" gorm:"not null;comment:昵称;size:100"`
	AvatarUrl       string `json:"avatarUrl" form:"avatarUrl" gorm:"comment:头像"`
	Telephone       string `json:"telephone" form:"telephone" gorm:"comment:手机号"`
	Gender          int    `json:"gender" form:"gender" gorm:"comment:性别"`
	Token           string `json:"token" form:"token" gorm:"comment:token;size:500"`
	Permissions     uint   `json:"permissions" form:"permissions" gorm:"comment:权限"`
	Count           uint   `json:"count" form:"count" gorm:"comment:访问次数"`
	City            string `json:"city" form:"city" gorm:"comment:城市"`
	LatestTime      string `json:"latestTime" form:"latestTime" gorm:"comment:上次登录时间;size:500"`
	AuthorityId     int    `json:"authorityId" gorm:"default:888;comment:用户角色ID"` // 用户角色ID
	SysUserId       int    `json:"sysUserId" form:"sysUserId" gorm:"comment:管理ID"`
	Enable          int    `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
	CaptainId       int    `json:"captainId" form:"captainId" gorm:"not null;comment:队长Id;"`
	InviteCode      string `json:"inviteCode" form:"inviteCode" gorm:"not null;comment:邀请码;size:6"`
	IsFirstPurchase int    `json:"isFirstPurchase" gorm:"default:1;comment:是否首购 0：是 1否"`
}

func (WXUser) TableName() string {
	return "wx_user"
}
