package business

import (
	"cooller/server/global"
	"cooller/server/model/business"
	"cooller/server/model/common/request"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"strings"
)

type VIPMemberService struct{}

var VIPMemberServiceApp = new(VIPMemberService)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateCustormerFormWeb
//@description: 创建客户
//@param: e model.ExaMember
//@return: err error

func (exa *VIPMemberService) CreateCustomerFormWeb(e *business.Customer) (err error) {
	db := global.GVA_DB.Model(&business.Customer{})
	var customer business.Customer
	result := db.Where("telephone = ?", e.Telephone).First(&customer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			e.UUID = uuid.Must(uuid.NewV4())
			err = global.GVA_DB.Create(e).Error
			return err
		}
		err = result.Error
	} else {
		err = db.Debug().Where("telephone = ?", e.Telephone).UpdateColumn("user_name", e.UserName).Error
		return err
	}
	return err
}

func (exa *VIPMemberService) CreateCustomerFormWechat(e *business.Customer) (err error) {
	db := global.GVA_DB.Model(&business.Customer{})
	var customer business.Customer
	result := db.Where("telephone = ?", e.Telephone).First(&customer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			e.UUID = uuid.Must(uuid.NewV4())
			err = global.GVA_DB.Create(&e).Error
			return err
		}
		err = result.Error
	} else {
		err = db.Debug().Where("telephone = ?", e.Telephone).UpdateColumn("open_id", e.OpenId).Error
		return err
	}
	return err
}

func (exa *VIPMemberService) CreateVIPMemberSynchronous(card *business.VIPCard, order *business.VIPOrder, statement *business.VIPStatement, statistics *business.VIPStatistics) (err error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return err
	}

	// 会员卡
	err = VIPMemberServiceApp.CreateVIPCard(card)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 订单
	err = VIPOrderServiceApp.CreateVIPOrder(order)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 流水
	err = VIPOrderServiceApp.CreateVIPStatement(statement)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 统计
	err = VIPOrderServiceApp.CreateVIPStatistics(statistics)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFileChunk
//@description: 删除客户
//@param: e model.ExaMember
//@return: err error

func (exa *VIPMemberService) DeleteVIPMemberById(id int, userId int) (err error) {
	var card business.VIPCard
	err = global.GVA_DB.Where("id = ? and sys_user_id = ?", id, userId).Delete(&card).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateExaMember
//@description: 更新客户
//@param: e *model.ExaMember
//@return: err error

//func (exa *VIPMemberService) UpdateVIPMember(e *business.Customer) (err error) {
//	err = global.GVA_DB.Save(e).Error
//	return err
//}

func (exa *VIPMemberService) UpdateVIPMemberSynchronous(member *business.VIPCard, statement *business.VIPStatement, statistics *business.VIPStatistics) (err error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return err
	}

	// 更新会员
	err = VIPMemberServiceApp.UpdateVIPCard(member)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 流水
	err = VIPOrderServiceApp.CreateVIPStatement(statement)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 统计
	err = VIPOrderServiceApp.CreateVIPStatistics(statistics)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
	//err = global.GVA_DB.Save(e).Error
	//return err
}

// 更新剩余次数remainTimes
//func (exa *VIPMemberService) UpdateVIPMemberRemainTimes(e *business.VIPMember, id int, num int) (err error) {
//	err = global.GVA_DB.Where("id = ?", id).Update("remainTimes", gorm.Expr("remainTimes - ?", num)).Error
//	return err
//}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetExaMember
//@description: 获取客户信息
//@param: id int
//@return: member model.ExaMember, err error

func (exa *VIPMemberService) GetVIPMemberByOpenIdWithoutCardList(id string) (member business.Customer, err error) {
	err = global.GVA_DB.Where("open_id = ?", id).First(&member).Error
	return
}

func (exa *VIPMemberService) GetVIPMemberWithTelephoneWithoutCardList(telephone int) (member business.Customer, err error) {
	err = global.GVA_DB.Where("telephone = ?", telephone).First(&member).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetMemberInfoList
//@description: 分页获取客户列表
//@param: sysUserAuthorityID string, info request.PageInfo
//@return: list interface{}, total int64, err error

func (exa *VIPMemberService) GetVIPMemberInfoList(userId int, info request.PageInfo) (list []business.VIPCard, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&business.VIPCard{})

	err = db.Where("sys_user_id = ?", userId).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Where("sys_user_id = ? and tmp = 0", userId).Preload("Customer").Preload("Combo").Find(&list).Error
	}
	return list, total, err
}

