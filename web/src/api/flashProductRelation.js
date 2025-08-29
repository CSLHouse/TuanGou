import request from '@/utils/request'
export function fetchList(params) {
  return request({
    url:'/flash/flashProductRelationList',
    method:'get',
    params:params
  })
}
export function createFlashProductRelation(data) {
  return request({
    url:'/flash/createFlashProductRelation',
    method:'post',
    data: data
  })
}
export function deleteFlashProductRelation(params) {
  return request({
    url:'/flash/deleteFlashProductRelation',
    method:'delete',
    params:params
  })
}
export function updateFlashProductRelation(data) {
  return request({
    url:'/flash/updateFlashProductRelation',
    method:'put',
    data:data
  })
}
