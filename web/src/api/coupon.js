import service from '@/utils/request'
export function fetchList(params) {
  return service({
    url:'/coupon/list',
    method:'get',
    params:params
  })
}

export function createCoupon(data) {
  return service({
    url:'/coupon/create',
    method:'post',
    data:data
  })
}

export function getCoupon(params) {
  return service({
    url:'/coupon/details',
    method:'get',
    params:params
  })
}

export function updateCoupon(data) {
  return service({
    url:'/coupon/update/',
    method:'post',
    data:data
  })
}

export function deleteCoupon(params) {
  return service({
    url:'/coupon/delete/',
    method:'post',
    params:params

  })
}


export function fetchCouponHistoryList(params) {
  return service({
    url:'/coupon/couponHistory',
    method:'get',
    params:params
  })
}