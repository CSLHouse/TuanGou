package service

import (
	"cooller/server/service/business"
	"cooller/server/service/example"
	"cooller/server/service/pay"
	"cooller/server/service/product"
	"cooller/server/service/system"
	"cooller/server/service/wechat"
)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	ExampleServiceGroup  example.ServiceGroup
	BusinessServiceGroup business.ServiceGroup
	WechatServiceGroup   wechat.ServiceGroup
	PaymentServiceGroup  pay.ServiceGroup
	ProductServiceGroup  product.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
