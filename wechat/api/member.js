import request from '@/utils/requestUtil'

export function memberLogin(data) {
	return request({
		method: 'POST',
		url: '/base/wxLogin',
		header: {
			'content-type': 'application/x-www-form-urlencoded;charset=utf-8'
		},
		data: data
	})
}

export function memberInfo() {
	return request({
		method: 'GET',
		url: '/sso/info'
	})
}

export function wxLogin(data) {
	return request({
		method: 'POST',
		url: '/base/wxLogin',
		header: {
			'content-type': 'application/json'
		},
		data: data
	})
}

export function getWXPhoneNumber(data) {
	return request({
		method: 'POST',
		url: '/base/phoneNumber',
		header: {
			'content-type': 'application/json'
		},
		data: data
	})
}

export function UpdateWxUserInfo(data) {
	return request({
		method: 'POST',
		url: '/base/wxUserInfo',
		header: {
			'content-type': 'application/json'
		},
		data: data
	})
}

export function GetWxUserInfo(params) {
	return request({
		method: 'GET',
		url: '/base/wxUserInfo',
		params: params
	})
}

export function wxRefreshLogin(data) {
	return request({
		method: 'POST',
		url: '/base/wxRefreshLogin',
		header: {
			'content-type': 'application/json'
		},
		data: data
	})
}

export function GetMemberCardList(params) {
	return request({
		method: 'GET',
		url: '/business/cardList',
		params: params
	})
}

export function GetCertificateList(params) {
	return request({
		method: 'GET',
		url: '/business/certificateList',
		params: params
	})
}

export function WXResetNickName(data) {
	return request({
		method: 'POST',
		url: '/user/resetNickName',
		header: {
			'content-type': 'application/json'
		},
		data: data
	})
}
