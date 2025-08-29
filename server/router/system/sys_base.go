package system

import (
	v1 "cooller/server/api/v1"
	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		baseRouter.POST("login", baseApi.Login)
		baseRouter.POST("captcha", baseApi.Captcha)
		baseRouter.POST("wxLogin", baseApi.WXLogin)
		baseRouter.GET("wxUserInfo", baseApi.GetWXUserInfo)
		baseRouter.POST("wxUserInfo", baseApi.CreateWXUserInfo)
		baseRouter.POST("wxRefreshLogin", baseApi.WXRefreshLogin)
		baseRouter.POST("phoneNumber", baseApi.ParsePhoneNumber)
		baseRouter.GET("checkPhone", baseApi.CheckPhoneNumber)
		baseRouter.POST("recordShare", baseApi.RecordShareScanAccount)
	}
	return baseRouter
}
