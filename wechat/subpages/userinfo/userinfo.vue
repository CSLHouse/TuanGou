<template>
	<view>
		<view class="user-section" >
			<image class="bg" src="/static/user-bg.jpg"></image>
			<text class="bg-upload-btn yticon icon-paizhao"></text>
			<view class="portrait-box" @click="handleUserinfo">
				<image class="portrait" :src="userInfo.avatarUrl || '/static/missing-face.png'"></image>
				<text class="pt-upload-btn yticon icon-paizhao"></text>
			</view>
			<view class="info-box">
				<text class="username">{{hasLogin ? userInfo.nickName || userInfo.telephone : '修改昵称'}}</text>
			</view>
		</view>
		<view v-if='!isCloseNickNameModel' >
			<div class="modal-mask" @click="closeNickNamePop"/>
			<div class="modal-dialog">
			  <div class="modal-content">
			    <image class="img" src="/static/pop.jpg"></image>
			    <div class="content-text">
			      <p class="info-bold-tip">完善信息可体验更多功能</p>
			      <p class="key-bold">99%用户选择使用微信昵称</p>
			    </div>
			  </div>
			  <view class="avatarUrl">
				  <text >头像：</text>
				  <button  open-type="chooseAvatar" @chooseavatar="handleChooseavatar">
					<image :src="avatarUrl || '/static/missing-face.png'" class="avatar-img"></image>
					<text class="yticon icon-you you-btn"></text>
				  </button>
			  </view>
			  <view class="nickname">
				  <text >昵称：</text>
				  <input type="nickname" class="weui-input" placeholder="请输入昵称" maxlength="15" v-model="nickName"
				   @change="getNickname" />
			  </view>
			  <div class="modal-footer">
			    <button class='btn' @click="handleConfirmNickName">
			    	确认
			    </button>
			  </div>
			</div>
		</view>
	</view>
</template>

<script>
	import {  
	    mapState,  
	    mapMutations  
	} from 'vuex';
	import common from '@/utils/common.js'
	import { getWXPhoneNumber, wxRefreshLogin, WXResetNickName} from '@/api/member.js';
	export default {
		data() {
			return {
				isCloseNickNameModel: true,
				nickName: '',
				avatarUrl: null,
			};
		},
		computed:{
			...mapState(['userInfo']),
		},
		methods:{
			handleChooseavatar(e) {
				this.avatarUrl = e.detail.avatarUrl;
			},
			getNickname(e) {
				this.nickName = e.detail.value
			},
			closeNickNamePop() {
				this.isCloseNickNameModel = true
			},
			async handleConfirmNickName () {
				let _this = this
				let isUpload = false
				if (_this.avatarUrl && _this.avatarUrl != '' && _this.nickName != '') {
					uni.uploadFile({
						url: common.baseUrl + "/fileUploadAndDownload/upload",
						filePath: _this.avatarUrl,
						name: 'file',
						header: {
							"x-token": _this.$store.state.token,
							"x-user_id": _this.$store.state.userInfo.id,
							"Access-Control-Allow-Origin": "*",
							"Access-Control-Allow-Methods": "*"
						},
						success: res =>{
							const response = JSON.parse(res.data)
							if (response.code == 0) {
								 _this.$store.state.userInfo.avatarUrl = response.data.file.url
								 _this.$store.state.userInfo.nickName = _this.nickName
								_this.$store.state.hadNickName = true
								uni.setStorage({ //缓存用户登陆状态
									key: 'UserInfo',  
									data: _this.$store.state.userInfo  
								})
								_this.isCloseNickNameModel = true
								_this.isUpload = true
								this.resetNick()
							}
						},
						fail: (error) => {
							this.$api.msg('设置失败')
						}
					})
				} else {
					if (_this.nickName != '') {
						_this.$store.state.userInfo.nickName = _this.nickName
						_this.$store.state.hadNickName = true
						_this.isCloseNickNameModel = true
						_this.isUpload = true
						this.resetNick()
					}
				}
			},
			async resetNick() {
				let _this = this
				if (_this.$store.state.userInfo && _this.isUpload) {
					_this.$store.state.userInfo.nickName = _this.nickName
					WXResetNickName(_this.$store.state.userInfo).then(res=>{
						if (res.code == 0) {
							_this.$store.state.hadNickName = true
							uni.setStorage({ //缓存用户登陆状态
							    key: 'UserInfo',  
							    data: _this.$store.state.userInfo  
							})
							_this.isCloseNickNameModel = true
							this.$api.msg('设置成功')
						}
						else {
							this.$api.msg('设置失败')
						}
					});
				} else {
					this.$api.msg('头像设置失败')
					_this.isCloseNickNameModel = true
				}
			},
			closeNickNamePop() {
				this.isCloseNickNameModel = true
			},
			handleUserinfo() {
				this.isCloseNickNameModel = false
			}
		}
	}
