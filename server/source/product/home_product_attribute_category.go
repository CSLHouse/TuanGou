package product

import (
	"context"
	productModel "cooller/server/model/product"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderHomeProductAttributeCategory = initOrderHomeProduct + 1

type initHomeProductAttributeCategory struct{}

// auto run
func init() {
	system.RegisterInit(initOrderHomeProductAttributeCategory, &initHomeProductAttributeCategory{})
}

func (i *initHomeProductAttributeCategory) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&productModel.ProductAttributeCategory{},
	)
}

func (i *initHomeProductAttributeCategory) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&productModel.ProductAttributeCategory{})
}

func (i initHomeProductAttributeCategory) InitializerName() string {
	return productModel.ProductAttributeCategory{}.TableName()
}

func (i *initHomeProductAttributeCategory) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []productModel.ProductAttributeCategory{
		{
			Name:           "体验卡",
			AttributeCount: 1,
			ParamCount:     0,
		},
		{
			Name:           "单人票",
			AttributeCount: 1,
			ParamCount:     0,
		},
		{
			Name:           "会员卡",
			AttributeCount: 1,
			ParamCount:     0,
		},
		{
			Name:           "挖宝",
			AttributeCount: 1,
			ParamCount:     0,
		},
		{
			Name:           "饮品",
			AttributeCount: 1,
			ParamCount:     0,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, productModel.ProductAttributeCategory{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initHomeProductAttributeCategory) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("id = ?", 1).First(&productModel.ProductAttributeCategory{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
