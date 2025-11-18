package product

import (
	"context"
	productModel "cooller/server/model/product"
	"cooller/server/service/system"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderHomeProductCategory = initOrderHomeProductAttribute + 1

type initHomeProductCategory struct{}

// auto run
func init() {
	system.RegisterInit(initOrderHomeProductCategory, &initHomeProductCategory{})
}

func (i *initHomeProductCategory) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&productModel.ProductCategory{},
	)
}

func (i *initHomeProductCategory) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&productModel.ProductCategory{})
}

func (i initHomeProductCategory) InitializerName() string {
	return productModel.ProductCategory{}.TableName()
}

func (i *initHomeProductCategory) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []productModel.ProductCategory{
		{
			ParentId:     0,
			Name:         "服装",
			Level:        0,
			ProductCount: 100,
			ProductUnit:  "件",
			NavStatus:    1,
			ShowStatus:   1,
			Sort:         1,
			Icon:         "",
			Keywords:     "服装",
			Description:  "服装",
		},
		{
			ParentId:     0,
			Name:         "饰品",
			Level:        0,
			ProductCount: 100,
			ProductUnit:  "件",
			NavStatus:    1,
			ShowStatus:   1,
			Sort:         1,
			Icon:         "",
			Keywords:     "饰品",
			Description:  "饰品",
		},
		{
			ParentId:     0,
			Name:         "鞋子",
			Level:        0,
			ProductCount: 100,
			ProductUnit:  "双",
			NavStatus:    1,
			ShowStatus:   1,
			Sort:         1,
			Icon:         "",
			Keywords:     "鞋子",
			Description:  "鞋子",
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, productModel.ProductCategory{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initHomeProductCategory) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("id = ?", 1).First(&productModel.ProductCategory{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
