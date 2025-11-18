package product

import (
	"context"
	wechatApi "cooller/server/api/v1/wechat"
	"cooller/server/client/consts"
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	payRequest "cooller/server/model/pay/request"
	payRes "cooller/server/model/pay/response"
	"cooller/server/model/product"
	productReq "cooller/server/model/product/request"
	productRes "cooller/server/model/product/response"
	"cooller/server/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ChangSZ/golib/mathutil"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderApi struct{}

// GenerateConfirmOrder 生成确认单信息
func (e *OrderApi) GenerateConfirmOrder(c *gin.Context) {
	var reqIds request.IdsTagReq
	err := c.ShouldBindJSON(&reqIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	if userId <= 0 {
		response.FailWithMessage("Not get userId!", c)
		return
	}

	// 获取购物车信息
	var productCartList []*product.CartCommonItem
	if reqIds.Tag == 1 {
		productCartList, err = orderService.GetProductTmpCartByIds(userId, reqIds.Ids)
		if err != nil {
			global.GVA_LOG.Error("获取购物车物品失败!", zap.Error(err))
			response.FailWithMessage("获取购物车物品失败", c)
			return
		}

	} else if reqIds.Tag == 2 {
		productCartList, err = orderService.GetProductCartByIds(userId, reqIds.Ids)
		if err != nil {
			global.GVA_LOG.Error("获取购物车物品失败!", zap.Error(err))
			response.FailWithMessage("获取购物车物品失败", c)
			return
		}
	}

	cartPromotionItemList, err := e.ListPromotion(productCartList)
	// 实际付款金额
	//payAmount := totalAmount - promotionAmount
	//if payAmount < 0 {
	//	global.GVA_LOG.Error("付款金额错误!", zap.Error(err))
	//	response.FailWithMessage("付款金额错误", c)
	//	return
	//}
	// 获取用户可用优惠券列表
	api := CouponApi{}
	couponHistoryDetailList, err := api.ListCart(userId, cartPromotionItemList, 1)
	if err != nil {
		global.GVA_LOG.Error("获取优惠券失败!", zap.Error(err))
		response.FailWithMessage("获取优惠券失败", c)
		return
	}
	// 获取用户收货地址列表
	addressList, err := accountService.GetMemberReceiveAddressList(userId)
	if err != nil {
		global.GVA_LOG.Error("获取收货地址列表失败!", zap.Error(err))
		response.FailWithMessage("获取收货地址列表失败", c)
		return
	}
	// 获取用户积分

	// 计算总金额、活动优惠、应付金额
	calcAmount := e.CalcCartAmount(cartPromotionItemList)

	var order productRes.GenerateOrderResModel
	order.CartPromotionItemList = cartPromotionItemList
	order.MemberReceiveAddressList = addressList
	order.CalcAmount = *calcAmount
	order.PickupType = 1
	order.CouponHistoryDetailList = couponHistoryDetailList
	response.OkWithData(order, c)
}

// CalcCartAmount 计算购物车中商品的价格
func (e *OrderApi) CalcCartAmount(cartPromotionItemList []*product.OrderItem) *productRes.CalcAmount {
	calcAmount := &productRes.CalcAmount{}
	var totalAmount float32
	var promotionAmount float32
	for _, cartPromotionItem := range cartPromotionItemList {
		totalAmount += cartPromotionItem.Price * float32(cartPromotionItem.Quantity)
		promotionAmount += cartPromotionItem.ReduceAmount * float32(cartPromotionItem.Quantity)
	}
	calcAmount.TotalAmount = totalAmount
	calcAmount.PromotionAmount = promotionAmount
	calcAmount.PayAmount = totalAmount - promotionAmount
	return calcAmount
}

// GenerateOrder
func (e *OrderApi) GenerateOrder(c *gin.Context) {
	var orderReq productReq.OrderCreateRequest
	err := c.ShouldBindJSON(&orderReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(orderReq.Ids) < 1 {
		global.GVA_LOG.Error("下单ids不可为空!", zap.Error(err))
		response.FailWithMessage("下单ids不可为空", c)
		return
	}
	// 校验收货地址
	if orderReq.MemberReceiveAddressId == 0 {
		global.GVA_LOG.Error("请选择收货地址！", zap.Error(err))
		response.FailWithMessage("请选择收货地址！", c)
		return
	}

	userId := utils.GetUserID(c)

	orderDescription := ""
	var goodsDetail []payRequest.GoodsDetail
	// 1. 获取当前要购买的商品列表（含数量），用于后续重复校验
	var productCommonCartList []*product.CartCommonItem
	if orderReq.BuyType == 1 { // 直接购买
		productCommonCartList, err = orderService.GetProductTmpCartByIds(userId, orderReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("获取购物车物品失败", zap.Error(err))
			response.FailWithMessage("获取购物车物品失败", c)
			return
		}

	} else if orderReq.BuyType == 2 {
		productCommonCartList, err = orderService.GetProductCartByIds(userId, orderReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("获取购物车物品失败", zap.Error(err))
			response.FailWithMessage("获取购物车物品失败", c)
			return
		}
	}
	// 2. 生成当前订单的商品特征（商品ID+数量，用于校验重复）
	currentProducts := make([]productKey, 0, len(productCommonCartList))
	for _, item := range productCommonCartList {
		currentProducts = append(currentProducts, productKey{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	// 3. 检查是否存在重复的未支付订单
	// 调用订单服务查询符合条件的未支付订单
	duplicateOrder, err := e.CheckDuplicateUnpaidOrder(userId, orderReq, currentProducts)
	if err != nil {
		global.GVA_LOG.Error("检查重复订单失败", zap.Error(err))
		response.FailWithMessage("系统错误", c)
		return
	}
	if duplicateOrder != nil {
		// 存在重复未支付订单，直接返回该订单的支付信息（无需创建新订单）
		var data payRes.GenerateOrderResponse
		data.OrderId = duplicateOrder.ID
		// 若原订单预支付信息已过期，可重新生成预支付单
		if duplicateOrder.PrepayId == "" || e.isPrepayExpired(duplicateOrder.PaymentTime) {
			// 重新生成预支付信息
			res, _, err := jspaymentService.PrepayWithRequestPayment(e.buildPayReq(duplicateOrder, orderReq, "商品订单支付"))
			if err != nil {
				global.GVA_LOG.Error("更新失败!", zap.Error(err))
				response.FailWithMessage(err.Error(), c)
				return
			}
			duplicateOrder.PrepayId = *res.PrepayId
			err = orderService.UpdateOrderPaySuccess(duplicateOrder.ID, *res.PrepayId)
			if err != nil {
				global.GVA_LOG.Error("更新订单预支付交易会话标识失败!", zap.Error(err))
				response.FailWithMessage(err.Error(), c)
				return
			}
			data.Payment = &res
		} else {
			// 预支付信息未过期，直接返回
			nonce, err := utils.GenerateNonce()
			if err != nil {
				global.GVA_LOG.Error("更新订单预支付交易会话标识失败!", zap.Error(err))
				response.FailWithMessage(err.Error(), c)
				return
			}

			data.Payment = &payRes.PrepayWithRequestPaymentResponse{
				PrepayId:  &duplicateOrder.PrepayId,
				Appid:     &orderReq.AppId,
				TimeStamp: utils.String(strconv.FormatInt(time.Now().Unix(), 10)),
				NonceStr:  utils.String(nonce),
				Package:   utils.String("prepay_id=" + duplicateOrder.PrepayId),
				SignType:  utils.String("RSA"),
			}
		}
		response.OkWithData(data, c)
		return
	}

	if !e.HasCommonStock(productCommonCartList) {
		global.GVA_LOG.Error("库存不足，无法下单", zap.Error(err))
		response.FailWithMessage("库存不足，无法下单", c)
		return
	}

	cartPromotionItemList, err := e.ListPromotion(productCommonCartList)
	if err != nil {
		global.GVA_LOG.Error("获取促销商品失败", zap.Error(err))
		response.FailWithMessage("获取促销商品失败", c)
		return
	}
	//orderItemList := make([]*product.CartPromotionItem, 0)

	totalItems := len(cartPromotionItemList)
	// 限制最多显示3个商品，超过则用“等N件”概括
	maxShow := 3
	for i, cartItem := range cartPromotionItemList {
		//
		if i >= maxShow {
			// 超出最大显示数量，补充剩余数量
			remaining := totalItems - maxShow
			orderDescription += fmt.Sprintf(" 等%d件商品", remaining)
			break
		}
		// 拼接商品信息（用“+”分隔多个商品）
		itemDesc := fmt.Sprintf("%s x%d", cartItem.ProductName, cartItem.Quantity)
		if i == 0 {
			orderDescription = itemDesc
		} else {
			orderDescription += " + " + itemDesc
		}
		var goodDetail payRequest.GoodsDetail
		goodDetail.MerchantGoodsId = utils.String(cartItem.ProductSN)
		goodDetail.GoodsName = utils.String(cartItem.ProductName)
		goodDetail.Quantity = utils.Int64(int64(cartItem.Quantity))
		goodDetail.UnitPrice = utils.Int64(int64(cartItem.Price * 100))
		goodsDetail = append(goodsDetail, goodDetail)
	}
	if orderDescription == "" {
		orderDescription = "商品订单支付"
	}

	//进行库存锁定
	//if err := e.LockCommonStock(productCommonCartList); err != nil {
	//	global.GVA_LOG.Error("锁定库存失败!", zap.Error(err))
	//	response.FailWithMessage("锁定库存失败", c)
	//	return
	//}

	// TODO: 优惠券计算
	// 比较注册日期和购买时间，计算注册赠券是否足够
	if orderReq.CouponId != 0 {
		couponHistoryDetail, err := e.GetUseCoupon(userId, cartPromotionItemList, orderReq.CouponId)
		if err != nil || couponHistoryDetail == nil {
			global.GVA_LOG.Error("获取优惠券失败!", zap.Error(err))
			response.FailWithMessage("获取优惠券失败", c)
			return
		}
		// 对下单商品的优惠券进行处理
		e.HandleCouponAmount(cartPromotionItemList, *couponHistoryDetail)
	}

	// 计算order_item的实付金额
	e.HandleRealAmount(cartPromotionItemList)
	// 根据商品合计、运费、活动优惠、优惠券、积分计算应付金额
	order := product.Order{
		TotalAmount:     e.CalcTotalAmount(cartPromotionItemList),
		PromotionAmount: e.CalcPromotionAmount(cartPromotionItemList),
		PromotionInfo:   e.GetOrderPromotionInfo(cartPromotionItemList),
	}

	if orderReq.CouponId != 0 {
		order.CouponId = orderReq.CouponId
		order.CouponAmount = e.CalcCouponAmount(cartPromotionItemList)
	}

	// 总金额+运费-促销优惠-优惠券优惠-积分抵扣
	order.DiscountAmount = 0
	order.FreightAmount = 0                                           // 运费
	order.IntegrationAmount = float32(orderReq.UseIntegration / 1000) // 1000积分抵1元
	order.PayAmount = order.TotalAmount + order.FreightAmount -
		order.PromotionAmount - order.CouponAmount - order.IntegrationAmount
	if order.PayAmount < 0 {
		global.GVA_LOG.Error("[GenerateOrder]获取价格计算错误!", zap.Error(err))
		response.FailWithMessage("获取价格计算错误", c)
		return
	}

	order.UserId = userId
	order.CouponId = orderReq.CouponId
	order.OrderSn = e.GenerateOrderSn(order)
	userName := utils.GetUserName(c)
	if len(userName) < 1 {
		userName = utils.GetTelephone(c)
	}
	order.UserName = userName
	// 支付方式：0->未支付；1->支付宝；2->微信
	order.PayType = orderReq.PayType
	// 订单来源：0->PC订单；1->app订单
	order.SourceType = 1
	// 订单状态：0->待付款；1->待发货；2->已发货；3->已完成；4->已关闭；5->无效订单
	order.Status = 0
	// 订单类型：0->正常订单；1->秒杀订单
	order.OrderType = 0

	//order.LogisticsCompany = ""
	//order.LogisticsSn = ""
	order.AutoConfirmDay = 7
	// 计算赠送积分
	order.Integration = e.CalcGifIntegration(cartPromotionItemList)
	order.Growth = 100
	//order.BillType = 0
	//order.BillHeader = ""
	//order.BillContent = ""
	//order.BillReceiverPhone = ""
	//order.BillReceiverEmail = ""

	// 收货人信息：姓名、电话、邮编、地址
	address, err := accountService.GetMemberReceiveAddressById(orderReq.MemberReceiveAddressId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order.MemberReceiveAddressId = address.ID
	//order.ReceiverPhone = address.Telephone
	//order.ReceiverName = address.Name
	//order.ReceiverPostCode = address.PostCode
	//order.ReceiverProvince = address.Province
	//order.ReceiverCity = address.City
	//order.ReceiverRegion = address.Region
	//order.ReceiverDetailAddress = address.DetailAddress

	order.Note = orderReq.Note
	// 0->未确认；1->已确认
	order.ConfirmStatus = 0
	order.DeleteStatus = 0

	order.UseIntegration = orderReq.UseIntegration
	order.PaymentTime = time.Now()

	order.LogisticsTime = time.Now()
	order.ReceiveTime = time.Now()
	//order.CommentTime = time.Now()
	//order.ModifyTime = time.Now()

	//var orderData wechat.Order
	//copy.AssignStruct(&order, orderData)
	order.OrderItemList = cartPromotionItemList
	for _, orderItem := range order.OrderItemList {
		//orderItem.OrderId = order.ID
		orderItem.OrderSn = e.GenerateOrderSn(order) //TODO: 唯一性有待优化
	}
	err = orderService.CreateOrder(&order)
	if err != nil {
		global.GVA_LOG.Error("创建订单数据失败!", zap.Error(err))
		response.FailWithMessage("创建订单数据失败", c)
		return
	}

	// 如使用优惠券c
	if order.CouponId != 0 {
		if err := couponService.UpdateCouponStatus(orderReq.CouponId, userId, 1); err != nil {
			global.GVA_LOG.Error("更新优惠券状态失败!", zap.Error(err))
			response.FailWithMessage("更新优惠券状态失败", c)
			return
		}
	}

	// 如使用积分需要扣除积分
	if order.UseIntegration != 0 {
		//	TODO: 更新积分
	}

	//删除购物车中的下单商品
	if orderReq.BuyType == 2 {
		err = productService.DeleteProductCartByIds(userId, orderReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("删除购物车失败!", zap.Error(err))
			response.FailWithMessage("删除购物车失败", c)
			return
		}
	}

	payReq := e.buildPayReq(&order, orderReq, orderDescription)
	res, _, err := jspaymentService.PrepayWithRequestPayment(payReq)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		fmt.Println("支付失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = orderService.UpdateOrderPaySuccess(order.ID, *res.PrepayId)
	if err != nil {
		global.GVA_LOG.Error("更新订单预支付交易会话标识失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	var data payRes.GenerateOrderResponse
	data.OrderId = order.ID
	data.Payment = &res
	response.OkWithData(data, c)
}

// 辅助函数：判断预支付信息是否过期（例如30分钟过期）
func (e *OrderApi) isPrepayExpired(paymentTime time.Time) bool {
	return time.Since(paymentTime) > 30*time.Minute
}

// 辅助函数：构建支付请求参数
func (e *OrderApi) buildPayReq(order *product.Order, orderReq productReq.OrderCreateRequest, orderDescription string) payRequest.PrepayRequest {
	var payReq payRequest.PrepayRequest
	payReq.Appid = utils.String(orderReq.AppId)
	payReq.Mchid = utils.String(consts.MachID)
	payReq.Description = utils.String(orderDescription)
	//OrderSn := e.GenerateOrderSn(*order)
	//global.GVA_LOG.Error("---OrderSn:" + OrderSn)
	payReq.OutTradeNo = utils.String(order.OrderSn)
	payReq.TimeExpire = utils.String(time.Now().Format(time.RFC3339))
	payReq.NotifyUrl = utils.String("https://wx.guocihuatai.com/order/notify")
	payReq.GoodsTag = utils.String("WXG")
	payReq.SettleInfo = &payRequest.SettleInfo{
		ProfitSharing: utils.Bool(false),
	}
	payReq.SupportFapiao = utils.Bool(false)
	payReq.Amount = &payRequest.Amount{
		Currency: utils.String("CNY"),
		Total:    utils.Int64(int64(order.PayAmount * 100)),
	}
	payReq.Payer = &payRequest.Payer{
		Openid: utils.String(orderReq.OpenId),
	}
	//payReq.Detail = &payRequest.Detail{
	//	CostPrice:   utils.Int64(int64(order.TotalAmount * 100)),
	//	GoodsDetail: goodsDetail,
	//	//InvoiceId:   utils.String("wx123"),
	//}
	//payReq.SceneInfo = &payRequest.SceneInfo{
	//	DeviceId:      utils.String("013467007045764"),
	//	PayerClientIp: utils.String(orderReq.IP),
	//	StoreInfo: &payRequest.StoreInfo{
	//		Address:  utils.String("浙江省绍兴市诸暨市大唐街道"),
	//		AreaCode: utils.String("476100"),
	//		Id:       utils.String("0001"),
	//		Name:     utils.String("怡驰针织"),
	//	},
	//}
	return payReq
}

// HasCommonStock 判断下单商品是否都有库存
func (e *OrderApi) HasCommonStock(cartItemList []*product.CartCommonItem) bool {
	for _, cartItem := range cartItemList {
		if cartItem.SkuStock.Stock <= 0 || cartItem.SkuStock.Stock < cartItem.Quantity {
			return false
		}
	}
	return true
}

// 进行库存锁定
func (e *OrderApi) LockCommonStock(cartItemList []*product.CartCommonItem) error {
	for _, cartPromotionItem := range cartItemList {
		count, err := productService.UpdateProductSkuStockForStock(cartPromotionItem.SkuStockId, cartPromotionItem.Quantity)
		if err != nil {
			return fmt.Errorf("修改库存时失败: %v", err)
		}
		if count == 0 {
			return fmt.Errorf("库存不足, 无法下单")
		}
	}
	return nil
}

func (e *OrderApi) GetUseCoupon(userId int, cartItemList []*product.OrderItem, couponId int) (couponHistory *product.CouponHistory, err error) {
	var couponApi CouponApi
	couponHistoryDetailList, err := couponApi.ListCart(userId, cartItemList, 1)
	if err != nil {
		return nil, err
	}
	for _, couponHistoryDetail := range couponHistoryDetailList {
		if couponHistoryDetail.CouponId == couponId {
			return couponHistoryDetail, nil
		}
	}
	return nil, nil
}

// ListPromotion
func (e *OrderApi) ListPromotion(cartItemList []*product.CartCommonItem) (cartPromotionItemList []*product.OrderItem, err error) {
	for _, cart := range cartItemList {
		var cartPromotionItem product.OrderItem
		cartPromotionItem.ID = cart.ID
		cartPromotionItem.CreatedAt = cart.CreatedAt
		cartPromotionItem.Quantity = cart.Quantity
		cartPromotionItem.ProductPic = cart.SkuStock.Pic
		cartPromotionItem.ProductName = cart.Product.Name
		cartPromotionItem.ProductId = cart.ProductId
		cartPromotionItem.ProductCategoryId = cart.Product.ProductCategoryId

		cartPromotionItem.Price = cart.Product.Price

		cartPromotionItem.ProductSkuCode = cart.SkuStock.SkuCode
		cartPromotionItem.MemberNickname = ""
		cartPromotionItem.ProductBrand = cart.Product.BrandName
		cartPromotionItem.ProductSN = cart.Product.ProductSN
		cartPromotionItem.ProductAttr = cart.SkuStock.SpData
		cartPromotionItem.CouponAmount = 0

		cartPromotionItem.GiftIntegration = 0
		cartPromotionItem.GiftGrowth = 0

		var homeApi wechatApi.HomeApi
		promotionMessage, reduceAmount := homeApi.CalculateProductPromotionPrice(cart.Product, nil)
		cartPromotionItem.ReduceAmount = reduceAmount
		cartPromotionItem.PromotionMessage = fmt.Sprintf("满减优惠：%s", promotionMessage)
		cartPromotionItem.RealStock = cart.SkuStock.Stock
		cartPromotionItem.IntegrationAmount = 0
		// 该商品经过优惠后的实际金额
		realAmount := cart.Product.Price*float32(cart.Quantity) - reduceAmount
		cartPromotionItem.RealAmount = realAmount
		cartPromotionItemList = append(cartPromotionItemList, &cartPromotionItem)
	}

	return cartPromotionItemList, nil
}

func (e *OrderApi) HandleCouponAmount(orderItemList []*product.OrderItem, couponHistoryDetail product.CouponHistory) {
	coupon := couponHistoryDetail.Coupon
	switch coupon.UseType {
	case 0: // 全场通用
		e.CalcPerCouponAmount(orderItemList, coupon)
	case 1: // 指定分类
		couponOrderItemList := e.GetCouponOrderItemByRelation(couponHistoryDetail, orderItemList, 0)
		e.CalcPerCouponAmount(couponOrderItemList, coupon)
	case 2: // 指定商品
		couponOrderItemList := e.GetCouponOrderItemByRelation(couponHistoryDetail, orderItemList, 1)
		e.CalcPerCouponAmount(couponOrderItemList, coupon)
	}
}

/**
 * 获取与优惠券有关系的下单商品
 *
 * @param couponHistoryDetail 优惠券详情
 * @param orderItemList       下单商品
 * @param refType             使用关系类型：0->相关分类；1->指定商品
 */
func (e *OrderApi) GetCouponOrderItemByRelation(couponHistoryDetail product.CouponHistory,
	orderItemList []*product.OrderItem, refType int32) []*product.OrderItem {
	result := make([]*product.OrderItem, 0)
	switch refType {
	case 0:
		categoryIdsMap := make(map[int]bool, 0)
		for _, productCategoryRelation := range couponHistoryDetail.Coupon.ProductCategoryRelationList {
			categoryIdsMap[productCategoryRelation.ProductCategoryId] = true
		}
		for i, orderItem := range orderItemList {
			if _, ok := categoryIdsMap[orderItem.ProductCategoryId]; ok {
				result = append(result, orderItem)
			} else {
				orderItemList[i].CouponAmount = 0
			}
		}
	case 1:
		productIdsMap := make(map[int]bool, 0)
		for _, productRelation := range couponHistoryDetail.Coupon.ProductRelationList {
			productIdsMap[productRelation.ProductId] = true
		}
		for i, orderItem := range orderItemList {
			if _, ok := productIdsMap[orderItem.ProductId]; ok {
				result = append(result, orderItem)
			} else {
				orderItemList[i].CouponAmount = 0
			}
		}
	}
	return result
}

/**
 * 对每个下单商品进行优惠券金额分摊的计算
 *
 * @param orderItemList 可用优惠券的下单商品商品
 */
func (e *OrderApi) CalcPerCouponAmount(orderItemList []*product.OrderItem, coupon product.Coupon) {
	totalAmount := e.CalcTotalAmount(orderItemList)
	for i, orderItem := range orderItemList {
		// (商品价格/可用商品总价)*优惠券面额
		couponAmount := mathutil.RoundHalfEven(float64(orderItem.Price/totalAmount)*float64(coupon.Amount), 3)
		orderItemList[i].CouponAmount = float32(couponAmount)
	}
}

// CalcTotalAmount 计算总金额
func (e *OrderApi) CalcTotalAmount(orderItemList []*product.OrderItem) float32 {
	var totalAmount float32
	for _, item := range orderItemList {
		totalAmount += item.Price * float32(item.Quantity)
	}
	return totalAmount
}

func (e *OrderApi) HandleRealAmount(orderItemList []*product.OrderItem) {
	for i, orderItem := range orderItemList {
		// 原价-促销优惠-优惠券抵扣-积分抵扣
		realAmount := orderItem.Price - orderItem.PromotionAmount -
			orderItem.CouponAmount - orderItem.IntegrationAmount
		orderItemList[i].RealAmount = realAmount
	}
}

// CalcPromotionAmount 计算订单活动优惠
func (e *OrderApi) CalcPromotionAmount(orderItemList []*product.OrderItem) float32 {
	var promotionAmount float32
	for _, orderItem := range orderItemList {
		if orderItem.PromotionAmount != 0 {
			promotionAmount += orderItem.PromotionAmount * float32(orderItem.Quantity)
		}
	}
	return promotionAmount
}

// GetOrderPromotionInfo 获取订单促销信息
func (e *OrderApi) GetOrderPromotionInfo(orderItemList []*product.OrderItem) string {
	promotionNameList := make([]string, 0, len(orderItemList))
	for _, orderItem := range orderItemList {
		promotionNameList = append(promotionNameList, orderItem.PromotionName)
	}
	return strings.Join(promotionNameList, ";")
}

// CalcCouponAmount 计算订单优惠券金额
func (e *OrderApi) CalcCouponAmount(orderItemList []*product.OrderItem) float32 {
	var couponAmount float32
	for _, orderItem := range orderItemList {
		if orderItem.CouponAmount != 0 {
			couponAmount += orderItem.CouponAmount * float32(orderItem.Quantity)
		}
	}
	return couponAmount
}

// CalcGifIntegration 计算该订单赠送的积分
func (e *OrderApi) CalcGifIntegration(orderItemList []*product.OrderItem) int {
	var sum int
	for _, orderItem := range orderItemList {
		sum += orderItem.GiftIntegration * orderItem.Quantity
	}
	return sum
}

// GenerateOrderSn 生成18位订单编号:8位日期+2位平台号码+2位支付方式+6位以上自增id
func (e *OrderApi) GenerateOrderSn(order product.Order) string {
	ctx := context.Background()
	var sb strings.Builder
	date := time.Now().Format("20060102")
	key := global.GVA_CONFIG.Wechat.AppID + ":" + global.GVA_CONFIG.Wechat.Secret + ":" + date
	increment := global.GVA_REDIS.Incr(ctx, key)
	if increment.Err() != nil {
		// 处理Redis操作错误（如连接失败等）
		global.GVA_LOG.Error("Redis Incr失败: %v", zap.Error(increment.Err()))
		return ""
	}
	sb.WriteString(date)
	sb.WriteString(fmt.Sprintf("%02d", order.SourceType))
	sb.WriteString(fmt.Sprintf("%02d", order.PayType))
	incrementStr := fmt.Sprintf("%06d", increment.Val())
	fmt.Println(incrementStr)
	// 限制自增ID最大20字节，超过则取后20位（保证总长度≤32）
	if len(incrementStr) > 6 {
		if len(incrementStr) > 20 {
			incrementStr = incrementStr[len(incrementStr)-20:] // 保留后20位
		} else {
			sb.WriteString(incrementStr)
		}
	} else {
		sb.WriteString(incrementStr[:6])
	}
	return sb.String()
}

type productKey struct {
	ProductId int
	Quantity  int
}

// CheckDuplicateUnpaidOrder 检查是否存在重复的未支付订单
func (e *OrderApi) CheckDuplicateUnpaidOrder(userId int, orderReq productReq.OrderCreateRequest, currentProducts []productKey) (*product.Order, error) {
	// 1. 先查询用户的未支付订单
	var unpaidOrders []*product.Order
	db := global.GVA_DB.Where("user_id = ? AND status = 0", userId)
	if err := db.Find(&unpaidOrders).Error; err != nil {
		return nil, err
	}

	// 2. 遍历未支付订单，检查是否与当前订单特征匹配
	for _, order := range unpaidOrders {
		// 基础条件匹配：收货地址、优惠券、支付方式、积分
		if order.MemberReceiveAddressId != orderReq.MemberReceiveAddressId ||
			order.CouponId != orderReq.CouponId ||
			order.PayType != orderReq.PayType ||
			order.UseIntegration != orderReq.UseIntegration {
			continue
		}

		// 3. 检查商品及数量是否完全匹配
		var orderItems []*product.OrderItem
		if err := global.GVA_DB.Where("order_sn = ?", order.OrderSn).Find(&orderItems).Error; err != nil {
			return nil, err
		}
		if len(orderItems) != len(currentProducts) {
			continue // 商品数量不同，不是重复订单
		}

		// 构建订单商品的ID-数量映射
		orderProductMap := make(map[int]int)
		for _, item := range orderItems {
			orderProductMap[item.ProductId] = item.Quantity
		}

		// 校验当前商品与订单商品是否完全一致
		match := true
		for _, p := range currentProducts {
			if orderProductMap[p.ProductId] != p.Quantity {
				match = false
				break
			}
		}
		if match {
			return order, nil // 找到重复订单
		}
	}

	return nil, nil // 无重复订单
}

func (e *OrderApi) GetOrderDetail(c *gin.Context) {
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

func (e *OrderApi) GetOrderList(c *gin.Context) {
	var stateInfo productReq.SearchInfo
	err := c.ShouldBindQuery(&stateInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(stateInfo, utils.StateInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	orderList, total, err := orderService.GetProductOrderListByStatus(stateInfo)
	if err != nil {
		global.GVA_LOG.Error("获取订单数据失败!", zap.Error(err))
		response.FailWithMessage("获取订单数据失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     orderList,
		Total:    total,
		Page:     stateInfo.Page,
		PageSize: stateInfo.PageSize,
	}, "获取成功", c)
}

func (e *OrderApi) GetOrderItemList(c *gin.Context) {
	var stateInfo productReq.SearchInfo
	err := c.ShouldBindQuery(&stateInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(stateInfo, utils.StateInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	orderList, total, err := orderService.GetProductOrderItemListByStatus(stateInfo)
	if err != nil {
		global.GVA_LOG.Error("获取订单数据失败!", zap.Error(err))
		response.FailWithMessage("获取订单数据失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     orderList,
		Total:    total,
		Page:     stateInfo.Page,
		PageSize: stateInfo.PageSize,
	}, "获取成功", c)
}

func (e *OrderApi) PaySuccess(c *gin.Context) {
	var paySuccess productReq.PaySuccessRequest
	err := c.ShouldBindJSON(&paySuccess)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//删除购物车中的下单商品
	//orderIdList := make([]int, 0)
	//orderIdList = append(order IdList, order.ID)
	//_, err = orderService.DeleteManyOrder(orderIdList)
	//if err != nil {
	//	global.GVA_LOG.Error("更新订单数据失败!", zap.Error(err))
	//	response.FailWithMessage("更新订单数据失败", c)
	//	return
	//}
	err = orderService.UpdateOrderStatus(&paySuccess, 1)
	if err != nil {
		global.GVA_LOG.Error("更新订单数据失败!", zap.Error(err))
		response.FailWithMessage("更新订单数据失败", c)
		return
	}
	response.OkWithMessage("支付成功", c)
}

// CancelOrders 取消订单 不付费版
func (e *OrderApi) CancelOrders(c *gin.Context) {
	var reqIds request.IdsReq
	err := c.ShouldBindJSON(&reqIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	if userId <= 0 {
		response.FailWithMessage("Not get userId!", c)
		return
	}
	err = orderService.UpdateOrdersStatus(reqIds.Ids, 5)
	if err != nil {
		global.GVA_LOG.Error("更新订单数据失败!", zap.Error(err))
		response.FailWithMessage("更新订单数据失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// CloseOrders 关闭订单
func (e *OrderApi) CloseOrders(c *gin.Context) {
	var reqIds request.IdsReq
	err := c.ShouldBindJSON(&reqIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	if userId <= 0 {
		response.FailWithMessage("Not get userId!", c)
		return
	}

	err = orderService.UpdateManyOrderStatus(reqIds.Ids, 4)
	if err != nil {
		global.GVA_LOG.Error("更新订单数据失败!", zap.Error(err))
		response.FailWithMessage("更新订单数据失败", c)
		return
	}

	response.OkWithMessage("关闭成功", c)
}

// DeleteOrders 删除订单
func (e *OrderApi) DeleteOrders(c *gin.Context) {
	var reqIds request.IdsReq
	err := c.ShouldBindJSON(&reqIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	if userId <= 0 {
		response.FailWithMessage("Not get userId!", c)
		return
	}
	_, err = orderService.DeleteManyOrder(reqIds.Ids)
	if err != nil {
		global.GVA_LOG.Error("更新订单数据失败!", zap.Error(err))
		response.FailWithMessage("更新订单数据失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (e *OrderApi) UpdateOrderReceiverInfo(c *gin.Context) {
	var address productReq.OrderReceiveAddress
	err := c.ShouldBindJSON(&address)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = orderService.UpdateOrderReceiverInfo(&address)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *OrderApi) UpdateOrderMoneyInfo(c *gin.Context) {
	var info productReq.OrderMoneyInfo
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = orderService.UpdateOrderMoneyInfo(&info)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *OrderApi) UpdateOrderNote(c *gin.Context) {
	var info productReq.OrderNoteInfo
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = orderService.UpdateOrderNoteInfo(&info)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *OrderApi) UpdateOrderCompletedStatus(c *gin.Context) {
	var reqIds request.IdsReq
	err := c.ShouldBindJSON(&reqIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = orderService.UpdateOrdersStatus(reqIds.Ids, 3)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *OrderApi) GetOrderSetting(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	setting, err := orderService.GetOrderSetting(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("更新订单数据失败!", zap.Error(err))
		response.FailWithMessage("更新订单数据失败", c)
		return
	}

	response.OkWithData(setting, c)
}

func (e *OrderApi) UpdateOrderSetting(c *gin.Context) {
	var home product.OrderSetting
	err := c.ShouldBindJSON(&home)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = orderService.UpdateOrderSetting(&home)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *OrderApi) GetProductCartList(c *gin.Context) {
	cartList, err := productService.GetProductCartList()
	if err != nil {
		global.GVA_LOG.Error("获取商品sku库存失败!", zap.Error(err))
		response.FailWithMessage("获取sku库存失败", c)
		return
	}
	response.OkWithData(cartList, c)
}

// CreateProductCart 创建商品购物车
func (e *OrderApi) CreateProductCart(c *gin.Context) {
	var cart product.CartItem
	err := c.ShouldBindJSON(&cart)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cart.UserId = utils.GetUserID(c)
	err = productService.CreateProductCart(&cart)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateProductCartQuantity 更新商品购物车数量
func (e *OrderApi) UpdateProductCartQuantity(c *gin.Context) {
	var quantityInfo request.QuantityInfo
	err := c.ShouldBindQuery(&quantityInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = productService.UpdateProductCartQuantity(userId, quantityInfo.ID, quantityInfo.Quantity)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteProductCartById 删除商品购物车
func (e *OrderApi) DeleteProductCartById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = productService.DeleteProductCartById(userId, reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (e *OrderApi) DeleteProductCartByIds(c *gin.Context) {
	var reqIds request.IdsReq
	err := c.ShouldBindQuery(&reqIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = productService.DeleteProductCartByIds(userId, reqIds.Ids)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// ClearProductCart 清空商品购物车
func (e *OrderApi) ClearProductCart(c *gin.Context) {
	userId := utils.GetUserID(c)
	err := productService.ClearProductCartUserId(userId)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// CreateProductTmpCart 创建商品购物车
func (e *OrderApi) CreateProductTmpCart(c *gin.Context) {
	var cart product.CartTmpItem
	err := c.ShouldBindJSON(&cart)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cart.UserId = utils.GetUserID(c)
	err = productService.CreateProductTmpCart(&cart)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	data := make(map[string]int)
	data["id"] = cart.ID
	response.OkWithData(data, c)
}

func (e *OrderApi) UpdateOrderLogistics(c *gin.Context) {
	var info productReq.UpdateLogisticsRequest
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(info.LogisticsInfos) < 1 {
		response.FailWithMessage("不可为空", c)
		return
	}
	var ids []int
	for _, item := range info.LogisticsInfos {
		if item.ID <= 0 {
			response.FailWithMessage("ID错误", c)
			return
		}
		ids = append(ids, item.ID)
	}

	err = orderService.UpdateOrderLogistics(&info, ids)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *OrderApi) UploadFileWx(c *gin.Context) {
	fmt.Println("---OssType:", global.GVA_CONFIG.System.OssType)
	var file product.AfterSalesUpload
	noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	userId := utils.GetUserID(c)

	file, err = orderService.UploadFile(header, noSave, userId) // 文件上传后拿到文件路径
	if err != nil {
		global.GVA_LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	response.OkWithDetailed(response.ItemId{Id: file.ID}, "上传成功", c)
}

// DealOrder Images:存储的是图片id的字符串
func (e *OrderApi) DealOrder(c *gin.Context) {
	var info productReq.OrderDealRequest
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	//imagesJson, err := json.Marshal(info.Images)
	//if err != nil {
	//	panic(err)
	//}
	//imagesStr := string(imagesJson)
	orderItem, err := orderService.GetProductOrderItemById(info.OrderItemId)
	if err != nil {
		global.GVA_LOG.Error("获取OrderItem失败!", zap.Error(err))
		response.FailWithMessage("获取OrderItem失败!", c)
		return
	}
	orderApply := product.AfterSalesApply{
		UserId:      userId,
		OrderItemId: info.OrderItemId,
		Content:     info.Content,
		Contact:     info.Contact,
		RealAmount:  orderItem.RealAmount,
		Images:      info.Images,
		Status:      0,
	}
	err = orderService.CreateDealOrderApply(&orderApply)
	if err != nil {
		global.GVA_LOG.Error("创建订单售后申请失败!", zap.Error(err))
		response.FailWithMessage("创建订单售后申请失败", c)
		return
	}
	err = orderService.UpdateOrderAfterSalesStatusById(info.OrderItemId, 1)
	if err != nil {
		global.GVA_LOG.Error("修改订单是否售后状态失败!", zap.Error(err))
		response.FailWithMessage("修改订单是否售后状态失败", c)
		return
	}
	response.OkWithMessage("创建订单售后申请成功", c)
}

type imageIdsList struct {
	ImageIds []int `json:"imageIds"`
}

func (e *OrderApi) GetDealOrderList(c *gin.Context) {
	var pageInfo productReq.OrderDealSearchRequest
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
	list, total, err := orderService.GetDealOrderList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	for _, item := range list {
		var imageIds imageIdsList
		err = json.Unmarshal([]byte(item.Images), &imageIds.ImageIds)
		if err != nil {
			global.GVA_LOG.Error("图片字符串组合json转换失败", zap.Error(err))
			response.FailWithMessage("图片字符串组合json转换失败"+err.Error(), c)
		}
		if len(imageIds.ImageIds) > 0 {
			images, err := orderService.GetOrderDealUploadImages(imageIds.ImageIds)
			if err != nil {
				global.GVA_LOG.Error("获取售后图片失败", zap.Error(err))
				response.FailWithMessage("获取售后图片失败"+err.Error(), c)
			}
			for _, image := range images {
				item.ImagesList = append(item.ImagesList, image.Url)
			}
		}
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (e *OrderApi) UpdateDealOrder(c *gin.Context) {
	var info productReq.UpdateDealOrderRequest
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if info.DealId < 1 || info.OrderItemId < 1 {
		global.GVA_LOG.Error("更新售后申请参数失败!", zap.Error(err))
		response.FailWithMessage("更新更新售后申请参数失败", c)
	}
	err = orderService.UpdateDealOrderSynchronous(info)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
