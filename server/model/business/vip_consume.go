package business

type Consume struct {
	ID     int    `json:"id" form:"id" gorm:"comment:会员卡序号"`
	CardID string `json:"cardId" form:"cardId" gorm:"comment:会员卡号"` // 客户名
	Number int    `json:"number" form:"number" gorm:"comment:消费次数"` // 管理ID
}

type ConsumeRecord struct {
	ID        int    `json:"Id" gorm:"primary_key;"`                          // 主键ID
	Telephone string `json:"telephone" form:"telephone" gorm:"comment:会员手机号"` // 管理ID
	//ComboId            int      `json:"comboId" form:"comboId" gorm:"comment:套餐ID"`            // 管理ID
	RemainTimes  int     `json:"remainTimes" form:"remainTimes" gorm:"comment:剩余次数/金额"` // 管理ID
	ConsumeTimes int     `json:"consumeTimes" form:"consumeTimes" gorm:"comment:消费次数/金额"`
	PunchDate    string  `json:"punchDate" form:"punchDate" gorm:"comment:登记日期"`
	State        int     `json:"state" form:"state" gorm:"comment:状态"`
	SysUserId    int     `json:"sysUserId" form:"sysUserId" gorm:"comment:管理角色ID"` // 管理角色ID
	Card         VIPCard `json:"member" gorm:"foreignKey:Telephone;references:Telephone;comment:会员信息"`
}

func (ConsumeRecord) TableName() string {
	return "bus_consume_record"
}

//type DouYinRecord struct {
//	ID        int     `json:"Id" gorm:"primary_key;"` // 主键ID
//	Amount    float32 `json:"amount" form:"amount" gorm:"comment:消费次数/金额"`
//	Date      string  `json:"date" form:"date" gorm:"comment:日期"`
//	SysUserId int     `json:"sysUserId" form:"sysUserId" gorm:"comment:管理角色ID"` // 管理角色ID
//}
//
//func (DouYinRecord) TableName() string {
//	return "bus_dou_yin_record"
//}
