package product

import (
	"cooller/server/service"
)

type ApiGroup struct {
	CouponApi
	OrderApi
	ProductApi
}

var (
	wechatService    = service.ServiceGroupApp.WechatServiceGroup.HomeService
	couponService    = service.ServiceGroupApp.ProductServiceGroup.CouponService
	orderService     = service.ServiceGroupApp.WechatServiceGroup.OrderService
	productService   = service.ServiceGroupApp.ProductServiceGroup.ProductService
	accountService   = service.ServiceGroupApp.WechatServiceGroup.AccountService
	jspaymentService = service.ServiceGroupApp.PaymentServiceGroup.PayMentService
)
