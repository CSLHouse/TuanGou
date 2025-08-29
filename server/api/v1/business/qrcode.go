package business

import (
	"cooller/server/global"
	"cooller/server/model/business"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	"cooller/server/model/example"
	"cooller/server/utils"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"image/color"
	"net/http"
	"os"
	"path"
)

type QrCodeApi struct{}

func (e *QrCodeApi) GetQrCodeList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	qrcodeList, total, err := qrcodeService.GetQrCodeInfoList(pageInfo, userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     qrcodeList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (e *QrCodeApi) CreateQrCode(c *gin.Context) {
	var qrcode business.QrCode
	err := c.ShouldBindJSON(&qrcode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	qrcode.SysUserId = userId
	err = e.GenerateQrCode(&qrcode, userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = qrcodeService.CreateQrCode(qrcode)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}
func (e *QrCodeApi) GenerateQrCode(qrcodeInfo *business.QrCode, userId int) error {
	if len(qrcodeInfo.Url) > 0 {
		n, err := snowflake.NewNode(1)
		if err != nil {
			global.GVA_LOG.Error("创建id失败!", zap.Error(err))
		}

		uniqueTag := n.Generate().Int64()
		if len(qrcodeInfo.Path) < 1 {
			qrcodeInfo.Path = fmt.Sprintf("./uploads/%d.png", uniqueTag)
		}
		if len(qrcodeInfo.RemoteUrl) > 0 {
			var file example.ExaFileUploadAndDownload
			file.Url = qrcodeInfo.RemoteUrl
			file.ID = qrcodeInfo.UploadId
			if err := fileUploadAndDownloadService.DeleteFile(file); err != nil {
				global.GVA_LOG.Error("删除失败!", zap.Error(err))
				return fmt.Errorf("删除失败!")
			}
		}
		err = qrcode.WriteColorFile(qrcodeInfo.Url, qrcode.Medium, 256, color.White, color.Black, qrcodeInfo.Path)
		if err != nil {
			global.GVA_LOG.Error("生成二维码失败!", zap.Error(err))
			return fmt.Errorf("生成二维码失败:", zap.Error(err))
		}

		var file example.ExaFileUploadAndDownload
		file, err = fileUploadAndDownloadService.UploadFileWithLocationPath(qrcodeInfo.Path, uniqueTag, userId) // 文件上传后拿到文件路径
		if err != nil {
			global.GVA_LOG.Error("修改数据库链接失败!", zap.Error(err))
			return fmt.Errorf("生成二维码失败:", zap.Error(err))
		}
		qrcodeInfo.RemoteUrl = file.Url
		qrcodeInfo.UploadId = file.ID

		liveCodeUrl := fmt.Sprintf("https://cs.coollerbaby.cn/qrcode/scan?id=%d", qrcodeInfo.ID)
		err = qrcode.WriteColorFile(liveCodeUrl, qrcode.Medium, 256, color.White, color.Black, qrcodeInfo.Path)
		if err != nil {
			global.GVA_LOG.Error("生成二维码失败!", zap.Error(err))
			return fmt.Errorf("生成二维码失败:", zap.Error(err))
		}
	}
	return nil
}
func (e *QrCodeApi) DeleteQrCodeById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = qrcodeService.DeleteQrCodeById(reqId.ID, userId)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (e *QrCodeApi) UpdateQrCode(c *gin.Context) {
	var qrcode business.QrCode
	err := c.ShouldBindJSON(&qrcode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(qrcode, utils.ComboVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	qrcode.SysUserId = userId
	err = e.GenerateQrCode(&qrcode, userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = qrcodeService.UpdateQrCode(&qrcode)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *QrCodeApi) DownloadQrCodeFile(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	qrcodeInfo, err := qrcodeService.GetQrCodeById(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	if len(qrcodeInfo.RemoteUrl) < 1 {
		global.GVA_LOG.Error("生成二维码失败!", zap.Error(err))
		response.FailWithMessage("生成二维码失败"+err.Error(), c)
		return
	}

	//fmt.Println(c.Request.URL)
	////filePath := c.Query("path")
	//requestURL := fmt.Sprintf("%v", c.Request.URL)
	//requestURLarray := strings.Split(requestURL, "url=")
	//if len(requestURLarray) < 2 {
	//	global.GVA_LOG.Error("downloadFile 失败 文件地址： " + requestURL)
	//	response.FailWithMessage("downloadFile 失败"+err.Error(), c)
	//	return
	//}
	filePath := qrcodeInfo.Path
	//filePath = "./" + filePath
	//打开文件
	fileTmp, errByOpenFile := os.Open(filePath)
	if errByOpenFile != nil {
		global.GVA_LOG.Error("downloadFile 失败 文件地址：", zap.Error(errByOpenFile))
		response.FailWithMessage("生成二维码失败"+errByOpenFile.Error(), c)
		return
	}
	defer fileTmp.Close()

	//获取文件的名称
	fileName := path.Base(filePath)
	isExist := utils.FileExist(filePath)
	if !isExist {
		global.GVA_LOG.Error("downloadFile 失败 文件不存在：" + qrcodeInfo.Path)
		response.FailWithMessage("downloadFile 失败 文件不存在", c)
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	//强制浏览器下载
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	//浏览器下载或预览
	c.Header("Content-Disposition", "inline;filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.File(filePath)
	return
}

func (e *QrCodeApi) ScanQrCodeFile(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	qrcodeInfo, err := qrcodeService.GetQrCodeById(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"group": qrcodeInfo.RemoteUrl,
	})

	err = qrcodeService.UpdateQrCodeCount(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("更新识别二维码次数识别!", zap.Error(err))
		return
	}
}

func CheckQrCodeExpired() {
	err := qrcodeService.UpdateExpiredQrCodeState()
	if err != nil {
		global.GVA_LOG.Error("过期二维码修改状态!", zap.Error(err))
		return
	}
	//TODO: 过期通知
	//qrcodeInfoList, err := qrcodeService.GetExpiredQrCodeList()
	//if err != nil {
	//	if err != nil {
	//		global.GVA_LOG.Error("获取过期二维码数据识别!", zap.Error(err))
	//		return
	//	}
	//}

}
