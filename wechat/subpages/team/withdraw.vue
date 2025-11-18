<template>
	<view class="wp">
		<view class="header">
			<view class="txt-xh">
				<view class="h3">
					¥{{count}}
				</view>
				<text>可提现金额(元)</text>
				<text class="s2">?</text>
			</view>
		</view>
		<view class="boe">
			<view class="m-about">
				<view class="h3">
					提现金额
				</view>
				<view class="num">
					<view class="desc">￥</view> <text @click="onAllClick()">全部提现</text>
					<input class="inp" placeholder="请输出金额" :value="withdrawCount" type="number" disabled="true">
				</view>
				<!-- <view class="boe">
					<text>提现至</text>
					<text class="s1">微信零钱</text>
				</view> -->
				<view class="g-queding-kx" @click="onClick()">
					<view class="txt">
						<view class="s1">全部结算</view>
					</view>
				</view>
			</view>
		</view>
		<view class="main">
			<view class="txt-xh1">
				<view class="tit">
					<view class="h3">
						说明：
					</view>
					<view class="info">
						1.<text>可提现金额</text>=“首推佣金”+“成团奖励”+“成员复购奖励”
					</view>
					<view class="info">
						2.<text>提交结算后，请保存、扫描下方二维码添加客服，并发送编号<text style="color:#09c18f">{{userId}}</text> 完成“提现”</text>
					</view>
					<!-- 	<view class="desc">
						<text>委托出让</text>
						<text>继续认养猪崽</text>
						<text style="margin-right: 0;">联系客服提现</text>
					</view> -->
					<view class="info">
						3.<text>每日奖励</text>：每日累计可提现奖励200元
					</view>
				</view>
			</view>
		</view>
		<view class="qr-wrapper">
			<!-- 二维码图片 -->
			<image class="qr-code" :src="qrCodePath" mode="widthFix" @longpress="saveQrCode"></image>
			<text class="hint-text">长按二维码保存到相册</text>
		</view>
	</view>
</template>

<script>
	import {
		teamSettlement
	} from '@/api/team.js';
	export default {
		data() {
			return {
				count: 0,
				withdrawCount: null,
				userId: null,
				qrCodePath: '/static/service_qrcode.jpg',
				isLoading: false,
			};
		},
		onLoad(options) {
			if (options.count) {
				this.count = parseFloat(options.count)
				this.withdrawCount = parseFloat(options.count)
			}
			const userInfo = wx.getStorageSync("UserInfo")
			this.userId = userInfo.id
		},
		methods: {
			onAllClick() {
				this.withdrawCount = this.count
			},
			onClick() {
				if (this.isLoading) {
					return;
				}
				if (this.withdrawCount <= 0) {
					this.$api.msg("结算金额不足")
					return
				}
				let _this = this
				const openId = wx.getStorageSync("OpenId")
				if (!openId || openId.length < 1) {
					this.$api.msg("程序异常，请重新登录")
					return
				}
				let data = {
					openId: openId,
					amount: this.withdrawCount
				}

				teamSettlement(data).then(response => {
					this.isLoading = false;
					if (response.code == 0) {
						_this.$api.msg("结算成功！")
						setTimeout(function() {
							uni.navigateBack()
						}, 1500)
					}
				}).catch(error => {
					// 接口调用失败也需要恢复按钮状态
					this.isLoading = false;
					console.error('结算失败', error);
					this.$api.msg("结算失败，请稍后重试")
				});
			},
			saveQrCode() {
				// 显示加载提示
				uni.showLoading({
					title: '保存中...'
				})

				// 检查是否有保存图片权限
				uni.getSetting({
					success: (res) => {
						// 如果没有权限，请求授权
						if (!res.authSetting['scope.writePhotosAlbum']) {
							uni.authorize({
								scope: 'scope.writePhotosAlbum',
								success: () => {
									this.saveImage()
								},
								fail: () => {
									// 用户拒绝授权，引导到设置页面开启
									uni.hideLoading()
									uni.showModal({
										title: '权限提示',
										content: '需要开启保存图片权限才能保存二维码',
										confirmText: '去设置',
										success: (modalRes) => {
											if (modalRes.confirm) {
												uni.openSetting()
											}
										}
									})
								}
							})
						} else {
							// 已有权限，直接保存
							this.saveImage()
						}
					}
				})
			},

			// 执行保存图片操作
			saveImage() {
				// 处理本地图片
				uni.getImageInfo({
					src: this.qrCodePath,
					success: (imageInfo) => {
						uni.saveImageToPhotosAlbum({
							filePath: imageInfo.path,
							success: () => {
								uni.hideLoading()
								uni.showToast({
									title: '保存成功',
									icon: 'success'
								})
							},
							fail: () => {
								uni.hideLoading()
								uni.showToast({
									title: '保存失败',
									icon: 'none'
								})
							}
						})
					},
					fail: () => {
						uni.hideLoading()
						uni.showToast({
							title: '获取图片失败',
							icon: 'none'
						})
					}
				})
			}

		}
	}