</script>

<style lang="scss">
	page{
		background: $page-color-base;
	}
	.user-section{
		display:flex;
		align-items:center;
		justify-content: center;
		height: 460upx;
		padding: 40upx 30upx 0;
		position:relative;
		.bg{
			position:absolute;
			left: 0;
			top: 0;
			width: 100%;
			height: 100%;
			filter: blur(1px);
			opacity: .7;
		}
		.portrait-box{
			width: 200upx;
			height: 200upx;
			border:6upx solid #fff;
			border-radius: 50%;
			position:relative;
			z-index: 2;
		}
		.portrait{
			position: relative;
			width: 100%;
			height: 100%;
			border-radius: 50%;
		}
		.yticon{
			position:absolute;
			line-height: 1;
			z-index: 5;
			font-size: 48upx;
			color: #fff;
			padding: 4upx 6upx;
			border-radius: 6upx;
			background: rgba(0,0,0,.4);
		}
		.pt-upload-btn{
			right: 0;
			bottom: 10upx;
		}
		.bg-upload-btn{
			right: 20upx;
			bottom: 16upx;
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
	  padding:10px 0px 10px 0px;
	  font-size:14px;
	  display: flex;
	  
	}
	.modal-footer {
	  box-sizing: border-box;
	  display: flex;
	  flex-direction: row;
	  border-top: 1px solid #e5e5e5;
	  font-size: 16px;
	  font-weight:bold;
	  /* height: 45px; */
	  line-height: 45px;
	  text-align: center;
	  background:#feb600;
	}
	.btn {
	  width: 100%;
	  height: 100%;
	  background:#feb600;
	  color:#FFFFFF;
	  font-weight:bold;
	}
	.img {
	  width: 560rpx;
	  height:140rpx;
	}
	.little-tip {
	  padding-top:15px;
	  padding-bottom:3px;
	  font-size: 14px;
	  font-weight:bold;
	  color: #feb600;
	}
	.little-content {
	  padding-top:5px;
	  font-size: 13px;
	  color:#606060;
	}
	.key-bold-tip {
	  padding-top:5px;
	  font-size: 15px;
	  font-weight:bold;
	  color: #feb600;
	}
	.key-bold {
	  padding-top:5px;
	  font-size: 14px;
	  /* font-weight:bold; */
	}
	.info-bold-tip {
		padding-top:5px;
		font-size: 15px;
		font-weight:bold;
		color: #feb600;
		text-align: center;
	}
	
	.avatarUrl{
		background: #fff;
		padding: 20rpx 20rpx 10rpx;
		/* display: flex; */
		align-items: center;
		justify-content: center;
		border-bottom: 1px solid #f5f5f5;
		height: 100rpx;
		button {
			background: rgba(0, 0, 0, 0);
			float: right;
			.avatar-img{
				height: 60rpx;
				width: 60rpx;
				border-radius: 50%;
			}
		}
		button::after {
			border: none;
		}
		.you-btn {
			margin-bottom: 10rpx;
			float: right;
		}
	}
	.nickname{
		background: #fff;
		padding: 20rpx 20rpx 10rpx;
		align-items: center;
		justify-content: center;
		height: 100rpx;
	
		.weui-input{
			float: right;
		}
	}
</style>
