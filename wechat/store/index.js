import Vue from 'vue'
import Vuex from 'vuex'
import {
	memberInfo,
	wxLogin,
	GetWxUserInfo
} from '@/api/member.js';

Vue.use(Vuex)

const store = new Vuex.Store({
	state: {
		hasLogin: false,
		userInfo: {},
		openId: '',
		token: "",
		hadNickName: false,
		formOpenId: '',
	},
	mutations: {
		refreshLoginSession() {
			let _this = this
			uni.login({
				provider: 'weixin', //使用微信登录
				success: function(loginRes) {
					wxLogin({
						code: loginRes.code
					}).then(res => {
						if (res.code == 0) {
							uni.setStorage({ //缓存用户登陆状态
								key: 'OpenId',
								data: res.data.openid
							})
							// wx.setStorageSync("WxCode", loginRes.code)
							// wx.setStorageSync("WxCodeTime", (new Date()).getTime())
							_this.state.openId = res.data.openid
						} else {
							uni.showToast({
								title: res.data,
								icon: 'none'
							})
						}
					}).catch(errors => {
						// uni.showModal({
						// 	title:'提示',
						// 	content:'登录失败',
						// 	showCancel:false
						// })
					});
				}
			});
		},
		login(state, provider) {
			state.hasLogin = true;
			state.userInfo = provider;
			if (provider.userName) {
				state.hadNickName = true
			}
			uni.setStorage({ //缓存用户登陆状态
				key: 'UserInfo',
				data: provider
			})
		},
		logout(state) {
			state.hasLogin = false
			state.hadNickName = false
			state.userInfo = {}
			state.token = null
			uni.removeStorage({
				key: 'UserInfo'
			});
			uni.removeStorage({
				key: 'Token'
			})
		},
		// setFormOpenId(state, provider) {
		// 	state.formOpenId = provider
		// }
	},
	getters: {
		// fetchFormOpenId(state) {
		// 	return state.formOpenId
		// }
	},
	actions: {
		// updateFormOpoenId({
		// 	commit
		// }, openId) {
		// 	commit('setFormOpenId', openId);
		// },
	}
})

export default store