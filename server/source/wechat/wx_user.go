package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"gorm.io/gorm"
)

const initOrderWXUser = initOrderHomeProductAttribute + 1

type initWXUser struct{}

// auto run
func init() {
	system.RegisterInit(initOrderWXUser, &initWXUser{})
}

func (i *initWXUser) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.WXUser{},
		&wechatModel.MemberPrice{},
	)
}

func (i *initWXUser) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.WXUser{})
}

func (i initWXUser) InitializerName() string {
	return wechatModel.WXUser{}.TableName()
}

func (i *initWXUser) InitializeData(ctx context.Context) (next context.Context, err error) {
	_, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	return next, err
}

func (i *initWXUser) DataInserted(ctx context.Context) bool {

	return true
}
