import request from '@/utils/requestUtil'

export function searchProductList(params) {
	return request({
		method: 'GET',
		url: '/product/productList',
		params: params
	})
}

export function fetchCategoryTreeList() {
	return request({
		method: 'GET',
		url: '/product/allCategory'
	})
}

export function fetchProductDetail(id) {
	return request({
		method: 'GET',
		url: '/product/productDetail?id='+id
	})
}
