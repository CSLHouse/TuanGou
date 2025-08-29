package business

import (
	"cooller/server/global"
)

type VIPCard struct {
	global.GVA_MODEL
	CardId      string   `json:"cardId" form:"cardId" gorm:"comment:会员卡号"` // 会员卡号
	Telephone   string   `json:"telephone" form:"telephone" gorm:"comment:用户手机号"`
	UserName    string   `json:"userName" gorm:"index;comment:会员卡用户名"`       // 用户登录名
	ComboId     int      `json:"comboId" form:"comboId" gorm:"comment:套餐ID"` // 管理ID
	Combo       VIPCombo `json:"combo" gorm:"foreignKey:ComboId;references:ID;comment:套餐"`
	RemainTimes int      `json:"remainTimes" form:"remainTimes" gorm:"comment:剩余次数/金额"` // 管理ID
	StartDate   string   `json:"startDate" form:"startDate" gorm:"comment:开始日期"`
	Deadline    string   `json:"deadline" form:"deadline" gorm:"comment:截止日期"` // 管理ID
	State       int      `json:"state" form:"state" gorm:"comment:状态"`         // 管理角色ID
	IsNew       bool     `json:"isNew" form:"isNew" gorm:"comment:新会员"`
	Collection  int      `json:"collection" form:"collection" gorm:"comment:实付金额"` // 管理角色ID
	SysUserId   int      `json:"sysUserId" form:"sysUserId" gorm:"comment:管理ID"`
	Customer    Customer `json:"customer" gorm:"foreignKey:CustomerId;comment:消费者"`
	CustomerId  int      `json:"customerId" form:"customerId" gorm:"comment:消费者ID"`
	StoreName   string   `json:"storeName" form:"storeName" gorm:"comment:所在商店名称"`
	Tmp         int      `json:"tmp" form:"tmp" gorm:"comment:是否是临时会员，0不是，1是"`
}

func (VIPCard) TableName() string {
	return "bus_vip_card"
}
