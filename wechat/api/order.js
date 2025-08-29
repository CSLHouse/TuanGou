import request from '@/utils/requestUtil'

export function generateConfirmOrder(data) {
	return request({
		method: 'POST',
		url: '/order/generateConfirmOrder',
		data: data
	})
}

export function generateOrder(data) {
	return request({
		method: 'POST',
		url: '/order/generateOrder',
		data: data
	})
}

export function fetchOrderList(params) {
	return request({
		method: 'GET',
		url: '/order/list',
		params: params
	})
}

export function payOrderSuccess(data) {
	return request({
		method: 'POST',
		url: '/order/paySuccess',
		// header: {
		// 	'content-type': 'application/x-www-form-urlencoded;charset=utf-8'
		// },
		data: data
	})
}

export function fetchOrderDetail(params) {
	return request({
		method: 'GET',
		url: `/order/detail`,
		params: params
	})
}

export function cancelUserOrder(data) {
	return request({
		method: 'POST',
		url: '/order/cancelOrder',
		// header: {
		// 	'content-type': 'application/x-www-form-urlencoded;charset=utf-8'
		// },
		data: data
	})
}

export function confirmReceiveOrder(data) {
	return request({
		method: 'POST',
		url: '/order/confirmReceiveOrder',
		header: {
			'content-type': 'application/x-www-form-urlencoded;charset=utf-8'
		},
		data: data
	})
}

export function deleteUserOrder(data) {
	return request({
		method: 'POST',
		url: '/order/deleteOrder',
		header: {
			'content-type': 'application/x-www-form-urlencoded;charset=utf-8'
		},
		data: data
	})
}