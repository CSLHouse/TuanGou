<template>
	<view class="content">
		<view class="top-area">
			<!-- 背景图片 -->
			<view class="top-bg"></view>
			<!-- 文字内容 -->
			<view class="top-content">
				<text class="line-1">首推佣金“12%+18%循环”+2人成团</text>
			</view>
		</view>

		<view class="title-area">
			<text class="title">所在队伍</text>
		</view>
		<view v-if="joindInfos.length < 1">
			<view class="wrapper">
				<text>暂无</text>
			</view>
		</view>
		<view v-else class="list b-b">
			<view class="index">1</view>
			<view v-for="(item, index) in joindInfos" :key="index" @click="checkAddress(item)">
				<view class="wrapper">
					<view class="avatar-container">
						<image :src="item.avatarUrl" class="avatar"></image>
						<view v-if="item.isCaptain" class="badge" :style="{ backgroundColor: 'green' }">
							队</view>
						<view v-else class="badge" :style="{ backgroundColor: 'red'}">
						</view>
					</view>
					<text class="name">{{item.name}}</text>
				</view>
			</view>
		</view>

		<view class="title-area">
			<text class="title">我的团队</text>
		</view>
		<view v-if="myTeamsInfos.length < 1">
			<view class="wrapper">
				<text>暂无</text>
			</view>
		</view>
		<view v-for="(teams, index) in myTeamsInfos" :key="index">
			<view class="list b-b" @click="onShowDetails(teams)">
				<view class="index">{{ index + 1 }}</view>
				<view v-for="(item, i) in teams" :key="i" @click="checkAddress(item)">
					<view class="wrapper">
						<view class="avatar-container">
							<image :src="item.avatarUrl" class="avatar"></image>
							<!-- 角标，根据需要显示 -->
							<view v-if="item.isActivated" class="badge" :style="{ backgroundColor: 'green' }">
								团</view>
							<view v-else class="badge" :style="{ backgroundColor: 'red'}">
							</view>
						</view>
						<text class="name">{{item.name}}</text>
					</view>
				</view>
				<button class="get-btn" @click="onShowDetails(teams)">
					<uni-icons type="right" size="18"></uni-icons>
				</button>
			</view>
		</view>
	</view>
</template>

<script>
	import {
		fetchTeamList
	} from "@/api/team.js"
	export default {
		data() {
			return {
				joindInfos: [{
						avatarUrl: "http://tmp/SOXFD5VA6OOh55b1fadce8cbd79820862bcc06178395.jpeg",
						name: "借我一支烟",
						isCaptain: true,
					},
					{
						avatarUrl: "http://tmp/SOXFD5VA6OOh55b1fadce8cbd79820862bcc06178395.jpeg",
						name: "借我一支烟"
					},
				],
				myTeamsInfos: [
					[{
							avatarUrl: "http://tmp/SOXFD5VA6OOh55b1fadce8cbd79820862bcc06178395.jpeg",
							name: "借我一支烟",
							isActivated: true,
						},
						{
							avatarUrl: "http://tmp/SOXFD5VA6OOh55b1fadce8cbd79820862bcc06178395.jpeg",
							name: "借我一支烟"
						},
					],
					[{
							avatarUrl: "http://tmp/SOXFD5VA6OOh55b1fadce8cbd79820862bcc06178395.jpeg",
							name: "借"
						},
						{
							avatarUrl: "http://tmp/SOXFD5VA6OOh55b1fadce8cbd79820862bcc06178395.jpeg",
							name: "借我一支烟sadghdfhtjjg"
						},
					]
				],

			}
		},
		onLoad(option) {
			this.loadData();
		},
		methods: {
			async loadData() {
				let _this = this
				if (_this.$store.state.openId) {
					fetchTeamList({
						openId: _this.$store.state.openId
					}).then(response => {
						this.joindInfos = response.data.joinedTeam
						this.myTeamsInfos = response.data.myTeams;
					});
				} else {
					uni.showToast({
						title: '请先登录',
						duration: 2000
					})
				}
			},
			onShowDetails(teams) {
				let userIds = [];
				teams.forEach(item => {
					if (item.id) {
						userIds.push(item.id);
					}
				})
				if (userIds.length < 1) {
					this.$api.msg("无队员！", 1000)
					return;
				}
				uni.navigateTo({
					url: `/subpages/team/teamDetails?userIds=${JSON.stringify(userIds)}`
				})
			}
		}
	}
</script>

