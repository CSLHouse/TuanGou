import service from '@/utils/request'
// export function fetchList(params) {
//   return service({
//     url:'/promo/newProductList',
//     method:'get',
//     params:params
//   })
// }

export function updateRecommendStatus(data) {
  return service({
    url:'/promo/newProductState',
    method:'put',
    data:data
  })
}

export function deleteNewProduct(params) {
  return service({
    url:'/promo/newProduct',
    method:'delete',
    params:params
  })
}

export function createNewProduct(data) {
  return service({
    url:'/promo/newProduct',
    method:'post',
    data:data
  })
}

export function updateNewProductSort(data) {
  return service({
    url:'/promo/newProductSort',
    method:'put',
    data:data
  })
}
