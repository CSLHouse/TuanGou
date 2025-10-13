import request from '@/utils/requestUtil'

export function fetchTeamList(params) {
	return request({
		method: 'GET',
		url: '/team/teamList',
		params: params,
	})
}

export function fetchTeamDetail(data) {
	return request({
		method: 'POST',
		url: `/team/detail`,
		data: data
	})
}

export function fetchTeamReward(params) {
	return request({
		method: 'GET',
		url: `/team/reward`,
		params: params
	})
}