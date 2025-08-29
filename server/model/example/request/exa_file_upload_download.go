package response

import "cooller/server/model/example"

type ExaFilesRequest struct {
	Files []example.ExaFileUploadAndDownload `json:"files"`
}

type FilesInfo struct {
	ID  int    `json:"id"`
	Key string `json:"key" gorm:"comment:编号"` // 编号
}
type FilesRequest struct {
	FilesInfo []FilesInfo `json:"filesInfo"`
}
