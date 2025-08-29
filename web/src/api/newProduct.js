import service from '@/utils/request'
export function fetchList(params) {
  return service({
    url:'/product/newProductList',
    method:'get',
    params:params
  })
}

export function updateRecommendStatus(data) {
  return service({
    url:'/product/newProductState',
    method:'put',
    data:data
  })
}

export function deleteNewProduct(params) {
  return service({
    url:'/product/newProduct',
    method:'delete',
    params:params
  })
}

export function createNewProduct(data) {
  return service({
    url:'/product/newProduct',
    method:'post',
    data:data
  })
}

export function updateNewProductSort(data) {
  return service({
    url:'/product/newProductSort',
    method:'put',
    data:data
  })
}
