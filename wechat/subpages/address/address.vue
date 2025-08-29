<template>
	<view class="content b-t">
		<view class="list b-b" v-for="(item, index) in addressList" :key="index" @click="checkAddress(item)">
			<view class="wrapper">
				<view class="address-box">
					<text v-if="item.defaultStatus==1" class="tag">默认</text>
					<text class="address">{{item.province}} {{item.city}} {{item.region}} {{item.detailAddress}}</text>
				</view>
				<view class="u-box">
					<text class="name">{{item.name}}</text>
					<text class="mobile">{{item.phoneNumber}}</text>
				</view>
			</view>
			<text class="yticon icon-bianji" @click.stop="addAddress('edit', item)"></text>
			<text class="yticon icon-iconfontshanchu1" @click.stop="handleDeleteAddress(item.id)"></text>
		</view>

		<button class="add-btn" @click="addAddress('add')">新增地址</button>
	</view>
</template>

<script>
	import {
		fetchAddressList,
		deleteAddress
	} from '@/api/address.js';
	export default {
		data() {
			return {
				source: 0,
				addressList: [],
			}
		},
		onLoad(option) {
			this.source = option.source;
			this.loadData();
		},
		methods: {
			async loadData() {
				let _this = this
				if (_this.$store.state.openId) {
					fetchAddressList({openId: _this.$store.state.openId}).then(response => {
						this.addressList = response.data;
					});
				} else {
					uni.showToast({ title: '请先登录', duration: 2000 })
				}
			},
			//选择地址
			checkAddress(item) {
				if (this.source == 1) {
					//this.$api.prePage()获取上一页实例，在App.vue定义
					this.$api.prePage().currentAddress = item;
					uni.navigateBack()
				}
			},
			addAddress(type, item) {
				if (type == 'edit') {
					uni.navigateTo({
						url: `/subpages/address/addressManage?type=${type}&id=${item.ID}`
					})
				} else {
					uni.navigateTo({
						url: `/subpages/address/addressManage?type=${type}`
					})
				}
			},
			//处理删除地址
			handleDeleteAddress(id){
				let superThis = this;
				uni.showModal({
				    title: '提示',
				    content: '是否要删除该地址',
				    success: function (res) {
				        if (res.confirm) {
				            deleteAddress(id).then(response=>{
								superThis.loadData();
								if (id == this.source) {
									this.$api.prePage().currentAddress = {};
								}
							});
				        } else if (res.cancel) {
				            uni.showModal({
				            	title:'提示',
				            	content:'取消',
				            	showCancel:false
				            })
				        }
				    }
				});
			},
			//添加或修改成功之后回调
			refreshList(data, type) {
				//添加或修改后事件，这里直接在最前面添加了一条数据，实际应用中直接刷新地址列表即可
				// this.addressList.unshift(data);
				this.loadData();
			}
		}
	}
</script>

<style lang='scss'>
	page {
		padding-bottom: 120upx;
	}

	.content {
		position: relative;
	}

	.list {
		display: flex;
		align-items: center;
		padding: 20upx 30upx;
		;
		background: #fff;
		position: relative;
	}

	.wrapper {
		display: flex;
		flex-direction: column;
		flex: 1;
	}

	.address-box {
		display: flex;
		align-items: center;

		.tag {
			font-size: 24upx;
			color: $base-color;
			margin-right: 10upx;
			background: #fffafb;
			border: 1px solid #ffb4c7;
			border-radius: 4upx;
			padding: 4upx 10upx;
			line-height: 1;
		}

		.address {
			font-size: 30upx;
			color: $font-color-dark;
		}
	}

	.u-box {
		font-size: 28upx;
		color: $font-color-light;
		margin-top: 16upx;

		.name {
			margin-right: 30upx;
		}
	}

	.icon-bianji {
		display: flex;
		align-items: center;
		height: 80upx;
		font-size: 40upx;
		color: $font-color-light;
		padding-left: 30upx;
	}
	
	.icon-iconfontshanchu1 {
		display: flex;
		align-items: center;
		height: 80upx;
		font-size: 40upx;
		color: $font-color-light;
		padding-left: 30upx;
	}

	.add-btn {
		position: fixed;
		left: 30upx;
		right: 30upx;
		bottom: 16upx;
		z-index: 95;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 690upx;
		height: 80upx;
		font-size: 32upx;
		color: #fff;
		background-color: $base-color;
		border-radius: 10upx;
		box-shadow: 1px 2px 5px rgba(219, 63, 96, 0.4);
	}
</style>