// 根据卡号、联系方式搜索会员
func (exa *VIPMemberService) SearchVIPMember(userId int, searchInfo request.MemberSearchInfo) (list []business.VIPCard, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	cmd := fmt.Sprintf("sys_user_id = %d", userId)
	if len(searchInfo.Telephone) > 1 {
		cmd += fmt.Sprintf(" and telephone like '%%%d%%'", searchInfo.Telephone)
	}
	if len(searchInfo.MemberName) > 1 {
		cmd += fmt.Sprintf(" and member_name like '%%%s%%'", strings.TrimSpace(searchInfo.MemberName))
	}
	if len(searchInfo.Deadline) >= 10 {
		cmd += fmt.Sprintf(" and deadline > '%s'", strings.TrimSpace(searchInfo.Deadline))
	}
	if searchInfo.State > 0 {
		cmd += fmt.Sprintf(" and state = %d", searchInfo.State)
	}
	if searchInfo.Tmp > 0 {
		cmd += fmt.Sprintf(" and tmp = %d", searchInfo.Tmp-100)
	}
	db := global.GVA_DB.Model(&business.VIPCard{})
	err = db.Where(cmd).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Preload("Customer").Preload("Combo").Where(cmd).Find(&list).Error
	}
	return list, total, err
}

// 根据卡号、联系方式搜索会员
func (exa *VIPMemberService) SearchVipCard(userId int, cardInfo request.CardInfo) (list []business.VIPCard, err error) {
	db := global.GVA_DB.Model(&business.VIPCard{})
	cmd := fmt.Sprintf("sys_user_id = %d and telephone like '%%%d%%' or card_id like '%%%d%%' and tmp = 0", userId, cardInfo.OnlyId, cardInfo.OnlyId)
	err = db.Where(cmd).Preload("Customer").Preload("Combo").Find(&list).Error
	return list, err
}

func (exa *VIPMemberService) CreateVIPCard(e *business.VIPCard) (err error) {
	db := global.GVA_DB.Model(&business.VIPCard{})
	var card business.VIPCard
	result := db.Where("telephone = ? and sys_user_id = ?", e.Telephone, e.SysUserId).First(&card)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			err = global.GVA_DB.Create(&e).Error
			return err
		}
		err = result.Error
	} else {
		err = db.Debug().Where("telephone = ? and sys_user_id = ?", e.Telephone, e.SysUserId).UpdateColumns(e).Error
		return err
	}

	return err
}

func (exa *VIPMemberService) UpdateVIPCard(e *business.VIPCard) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *VIPMemberService) GetVIPCardById(id int) (card business.VIPCard, err error) {
	err = global.GVA_DB.Where("id = ?", id).Preload("Combo").First(&card).Error
	return card, err
}

func (exa *VIPMemberService) GetVIPCardByTelephone(id int) (card []business.VIPCard, err error) {
	err = global.GVA_DB.Where("telephone = ?", id).Preload("Combo").Find(&card).Error
	return card, err
}

func (exa *VIPMemberService) CreateVIPCertificate(e *business.VIPCertificate) (err error) {
	db := global.GVA_DB.Model(&business.VIPCertificate{})
	var card business.VIPCertificate
	result := db.Where("telephone = ? and sys_user_id = ?", e.Telephone, e.SysUserId).First(&card)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			e.IsFirst = true
			e.Count = 1
			err = global.GVA_DB.Create(&e).Error
			return err
		}
		err = result.Error
	} else {
		result = db.Where("telephone = ? and sys_user_id = ?", e.Telephone, e.SysUserId).First(&card)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				e.Count = 1
				err = global.GVA_DB.Create(&e).Error
				return err
			}
			err = result.Error
		} else {
			err = db.Where("telephone = ? and sys_user_id = ?", e.Telephone, e.SysUserId).UpdateColumn("count", gorm.Expr("count+?", 1)).Error
			return err
		}
	}

	return err
}

func (exa *VIPMemberService) GetVIPCertificateByTelephone(id int) (card []business.VIPCertificate, err error) {
	err = global.GVA_DB.Where("telephone = ?", id).Find(&card).Error
	return card, err
}
