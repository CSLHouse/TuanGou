package product

import (
	v1 "cooller/server/api/v1"
	//"cooller/server/middleware"
	"github.com/gin-gonic/gin"
)

type LogisticsRouter struct{}

func (s *LogisticsRouter) InitLogisticsRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	logisticsPublicRouterWithoutRecord := RouterPub.Group("logistics")
	//logisticsRouterRouter := Router.Group("logistics").Use(middleware.OperationRecord())
	logisticsApi := v1.ApiGroupApp.ProductApiGroup.LogisticsApi
	{
		logisticsPublicRouterWithoutRecord.POST("info", logisticsApi.QueryLogisticsInfo)

	}
}
