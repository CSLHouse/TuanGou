package router

import (
	"cooller/server/router/business"
	"cooller/server/router/example"
	"cooller/server/router/pay"
	"cooller/server/router/product"
	"cooller/server/router/system"
	"cooller/server/router/wechat"
)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
	Business business.RouterGroup
	Wechat   wechat.RouterGroup
	Pay      pay.RouterGroup
	Product  product.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
