package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

const initOrderFlashPromotion = initOrderBrand + 1

type initFlashPromotion struct{}

// auto run
func init() {
	system.RegisterInit(initOrderFlashPromotion, &initFlashPromotion{})
}

func (i *initFlashPromotion) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.FlashPromotion{},
		&wechatModel.FlashPromotionProductRelation{},
		&wechatModel.FlashPromotionSession{},
	)
}

func (i *initFlashPromotion) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.FlashPromotion{})
}

func (i initFlashPromotion) InitializerName() string {
	return wechatModel.FlashPromotion{}.TableName()
}

func (i *initFlashPromotion) InitializeData(ctx context.Context) (next context.Context, err error) {

	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	startTie, err := time.Parse("2006-01-02 15:04:05", "2023-12-09 00:00:00")
	endTime, err := time.Parse("2006-01-02 15:04:05", "2024-02-25 00:00:00")
	entities := []wechatModel.FlashPromotion{
		{
			Title:     "春节特卖活动",
			StartDate: startTie,
			EndDate:   endTime,
			Status:    1,
		},
	}

	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.FlashPromotion{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initFlashPromotion) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("title = ?", "春节特卖活动").First(&wechatModel.FlashPromotion{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
