package system

import (
	"cooller/server/global"
	"github.com/gofrs/uuid/v5"
)

type SysUser struct {
	global.GVA_MODEL
	UUID     uuid.UUID `json:"uuid" gorm:"index;comment:用户UUID"`    // 用户UUID
	UserName string    `json:"userName" gorm:"index;comment:用户登录名"` // 用户登录名
	Password string    `json:"-"  gorm:"comment:用户登录密码"`            // 用户登录密码
	NickName string    `json:"nickName" gorm:"comment:用户昵称"`        // 用户昵称
	//OpenId      string         `json:"openId" form:"openId" gorm:"primaryKey"`
	//Gender      int            `json:"gender" form:"gender" gorm:"comment:性别"`
	SideMode     string         `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`                                          // 用户侧边主题
	AvatarUrl    string         `json:"avatarUrl" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	BaseColor    string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	ActiveColor  string         `json:"activeColor" gorm:"default:#1890ff;comment:活跃颜色"`                                      // 活跃颜色
	AuthorityId  int            `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	Authority    SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Authorities  []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	Phone        string         `json:"phone"  gorm:"comment:用户手机号"`                        // 用户手机号
	Email        string         `json:"email"  gorm:"comment:用户邮箱"`                         // 用户邮箱
	Position     int            `json:"position" form:"position" gorm:"comment:职位 1店长 2店员"` // 职位
	Count        uint           `json:"count" form:"count" gorm:"comment:访问次数"`
	Enable       int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
	WriteOffNum  int            `json:"writeOffNum" gorm:"default:0;comment:核销次数"`       //核销次数
	DiscountNum  int            `json:"discountNum" gorm:"default:0;comment:核销次数"`       //优惠次数
	IsMembership int            `json:"isMembership" gorm:"default:1;comment:是否是会员制 1是，2不是"`
	PayOnline    int            `json:"payOnline" gorm:"default:2;comment:是否线上付费 1是，2不是"`
}

func (SysUser) TableName() string {
	return "sys_users"
}
