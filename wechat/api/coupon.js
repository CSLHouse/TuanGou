import request from '@/utils/requestUtil'

export function fetchProductCouponList(params) {
	return request({
		method: 'GET',
		url: `/coupon/listByProduct`,
		params: params
	})
}

export function addMemberCoupon(params) {
	return request({
		method: 'POST',
		url: `/coupon/add`,
		params: params
	})
}

export function fetchMemberCouponList(useStatus) {
	return request({
		method: 'GET',
		url: '/coupon/listWithState',
		params: {
			useStatus: useStatus
		}
	})
}