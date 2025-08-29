package request

import (
	"cooller/server/model/system"
)

// Register User register structure
type Register struct {
	UserName     string `json:"userName" example:"用户名"`
	Password     string `json:"passWord" example:"密码"`
	NickName     string `json:"nickName" example:"昵称"`
	AvatarUrl    string `json:"avatarUrl" example:"头像链接"`
	AuthorityId  int    `json:"authorityId" swaggertype:"string" example:"int 角色id"`
	Enable       int    `json:"enable" swaggertype:"string" example:"int 是否启用"`
	AuthorityIds []int  `json:"authorityIds" swaggertype:"string" example:"[]int 角色id"`
	Phone        string `json:"phone" example:"电话号码"`
	Email        string `json:"email" example:"电子邮箱"`
}

// User login structure
type Login struct {
	UserName  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

// Modify password structure
type ChangePasswordReq struct {
	ID          int    `json:"-"`           // 从 JWT 中提取 user id，避免越权
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}

// Modify  user's auth structure
type SetUserAuth struct {
	AuthorityId int `json:"authorityId"` // 角色ID
}

// Modify  user's auth structure
type SetUserAuthorities struct {
	ID           int
	AuthorityIds []int `json:"authorityIds"` // 角色ID
}

type ChangeUserInfo struct {
	ID           int                   `gorm:"primarykey"`                                                                           // 主键ID
	NickName     string                `json:"nickName" gorm:"comment:用户昵称"`                                                         // 用户昵称
	Phone        string                `json:"phone"  gorm:"comment:用户手机号"`                                                          // 用户手机号
	AuthorityIds []int                 `json:"authorityIds" gorm:"-"`                                                                // 角色ID
	Email        string                `json:"email"  gorm:"comment:用户邮箱"`                                                           // 用户邮箱
	AvatarUrl    string                `json:"avatarUrl" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	SideMode     string                `json:"sideMode"  gorm:"comment:用户侧边主题"`                                                      // 用户侧边主题
	Enable       int                   `json:"enable" gorm:"comment:冻结用户"`                                                           //冻结用户
	IsMembership int                   `json:"isMembership" gorm:"comment:是否是会员制 1是，2不是"`
	PayOnline    int                   `json:"payOnline" gorm:"comment:是否线上付费 1是，2不是"`
	Authorities  []system.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
}
