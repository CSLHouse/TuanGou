package product

type RouterGroup struct {
	CouponRouter
	ProductRouter
	OrderRouter
	LogisticsRouter
}
