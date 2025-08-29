package business

import (
	"cooller/server/global"
	"cooller/server/model/business"
	businessRes "cooller/server/model/business/response"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	"cooller/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const ( // 须与前端保持一致
	VIPORDERCARD    int32 = 1 // 次卡
	VIPCYCLECARD    int32 = 2 // 期卡
	VIPRECHARGECARD int32 = 3 // 充值卡
)

type ComboApi struct{}

// CreateVIPCombo
// @Tags      ExaCustomer
// @Summary   创建客户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      example.ExaCustomer            true  "客户用户名, 客户手机号码"
// @Success   200   {object}  response.Response{msg=string}  "创建客户"
// @Router    /customer/customer [post]
func (e *ComboApi) CreateVIPCombo(c *gin.Context) {
	var combo business.VIPCombo
	err := c.ShouldBindJSON(&combo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	combo.SysUserId = utils.GetUserID(c)
	combo.State = 1
	err = comboService.CreateVIPCombo(combo)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteExaCustomer
// @Tags      ExaCustomer
// @Summary   删除客户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      example.ExaCustomer            true  "客户ID"
// @Success   200   {object}  response.Response{msg=string}  "删除客户"
// @Router    /customer/customer [delete]
func (e *ComboApi) DeleteVIPComboById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = comboService.DeleteVIPComboById(reqId.ID, utils.GetUserID(c))
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateExaCustomer
// @Tags      ExaCustomer
// @Summary   更新客户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      example.ExaCustomer            true  "客户ID, 客户信息"
// @Success   200   {object}  response.Response{msg=string}  "更新客户信息"
// @Router    /customer/customer [put]
func (e *ComboApi) UpdateVIPCombo(c *gin.Context) {
	var combo business.VIPCombo
	err := c.ShouldBindJSON(&combo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(combo, utils.ComboVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = comboService.UpdateVIPCombo(&combo)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetExaCustomer
// @Tags      ExaCustomer
// @Summary   获取单一客户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     example.ExaCustomer                                                true  "客户ID"
// @Success   200   {object}  response.Response{data=exampleRes.ExaCustomerResponse,msg=string}  "获取单一客户信息,返回包括客户详情"
// @Router    /customer/customer [get]
func (e *ComboApi) GetVIPComboById(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := comboService.GetVIPComboById(reqId.ID, utils.GetUserID(c))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	var comboData businessRes.VipComboResModel
	comboData.StoreName = utils.GetUserName(c)
	comboData.ComboName = data.ComboName
	comboData.ComboType = data.ComboType
	comboData.ComboPrice = data.ComboPrice
	comboData.Times = data.Times
	comboData.State = data.State
	response.OkWithDetailed(businessRes.VipComboResponse{Combo: comboData}, "获取成功", c)
}

// GetVIPComboList
// @Tags      ExaCustomer
// @Summary   分页获取权限客户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限客户列表,返回包括列表,总数,页码,每页数量"
// @Router    /customer/customerList [get]
func (e *ComboApi) GetVIPComboList(c *gin.Context) {
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
	comboList, total, err := comboService.GetVIPComboInfoList(utils.GetUserID(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     comboList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (e *ComboApi) GetAllVIPCombos(c *gin.Context) {
	comboList, err := comboService.GetAllVIPComboInfoList(utils.GetUserID(c))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.AllResult{
		List: comboList,
	}, "获取成功", c)
}
