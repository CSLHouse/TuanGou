package business

import "cooller/server/service"

type ApiGroup struct {
	QrCodeApi
}

var (
	qrcodeService                = service.ServiceGroupApp.BusinessServiceGroup.QrCodeService
	fileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
)
