package product

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/product"
	couponRes "cooller/server/model/product/request"
	"gorm.io/gorm"
	"time"
)

type CouponService struct{}

func (exa *CouponService) CreateCoupon(e *product.Coupon) (err error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return err
	}

	if len(e.ProductRelationList) > 0 {
		for _, productItem := range e.ProductRelationList {
			productItem.CouponId = e.ID
		}
		err = tx.Model(product.CouponProductRelation{}).CreateInBatches(e.ProductRelationList, len(e.ProductRelationList)).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(e.ProductCategoryRelationList) > 0 {
		for _, productCategoryItem := range e.ProductCategoryRelationList {
			productCategoryItem.CouponId = e.ID
		}
		err = tx.Model(product.CouponProductCategoryRelation{}).CreateInBatches(e.ProductCategoryRelationList, len(e.ProductCategoryRelationList)).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = global.GVA_DB.Debug().Create(&e).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (exa *CouponService) UpdateCoupon(e *product.Coupon) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *CouponService) UpdateCouponSynchronous(e *product.Coupon) (err error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return err
	}

	// 清除该优惠券已有的所有商品关联
	err = exa.DeleteCouponAndProductRelation(e.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 指定商品时处理
	err = exa.SetCouponAndProductRelation(e, tx)
	if err != nil {
		return err
	}
	// 指定分类时处理
	err = exa.DeleteCouponAndProductCategoryRelation(e.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = exa.SetCouponAndProductCategoryRelation(e, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 保存其他
	err = tx.Where("id = ?", e.ID).Save(e).Error
	return tx.Commit().Error
}

// 优惠券与商品绑定
// TODO: many2many 待优化
func (exa *CouponService) SetCouponAndProductRelation(e *product.Coupon, tx *gorm.DB) (err error) {
	if len(e.ProductRelationList) > 0 {
		for _, productItem := range e.ProductRelationList {
			productItem.CouponId = e.ID
		}

		err = tx.Model(product.CouponProductRelation{}).CreateInBatches(e.ProductRelationList, len(e.ProductRelationList)).Error
		return err
	}
	return nil
}

// 优惠券与商品分类绑定
func (exa *CouponService) SetCouponAndProductCategoryRelation(e *product.Coupon, tx *gorm.DB) (err error) {
	if len(e.ProductCategoryRelationList) > 0 {
		for _, productItem := range e.ProductCategoryRelationList {
			productItem.CouponId = e.ID
		}
		err = tx.Model(product.CouponProductCategoryRelation{}).CreateInBatches(e.ProductCategoryRelationList, len(e.ProductCategoryRelationList)).Error
	}
	return nil
}

func (exa *CouponService) UpdateCouponCount(id int) (err error) {
	db := global.GVA_DB.Model(&product.Coupon{})
	err = db.Where("id = ?", id).Updates(map[string]interface{}{
		"count":         gorm.Expr("count - ?", 1),
		"receive_count": gorm.Expr("receive_count + ?", 1),
	}).Error
	return err
}

func (exa *CouponService) DeleteCoupon(id int) (err error) {
	var coupon product.Coupon
	err = global.GVA_DB.Where("id = ?", id).Delete(&coupon).Error
	return err
}

func (exa *CouponService) GetCouponList(info request.PageInfo) (list []product.Coupon, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&product.Coupon{})
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Preload("ProductRelationList").Find(&list).Error
		return list, total, err
	}
}

func (exa *CouponService) GetCouponById(id int) (coupon product.Coupon, err error) {
	db := global.GVA_DB.Model(&product.Coupon{})
	err = db.Where("id = ?", id).Preload("ProductRelationList").First(&coupon).Error
	coupon.ProductRelationList, err = exa.GetCouponProductList(id)
	if err != nil {
		return coupon, err
	}
	coupon.ProductCategoryRelationList, err = exa.GetCouponProductCategoryList(id)
	return coupon, err
}

// GetCouponByType 根据优惠券类型获取优惠券列表 t=0,全场赠券，t=1会员赠券，t=2购物赠券，t=3注册赠券
func (exa *CouponService) GetCouponByType(t int) (couponList []product.Coupon, err error) {
	db := global.GVA_DB.Model(&product.Coupon{})
	err = db.Where("type = ?", t).Find(&couponList).Error
	return couponList, err
}

func (exa *CouponService) GetCouponHistoryList(info couponRes.SearchInfoCoupon) (list []product.CouponHistory, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&product.CouponHistory{})
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Find(&list).Error
	}
	return list, total, err
}

func (exa *CouponService) CreateCouponHistory(e *product.CouponHistory) (err error) {
	err = global.GVA_DB.Create(e).Error
	return err
}

func (exa *CouponService) DeleteCouponAndProductRelation(id int) (err error) {
	var coupon product.CouponProductRelation
	err = global.GVA_DB.Where("coupon_id = ?", id).Delete(&coupon).Error
	return err
}

func (exa *CouponService) DeleteCouponAndProductCategoryRelation(id int) (err error) {
	var coupon product.CouponProductCategoryRelation
	err = global.GVA_DB.Where("coupon_id = ?", id).Delete(&coupon).Error
	return err
}

func (exa *CouponService) GetCouponProductList(id int) (list []*product.CouponProductRelation, err error) {
	err = global.GVA_DB.Model(product.CouponProductRelation{}).Where("coupon_id = ?", id).Find(&list).Error
	return list, err
}

func (exa *CouponService) GetCouponProductCategoryList(id int) (list []*product.CouponProductCategoryRelation, err error) {
	err = global.GVA_DB.Model(product.CouponProductCategoryRelation{}).Where("coupon_id = ?", id).Find(&list).Error
	return list, err
}

func (exa *CouponService) GetCouponHistoryCount(couponId int, userId int) (total int64, err error) {
	db := global.GVA_DB.Model(&product.CouponHistory{})
	err = db.Where("coupon_id = ? and member_id = ?", couponId, userId).Count(&total).Error
	if err != nil {
		return total, err
	}
	return total, err
}

func (exa *CouponService) GetUserCouponHistoryById(userId int) (list []*product.CouponHistory, err error) {
	db := global.GVA_DB.Model(&product.CouponHistory{})
	err = db.Where(" member_id = ?", userId).Preload("Coupon").Find(&list).Error

	// TODO: 优化获取方式
	for _, v := range list {
		v.Coupon.ProductRelationList, err = exa.GetCouponProductList(v.CouponId)
		v.Coupon.ProductCategoryRelationList, err = exa.GetCouponProductCategoryList(v.CouponId)
	}
	return list, err
}

func (exa *CouponService) UpdateCouponStatus(couponId int, userId int, useStatus int) error {
	if couponId == 0 {
		return nil
	}
	useStatusQ := 0
	if useStatus == 0 {
		useStatusQ = 1
	}
	now := time.Now()
	db := global.GVA_DB.Model(&product.CouponHistory{})
	err := db.Where("coupon_id = ? and member_id = ? and use_status = ?", couponId, userId, useStatusQ).
		Updates(product.CouponHistory{UseTime: &now, UseStatus: useStatus}).Error
	return err
}