<style lang='scss'>
	page {
		padding: 20upx;
	}

	/* 顶部区域样式 */
	.top-area {
		position: relative;
		width: 100%;
		height: 100rpx;
		transform: translateX(-10px);

		.top-bg {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background-image: url('/static/user-bg.jpg');
			background-size: cover;
			background-position: center;
			background-repeat: no-repeat;
			filter: brightness(0.7);
		}

		.top-content {
			position: relative;
			height: 100%;
			display: flex;
			flex-direction: column;
			justify-content: center;
			align-items: center;
			color: #ffffff;
			text-align: center;

			.line-1 {
				font-size: 28rpx;
				font-weight: bold;
				margin-bottom: 15rpx;
				text-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.5);
			}

			/* .line-2 {
				font-size: 30rpx;
				margin-bottom: 10rpx;
				text-shadow: 0 1rpx 2rpx rgba(0, 0, 0, 0.5);
			}

			.line-3 {
				font-size: 26rpx;
				text-shadow: 0 1rpx 2rpx rgba(0, 0, 0, 0.5);
			} */
		}
	}

	/* 主内容区域样式 */
	.main-content {
		flex: 1;
		padding: 30rpx;
		background-color: #f5f5f5;

		.content-text {
			font-size: 28rpx;
			color: #333;
		}
	}

	.content {
		position: relative;
	}

	.title-area {
		display: flex;
		width: 100%;
		flex-direction: column;
		/*元素的排列方向为垂直*/
		justify-content: center;
		/*水平居中对齐*/
		align-items: center;
		/*垂直居中对齐*/
		/* 		border-bottom: 1px solid #ccc; */
		transform: translateX(-10px);
		background-color: $uni-color-primary;
	}

	.title {
		margin: 10px 0 10px 0;
		color: #ffffff;
		font-size: $font-lg;
	}

	.list {
		display: flex;
		align-items: center;
		padding: 20upx 30upx;
		background: #fff;
		position: relative;
	}

	.wrapper {
		display: flex;
		width: 120px;
		flex-direction: column;
		/*元素的排列方向为垂直*/
		justify-content: center;
		/*水平居中对齐*/
		align-items: center;
		/*垂直居中对齐*/
		/* border: 1px solid #ccc; */
	}

	/* 头像容器，用于定位角标 */
	.avatar-container {
		position: relative; // 关键：作为角标的定位容器
	}

	/* 头像样式 */
	.avatar {
		width: 80rpx;
		height: 80rpx;
		border-radius: 50%;
		object-fit: cover;
		cursor: pointer;
		display: block;

	}

	.name {
		font-size: 24upx;
		color: $font-color-light;
		margin-top: 6px;
		/* 固定宽度，根据需要调整 */
		max-width: 120px;
		/* 强制不换行 */
		white-space: nowrap;
		/* 超出部分隐藏 */
		overflow: hidden;
		/* 显示省略号 */
		text-overflow: ellipsis;
		/* 使宽度设置生效 */
		/* display: inline-block; */
	}

	.index {
		width: 40rpx;
		height: 40rpx;
		border-radius: 50%;
		background-color: $uni-color-primary;
		margin-right: 20px;

		display: flex;
		justify-content: center;
		align-items: center;

		color: #ffffff;
		font-size: 24rpx;
		font-weight: bold;
	}

	/* 角标样式 */
	.badge {
		position: absolute; // 绝对定位，相对于头像容器
		top: -10upx; // 向上偏移，部分超出头像
		right: -10upx; // 向右偏移，部分超出头像
		/* background-color: #ff4d4f; // 红色背景 */
		color: white; // 白色文字
		font-size: 20upx; // 小号字体
		width: 30upx; // 角标宽度
		height: 30upx; // 角标高度
		border-radius: 50%; // 圆形角标
		text-align: center; // 文字居中
		line-height: 30upx; // 文字垂直居中
		border: 2upx solid white; // 白色边框，与头像区分
	}

	/* “获取”按钮样式（带箭头） */
	.get-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0;
		margin: 0;
		background: transparent;
		border: none;
		/* 清除默认边框 */
		border-width: 0;
		outline: none;
		/* 清除聚焦边框 */
		font-size: $uni-font-size-base;
		color: $uni-color-primary;
		line-height: 1;
		/* 清除默认行高 */
		margin: 0 20px 0 auto;
	}

	/* 清除伪元素可能带来的边框 */
	.get-btn::after {
		border: none !important;
	}
</style>