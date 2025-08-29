package wechat

import (
	"cooller/server/global"
	"github.com/gofrs/uuid/v5"
)

// WXUser 微信用户
type WXUser struct {
	global.GVA_MODEL
	UUID        uuid.UUID `json:"uuid" gorm:"index;comment:用户UUID"`    // 用户UUID
	UserName    string    `json:"userName" gorm:"index;comment:用户登录名"` // 用户登录名
	OpenId      string    `json:"openId" form:"openId" gorm:"primaryKey"`
	SessionKey  string    `json:"sessionKey" form:"sessionKey"`
	NickName    string    `json:"nickName" gorm:"not null;comment:昵称;size:100"`
	AvatarUrl   string    `json:"avatarUrl" form:"avatarUrl" gorm:"comment:头像"`
	PhoneNumber string    `json:"phoneNumber" form:"phoneNumber" gorm:"comment:手机号"`
	Gender      int       `json:"gender" form:"gender" gorm:"comment:性别"`
	Token       string    `json:"token" form:"token" gorm:"comment:token;size:500"`
	Permissions uint      `json:"permissions" form:"permissions" gorm:"comment:权限"`
	Count       uint      `json:"count" form:"count" gorm:"comment:访问次数"`
	LatestTime  string    `json:"latestTime" form:"latestTime" gorm:"comment:上次登录时间;size:500"`
	AuthorityId int       `json:"authorityId" gorm:"default:888;comment:用户角色ID"`   // 用户角色ID
	Enable      int       `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
}

func (WXUser) TableName() string {
	return "wx_user"
}

// MemberReceiveAddress 会员收货地址表
type MemberReceiveAddress struct {
	global.GVA_MODEL
	UserId        int    `json:"userId" `
	Name          string `json:"name" gorm:"not null;comment:收货人名称;size:100"`
	PhoneNumber   string `json:"phoneNumber" gorm:"not null；comment:手机号;size:11"`
	DefaultStatus uint   `json:"defaultStatus" gorm:"comment:是否默认;0->是；1->否"`
	PostCode      string `json:"postCode" gorm:"comment:邮政编码"`
	Province      string `json:"province" gorm:"comment:省份/直辖市"`
	City          string `json:"city" gorm:"not null；comment:城市"`
	Region        string `json:"region" gorm:"not null；comment:区"`
	DetailAddress string `json:"detailAddress" gorm:"comment:详细地址(楼层、门牌号)"`

	//SysUserAuthorityID uint   `json:"sys_user_authority_id" form:"sys_user_authority_id" gorm:"comment:管理角色ID"`
}

func (MemberReceiveAddress) TableName() string {
	return "ums_member_receive_address"
}

// MemberPrice 会员价格
type MemberPrice struct {
	global.GVA_MODEL
	ProductId       int     `json:"productId" gorm:"null;default null"`
	MemberLevelId   int     `json:"memberLevelId" gorm:"null;default null;"`
	MemberPrice     float32 `json:"memberPrice" gorm:"null;default null;comment:会员价格;"`
	MemberLevelName string  `json:"price" gorm:"null;default null;"`
}

func (MemberPrice) TableName() string {
	return "pms_member_price"
}
