package product

import (
	v1 "cooller/server/api/v1"
	"cooller/server/middleware"
	"github.com/gin-gonic/gin"
)

type OrderRouter struct{}

func (s *OrderRouter) InitOrderRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	couponPublicRouterWithoutRecord := RouterPub.Group("order")
	wxOrderRouter := Router.Group("order").Use(middleware.OperationRecord())
	wxOrderApi := v1.ApiGroupApp.ProductApiGroup.OrderApi
	{
		wxOrderRouter.POST("generateConfirmOrder", wxOrderApi.GenerateConfirmOrder)
		wxOrderRouter.POST("generateOrder", wxOrderApi.GenerateOrder)
		wxOrderRouter.POST("paySuccess", wxOrderApi.PaySuccess)
		wxOrderRouter.POST("cancelOrder", wxOrderApi.CancelOrders)
		wxOrderRouter.POST("settingUpdate", wxOrderApi.UpdateOrderSetting)
		wxOrderRouter.POST("closeOrders", wxOrderApi.CloseOrders)
		wxOrderRouter.DELETE("delete", wxOrderApi.DeleteOrders)
		wxOrderRouter.POST("update/receiverInfo", wxOrderApi.UpdateOrderReceiverInfo)
		wxOrderRouter.POST("update/moneyInfo", wxOrderApi.UpdateOrderMoneyInfo)
		wxOrderRouter.POST("update/note", wxOrderApi.UpdateOrderNote)
		wxOrderRouter.POST("update/complete", wxOrderApi.UpdateOrderCompletedStatus)
		wxOrderRouter.POST("cart", wxOrderApi.CreateProductCart)
		wxOrderRouter.PUT("cart", wxOrderApi.UpdateProductCartQuantity)
		wxOrderRouter.DELETE("cart", wxOrderApi.DeleteProductCartById)
		wxOrderRouter.DELETE("cart/clear", wxOrderApi.ClearProductCart)
		wxOrderRouter.DELETE("carts", wxOrderApi.DeleteProductCartByIds)
		wxOrderRouter.POST("tmpCart", wxOrderApi.CreateProductTmpCart)
		wxOrderRouter.PUT("logistics", wxOrderApi.UpdateOrderLogistics)
	}
	{
		couponPublicRouterWithoutRecord.GET("detail", wxOrderApi.GetOrderDetail)
		couponPublicRouterWithoutRecord.GET("list", wxOrderApi.GetOrderList)
		couponPublicRouterWithoutRecord.GET("setting", wxOrderApi.GetOrderSetting)
		couponPublicRouterWithoutRecord.GET("cart/list", wxOrderApi.GetProductCartList)
	}
}
