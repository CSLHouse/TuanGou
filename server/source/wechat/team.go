package wechat

import (
	"context"
	wechatModel "cooller/server/model/wechat"
	"cooller/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderTeam = initOrderBoutiqueGroup + 1

type initTeam struct{}

// auto run
func init() {
	system.RegisterInit(initOrderTeam, &initTeam{})
}

func (i *initTeam) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&wechatModel.TeamRecord{},
		&wechatModel.TeamSequenceNum{},
		&wechatModel.TeamConsumeRecord{},
		&wechatModel.TeamSettlement{},
	)
}

func (i *initTeam) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&wechatModel.TeamRecord{})
}

func (i initTeam) InitializerName() string {
	return wechatModel.TeamRecord{}.TableName()
}

func (i *initTeam) InitializeData(ctx context.Context) (next context.Context, err error) {

	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []wechatModel.TeamRecord{
		{
			UserId:       1,
			CaptainId:    0,
			InviteCode:   "",
			TeamId:       1,
			TeamSequence: 1,
			IsActivated:  0,
			IsSettled:    0,
		},
	}

	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, wechatModel.TeamRecord{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initTeam) DataInserted(ctx context.Context) bool {
	//db, ok := ctx.Value("db").(*gorm.DB)
	//if !ok {
	//	return false
	//}
	//if errors.Is(db.Where("id = ?", 1).First(&wechatModel.TeamRecord{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
	//	return false
	//}
	return true
}
