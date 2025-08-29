package response

import "cooller/server/model/wechat"

// PrepayResponse
type PrepayResponse struct {
	// 预支付交易会话标识
	PrepayId *string `json:"prepay_id"`
}

// PrepayWithRequestPaymentResponse 预下单ID，并包含了调起支付的请求参数
type PrepayWithRequestPaymentResponse struct {
	// 预支付交易会话标识
	PrepayId *string `json:"prepay_id"` // revive:disable-line:var-naming
	// 应用ID
	Appid *string `json:"appId"`
	// 时间戳
	TimeStamp *string `json:"timeStamp"`
	// 随机字符串
	NonceStr *string `json:"nonceStr"`
	// 订单详情扩展字符串
	Package *string `json:"package"`
	// 签名方式
	SignType *string `json:"signType"`
	// 签名
	PaySign *string `json:"paySign"`
}

type GenerateOrderResponse struct {
	OrderId int                               `json:"orderId"`
	Payment *PrepayWithRequestPaymentResponse `json:"payment"`
}

type GenerateOrderDetailResponse struct {
	Order   wechat.Order                      `json:"order"`
	Payment *PrepayWithRequestPaymentResponse `json:"payment"`
}
