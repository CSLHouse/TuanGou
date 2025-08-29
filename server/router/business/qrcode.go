package business

import (
	v1 "cooller/server/api/v1"
	"github.com/gin-gonic/gin"
)

type QrCodeRouter struct{}

func (e *QrCodeRouter) InitQrCodeRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	qrcodePublicRouterWithoutRecord := RouterPub.Group("qrcode")

	qrcodeApi := v1.ApiGroupApp.BusinessApiGroup.QrCodeApi
	{
		qrcodePublicRouterWithoutRecord.GET("list", qrcodeApi.GetQrCodeList)
		qrcodePublicRouterWithoutRecord.POST("create", qrcodeApi.CreateQrCode)
		qrcodePublicRouterWithoutRecord.PUT("update", qrcodeApi.UpdateQrCode)
		qrcodePublicRouterWithoutRecord.DELETE("delete", qrcodeApi.DeleteQrCodeById)
		qrcodePublicRouterWithoutRecord.GET("download", qrcodeApi.DownloadQrCodeFile)
		qrcodePublicRouterWithoutRecord.GET("scan", qrcodeApi.ScanQrCodeFile)
	}
}
