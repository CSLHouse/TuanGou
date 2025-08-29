package wechat

import (
	v1 "cooller/server/api/v1"
	"github.com/gin-gonic/gin"
)

type AccountRouter struct{}

func (s *AccountRouter) InitAccountRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	wxBaseRouter := Router.Group("base")
	wxBaseApi := v1.ApiGroupApp.WechatApiGroup.WXAccountApi
	{
		//wxBaseRouter.POST("wxLogin", wxBaseApi.WXLogin)
		//wxBaseRouter.GET("userInfo", wxBaseApi.GetUserInfo)
		//wxBaseRouter.POST("userInfo", wxBaseApi.UpdateUserInfo)
		wxBaseRouter.POST("address", wxBaseApi.CreateMemberReceiveAddress)
		wxBaseRouter.GET("addressList", wxBaseApi.GetMemberReceiveAddressList)
		wxBaseRouter.PUT("address", wxBaseApi.UpdateMemberReceiveAddress)
		wxBaseRouter.DELETE("address", wxBaseApi.DeleteMemberReceiveAddress)
		wxBaseRouter.GET("address", wxBaseApi.GetMemberReceiveAddressById)
	}
	return wxBaseRouter
}
