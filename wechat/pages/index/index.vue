<template>
	<view class="container">
		<!-- 小程序头部兼容 -->
		<!-- #ifdef MP -->
		<!-- <view class="mp-search-box">
			<input class="ser-input" type="text" value="输入关键字搜索" disabled />
		</view> -->
		<!-- #endif -->

		<!-- 头部轮播 -->
		<view class="carousel-section">
			<!-- 标题栏和状态栏占位符 -->
			<view class="titleNview-placing"></view>
			<!-- 背景色区域 -->
			<view class="titleNview-background"></view>
			<swiper class="carousel" circular indicator-dots autoplay :interval="5000" :duration="500">
				<swiper-item v-for="(item, index) in advertiseList" :key="index" class="carousel-item"
					@click="navToAdvertisePage(item)">
					<image :src="item.pic" />
				</swiper-item>
			</swiper>
			<!-- 自定义swiper指示器 -->
			<view class="swiper-dots">
				<text class="num">{{swiperCurrent+1}}</text>
				<text class="sign">/</text>
				<text class="num">{{swiperLength}}</text>
			</view>
		</view>
		<!-- 头部功能区 -->
		<!-- <view class="cate-section">
			<view class="cate-item">
				<image src="/static/temp/c3.png"></image>
				<text>专题</text>
			</view>
			<view class="cate-item">
				<image src="/static/temp/c5.png"></image>
				<text>话题</text>
			</view>
			<view class="cate-item">
				<image src="/static/temp/c6.png"></image>
				<text>优选</text>
			</view>
			<view class="cate-item">
				<image src="/static/temp/c7.png"></image>
				<text>特惠</text>
			</view>
		</view> -->

		<!-- 品牌制造商直供 -->
		<view class="f-header m-t" @click="navToRecommendBrandPage()">
			<image src="/static/icon_home_brand.png"></image>
			<view class="tit-box">
				<text class="tit">品牌制造商直供</text>
				<text class="tit2">工厂直达消费者，剔除品牌溢价</text>
			</view>
		</view>
		<view class="guess-section">
			<view v-for="(item, index) in brandList" :key="index" class="guess-item"
				@click="navToBrandDetailPage(item)">
				<view class="image-wrapper-brand">
					<image :src="item.logo" mode="aspectFit"></image>
				</view>
				<text class="title clamp">{{item.name}}</text>
				<text class="title2">商品数量：{{item.productCount}}</text>
			</view>
		</view>

		<!-- 秒杀专区 -->
		<view class="f-header m-t" v-if="homeFlashPromotion!==null">
			<image src="/static/icon_flash_promotion.png"></image>
			<view class="tit-box">
				<text class="tit">秒杀专区</text>
				<text class="tit2">下一场 {{homeFlashPromotion.nextStartTime | formatTime}} 开始</text>
			</view>
			<view class="tit-box">
				<text class="tit2" style="text-align: right;">本场结束剩余：</text>
				<view style="text-align: right;">
					<text class="hour timer">{{cutDownTime.endHour}}</text>
					<text>:</text>
					<text class="minute timer">{{cutDownTime.endMinute}}</text>
					<text>:</text>
					<text class="second timer">{{cutDownTime.endSecond}}</text>
				</view>
			</view>
			<text class="yticon icon-you" v-show="false"></text>
		</view>

		<view class="guess-section">
			<view v-for="(item, index) in homeFlashPromotion.productList" :key="index" class="guess-item"
				@click="navToDetailPage(item)">
				<view class="image-wrapper">
					<image :src="item.pic" mode="aspectFill"></image>
				</view>
				<text class="title clamp">{{item.name}}</text>
				<text class="title2 clamp">{{item.subTitle}}</text>
				<text class="price">￥{{item.price}}</text>
			</view>
		</view>

		<!-- 团购楼层 -->
		<view class="f-header m-t" v-if="groupBuyProducts.groups != null && groupBuyProducts.groups.length > 0">
			<image src="/static/temp/h1.png"></image>
			<view class="tit-box">
				<text class="tit">精品团购</text>
			</view>
			<text class="yticon icon-you"></text>
		</view>
		<view class="group-section" v-if="groupBuyProducts.groups != null && groupBuyProducts.groups.length > 0">
			<swiper class="g-swiper" :duration="500">
				<swiper-item class="g-swiper-item" v-for="(item, index) in groupBuyProducts.groups" :key="index">
					<view class="g-item left" @click="navToDetailPageById(item.productId)">
						<image :src="item.pic" mode="aspectFill"></image>
						<view class="t-box">
							<text class="title clamp">{{item.name}}</text>
							<view class="price-box">
								<text class="price">￥{{item.price}}</text>
								<text class="m-price">￥{{item.originalPrice}}</text>
							</view>

							<view class="pro-box">
								<view class="progress-box">
									<progress :percent="item.percent" activeColor="#fa436a" active stroke-width="6" />
								</view>
								<text>火热团购</text>
							</view>
						</view>
					</view>
					<view class="g-item right" @click="navToDetailPageById(item.productId)">
						<image :src="item.pic" mode="aspectFill"></image>
						<view class="t-box">
							<text class="title clamp">{{item.name}}</text>
							<view class="price-box">
								<text class="price">￥{{item.price}}</text>
								<text class="m-price">￥{{item.originalPrice}}</text>
							</view>
							<view class="pro-box">
								<view class="progress-box">
									<progress :percent="item.percent" activeColor="#fa436a" active stroke-width="6" />
								</view>
								<text>热门推荐</text>
							</view>
						</view>
					</view>
				</swiper-item>

			</swiper>
		</view>

		<!-- 新鲜好物 -->
		<view class="f-header m-t" @click="navToNewProudctListPage()">
			<image src="/static/icon_new_product.png"></image>
			<view class="tit-box">
				<text class="tit">新鲜好物</text>
				<text class="tit2">为你寻觅世间好物</text>
			</view>
			<text class="yticon icon-you"></text>
		</view>
		<view class="seckill-section">
			<scroll-view class="floor-list" scroll-x>
				<view class="scoll-wrapper">
					<view v-for="(item, index) in newProductList" :key="index" class="floor-item"
						@click="navToDetailPage(item)">
						<image :src="item.product.pic" mode="aspectFill"></image>
						<text class="title clamp">{{item.product.name}}</text>
						<text class="title2 clamp">{{item.product.subTitle}}</text>
						<text class="price">￥{{item.product.price}}</text>
					</view>
				</view>
			</scroll-view>
		</view>

		<!-- 人气推荐楼层 -->
		<view class="f-header m-t" @click="navToHotProudctListPage()">
			<image src="/static/icon_hot_product.png"></image>
			<view class="tit-box">
				<text class="tit">人气推荐</text>
				<text class="tit2">大家都赞不绝口的</text>
			</view>
			<text class="yticon icon-you"></text>
		</view>

		<view class="hot-section">
			<view v-for="(item, index) in hotProductList" :key="index" class="guess-item"
				@click="navToDetailPage(item)">
				<view class="image-wrapper">
					<image :src="item.product.pic" mode="aspectFill"></image>
				</view>
				<view class="txt">
					<text class="title clamp">{{item.product.name}}</text>
					<text class="title2">{{item.product.subTitle}}</text>
					<text class="price">￥{{item.product.price}}</text>
				</view>
			</view>
		</view>

		<!-- 猜你喜欢-->
		<view class="f-header m-t">
			<image src="/static/icon_recommend_product.png"></image>
			<view class="tit-box">
				<text class="tit">猜你喜欢</text>
				<text class="tit2">你喜欢的都在这里了</text>
			</view>
			<text class="yticon icon-you" v-show="false"></text>
		</view>

		<view class="guess-section">
			<view v-for="(item, index) in recommendProductList" :key="index" class="guess-item"
				@click="navToDetailPage(item)">
				<view class="image-wrapper">
					<image :src="item.pic" mode="aspectFill"></image>
				</view>
				<text class="title clamp">{{item.name}}</text>
				<text class="title2 clamp">{{item.subTitle}}</text>
				<text class="price">￥{{item.price}}</text>
			</view>
		</view>
		<uni-load-more :status="loadingType"></uni-load-more>

		<!-- <view v-if='!hasLogin && !isCloseModel'>
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
		</view> -->
		<one-click-register></one-click-register>

	</view>
