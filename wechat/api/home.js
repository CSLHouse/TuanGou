import request from '@/utils/requestUtil'

export function fetchContent(params) {
	return request({
		method: 'GET',
		url: '/product/content',
		params: params
	})
}

export function fetchRecommendProductList(params) {
	return request({
		method: 'GET',
		url: '/product/recommendProductList',
		params:params
	})
}

export function fetchProductCateList(parentId) {
	return request({
		method: 'GET',
		url: '/product/productCategory',
		params: {"tag": parentId}
	})
}

export function fetchNewProductList(params) {
	return request({
		method: 'GET',
		url: '/product/newProductList',
		params:params
	})
}

export function fetchHotProductList(params) {
	return request({
		method: 'GET',
		url: '/product/hotProductList',
		params:params
	})
}

export function recordShareCount(data) {
	return request({
		method: 'POST',
		url: '/base/recordShare',
		data: data
	})
}