import service from '@/utils/request'

export const getSettlementList = (params) => {
  return service({
    url: '/team/settlementList',
    method: 'get',
    params
  })
}


export const updateSettlementState = (data) => {
  return service({
    url: '/team/settlementUpdate',
    method: 'post',
    data
  })
}