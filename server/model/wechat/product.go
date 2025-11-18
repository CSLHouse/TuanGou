package wechat

import (
	"cooller/server/global"
	productModel "cooller/server/model/product"
	"time"
)

// Advertise 首页轮播广告表
type Advertise struct {
	global.GVA_MODEL
	Name       string    `json:"name" gorm:"not null;comment:名称;size:100"`                  // 套餐ID
	Type       int       `json:"type" form:"type" gorm:"comment:轮播位置：0->PC首页轮播；1->app首页轮播"` // 套餐名称
	Pic        string    `json:"pic" form:"pic" gorm:"comment:图片"`                          // 套餐类型
	StartTime  time.Time `json:"startTime" form:"startTime" gorm:"comment:开始时间"`            // 套餐价格
	EndTime    time.Time `json:"endTime" form:"endTime" gorm:"comment:结束时间"`                // 天数/次数/金额
	State      int       `json:"state" form:"state" gorm:"comment:上下线状态：0->下线；1->上线"`       // 状态
	ClickCount int       `json:"clickCount" form:"clickCount" gorm:"comment:点击数"`
	OrderCount int       `json:"orderCount" form:"orderCount" gorm:"comment:下单数"`
	Url        string    `json:"url" form:"url" gorm:"comment:链接地址;size:500"`
	Note       string    `json:"note" form:"note" gorm:"comment:备注;size:500"`
	Sort       int       `json:"sort" form:"sort" gorm:"comment:排序"`
	//SysUserAuthorityID int   `json:"sys_user_authority_id" form:"sys_user_authority_id" gorm:"comment:管理角色ID"`
}

func (Advertise) TableName() string {
	return "sms_home_advertise"
}

//// Brand 品牌表
//type Brand struct {
//	global.GVA_MODEL
//	Name                string `json:"name" gorm:"not null;comment:名称;size:100"`
//	FirstLetter         string `json:"firstLetter" form:"firstLetter" `
//	Sort                int    `json:"sort" form:"sort" gorm:"comment:排序"`
//	FactoryStatus       int    `json:"factoryStatus" form:"factoryStatus" gorm:"comment:是否是品牌制造商：0->不是；1->是"`
//	ShowStatus          int    `json:"showStatus" form:"showStatus" gorm:"comment:是否显示"`
//	ProductCount        int    `json:"productCount" form:"productCount" gorm:"comment:产品数量"`
//	ProductCommentCount int    `json:"productCommentCount" form:"productCommentCount" gorm:"comment:产品评论数量"`
//	Logo                string `json:"logo" gorm:"not null;comment:品牌logo"`
//	BigPic              string `json:"bigPic" gorm:"not null;comment:专区大图"`
//	BrandStory          string `json:"brandStory" gorm:"not null;comment:品牌故事"`
//	//ProductList         []*Product `json:"productList" gorm:"foreignKey:BrandId;"`
//	//SysUserAuthorityID int   `json:"sys_user_authority_id" form:"sys_user_authority_id" gorm:"comment:管理角色ID"`
//}
//
//func (Brand) TableName() string {
//	return "sms_home_brand"
//}

// FlashPromotion 限时购表
type FlashPromotion struct {
	global.GVA_MODEL
	Title     string    `json:"title" gorm:"not null;comment:名称;size:100"`
	StartDate time.Time `json:"startDate" form:"startDate" gorm:"comment:开始日期"`     // 套餐价格
	EndDate   time.Time `json:"endDate" form:"endDate" gorm:"comment:结束日期"`         // 天数/次数/金额
	Status    int       `json:"status" form:"status" gorm:"comment:状态：0->下线；1->上线"` // 状态
	//SysUserAuthorityID int   `json:"sys_user_authority_id" form:"sys_user_authority_id" gorm:"comment:管理角色ID"`
}

func (FlashPromotion) TableName() string {
	return "sms_flash_promotion"
}

