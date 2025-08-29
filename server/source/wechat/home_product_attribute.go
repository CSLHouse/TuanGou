package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderHomeProductAttribute = initOrderHomeProductAttributeCategory + 1

type initHomeProductAttribute struct{}

// auto run
func init() {
	system.RegisterInit(initOrderHomeProductAttribute, &initHomeProductAttribute{})
}

func (i *initHomeProductAttribute) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.ProductAttribute{},
		&wechatModel.ProductAttributeValue{},
	)
}

func (i *initHomeProductAttribute) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.ProductAttribute{})
}

func (i initHomeProductAttribute) InitializerName() string {
	return wechatModel.ProductAttribute{}.TableName()
}

func (i *initHomeProductAttribute) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []wechatModel.ProductAttribute{
		{
			ProductAttributeCategoryId: 1,
			Name:                       "次数",
			SelectType:                 0,
			InputType:                  1,
			InputList:                  "1次",
			Sort:                       0,
			FilterType:                 0,
			SearchType:                 0,
			RelatedStatus:              0,
			HandAddStatus:              0,
			Type:                       0,
		},
		{
			ProductAttributeCategoryId: 2,
			Name:                       "次数",
			SelectType:                 0,
			InputType:                  1,
			InputList:                  "1次",
			Sort:                       100,
			FilterType:                 0,
			SearchType:                 0,
			RelatedStatus:              0,
			HandAddStatus:              1,
			Type:                       0,
		},
		{
			ProductAttributeCategoryId: 3,
			Name:                       "次数",
			SelectType:                 0,
			InputType:                  1,
			InputList:                  "10次,22次,35次",
			Sort:                       0,
			FilterType:                 0,
			SearchType:                 0,
			RelatedStatus:              0,
			HandAddStatus:              1,
			Type:                       0,
		},
		{
			ProductAttributeCategoryId: 4,
			Name:                       "款式",
			SelectType:                 1,
			InputType:                  0,
			InputList:                  "独角兽,升级版国玺寻宝,斗罗大陆,水果排队,奥特曼,大宝剑,蛋仔派对,海绵宝宝,新古代兵器",
			Sort:                       100,
			FilterType:                 0,
			SearchType:                 0,
			RelatedStatus:              0,
			HandAddStatus:              1,
			Type:                       0,
		},
		{
			ProductAttributeCategoryId: 5,
			Name:                       "款式",
			SelectType:                 0,
			InputType:                  1,
			InputList:                  "纯净水,冰红茶",
			Sort:                       100,
			FilterType:                 0,
			SearchType:                 0,
			RelatedStatus:              0,
			HandAddStatus:              1,
			Type:                       0,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.ProductAttribute{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initHomeProductAttribute) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("name = ?", "款式").First(&wechatModel.ProductAttribute{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
