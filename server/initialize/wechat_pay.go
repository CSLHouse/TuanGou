package initialize

import (
	"cooller/server/client/consts"
	"cooller/server/global"
	wechatPay "github.com/go-pay/gopay/wechat/v3"
	"go.uber.org/zap"
)

func WeChatPay() {
	payClient, err := wechatPay.NewClientV3(consts.MachID, consts.MchCertificateSerialNumber, consts.MchAPIv3Key, consts.PrivateKey)
	if err != nil {
		global.GVA_LOG.Error("初始化NewClientV3失败!", zap.Error(err))
		return
	}
	global.GVA_WECHAT_PAY_CLIENT = payClient
	// 注意：以下两种自动验签方式二选一
	// 微信支付公钥自动同步验签（新微信支付用户推荐）
	err = payClient.AutoVerifySignByPublicKey([]byte(consts.PublicKey), "微信支付公钥ID，不能删除 PUB_KEY_ID_ 前缀，否则会出错")
	if err != nil {
		global.GVA_LOG.Error("初始化NewClientV3失败!", zap.Error(err))
		return
	}
}
