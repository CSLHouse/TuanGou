import request from '@/utils/requestUtil'

export function fetchContent(params) {
	return request({
		method: 'GET',
		url: '/promo/content',
		params: params
	})
}

export function fetchRecommendProductList(params) {
	return request({
		method: 'GET',
		url: '/promo/recommendProductList',
		params: params
	})
}

export function fetchProductCateList(parentId) {
	return request({
		method: 'GET',
		url: '/product/productCategory',
		params: {
			"tag": parentId
		}
	})
}

export function fetchNewProductList(params) {
	return request({
		method: 'GET',
		url: '/promo/newProductList',
		params: params
	})
}

export function fetchHotProductList(params) {
	return request({
		method: 'GET',
		url: '/promo/recommendProductList',
		params: params
	})
}