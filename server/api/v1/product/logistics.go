package product

import (
	"cooller/server/global"
	"cooller/server/middleware"
	"cooller/server/model/common/response"
	"cooller/server/model/product"
	"cooller/server/utils"
	date_conversion "cooller/server/utils/timer"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogisticsApi struct{}

func (e *LogisticsApi) QueryLogisticsInfo(c *gin.Context) {
	var searchInfo product.ApiRequestBody
	err := c.ShouldBindJSON(&searchInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	if userId <= 0 {
		response.FailWithMessage("Not get userId!", c)
		return
	}

	// 步骤1：初始化物流Client
	client, err := middleware.NewLogisticsClient()
	if err != nil {
		global.GVA_LOG.Error("Client初始化失败!", zap.Error(err))
		response.FailWithMessage("Client初始化失败!", c)
		return
	}
	// 步骤2：自动识别快递公司（核心：无需用户手动输入CpCode）
	fmt.Printf("正在识别运单号【%s】所属快递公司...\n", searchInfo.MailNo)
	cpCode, cpName, err := client.DetectExpressCompany(searchInfo.MailNo)
	if err != nil {
		fmt.Printf("快递公司识别失败：%v\n", err)
		return
	}
	fmt.Printf("识别成功：快递公司【%s】（编码：%s）\n", cpName, cpCode)

	// 步骤3：构造物流查询参数
	reqBody := product.ApiRequestBody{
		CpCode:    cpCode,
		MailNo:    searchInfo.MailNo,
		Tel:       searchInfo.Tel,
		OrderType: searchInfo.OrderType,
	}

	// 步骤3：发送查询请求
	resp, err := client.QueryLogistics(reqBody)
	if err != nil {
		global.GVA_LOG.Error("物流查询失败!", zap.Error(err))
		response.FailWithMessage("物流查询失败!", c)
		return
	}

	fmt.Println("=======================================")
	fmt.Printf("查询成功（TraceId：%s）\n", resp.TraceId)
	fmt.Printf("快递公司：%s（编码：%s）\n", resp.LogisticsTrace.LogisticsCompanyName, resp.LogisticsTrace.CpCode)
	fmt.Printf("运单号：%s\n", resp.LogisticsTrace.MailNo)
	fmt.Printf("当前物流状态：%s（%s）\n", resp.LogisticsTrace.LogisticsStatus, resp.LogisticsTrace.LogisticsStatusDesc)
	fmt.Printf("最新物流消息：%s（%s）\n", resp.LogisticsTrace.TheLastMessage, date_conversion.FormatTime(resp.LogisticsTrace.LogisticsTraceDetailList[0].Time))
	fmt.Println("\n物流明细：")
	for i, detail := range resp.LogisticsTrace.LogisticsTraceDetailList {
		fmt.Printf("%d. 时间：%s | 区域：%s | 状态：%s\n",
			i+1,
			date_conversion.FormatTime(detail.Time),
			detail.AreaName,
			detail.Desc,
		)
	}
	fmt.Println("=======================================")

	var data product.QueryResponse
	data.No = resp.LogisticsTrace.MailNo
	data.CpCode = resp.LogisticsTrace.CpCode
	data.Type = resp.LogisticsTrace.LogisticsCompanyName
	dataList := make([]product.LogisticsInfo, 0)

	for i := len(resp.LogisticsTrace.LogisticsTraceDetailList) - 1; i >= 0; i-- {
		detail := resp.LogisticsTrace.LogisticsTraceDetailList[i]
		var info product.LogisticsInfo
		info.Status = detail.LogisticsStatus
		info.Time = date_conversion.FormatTime(detail.Time)
		info.Content = detail.Desc
		info.AreaName = detail.AreaName
		dataList = append(dataList, info)
	}
	//for _, detail := range resp.LogisticsTrace.LogisticsTraceDetailList {
	//	var info product.LogisticsInfo
	//	info.Status = detail.LogisticsStatus
	//	info.Time = date_conversion.FormatTime(detail.Time)
	//	info.Content = detail.Desc
	//	info.AreaName = detail.AreaName
	//	dataList = append(dataList, info)
	//}
	data.DataList = dataList
	response.OkWithData(data, c)
}
