package request

type DetailsInfo struct {
	UserIds []int `json:"userIds" binding:"required"`
}

type SettlementInfo struct {
	OpenId string  `json:"openId" binding:"required"`
	Amount float32 `json:"amount" binding:"required"`
}
