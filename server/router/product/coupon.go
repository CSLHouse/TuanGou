package product

import (
	v1 "cooller/server/api/v1"
	"cooller/server/middleware"
	"github.com/gin-gonic/gin"
)

type CouponRouter struct{}

func (s *CouponRouter) InitCouponRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	couponPublicRouterWithoutRecord := RouterPub.Group("coupon")
	couponRouter := Router.Group("coupon").Use(middleware.OperationRecord())
	couponApi := v1.ApiGroupApp.ProductApiGroup.CouponApi
	{
		couponRouter.POST("create", couponApi.CreateCoupon)
		couponRouter.PUT("update", couponApi.UpdateCoupon)
		couponRouter.DELETE("delete", couponApi.DeleteCouponById)
		couponRouter.POST("add", couponApi.AddUserCoupon)
	}
	{

		couponPublicRouterWithoutRecord.GET("list", couponApi.GetCouponList)
		couponPublicRouterWithoutRecord.GET("listWithState", couponApi.GetCouponListWithState)
		couponPublicRouterWithoutRecord.GET("details", couponApi.GetCouponById)
		couponPublicRouterWithoutRecord.GET("couponHistory", couponApi.GetCouponHistoryList)
		couponPublicRouterWithoutRecord.GET("listByProduct", couponApi.GetCouponListByProduct)
	}
}
