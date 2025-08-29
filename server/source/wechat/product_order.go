package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"gorm.io/gorm"
)

const initOrderProductOrder = initOrderReceiveAddress + 1

type initProductOrder struct{}

// auto run
func init() {
	system.RegisterInit(initOrderProductOrder, &initProductOrder{})
}

func (i *initProductOrder) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.Order{},
		&wechatModel.OrderItem{},
	)
}

func (i *initProductOrder) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.Order{})
}

func (i initProductOrder) InitializerName() string {
	return wechatModel.Order{}.TableName()
}

func (i *initProductOrder) InitializeData(ctx context.Context) (next context.Context, err error) {
	return next, err
}

func (i *initProductOrder) DataInserted(ctx context.Context) bool {
	return true
}
