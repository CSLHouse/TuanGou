package wechat

import "cooller/server/router/product"

type RouterGroup struct {
	FlashRouter
	AccountRouter
	product.OrderRouter
	WechatRouter
}
