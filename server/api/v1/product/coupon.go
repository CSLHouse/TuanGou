package product

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	"cooller/server/model/product"
	couponRes "cooller/server/model/product/request"
	"cooller/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	var flash product.Coupon
	err := c.ShouldBindJSON(&flash)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = couponService.UpdateCoupon(&flash)
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

func (e *CouponApi) GetCouponListByProduct(c *gin.Context) {
	var productId request.GetById
	err := c.ShouldBindQuery(&productId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 获取指定商品优惠券

	// 获取指定分类优惠券
	err = couponService.DeleteCoupon(productId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
