package wechat

import (
	"cooller/server/service"
)

type ApiGroup struct {
	WXAccountApi
	FlashApi
	HomeApi
}

var (
	wechatService                = service.ServiceGroupApp.WechatServiceGroup.HomeService
	accountService               = service.ServiceGroupApp.WechatServiceGroup.AccountService
	orderService                 = service.ServiceGroupApp.WechatServiceGroup.OrderService
	jspaymentService             = service.ServiceGroupApp.PaymentServiceGroup.PayMentService
	fileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
	productService               = service.ServiceGroupApp.ProductServiceGroup.ProductService
)
