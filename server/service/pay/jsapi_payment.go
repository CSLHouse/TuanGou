package pay

import (
	"context"
	"cooller/server/client"
	"cooller/server/client/consts"
	"cooller/server/client/option"
	"cooller/server/model/pay"
	payRequest "cooller/server/model/pay/request"
	payRes "cooller/server/model/pay/response"
	"cooller/server/utils"
	"fmt"
	"log"
	"strconv"
	"time"
)

type PayMentService struct{}

// PrepayWithRequestPayment Jsapi支付下单，并返回调起支付的请求参数
func (a *PayMentService) PrepayWithRequestPayment(req payRequest.PrepayRequest) (resp payRes.PrepayWithRequestPaymentResponse, result *client.APIResult, err error) {
	ctx := context.Background()
	opts := []client.ClientOption{
		option.WithWechatPayAutoAuthCipher(),
	}
	defaultClient, err := client.NewClient(opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}
	jsapi := JsapiApiService{Client: defaultClient}
	prepayResp, result, err := jsapi.Prepay(ctx, req)
	if err != nil {
		return resp, result, err
	}

	resp.PrepayId = prepayResp.PrepayId
	resp.SignType = utils.String("RSA")
	resp.Appid = req.Appid
	resp.TimeStamp = utils.String(strconv.FormatInt(time.Now().Unix(), 10))
	nonce, err := utils.GenerateNonce()
	if err != nil {
		return resp, nil, fmt.Errorf("generate request for payment err:%s", err.Error())
	}
	resp.NonceStr = utils.String(nonce)
	resp.Package = utils.String("prepay_id=" + *prepayResp.PrepayId)
	message := fmt.Sprintf("%s\n%s\n%s\n%s\n", *resp.Appid, *resp.TimeStamp, *resp.NonceStr, *resp.Package)
	signatureResult, err := jsapi.Client.Sign(ctx, message)
	if err != nil {
		return resp, nil, fmt.Errorf("generate sign for payment err:%s", err.Error())
	}
	resp.PaySign = utils.String(signatureResult.Signature)
	return resp, result, nil
}

func (a *PayMentService) QueryOrderByOutTradeNo(req payRequest.QueryOrderByOutTradeNoRequest, prepayId string) (transaction *pay.Transaction, resp *payRes.PrepayWithRequestPaymentResponse, result *client.APIResult, err error) {
	ctx := context.Background()
	opts := []client.ClientOption{
		option.WithWechatPayAutoAuthCipher(),
	}
	defaultClient, err := client.NewClient(opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}
	jsapi := JsapiApiService{Client: defaultClient}
	queryResp, result, err := jsapi.QueryOrderByOutTradeNo(ctx, req)
	if err != nil {
		return nil, nil, result, err
	}
	transaction = queryResp
	if *transaction.TradeState == "NOTPAY" {
		resp = new(payRes.PrepayWithRequestPaymentResponse)
		resp.PrepayId = utils.String(prepayId)
		resp.SignType = utils.String("RSA")
		resp.Appid = utils.String(*queryResp.Appid)
		resp.TimeStamp = utils.String(strconv.FormatInt(time.Now().Unix(), 10))
		nonce, err := utils.GenerateNonce()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("generate request for payment err:%s", err.Error())
		}
		resp.NonceStr = utils.String(nonce)
		resp.Package = utils.String("prepay_id=" + prepayId)
		message := fmt.Sprintf("%s\n%s\n%s\n%s\n", *resp.Appid, *resp.TimeStamp, *resp.NonceStr, *resp.Package)
		signatureResult, err := jsapi.Client.Sign(ctx, message)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("generate sign for payment err:%s", err.Error())
		}
		resp.PaySign = utils.String(signatureResult.Signature)
	}

	return transaction, resp, result, nil
}

func (a *PayMentService) CloseOrder(outTradeNo string) (err error) {
	ctx := context.Background()
	opts := []client.ClientOption{
		option.WithWechatPayAutoAuthCipher(),
	}
	defaultClient, err := client.NewClient(opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}
	jsapi := JsapiApiService{Client: defaultClient}
	req := payRequest.CloseOrderRequest{}
	req.OutTradeNo = utils.String(outTradeNo)
	req.Mchid = utils.String(consts.MachID)
	_, err = jsapi.CloseOrder(ctx, req)
	if err != nil {
		return err
	}
	return err
}
