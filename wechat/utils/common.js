// const baseUrl = "http://wx.tiepiwa.cn:8888"
const baseUrl = "http://127.0.0.1:8888"
// 保留小数点数值后两位，尾数四舍五入
export function numFilter(value) {
	// 截取当前数据到小数点后两位
	let realVal = parseFloat(value).toFixed(2)
	return realVal
}

export function strToJson(str) {
	var json = eval('(' + str + ')');
	return json;
}

export default {
	numFilter,
	baseUrl,
	strToJson,
}