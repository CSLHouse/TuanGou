<template>
	<view class="container">
		<view class="user-section">
			<image class="bg" src="/static/user-bg.jpg"></image>
			<view class="user-info-box">
				<view class="portrait-box" @click="handleUserInfo">
					<image class="portrait" :src="userInfo.avatarUrl || '/static/missing-face.png'"></image>
				</view>
				<view class="info-box">
					<text class="username"
						@click="goLogin">{{hasLogin ? userInfo.userName || userInfo.telephone : '立即登录'}}</text>
				</view>
			</view>
		</view>
		<!-- 积分、成长值、优惠券 -->
		<view class="cover-container" :style="[{
				transform: coverTransform,
				transition: coverTransition
			}]" @touchstart="coverTouchstart" @touchmove="coverTouchmove" @touchend="coverTouchend">
			<image class="arc" src="/static/arc.png"></image>

			<view class="tj-sction">
				<view class="tj-item">
					<text class="num">{{rewardAmount || '暂无'}}</text>
					<text>分红</text>
				</view>
				<view class="tj-item">
					<text class="num">{{userInfo.integration || '暂无'}}</text>
					<text>积分</text>
				</view>
				<view class="tj-item" @click="navTo('/subpages/coupon/couponList')">
					<text class="num">{{couponCount || '暂无'}}</text>
					<text>优惠券</text>
				</view>
			</view>
			<view class="order-section">
				<view class="order-item" @click="navTo('/subpages/order/order?state=0')" hover-class="common-hover"
					:hover-stay-time="50">
					<text class="yticon icon-shouye"></text>
					<text>全部订单</text>
				</view>
				<view class="order-item" @click="navTo('/subpages/order/order?state=1')" hover-class="common-hover"
					:hover-stay-time="50">
					<text class="yticon icon-daifukuan"></text>
					<text>待付款</text>
				</view>
				<view class="order-item" @click="navTo('/subpages/order/order?state=2')" hover-class="common-hover"
					:hover-stay-time="50">
					<text class="yticon icon-yishouhuo"></text>
					<text>待收货</text>
				</view>
				<view class="order-item" @click="navTo('/subpages/order/order?state=4')" hover-class="common-hover"
					:hover-stay-time="50">
					<text class="yticon icon-shouhoutuikuan"></text>
					<text>退款/售后</text>
				</view>
			</view>
			<!-- 浏览历史 -->
			<view class="history-section icon">
				<list-cell icon="icon-shoucang" iconColor="#5fcda2" title="团队管理"
					@eventClick="navTo('/subpages/team/team')"></list-cell>
				<list-cell icon="icon-dizhi" iconColor="#5fcda2" title="地址管理"
					@eventClick="navTo('/subpages/address/address')"></list-cell>
				<!-- <list-cell icon="icon-lishijilu" iconColor="#e07472" title="我的足迹" @eventClick="navTo('/subpages/user/readHistory')"></list-cell>
				<list-cell icon="icon-shoucang" iconColor="#5fcda2" title="我的关注" @eventClick="navTo('/subpages/user/brandAttention')"></list-cell>
				<list-cell icon="icon-shoucang_xuanzhongzhuangtai" iconColor="#54b4ef" title="我的收藏" @eventClick="navTo('/subpages/user/productCollection')"></list-cell>
				<list-cell icon="icon-pingjia" iconColor="#ee883b" title="我的评价"></list-cell> -->
				<list-cell icon="icon-pingjia" iconColor="#ee883b" title="收益榜" @eventClick="handleTmp()"></list-cell>
				<list-cell icon="icon-shezhi1" iconColor="#e07472" title="设置" border=""
					@eventClick="navTo('/subpages/set/set')"></list-cell>
				<list-cell icon="icon-share" iconColor="#5fcda2" title="推广二维码" border=""
					@eventClick="navTo('/subpages/qrcode/qrcode')"></list-cell>
			</view>
		</view>

		<one-click-register></one-click-register>
	</view>
</template>

