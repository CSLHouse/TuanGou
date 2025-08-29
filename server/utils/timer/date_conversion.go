package timer

import (
	"fmt"
	"strings"
	"time"
)

func ParseStringDate(dateStr string) time.Time {
	index := strings.Index(dateStr, " ")
	if index < 1 {
		dateStr = dateStr + " 00:00:00"
	}
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	dateTime, _ := time.ParseInLocation(timeLayout, dateStr, loc)
	//dateTime := theTime.Unix()
	return dateTime
}

func DateYearLater(dateStr string, years int) string {
	dateTime := ParseStringDate(dateStr)
	addNYear := dateTime.AddDate(years, 0, 0)
	//dateStr = addNYear.Format("2006-01-02 15:04:05")
	dateStr = fmt.Sprintf("%d-%02d-%02d", addNYear.Year(), addNYear.Month(), addNYear.Day())
	return dateStr
}

// 判断是否逾期，未逾期返回false，逾期返回true
func CheckOverTime(dateStr string) bool {
	dateTime := ParseStringDate(dateStr)
	nowTime := time.Now()
	nowZero := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, nowTime.Location())
	subDay := dateTime.Sub(nowZero)
	if subDay > 0 {
		return false
	}
	return true
}

func BuildTheDayStr() string {
	nowTime := time.Now()
	return fmt.Sprintf("%d-%02d-%02d", nowTime.Year(), nowTime.Month(), nowTime.Day())
}

func BuildTheTimeStr() string {
	nowTime := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d", nowTime.Year(), nowTime.Month(), nowTime.Day(), nowTime.Hour(), nowTime.Minute())
}
