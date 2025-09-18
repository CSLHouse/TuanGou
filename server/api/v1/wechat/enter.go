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
	fileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
	productService               = service.ServiceGroupApp.ProductServiceGroup.ProductService
	wxAccountServer              = service.ServiceGroupApp.WechatServiceGroup.AccountService
	jwtService                   = service.ServiceGroupApp.SystemServiceGroup.JwtService
	couponService                = service.ServiceGroupApp.ProductServiceGroup.CouponService
)
