package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"gorm.io/gorm"
)

const initOrderReceiveAddress = initOrderWXUser + 1

type initMemberReceiveAddress struct{}

// auto run
func init() {
	system.RegisterInit(initOrderReceiveAddress, &initMemberReceiveAddress{})
}

func (i *initMemberReceiveAddress) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.MemberReceiveAddress{},
	)
}

func (i *initMemberReceiveAddress) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.MemberReceiveAddress{})
}

func (i initMemberReceiveAddress) InitializerName() string {
	return wechatModel.MemberReceiveAddress{}.TableName()
}

func (i *initMemberReceiveAddress) InitializeData(ctx context.Context) (next context.Context, err error) {
	return next, err
}

func (i *initMemberReceiveAddress) DataInserted(ctx context.Context) bool {
	return true
}
