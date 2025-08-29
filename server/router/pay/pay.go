package pay

import (
	v1 "cooller/server/api/v1"
	"github.com/gin-gonic/gin"
)

type PayRouter struct{}

func (e *PayRouter) InitPayRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	payPublicRouterWithoutRecord := RouterPub.Group("pay")
	payRouterWithoutRecord := Router.Group("pay")
	payApi := v1.ApiGroupApp.PayApiGroup.PayApi
	{
		payRouterWithoutRecord.POST("generateOrder", payApi.GenerateOrder)
		payRouterWithoutRecord.GET("detail", payApi.GetOrderDetail)
		//payRouterWithoutRecord.POST("cancelOrder", payApi.CancelOrder)
	}
	{
		payPublicRouterWithoutRecord.POST("notify", payApi.OrderNotify)

	}
}
