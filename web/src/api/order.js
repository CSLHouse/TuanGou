import request from '@/utils/request'
export function fetchList(params) {
  return request({
    url:'/order/list',
    method:'get',
    params:params
  })
}

export function closeOrder(data) {
  return request({
    url:'/order/closeOrders',
    method:'post',
    data:data
  })
}

export function deleteOrder(data) {
  return request({
    url:'/order/delete',
    method:'delete',
    data:data
  })
}

export function cancelOrder(data) {
  return request({
    url:'/order/cancelOrder',
    method:'post',
    data:data
  })
}
export function deliveryOrder(data) {
  return request({
    url:'/order/update/delivery',
    method:'post',
    data:data
  });
}

export function getOrderDetail(params) {
  return request({
    url:'/order/detail',
    method:'get',
    params: params,
  });
}

export function updateReceiverInfo(data) {
  return request({
    url:'/order/update/receiverInfo',
    method:'post',
    data:data
  });
}

export function updateMoneyInfo(data) {
  return request({
    url:'/order/update/moneyInfo',
    method:'post',
    data:data
  });
}

export function updateOrderNote(data) {
  return request({
    url:'/order/update/note',
    method:'post',
    data:data
  })
}

export function updateOrderCompletedStatus(data) {
  return request({
    url:'/order/update/complete',
    method:'post',
    data:data
  })
}

export function getOrderSetting(params) {
  return request({
    url:'/order/setting',
    method:'get',
    params:params
  })
}

export function updateOrderSetting(data) {
  return request({
    url:'/order/settingUpdate',
    method:'post',
    data:data
  });
}
