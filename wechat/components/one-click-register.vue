<template>
	<view v-if='!hasLogin && !isCloseModel'>
		<div class="modal-mask" @click="closePop">
		</div>
		<div class="modal-dialog">
			<div class="modal-content">
				<image class="img" src="/static/pop.jpg"></image>
				<div class="content-text">
					<p class="little-tip">我们的生活圈：</p>
					<p class="little-content">
						注册成为会员，一店消费，多家优惠，欢迎体验
					</p>
				</div>
			</div>
			<div class="modal-footer">
				<button class='btn' open-type='getPhoneNumber' @getphonenumber="decryptPhoneNumber">
					一键注册
				</button>
			</div>
		</div>
	</view>
</template>

<script>
	import {
		mapState,
		mapMutations
	} from 'vuex';
	import {
		getWXPhoneNumber,
		wxRefreshLogin,
	} from '@/api/member.js';
	export default {
		name: "one-click-register",
		data() {
			return {
				isCloseModel: false,
			};
		},
		// props: {
		// 	hasLogin: {
		// 		type: Boolean,
		// 		default: false
		// 	}
		// },
		computed: {
			...mapState(['hasLogin']),
		},
		methods: {
			...mapMutations(['login']),
			closePop() {
				this.isCloseModel = true
			},
			decryptPhoneNumber: function(e) {
				let _this = this
				console.log("--------decryptPhoneNumber---------", e.detail)
				if (e.detail.errMsg == "getPhoneNumber:ok") {
					console.log('-----getWXPhoneNumber---openId---', _this.$store.state.openId)
					if (_this.$store.state.openId && _this.$store.state.openId.length > 0) {
						getWXPhoneNumber({
							openId: _this.$store.state.openId,
							code: e.detail.code
						}).then(res => {
							console.log('-----getWXPhoneNumber------', res)
							if (res.code == 0) {
								_this.getToken()
								_this.$api.msg('注册成功')
							} else {
								_this.$api.msg('注册会员失败')
							}
						});
					}
				}
			},
			getToken() {
				let _this = this
				wxRefreshLogin({
					openId: _this.$store.state.openId
				}).then(res => {
					if (res.code == 0) {
						const userinfo = res.data
						wx.setStorageSync("Token", userinfo.token)
						wx.setStorageSync("TokenTime", userinfo.expiresAt)
						_this.$store.state.token = userinfo.token
						_this.login(userinfo.customer);
						_this.isCloseModel = true
						console.log("----isCloseModel---", _this.isCloseModel)
					}
				}).catch(errors => {
					uni.showModal({
						title: '提示',
						content: '网络错误',
						showCancel: false
					})
				});
			},
		}
	}
</script>

<style>
	.modal-mask {
		width: 100%;
		height: 100%;
		position: fixed;
		top: 0;
		left: 0;
		background: #000;
		opacity: 0.5;
		overflow: hidden;
		z-index: 9000;
		color: #fff;
	}

	.modal-dialog {
		box-sizing: border-box;
		width: 560rpx;
		overflow: hidden;
		position: fixed;
		top: 40%;
		left: 0;
		z-index: 9999;
		background: #fff;
		margin: -150rpx 95rpx;
		border-radius: 16rpx;
	}

	.modal-content {
		box-sizing: border-box;
		display: flex;
		padding: 0rpx 53rpx 50rpx 53rpx;
		font-size: 32rpx;
		align-items: center;
		justify-content: center;
		flex-direction: column;
	}

	.content-tip {
		text-align: center;
		font-size: 36rpx;
		color: #333333;
	}

	.content-text {
		/* height:230px; */
		padding: 10px 0px 10px 0px;
		font-size: 14px;
	}

	.modal-footer {
		box-sizing: border-box;
		display: flex;
		flex-direction: row;
		border-top: 1px solid #e5e5e5;
		font-size: 16px;
		font-weight: bold;
		/* height: 45px; */
		line-height: 45px;
		text-align: center;
		background: #d60000;
	}

	.btn {
		width: 100%;
		height: 100%;
		background: #d60000;
		color: #FFFFFF;
		font-weight: bold;
	}

	.img {
		width: 560rpx;
		height: 240rpx;
	}

	.little-tip {
		padding-top: 15px;
		padding-bottom: 3px;
		font-size: 14px;
		font-weight: bold;
		color: #d60000;
	}

	.little-content {
		padding-top: 5px;
		font-size: 13px;
		color: #606060;
	}

	.key-bold-tip {
		padding-top: 5px;
		font-size: 15px;
		font-weight: bold;
		color: #d60000;
	}

	.key-bold {
		padding-top: 5px;
		font-size: 14px;
		/* font-weight:bold; */
	}

	.info-bold-tip {
		padding-top: 5px;
		font-size: 15px;
		font-weight: bold;
		color: #d60000;
		text-align: center;
	}
</style>