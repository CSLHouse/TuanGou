<template>
	<view class="content">
		<view class="coupon-item" v-for="(item,index) in couponList" :key="index">
			<view class="con">
				<view class="left">
					<text class="title">{{item.name}}</text>
					<text class="time">有效期至{{item.endTime | formatDateTime}}</text>
				</view>
				<view class="right">
					<text class="price">{{item.amount}}</text>
					<text>满{{item.minPoint}}可用</text>
				</view>

				<view class="circle l"></view>
				<view class="circle r"></view>
			</view>
			<text class="tips">{{item.useType | formatCouponUseType}}</text>
		</view>
		<view class="footer">
			<button>全部提取</button>
		</view>
	</view>
</template>

<script>
	import {
		fetchSettlementList
	} from '@/api/team.js';
	export default {
		data() {
			return {
				settlementList: [],

			};
		},
		onLoad() {
			this.loadData();
		},
		filters: {
			formatDateTime(time) {
				if (time == null || time === '') {
					return 'N/A';
				}
				let date = new Date(time);
				return formatDate(date, 'yyyy-MM-dd HH:mm:ss')
			},
			formatCouponUseType(useType) {
				if (useType == 0) {
					return "全场通用";
				} else if (useType == 1) {
					return "指定分类商品可用";
				} else if (useType == 2) {
					return "指定商品可用";
				}
				return null;
			},
		},
		methods: {
			loadData() {
				fetchMemberCouponList(this.useStatus).then(response => {
					this.couponList = response.data;
				});
			},
			tabClick(index) {
				this.tabCurrentIndex = index;
				this.useStatus = this.navList[index].useStatus;
				this.loadData();
			},
		}
	}
</script>

<style lang='scss'>
	page {
		background: $page-color-base;
		padding-bottom: 100upx;
	}

	/* 优惠券列表 */

	.coupon-item {

		display: flex;

		flex-direction: column;

		margin: 20upx 24upx;

		background: #fff;



		.con {

			display: flex;

			align-items: center;

			position: relative;

			height: 120upx;

			padding: 0 30upx;



			&:after {

				position: absolute;

				left: 0;

				bottom: 0;

				content: '';

				width: 100%;

				height: 0;

				border-bottom: 1px dashed #f3f3f3;

				transform: scaleY(50%);

			}

		}



		.left {

			display: flex;

			flex-direction: column;

			justify-content: center;

			flex: 1;

			overflow: hidden;

			height: 100upx;

		}



		.title {

			font-size: 32upx;

			color: $font-color-dark;

			margin-bottom: 10upx;

		}



		.time {

			font-size: 24upx;

			color: $font-color-light;

		}



		.right {

			display: flex;

			flex-direction: column;

			justify-content: center;

			align-items: center;

			font-size: 26upx;

			color: $font-color-base;

			height: 100upx;

		}



		.price {

			font-size: 44upx;

			color: $base-color;



			&:before {

				content: '￥';

				font-size: 34upx;

			}

		}



		.tips {

			font-size: 24upx;

			color: $font-color-light;

			line-height: 60upx;

			padding-left: 30upx;

		}



		.circle {

			position: absolute;

			left: -6upx;

			bottom: -10upx;

			z-index: 10;

			width: 20upx;

			height: 20upx;

			background: #f3f3f3;

			border-radius: 100px;



			&.r {

				left: auto;

				right: -6upx;

			}



		}

	}
</style>