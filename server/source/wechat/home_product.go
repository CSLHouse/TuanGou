package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

const initOrderHomeProduct = initOrderRecommendProduct + 1

type initHomeProduct struct{}

// auto run
func init() {
	system.RegisterInit(initOrderHomeProduct, &initHomeProduct{})
}

func (i *initHomeProduct) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.Product{},
	)
}

func (i *initHomeProduct) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.Product{})
}

func (i initHomeProduct) InitializerName() string {
	return wechatModel.Product{}.TableName()
}

func (i *initHomeProduct) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	startTime, err := time.Parse("2006-01-02 15:04:05", "2023-12-09 00:00:00")
	endTime, err := time.Parse("2006-01-02 15:04:05", "2024-02-25 00:00:00")
	entities := []wechatModel.Product{
		{
			BrandId:                    1,
			ProductCategoryId:          1,
			FeightTemplateId:           0,
			ProductAttributeCategoryId: 1,
			Name:                       "猪迪克体验卡",
			Pic:                        "https://cooller.oss-cn-beijing.aliyuncs.com/resource/uploads/1/1739833092480176128_zhudike_banner.png",
			ProductSN:                  "6946605",
			DeleteStatus:               0,
			PublishStatus:              1,
			NewStatus:                  1,
			RecommandStatus:            1,
			VerifyStatus:               0,
			Sort:                       100,
			Sale:                       100,
			Price:                      9.9,
			PromotionPrice:             9.9,
			GiftGrowth:                 100,
			GiftPoint:                  100,
			UsePointLimit:              0,
			SubTitle:                   "体验卡 9.9畅玩一整天",
			Description:                "",
			OriginalPrice:              25.8,
			Stock:                      100,
			LowStock:                   0,
			Unit:                       "张",
			Weight:                     0.00,
			PreviewStatus:              1,
			ServiceIds:                 "0",
			Keywords:                   "",
			Note:                       "",
			AlbumPics:                  "",
			DetailTitle:                "",
			DetailDesc:                 "",
			DetailHTML:                 "",
			DetailMobileHTML:           "",
			PromotionStartTime:         &startTime,
			PromotionEndTime:           &endTime,
			PromotionPerLimit:          0,
			PromotionType:              1,
			BrandName:                  "猪迪克星动乐园",
			ProductCategoryName:        "门票",
		},
		{
			BrandId:                    1,
			ProductCategoryId:          1,
			FeightTemplateId:           0,
			ProductAttributeCategoryId: 1,
			Name:                       "猪迪克单人票",
			Pic:                        "https://cooller.oss-cn-beijing.aliyuncs.com/resource/uploads/1/1742477305617321984_%E6%96%B0%E5%BB%BA%E9%A1%B9%E7%9B%AE%20%284%29.jpg",
			ProductSN:                  "6946605",
			DeleteStatus:               0,
			PublishStatus:              1,
			NewStatus:                  1,
			RecommandStatus:            1,
			VerifyStatus:               0,
			Sort:                       100,
			Sale:                       100,
			Price:                      15.8,
			PromotionPrice:             15.8,
			GiftGrowth:                 100,
			GiftPoint:                  100,
			UsePointLimit:              0,
			SubTitle:                   "一大一小 无需预约 空调开放",
			Description:                "",
			OriginalPrice:              25.8,
			Stock:                      100,
			LowStock:                   0,
			Unit:                       "张",
			Weight:                     0.00,
			PreviewStatus:              1,
			ServiceIds:                 "0",
			Keywords:                   "",
			Note:                       "",
			AlbumPics:                  "",
			DetailTitle:                "",
			DetailDesc:                 "",
			DetailHTML:                 "",
			DetailMobileHTML:           "",
			PromotionStartTime:         &startTime,
			PromotionEndTime:           &endTime,
			PromotionPerLimit:          0,
			PromotionType:              0,
			BrandName:                  "猪迪克星动乐园",
			ProductCategoryName:        "门票",
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.Product{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initHomeProduct) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("name = ?", "猪迪克单人票").First(&wechatModel.Product{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
