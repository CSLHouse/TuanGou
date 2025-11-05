<template>
	<view class="page">
		<br />
		<text class="bill-id">-　订单号:　{{ orderItemId }} </text>

		<view class="feedback-section">
			<view class="feedback-title">
				<text>问题和意见</text>
				<text class="feedback-quick" @tap="chooseMsg">快速键入</text>
			</view>
			<view class="feedback-body">
				<textarea placeholder="请详细描述你的问题和意见..." v-model="sendData.content" class="feedback-textarea" autosize
					maxlength="500"></textarea>
				<text class="word-count">{{ sendData.content.length }}/500</text>
			</view>
		</view>

		<view class="feedback-section">
			<view class="feedback-title">
				<text>图片(选填,提供问题截图,总大小10M以下)</text>
			</view>
			<view class="feedback-body feedback-uploader">
				<view class="uni-uploader">
					<view class="uni-uploader-head">
						<view class="uni-uploader-title">点击预览图片</view>
						<view class="uni-uploader-info">{{ imageList.length }}/5</view>
					</view>
					<view class="uni-uploader-body">
						<view class="uni-uploader__files">
							<block v-for="(image, index) in imageList" :key="index">
								<view class="uni-uploader__file" style="position: relative;">
									<image class="uni-uploader__img" :src="image" @tap="previewImage(index)"
										mode="aspectFill" lazy-load></image>
									<view class="close-view" @click.stop="close(index)">×</view>
								</view>
							</block>
							<view class="uni-uploader__input-box" v-show="imageList.length < 5">
								<view class="uni-uploader__input" @tap="chooseImg"></view>
							</view>
						</view>
					</view>
				</view>
			</view>
		</view>

		<view class="feedback-section">
			<view class="feedback-title"><text>手机/QQ/邮箱</text></view>
			<view class="feedback-body">
				<input class="feedback-input" v-model="sendData.contact" placeholder="(请填写您的手机、QQ或E-mail,方便我们联系您 )"
					type="text" maxlength="50" />
			</view>
		</view>

		<button type="default" class="feedback-submit" @tap="send" :disabled="isSubmitting" :loading="isSubmitting">
			{{ isSubmitting ? '提交中...' : '提交' }}
		</button>
		<view class="feedback-note"><text>您反馈的信息我们将以最快速度处理</text></view>
	</view>
</template>

<script>
	import {
		UploadFileWx,
		UploadFile,
		DealOrder
	} from '@/api/order.js';
	export default {
		data() {
			return {
				orderItemId: "",
				index: 0,
				msgContents: [],
				imageList: [],
				sendData: {
					content: '', // 反馈内容
					contact: '' // 联系方式
				},
				isSubmitting: false // 提交状态控制
			};
		},
		onLoad(options) {
			this.orderItemId = options.orderItemId || '';
			this.index = parseInt(options.index || 0);

			// 初始化快速输入选项
			this.msgContents = [
				`订单号：${this.orderItemId}`,
				"申请退款",
				"订单损坏",
				"损坏包赔",
				"物流问题",
				"商品与描述不符"
			];
		},
		methods: {
			/**
			 * 移除图片
			 * @param {Number} index - 图片索引
			 */
			close(index) {
				if (index >= 0 && index < this.imageList.length) {
					this.imageList.splice(index, 1);
				}
			},

			/**
			 * 快速输入选择
			 */
			chooseMsg() {
				uni.showActionSheet({
					itemList: this.msgContents,
					success: res => {
						const selectedText = this.msgContents[res.tapIndex];
						if (selectedText) {
							// 限制最大长度
							const maxLength = 500;
							const newContent = this.sendData.content ?
								`${this.sendData.content}\n${selectedText}` :
								selectedText;

							if (newContent.length <= maxLength) {
								this.sendData.content = newContent;
							} else {
								uni.showToast({
									title: '内容已达最大长度',
									icon: 'none',
									duration: 1500
								});
							}
						}
					},
					fail: err => {
						console.log('选择快速输入失败：', err);
					}
				});
			},

			/**
			 * 选择图片
			 */
			chooseImg() {
				const remainCount = 5 - this.imageList.length;
				if (remainCount <= 0) return;

				uni.chooseImage({
					sourceType: ['camera', 'album'],
					sizeType: 'compressed',
					count: remainCount,
					success: res => {
						// 过滤过大的图片
						const validImages = res.tempFilePaths.filter(path => {
							const size = res.tempFiles.find(f => f.path === path)?.size || 0;
							if (size > 10 * 1024 * 1024) { // 10M
								uni.showToast({
									title: '单张图片不能超过10M',
									icon: 'none',
									duration: 1500
								});
								return false;
							}
							return true;
						});
						this.imageList = this.imageList.concat(validImages);
					},
					fail: err => {
						console.log('选择图片失败：', err);
					}
				});
			},

			/**
			 * 预览图片
			 * @param {Number} index - 图片索引
			 */
			previewImage(index) {
				if (!this.imageList.length) return;
				uni.previewImage({
					urls: this.imageList,
					current: this.imageList[index],
					fail: err => {
						console.log('预览图片失败：', err);
					}
				});
			},

			/**
			 * 提交反馈
			 */
			async send() {
				// 防止重复提交
				if (this.isSubmitting) return;

				// 验证内容
				const content = this.sendData.content.trim();
				if (!content) {
					return uni.showModal({
						content: '请输入问题和意见',
						showCancel: false
					});
				}
				if (this.orderItemId < 1) {
					return uni.showModal({
						content: '订单号错误',
						showCancel: false
					});
				}
				if (this.sendData.contact.length < 7) {
					return uni.showModal({
						content: '请填写您的手机、QQ或E-mail，方便我们联系您',
						showCancel: false
					});
				}

				this.isSubmitting = true;
				uni.showLoading({
					title: '上传中...',
					mask: true
				});

				try {
					// 处理图片上传
					let simgs = [];
					if (this.imageList.length) {
						const imgs = this.imageList.map((value, index) => ({
							name: `images${index}`,
							uri: value
						}));

						// 区分平台处理上传
						const wxUploadRes = await UploadFileWx(imgs);
						if (wxUploadRes) {
							simgs = wxUploadRes || [];
						}

						// if (uni.getSystemInfoSync().platform === 'weixin') {
						// 	const wxUploadRes = await UploadFileWx(imgs);
						// 	simgs = wxUploadRes.success || [];
						// } else {
						// 	const res = await UploadFile(imgs);
						// 	if (res.fail) {
						// 		throw new Error(res.fail.error || '图片上传失败');
						// 	}
						// 	simgs = res.success || [];
						// }
					}
					let orderItemId = parseInt(this.orderItemId)
					// 提交反馈
					const resp = await DealOrder(
						orderItemId,
						content,
						this.sendData.contact.trim(),
						JSON.stringify(simgs)

					);
					if (resp?.code === 0) {
						uni.showModal({
							content: '成功反馈，我们会尽快处理',
							showCancel: false,
							success: () => {
								this.$api.prePage().$emit('feedback', this.index);
								uni.navigateBack();
							}
						});
					} else {
						throw new Error(resp?.message || '提交反馈失败，请稍后重试');
					}
				} catch (err) {
					uni.showToast({
						title: err.message || '操作失败',
						icon: 'none',
						duration: 2000
					});
				} finally {
					this.isSubmitting = false;
					uni.hideLoading();
				}
			},
		}
	};
