package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderHome = system.InitWechatInternal + 1

type initHome struct{}

// auto run
func init() {
	system.RegisterInit(initOrderHome, &initHome{})
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
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []wechatModel.Advertise{
		{
			Name:       "猪迪克推荐广告",
			Type:       1,
			Pic:        "https://cooller.oss-cn-beijing.aliyuncs.com/resource/uploads/1/1739833092480176128_zhudike_banner.png",
			StartTime:  "2024-01-01 17:04:03",
			EndTime:    "2024-11-08 17:04:05",
			State:      1,
			ClickCount: 0,
			OrderCount: 0,
			Url:        "/subpages/brand/brandDetail?id=1",
			Note:       "猪迪克星动乐园",
			Sort:       999,
		},
		{
			Name:       "小米推荐广告",
			Type:       1,
			Pic:        "http://macro-oss.oss-cn-shenzhen.aliyuncs.com/mall/images/20221108/xiaomi_banner_01.png",
			StartTime:  "2022-11-08 17:04:03",
			EndTime:    "2023-11-08 17:04:05",
			State:      1,
			ClickCount: 0,
			OrderCount: 0,
			Url:        "/subpages/brand/brandDetail?id=2",
			Note:       "夏季大热促销",
			Sort:       0,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.Advertise{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initHome) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("name = ?", "猪迪克推荐广告").First(&wechatModel.Advertise{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
