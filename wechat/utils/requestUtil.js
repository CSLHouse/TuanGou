import Request from '@/js_sdk/luch-request/request.js'
import store from '../store'
import {
	wxRefreshLogin
} from '@/api/member.js';
import common from '@/utils/common.js'

const http = new Request()

// 全局配置
http.setConfig((config) => {
	config.baseUrl = common.baseUrl
	config.header = {
		"Access-Control-Allow-Origin": "*",
		"Access-Control-Allow-Methods": "*",
		...config.header,
	}
	config.timeout = 15000 // 超时时间设置
	return config
})

/**
 * 自定义状态码验证
 * @param { Number } statusCode - 响应状态码
 * @return { Boolean } 是否验证通过
 */
http.validateStatus = (statusCode) => {
	return statusCode >= 200 && statusCode < 300
}

// 请求拦截器
http.interceptor.request((config, cancel) => {
	const token = store.state.token
	// 设置请求头
	config.header = {
		'x-token': token || '',
		'x-user_id': store.state.userInfo?.id || '',
		...config.header
	}

	// 取消请求示例
	// if (!token && config.url !== '/auth/login') {
	// 	cancel('未登录，请先登录')
	// }

	return config
})

// 响应拦截器
http.interceptor.response((response) => {
	// 处理新token
	if (response.header['new-token']) {
		uni.setStorageSync('Token', response.header['new-token'])
		store.state.token = response.header['new-token']
	}

	const res = response.data

	// 业务逻辑错误处理
	if (res.code !== 0) {
		// 登录状态过期处理
		if (res.code === 401 || (res.code === 7 && res.data?.reload)) {
			// 防止重复弹窗
			if (!store.state.isRefreshingToken) {
				store.state.isRefreshingToken = true

				wxRefreshLogin({
					openId: store.state.openId
				}).then(refreshRes => {
					store.state.isRefreshingToken = false
					if (refreshRes.code === 0) {
						const userinfo = refreshRes.data
						uni.setStorageSync("Token", userinfo.token)
						uni.setStorageSync("TokenTime", userinfo.expiresAt)
						store.state.token = userinfo.token
						store.state.userInfo = userinfo.customer
						store.state.hasLogin = true
					} else {
						// 刷新token失败，需要重新登录
						uni.showModal({
							title: '提示',
							content: '登录已过期，请重新登录',
							showCancel: false,
							success: () => {
								store.dispatch('logout')
								uni.reLaunch({
									url: '/pages/index/index'
								})
							}
						})
					}
				}).catch(() => {
					store.state.isRefreshingToken = false
					uni.showModal({
						title: '提示',
						content: '登录已过期，请重新登录',
						showCancel: false,
						success: () => {
							store.dispatch('logout')
							uni.reLaunch({
								url: '/pages/index/index'
							})
						}
					})
				})
			}
			return Promise.reject(res)
		}

		// 显示错误信息
		uni.showToast({
			title: res.message || '操作失败',
			icon: 'none',
			duration: 2000
		})

		return Promise.reject(res)
	}

	return res
}, (response) => {
	// 网络错误处理
	let errorMsg = '网络异常，请稍后重试'
	if (response.errMsg.includes('timeout')) {
		errorMsg = '请求超时，请稍后重试'
	}

	uni.showToast({
		title: errorMsg,
		icon: 'none',
		duration: 2000
	})

	return Promise.reject(response)
})

export function request(options = {}) {
	return http.request(options)
}

export default request