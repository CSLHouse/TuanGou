// date.js
export function formatDate(date, format) {
  const o = {
      'M+': date.getMonth() + 1, // 月份
      'd+': date.getDate(), // 日
      'h+': date.getHours() % 12 === 0 ? 12 : date.getHours() % 12, // 小时
      'H+': date.getHours(), // 小时
      'm+': date.getMinutes(), // 分
      's+': date.getSeconds(), // 秒
    };
    const re = /(y+)/
    if (re.test(format)) {
      const t = re.exec(format)[1]
      format = format.replace(t, (date.getFullYear() + '').slice(4 - t.length));
    }
    for (let k in o) {
      const reg = new RegExp('(' + k + ')')
      if (reg.test(format)) {
        const t = reg.exec(format)[1]
        
        format = format.replace(
          t, t.length === 1 ? o[k] : ('00' + o[k]).slice(('' + o[k]).length)
        );
      }
    }
    return format;
}

function padLeftZero(str) {
  return ('00' + str).substr(str.length);
}

export function str2Date(dateStr, separator) {
  if (!separator) {
    separator = "-";
  }
  let dateArr = dateStr.split(separator);
  let year = parseInt(dateArr[0]);
  let month;
  //处理月份为04这样的情况
  if (dateArr[1].indexOf("0") == 0) {
    month = parseInt(dateArr[1].substring(1));
  } else {
    month = parseInt(dateArr[1]);
  }
  let day = parseInt(dateArr[2]);
  let date = new Date(year, month - 1, day);
  return date;
}
