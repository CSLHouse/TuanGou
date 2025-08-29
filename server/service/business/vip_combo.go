package business

import (
	"cooller/server/global"
	"cooller/server/model/business"
	"cooller/server/model/common/request"
)

type VIPComboService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateVIPCombo
//@description: 创建套餐
//@param: e model.VIPCombo
//@return: err error

func (exa *VIPComboService) CreateVIPCombo(e business.VIPCombo) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteVIPComboById
//@description: 删除套餐
//@param: e model.VIPCombo
//@return: err error

func (exa *VIPComboService) DeleteVIPComboById(id int, userId int) (err error) {
	var combo business.VIPCombo
	err = global.GVA_DB.Where("id = ? and sys_user_id = ?", id, userId).Delete(&combo).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateVIPCombo
//@description: 更新套餐
//@param: e *model.VIPCombo
//@return: err error

func (exa *VIPComboService) UpdateVIPCombo(e *business.VIPCombo) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetVIPCombo
//@description: 获取套餐信息
//@param: id int
//@return: customer model.VIPCombo, err error

func (exa *VIPComboService) GetVIPComboById(id int, userId int) (customer business.VIPCombo, err error) {
	err = global.GVA_DB.Where("id = ?  and sys_user_id = ?", id, userId).First(&customer).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetVIPComboInfoList
//@description: 获取套餐列表
//@param: sysUserAuthorityID string, info request.PageInfo
//@return: list interface{}, total int64, err error

func (exa *VIPComboService) GetVIPComboInfoList(sysUserId int, info request.PageInfo) (list []business.VIPCombo, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&business.VIPCombo{})
	err = db.Where("sys_user_id = ?", sysUserId).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Where("sys_user_id = ?", sysUserId).Find(&list).Error
	}
	return list, total, err
}

func (exa *VIPComboService) GetAllVIPComboInfoList(sysUserId int) (list []business.VIPCombo, err error) {
	db := global.GVA_DB.Model(&business.VIPCombo{})
	err = db.Where("sys_user_id = ?", sysUserId).Find(&list).Error
	return list, err
}
