import request from '@/utils/requestUtil'

export function fetchAddressList(params) {
	return request({
		method: 'GET',
		url: '/base/addressList',
		params: params,
	})
}

export function fetchAddressDetail(id) {
	return request({
		method: 'GET',
		url: `/base/address`,
		params: {id: id}
	})
}

export function addAddress(data) {
	return request({
		method: 'POST',
		url: '/base/address',
		data:data
	})
}

export function updateAddress(data) {
	return request({
		method: 'PUT',
		url: `/base/address`,
		data:data
	})
}

export function deleteAddress(id) {
	return request({
		method: 'DELETE',
		url: `/base/address`,
		data: {id: id}
	})
}

