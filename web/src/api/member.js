import service from '@/utils/request'

// @Summary 创建套餐
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "创建套餐"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
export const createVIPMember = (data) => {
  return service({
    url: '/business/member',
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
export const updateVIPMember = (data) => {
  return service({
    url: '/business/member',
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
export const deleteVIPMember = (params) => {
  return service({
    url: '/business/member',
    method: 'delete',
    params: params
  })
}

// @Tags SysApi
// @Summary 获取单一套餐
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.ExaCustomer true "获取单一套餐"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
export const getVIPMember = (params) => {
  return service({
    url: '/business/member',
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
export const getVIPMemberList = (params) => {
  return service({
    url: '/business/memberList',
    method: 'get',
    params
  })
}


export const getAllVIPMemberList = (params) => {
  return service({
    url: '/business/allMember',
    method: 'get',
    params
  })
}

export const searchVIPMembers = (params) => {
  return service({
    url: '/business/memberSearch',
    method: 'post',
    params
  })
}

export const searchVIPCards = (params) => {
  return service({
    url: '/business/searchCard',
    method: 'post',
    params
  })
}

export const renewVIPCards = (params) => {
  return service({
    url: '/business/renew',
    method: 'post',
    params
  })
}

export const consumeVIPCard = (params) => {
  return service({
    url: '/business/consume',
    method: 'post',
    params: params
  })
}

export const getVIPConsumeList = (params) => {
  return service({
    url: '/business/consumeList',
    method: 'get',
    params
  })
}

export const getVIPOrderList = (params) => {
  return service({
    url: '/business/orderList',
    method: 'get',
    params
  })
}

export const getStatisticsList = (params) => {
  return service({
    url: '/business/statistics',
    method: 'get',
    params
  })
}

export const getStatementList = (params) => {
  return service({
    url: '/business/statement',
    method: 'get',
    params
  })
}