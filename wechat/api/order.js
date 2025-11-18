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

export function fetchOrderItemList(params) {
	return request({
		method: 'GET',
		url: '/order/item/list',
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

export function queryLogisticsInfo(data) {
	return request({
		method: 'POST',
		url: '/logistics/info',
		data: data
	})
}


import common from '@/utils/common.js'
import store from '../store'
/**
 * 微信平台图片上传
 * @param {Array} imgs - 图片信息数组，包含name和uri
 * @returns {Promise}
 */
export function UploadFileEx(file) {
	return new Promise((resolve, reject) => {
		const token = store.state.token
		let fromData = {
			url: common.baseUrl + "/order/upload/file",
			dataType: 'json',
			header: {
				'x-token': token || '',
				'x-user_id': store.state.userInfo?.id || '',
			},
			success: res => {
				let resp = JSON.parse(res.data);
				if (resp.code == 0) { // 正常返回 
					resolve(resp)
				} else { // 错误返回
					resolve(resp);
				}
			},
			fail: res => {
				reject(res);
				uni.showToast({
					title: res,
					icon: 'none'
				})
			}
		};

		fromData.filePath = file.uri;
		fromData.name = "file";
		uni.uploadFile(fromData);
	});
}

/**
 * 非微信平台图片上传
 * @param {Array} file - 图片信息数组，包含name和uri
 * @returns {Promise}
 */
export function UploadFile(file) {
	return new Promise((resolve, reject) => {
		const token = store.state.token
		let fromData = {
			url: common.baseUrl + "/upload/file",
			dataType: 'json',
			header: {
				'x-token': token || '',
				'x-user_id': store.state.userInfo?.id || '',
			},
			success: res => {
				let resp = JSON.parse(res.data);
				if (resp.state) { // 正常返回 
					resolve({
						"success": resp.data,
						"fail": null
					})
				} else { // 错误返回
					resolve({
						"success": null,
						"fail": resp.data
					});
				}
			},
			fail: res => {
				reject(res);
				uni.showToast({
					title: res,
					icon: 'none'
				})
			}
		};

		fromData.files = files;
		uni.uploadFile(fromData);
	});
}

/**
 * 提交订单反馈
 * @param {string} orderItemId - 订单号
 * @param {number} type - 处理类型（1表示反馈）
 * @param {string} content - 反馈内容
 * @param {string} contact - 联系方式
 * @param {Array} images - 上传后的图片列表
 * @returns {Promise}
 */
export function DealOrder(orderItemId, content, contact, images) {
	return request({
		method: 'POST',
		url: '/order/deal',
		data: {
			orderItemId,
			content,
			contact,
			images
		}
	});
}


// 上传多个文件
export async function UploadFileWx(files) {
	const resp = [];
	try {
		// 循环上传每个文件（串行上传，也可根据需求改为并行Promise.all）
		for (const file of files) {
			const res = await UploadFileEx(file);
			if (res.code === 0 && res.data.id && res.data.id > 0) {
				resp.push(res.data.id);
			} else {
				// 单个文件上传失败则整体失败
				throw new Error(`文件${file.name}上传失败: ${res.msg || '未知错误'}`);
			}
		}
		return resp;
	} catch (error) {
		console.error('多文件上传失败:', error);
		return null; // 或根据需求抛出错误让上层处理
	}
}