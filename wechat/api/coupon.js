import request from '@/utils/requestUtil'

export function fetchProductCouponList(params) {
	return request({
		method: 'GET',
		url: `coupon/listByProduct`,
		params: params
	})
}

export function addMemberCoupon(couponId) {
	return request({
		method: 'POST',
		url: `/member/coupon/add/${couponId}`,
	})
}

export function fetchMemberCouponList(useStatus) {
	return request({
		method: 'GET',
		url: '/member/coupon/list',
		params: {
			useStatus: useStatus
		}
	})
}