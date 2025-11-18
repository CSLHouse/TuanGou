<template>
	<view>
		<lee-logistics :list="dataList" :cardInfo="cardInfo"></lee-logistics>
	</view>
</template>

<script>
	import {
		queryLogisticsInfo
	} from '@/api/order.js';
	import leeLogistics from '@/subpages/components/logistics/logistics.vue'
	export default {
		components: {
			leeLogistics
		},
		data() {
			return {
				orderId: null,
				dataList: [],
				cardInfo: {
					src: '',
					type: '',
					no: ''
				},
				statusList: {
					SIGN: "已签收",
					DELIVERING: "派送中",
					TRANSPORT: "运输中",
					ACCEPT: "已发货",
				},
				typeList: {
					JT: 'https://img.chinasongzhuang.cn/uploadimg/ico/2021/1027/1635301866658937.png', //韵达快递
					SF: 'https://image.suning.cn/uimg/b2c/newcatentries/0070996657-000000012347697358_1.jpg_200w_200h_4e', //顺丰速运
					ZTO: 'https://t15.baidu.com/it/u=555893266,1114630941&fm=224&app=112&f=JPEG?w=268&h=204', //中通快递
					YTO: 'https://img2.baidu.com/it/u=222406731,1463066436&fm=253&fmt=auto&app=138&f=JPEG?w=197&h=197', //圆通速递
					STO: 'https://img.duoziwang.com/2019/05/08141106806488.jpg', // 申通快递
					EMS: 'https://upload.zoular.com/article/371/d7a82ee9d1a30a7c05728b3b1a9a8888.jpeg', //EMS
					JD: 'https://www.bugela.com/cjpic/frombd/0/253/552086147/785087467.jpg', //京东快递
					DEPPON: '', //德邦快递
					POSTB: '', //邮政包裹
				}
			}
		},
		onLoad(option) {
			//商品数据
			this.orderId = option.orderId;
			this.loadData();
		},
		methods: {
			async loadData() {
				let params = {
					mailNo: "JT3139592223364"
				}
				let _this = this
				queryLogisticsInfo(params).then(response => {
					if (response.code == 0) {
						_this.dataList = response.data.dataList
						_this.cardInfo.no = response.data.no
						_this.cardInfo.type = response.data.type
						_this.cardInfo.src = _this.typeList[response.data.cpCode]
						if (_this.cardInfo.src == '') {
							_this.cardInfo.src =
								'https://img1.baidu.com/it/u=1792093058,758482164&fm=253&fmt=auto&app=138&f=JPEG?w=247&h=247'
						}
						_this.dataList = _this.dataList.map((item, index, arr) => {
							item.status = _this.statusList[item.status]
							return item
						})
						_this.$forceUpdate()
					}
				});
			},
		}
	}
</script>

<style>
	page {
		background: #eee;
	}

	.flex {
		display: flex;
		padding: 20rpx 0;
	}
</style>