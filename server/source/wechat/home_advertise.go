package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"

	"gorm.io/gorm"
)

const initAdvertiseOrderHome = system.InitWechatInternal + 1

type initHome struct{}

// auto run
func init() {
	system.RegisterInit(initAdvertiseOrderHome, &initHome{})
}

func (i *initHome) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.Advertise{},
	)
}

func (i *initHome) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.Advertise{})
}

func (i initHome) InitializerName() string {
	return wechatModel.Advertise{}.TableName()
}

func (i *initHome) InitializeData(ctx context.Context) (next context.Context, err error) {
	_, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	return next, err
}

func (i *initHome) DataInserted(ctx context.Context) bool {
	return true
}
