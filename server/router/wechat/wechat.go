package wechat

import (
	v1 "cooller/server/api/v1"
	"cooller/server/middleware"
	"github.com/gin-gonic/gin"
)

type WechatRouter struct{}

func (e *WechatRouter) InitWechatRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {

	homePublicRouterWithoutRecord := RouterPub.Group("promo")

	wechatRouter := Router.Group("promo").Use(middleware.OperationRecord())
	homeApi := v1.ApiGroupApp.WechatApiGroup.HomeApi
	{
		wechatRouter.POST("advertise", homeApi.CreateHomeAdvertise)
		wechatRouter.PUT("advertise", homeApi.UpdateHomeAdvertise)
		wechatRouter.DELETE("advertise", homeApi.DeleteHomeAdvertise)
		wechatRouter.PUT("advertiseState", homeApi.UpdateHomeAdvertiseOnlineState)

		wechatRouter.POST("recommendProduct", homeApi.CreateRecommendProducts)
		wechatRouter.PUT("recommendProduct", homeApi.UpdateRecommendProducts)
		wechatRouter.DELETE("recommendProduct", homeApi.DeleteRecommendProducts)
		wechatRouter.POST("updateRecommendSort", homeApi.UpdateRecommendProductSortById)

		wechatRouter.POST("newProduct", homeApi.CreateNewProduct)
		wechatRouter.DELETE("newProduct", homeApi.DeleteNewProduct)
		wechatRouter.PUT("newProductSort", homeApi.UpdateNewProductSort)
		wechatRouter.PUT("newProductState", homeApi.UpdateNewProductState)
	}
	{
		homePublicRouterWithoutRecord.GET("content", homeApi.GetAllWechatContent)
		homePublicRouterWithoutRecord.GET("advertiseList", homeApi.GetHomeAdvertiseList)
		homePublicRouterWithoutRecord.GET("recommendProductList", homeApi.GetRecommendProductList) // 获取推荐商品列表
		homePublicRouterWithoutRecord.GET("recommendProduct", homeApi.GetRecommendProductListByCondition)
		homePublicRouterWithoutRecord.GET("newProductList", homeApi.GetNewProductList)
	}
}
