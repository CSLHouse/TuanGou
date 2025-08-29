package business

import (
	"bytes"
	"cooller/server/global"
	"cooller/server/model/business"
	"cooller/server/model/common/request"
	date_conversion "cooller/server/utils/timer"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type VIPOrderService struct{}

var VIPOrderServiceApp = new(VIPOrderService)

func (exa *VIPOrderService) CreateVIPOrder(e *business.VIPOrder) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetVIPOrderInfoList
//@description: 获取套餐列表
//@param: sysUserAuthorityID string, info request.PageInfo
//@return: list interface{}, total int64, err error

func (exa *VIPOrderService) GetVIPOrderInfoList(userId int, searchInfo request.OrderSearchInfo) (list interface{}, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	var orderList []business.VIPOrder
	cmd := fmt.Sprintf("sys_user_id = %d", userId)
	if searchInfo.Telephone >= 1000 {
		cmd += fmt.Sprintf(" and telephone like '%%%d%%'", searchInfo.Telephone)
	}
	if len(searchInfo.OrderId) > 1 {
		cmd += fmt.Sprintf(" and order_id like '%%%s%%'", searchInfo.OrderId)
	}
	if searchInfo.Type > 0 {
		cmd += fmt.Sprintf(" and type = %d", searchInfo.Type)
	}
	if len(searchInfo.BuyDate) > 1 {
		cmd += fmt.Sprintf(" and buy_date = %s", searchInfo.BuyDate)
	}

	if limit > 0 && offset > 0 {
		cmd += fmt.Sprintf(" limit %d offset %d", limit, offset)
	}
	db := global.GVA_DB.Model(&business.VIPOrder{})
	err = db.Where(cmd).Count(&total).Error
	if err != nil {
		return orderList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Preload("Card").Preload("Card.Combo").Where(cmd).Find(&orderList).Error
	}

	//err = global.GVA_DB.Where("sysUserAuthorityID = ? and telephone like ?", sysUserAuthorityID, telephone+"%").First(&member).Error
	return orderList, total, err
}

func (exa *VIPOrderService) CreateVIPStatement(e *business.VIPStatement) (err error) {
	var sql bytes.Buffer
	sql.WriteString("insert into bus_statement(updated_at, recharge, card_number,new_member,consume_number,sys_user_id) values (")
	sql.WriteString("\"")
	sql.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	sql.WriteString("\",")
	sql.WriteString(strconv.Itoa(e.Recharge))
	sql.WriteString(",")
	sql.WriteString(strconv.Itoa(e.CardNumber))
	sql.WriteString(",")
	sql.WriteString(strconv.Itoa(e.NewMember))
	sql.WriteString(",")
	sql.WriteString(strconv.Itoa(e.ConsumeNumber))
	sql.WriteString(",")
	sql.WriteString(strconv.Itoa(e.SysUserId))
	sql.WriteString(") ON DUPLICATE KEY UPDATE ")
	sql.WriteString("recharge=recharge+")
	sql.WriteString(strconv.Itoa(e.Recharge))
	sql.WriteString(",card_number=card_number+")
	sql.WriteString(strconv.Itoa(e.CardNumber))
	sql.WriteString(",new_member=new_member+")
	sql.WriteString(strconv.Itoa(e.NewMember))
	sql.WriteString(",consume_number=consume_number+")
	sql.WriteString(strconv.Itoa(e.ConsumeNumber))
	sql.WriteString(";")
	err = global.GVA_DB.Exec(sql.String()).Error
	if err != nil {
		global.GVA_LOG.Error("创建统计失败!", zap.Error(err))
		return err
	}
	//err = global.GVA_DB.Create(&e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetVIPStatementInfoList
//@description: 获取流水列表
//@param: sysUserAuthorityID string, info request.PageInfo
//@return: list interface{}, total int64, err error

func (exa *VIPOrderService) GetVIPStatementInfoList(userId int, searchInfo request.StatisticsSearchInfo) (list []*business.VIPStatement, err error) {

	var cmdList []interface{}
	var cmdString = "sys_user_id = ?"
	cmdList = append(cmdList, userId)

	if len(searchInfo.StartDate) > 0 {
		cmdList = append(cmdList, date_conversion.ParseStringDate(searchInfo.StartDate))
		cmdString += " and updated_at >= ?"
	}
	if len(searchInfo.EndDate) > 0 {
		endDate := date_conversion.ParseStringDate(searchInfo.EndDate)
		cmdList = append(cmdList, endDate)
		cmdString += " and updated_at <= ?"
	}

	db := global.GVA_DB.Model(&business.VIPStatement{})
	err = db.Debug().Where(cmdString, cmdList...).Find(&list).Error

	return list, err
}

func (exa *VIPOrderService) BuildVIPStatistics(e *business.VIPStatistics) (cmd string) {
	var sql bytes.Buffer
	sql.WriteString("insert into bus_statistics(total_stream, total_order, total_member,total_consumer,sys_user_id) values (")
	sql.WriteString("\"")
	sql.WriteString(strconv.FormatFloat(e.TotalStream, 'f', 2, 64))
	sql.WriteString("\",")
	sql.WriteString(strconv.Itoa(int(e.TotalOrder)))
	sql.WriteString(",")
	sql.WriteString(strconv.Itoa(int(e.TotalMember)))
	sql.WriteString(",")
	sql.WriteString(strconv.Itoa(int(e.TotalConsumer)))
	sql.WriteString(",")
	sql.WriteString(strconv.Itoa(int(e.SysUserId)))
	sql.WriteString(") ON DUPLICATE KEY UPDATE ")
	sql.WriteString("total_stream=total_stream+")
	sql.WriteString(strconv.Itoa(int(e.TotalStream)))
	sql.WriteString(",total_order=total_order+")
	sql.WriteString(strconv.Itoa(int(e.TotalOrder)))
	sql.WriteString(",total_member=total_member+")
	sql.WriteString(strconv.Itoa(int(e.TotalMember)))
	sql.WriteString(",total_consumer=total_consumer+")
	sql.WriteString(strconv.Itoa(int(e.TotalConsumer)))

	sql.WriteString(";")
	return sql.String()
}

func (exa *VIPOrderService) CreateVIPStatistics(e *business.VIPStatistics) (err error) {
	sql := VIPOrderServiceApp.BuildVIPStatistics(e)
	//fmt.Println("------sql:", sql.String())
	err = global.GVA_DB.Exec(sql).Error

	if err != nil {
		global.GVA_LOG.Error("创建统计失败!", zap.Error(err))
		return err
	}
	//err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *VIPOrderService) GetVIPStatisticsInfoList(userId int) (statistics business.VIPStatistics, err error) {
	err = global.GVA_DB.Where("sys_user_id = ? ", userId).First(&statistics).Error
	return statistics, err
}
