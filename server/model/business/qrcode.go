package business

import "cooller/server/global"

type QrCode struct {
	global.GVA_MODEL
	Title     string `json:"title" gorm:"comment:二维码标题"` // 文件名
	Url       string `json:"url" gorm:"comment:地址"`
	Count     int    `json:"count" form:"count" gorm:"comment:扫码次数"`
	Path      string `json:"path" gorm:"comment:文件地址"`      // 文件地址
	RemoteUrl string `json:"remoteUrl" gorm:"comment:文件地址"` // 服务器即远程文件地址
	UploadId  int    `json:"uploadId" gorm:"comment:地址"`    // 对应的exa_file_upload_and_downloads表的ID
	IsExpired int    `json:"isExpired" gorm:"comment:是否过期 0未过期1过期;"`
	SysUserId int    `json:"sysUserId" form:"sysUserId" gorm:"comment:管理ID"`
}

func (QrCode) TableName() string {
	return "bus_qr_code"
}
