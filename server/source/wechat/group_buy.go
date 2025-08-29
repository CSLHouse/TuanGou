package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderBoutiqueGroup = initOrderHomeProductAttribute + 1

type initBoutiqueGroup struct{}

// auto run
func init() {
	system.RegisterInit(initOrderBoutiqueGroup, &initBoutiqueGroup{})
}

func (i *initBoutiqueGroup) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.GroupBuyProduct{},
	)
}

func (i *initBoutiqueGroup) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.GroupBuyProduct{})
}

func (i initBoutiqueGroup) InitializerName() string {
	return wechatModel.GroupBuyProduct{}.TableName()
}

func (i *initBoutiqueGroup) InitializeData(ctx context.Context) (next context.Context, err error) {

	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []wechatModel.GroupBuyProduct{
		{
			ProductId: 1,
			Price:     9.8,
			Percent:   10,
			Status:    1,
		},
		{
			ProductId: 2,
			Price:     12.8,
			Percent:   10,
			Status:    1,
		},
	}

	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.GroupBuyProduct{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initBoutiqueGroup) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("product_id = ?", 1).First(&wechatModel.GroupBuyProduct{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
