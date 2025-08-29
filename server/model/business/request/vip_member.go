package request

type RenewCardRequest struct {
	ID         int `json:"id" form:"id" gorm:"comment:会员卡编号"`
	CardID     int `json:"cardId" form:"cardId" gorm:"unique;comment:会员卡号"`   // 客户名
	ComboId    int `json:"comboId" form:"comboId" gorm:"comment:套餐ID;size:9"` // 套餐ID
	Times      int `json:"times" form:"times" gorm:"comment:赠送次数/金额"`         // 管理ID
	Collection int `json:"collection" form:"collection" gorm:"comment:实付金额"`  // 管理角色ID
}
