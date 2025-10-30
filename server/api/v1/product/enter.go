package product

import (
	"cooller/server/service"
)

type ApiGroup struct {
	CouponApi
	OrderApi
	ProductApi
	LogisticsApi
}

var (
	wechatService    = service.ServiceGroupApp.WechatServiceGroup.HomeService
	couponService    = service.ServiceGroupApp.ProductServiceGroup.CouponService
	productService   = service.ServiceGroupApp.ProductServiceGroup.ProductService
	accountService   = service.ServiceGroupApp.WechatServiceGroup.AccountService
	jspaymentService = service.ServiceGroupApp.PaymentServiceGroup.PayMentService
	orderService     = service.ServiceGroupApp.ProductServiceGroup.OrderService
)
