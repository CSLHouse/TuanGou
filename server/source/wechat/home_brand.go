package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderBrand = initOrderHome + 1

type initBrand struct{}

// auto run
func init() {
	system.RegisterInit(initOrderBrand, &initBrand{})
}

func (i *initBrand) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.Brand{},
	)
}

func (i *initBrand) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.Brand{})
}

func (i initBrand) InitializerName() string {
	return wechatModel.Brand{}.TableName()
}

func (i *initBrand) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []wechatModel.Brand{
		{
			Name:                "猪迪克星动乐园",
			FirstLetter:         "Z",
			Sort:                500,
			FactoryStatus:       1,
			ShowStatus:          1,
			ProductCount:        100,
			ProductCommentCount: 100,
			Logo:                "https://cooller.oss-cn-beijing.aliyuncs.com/resource/uploads/1/1739925089551388672_logo.jpg",
			BigPic:              "https://cooller.oss-cn-beijing.aliyuncs.com/resource/uploads/1/1739601433512120320_%E5%BE%AE%E4%BF%A1%E5%9B%BE%E7%89%87_20231109032800.png",
			BrandStory:          "猪迪克星动乐园是以《我的朋友猪迪克》动画片为背景打造的儿童乐园。有免费玩具、太空沙、考古、手工艺品、淘气堡、海洋球、互动投影、双排滑梯、极限蹦床、沙滩乐园、萌宠转马、积木墙等项目，旨在让孩子们远离电视，健康快乐成长。",
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.Brand{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initBrand) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("name = ?", "猪迪克星动乐园").First(&wechatModel.Brand{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
