package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"

	"gorm.io/gorm"
)

const initOrderNewProduct = initOrderReceiveAddress + 1

type initNewProduct struct{}

// auto run
func init() {
	system.RegisterInit(initOrderNewProduct, &initNewProduct{})
}

func (i *initNewProduct) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.NewProduct{},
	)
}

func (i *initNewProduct) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.NewProduct{})
}

func (i initNewProduct) InitializerName() string {
	return wechatModel.NewProduct{}.TableName()
}

func (i *initNewProduct) InitializeData(ctx context.Context) (next context.Context, err error) {
	_, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	return next, err
}

func (i *initNewProduct) DataInserted(ctx context.Context) bool {
	return true
}
