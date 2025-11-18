package product

import (
	"context"
	productModel "cooller/server/model/product"
	"cooller/server/service/system"

	"gorm.io/gorm"
)

const initCouponBoutiqueGroup = initOrderHomeProductAttributeCategory + 1

type initCouponGroup struct{}

// auto run
func init() {
	system.RegisterInit(initCouponBoutiqueGroup, &initCouponGroup{})
}

func (i *initCouponGroup) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&productModel.Coupon{},
		&productModel.CouponHistory{},
		&productModel.CouponProductRelation{},
		&productModel.CouponProductCategoryRelation{},
	)
}

func (i *initCouponGroup) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&productModel.Coupon{})
}

func (i initCouponGroup) InitializerName() string {
	return productModel.Coupon{}.TableName()
}

func (i *initCouponGroup) InitializeData(ctx context.Context) (next context.Context, err error) {

	_, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	return next, err
}

func (i *initCouponGroup) DataInserted(ctx context.Context) bool {
	return true
}
