package pay

import (
	"cooller/server/service"
)

type ApiGroup struct {
	PayApi
}

var (
	jspaymentService = service.ServiceGroupApp.PaymentServiceGroup.PayMentService
	accountService   = service.ServiceGroupApp.WechatServiceGroup.AccountService
	productService   = service.ServiceGroupApp.ProductServiceGroup.ProductService
	orderService     = service.ServiceGroupApp.ProductServiceGroup.OrderService
)
