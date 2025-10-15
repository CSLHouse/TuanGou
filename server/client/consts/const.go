// Copyright 2021 Tencent Inc. All rights reserved.

// Package consts 微信支付 API v3 Go SDK 常量
package consts

import "time"

// 微信支付 API 地址
const (
	WechatPayAPIServer       = "https://api.mch.weixin.qq.com"  // 微信支付 API 地址
	WechatPayAPIServerBackup = "https://api2.mch.weixin.qq.com" // 微信支付 API 备份地址
)

// SDK 相关信息
const (
	Version         = "0.2.18"                     // SDK 版本
	UserAgentFormat = "WechatPay-Go/%s (%s) GO/%s" // UserAgent中的信息
	//UserAgentTest   = "https://zh.wikipedia.org/wiki/User_agent"
)

// HTTP 请求报文 Header 相关常量
const (
	Authorization = "Authorization"  // Header 中的 Authorization 字段
	Accept        = "Accept"         // Header 中的 Accept 字段
	ContentType   = "Content-Type"   // Header 中的 ContentType 字段
	ContentLength = "Content-Length" // Header 中的 ContentLength 字段
	UserAgent     = "User-Agent"     // Header 中的 UserAgent 字段
)

// 常用 ContentType
const (
	ApplicationJSON = "application/json"
	ImageJPG        = "image/jpg"
	ImagePNG        = "image/png"
	VideoMP4        = "video/mp4"
)

// 请求报文签名相关常量
const (
	//SignatureMessageFormat = "%s\n%s\n%d\n%s\n%s\n" // 数字签名原文格式
	SignatureMessageFormat = "%s\n%s\n%d\n%s\n%s\n" // 数字签名原文格式
	// HeaderAuthorizationFormat 请求头中的 Authorization 拼接格式
	HeaderAuthorizationFormat = "%s mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\""
)

// HTTP 应答报文 Header 相关常量
const (
	WechatPayTimestamp = "Wechatpay-Timestamp" // 微信支付回包时间戳
	WechatPayNonce     = "Wechatpay-Nonce"     // 微信支付回包随机字符串
	WechatPaySignature = "Wechatpay-Signature" // 微信支付回包签名信息
	WechatPaySerial    = "Wechatpay-Serial"    // 微信支付回包平台序列号
	RequestID          = "Request-Id"          // 微信支付回包请求ID
)

// 时间相关常量
const (
	FiveMinute     = 5 * 60           // 回包校验最长时间（秒）
	DefaultTimeout = 30 * time.Second // HTTP 请求默认超时时间
)

// 微信支付
const (
	MachID                     = "1722992716"                               // 商户号
	MchCertificateSerialNumber = "41A714460BA39A5DA85A1F315D68FE06AAD9950B" // 商户证书序列号
	MchAPIv3Key                = "u2z3SsLS8EbajEmc4spCCHn7CMYnq6LJ"         // 商户APIv3密钥
	PrivateKey                 = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDoFgamKnaa5Mw2\ndOe9ihDYjii6l4fYBOzxLQLbMqnlIsMO2kVVPxya1kFoHnG0GuvcCmZVg2WRYqh6\nRHW+zGtOEWjKWjOjlJlbDkDS2u6aITnNsIchj0TMoZJ3+ya4yi5kgwLPkDTlU9gZ\nFf53Wx2SPqQe1+Iub5dtwuJVOiqHoc0Tu8W+vjOmOwoEsdoxcpu+VnSqs3vOqi1o\nfnoysVu/1OjN92ywd+qunvbai203TFHqNZdlOXs7D8e51bl+EKUKXHHCSZW/7OVA\nZN53qhOh70cTb5xzojhebD/EPT6rk/46VwMGh8YLpgPds1yFHwrkSqOh41rrm1fs\n3IItF5LRAgMBAAECggEBAJvowAq8l/SyIyWqH0vLBkoUJfRDe3R4ypux+ys5u7w8\nQAEoTE9aGGND76h8WzH8q9mj+t8TaBXkyv1oJBlKxgMp2tWlBu+rEeKDjrSXcpaE\nH4q4Pe8jibpWCEklpPuahtPq1uaNH5u3WYJwrgOMaC8oeaRmroE3YejQdK+AZbW6\nin3u5kqOdNDZW/JL07sekydTUsgLal8jbvZRsDzLAmNELviFDeTGnPfVHcpp1IgA\nqMS7y1vfBoxAPyc11LNJKJtu6vCE2h5yqHrQhkWdjpYwjCNVn7O8zMkDjAIt+Hcg\nf+7RyNQBnE9/m6oJXnXhf6MuOKMrVlOu9DvFWLyIXWECgYEA/gWD29tROr6mZaez\n9yiTiklE9MwxYPcQzPzgu73jIMBZDU8n0djt0+BvZ9EeYrkIF/NcPr+FHj9wKV9w\nj+/ymoXr8LMmsGbxZmAT/g4ii4UB3W3BdklLpm+FjhZ0guMF4euzBaZtq25f1GwQ\n9vEusizN48RWjQSUtpLWAlToIj0CgYEA6eTGQd5fxsYfHxRj4Jk885TrM/SdX6OX\n2kOVbPczJ+yMNQPdI6EHbuho5+x7erUI1IMY75CtGp4bLHrzt5ddFc5e8/wreYje\nPfamfJJ9TuW6K9ZN95fNvAKBrKIbtMmngNMpCH6RlWFk9jqWAixg41AmAZNaoEAx\ngjoX3vN/ICUCgYAfKb4LctR8SHdRcUl6wNeY13RGM+a65pSBEWTgo2MB2ZcTMurH\nq71BiP8h1V3M7rY3efuPFx/VniK8cKD2h5FTs6pGVTQh3/8teAv71vAV/bNE3vok\n8Mj/Gh9gVxDkHcXS+X906f36Uggfn0JViTEZWrXHg4a6th0oaMsobhIsYQKBgCW9\nq5wPVOcPKxBpyt//+gzX98fvbcFz9Vnb0+28Sb3kdo9La1CHeFqWF+9sglQ/iAg6\niziE4NnNr1bTFCaIvxV4smuuQhmfUzUuapjpTlz/xrWyI+ySyzjOMrx3f/8BFw07\nhYCAY9910sPEmlYJcSzczvUsINCA6zw3QYjRQ68dAoGBANgd3XGrp4fXBMNvtEyu\neDrB3OAjNL2mgtWE1/g0Pf2jXenZYZvR5gsIooCTlcuv5yaulYjicTTTA4Va7G+I\negcTkpT7NNnL9yfnUVr/z/MG8DQcIIbpRR/TjKoNrl8kynMoe9o74xsRnQtvdFFc\nt9L1zFH5fIBhOJNI0xHilB8X\n-----END PRIVATE KEY-----"
	PublicKey                  = "sfds"
)
