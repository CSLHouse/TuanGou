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
	wechatReq "cooller/server/model/wechat/request"
	wechatRes "cooller/server/model/wechat/response"
	"cooller/server/utils"
	"fmt"
	"github.com/ChangSZ/golib/mathutil"
	"github.com/ChangSZ/golib/repo/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"time"
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

	var order wechatRes.GenerateOrderResModel
	order.CartPromotionItemList = cartPromotionItemList
	order.MemberReceiveAddressList = addressList
	order.CalcAmount = *calcAmount
	order.PickupType = 1
	order.CouponHistoryDetailList = couponHistoryDetailList
	response.OkWithData(order, c)
}

// CalcCartAmount 计算购物车中商品的价格
func (e *OrderApi) CalcCartAmount(cartPromotionItemList []*product.CartPromotionItem) *wechatRes.CalcAmount {
	calcAmount := &wechatRes.CalcAmount{}
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

func (e *OrderApi) GenerateOrder(c *gin.Context) {
	var orderReq wechatReq.OrderCreateRequest
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
	orderItemList := make([]*product.CartPromotionItem, 0)

	for _, cartItem := range cartPromotionItemList {
		var orderItem product.CartPromotionItem

		//orderItem.OrderId = order.ID
		orderItem.PromotionName = cartItem.PromotionName
		orderItem.ProductId = cartItem.ProductId
		orderItem.ProductSkuId = cartItem.ProductSkuId
		orderItem.UserId = cartItem.UserId
		orderItem.Quantity = cartItem.Quantity
		orderItem.Price = cartItem.Price
		orderItem.ProductPic = cartItem.ProductPic
		orderItem.ProductName = cartItem.ProductName
		orderItem.ProductSubTitle = cartItem.ProductSubTitle
		orderItem.ProductSkuCode = ""
		orderItem.MemberNickname = utils.GetUserName(c)
		orderItem.DeleteStatus = 0
		orderItem.ProductCategoryId = cartItem.ProductCategoryId
		orderItem.ProductBrand = cartItem.ProductBrand
		orderItem.ProductSN = cartItem.ProductSN
		orderItem.ProductAttr = cartItem.ProductAttr
		orderItem.CouponAmount = 0
		orderItem.IntegrationAmount = 0
		orderItem.RealAmount = cartItem.RealAmount
		orderItem.GiftIntegration = 0
		orderItem.GiftGrowth = 0
		orderItemList = append(orderItemList, &orderItem)

		//
		orderDescription = fmt.Sprintf("%s x%d ", cartItem.ProductName, cartItem.Quantity)
		var goodDetail payRequest.GoodsDetail
		goodDetail.MerchantGoodsId = utils.String(cartItem.ProductSN)
		goodDetail.GoodsName = utils.String(cartItem.ProductName)
		goodDetail.Quantity = utils.Int64(int64(cartItem.Quantity))
		goodDetail.UnitPrice = utils.Int64(int64(cartItem.Price * 100))
		goodsDetail = append(goodsDetail, goodDetail)
	}
	//var order product.Order
	//order.TotalAmount = 0
	//order.PayAmount = 0
	//order.PromotionAmount = 0
	//order.OrderItemList = orderItemList

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
		e.HandleCouponAmount(orderItemList, *couponHistoryDetail)
	}

	// 计算order_item的实付金额
	e.HandleRealAmount(orderItemList)
	// 根据商品合计、运费、活动优惠、优惠券、积分计算应付金额
	order := product.Order{
		TotalAmount:     e.CalcTotalAmount(orderItemList),
		PromotionAmount: e.CalcPromotionAmount(orderItemList),
		PromotionInfo:   e.GetOrderPromotionInfo(orderItemList),
	}

	if orderReq.CouponId != 0 {
		order.CouponId = orderReq.CouponId
		order.CouponAmount = e.CalcCouponAmount(orderItemList)
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

	//order.DeliveryCompany = ""
	//order.DeliverySn = ""
	order.AutoConfirmDay = 7
	// 计算赠送积分
	order.Integration = e.CalcGifIntegration(orderItemList)
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
	order.ReceiverPhone = address.Telephone
	order.ReceiverName = address.Name
	order.ReceiverPostCode = address.PostCode
	order.ReceiverProvince = address.Province
	order.ReceiverCity = address.City
	order.ReceiverRegion = address.Region
	order.ReceiverDetailAddress = address.DetailAddress

	order.Note = orderReq.Note
	// 0->未确认；1->已确认
	order.ConfirmStatus = 0
	order.DeleteStatus = 0

	order.UseIntegration = orderReq.UseIntegration
	order.PaymentTime = time.Now()

	order.DeliveryTime = time.Now()
	order.ReceiveTime = time.Now()
	order.CommentTime = time.Now()
	order.ModifyTime = time.Now()

	//var orderData wechat.Order
	//copy.AssignStruct(&order, orderData)
	for _, orderItem := range order.OrderItemList {
		orderItem.OrderId = order.ID
		orderItem.OrderSn = order.OrderSn
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
	//orderIdList := make([]int, 0)
	//orderIdList = append(orderIdList, order.ID)
	//_, err = orderService.DeleteManyOrder(orderIdList)
	//if err != nil {
	//	global.GVA_LOG.Error("更新订单数据失败!", zap.Error(err))
	//	response.FailWithMessage("更新订单数据失败", c)
	//	return
	//}

	var payReq payRequest.PrepayRequest
	payReq.Appid = utils.String(orderReq.AppId)
	payReq.Mchid = utils.String(consts.MachID)
	payReq.Description = utils.String(orderDescription)
	payReq.OutTradeNo = utils.String(order.OrderSn)
	payReq.TimeExpire = utils.String(time.Now().Format(time.RFC3339))
	payReq.NotifyUrl = utils.String("https://cs.coollerbaby.cn/pay/notify")
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
	payReq.Detail = &payRequest.Detail{
		CostPrice:   utils.Int64(int64(order.TotalAmount * 100)),
		GoodsDetail: goodsDetail,
		//InvoiceId:   utils.String("wx123"),
	}
	payReq.SceneInfo = &payRequest.SceneInfo{
		DeviceId:      utils.String("013467007045764"),
		PayerClientIp: utils.String(orderReq.IP),
		StoreInfo: &payRequest.StoreInfo{
			Address:  utils.String("浙江省绍兴市诸暨市大唐街道"),
			AreaCode: utils.String("476100"),
			Id:       utils.String("0001"),
			Name:     utils.String("怡驰针织"),
		},
	}
	res, _, err := jspaymentService.PrepayWithRequestPayment(payReq)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		fmt.Println("支付失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = orderService.UpdateOrderPrepayId(order.ID, *res.PrepayId)
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

func (e *OrderApi) GetUseCoupon(userId int, cartItemList []*product.CartPromotionItem, couponId int) (couponHistory *product.CouponHistory, err error) {
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
func (e *OrderApi) ListPromotion(cartItemList []*product.CartCommonItem) (cartPromotionItemList []*product.CartPromotionItem, err error) {
	for _, cart := range cartItemList {
		var cartPromotionItem product.CartPromotionItem
		cartPromotionItem.ID = cart.ID
		cartPromotionItem.CreatedAt = cart.CreatedAt
		cartPromotionItem.Quantity = cart.Quantity
		cartPromotionItem.ProductPic = cart.SkuStock.Pic
		cartPromotionItem.ProductName = cart.Product.Name
		cartPromotionItem.ProductId = cart.ProductId
		cartPromotionItem.ProductCategoryId = cart.Product.ProductCategoryId
		promotionMessage, reduceAmount := wechatApi.CalculateProductPromotionPrice(cart.Product, nil)
		cartPromotionItem.ReduceAmount = reduceAmount
		cartPromotionItem.Price = cart.Product.Price
		cartPromotionItem.PromotionMessage = fmt.Sprintf("满减优惠：%s", promotionMessage)
		cartPromotionItem.RealStock = cart.SkuStock.Stock
		// 该商品经过优惠后的实际金额
		realAmount := cart.Product.Price*float32(cart.Quantity) - reduceAmount
		cartPromotionItem.RealAmount = realAmount
		cartPromotionItemList = append(cartPromotionItemList, &cartPromotionItem)
	}

	return cartPromotionItemList, nil
}

func (e *OrderApi) HandleCouponAmount(orderItemList []*product.CartPromotionItem, couponHistoryDetail product.CouponHistory) {
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
	orderItemList []*product.CartPromotionItem, refType int32) []*product.CartPromotionItem {
	result := make([]*product.CartPromotionItem, 0)
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
func (e *OrderApi) CalcPerCouponAmount(orderItemList []*product.CartPromotionItem, coupon product.Coupon) {
	totalAmount := e.CalcTotalAmount(orderItemList)
	for i, orderItem := range orderItemList {
		// (商品价格/可用商品总价)*优惠券面额
		couponAmount := mathutil.RoundHalfEven(float64(orderItem.Price/totalAmount)*float64(coupon.Amount), 3)
		orderItemList[i].CouponAmount = float32(couponAmount)
	}
}

// CalcTotalAmount 计算总金额
func (e *OrderApi) CalcTotalAmount(orderItemList []*product.CartPromotionItem) float32 {
	var totalAmount float32
	for _, item := range orderItemList {
		totalAmount += item.Price * float32(item.Quantity)
	}
	return totalAmount
}

func (e *OrderApi) HandleRealAmount(orderItemList []*product.CartPromotionItem) {
	for i, orderItem := range orderItemList {
		// 原价-促销优惠-优惠券抵扣-积分抵扣
		realAmount := orderItem.Price - orderItem.PromotionAmount -
			orderItem.CouponAmount - orderItem.IntegrationAmount
		orderItemList[i].RealAmount = realAmount
	}
}

// CalcPromotionAmount 计算订单活动优惠
func (e *OrderApi) CalcPromotionAmount(orderItemList []*product.CartPromotionItem) float32 {
	var promotionAmount float32
	for _, orderItem := range orderItemList {
		if orderItem.PromotionAmount != 0 {
			promotionAmount += orderItem.PromotionAmount * float32(orderItem.Quantity)
		}
	}
	return promotionAmount
}

// GetOrderPromotionInfo 获取订单促销信息
func (e *OrderApi) GetOrderPromotionInfo(orderItemList []*product.CartPromotionItem) string {
	promotionNameList := make([]string, 0, len(orderItemList))
	for _, orderItem := range orderItemList {
		promotionNameList = append(promotionNameList, orderItem.PromotionName)
	}
	return strings.Join(promotionNameList, ";")
}

// CalcCouponAmount 计算订单优惠券金额
func (e *OrderApi) CalcCouponAmount(orderItemList []*product.CartPromotionItem) float32 {
	var couponAmount float32
	for _, orderItem := range orderItemList {
		if orderItem.CouponAmount != 0 {
			couponAmount += orderItem.CouponAmount * float32(orderItem.Quantity)
		}
	}
	return couponAmount
}

// CalcGifIntegration 计算该订单赠送的积分
func (e *OrderApi) CalcGifIntegration(orderItemList []*product.CartPromotionItem) int {
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
	increment := redis.Cache().Incr(ctx, key)
	sb.WriteString(date)
	sb.WriteString(fmt.Sprintf("%02d", order.SourceType))
	sb.WriteString(fmt.Sprintf("%02d", order.PayType))
	incrementStr := fmt.Sprintf("%06d", increment)
	if len(incrementStr) > 6 {
		sb.WriteString(incrementStr)
	} else {
		sb.WriteString(incrementStr[:6])
	}
	return sb.String()
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
	var stateInfo request.StateInfo
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

func (e *OrderApi) PaySuccess(c *gin.Context) {
	var paySuccess wechatReq.PaySuccessRequest
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
	var address wechatReq.OrderReceiveAddress
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
	var info wechatReq.OrderMoneyInfo
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
	var info wechatReq.OrderNoteInfo
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