</template>

<script>
	import {
		mapState,
		mapMutations
	} from 'vuex';
	import {
		fetchContent,
		fetchRecommendProductList,
		recordShareCount
	} from '@/api/home.js';
	import {
		getWXPhoneNumber,
		wxRefreshLogin,
		WXResetNickName
	} from '@/api/member.js';
	import {
		formatDate
	} from '@/utils/date';
	import uniLoadMore from '@/components/uni-load-more/uni-load-more.vue';
	import oneClickRegister from '@/components/one-click-register.vue';
	export default {
		components: {
			uniLoadMore,
			oneClickRegister,
		},
		data() {
			return {
				titleNViewBackground: '',
				swiperCurrent: 0,
				swiperLength: 0,
				carouselList: [],
				goodsList: [],
				advertiseList: [],
				brandList: [],
				homeFlashPromotion: [],
				newProductList: [],
				hotProductList: [],
				recommendProductList: [],
				recommendParams: {
					page: 1,
					pageSize: 4
				},
				loadingType: 'more',
				isCloseModel: false,
				isCloseNickNameModel: false,
				nickName: '',
				groupBuyProducts: [],
			};
		},
		onLoad(options) {
			if (options.refCode && options.refCode.length > 0) {
				console.log("--------refCode------", options.refCode)
				this.$store.state.formOpenId = options.refCode
			}
			this.loadData(options);
		},
		onShow() {},
		//下拉刷新
		onPullDownRefresh() {
			this.recommendParams.page = 1;
			this.loadData();
		},
		//加载更多
		onReachBottom() {
			this.recommendParams.page++;
			this.loadingType = 'loading';
			fetchRecommendProductList(this.recommendParams).then(response => {
				let addProductList = response.data.list;
				if (!addProductList) {
					//没有更多了
					this.recommendParams.page;
					this.loadingType = 'nomore';
				} else {
					this.recommendProductList = this.recommendProductList.concat(addProductList);
					this.loadingType = 'more';
				}
			})
		},
		computed: {
			...mapState(['hasLogin', 'userInfo']),
			cutDownTime() {
				let endTime = new Date(this.homeFlashPromotion.endTime);
				let endDateTime = new Date();
				let startDateTime = new Date();
				endDateTime.setHours(endTime.getHours());
				endDateTime.setMinutes(endTime.getMinutes());
				endDateTime.setSeconds(endTime.getSeconds());
				let offsetTime = (endDateTime.getTime() - startDateTime.getTime());
				if (offsetTime < 0) {
					offsetTime = 0
				}
				let endHour = Math.floor(offsetTime / (60 * 60 * 1000));
				let offsetMinute = offsetTime % (60 * 60 * 1000);
				let endMinute = Math.floor(offsetMinute / (60 * 1000));
				let offsetSecond = offsetTime % (60 * 1000);
				let endSecond = Math.floor(offsetSecond / 1000);
				return {
					endHour: endHour,
					endMinute: endMinute,
					endSecond: endSecond
				}
			}
		},
		filters: {
			formatTime(time) {
				if (time == null || time === '0001-01-01T00:00:00Z') {
					return 'N/A';
				}
				let date = new Date(time);
				return formatDate(date, 'HH:mm:ss')
			},
		},
		methods: {
			...mapMutations(['login']),
			// getNickname(e) {
			// 	this.nickName = e.detail.value
			// },
			// checkNickName() {
			// 	if (!this.nickName) {
			// 		this.$api.msg('请输入昵称')
			// 		return false
			// 	}
			// 	let str = this.nickName.trim();
			// 	if (str.length == 0) {
			// 		this.$api.msg('请输入正确的昵称')
			// 		return false
			// 	}
			// 	this.nickName = str
			// 	return true
			// },

			// confirmNickName() {
			// 	let _this = this
			// 	if (this.$store.state.userInfo) {
			// 		_this.$store.state.userInfo.nickName = this.nickName
			// 		WXResetNickName(this.$store.state.userInfo).then(res => {
			// 			if (res.code == 0) {
			// 				_this.$store.state.hadNickName = true
			// 				uni.setStorage({ //缓存用户登陆状态
			// 					key: 'UserInfo',
			// 					data: _this.$store.state.userInfo
			// 				})
			// 				_this.isCloseNickNameModel = true
			// 				this.$api.msg('设置成功')
			// 			} else {
			// 				this.$api.msg('设置失败')
			// 			}
			// 		});
			// 	}

			// },
			// closePop() {
			// 	this.isCloseModel = true
			// },
			// decryptPhoneNumber: function(e) {
			// 	let _this = this
			// 	console.log("--------decryptPhoneNumber---------", e.detail)
			// 	if (e.detail.errMsg == "getPhoneNumber:ok") {
			// 		if (_this.$store.state.openId && _this.$store.state.openId.length > 0) {
			// 			getWXPhoneNumber({
			// 				openId: _this.$store.state.openId,
			// 				code: e.detail.code
			// 			}).then(res => {
			// 				if (res.code == 0) {
			// 					_this.getToken()
			// 					this.$api.msg('注册成功')
			// 				} else {
			// 					this.$api.msg('注册会员失败')
			// 				}
			// 			});
			// 		}
			// 	}
			// },
			// getToken() {
			// 	let _this = this
			// 	wxRefreshLogin({
			// 		openId: _this.$store.state.openId
			// 	}).then(res => {
			// 		if (res.code == 0) {
			// 			const userinfo = res.data
			// 			wx.setStorageSync("Token", userinfo.token)
			// 			wx.setStorageSync("TokenTime", userinfo.expiresAt)
			// 			_this.$store.state.token = userinfo.token
			// 			this.login(userinfo.customer);
			// 		}
			// 	}).catch(errors => {
			// 		uni.showModal({
			// 			title: '提示',
			// 			content: '网络错误',
			// 			showCancel: false
			// 		})
			// 	});
			// },
			/**
			 * 加载数据
			 */
			async loadData() {
				fetchContent().then(response => {
					console.log("---fetchContent-----", response.data)
					this.advertiseList = response.data.advertiseList;
					this.swiperLength = this.advertiseList.length;
					this.brandList = response.data.brandList;
					this.homeFlashPromotion = response.data.homeFlashPromotion;
					this.newProductList = response.data.newProductList;
					this.hotProductList = response.data.hotProductList;
					this.groupBuyProducts = response.data.groupBuy
					fetchRecommendProductList(this.recommendParams).then(response => {
						this.recommendProductList = response.data.list;
						uni.stopPullDownRefresh();
					})
				});
			},

			//商品详情页
			navToDetailPage(item) {
				let id = 0
				if (item.product && item.product.id > 0) {
					id = item.product.id;
				} else {
					id = item.id
				}
				if (id < 1) {
					this.$api.msg('参数错误');
					return
				}
				uni.navigateTo({
					url: `/subpages/product/product?id=${id}`
				})
			},
			navToDetailPageById(id) {
				uni.navigateTo({
					url: `/subpages/product/product?id=${id}`
				})
			},
			//广告详情页
			navToAdvertisePage(item) {
				uni.navigateTo({
					url: '' + item.url
				})
			},
			//品牌详情页
			navToBrandDetailPage(item) {
				let id = item.id;
				uni.navigateTo({
					url: `/subpages/brand/brandDetail?id=${id}`
				})
			},
			//推荐品牌列表页
			navToRecommendBrandPage() {
				uni.navigateTo({
					url: `/subpages/brand/list`
				})
			},
			//新鲜好物列表页
			navToNewProudctListPage() {
				uni.navigateTo({
					url: `/subpages/product/newProductList`
				})
			},
			//人气推荐列表页
			// navToHotProudctListPage() {
			// 	uni.navigateTo({
			// 		url: `/subpages/product/hotProductList`
			// 	})
			// },
		},
		// #ifndef MP
		// 标题栏input搜索框点击
		onNavigationBarSearchInputClicked: async function(e) {
			this.$api.msg('点击了搜索框');
		},
		//点击导航栏 buttons 时触发
		onNavigationBarButtonTap(e) {
			const index = e.index;
			if (index === 0) {
				this.$api.msg('点击了扫描');
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
		}
		// #endif
	}
</script>

<style lang="scss">
	/* #ifdef MP */
	.mp-search-box {
		position: absolute;
		left: 0;
		top: 30upx;
		z-index: 9999;
		width: 100%;
		padding: 0 80upx;

		.ser-input {
			flex: 1;
			height: 56upx;
			line-height: 56upx;
			text-align: center;
			font-size: 28upx;
			color: $font-color-base;
			border-radius: 20px;
			background: rgba(255, 255, 255, .6);
		}
	}

	page {
		.cate-section {
			position: relative;
			z-index: 5;
			border-radius: 16upx 16upx 0 0;
			margin-top: -20upx;
		}

		.carousel-section {
			padding: 0;

			.titleNview-placing {
				padding-top: 0;
				height: 0;
			}

			.carousel {
				.carousel-item {
					padding: 0;
				}
			}

			.swiper-dots {
				left: 45upx;
				bottom: 40upx;
			}
		}
	}

	/* #endif */


	page {
		background: #f5f5f5;
	}

	.m-t {
		margin-top: 16upx;
	}

	/* 头部 轮播图 */
	.carousel-section {
		position: relative;
		padding-top: 10px;
		margin-bottom: 10px;

		.titleNview-placing {
			height: var(--status-bar-height);
			padding-top: 44px;
			box-sizing: content-box;
		}

		.titleNview-background {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 426upx;
			transition: .4s;
		}
	}

	.carousel {
		width: 100%;
		height: 350upx;

		.carousel-item {
			width: 100%;
			height: 100%;
			padding: 0 28upx;
			overflow: hidden;
		}

		image {
			width: 100%;
			height: 100%;
			border-radius: 10upx;
		}
	}

	.swiper-dots {
		display: flex;
		position: absolute;
		left: 60upx;
		bottom: 15upx;
		width: 72upx;
		height: 36upx;
		background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAMgAAABkCAYAAADDhn8LAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyZpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuNi1jMTMyIDc5LjE1OTI4NCwgMjAxNi8wNC8xOS0xMzoxMzo0MCAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wTU09Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9tbS8iIHhtbG5zOnN0UmVmPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvc1R5cGUvUmVzb3VyY2VSZWYjIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtcE1NOkRvY3VtZW50SUQ9InhtcC5kaWQ6OTk4MzlBNjE0NjU1MTFFOUExNjRFQ0I3RTQ0NEExQjMiIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6OTk4MzlBNjA0NjU1MTFFOUExNjRFQ0I3RTQ0NEExQjMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTcgKFdpbmRvd3MpIj4gPHhtcE1NOkRlcml2ZWRGcm9tIHN0UmVmOmluc3RhbmNlSUQ9InhtcC5paWQ6Q0E3RUNERkE0NjExMTFFOTg5NzI4MTM2Rjg0OUQwOEUiIHN0UmVmOmRvY3VtZW50SUQ9InhtcC5kaWQ6Q0E3RUNERkI0NjExMTFFOTg5NzI4MTM2Rjg0OUQwOEUiLz4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz4Gh5BPAAACTUlEQVR42uzcQW7jQAwFUdN306l1uWwNww5kqdsmm6/2MwtVCp8CosQtP9vg/2+/gY+DRAMBgqnjIp2PaCxCLLldpPARRIiFj1yBbMV+cHZh9PURRLQNhY8kgWyL/WDtwujjI8hoE8rKLqb5CDJaRMJHokC6yKgSCR9JAukmokIknCQJpLOIrJFwMsBJELFcKHwM9BFkLBMKFxNcBCHlQ+FhoocgpVwwnv0Xn30QBJGMC0QcaBVJiAMiec/dcwKuL4j1QMsVCXFAJE4s4NQA3K/8Y6DzO4g40P7UcmIBJxbEesCKWBDg8wWxHrAiFgT4fEGsB/CwIhYE+AeBAAdPLOcV8HRmWRDAiQVcO7GcV8CLM8uCAE4sQCDAlHcQ7x+ABQEEAggEEAggEEAggEAAgQACASAQQCCAQACBAAIBBAIIBBAIIBBAIABe4e9iAe/xd7EAJxYgEGDeO4j3EODp/cOCAE4sYMyJ5cwCHs4rCwI4sYBxJ5YzC84rCwKcXxArAuthQYDzC2JF0H49LAhwYUGsCFqvx5EF2T07dMaJBetx4cRyaqFtHJ8EIhK0i8OJBQxcECuCVutxJhCRoE0cZwMRyRcFefa/ffZBVPogePihhyCnbBhcfMFFEFM+DD4m+ghSlgmDkwlOgpAl4+BkkJMgZdk4+EgaSCcpVX7bmY9kgXQQU+1TgE0c+QJZUUz1b2T4SBbIKmJW+3iMj2SBVBWz+leVfCQLpIqYbp8b85EskIxyfIOfK5Sf+wiCRJEsllQ+oqEkQfBxmD8BBgA5hVjXyrBNUQAAAABJRU5ErkJggg==);
		background-size: 100% 100%;

		.num {
			width: 36upx;
			height: 36upx;
			border-radius: 50px;
			font-size: 24upx;
			color: #fff;
			text-align: center;
			line-height: 36upx;
		}

		.sign {
			position: absolute;
			top: 0;
			left: 50%;
			line-height: 36upx;
			font-size: 12upx;
			color: #fff;
			transform: translateX(-50%);
		}
	}

	/* 分类 */
	.cate-section {
		display: flex;
		justify-content: space-around;
		align-items: center;
		flex-wrap: wrap;
		padding: 30upx 22upx;
		background: #fff;

		.cate-item {
			display: flex;
			flex-direction: column;
			align-items: center;
			font-size: $font-sm + 2upx;
			color: $font-color-dark;
		}

		/* 原图标颜色太深,不想改图了,所以加了透明度 */
		image {
			width: 88upx;
			height: 88upx;
			margin-bottom: 14upx;
			border-radius: 50%;
			opacity: .7;
			box-shadow: 4upx 4upx 20upx rgba(250, 67, 106, 0.3);
		}
	}

	.ad-1 {
		width: 100%;
		height: 210upx;
		padding: 10upx 0;
		background: #fff;

		image {
			width: 100%;
			height: 100%;
		}
	}

	/* 秒杀专区 */
	.seckill-section {
		padding: 4upx 30upx 24upx;
		background: #fff;

		.s-header {
			display: flex;
			align-items: center;
			height: 92upx;
			line-height: 1;

			.s-img {
				width: 140upx;
				height: 30upx;
			}

			.tip {
				font-size: $font-base;
				color: $font-color-light;
				margin: 0 20upx 0 40upx;
			}

			.timer {
				display: inline-block;
				width: 40upx;
				height: 36upx;
				text-align: center;
				line-height: 36upx;
				margin-right: 14upx;
				font-size: $font-sm+2upx;
				color: #fff;
				border-radius: 2px;
				background: rgba(0, 0, 0, .8);
			}

			.icon-you {
				font-size: $font-lg;
				color: $font-color-light;
				flex: 1;
				text-align: right;
			}
		}

		.floor-list {
			white-space: nowrap;
		}

		.scoll-wrapper {
			display: flex;
			align-items: flex-start;
		}

		.floor-item {
			width: 300upx;
			margin-right: 20upx;
			font-size: $font-sm+2upx;
			color: $font-color-dark;
			line-height: 1.8;

			image {
				width: 300upx;
				height: 300upx;
				border-radius: 6upx;
			}

			.price {
				color: $uni-color-primary;
			}
		}

		.title2 {
			font-size: $font-sm;
			color: $font-color-light;
			line-height: 40upx;
		}
	}

	.f-header {
		display: flex;
		align-items: center;
		height: 140upx;
		padding: 6upx 30upx 8upx;
		background: #fff;

		image {
			flex-shrink: 0;
			width: 80upx;
			height: 80upx;
			margin-right: 20upx;
		}

		.tit-box {
			flex: 1;
			display: flex;
			flex-direction: column;
		}

		.tit {
			font-size: $font-lg +2upx;
			color: #font-color-dark;
			line-height: 1.3;
		}

		.tit2 {
			font-size: $font-sm;
			color: $font-color-light;
		}

		.icon-you {
			font-size: $font-lg +2upx;
			color: $font-color-light;
		}

		.timer {
			display: inline-block;
			width: 40upx;
			height: 36upx;
			text-align: center;
			line-height: 36upx;
			margin-right: 14upx;
			font-size: $font-sm+2upx;
			color: #fff;
			border-radius: 2px;
			background: rgba(0, 0, 0, .8);
		}
	}

	/* 团购楼层 */
	.group-section {
		background: #fff;

		.g-swiper {
			height: 650upx;
			padding-bottom: 30upx;
		}

		.g-swiper-item {
			width: 100%;
			padding: 0 30upx;
			display: flex;
		}

		image {
			width: 100%;
			height: 460upx;
			border-radius: 4px;
		}

		.g-item {
			display: flex;
			flex-direction: column;
			overflow: hidden;
		}

		.left {
			flex: 1.2;
			margin-right: 24upx;

			.t-box {
				padding-top: 20upx;
			}
		}

		.right {
			flex: 0.8;
			flex-direction: column-reverse;

			.t-box {
				padding-bottom: 20upx;
			}
		}

		.t-box {
			height: 160upx;
			font-size: $font-base+2upx;
			color: $font-color-dark;
			line-height: 1.6;
		}

		.price {
			color: $uni-color-primary;
		}

		.m-price {
			font-size: $font-sm+2upx;
			text-decoration: line-through;
			color: $font-color-light;
			margin-left: 8upx;
		}

		.pro-box {
			display: flex;
			align-items: center;
			margin-top: 10upx;
			font-size: $font-sm;
			color: $font-base;
			padding-right: 10upx;
		}

		.progress-box {
			flex: 1;
			border-radius: 10px;
			overflow: hidden;
			margin-right: 8upx;
		}
	}

	/* 分类推荐楼层 */
	.hot-floor {
		width: 100%;
		overflow: hidden;
		margin-bottom: 20upx;

		.floor-img-box {
			width: 100%;
			height: 320upx;
			position: relative;

			&:after {
				content: '';
				position: absolute;
				left: 0;
				top: 0;
				width: 100%;
				height: 100%;
				background: linear-gradient(rgba(255, 255, 255, .06) 30%, #f8f8f8);
			}
		}

		.floor-img {
			width: 100%;
			height: 100%;
		}

		.floor-list {
			white-space: nowrap;
			padding: 20upx;
			padding-right: 50upx;
			border-radius: 6upx;
			margin-top: -140upx;
			margin-left: 30upx;
			background: #fff;
			box-shadow: 1px 1px 5px rgba(0, 0, 0, .2);
			position: relative;
			z-index: 1;
		}

		.scoll-wrapper {
			display: flex;
			align-items: flex-start;
		}

		.floor-item {
			width: 180upx;
			margin-right: 20upx;
			font-size: $font-sm+2upx;
			color: $font-color-dark;
			line-height: 1.8;

			image {
				width: 180upx;
				height: 180upx;
				border-radius: 6upx;
			}

			.price {
				color: $uni-color-primary;
			}
		}

		.more {
			display: flex;
			align-items: center;
			justify-content: center;
			flex-direction: column;
			flex-shrink: 0;
			width: 180upx;
			height: 180upx;
			border-radius: 6upx;
			background: #f3f3f3;
			font-size: $font-base;
			color: $font-color-light;

			text:first-child {
				margin-bottom: 4upx;
			}
		}
	}

	/* 猜你喜欢 */
	.guess-section {
		display: flex;
		flex-wrap: wrap;
		padding: 0 30upx;
		background: #fff;

		.guess-item {
			display: flex;
			flex-direction: column;
			width: 48%;
			padding-bottom: 40upx;

			&:nth-child(2n+1) {
				margin-right: 4%;
			}
		}

		.image-wrapper {
			width: 100%;
			height: 330upx;
			border-radius: 3px;
			overflow: hidden;

			image {
				width: 100%;
				height: 100%;
				opacity: 1;
			}
		}

		.image-wrapper-brand {
			width: 100%;
			height: 150upx;
			border-radius: 3px;
			overflow: hidden;
			border-style: solid;
			border-color: rgba(250, 250, 255, 0.9);

			image {
				width: 100%;
				height: 100%;
				opacity: 1;
			}
		}

		.title {
			font-size: $font-lg;
			color: $font-color-dark;
			line-height: 80upx;
		}

		.title2 {
			font-size: $font-sm;
			color: $font-color-light;
			line-height: 40upx;
		}

		.price {
			font-size: $font-lg;
			color: $uni-color-primary;
			line-height: 1;
		}
	}

	.hot-section {
		display: flex;
		flex-wrap: wrap;
		padding: 0 30upx;
		background: #fff;

		.guess-item {
			display: flex;
			flex-direction: row;
			width: 100%;
			padding-bottom: 40upx;
		}

		.image-wrapper {
			width: 30%;
			height: 250upx;
			border-radius: 3px;
			overflow: hidden;

			image {
				width: 100%;
				height: 100%;
				opacity: 1;
			}
		}

		.title {
			font-size: $font-lg;
			color: $font-color-dark;
			line-height: 80upx;
		}

		.title2 {
			font-size: $font-sm;
			color: $font-color-light;
			line-height: 40upx;
			height: 80upx;
			overflow: hidden;
			text-overflow: ellipsis;
			display: block;
		}

		.price {
			font-size: $font-lg;
			color: $uni-color-primary;
			line-height: 80upx;
		}

		.txt {
			width: 70%;
			display: flex;
			flex-direction: column;
			padding-left: 40upx;
		}
	}

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

	.weui-input {
		margin-top: 40px;
		// width: 200px;
		height: 40px;
		background: #f4f4f6;
		line-height: 40px;
		text-align: center;
	}
</style>