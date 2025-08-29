import request from '@/utils/requestUtil'

export function addCartItem(data) {
	return request({
		method: 'POST',
		url: '/order/cart',
		data: data
	})
}

export function fetchCartList() {
	return request({
		method: 'GET',
		url: '/order/cart/list'
	})
}

export function deletCartItem(params) {
	return request({
		method: 'DELETE',
		url: '/order/cart',
		params:params
	})
}

export function deletCartItemWithList(params) {
	return request({
		method: 'DELETE',
		url: '/order/carts',
		params:params
	})
}

export function updateQuantity(params) {
	return request({
		method: 'PUT',
		url: '/order/cart',
		params:params
	})
}

export function clearCartList() {
	return request({
		method: 'POST',
		url: '/order/cart/clear'
	})
}


export function addCartTmpItem(data) {
	return request({
		method: 'POST',
		url: '/order/tmpCart',
		data: data
	})
}