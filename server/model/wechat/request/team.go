package request

type DetailsInfo struct {
	UserIds []int `json:"userIds" binding:"required"`
}

type SettlementInfo struct {
	OpenId string  `json:"openId" binding:"required"`
	Amount float32 `json:"amount" binding:"required"`
}

type SettlementSearchInfo struct {
	CreatedAt      string
	UserId         int    `json:"userId" form:"userId" gorm:"comment:用户序号"`
	SettlementNo   string `json:"settlementNo" form:"settlementNo" gorm:"not null;unique;comment:结算单号"`
	SettlementTime string `json:"settlementTime" form:"settlementTime" gorm:"not null;comment:结算时间"`
	Status         int    `json:"status" form:"status" gorm:"not null;default:0;comment:结算状态，0:待处理,1:已完成"`
	Page           int    `json:"page" form:"page"`         // 页码
	PageSize       int    `json:"pageSize" form:"pageSize"` // 每页大小
}