</script>

<style>
	.header {
		background-color: #1ec796;
		position: relative;
		height: 400rpx;
	}

	.txt-xh {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translateY(-50%) translateX(-50%);
	}

	.txt-xh .h3 {
		font-weight: bold;
		color: #fff;
		margin-bottom: 5rpx;
		font-size: 52rpx;
	}

	.txt-xh text {
		color: #fff;
		font-size: 24rpx;
	}

	.txt-xh .s2 {
		border: 1px solid #fff;
		border-radius: 50%;
		display: inline-block;
		margin-left: 15rpx;
		font-size: 14rpx;
		width: 30rpx;
		height: 30rpx;
		text-align: center;
		line-height: 30rpx;
	}

	.m-about {
		margin: 0 40rpx;
		padding: 20rpx 40rpx;
		margin-top: -80rpx;
		background-color: #fff;
		border-radius: 10rpx;
		position: relative;
		z-index: 99;
		box-shadow: 0 0 20rpx rgba(0, 0, 0, 0.1);
		margin-bottom: 100rpx;
	}

	.m-about .h3 {
		font-size: 28rpx;
		margin-bottom: 25rpx;
	}

	.m-about .num {
		margin-bottom: 35rpx;
		border-bottom: 1px solid #f2f2f2;
	}

	.m-about .num .desc {
		float: left;
		font-size: 24rpx;
		margin-bottom: 15rpx;
	}

	.m-about .num text {
		float: right;
		color: #09c18f;
		font-size: 24rpx;
		margin-bottom: 15rpx;
	}

	.m-about .boe {
		overflow: hidden;
		margin-bottom: 100rpx;
		color: #666666;
		font-size: 24rpx;
	}

	.m-about .boe text {
		float: left;
	}

	.m-about .boe .s1 {
		float: right;
	}

	uni-input {
		height: 1.4em;
		font-size: 24rpx;
		line-height: 1.4em;
		margin-bottom: 15rpx;
	}

	.g-queding-kx {
		text-align: center;
		background-color: #09c18f;
		height: 80rpx;
		line-height: 80rpx;
		font-size: 28rpx;
		border-radius: 10rpx;
		color: #fff;
		width: 100%;
	}

	.txt-xh1 {
		padding: 0 40rpx;
	}

	.txt-xh1 .h3 {
		display: block;
		font-size: 24rpx;
		margin-bottom: 20rpx;
	}

	.txt-xh1 .info {
		display: inline-block;
		font-size: 24rpx;
		letter-spacing: 1rpx;
		margin-bottom: 15rpx;
	}

	.txt-xh1 .info text {
		color: #000;
	}

	.txt-xh1 .info {
		color: #666666;

	}

	.txt-xh1 .desc {
		font-size: 28rpx;
		color: #1ec796;
		margin-bottom: 20rpx;
		margin-top: 5rpx;

	}

	.txt-xh1 .desc text {
		margin-right: 85rpx;

	}

	.qr-wrapper {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 40rpx;
		background-color: #ffffff;
		border-radius: 16rpx;
		box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
	}

	.qr-code {
		width: 400rpx;
		height: 400rpx;
		margin-bottom: 20rpx;
	}

	.hint-text {
		font-size: 28rpx;
		color: #666666;
	}
</style>