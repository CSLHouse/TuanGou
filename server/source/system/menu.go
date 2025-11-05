package system

import (
	"context"
	. "cooller/server/model/system"
	"cooller/server/service/system"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenu = initOrderAuthority + 1

type initMenu struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenu, &initMenu{})
}

func (i initMenu) InitializerName() string {
	return SysBaseMenu{}.TableName()
}

func (i *initMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&SysBaseMenu{},
		&SysBaseMenuParameter{},
		&SysBaseMenuBtn{},
	)
}

func (i *initMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&SysBaseMenu{}) &&
		m.HasTable(&SysBaseMenuParameter{}) &&
		m.HasTable(&SysBaseMenuBtn{})
}

func (i *initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	// ParentId是menu在数据库中的id顺序，0-表示在根目录
	entities := []SysBaseMenu{
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "dashboard", Name: "dashboard", Component: "view/dashboard/index.vue", Sort: 1, Meta: Meta{Title: "仪表盘", Icon: "odometer"}},
		{MenuLevel: 0, Hidden: true, ParentId: "0", Path: "person", Name: "person", Component: "view/person/person.vue", Sort: 2, Meta: Meta{Title: "个人信息", Icon: "message"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "admin", Name: "superAdmin", Component: "view/superAdmin/index.vue", Sort: 3, Meta: Meta{Title: "超级管理员", Icon: "user"}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "authority", Name: "authority", Component: "view/superAdmin/authority/authority.vue", Sort: 1, Meta: Meta{Title: "角色管理", Icon: "avatar"}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "menu", Name: "menu", Component: "view/superAdmin/menu/menu.vue", Sort: 2, Meta: Meta{Title: "菜单管理", Icon: "tickets", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "user", Name: "user", Component: "view/superAdmin/user/user.vue", Sort: 4, Meta: Meta{Title: "员工管理", Icon: "coordinate"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "2", Path: "dictionary", Name: "dictionary", Component: "view/superAdmin/dictionary/sysDictionary.vue", Sort: 5, Meta: Meta{Title: "字典管理", Icon: "notebook"}},
		//{MenuLevel: 0, Hidden: true, ParentId: "2", Path: "dictionaryDetail/:id", Name: "dictionaryDetail", Component: "view/superAdmin/dictionary/sysDictionaryDetail.vue", Sort: 1, Meta: Meta{Title: "字典详情-${id}", Icon: "list", ActiveName: "dictionary"}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "operation", Name: "operation", Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 6, Meta: Meta{Title: "操作历史", Icon: "pie-chart"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "systemTools", Name: "systemTools", Component: "view/systemTools/index.vue", Sort: 4, Meta: Meta{Title: "系统工具", Icon: "tools"}},
		{MenuLevel: 0, Hidden: false, ParentId: "8", Path: "system", Name: "system", Component: "view/systemTools/system/system.vue", Sort: 1, Meta: Meta{Title: "系统管理", Icon: "operation"}},
		{MenuLevel: 0, Hidden: false, ParentId: "8", Path: "setting", Name: "setting", Component: "view/systemTools/setting/setting.vue", Sort: 2, Meta: Meta{Title: "系统配置", Icon: "setting"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "state", Name: "state", Component: "view/system/state.vue", Sort: 9, Meta: Meta{Title: "服务器状态", Icon: "cloudy"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "members", Name: "members", Component: "view/members/index.vue", Sort: 2, Meta: Meta{Title: "会员管理", Icon: "chat-dot-square"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "12", Path: "newMember", Name: "newMember", Component: "view/members/newMember/newMember.vue", Sort: 1, Meta: Meta{Title: "新会员办理", Icon: "shop"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "12", Path: "oldMember", Name: "oldMember", Component: "view/members/oldMember/oldMember.vue", Sort: 2, Meta: Meta{Title: "老会员续卡", Icon: "shopping-cart"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "12", Path: "showMembers", Name: "showMembers", Component: "view/members/showMembers/showMembers.vue", Sort: 3, Meta: Meta{Title: "会员查看", Icon: "shopping-cart-full"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "12", Path: "vipCombo", Name: "vipCombo", Component: "view/members/vipCombo/vipCombo.vue", Sort: 4, Meta: Meta{Title: "Vip套餐管理", Icon: "set-up"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "12", Path: "statistics", Name: "Statistics", Component: "view/members/Statistics/Statistics.vue", Sort: 5, Meta: Meta{Title: "数据统计", Icon: "Collection"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "12", Path: "orderView", Name: "orderView", Component: "view/members/orderView/orderView.vue", Sort: 6, Meta: Meta{Title: "会员订单", Icon: "Reading"}},

		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "shop", Name: "shop", Component: "view/shop/index.vue", Sort: 4, Meta: Meta{Title: "门店管理", Icon: "Shop"}},
		{MenuLevel: 0, Hidden: false, ParentId: "12", Path: "qrcode", Name: "qrcode", Component: "view/shop/qrcode/qrcode.vue", Sort: 2, Meta: Meta{Title: "二维码", Icon: "Reading"}},

		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "orderManage", Name: "orderManage", Component: "view/orderManage/index.vue", Sort: 5, Meta: Meta{Title: "订单管理", Icon: "Management"}},
		{MenuLevel: 0, Hidden: false, ParentId: "14", Path: "orderList", Name: "orderList", Component: "view/orderManage/orderList/orderList.vue", Sort: 1, Meta: Meta{Title: "订单列表", Icon: "Collection"}},
		{MenuLevel: 0, Hidden: false, ParentId: "14", Path: "orderSetting", Name: "orderSetting", Component: "view/orderManage/orderSetting/orderSetting.vue", Sort: 2, Meta: Meta{Title: "订单设置", Icon: "Reading"}},
		{MenuLevel: 0, Hidden: true, ParentId: "14", Path: "orderDetail", Name: "orderDetail", Component: "view/orderManage/orderDetail/orderDetail.vue", Sort: 3, Meta: Meta{Title: "订单详情", Icon: "Reading"}},
		{MenuLevel: 0, Hidden: true, ParentId: "14", Path: "afterSales", Name: "afterSales", Component: "view/orderManage/afterSales/afterSales.vue", Sort: 4, Meta: Meta{Title: "售后管理", Icon: "Ticket"}},

		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "product", Name: "Product", Component: "view/product/index.vue", Sort: 6, Meta: Meta{Title: "商品", Icon: "Goods"}},
		{MenuLevel: 0, Hidden: false, ParentId: "18", Path: "list", Name: "list", Component: "view/product/list/index.vue", Sort: 1, Meta: Meta{Title: "商品列表", Icon: "Box"}},
		{MenuLevel: 0, Hidden: false, ParentId: "18", Path: "add", Name: "add", Component: "view/product/add/add.vue", Sort: 2, Meta: Meta{Title: "添加商品", Icon: "Sell"}},
		{MenuLevel: 0, Hidden: false, ParentId: "18", Path: "productCategories", Name: "ProductCategories", Component: "view/product/productCategories/ProductCategories.vue", Sort: 3, Meta: Meta{Title: "商品分类", Icon: "Paperclip"}},
		{MenuLevel: 0, Hidden: false, ParentId: "18", Path: "productType", Name: "ProductType", Component: "view/product/productType/ProductType.vue", Sort: 4, Meta: Meta{Title: "商品类型", Icon: "MagicStick"}},
		{MenuLevel: 0, Hidden: false, ParentId: "18", Path: "brand", Name: "Brand", Component: "view/product/brand/Brand.vue", Sort: 5, Meta: Meta{Title: "品牌管理", Icon: "Medal"}},
		{MenuLevel: 0, Hidden: true, ParentId: "18", Path: "attribute", Name: "Attribute", Component: "view/product/attribute/Attribute.vue", Sort: 6, Meta: Meta{Title: "商品属性参数", Icon: "Medal"}},
		{MenuLevel: 0, Hidden: true, ParentId: "18", Path: "update", Name: "updateProduct", Component: "view/product/update/update.vue", Sort: 7, Meta: Meta{Title: "修改商品", Icon: "Sell"}},
		
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "marketing", Name: "Marketing", Component: "view/marketing/index.vue", Sort: 7, Meta: Meta{Title: "营销", Icon: "Van"}},
		{MenuLevel: 0, Hidden: false, ParentId: "26", Path: "flashPromotion", Name: "FlashPromotion", Component: "view/marketing/flashPromotion/index.vue", Sort: 1, Meta: Meta{Title: "秒杀活动列表", Icon: "Clock"}},
		{MenuLevel: 0, Hidden: true, ParentId: "26", Path: "productRelation", Name: "productRelation", Component: "view/marketing/flashPromotion/productRelationList.vue", Sort: 2, Meta: Meta{Title: "秒杀商品列表", Icon: "Present"}},
		{MenuLevel: 0, Hidden: true, ParentId: "26", Path: "selectSession", Name: "selectSession", Component: "view/marketing/flashPromotion/selectSessionList.vue", Sort: 3, Meta: Meta{Title: "秒杀时间段选择", Icon: "Box"}},
		{MenuLevel: 0, Hidden: true, ParentId: "26", Path: "session", Name: "session", Component: "view/marketing/flashPromotion/sessionList.vue", Sort: 4, Meta: Meta{Title: "秒杀时间段列表", Icon: "Box"}},
		{MenuLevel: 0, Hidden: false, ParentId: "26", Path: "recommend", Name: "RecommendProduct", Component: "view/marketing/recommend/RecommendProduct.vue", Sort: 5, Meta: Meta{Title: "人气推荐", Icon: "Medal"}},
		{MenuLevel: 0, Hidden: false, ParentId: "26", Path: "brandRecommend", Name: "BrandRecommend", Component: "view/marketing/brandRecommend/brandRecommend.vue", Sort: 6, Meta: Meta{Title: "品牌推荐", Icon: "Sell"}},
		{MenuLevel: 0, Hidden: false, ParentId: "26", Path: "advertise", Name: "Advertise", Component: "view/marketing/advertise/Advertise.vue", Sort: 7, Meta: Meta{Title: "广告列表", Icon: "MagicStick"}},
		{MenuLevel: 0, Hidden: false, ParentId: "26", Path: "newRecommend", Name: "NewRecommend", Component: "view/marketing/newRecommend/NewRecommend.vue", Sort: 8, Meta: Meta{Title: "新品推荐", Icon: "Medal"}},
		{MenuLevel: 0, Hidden: false, ParentId: "26", Path: "coupon", Name: "coupon", Component: "view/marketing/coupon/index.vue", Sort: 9, Meta: Meta{Title: "优惠券", Icon: "Ticket"}},
		{MenuLevel: 0, Hidden: true, ParentId: "26", Path: "addCoupon", Name: "addCoupon", Component: "view/marketing/coupon/add.vue", Sort: 10, Meta: Meta{Title: "添加优惠券", Icon: "Ticket"}},
		{MenuLevel: 0, Hidden: true, ParentId: "26", Path: "couponHistory", Name: "couponHistory", Component: "view/marketing/coupon/history.vue", Sort: 11, Meta: Meta{Title: "优惠券领取详情", Icon: "Ticket"}},
		{MenuLevel: 0, Hidden: true, ParentId: "26", Path: "updateCoupon", Name: "updateCoupon", Component: "view/marketing/coupon/update.vue", Sort: 12, Meta: Meta{Title: "修改优惠券", Icon: "Ticket"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "settlement", Name: "Settlement", Component: "view/settlement/index.vue", Sort: 8, Meta: Meta{Title: "结算", Icon: "Ticket"}},
		{MenuLevel: 0, Hidden: false, ParentId: "39", Path: "settlementList", Name: "settlementList", Component: "view/settlement/settlementList/settlementList.vue", Sort: 1, Meta: Meta{Title: "结算列表", Icon: "Document"}},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ?", "person").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
