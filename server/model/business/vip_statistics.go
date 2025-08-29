package business

import "cooller/server/global"

// VIPStatement 流水
type VIPStatement struct {
	//Date          string `json:"date" form:"date" gorm:"primary_key;comment:日期;"`
	global.GVA_MODEL
	Recharge      int `json:"recharge" form:"recharge" gorm:"comment:会员卡流水"`          // 套餐名称
	CardNumber    int `json:"cardNumber" form:"cardNumber" gorm:"comment:会员卡单量"`      // 套餐类型
	NewMember     int `json:"newMember" form:"newMember" gorm:"comment:新增会员"`         // 套餐价格
	ConsumeNumber int `json:"consumeNumber" form:"consumeNumber" gorm:"comment:入店统计"` // 天数/次数/金额
	SysUserId     int `json:"sysUserId" form:"sysUserId" gorm:"comment:管理角色ID"`       // 管理角色ID
}

func (VIPStatement) TableName() string {
	return "bus_statement"
}

type VIPStatistics struct {
	TotalStream   float64 `json:"totalStream" form:"totalStream" gorm:"comment:总流水;"`
	TotalOrder    int     `json:"totalOrder" form:"totalOrder" gorm:"comment:会员卡总单量"`
	TotalMember   int     `json:"totalMember" form:"totalMember" gorm:"comment:总会员"`
	TotalConsumer int     `json:"totalConsumer" form:"totalConsumer" gorm:"comment:入店总人数"`
	SysUserId     int     `json:"sysUserId" form:"sysUserId" gorm:"comment:管理角色ID"`
}

func (VIPStatistics) TableName() string {
	return "bus_statistics"
}