<script>
	import listCell from '@/components/mix-list-cell';
	import oneClickRegister from '@/components/one-click-register.vue';
	import {
		fetchMemberCouponList
	} from '@/api/coupon.js';
	import {
		fetchTeamReward
	}
	from "@/api/team.js"
	import {
		WXResetNickName
	} from '@/api/member.js';
	import common from '@/utils/common.js'
	import {
		mapState,
		mapMutations
	} from 'vuex';
	let startY = 0,
		moveY = 0,
		pageAtTop = true;
	export default {
		components: {
			listCell,
			oneClickRegister,
		},
		data() {
			return {
				coverTransform: 'translateY(0px)',
				coverTransition: '0s',
				moving: false,
				couponCount: null,
				nickName: '',
				avatarUrl: null,
				rewardAmount: 0.0
			}
		},
		onLoad() {},
		onShow() {
			if (this.hasLogin) {
				this.fetchCounponCount()
				this.fetchTeamReward()
			} else {
				this.couponCount = null;
			}
		},
		// #ifndef MP
		onNavigationBarButtonTap(e) {
			const index = e.index;
			if (index === 0) {
				this.navTo('/subpages/set/set');
			} else if (index === 1) {
				// #ifdef APP-PLUS
				const pages = getCurrentPages();
				const page = pages[pages.length - 1];
				const currentWebview = page.$getAppWebview();
				currentWebview.hideTitleNViewButtonRedDot({
					index
				});
				// #endif
				uni.navigateTo({
					url: '/subpages/notice/notice'
				})
			}
		},
		// #endif
		computed: {
			...mapState(['hasLogin', 'userInfo', 'hadNickName']),
		},
		methods: {
			...mapMutations(['login']),
			// getNickname(e) {
			// 	this.nickName = e.detail.value
			// },
			// async handleConfirmNickName() {
			// 	let _this = this
			// 	let isUpload = false
			// 	if (_this.avatarUrl && _this.avatarUrl != '') {
			// 		uni.uploadFile({
			// 			url: common.baseUrl + "/fileUploadAndDownload/upload",
			// 			filePath: _this.avatarUrl,
			// 			name: 'file',
			// 			header: {
			// 				"x-token": _this.$store.state.token,
			// 				"x-user_id": _this.$store.state.userInfo.id,
			// 				"Access-Control-Allow-Origin": "*",
			// 				"Access-Control-Allow-Methods": "*"
			// 			},
			// 			success: res => {
			// 				const response = JSON.parse(res.data)
			// 				if (response.code == 0) {
			// 					_this.$store.state.userInfo.avatarUrl = response.data.file.url
			// 					_this.$store.state.hadNickName = true
			// 					uni.setStorage({ //缓存用户登陆状态
			// 						key: 'UserInfo',
			// 						data: _this.$store.state.userInfo
			// 					})
			// 					_this.isCloseNickNameModel = true
			// 					_this.isUpload = true
			// 					this.resetNick()
			// 				}
			// 			},
			// 			fail: (error) => {
			// 				this.$api.msg("设置失败", 2000)
			// 			}
			// 		})
			// 	}
			// },
			// async resetNick() {
			// 	let _this = this
			// 	if (_this.$store.state.userInfo && _this.isUpload) {
			// 		_this.$store.state.userInfo.nickName = _this.nickName
			// 		WXResetNickName(_this.$store.state.userInfo).then(res => {
			// 			if (res.code == 0) {
			// 				this.$api.msg("设置成功", 2000)
			// 				_this.$store.state.hadNickName = true
			// 				uni.setStorage({ //缓存用户登陆状态
			// 					key: 'UserInfo',
			// 					data: _this.$store.state.userInfo
			// 				})
			// 				_this.isCloseNickNameModel = true
			// 			} else {
			// 				this.$api.msg("设置失败", 2000)
			// 			}
			// 		});
			// 	} else {
			// 		this.$api.msg("头像设置失败", 2000)
			// 	}
			// },
			async fetchCounponCount() {
				if (this.hasLogin) {
					fetchMemberCouponList(0).then(response => {
						if (response.data != null && response.data.list.length > 0) {
							this.couponCount = response.data.list.length;
						}
					});
				} else {
					this.couponCount = null;
				}
			},
			async fetchTeamReward() {
				if (this.hasLogin) {
					fetchTeamReward().then(response => {
						console.log("-----------fetchTeamReward-------", response)
						if (response.data != null) {
							console.log("-----response---", response)
							this.rewardAmount = response.data;
						}
					});
				} else {}
			},

			goLogin() {
				this.isCloseModel = false
			},
			/**
			 * 统一跳转接口,拦截未登录路由
			 * navigator标签现在默认没有转场动画，所以用view
			 */
			navTo(url) {
				if (!this.hasLogin) {
					this.$api.msg("请先登录", 2000)
					return
				}
				uni.navigateTo({
					url
				})
			},

			/**
			 *  会员卡下拉和回弹
			 *  1.关闭bounce避免ios端下拉冲突
			 *  2.由于touchmove事件的缺陷（以前做小程序就遇到，比如20跳到40，h5反而好很多），下拉的时候会有掉帧的感觉
			 *    transition设置0.1秒延迟，让css来过渡这段空窗期
			 *  3.回弹效果可修改曲线值来调整效果，推荐一个好用的bezier生成工具 http://cubic-bezier.com/
			 */
			coverTouchstart(e) {
				if (pageAtTop === false) {
					return;
				}
				this.coverTransition = 'transform .1s linear';
				startY = e.touches[0].clientY;
			},
			coverTouchmove(e) {
				moveY = e.touches[0].clientY;
				let moveDistance = moveY - startY;
				if (moveDistance < 0) {
					this.moving = false;
					return;
				}
				this.moving = true;
				if (moveDistance >= 80 && moveDistance < 100) {
					moveDistance = 80;
				}

				if (moveDistance > 0 && moveDistance <= 80) {
					this.coverTransform = `translateY(${moveDistance}px)`;
				}
			},
			coverTouchend() {
				if (this.moving === false) {
					return;
				}
				this.moving = false;
				this.coverTransition = 'transform 0.3s cubic-bezier(.21,1.93,.53,.64)';
				this.coverTransform = 'translateY(0px)';
			},
			handleUserInfo() {
				if (!this.hadNickName) {
					this.isCloseNickNameModel = false
				}
			},
			handleChooseavatar(e) {
				this.avatarUrl = e.detail.avatarUrl;
			},
			handleTmp() {
				this.$api.msg("待开放")
			}
		}
	}
