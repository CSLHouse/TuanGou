package wechat

import (
	v1 "cooller/server/api/v1"
	"cooller/server/middleware"
	"github.com/gin-gonic/gin"
)

type FlashRouter struct{}

func (s *FlashRouter) InitFlashRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	flashPublicRouterWithoutRecord := RouterPub.Group("flash")
	flashRouter := Router.Group("flash").Use(middleware.OperationRecord())
	flashApi := v1.ApiGroupApp.WechatApiGroup.FlashApi
	{
		flashRouter.POST("create", flashApi.CreateFlashPromotion)
		flashRouter.PUT("update", flashApi.UpdateFlashPromotion)
		flashRouter.DELETE("delete", flashApi.DeleteFlashPromotionById)
		flashRouter.POST("updateStatus", flashApi.UpdateFlashPromotionStatus)
		flashRouter.POST("createFlashProductRelation", flashApi.CreateFlashPromotionProductRelation)
		flashRouter.PUT("updateFlashProductRelation", flashApi.UpdateFlashPromotionProductRelation)
		flashRouter.DELETE("deleteFlashProductRelation", flashApi.DeleteFlashPromotionProductRelationById)
		flashRouter.POST("createFlashSession", flashApi.CreateFlashSession)
		flashRouter.PUT("updateFlashSession", flashApi.UpdateFlashSession)
		flashRouter.DELETE("deleteFlashSession", flashApi.DeleteFlashSessionById)
		flashRouter.POST("updateFlashSessionStatus", flashApi.UpdateFlashSessionStatus)
	}
	{
		flashPublicRouterWithoutRecord.GET("list", flashApi.GetFlashPromotionList)
		flashPublicRouterWithoutRecord.GET("flashProductRelationList", flashApi.GetFlashPromotionProductRelationList)
		flashPublicRouterWithoutRecord.GET("flashSessionList", flashApi.GetFlashSessionList)
		flashPublicRouterWithoutRecord.GET("flashSessionSelectList", flashApi.GetFlashSessionSelectList)
	}
}
