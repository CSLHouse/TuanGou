package request

type SearchInfoCoupon struct {
	Page      int    `json:"page" form:"page"`         // 页码
	PageSize  int    `json:"pageSize" form:"pageSize"` // 每页大小
	UseStatus int    `json:"useStatus" form:"useStatus"`
	OrderSn   string `json:"orderSn" form:"orderSn"`
	CouponId  string `json:"couponId" form:"couponId"`
}
