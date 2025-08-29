<template>
	<view class="group-buy-button">
		<!-- 折叠的头像 -->
		<view class="avatar-stack">
			<view class="avatar" v-for="(avatar, index) in avatars" :key="index"
				:style="{zIndex: avatars.length - index}">
				<image :src="avatar" mode="aspectFill"></image>
			</view>
			<view class="more-avatars" v-if="maxAvatars > maxVisibleAvatars">
				<text>+{{maxAvatars - maxVisibleAvatars}}</text>
				<!-- <text>......{{maxAvatars}}</text> -->
			</view>
		</view>
	</view>
</template>

<script>
	export default {
		name: "folding-avatar",
		data() {
			return {
				// avatars: [
				// 	'https://picsum.photos/seed/user1/200/200',
				// 	'https://picsum.photos/seed/user2/200/200',
				// 	'https://picsum.photos/seed/user3/200/200',
				// 	'https://picsum.photos/seed/user4/200/200',
				// 	'https://picsum.photos/seed/user5/200/200'
				// ],
				maxVisibleAvatars: 3
			};
		},
		props: {
			avatars: {
				type: Array,
				default: () => [
					'https://picsum.photos/seed/user1/200/200',
					'https://picsum.photos/seed/user2/200/200',
					'https://picsum.photos/seed/user3/200/200',
					'https://picsum.photos/seed/user4/200/200',
					'https://picsum.photos/seed/user5/200/200'
				]
			},
			maxAvatars: {
				type: Number,
				default: 0
			}
		},
	}
</script>

<style>
	.group-buy-button {
		display: flex;
		align-items: center;
		cursor: pointer;
		transition: all 0.3s ease;
		position: relative;
		overflow: hidden;
	}

	.group-buy-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 6px 16px rgba(255, 78, 80, 0.4);
	}

	.avatar-stack {
		display: flex;
		margin-right: 12px;
		position: relative;
	}

	.avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		border: 2px solid white;
		overflow: hidden;
		margin-right: -8px;
		transition: all 0.3s ease;
		transform-origin: center;
	}

	.avatar image {
		width: 100%;
		height: 100%;
	}

	.avatar:hover {
		z-index: 10 !important;
		transform: scale(1.1);
	}

	.more-avatars {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background-color: rgba(255, 255, 255, 0.8);
		border: 2px solid white;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-right: -8px;
		font-size: 12px;
		color: #FF4E50;
		font-weight: bold;
	}

	.button-text {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
	}

	.main-text {
		color: white;
		font-size: 16px;
		font-weight: bold;
		text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.2);
	}

	.sub-text {
		color: white;
		font-size: 12px;
		opacity: 0.9;
	}
</style>