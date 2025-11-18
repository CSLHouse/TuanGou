package product

import (
	"context"
	productModel "cooller/server/model/product"
	"cooller/server/service/system"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderHomeProductAttribute = initOrderHomeProduct + 1

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
		&productModel.ProductAttribute{},
		&productModel.ProductAttributeValue{},
	)
}

func (i *initHomeProductAttribute) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&productModel.ProductAttribute{})
}

func (i initHomeProductAttribute) InitializerName() string {
	return productModel.ProductAttribute{}.TableName()
}

func (i *initHomeProductAttribute) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []productModel.ProductAttribute{
		{
			ProductAttributeCategoryId: 1,
			Name:                       "颜色",
			SelectType:                 2,
			InputType:                  1,
			InputList:                  "黑色,红色",
			Sort:                       100,
			FilterType:                 0,
			SearchType:                 0,
			RelatedStatus:              0,
			HandAddStatus:              1,
			Type:                       0,
		},
		{
			ProductAttributeCategoryId: 1,
			Name:                       "尺寸",
			SelectType:                 2,
			InputType:                  1,
			InputList:                  "4-8岁,8-12岁,12-16岁",
			Sort:                       0,
			FilterType:                 0,
			SearchType:                 0,
			RelatedStatus:              0,
			HandAddStatus:              0,
			Type:                       0,
		},
		{
			ProductAttributeCategoryId: 1,
			Name:                       "适用人群",
			SelectType:                 1,
			InputType:                  1,
			InputList:                  "儿童,青年,中年,老年",
			Sort:                       0,
			FilterType:                 0,
			SearchType:                 0,
			RelatedStatus:              0,
			HandAddStatus:              0,
			Type:                       1,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, productModel.ProductAttribute{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initHomeProductAttribute) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("name = ?", "颜色").First(&productModel.ProductAttribute{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
