<template>
	<view v-if='!hasLogin && !isCloseModel'>
		<div class="modal-mask" @click="closePop">
		</div>
		<div class="modal-dialog">
			<div class="modal-content">
				<!-- <image class="img" src="/static/pop.jpg"></image>
				<div class="content-text">
					<p class="little-tip">我们的生活圈：</p>
					<p class="little-content">
						注册成为会员，一店消费，多家优惠，欢迎体验
					</p>
				</div> -->
				<!-- 标题 -->
				<div class="title">登录后可体验更多功能</div>
				<div class="content-text">
					<p class="little-content">
						99%+用户选择使用微信头像和昵称，
						便于订单发货和售后沟通
					</p>
				</div>
				<div class="form-item">
					<label class="form-label">微信头像</label>
					<div class="form-value">
						<image :src="defaultAvatarUrl" class="avatar"></image>
						<!-- 未获取头像：显示“选择微信头像”按钮 -->
						<button class="get-btn" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
							<!-- <icon type="right" size="24rpx" color="#999" class="btn-arrow"></icon> -->
							<uni-icons type="right" size="18"></uni-icons>
						</button>
					</div>
				</div>
				<!-- 2. 微信昵称（左侧固定文字，右侧显示/输入区域） -->
				<div class="form-item">
					<label class="form-label">微信昵称</label>
					<div class="form-value">
						<!-- 已获取昵称：显示昵称 -->
						<!-- 			<span v-if="userInfo.nickName" class="nickname-text">{{ userInfo.nickName }} </span> -->
						<!-- 未获取昵称：显示输入框 -->
						<input type="nickname" placeholder="请输入昵称" v-model="customNickname" class="form-input" />
					</div>
				</div>
				<!-- 3. 城市 -->
				<div class="form-item">
					<label class="form-label">城市</label>
					<div class="form-value">
						<input type="text" placeholder="请输入所在城市" v-model="city" class="form-input" />
					</div>
				</div>
				<!-- 3. 电话号码（左侧固定文字，右侧获取/输入区域） -->
				<div class="form-item">
					<label class="form-label">电话号码</label>
					<div class="form-value">
						<!-- 已获取手机号：显示脱敏号码 -->
						<span v-if="phoneNumber"
							class="phone-text">{{ phoneNumber.replace(/^(\d{3})(\d{4})(\d{4})$/, '$1****$3') }}</span>
						<!-- 未获取手机号：显示“一键获取”按钮（带箭头） -->
						<button v-else class="get-btn" open-type="getPhoneNumber" @getphonenumber="onGetPhoneNumber">
							<uni-icons type="right" size="16" color="#999" class="btn-arrow"></uni-icons>
						</button>
					</div>
				</div>
				<!-- 4. 协议勾选（左侧勾选按钮，右侧协议文字） -->
				<div class="agreement-item">
					<!-- 核心：单个勾选按钮（uni-checkbox） -->
					<view class="custom-checkbox" @click="onRadioEvent"
						:class="{ 'custom-checkbox-active': isAgreementChecked }">
						<!-- 选中对勾（条件显示 + 动画） -->
						<view class="check-icon" :class="{ 'check-icon-active': isAgreementChecked }"></view>
					</view>
					<!-- 协议文字（含链接跳转） -->
					<label class="agreement-text">
						请阅读并勾选
						<span class="link" @click="openAgreement('user')">《用户协议》</span>
						和
						<span class="link" @click="openAgreement('privacy')">《隐私政策》</span>
					</label>
				</div>
			</div>
			<div class="modal-footer">
				<button class="confirm-btn" @click="submitInfo">
					确定
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
		CreateWxUserInfo
	} from '@/api/member.js';
	export default {
		name: "one-click-register",
		data() {
			return {
				isCloseModel: false,
				hasGotUserInfo: false,
				phoneNumber: '',
				isAgreementChecked: false,
				customNickname: '',
				defaultAvatarUrl: "/static/people.jpg",
				city: "",
				hasGotPhone: false,
			};
		},
		computed: {
			...mapState(['hasLogin']),
			isInfoComplete() {
				return this.hasGotUserInfo && this.hasGotPhone && this.city.length > 0;
			}
		},
		methods: {
			...mapMutations(['login']),
			closePop() {
				this.isCloseModel = true
			},
			// 获取用户头像、昵称、省份信息
			onChooseAvatar(e) {
				console.log("---onChooseAvatar---", e)

				this.defaultAvatarUrl = e.detail.avatarUrl
				this.hasGotUserInfo = true
			},

			// onGetLocation() {
			// 	uni.getLocation({
			// 		success: (res) => {
			// 			console.log("--success--res----", res)
			// 		},
			// 		fail: (res) => {
			// 			console.log("----fail---res---", res)
			// 		}
			// 	})
			// },
			// 获取手机号（需要后端解密）
			onGetPhoneNumber(e) {
				let _this = this
				if (e.detail.errMsg == "getPhoneNumber:ok") {
					if (_this.$store.state.openId && _this.$store.state.openId.length > 0) {
						// 这里需要将code发送到后端，调用微信接口解密
						this.getPhoneNumberFromServer(e.detail.code)
					} else {
						uni.showToast({
							title: '获取openId失败',
							icon: 'none',
							duration: 1500
						})
					}
				} else {
					uni.showToast({
						title: '请允许获取手机号',
						icon: 'none',
						duration: 1500
					})
				}
			},
			// 从服务器获取解密后的手机号
			getPhoneNumberFromServer(code) {
				let _this = this
				getWXPhoneNumber({
					openId: _this.$store.state.openId,
					code: code
				}).then(res => {
					if (res.code == 0) {
						_this.phoneNumber = res.data.phoneNumber

						this.hasGotPhone = true
					} else {
						_this.$api.msg('获取手机号失败')
					}
				});
			},
			onRadioEvent() {
				this.isAgreementChecked = !this.isAgreementChecked
			},
			// 打开用户协议/隐私政策（跳转webview）
			openAgreement(type) {
				const url = type === 'user' ?
					'/pages/webview/webview?url=https://你的域名/user-agreement.html' :
					'/pages/webview/webview?url=https://你的域名/privacy-policy.html';
				uni.navigateTo({
					url
				});
			},
			// 确认提交
			submitInfo() {
				if (!this.isAgreementChecked) {
					uni.showToast({
						title: '请勾选协议',
						icon: 'none',
						duration: 1500
					});
					return;
				}
				if (!this.isInfoComplete) {
					uni.showToast({
						title: '请完善头像、昵称和手机号',
						icon: 'none',
						duration: 1500
					});
					return;
				}

				// 这里可以添加提交逻辑，比如将获取到的信息发送给后端
				let _this = this
				CreateWxUserInfo({
					openId: _this.$store.state.openId,
					avatarUrl: _this.defaultAvatarUrl,
					nickName: _this.customNickname,
					city: _this.city,
					telephone: _this.phoneNumber,
				}).then(res => {
					if (res.code == 0) {
						const userinfo = res.data
						wx.setStorageSync("Token", userinfo.token)
						wx.setStorageSync("TokenTime", userinfo.expiresAt)
						_this.$store.state.token = userinfo.token
						_this.login(userinfo.customer);
						_this.isCloseModel = true

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

<style scoped lang="scss">
	/* 遮罩层 */
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

	/* 弹窗主体 */
	.modal-dialog {
		width: 560rpx;
		overflow: hidden;
		position: fixed;
		top: 30%;
		// left: 50%;
		// transform: translate(-50%, -50%);
		background: #fff;
		border-radius: 16rpx;
		z-index: 9999;
		overflow: hidden;
	}

	/* 弹窗内容区 */
	.modal-content {
		padding: 30rpx 40rpx;
	}

	/* 标题样式 */
	.title {
		font-size: 32rpx;
		font-weight: 600;
		color: #333;
		text-align: center;
		margin-bottom: 10rpx;
	}

	.tip {
		font-size: 24rpx;
		color: #666;
		text-align: center;
		margin-bottom: 5rpx;
	}

	.sub-tip {
		font-size: 22rpx;
		color: #999;
		text-align: center;
		margin-bottom: 40rpx;
	}

	/* 表单项样式（左侧文字+右侧内容） */
	.form-item {
		display: flex;
		justify-content: space-between;
		/* 关键：label居左，value居右 */
		align-items: center;
		/* 垂直居中对齐 */
		width: 100%;
		/* 占满父容器宽度 */
		margin-bottom: 10rpx;
		/* 表单项之间的间距 */
	}

	/* 左侧固定文字 */
	.form-label {
		width: 160rpx;
		/* 固定宽度，所有label对齐 */
		font-size: $uni-font-size-base;
		color: $font-color-dark;
		text-align: left;
		flex-shrink: 0;
		/* 禁止收缩，避免宽度被挤压 */
		line-height: 1;
		/* 清除默认行高影响 */
	}

	/* 右侧内容容器 */
	.form-value {
		flex: 1;
		/* 占满剩余宽度，确保与label分列 */
		display: flex;
		justify-content: flex-end;
		/* 核心：子组件靠右对齐 */
		align-items: center;
		/* 子组件垂直居中 */
		gap: 12rpx;
		/* 子组件之间间距，避免拥挤 */
		padding-left: 20rpx;
		/* 与label保持间距，避免紧贴 */
		min-width: 0;
		/* 关键：解决flex子元素溢出导致对齐失效 */
	}

	/* 头像样式 */
	.avatar {
		width: 60rpx;
		height: 60rpx;
		border-radius: 50%;
		object-fit: cover;
		cursor: pointer;
		display: block;
		/* 清除inline默认间距 */
	}

	/* 昵称文本：居右显示时避免文字溢出 */
	.nickname-text {
		font-size: 26rpx;
		color: #333;
		max-width: 260rpx;
		/* 限制最大宽度，超出时省略 */
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* 手机号文本：同昵称文本样式 */
	.phone-text {
		font-size: 26rpx;
		color: #333;
	}

	/* 输入框：居右显示+自适应宽度 */
	.form-input {

		/* 可选：使背景色透明 */
		width: 100%;
		/* 占满form-value宽度 */
		height: 60rpx;
		padding: 0 20rpx;
		border-bottom: 1px solid #eee;
		border-radius: 8rpx;
		font-size: 26rpx;
		color: #333;
		text-align: right;
		/* 输入框文字居左（符合用户输入习惯） */
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
		outline: none;
		/* 清除聚焦边框 */
		font-size: $uni-font-size-base;
		color: $uni-color-primary;
		line-height: 1;
		/* 清除默认行高 */
	}

	.btn-arrow {
		margin-left: 8rpx;
	}

	/* 协议勾选区域样式 */
	.agreement-item {
		display: flex;
		align-items: center;
		margin-top: 10rpx;
		margin-bottom: 30rpx;
	}

	/* 圆形Checkbox容器（单独可点击） */
	.custom-checkbox {
		width: 28rpx;
		height: 28rpx;
		border: 2px solid $border-color-dark;
		/* 未选中边框 */
		border-radius: 50%;
		/* 圆形 */
		display: flex;
		align-items: center;
		justify-content: center;
		margin-right: 12rpx;
		cursor: pointer;
		/* 提示可点击 */
		transition: all 0.2s ease;
		/* 状态切换动画 */

		/* 点击反馈：轻微缩小 */
		&:active {
			transform: scale(0.95);
		}
	}

	/* 选中状态：红色填充 + 红色边框（通过条件类名） */
	.custom-checkbox-active {
		background-color: $uni-color-primary;
		border-color: $uni-color-primary;
	}

	/* 选中对勾样式（View实现） */
	.check-icon {
		width: 12rpx;
		height: 20rpx;
		border: 2px solid transparent;
		/* 未选中时透明 */
		border-top: none;
		border-left: none;
		transform: rotate(45deg);
		transition: all 0.2s ease;
	}

	/* 选中时：显示白色对勾 */
	.check-icon-active {
		border-color: $uni-color-white;
	}

	/* 协议文字样式 */
	.agreement-text {
		font-size: 24rpx;
		color: #666;
		line-height: 1.4;
	}

	/* 协议链接样式 */
	.link {
		color: #007aff;
		text-decoration: underline;
	}

	/* 底部确定按钮 */
	.modal-footer {
		border-top: 1px solid #eee;
	}

	.confirm-btn {
		width: 100%;
		height: 90rpx;
		background: #d60000;
		color: #fff;
		font-size: 30rpx;
		font-weight: 500;
		border-radius: 0;
	}

	/* 按钮禁用状态 */
	.confirm-btn[disabled] {
		background: #f5f5f5;
		color: #ccc;
	}

	// 去掉烦人的边框
	button::after {
		border: 0; // 或者 border: none;
	}

	checkbox .wx-checkbox-input {
		width: 34rpx;
		height: 34rpx;
		border-radius: 50%;
	}

	/*checkbox选中后样式  */
	checkbox .wx-checkbox-input.wx-checkbox-input-checked {
		background: #0394F0;
		border-color: #0394F0;
	}

	/*checkbox选中后图标样式  */
	checkbox .wx-checkbox-input.wx-checkbox-input-checked::before {
		width: 20rpx;
		height: 20rpx;
		line-height: 20rpx;
		text-align: center;
		font-size: 22rpx;
		color: #fff;
		background: transparent;
		transform: translate(-50%, -50%) scale(1);
		-webkit-transform: translate(-50%, -50%) scale(1);
	}
</style>