package wechat

import (
	v1 "cooller/server/api/v1"
	"cooller/server/middleware"
	"github.com/gin-gonic/gin"
)

type TeamRouter struct{}

func (e *TeamRouter) InitTeamRouter(Router *gin.RouterGroup) {

	teamRouter := Router.Group("team").Use(middleware.OperationRecord())
	teamRouterWithoutRecord := Router.Group("team")
	teamApi := v1.ApiGroupApp.WechatApiGroup.TeamApi
	{

		//teamRouter.POST("create", teamApi.CreateProduct)
		teamRouter.POST("settlement", teamApi.GenerateSettlement)
		teamRouter.POST("settlementUpdate", teamApi.UpdateSettlementStatus)
	}
	{
		teamRouterWithoutRecord.GET("teamList", teamApi.GetTeamRecordList)
		teamRouterWithoutRecord.POST("detail", teamApi.GetTeamConsumeDetails)
		teamRouterWithoutRecord.GET("reward", teamApi.GetTeamReward)
		teamRouterWithoutRecord.GET("settlementList", teamApi.GetTeamSettlementList)
	}
}
