<template>
	<view class="container">
		<!-- 二维码组件（默认隐藏，仅用于生成） -->
		<uqrcode ref="uqrcodeRef" canvas-id="auto-qrcode" :value="qrContent" :size="300" :options="qrOptions"
			:auto="true" @complete="handleComplete"></uqrcode>

		<!-- 状态提示 -->
		<view class="status">{{ statusText }}</view>
		<button @click="exportAndDownload" class="download-btn" :disabled="!canvasReady">下载二维码</button>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				// 二维码内容（可动态修改）
				qrContent: '/pages/index/index?refCode=' + wx.getStorageSync("OpenId"),
				// 二维码样式配置
				qrOptions: {
					color: "#333333",
					backgroundColor: "#ffffff",
					margin: 10
					// 如需添加logo，使用绝对路径
					// logo: { imageSrc: "/static/logo.png" }
				},
				// 状态文本
				statusText: "正在生成二维码...",
				canvasReady: false, // 标记canvas是否可用
			};
		},
		onReady() {
			// 页面就绪后自动触发生成（如需手动触发可移至按钮事件）
			this.$nextTick(() => {
				this.$refs.uqrcodeRef?.remake();
			});
		},
		methods: {
			// 二维码生成完成回调
			handleComplete(e) {
				if (!e.success) {
					this.statusText = `生成失败：${e.error}`;
					return;
				}
				this.statusText = "生成成功，正在准备下载...";
				this.canvasReady = true;
				// 生成成功后立即导出临时文件并下载
				// this.exportAndDownload();
			},

			// 导出临时文件并下载
			exportAndDownload() {
				// 调用组件方法获取临时文件路径
				this.$refs.uqrcodeRef.toTempFilePath({
					success: (res) => {
						const tempFilePath = res.tempFilePath;
						// 根据平台执行下载
						this.downloadByPlatform(tempFilePath);
					},
					fail: (err) => {
						this.statusText = `导出失败：${err.errMsg}`;
					}
				});
			},

			// 按平台处理下载逻辑
			downloadByPlatform(tempFilePath) {
				// H5平台：通过a标签下载
				// #ifdef H5
				uni.showToast({
					title: "H5请长按图片保存",
					icon: "none",
					duration: 3000
				});
				// 创建临时图片显示（供用户长按保存）
				const img = new Image();
				img.src = tempFilePath;
				img.onload = () => {
					const a = document.createElement("a");
					a.href = tempFilePath;
					a.download = "qrcode.png";
					document.body.appendChild(a);
					a.click();
					document.body.removeChild(a);
					this.statusText = "已触发下载（部分浏览器可能需要手动保存）";
				};
				// #endif

				// 微信小程序/APP等平台：保存到相册
				// #ifndef H5
				this.$refs.uqrcodeRef.save({
					filePath: tempFilePath,
					success: () => {
						this.statusText = "二维码已保存到相册";
						uni.showToast({
							title: "保存成功",
							icon: "success"
						});
					},
					fail: (err) => {
						this.statusText = `保存失败：${err.errMsg}`;
						// 处理授权失败
						if (err.errMsg.includes("auth deny")) {
							uni.showModal({
								title: "授权提示",
								content: "需要获取相册权限才能保存图片",
								success: (res) => {
									if (res.confirm) {
										uni.openSetting();
									}
								}
							});
						}
					}
				});
				// #endif
			}
		}
	};
</script>

<style scoped>
	.container {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 50rpx;
	}

	.status {
		margin-top: 30rpx;
		font-size: 28rpx;
		color: #666;
	}

	.download-btn {
		margin-top: 30rpx;
		background-color: #fa436a;
		color: white;
	}

	.download-btn:disabled {
		background-color: #ccc;
		cursor: not-allowed;
	}
</style>