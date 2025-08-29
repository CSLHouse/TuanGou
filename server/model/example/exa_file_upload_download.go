package example

import (
	"cooller/server/global"
)

type ExaFileUploadAndDownload struct {
	global.GVA_MODEL
	Name      string `json:"name" gorm:"comment:文件名"`  // 文件名
	Url       string `json:"url" gorm:"comment:文件地址"`  // 文件地址
	Tag       string `json:"tag" gorm:"comment:文件标签"`  // 文件标签
	Key       string `json:"key" gorm:"comment:编号"`    // 编号
	FileId    int64  `json:"fileId" gorm:"comment:编号"` // 编号
	SysUserId int    `json:"sysUserId" form:"sysUserId" gorm:"comment:管理ID"`
}

func (ExaFileUploadAndDownload) TableName() string {
	return "exa_file_upload_and_downloads"
}
