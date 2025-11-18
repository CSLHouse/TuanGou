package wechat

import (
	"context"
	"cooller/server/model/product"
	"cooller/server/service/system"

	"gorm.io/gorm"
)

const initOrderReceiveAddress = initOrderRecommendProduct + 1

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
		&product.MemberReceiveAddress{},
	)
}

func (i *initMemberReceiveAddress) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&product.MemberReceiveAddress{})
}

func (i initMemberReceiveAddress) InitializerName() string {
	return product.MemberReceiveAddress{}.TableName()
}

func (i *initMemberReceiveAddress) InitializeData(ctx context.Context) (next context.Context, err error) {
	return next, err
}

func (i *initMemberReceiveAddress) DataInserted(ctx context.Context) bool {
	return true
}
