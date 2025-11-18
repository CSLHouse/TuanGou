package product

import (
	"context"
	productModel "cooller/server/model/product"
	"cooller/server/service/system"

	"gorm.io/gorm"
)

const initOrderHomeProduct = initOrderBrand + 1

type initHomeProduct struct{}

// auto run
func init() {
	system.RegisterInit(initOrderHomeProduct, &initHomeProduct{})
}

func (i *initHomeProduct) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&productModel.Product{},
		&productModel.ProductFullReduction{},
		&productModel.ProductLadder{},
		&productModel.SkuStock{},
	)
}

func (i *initHomeProduct) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&productModel.Product{})
}

func (i initHomeProduct) InitializerName() string {
	return productModel.Product{}.TableName()
}

func (i *initHomeProduct) InitializeData(ctx context.Context) (next context.Context, err error) {
	_, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	return next, err
}

func (i *initHomeProduct) DataInserted(ctx context.Context) bool {
	
	return true
}
