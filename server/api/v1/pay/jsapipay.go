package pay

import (
	"cooller/server/client/consts"
	"cooller/server/client/notify"
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	payModel "cooller/server/model/pay"
	payRequest "cooller/server/model/pay/request"
	payRes "cooller/server/model/pay/response"
	"cooller/server/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

type PayApi struct{}

//func (e *PayApi) GenerateOrder(c *gin.Context) {
//	var orderReq wechatReq.OrderCreateRequest
//	err := c.ShouldBindJSON(&orderReq)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	if len(orderReq.Ids) < 1 {
//		global.GVA_LOG.Error("下单ids不可为空!", zap.Error(err))
//		response.FailWithMessage("下单ids不可为空", c)
//		return
//	}
//	var order product.Order
//	var payReq payRequest.PrepayRequest
//	payReq.Appid = utils.String(orderReq.AppId)
//	payReq.Mchid = utils.String(consts.MachID)
//	orderDescription := ""
//	var goodsDetail []payRequest.GoodsDetail
//
//	order.TotalAmount = 0
//	order.PayAmount = 0
//	order.PromotionAmount = 0
//
//	userId := utils.GetUserID(c)
//	var productCartList []product.CartCommonItem
//	if orderReq.BuyType == 1 { // 直接购买
//		productCartList, err = orderService.GetProductTmpCartByIds(userId, orderReq.Ids)
//		if err != nil || len(productCartList) < 1 {
//			global.GVA_LOG.Error("获取购物车物品失败!", zap.Error(err))
//			response.FailWithMessage("获取购物车物品失败", c)
//			return
//		}
//
//	} else if orderReq.BuyType == 2 { // 加入购物车
//		productCartList, err = orderService.GetProductCartByIds(userId, orderReq.Ids)
//		if err != nil || len(productCartList) < 1 {
//			global.GVA_LOG.Error("获取购物车物品失败!", zap.Error(err))
//			response.FailWithMessage("获取购物车物品失败", c)
//			return
//		}
//	}
//
//	var productApi product2.OrderApi
//	if !productApi.HasCommonStock(productCartList) {
//		global.GVA_LOG.Error("获取购物车物品失败!", zap.Error(err))
//		response.FailWithMessage("获取购物车物品失败", c)
//		return
//	}
//
//	for _, cartItem := range productCartList {
//		_, promotionMessage, reduceAmount := wechatApi.CalculateProductPromotionPrice(cartItem.Product, nil)
//		var orderItem product.OrderItem
//		order.PromotionAmount += reduceAmount
//		order.PromotionInfo = promotionMessage
//		orderItem.PromotionAmount = reduceAmount
//		orderItem.PromotionName = promotionMessage
//		// 计算优惠前总金额
//		order.TotalAmount += order.TotalAmount + cartItem.SkuStock.PromotionPrice*float32(cartItem.Quantity)
//		// 该商品经过优惠后的实际金额
//		realAmount := cartItem.SkuStock.PromotionPrice*float32(cartItem.Quantity) - reduceAmount
//		if realAmount < 0 {
//			global.GVA_LOG.Error("[GenerateOrder]获取价格计算失败!", zap.Error(err))
//			response.FailWithMessage("[GenerateOrder]获取价格计算失败", c)
//			return
//		}
//		//orderItem.OrderId = order.ID
//		orderItem.ProductId = cartItem.ProductId
//		orderItem.ProductSkuId = cartItem.SkuStock.SkuCode
//		orderItem.UserId = cartItem.UserId
//		orderItem.Quantity = cartItem.Quantity
//		orderItem.Price = cartItem.SkuStock.PromotionPrice
//		orderItem.ProductPic = cartItem.Product.Pic
//		orderItem.ProductName = cartItem.Product.Name
//		orderItem.ProductSubTitle = cartItem.Product.SubTitle
//		orderItem.ProductSkuCode = ""
//		orderItem.MemberNickname = utils.GetUserName(c)
//		orderItem.DeleteStatus = 0
//		orderItem.ProductCategoryId = cartItem.Product.ProductCategoryId
//		orderItem.ProductBrand = cartItem.Product.BrandName
//		orderItem.ProductSN = cartItem.Product.ProductSN
//		orderItem.ProductAttr = cartItem.SkuStock.SpData
//		orderItem.CouponAmount = 0
//		orderItem.IntegrationAmount = 0
//		orderItem.RealAmount = realAmount
//		orderItem.GiftIntegration = 0
//		orderItem.GiftGrowth = 0
//		order.OrderItemList = append(order.OrderItemList, &orderItem)
//
//		orderDescription = fmt.Sprintf("%s x%d ", cartItem.Product.Name, cartItem.Quantity)
//		var goodDetail payRequest.GoodsDetail
//		goodDetail.MerchantGoodsId = utils.String(cartItem.Product.ProductSN)
//		goodDetail.GoodsName = utils.String(cartItem.Product.Name)
//		goodDetail.Quantity = utils.Int64(int64(cartItem.Quantity))
//		goodDetail.UnitPrice = utils.Int64(int64(cartItem.SkuStock.PromotionPrice * 100))
//		goodsDetail = append(goodsDetail, goodDetail)
//	}
//
//	n, err := snowflake.NewNode(1)
//	if err != nil {
//		global.GVA_LOG.Error("创建id失败!", zap.Error(err))
//	}
//	order.UserId = userId
//	order.CouponId = orderReq.CouponId
//	order.OrderSn = fmt.Sprintf("%d", n.Generate())
//	userName := utils.GetNickName(c)
//	if len(userName) < 1 {
//		userName = utils.GetTelephone(c)
//	}
//	order.UserName = userName
//	order.FreightAmount = 0
//	order.PayAmount += order.TotalAmount - order.PromotionAmount
//	order.IntegrationAmount = float32(orderReq.UseIntegration / 1000) // 1000积分抵1元
//	order.CouponAmount = 0
//	order.DiscountAmount = 0
//	order.PayType = orderReq.PayType
//	order.SourceType = 1
//	order.Status = 0
//	order.OrderType = 0
//	//order.LogisticsCompany = ""
//	//order.LogisticsSn = ""
//	order.AutoConfirmDay = 7
//	order.Integration = 0
//	order.Growth = 0
//	address, err := accountService.GetMemberReceiveAddressById(orderReq.MemberReceiveAddressId)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	order.ReceiverPhone = address.Telephone
//	order.ReceiverName = address.Name
//	order.ReceiverPostCode = address.PostCode
//	order.ReceiverProvince = address.Province
//	order.ReceiverCity = address.City
//	order.ReceiverRegion = address.Region
//	order.ReceiverDetailAddress = address.DetailAddress
//	order.Note = orderReq.Note
//	order.ConfirmStatus = 0
//	order.DeleteStatus = 0
//	order.UseIntegration = orderReq.UseIntegration
//	order.PaymentTime = time.Now()
//
//	order.LogisticsTime = time.Now()
//	order.ReceiveTime = time.Now()
//	order.CommentTime = time.Now()
//	order.ModifyTime = time.Now()
//	err = orderService.CreateOrder(&order)
//	if err != nil {
//		global.GVA_LOG.Error("创建订单数据失败!", zap.Error(err))
//		response.FailWithMessage("创建订单数据失败", c)
//		return
//	}
//
//	payReq.Description = utils.String(orderDescription)
//	payReq.OutTradeNo = utils.String(order.OrderSn)
//	payReq.TimeExpire = utils.String(time.Now().Format(time.RFC3339))
//	payReq.NotifyUrl = utils.String("https://cs.coollerbaby.cn/pay/notify")
//	payReq.GoodsTag = utils.String("WXG")
//	payReq.SettleInfo = &payRequest.SettleInfo{
//		ProfitSharing: utils.Bool(false),
//	}
//	payReq.SupportFapiao = utils.Bool(false)
//	payReq.Amount = &payRequest.Amount{
//		Currency: utils.String("CNY"),
//		Total:    utils.Int64(int64(order.PayAmount * 100)),
//	}
//	payReq.Payer = &payRequest.Payer{
//		Openid: utils.String(orderReq.OpenId),
//	}
//	payReq.Detail = &payRequest.Detail{
//		CostPrice:   utils.Int64(int64(order.TotalAmount * 100)),
//		GoodsDetail: goodsDetail,
//		//InvoiceId:   utils.String("wx123"),
//	}
//	payReq.SceneInfo = &payRequest.SceneInfo{
//		DeviceId:      utils.String("013467007045764"),
//		PayerClientIp: utils.String(orderReq.IP),
//		StoreInfo: &payRequest.StoreInfo{
//			Address:  utils.String("河南省商丘市中骏雍景台27栋113商铺"),
//			AreaCode: utils.String("476100"),
//			Id:       utils.String("0001"),
//			Name:     utils.String("猪迪克星动乐园"),
//		},
//	}
//	res, _, err := jspaymentService.PrepayWithRequestPayment(payReq)
//	if err != nil {
//		global.GVA_LOG.Error("更新失败!", zap.Error(err))
//		fmt.Println("支付失败!", zap.Error(err))
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//
//	//err = orderService.UpdateOrderStatusById(order.ID, 0)
//	//if err != nil {
//	//	global.GVA_LOG.Error("更新订单支付状态失败!", zap.Error(err))
//	//	response.FailWithMessage(err.Error(), c)
//	//	return
//	//}
//	err = orderService.UpdateOrderPrepayId(order.ID, *res.PrepayId)
//	if err != nil {
//		global.GVA_LOG.Error("更新订单预支付交易会话标识失败!", zap.Error(err))
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	var data payRes.GenerateOrderResponse
//	data.OrderId = order.ID
//	data.Payment = &res
//	response.OkWithData(data, c)
//}
//
//// HasStock 判断下单商品是否都有库存
//func (e *PayApi) HasStock(cartItemList []product.CartTmpItem) bool {
//	for _, cartItem := range cartItemList {
//		if cartItem.SkuStock.Stock <= 0 || cartItem.SkuStock.Stock < cartItem.Quantity {
//			return false
//		}
//	}
//	return true
//}
//
//// 进行库存锁定
//func (e *PayApi) LockStock(cartItemList []product.CartItem) error {
//	for _, cartPromotionItem := range cartItemList {
//		count, err := productService.UpdateProductSkuStockForStock(cartPromotionItem.SkuStockId, cartPromotionItem.ProductId)
//		if err != nil {
//			return fmt.Errorf("修改库存时失败: %v", err)
//		}
//		if count == 0 {
//			return fmt.Errorf("库存不足, 无法下单")
//		}
//	}
//	return nil
//}

func (e *PayApi) GetOrderDetail(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order, err := orderService.GetProductOrderById(reqId.ID)
	if err != nil {
		response.FailWithMessage("获取订单数据失败", c)
		return
	}

	var queryReq payRequest.QueryOrderByOutTradeNoRequest
	queryReq.Mchid = utils.String(consts.MachID)
	queryReq.OutTradeNo = utils.String(order.OrderSn)
	_, res, _, err := jspaymentService.QueryOrderByOutTradeNo(queryReq, order.PrepayId)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	var data payRes.GenerateOrderDetailResponse
	data.Order = order
	data.Payment = res
	response.OkWithData(data, c)
}

// CancelOrder 取消订单 付费版
//func (e *PayApi) CancelOrder(c *gin.Context) {
//	var reqIds request.IdsReq
//	err := c.ShouldBindJSON(&reqIds)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	userId := utils.GetUserID(c)
//	if userId <= 0 {
//		response.FailWithMessage("Not get userId!", c)
//		return
//	}
//	orderList, err := orderService.DeleteManyOrder(reqIds.Ids)
//	if err != nil {
//		global.GVA_LOG.Error("删除订单数据失败!", zap.Error(err))
//		response.FailWithMessage("删除订单数据失败", c)
//		return
//	}
//	for _, order := range orderList {
//		err = jspaymentService.CloseOrder(order.OrderSn)
//		if err != nil {
//			global.GVA_LOG.Error("删除订单数据失败!", zap.Error(err))
//			response.FailWithMessage("删除订单数据失败", c)
//			return
//		}
//	}
//
//	response.OkWithMessage("删除成功", c)
//}

func (e *PayApi) OrderNotify(c *gin.Context) {
	request := notify.Request{}
	c.ShouldBind(&request)
	mapstructure.Decode(c.Params, &request)
	if request.EventType == "TRANSACTION.SUCCESS" {
		//plaintext, err := wepay.DecryptAES256GCM(
		//	aesKey, request.Resource.AssociatedData, request.Resource.Nonce, request.Resource.Ciphertext,
		//)
		plaintext, err := utils.DecryptAES256GCM(
			consts.MchAPIv3Key, request.Resource.AssociatedData, request.Resource.Nonce, request.Resource.Ciphertext,
		)
		if err != nil {
			fmt.Println(err)
			zap.S().Error("DecryptAES256GCM err" + err.Error())
		}
		transaction := payModel.Transaction{}
		json.Unmarshal([]byte(plaintext), &transaction)
		go func() {
			// 执行service层代码
			err = orderService.UpdateOrderStatusByOrderSn(transaction.OutTradeNo, 1)
			if err != nil {
				global.GVA_LOG.Error("支付回调更新订单状态失败!", zap.Error(err))
			}
		}()
		tmp := make(map[string]interface{})
		tmp["code"] = "SUCCESS"
		tmp["message"] = "成功"
		tmpJson, _ := json.Marshal(tmp)
		c.Writer.Write(tmpJson)
	} else {
		tmp := make(map[string]interface{})
		tmp["code"] = "500"
		tmp["message"] = "失败"
		tmpJson, _ := json.Marshal(tmp)
		c.Writer.Write(tmpJson)
	}
}
