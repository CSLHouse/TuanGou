package v1

import (
	"cooller/server/api/v1/business"
	"cooller/server/api/v1/example"
	"cooller/server/api/v1/pay"
	"cooller/server/api/v1/product"
	"cooller/server/api/v1/system"
	"cooller/server/api/v1/wechat"
)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	BusinessApiGroup business.ApiGroup
	WechatApiGroup   wechat.ApiGroup
	PayApiGroup      pay.ApiGroup
	ProductApiGroup  product.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
