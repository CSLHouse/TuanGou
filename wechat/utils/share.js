export default {
	data() {
		return {

		}
	},
	onLoad: function() {
		wx.showShareMenu({
			withShareTicket: true,
			menus: ["shareAppMessage", "shareTimeline"]
		})
	},
	onShareAppMessage(res) { //发送给朋友
		let that = this;
		let imageUrl = that.shareUrl || '';
		// console.log("----------scope:", this.$scope)
		if (res.from === 'button') {
		//这块需要传参，不然链接地址进去获取不到数据
			let path = `/` + that.$scope.route + `?id=` + that.$scope.options.id;
			return {
				title: '好物分享~',
				path: path,
				// imageUrl: imageUrl
			};
		}
		if (res.from === 'menu') {
			const openId = wx.getStorageSync("OpenId")
			return {
				title: '猪迪克',
				path: '/pages/index/index?refCode=' + openId,
				// imageUrl: 'https://cdn.uviewui.com/uview/swiper/1.jpg'
			};
		}
	},
	// 分享到朋友圈
	onShareTimeline() {
		const openId = wx.getStorageSync("OpenId")
		return {
			title: '猪迪克',
			path: '/pages/index/index?refCode=' + openId,
			// imageUrl: 'https://cdn.uviewui.com/uview/swiper/1.jpg'
		};
	},
	methods: {

	}
}

