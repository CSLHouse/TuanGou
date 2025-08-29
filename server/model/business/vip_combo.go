package business

import (
	"cooller/server/global"
)

type VIPCombo struct {
	global.GVA_MODEL
	//ComboId    int    `json:"comboId" gorm:"not null;unique;primary_key;comment:套餐ID;size:90"` // 套餐ID
	ComboName  string `json:"comboName" form:"comboName" gorm:"comment:套餐名称"`   // 套餐名称
	ComboType  int    `json:"comboType" form:"comboType" gorm:"comment:套餐类型"`   // 套餐类型
	ComboPrice int    `json:"comboPrice" form:"comboPrice" gorm:"comment:套餐价格"` // 套餐价格
	Times      int    `json:"times" form:"times" gorm:"comment:天数/次数/金额"`       // 天数/次数/金额
	State      int    `json:"state" form:"state" gorm:"comment:状态"`             // 状态
	SysUserId  int    `json:"sysUserId" form:"sysUserId" gorm:"comment:管理ID"`
}

func (VIPCombo) TableName() string {
	return "bus_combo"
}
