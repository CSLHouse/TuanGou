package product

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	"cooller/server/model/product"
	couponRes "cooller/server/model/product/request"
	"cooller/server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type CouponApi struct{}

func (e *CouponApi) GetCouponList(c *gin.Context) {
	var pageInfo request.PageInfo
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
	list, total, err := couponService.GetCouponList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (e *CouponApi) GetCouponListWithState(c *gin.Context) {
	var searchInfo couponRes.SearchByState
	err := c.ShouldBindQuery(&searchInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	list, err := couponService.GetCouponHistoryListByState(userId, searchInfo.UseStatus)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.AllResult{
		List: list,
	}, "获取成功", c)
}

func (e *CouponApi) GetCouponById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	coupon, err := couponService.GetCouponById(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(coupon, "获取成功", c)
}

func (e *CouponApi) CreateCoupon(c *gin.Context) {
	var coupon product.Coupon

	err := c.ShouldBindJSON(&coupon)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = couponService.CreateCoupon(&coupon)
	if err != nil {
		global.GVA_LOG.Error("创建产品失败!", zap.Error(err))
		response.FailWithMessage("创建产品失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (e *CouponApi) UpdateCoupon(c *gin.Context) {
	var coupon product.Coupon
	err := c.ShouldBindJSON(&coupon)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = couponService.UpdateCouponSynchronous(&coupon)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)

}

func (e *CouponApi) DeleteCouponById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = couponService.DeleteCoupon(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (e *CouponApi) GetCouponHistoryList(c *gin.Context) {
	var pageInfo couponRes.SearchInfoCoupon
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
	list, total, err := couponService.GetCouponHistoryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// TODO: 修复many2many获取CouponList失败
func (e *CouponApi) GetCouponListByProduct(c *gin.Context) {
	var productId request.GetById
	err := c.ShouldBindQuery(&productId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 获取指定商品优惠券
	productData, err := productService.GetProductCouponByID(productId.ID)
	if err != nil {
		global.GVA_LOG.Error("获取优惠券失败!", zap.Error(err))
		response.FailWithMessage("获取优惠券失败", c)
		return
	}
	couponList := productData.CouponList

	// 获取指定分类优惠券
	productCategoryData, err := productService.GetProductCategoryById(productData.ProductCategoryId)
	if err != nil {
		global.GVA_LOG.Error("获取优惠券失败!", zap.Error(err))
		response.FailWithMessage("获取优惠券失败", c)
		return
	}
	couponList = append(couponList, productCategoryData.CouponList...)

	response.OkWithDetailed(couponList, "", c)
}

// AddUserCoupon
func (e *CouponApi) AddUserCoupon(c *gin.Context) {
	var couponId request.GetById
	err := c.ShouldBindQuery(&couponId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	coupon, err := couponService.GetCouponById(couponId.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if coupon.ID == 0 {
		response.FailWithMessage("优惠券不存在", c)
		return
	}
	if coupon.Count <= 0 {
		response.FailWithMessage("优惠券已经领完了", c)
		return
	}
	if time.Now().Before(coupon.EnableTime) {
		response.FailWithMessage("优惠券还没到领取时间", c)
		return
	}

	// 判断用户领取的优惠券数量是否超过限制
	userId := utils.GetUserID(c)
	count, err := couponService.GetCouponHistoryCount(couponId.ID, userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if count >= int64(coupon.PerLimit) {
		response.WarningWithMessage("您已经领取过该优惠券", c)
		return
	}
	//currMember, err := accountService.GetWXAccountByOpenID(utils.Get)
	// 生成领取优惠券历史
	var couponHistory product.CouponHistory
	couponHistory.CouponId = couponId.ID
	couponHistory.CouponCode = e.generateCouponCode(int64(userId))
	couponHistory.MemberId = userId
	couponHistory.MemberNickname = utils.GetUserName(c)
	couponHistory.GetType = 1
	couponHistory.UseStatus = 0
	couponHistory.UseTime = nil
	err = couponService.CreateCouponHistory(&couponHistory)
	if err != nil {
		global.GVA_LOG.Error("创建产品失败!", zap.Error(err))
		response.FailWithMessage("创建产品失败", c)
		return
	}
	// 更新优惠券记录
	err = couponService.UpdateCouponCount(couponId.ID)
	if err != nil {
		global.GVA_LOG.Error("领取优惠券失败!", zap.Error(err))
		response.FailWithMessage("领取优惠券失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// generateCouponCode 16位优惠码生成：时间戳后8位+4位随机数+用户id后4位
func (e *CouponApi) generateCouponCode(memberId int64) string {
	var sb string
	currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)
	timeMillisStr := fmt.Sprintf("%d", currentTimeMillis)
	sb += timeMillisStr[len(timeMillisStr)-8:]

	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 4; i++ {
		sb += fmt.Sprintf("%d", rand.Intn(10))
	}

	memberIdStr := fmt.Sprintf("%d", memberId)
	if len(memberIdStr) <= 4 {
		sb += fmt.Sprintf("%04d", memberId)
	} else {
		sb += memberIdStr[len(memberIdStr)-4:]
	}
	return sb
}

func (e *CouponApi) ListCart(userId int, cartItemList []*product.CartPromotionItem, enable int) (couponHistorylList []*product.CouponHistory, err error) {
	// 获取该用户所有优惠券
	allList, err := couponService.GetUserCouponHistoryById(userId)
	if err != nil {
		global.GVA_LOG.Error("获取所有优惠券失败!", zap.Error(err))
		return couponHistorylList, err
	}
	// 根据优惠券使用类型来判断优惠券是否可用
	now := time.Now()
	enableList := make([]*product.CouponHistory, 0)
	disableList := make([]*product.CouponHistory, 0)
	for _, v := range allList {
		useType := v.Coupon.UseType
		minPoint := v.Coupon.MinPoint
		endTime := v.Coupon.EndTime
		switch useType {
		case 0: // 全场通用
			// 判断是否满足优惠起点
			// 计算购物车商品的总价
			totalAmount := e.calcTotalAmount(cartItemList)
			if now.Before(endTime) && totalAmount-minPoint >= 0 {
				enableList = append(enableList, v)
			} else {
				disableList = append(disableList, v)
			}
		case 1: // 指定分类
			// 计算指定分类商品的总价
			productCategoryIds := make([]int, 0)
			for _, item := range v.Coupon.ProductCategoryRelationList {
				productCategoryIds = append(productCategoryIds, item.ProductCategoryId)
			}
			totalAmount := e.calcTotalAmountByproductCategoryId(cartItemList, productCategoryIds)
			if now.Before(endTime) && totalAmount > 0 && totalAmount-minPoint >= 0 {
				enableList = append(enableList, v)
			} else {
				disableList = append(disableList, v)
			}
		case 2: // 指定商品
			// 计算指定商品的总价
			productIds := make([]int, 0)
			for _, item := range v.Coupon.ProductRelationList {
				productIds = append(productIds, item.ProductId)
			}
			totalAmount := e.calcTotalAmountByProductId(cartItemList, productIds)
			if now.Before(endTime) && totalAmount > 0 && totalAmount-minPoint >= 0 {
				enableList = append(enableList, v)
			} else {
				disableList = append(disableList, v)
			}
		}
	}
	if enable == 1 {
		return enableList, nil
	}
	return disableList, err
}

func (e *CouponApi) calcTotalAmount(cartItemList []*product.CartPromotionItem) float32 {
	total := float32(0)
	for _, item := range cartItemList {
		realPrice := item.Price - item.ReduceAmount
		total += realPrice * float32(item.Quantity)
	}
	return total
}

func (e *CouponApi) calcTotalAmountByproductCategoryId(cartItemList []*product.CartPromotionItem, productCategoryIds []int) float32 {
	cateIdsMap := make(map[int]bool)
	for _, v := range productCategoryIds {
		cateIdsMap[v] = true
	}

	total := float32(0.0)
	for _, item := range cartItemList {
		if _, ok := cateIdsMap[item.ProductCategoryId]; !ok {
			continue
		}
		realPrice := item.Price - item.ReduceAmount
		total += realPrice * float32(item.Quantity)
	}
	return total
}

func (e *CouponApi) calcTotalAmountByProductId(cartItemList []*product.CartPromotionItem, productIds []int) float32 {
	prodIdsMap := make(map[int]bool)
	for _, v := range productIds {
		prodIdsMap[v] = true
	}

	total := float32(0.0)
	for _, item := range cartItemList {
		if _, ok := prodIdsMap[item.ProductId]; !ok {
			continue
		}
		realPrice := item.Price - item.ReduceAmount
		total += realPrice * float32(item.Quantity)
	}
	return total
}
