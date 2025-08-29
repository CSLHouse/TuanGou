import service from '@/utils/request'

// @Summary 创建套餐
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "创建套餐"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
export const createExaVIPCombo = (data) => {
  return service({
    url: '/business/combo',
    method: 'post',
    data
  })
}

// @Tags SysApi
// @Summary 更新套餐
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "更新套餐"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
export const updateExaVIPCombo = (data) => {
  return service({
    url: '/business/combo',
    method: 'put',
    data
  })
}

// @Tags SysApi
// @Summary 删除套餐
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "删除套餐"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
export const deleteExaVIPCombo = (params) => {
  return service({
    url: '/business/combo',
    method: 'delete',
    params: params,
  })
}

// @Tags SysApi
// @Summary 获取单一套餐
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "获取单一套餐"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
export const getExaVIPCombo = (params) => {
  return service({
    url: '/business/combo',
    method: 'get',
    params
  })
}

// @Tags SysApi
// @Summary 获取套餐列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取套餐列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
export const getExaVIPComboList = (params) => {
  return service({
    url: '/business/comboList',
    method: 'get',
    params
  })
}


export const getAllVIPComboList = (params) => {
  return service({
    url: '/business/allCombo',
    method: 'get',
    params
  })
}