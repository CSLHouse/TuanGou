import request from '@/utils/request'

export function fetchList(params) {
  return request({
    url: '/flash/flashSessionList',
    method: 'get',
    params: params
  })
}

export function fetchSelectList(params) {
  return request({
    url: '/flash/flashSessionSelectList',
    method: 'get',
    params: params
  })
}

export function updateStatus(data) {
  return request({
    url: '/flash/updateFlashSessionStatus',
    method: 'post',
    data: data
  })
}

export function deleteSession(params) {
  return request({
    url: '/flash/deleteFlashSession',
    method: 'delete',
    params: params
  })
}

export function createSession(data) {
  return request({
    url: '/flash/createFlashSession',
    method: 'post',
    data: data
  })
}

export function updateSession(data) {
  return request({
    url: '/flash/updateFlashSession',
    method: 'put',
    data: data
  })
}
