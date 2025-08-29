package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderNewProduct = initOrderFlashPromotion + 1

type initNewProduct struct{}

// auto run
func init() {
	system.RegisterInit(initOrderNewProduct, &initNewProduct{})
}

func (i *initNewProduct) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.NewProduct{},
		&wechatModel.ProductFullReduction{},
		&wechatModel.ProductLadder{},
		&wechatModel.SkuStock{},
		&wechatModel.CartItem{},
		&wechatModel.CartTmpItem{},
	)
}

func (i *initNewProduct) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.NewProduct{})
}

func (i initNewProduct) InitializerName() string {
	return wechatModel.NewProduct{}.TableName()
}

func (i *initNewProduct) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []wechatModel.NewProduct{
		{
			ProductId:       37,
			ProductName:     "Apple iPhone 14 (A2884) 128GB 支持移动联通电信5G 双卡双待手机",
			RecommendStatus: 1,
			Sort:            197,
		},
		{
			ProductId:       38,
			ProductName:     "Apple iPad 10.9英寸平板电脑 2022年款（64GB WLAN版/A14芯片/1200万像素/iPadOS MPQ03CH/A ）",
			RecommendStatus: 1,
			Sort:            0,
		},
		{
			ProductId:       39,
			ProductName:     "小米 Xiaomi Book Pro 14 2022 锐龙版 2.8K超清大师屏 高端轻薄笔记本电脑(新R5-6600H标压 16G 512G win11)",
			RecommendStatus: 1,
			Sort:            198,
		},
		{
			ProductId:       40,
			ProductName:     "小米12 Pro 天玑版 天玑9000+处理器 5000万疾速影像 2K超视感屏 120Hz高刷 67W快充 12GB+256GB 黑色 5G手机",
			RecommendStatus: 1,
			Sort:            200,
		},
		{
			ProductId:       41,
			ProductName:     "Redmi K50 天玑8100 2K柔性直屏 OIS光学防抖 67W快充 5500mAh大电量 墨羽 12GB+256GB 5G智能手机 小米 红米",
			RecommendStatus: 1,
			Sort:            199,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.NewProduct{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initNewProduct) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("product_id = ?", 41).First(&wechatModel.NewProduct{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
