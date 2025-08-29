package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderRecommendProduct = initOrderBrand + 1

type initRecommendProduct struct{}

// auto run
func init() {
	system.RegisterInit(initOrderRecommendProduct, &initRecommendProduct{})
}

func (i *initRecommendProduct) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.RecommendProduct{},
	)
}

func (i *initRecommendProduct) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.RecommendProduct{})
}

func (i initRecommendProduct) InitializerName() string {
	return wechatModel.RecommendProduct{}.TableName()
}

func (i *initRecommendProduct) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []wechatModel.RecommendProduct{
		{
			ProductId:       1,
			ProductName:     "猪迪克体验卡",
			RecommendStatus: 1,
			Sort:            0,
		},
		{
			ProductId:       2,
			ProductName:     "猪迪克单人票",
			RecommendStatus: 1,
			Sort:            0,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.RecommendProduct{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initRecommendProduct) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("product_id = ?", 1).First(&wechatModel.RecommendProduct{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