</script>

<style scoped>
	page {
		background-color: #efeff4;
		padding: 20rpx;
		min-height: 100vh;
		box-sizing: border-box;
	}

	.bill-id {
		display: block;
		font-size: 28rpx;
		color: #666;
		margin-bottom: 20rpx;
		padding-left: 10rpx;
		line-height: 1.5;
	}

	.feedback-section {
		background-color: #fff;
		border-radius: 16rpx;
		padding: 20rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
	}

	.feedback-title {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: 32rpx;
		color: #333;
		font-weight: 500;
		margin-bottom: 15rpx;
		padding-left: 5rpx;
	}

	.feedback-quick {
		font-size: 28rpx;
		color: #007aff;
		text-decoration: underline;
		cursor: pointer;
	}

	.feedback-body {
		width: 100%;
	}

	.feedback-textarea {
		width: 100%;
		min-height: 200rpx;
		padding: 15rpx;
		font-size: 28rpx;
		border: 1px solid #eee;
		border-radius: 10rpx;
		resize: none;
		box-sizing: border-box;
		line-height: 1.6;
	}

	.feedback-textarea::placeholder {
		color: #ccc;
	}

	.word-count {
		display: block;
		text-align: right;
		font-size: 24rpx;
		color: #999;
		margin-top: 8rpx;
		padding-right: 5rpx;
	}

	.feedback-uploader {
		padding: 10rpx 0;
	}

	.uni-uploader {
		width: 100%;
	}

	.uni-uploader-head {
		display: flex;
		justify-content: space-between;
		margin-bottom: 15rpx;
		font-size: 26rpx;
		color: #666;
	}

	.uni-uploader__files {
		display: flex;
		flex-wrap: wrap;
		gap: 15rpx;
	}

	.uni-uploader__file {
		width: 160rpx;
		height: 160rpx;
		border-radius: 10rpx;
		overflow: hidden;
		border: 1px solid #eee;
	}

	.uni-uploader__img {
		width: 100%;
		height: 100%;
	}

	.close-view {
		text-align: center;
		line-height: 36rpx;
		height: 36rpx;
		width: 36rpx;
		border-radius: 50%;
		background: rgba(255, 80, 83, 0.9);
		color: #ffffff;
		position: absolute;
		top: -10rpx;
		right: -10rpx;
		font-size: 24rpx;
		z-index: 10;
		box-shadow: 0 2rpx 5rpx rgba(0, 0, 0, 0.2);
		cursor: pointer;
	}

	.uni-uploader__input-box {
		width: 160rpx;
		height: 160rpx;
		border: 1px dashed #ddd;
		border-radius: 10rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: #f9f9f9;
	}

	.uni-uploader__input {
		width: 100%;
		height: 100%;
	}

	.feedback-input {
		width: 100%;
		height: 80rpx;
		padding: 0 15rpx;
		font-size: 28rpx;
		border: 1px solid #eee;
		border-radius: 10rpx;
		box-sizing: border-box;
	}

	.feedback-input::placeholder {
		color: #ccc;
	}

	.feedback-submit {
		width: 100%;
		height: 90rpx;
		line-height: 90rpx;
		font-size: 32rpx;
		border-radius: 16rpx;
		background-color: #007aff;
		color: #fff;
		margin: 30rpx 0;
	}

	.feedback-submit:disabled {
		opacity: 0.7;
	}

	.feedback-note {
		text-align: center;
		font-size: 26rpx;
		color: #999;
		margin-top: 10rpx;
		line-height: 1.5;
	}
</style>