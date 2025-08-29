package product

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	"cooller/server/model/wechat"
	request2 "cooller/server/model/wechat/request"
	"cooller/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type FlashApi struct{}

func (e *FlashApi) GetFlashPromotionList(c *gin.Context) {
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
	homeList, total, err := wechatService.GetFlashPromotionList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     homeList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (e *FlashApi) CreateFlashPromotion(c *gin.Context) {
	var flash wechat.FlashPromotion

	err := c.ShouldBindJSON(&flash)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = wechatService.CreateFlashPromotion(&flash)
	if err != nil {
		global.GVA_LOG.Error("创建产品失败!", zap.Error(err))
		response.FailWithMessage("创建产品失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (e *FlashApi) UpdateFlashPromotion(c *gin.Context) {
	var flash wechat.FlashPromotion

	err := c.ShouldBindJSON(&flash)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = wechatService.UpdateFlashPromotion(&flash)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)

}

func (e *FlashApi) DeleteFlashPromotionById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = wechatService.DeleteFlashPromotionById(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (e *FlashApi) UpdateFlashPromotionStatus(c *gin.Context) {
	var statusInfo request.StatusUpdateInfo
	err := c.ShouldBindJSON(&statusInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = wechatService.UpdateFlashPromotionStatus(statusInfo.ID, statusInfo.Status)
	if err != nil {
		global.GVA_LOG.Error("更新状态失败!", zap.Error(err))
		response.FailWithMessage("更新状态失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新状态成功", c)

}

func (e *FlashApi) GetFlashPromotionProductRelationList(c *gin.Context) {
	var pageInfo request2.FlashProductRelationInfo
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
	homeList, total, err := wechatService.GetFlashPromotionProductRelationList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     homeList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

type Flashes struct {
	Flashes []wechat.FlashPromotionProductRelation `json:"flashes" form:"flashes"`
}

func (e *FlashApi) CreateFlashPromotionProductRelation(c *gin.Context) {
	var flashes Flashes
	err := c.ShouldBindJSON(&flashes)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = wechatService.CreateFlashPromotionProductRelation(flashes.Flashes)
	if err != nil {
		global.GVA_LOG.Error("创建产品失败!", zap.Error(err))
		response.FailWithMessage("创建产品失败", c)
		return
	}

	for _, falshRelation := range flashes.Flashes {
		flashSession, err := wechatService.GetFlashSessionById(falshRelation.FlashPromotionSessionId)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		err = wechatService.UpdateProductPromotionType(&falshRelation, &flashSession)
		if err != nil {
			global.GVA_LOG.Error("创建产品失败!", zap.Error(err))
			response.FailWithMessage("创建产品失败", c)
			return
		}
	}

	response.OkWithMessage("创建成功", c)
}

func (e *FlashApi) UpdateFlashPromotionProductRelation(c *gin.Context) {
	var flash wechat.FlashPromotionProductRelation

	err := c.ShouldBindJSON(&flash)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = wechatService.UpdateFlashPromotionProductRelation(&flash)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)

}

func (e *FlashApi) DeleteFlashPromotionProductRelationById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = wechatService.DeleteFlashPromotionProductRelationById(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (e *FlashApi) GetFlashSessionList(c *gin.Context) {
	homeList, err := wechatService.GetFlashSessionList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.AllResult{
		List: homeList,
	}, "获取成功", c)
}

func (e *FlashApi) CreateFlashSession(c *gin.Context) {
	var flash wechat.FlashPromotionSession

	err := c.ShouldBindJSON(&flash)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	flash.CreateDate = time.Now()
	err = wechatService.CreateFlashSession(&flash)
	if err != nil {
		global.GVA_LOG.Error("创建产品失败!", zap.Error(err))
		response.FailWithMessage("创建产品失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (e *FlashApi) UpdateFlashSession(c *gin.Context) {
	var flash wechat.FlashPromotionSession

	err := c.ShouldBindJSON(&flash)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = wechatService.UpdateFlashSession(&flash)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)

}

func (e *FlashApi) DeleteFlashSessionById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = wechatService.DeleteFlashSessionById(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (e *FlashApi) UpdateFlashSessionStatus(c *gin.Context) {
	var statusInfo request.StatusUpdateInfo
	err := c.ShouldBindJSON(&statusInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = wechatService.UpdateFlashSessionStatus(statusInfo.ID, statusInfo.Status)
	if err != nil {
		global.GVA_LOG.Error("更新状态失败!", zap.Error(err))
		response.FailWithMessage("更新状态失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新状态成功", c)
}

func (e *FlashApi) GetFlashSessionSelectList(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	homeList, err := wechatService.GetFlashSessionSelectList(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	for _, session := range homeList {
		productCount := len(session.ProductRelation)
		session.ProductCount = productCount
	}
	response.OkWithDetailed(response.AllResult{
		List: homeList,
	}, "获取成功", c)
}
