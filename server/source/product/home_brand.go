package product

import (
	"context"
	productModel "cooller/server/model/product"
	"cooller/server/service/system"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderBrand = system.InitProductInternal + 1

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
		&productModel.Brand{},
	)
}

func (i *initBrand) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&productModel.Brand{})
}

func (i initBrand) InitializerName() string {
	return productModel.Brand{}.TableName()
}

func (i *initBrand) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []productModel.Brand{
		{
			Name:                "自营",
			FirstLetter:         "Z",
			Sort:                500,
			FactoryStatus:       1,
			ShowStatus:          1,
			ProductCount:        100,
			ProductCommentCount: 100,
			Logo:                "https://cooller.oss-cn-beijing.aliyuncs.com/resource/uploads/1/1739925089551388672_logo.jpg",
			BigPic:              "https://cooller.oss-cn-beijing.aliyuncs.com/resource/uploads/1/1739601433512120320_%E5%BE%AE%E4%BF%A1%E5%9B%BE%E7%89%87_20231109032800.png",
			BrandStory:          "团购供销社，一手货源，物美价廉。",
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, productModel.Brand{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initBrand) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("name = ?", "自营").First(&productModel.Brand{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
