package example

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"mime/multipart"
	"strconv"
	"strings"

	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/example"
	exampleReq "cooller/server/model/example/request"
	"cooller/server/utils/upload"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Upload
//@description: 创建文件上传记录
//@param: file model.ExaFileUploadAndDownload
//@return: error

func (e *FileUploadAndDownloadService) Upload(file *example.ExaFileUploadAndDownload) error {
	return global.GVA_DB.Create(file).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: FindFile
//@description: 查询文件记录
//@param: id uint
//@return: model.ExaFileUploadAndDownload, error

func (e *FileUploadAndDownloadService) FindFile(id int) (example.ExaFileUploadAndDownload, error) {
	var file example.ExaFileUploadAndDownload
	err := global.GVA_DB.Where("id = ?", id).First(&file).Error
	return file, err
}

func (e *FileUploadAndDownloadService) FindFileByFileId(fileId int64) (example.ExaFileUploadAndDownload, error) {
	var file example.ExaFileUploadAndDownload
	err := global.GVA_DB.Where("file_id = ?", fileId).First(&file).Error
	return file, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFile
//@description: 删除文件记录
//@param: file model.ExaFileUploadAndDownload
//@return: err error

func (e *FileUploadAndDownloadService) DeleteFile(file example.ExaFileUploadAndDownload) (err error) {
	var fileFromDb example.ExaFileUploadAndDownload
	s := strings.Split(file.Url, "/")
	fileName := s[len(s)-1]
	fileIdL := strings.Split(fileName, ".")

	fileId, err := strconv.ParseInt(fileIdL[0], 10, 64)
	if err != nil {
		return err
	}

	fileFromDb, err = e.FindFileByFileId(fileId)
	if err != nil {
		return
	}
	oss := upload.NewOss()
	if err = oss.DeleteFile(fileFromDb.Key); err != nil {
		return errors.New("文件删除失败")
	}
	err = global.GVA_DB.Debug().Where("id = ?", file.ID).Unscoped().Delete(&file).Error
	return err
}

func (e *FileUploadAndDownloadService) DeleteFiles(filesInfo *exampleReq.FilesRequest) (err error) {

	var filesPathList []string
	var filesIds []int
	for _, file := range filesInfo.FilesInfo {
		// TODO: 截取绝对路径
		filesPathList = append(filesPathList, file.Key)
		filesIds = append(filesIds, file.ID)
	}

	oss := upload.NewOss()
	if err = oss.DeleteFiles(filesPathList); err != nil {
		return errors.New("文件删除失败")
	}

	var file example.ExaFileUploadAndDownload
	err = global.GVA_DB.Where("id in ?", filesIds).Unscoped().Delete(&file).Error
	return err
}

// EditFileName 编辑文件名或者备注
func (e *FileUploadAndDownloadService) EditFileName(file example.ExaFileUploadAndDownload) (err error) {
	var fileFromDb example.ExaFileUploadAndDownload
	return global.GVA_DB.Where("id = ?", file.ID).First(&fileFromDb).Update("name", file.Name).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: list interface{}, total int64, err error

func (e *FileUploadAndDownloadService) GetFileRecordInfoList(info request.PageInfo, userId int) (fileLists []example.ExaFileUploadAndDownload, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	db := global.GVA_DB.Model(&example.ExaFileUploadAndDownload{})
	if len(keyword) > 0 {
		db = db.Where("name LIKE ? and sys_user_id = ?", "%"+keyword+"%", userId)
	}
	err = db.Where("sys_user_id = ?", userId).Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Where("sys_user_id = ?", userId).Order("updated_at desc").Find(&fileLists).Error
	return fileLists, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, noSave string
//@return: file model.ExaFileUploadAndDownload, err error

func (e *FileUploadAndDownloadService) UploadFile(header *multipart.FileHeader, noSave string, userId int) (file example.ExaFileUploadAndDownload, err error) {
	n, err := snowflake.NewNode(1)
	if err != nil {
		global.GVA_LOG.Error("创建id失败!", zap.Error(err))
	}
	s := strings.Split(header.Filename, ".")
	tag := s[len(s)-1]

	uuid := n.Generate()
	fileName := uuid.String()
	header.Filename = fmt.Sprintf("%s.%s", fileName, tag)
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header, userId)
	if uploadErr != nil {
		panic(err)
	}

	f := example.ExaFileUploadAndDownload{
		Url:       filePath,
		Name:      header.Filename,
		Tag:       tag,
		Key:       key,
		SysUserId: userId,
		FileId:    uuid.Int64(),
	}
	if noSave == "0" {
		err = e.Upload(&f)
		fmt.Println(f)
		return f, err
	}
	return f, nil
}

func (e *FileUploadAndDownloadService) CheckFile(fileName string, userId int) bool {
	var file example.ExaFileUploadAndDownload
	err := global.GVA_DB.Where("name = ? and sys_user_id = ?", fileName, userId).First(&file).Error
	fmt.Println(err)
	if err == nil {
		return true
	}
	return false
}

func (e *FileUploadAndDownloadService) UploadFileWithLocationPath(locationPath string, uniqueTag int64, userId int) (file example.ExaFileUploadAndDownload, err error) {
	oss := upload.NewOss()
	fileName := fmt.Sprintf("%d.png", uniqueTag)
	filePath, key, uploadErr := oss.UploadFileWithLocationPath(locationPath, fileName, userId)
	if uploadErr != nil {
		panic(err)
	}

	s := strings.Split(fileName, ".")

	f := example.ExaFileUploadAndDownload{
		Url:       filePath,
		Name:      fileName,
		Tag:       s[len(s)-1],
		Key:       key,
		SysUserId: userId,
		FileId:    uniqueTag,
	}
	err = e.Upload(&f)
	return f, err
}
