<template>
	<view>
		<view>
			<text class="title">展示队员n天贡献的奖励</text>
		</view>
		<!-- 1. 图表容器：沿用你示例的rpx适配，确保高度足够 -->
		<view style="width: 100%; height: 750rpx;  box-sizing: border-box;">
			<!-- lime-echart组件：ref用于获取实例，@finished触发初始化 -->
			<l-echart ref="chartRef" @finished="init" class="team-profit-chart" />
		</view>
	</view>
</template>

<script>
	// 2. 引入ECharts（路径与你示例一致，确保文件存在）
	import * as echarts from '@/uni_modules/lime-echart/static/echarts.min'
	// 3. 引入原业务接口（获取队员收益数据）
	import {
		fetchTeamDetail
	} from "@/api/team.js"

	export default {
		data() {
			return {
				userIds: [], // 队员ID列表（从页面参数获取）
				// chartInstance: null, // ECharts实例（不存data会被Vue响应式监听，这里存但仅用于更新）
				// 4. ECharts配置项：适配“队员收益折线图”业务
				option: {
					// 提示框：hover显示数据（保留axis类型，适配折线图）
					tooltip: {
						trigger: 'axis',
						axisPointer: {
							type: 'shadow'
						},
						confine: true, // 限制在容器内显示（避免溢出）
						formatter: (params) => {
							const date = params[0].axisValue;
							// 用换行符'\n'分隔，配合rich样式定义行高
							let text = `日期：${date}\n`; // \n在非HTML模式下生效
							params.forEach(item => {
								text += `${item.seriesName}：${item.value.toFixed(2)}\n`;
							});
							return text;
						}
					},
					// 图例：显示队员名称（后续从数据动态填充）
					legend: {
						data: [], // 初始空，loadData后赋值为队员名数组
						top: '5%', // 图例位置（顶部5%，避免遮挡图表）
						textStyle: {
							color: '#666'
						}
					},
					// 网格：调整边距（确保Y轴标签不被截断，参考原代码逻辑）
					// grid: {
					// 	left: '5%', // 左侧留足空间（Y轴标签+小数）
					// 	right: '5%',
					// 	bottom: '15%', // 底部留空间（X轴日期标签）
					// 	top: '15%',
					// 	containLabel: true // 关键：包含标签，避免截断
					// },
					// X轴：日期分类（原业务是日期，所以type设为category）
					xAxis: [{
						type: 'category',
						data: [], // 初始空，loadData后赋值为日期列表（如["09-17","09-18"]）
						axisLine: {
							lineStyle: {
								color: '#999'
							}
						}, // 轴线样式
						axisLabel: {
							color: '#666',
							rotate: 30, // 日期标签旋转30度（避免重叠）
							interval: 0 // 显示所有标签（不省略）
						}
					}],
					// Y轴：收益数值（解决原问题：显示1位小数）
					yAxis: [{
						type: 'value',
						axisLine: {
							lineStyle: {
								color: '#999'
							}
						}, // 轴线样式
						axisLabel: {
							color: '#666',
							// 强制显示1位小数（即使是整数，如1→1.0）
							formatter: (value) => value.toFixed(2)
						}
					}],
					// 系列：队员收益折线（初始空，loadData后赋值为队员数据）
					series: []
				}
			};
		},

		// 5. 页面加载时：获取队员ID，加载数据
		onLoad(option) {
			this.userIds = JSON.parse(option.userIds); // 从页面参数获取队员ID
			this.loadTeamProfitData(); // 加载队员收益数据
		},

		methods: {
			// 6. 核心：加载队员收益数据（复用原processChartData逻辑）
			async loadTeamProfitData() {
				try {
					// 调用原接口：传入队员ID获取数据
					const response = await fetchTeamDetail({
						userIds: this.userIds
					});
					if (response.code === 0) {
						const detailsMap = response.data.details; // 接口返回的队员数据（格式同原代码）
						// 处理数据：转为ECharts需要的格式（复用原逻辑，调整字段映射）
						const {
							categories,
							series
						} = this.processChartData(detailsMap);

						// 更新ECharts配置：填充X轴日期、图例、系列数据
						this.option.xAxis[0].data = categories; // X轴：日期列表
						this.option.legend.data = series.map(item => item.name); // 图例：队员名
						this.option.series = series; // 系列：队员收益折线数据

						// 若图表已初始化，立即更新；未初始化则等待init触发
						if (this.chartInstance) {
							this.chartInstance.setOption(this.option);
						}
					}
				} catch (error) {
					console.error("加载队员收益数据失败：", error);
				}
			},

			// 7. 数据处理：复用原逻辑，适配ECharts的series格式
			processChartData(detailsMap) {
				// 步骤1：提取所有日期并去重、排序（同原逻辑）
				const allDates = new Set();
				Object.values(detailsMap).forEach(userData => {
					userData.forEach(item => allDates.add(item.date));
				});
				const categories = Array.from(allDates).sort((a, b) =>
					new Date(a).getTime() - new Date(b).getTime()
				);

				// 步骤2：生成ECharts的series数组（每个队员一条折线）
				const series = Object.entries(detailsMap).map(([userName, userData]) => {
					// 日期→收益的映射（方便匹配）
					const dateValueMap = {};
					userData.forEach(item => dateValueMap[item.date] = item.value || 0);

					// 生成与日期对应的收益数据（缺失补0）
					const data = categories.map(date => dateValueMap[date] || 0);

					// 返回ECharts的series项（折线样式同原代码：平滑、圆点）
					return {
						name: userName, // 队员名（图例和提示框显示）
						type: 'line', // 图表类型：折线图（原业务需求）
						data: data, // 队员收益数据
						smooth: true, // 平滑折线（同原代码）
						symbol: 'circle', // 数据点样式：圆形（同原代码）
						symbolSize: 8, // 数据点大小（同原代码）
						lineStyle: {
							width: 2
						} // 折线宽度（优化显示）
					};
				});

				return {
					categories,
					series
				}; // 返回ECharts需要的格式
			},

			// 8. ECharts初始化：组件渲染完成后触发（@finished事件）
			async init() {
				// 获取ECharts实例（lime-echart的init方法）
				this.chartInstance = await this.$refs.chartRef.init(echarts);
				// 设置配置项：此时option已填充数据（若loadData先完成）
				this.chartInstance.setOption(this.option);

				// 监听窗口 resize（可选：适配屏幕旋转）
				const query = uni.createSelectorQuery().in(this);
				query.select('.team-profit-chart') // ✅ 改为自定义类名
					.boundingClientRect(rect => {
						// ✅ 关键：判断rect是否存在，避免null报错
						if (rect) {
							this.chartInstance.resize({
								width: rect.width,
								height: rect.height
							});
						} else {
							console.warn("未找到图表元素，跳过resize");
							// 可选：若未找到元素，手动设置尺寸（比如用容器尺寸）
							this.chartInstance.resize({
								width: uni.getSystemInfoSync().windowWidth * 0.95, // 屏幕宽度的95%
								height: 300 // 固定高度（与容器高度一致）
							});
						}
					})
					.exec();
			}
		},

		// 9. 页面卸载：销毁ECharts实例（避免内存泄漏）
		onUnload() {
			if (this.chartInstance) {
				this.chartInstance.dispose();
				this.chartInstance = null;
			}
		}
	};
</script>

<style scoped>
	.title {
		display: flex;
		align-items: center;
		/* 垂直居中 */
		justify-content: center;
		/* 水平居中 */
		font-size: $font-lg;
		color: $uni-color-primary;
		margin-top: 10px;
	}
</style>