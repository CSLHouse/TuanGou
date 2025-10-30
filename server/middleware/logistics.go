package middleware

import (
	"bytes"
	"cooller/server/global"
	"cooller/server/model/product"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LogisticsClient struct {
	apiBaseURL string       // API请求地址
	token      string       // X-APISpace-Token
	httpClient *http.Client // HTTP客户端（复用连接、设置超时）
}

func NewLogisticsClient() (*LogisticsClient, error) {
	// 初始化HTTP客户端（设置5秒超时，避免长时间阻塞）
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &LogisticsClient{
		apiBaseURL: global.GVA_CONFIG.Logistics.URL,
		token:      global.GVA_CONFIG.Logistics.Token,
		httpClient: httpClient,
	}, nil
}

// DetectExpressCompany：传入运单号，自动识别快递公司编码
// 返回：识别到的CpCode、快递公司名称、错误信息
func (c *LogisticsClient) DetectExpressCompany(mailNo string) (string, string, error) {
	// 1. 构造识别API的请求体
	reqBody := product.ExpressDetectReq{MailNo: mailNo}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("识别请求体编码失败：%w", err)
	}

	// 2. 创建POST请求
	req, err := http.NewRequest(http.MethodPost, global.GVA_CONFIG.Logistics.CpURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", "", fmt.Errorf("创建识别请求失败：%w", err)
	}

	// 3. 设置请求头（与物流查询API共用Token）
	req.Header.Set("X-APISpace-Token", global.GVA_CONFIG.Logistics.Token)
	req.Header.Set("Content-Type", "application/json")

	// 4. 发送请求并解析响应
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("发送识别请求失败：%w", err)
	}
	defer resp.Body.Close()

	// 5. 校验HTTP状态码
	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("识别API返回错误：状态码%d，内容：%s", resp.StatusCode, string(errBody))
	}

	// 6. 解析识别结果
	var detectResp product.ExpressDetectResp
	if err := json.NewDecoder(resp.Body).Decode(&detectResp); err != nil {
		return "", "", fmt.Errorf("识别响应解析失败：%w", err)
	}

	// 7. 业务逻辑校验
	if !detectResp.Success {
		return "", "", fmt.Errorf("运单号识别失败：%s", detectResp.ErrorMsg)
	}
	if len(detectResp.ExpressCompanyList) == 0 {
		return "", "", errors.New("未识别到对应的快递公司，请确认运单号是否正确")
	}

	// 8. 返回第一个识别结果（通常最准确）
	firstResult := detectResp.ExpressCompanyList[0]
	return firstResult.CpCode, firstResult.CompanyName, nil
}

// QueryLogistics：发送物流查询请求
// 参数：reqBody（请求参数）
// 返回：API响应数据、错误信息
func (c *LogisticsClient) QueryLogistics(reqBody product.ApiRequestBody) (*product.ApiResponse, error) {
	// 步骤1：编码请求体为JSON格式
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("请求体JSON编码失败：%w", err)
	}

	// 步骤2：创建HTTP POST请求
	req, err := http.NewRequest(
		http.MethodPost,
		c.apiBaseURL,
		bytes.NewBuffer(reqJSON), // 请求体（JSON格式）
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%w", err)
	}

	// 步骤3：设置请求头（必填：X-APISpace-Token、Content-Type）
	req.Header.Set("X-APISpace-Token", global.GVA_CONFIG.Logistics.Token)
	req.Header.Set("Content-Type", "application/json")

	// 步骤4：发送请求并获取响应
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败：%w", err)
	}
	defer resp.Body.Close() // 确保响应体被关闭，避免资源泄漏

	// 步骤5：校验HTTP响应状态码（200为成功）
	if resp.StatusCode != http.StatusOK {
		// 读取错误响应内容（便于调试）
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回非成功状态码：%d，响应内容：%s", resp.StatusCode, string(errBody))
	}

	// 步骤6：解析响应体为ApiResponse结构体
	var apiResp product.ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("响应体JSON解析失败：%w", err)
	}

	// 步骤7：校验API业务逻辑是否成功（根据API的Success字段判断）
	if !apiResp.Success {
		return &apiResp, fmt.Errorf("物流查询失败：%s（TraceId：%s）", apiResp.ErrorMsg, apiResp.TraceId)
	}

	return &apiResp, nil
}
