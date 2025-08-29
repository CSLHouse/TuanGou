package response

import (
	"cooller/server/model/wechat"
	"time"
)

type CartPromotionItem struct {
	ID        int       `json:"id"`
	CreatedAt time.Time // 创建时间
	//ProductId         int       `json:"productId" gorm:"null;default null"`
	//ProductSkuId      int       `json:"productSkuId" gorm:"null;default null;comment:商品sku编号;"`
	//UserId            int       `json:"user_id" gorm:" not null;"`
	Quantity    int     `json:"quantity" gorm:"null;default null;comment:购买数量;"`
	Price       float32 `json:"price" gorm:"null;default null;comment:销售价格;"`
	ProductPic  string  `json:"productPic" gorm:"null;default null;comment:商品主图;"`
	ProductName string  `json:"productName" gorm:"null;default null;comment:商品名称;"`
	//ProductSubTitle   string    `json:"productSubTitle" gorm:"null;default null;comment:商品副标题（卖点）;"`
	//ProductSkuCode    string    `json:"productSkuCode" gorm:"null;default null;comment:商品sku条码;"`
	//MemberNickname    string    `json:"memberNickname" gorm:"null;default null;comment:会员昵称;"`
	//DeleteStatus      int       `json:"deleteStatus" gorm:"null;default null;comment:是否删除;"`
	//ProductCategoryId int       `json:"productCategoryId" gorm:"null;default null;comment:商品分类;"`
	//ProductBrand      string    `json:"productBrand" gorm:"null;default null;comment:品牌;"`
	//ProductSn         string    `json:"productSn" gorm:"null;default null;"`
	//ProductAttr       string    `json:"productAttr" gorm:"null;default null;comment:商品销售属性:[{\"key\":\"颜色\",\"value\":\"颜色\"},{\"key\":\"容量\",\"value\":\"4G\"}];"`
	PromotionMessage string  `json:"promotionMessage" gorm:"null;default null;comment:优惠信息;"`
	ReduceAmount     float32 `json:"reduceAmount" gorm:"null;default null;comment:优惠金额;"`
	RealStock        int     `json:"realStock" gorm:"null;default null;comment:剩余库存;"`
	Integration      int     `json:"integration" gorm:"null;default null;comment:积分;"`
	Growth           int     `json:"growth" gorm:"null;default null;"`
}

type CalcAmount struct {
	TotalAmount     float32 `json:"totalAmount" gorm:"null;default null;comment:销售总价格;"`
	FreightAmount   float32 `json:"freightAmount" gorm:"null;default null;comment:运费;"`
	PromotionAmount float32 `json:"promotionAmount" gorm:"null;default null;comment:优惠金额;"`
	PayAmount       float32 `json:"payAmount" gorm:"null;default null;comment:实际支付金额;"`
}

type GenerateOrderResModel struct {
	CartPromotionItemList    []*CartPromotionItem          `json:"cartPromotionItemList" gorm:"null;default null"`
	MemberReceiveAddressList []wechat.MemberReceiveAddress `json:"memberReceiveAddressList" gorm:"null;default null"`
	//CouponHistoryDetailList   []int                         `json:"couponHistoryDetailList" gorm:"null;default null"`
	//IntegrationConsumeSetting int                           `json:"integrationConsumeSetting" gorm:"null;default null"`
	//MemberIntegration         int                           `json:"memberIntegration" gorm:"null;default null"`
	CalcAmount CalcAmount `json:"calcAmount" gorm:"null;default null"`
	PickupType int        `json:"pickupType" gorm:"null;default null;comment:是否包含不自提的物品：0不包含，1包含;"`
}
