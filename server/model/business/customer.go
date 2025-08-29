package business

import (
	"cooller/server/global"
	"github.com/gofrs/uuid/v5"
)

type Customer struct {
	global.GVA_MODEL
	UUID        uuid.UUID  `json:"uuid" gorm:"index;comment:用户UUID"`    // 用户UUID
	UserName    string     `json:"userName" gorm:"index;comment:用户登录名"` // 用户登录名
	Password    string     `json:"-"  gorm:"comment:用户登录密码"`            // 用户登录密码
	NickName    string     `json:"nickName" gorm:"comment:用户昵称"`        // 用户昵称
	OpenId      string     `json:"openId" form:"openId" gorm:"primaryKey"`
	Gender      int        `json:"gender" form:"gender" gorm:"comment:性别"`
	AvatarUrl   string     `json:"avatarUrl" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	Telephone   string     `json:"telephone" form:"telephone" gorm:"primaryKey;comment:用户手机号"`                           // 用户手机号
	Count       uint       `json:"count" form:"count" gorm:"comment:访问次数"`
	Enable      int        `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
	ShareCount  int        `json:"shareCount" form:"shareCount" gorm:"comment:分享次数"`
	CardList    []*VIPCard `json:"cardList" form:"cardList" gorm:"foreignKey:CustomerId;comment:会员卡号"`
	AuthorityId int        `json:"authorityId" gorm:"default:9528;comment:用户角色ID"`
	SysUserId   int        `json:"sysUserId" form:"sysUserId" gorm:"comment:管理ID"`
}

func (Customer) TableName() string {
	return "bus_customer"
}
