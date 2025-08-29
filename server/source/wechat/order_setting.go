package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderOrderSetting = initOrderProductOrder + 1

type initOrderSetting struct{}

// auto run
func init() {
	system.RegisterInit(initOrderOrderSetting, &initOrderSetting{})
}

func (i *initOrderSetting) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.OrderSetting{},
	)
}

func (i *initOrderSetting) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.OrderSetting{})
}

func (i initOrderSetting) InitializerName() string {
	return wechatModel.OrderSetting{}.TableName()
}

func (i *initOrderSetting) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []wechatModel.OrderSetting{
		{
			FlashOrderOvertime:  60,
			NormalOrderOvertime: 120,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.OrderSetting{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initOrderSetting) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("flash_order_overtime = ?", 60).First(&wechatModel.OrderSetting{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
