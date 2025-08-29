import request from '@/utils/request'
export function fetchList(params) {
  return request({
    url:'/flash/list',
    method:'get',
    params:params
  })
}
export function updateStatus(id,params) {
  return request({
    url:'/flash/update/status/'+id,
    method:'post',
    params:params
  })
}
export function deleteFlash(params) {
  return request({
    url:'/flash/delete',
    method:'delete',
    params: params
  })
}
export function createFlash(data) {
  return request({
    url:'/flash/create',
    method:'post',
    data:data
  })
}
export function updateFlash(data) {
  return request({
    url:'/flash/update',
    method:'put',
    data:data
  })
}
