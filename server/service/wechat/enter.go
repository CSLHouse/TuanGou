package wechat

import "cooller/server/service/product"

type ServiceGroup struct {
	HomeService
	AccountService
	product.OrderService
}
