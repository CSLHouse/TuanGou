package product

import "cooller/server/service"

type ApiGroup struct {
	FlashApi
	CouponApi
}

var (
	wechatService = service.ServiceGroupApp.WechatServiceGroup.HomeService
	couponService = service.ServiceGroupApp.CouponServiceGroup.CouponService
)
