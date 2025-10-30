package product

// ExpressDetectReq：运单号识别API的请求体
type ExpressDetectReq struct {
	MailNo string `json:"mailNo"` // 仅需传入运单号
}

// ExpressDetectResp：运单号识别API的响应体
type ExpressDetectResp struct {
	TraceId            string                `json:"traceId"`
	Success            bool                  `json:"success"`            // 请求是否成功
	ExpressCompanyList []ExpressDetectResult `json:"expressCompanyList"` // 识别结果列表（可能多个候选）
	ErrorMsg           string                `json:"errorMsg,omitempty"` // 错误信息
}

// ExpressDetectResult：单条识别结果（包含快递公司编码和名称）
type ExpressDetectResult struct {
	CpCode      string `json:"cpCode"`      // 快递公司编码（关键：用于后续物流查询）
	CompanyName string `json:"companyName"` // 快递公司名称（如“顺丰速运”）
}

type ApiRequestBody struct {
	CpCode    string `json:"cpCode"`              // 快递公司编码（如EMS、SF等）
	MailNo    string `json:"mailNo"`              // 快递运单号
	Tel       string `json:"tel,omitempty"`       // 收件人手机号（可选）
	OrderType string `json:"orderType,omitempty"` // 时间排序方式（可选）
}

// 2. 响应体结构体：对应API返回的物流数据结构
type ApiResponse struct {
	TraceId        string          `json:"traceId"`        // 跟踪ID（用于调试）
	Success        bool            `json:"success"`        // 请求是否成功
	LogisticsTrace *LogisticsTrace `json:"logisticsTrace"` // 物流详细信息（请求成功时返回）
	// 若API返回错误信息，可新增ErrorMsg字段接收（根据实际API调整）
	ErrorMsg string `json:"errorMsg,omitempty"`
}

// LogisticsTrace：物流总体信息
type LogisticsTrace struct {
	CpCode                   string                     `json:"cpCode"`                   // 快递公司编码
	MailNo                   string                     `json:"mailNo"`                   // 运单号
	TheLastTime              string                     `json:"theLastTime"`              // 最新物流时间
	TheLastMessage           string                     `json:"theLastMessage"`           // 最新物流消息
	LogisticsCompanyName     string                     `json:"logisticsCompanyName"`     // 快递公司名称
	Courier                  string                     `json:"courier"`                  // 配送员
	TakeTime                 string                     `json:"takeTime"`                 // 全程用时
	CourierPhone             string                     `json:"courierPhone"`             // 配送员电话
	LogisticsStatusDesc      string                     `json:"logisticsStatusDesc"`      // 物流状态描述
	LogisticsStatus          string                     `json:"logisticsStatus"`          // 物流状态（如运输中、已签收）
	CpMobile                 string                     `json:"cpMobile"`                 // 公司电话
	CpUrl                    string                     `json:"cpUrl"`                    // 公司网址
	LogisticsTraceDetailList []LogisticsTraceDetailList `json:"logisticsTraceDetailList"` // 物流明细列表

}

// LogisticsTraceDetailList：单条物流明细
type LogisticsTraceDetailList struct {
	AreaCode           string `json:"areaCode"`           // 区域编码
	AreaName           string `json:"areaName"`           // 区域名称（如XX市XX区）
	Desc               string `json:"desc"`               // 物流描述（如“已收取快件”）
	Time               int64  `json:"time"`               // 时间戳（毫秒级）
	SubLogisticsStatus string `json:"subLogisticsStatus"` // 子物流状态
	LogisticsStatus    string `json:"logisticsStatus"`    // 物流状态
}

type QueryResponse struct {
	CpCode          string          `json:"cpCode"`          // 快递公司编码
	Type            string          `json:"type"`            // 快递公司名称
	No              string          `json:"no"`              // 运单号
	LogisticsStatus string          `json:"logisticsStatus"` // 物流状态（如运输中、已签收）
	TakeTime        string          `json:"takeTime"`        // 全程用时
	DataList        []LogisticsInfo `json:"dataList"`
}

type LogisticsInfo struct {
	Time     string `json:"time"`     // 时间戳（毫秒级）
	Status   string `json:"status"`   // 物流状态
	Content  string `json:"content"`  // 物流描述（如“已收取快件”）
	AreaName string `json:"areaName"` // 区域名称（如XX市XX区）
}