// FlashPromotionProductRelation 商品限时购与商品关系表
type FlashPromotionProductRelation struct {
	global.GVA_MODEL
	FlashPromotionId        int                  `json:"flashPromotionId" gorm:"not null;comment:编号;size:100"`
	FlashPromotionSessionId int                  `json:"flashPromotionSessionId" form:"flashPromotionSessionId" gorm:"comment:编号"`
	ProductId               int                  `json:"productId" form:"productId" gorm:"comment:结束日期"`
	FlashPromotionPrice     float32              `json:"flashPromotionPrice" form:"flashPromotionPrice" gorm:"comment:限购价格"`
	FlashPromotionCount     int                  `json:"flashPromotionCount" form:"flashPromotionCount" gorm:"comment:限时购数量"`
	FlashPromotionLimit     int                  `json:"flashPromotionLimit" form:"flashPromotionLimit" gorm:"comment:每人限购数量"`
	Sort                    int                  `json:"sort" form:"sort" gorm:"comment:排序"`
	Product                 productModel.Product `json:"product" form:"product" gorm:"foreignKey:ProductId;comment:限时购商品"`
}

func (FlashPromotionProductRelation) TableName() string {
	return "sms_flash_promotion_product_relation"
}

// FlashPromotionSession 限时购场次表
type FlashPromotionSession struct {
	ID              int                              `json:"id" gorm:"primarykey"` // 主键ID
	Name            string                           `json:"name" gorm:"not null;comment:场次名称;size:100"`
	StartTime       time.Time                        `json:"startTime" form:"startTime" gorm:"comment:每日开始时间"`
	EndTime         time.Time                        `json:"endTime" form:"endTime" gorm:"comment:每日结束时间"`
	Status          int                              `json:"status" form:"status" gorm:"comment:启用状态：0->不启用；1->启用"` // 状态
	CreateDate      time.Time                        `json:"createDate" form:"createDate" gorm:"comment:创建时间"`
	ProductCount    int                              `json:"productCount" form:"productCount" gorm:"comment:商品数量"` // 状态
	ProductRelation []*FlashPromotionProductRelation `json:"productRelation" form:"productRelation" gorm:"foreignKey:FlashPromotionSessionId;comment:限购商品列表"`
}

func (FlashPromotionSession) TableName() string {
	return "sms_flash_promotion_session"
}

// NewProduct 新鲜好物表
type NewProduct struct {
	global.GVA_MODEL
	ProductId       int    `json:"productId" gorm:"not null;comment:物品序号"`
	ProductName     string `json:"productName" form:"productName" gorm:"comment:物品名称"`                     // 套餐价格
	RecommendStatus int    `json:"recommendStatus" form:"recommendStatus" gorm:"comment:推荐状态：0->下线；1->上线"` // 状态
	Sort            int    `json:"sort" form:"sort" gorm:"comment:排序"`
	//SysUserAuthorityID int    `json:"sys_user_authority_id" form:"sys_user_authority_id" gorm:"comment:管理角色ID"`
	Product productModel.Product `json:"product" gorm:"foreignKey:ProductId;"`
}

func (NewProduct) TableName() string {
	return "sms_new_product"
}

// RecommendProduct 人气推荐商品表
type RecommendProduct struct {
	global.GVA_MODEL
	ProductName     string `json:"productName" form:"productName" gorm:"comment:物品名称"`                     // 套餐价格
	RecommendStatus int    `json:"recommendStatus" form:"recommendStatus" gorm:"comment:推荐状态：0->下线；1->上线"` // 状态
	Sort            int    `json:"sort" form:"sort" gorm:"comment:排序"`
	//SysUserAuthorityID int   `json:"sys_user_authority_id" form:"sys_user_authority_id" gorm:"comment:管理角色ID"`
	ProductId int                  `json:"productId" gorm:"not null;comment:物品序号;size:100"`
	Product   productModel.Product `json:"product" gorm:"foreignKey:ProductId;"`
}

func (RecommendProduct) TableName() string {
	return "sms_home_recommend_product"
}

// GroupBuyProduct 精品商品团购列表
type GroupBuyProduct struct {
	global.GVA_MODEL
	ProductId int                  `json:"productId" gorm:"not null;comment:物品序号"`
	Price     float32              `json:"price" form:"price" gorm:"comment:商品价格"` // 状态
	Percent   int                  `json:"percent" form:"percent" gorm:"comment:团购进度条"`
	Sort      int                  `json:"sort" form:"sort" gorm:"comment:排序"`
	Status    int                  `json:"status" gorm:"not null；comment:显示状态：0->不显示；1->显示;size:1"`
	Product   productModel.Product `json:"product" gorm:"foreignKey:ProductId;"`
	//Type      int     `json:"type" gorm:"not null；comment:推荐类型；0->主推荐；1->附加推荐;size:1"`
}

func (GroupBuyProduct) TableName() string {
	return "sms_group_buy"
}
