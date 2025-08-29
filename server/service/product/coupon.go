package product

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/product"
	couponRes "cooller/server/model/product/request"
)

type CouponService struct{}

func (exa *CouponService) CreateCoupon(e *product.Coupon) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *CouponService) UpdateCoupon(e *product.Coupon) (err error) {
	err = global.GVA_DB.Save(e).Error
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
		err = db.Limit(limit).Offset(offset).Find(&list).Error
	}
	return list, total, err
}

func (exa *CouponService) GetCouponById(id int) (coupon product.Coupon, err error) {
	db := global.GVA_DB.Model(&product.Coupon{})
	err = db.Where("id = ?", id).First(&coupon).Error
	return coupon, err
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