</script>
<style lang='scss'>
	%flex-center {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}

	%section {
		display: flex;
		justify-content: space-around;
		align-content: center;
		background: #fff;
		border-radius: 10upx;
	}

	.user-section {
		height: 520upx;
		padding: 100upx 30upx 0;
		position: relative;

		.bg {
			position: absolute;
			left: 0;
			top: 0;
			width: 100%;
			height: 100%;
			filter: blur(1px);
			/* opacity: .9; */
		}
	}

	.user-info-box {
		height: 180upx;
		display: flex;
		align-items: center;
		position: relative;
		z-index: 1;

		.portrait {
			width: 130upx;
			height: 130upx;
			border: 5upx solid #fff;
			border-radius: 50%;
		}

		.username {
			font-size: 24upx;
			color: $font-color-dark;
			margin-left: 20upx;
		}
	}

	.vip-card-box {
		display: flex;
		flex-direction: column;
		color: #f7d680;
		height: 240upx;
		background: linear-gradient(left, rgba(0, 0, 0, .7), rgba(0, 0, 0, .8));
		border-radius: 16upx 16upx 0 0;
		overflow: hidden;
		position: relative;
		padding: 20upx 24upx;

		.card-bg {
			position: absolute;
			top: 20upx;
			right: 0;
			width: 380upx;
			height: 260upx;
		}

		.b-btn {
			position: absolute;
			right: 20upx;
			top: 16upx;
			width: 132upx;
			height: 40upx;
			text-align: center;
			line-height: 40upx;
			font-size: 22upx;
			color: #36343c;
			border-radius: 20px;
			background: linear-gradient(left, #f9e6af, #ffd465);
			z-index: 1;
		}

		.tit {
			font-size: $font-base+2upx;
			color: #f7d680;
			margin-bottom: 28upx;

			.yticon {
				color: #f6e5a3;
				margin-right: 16upx;
			}
		}

		.e-b {
			font-size: $font-sm;
			color: #d8cba9;
			margin-top: 10upx;
		}
	}

	.cover-container {
		background: $page-color-base;
		margin-top: -150upx;
		padding: 0 30upx;
		position: relative;
		background: #f5f5f5;
		padding-bottom: 20upx;

		.arc {
			position: absolute;
			left: 0;
			top: -34upx;
			width: 100%;
			height: 36upx;
		}
	}

	.tj-sction {
		@extend %section;

		.tj-item {
			@extend %flex-center;
			flex-direction: column;
			height: 140upx;
			font-size: $font-sm;
			color: #75787d;
		}

		.num {
			font-size: $font-lg;
			color: $font-color-dark;
			margin-bottom: 8upx;
		}
	}

	.order-section {
		@extend %section;
		padding: 28upx 0;
		margin-top: 20upx;

		.order-item {
			@extend %flex-center;
			width: 120upx;
			height: 120upx;
			border-radius: 10upx;
			font-size: $font-sm;
			color: $font-color-dark;
		}

		.yticon {
			font-size: 48upx;
			margin-bottom: 18upx;
			color: #fa436a;
		}

		.icon-shouhoutuikuan {
			font-size: 44upx;
		}
	}

	.history-section {
		padding: 30upx 0 0;
		margin-top: 20upx;
		background: #fff;
		border-radius: 10upx;

		.sec-header {
			display: flex;
			align-items: center;
			font-size: $font-base;
			color: $font-color-dark;
			line-height: 40upx;
			margin-left: 30upx;

			.yticon {
				font-size: 44upx;
				color: #5eba8f;
				margin-right: 16upx;
				line-height: 40upx;
			}
		}

		.h-list {
			white-space: nowrap;
			padding: 30upx 30upx 0;

			image {
				display: inline-block;
				width: 160upx;
				height: 160upx;
				margin-right: 20upx;
				border-radius: 10upx;
			}
		}
	}
</style>