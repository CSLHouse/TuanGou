package response

type VipComboResModel struct {
	StoreName  string `json:"storeName" form:"storeName" gorm:"comment:分店名称"`   // 分店名称
	ComboName  string `json:"comboName" form:"comboName" gorm:"comment:套餐名称"`   // 套餐名称
	ComboType  int    `json:"comboType" form:"comboType" gorm:"comment:套餐类型"`   // 套餐类型
	ComboPrice int    `json:"comboPrice" form:"comboPrice" gorm:"comment:套餐价格"` // 套餐价格
	Times      int    `json:"amount" form:"amount" gorm:"comment:天数/次数/金额"`     // 天数/次数/金额
	State      int    `json:"state" form:"state" gorm:"comment:状态"`             // 状态
}
type VipComboResponse struct {
	Combo VipComboResModel `json:"combo"`
}
