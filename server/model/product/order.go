package product

import (
	"cooller/server/global"
	"time"
)

// Order 订单表
type Order struct {
	global.GVA_MODEL
	UserId            int     `json:"user_id" gorm:" not null;"`
	CouponId          int     `json:"couponId" gorm:" not null;"`
	OrderSn           string  `json:"orderSn" gorm:"null;default null;comment:订单编号;"`
	UserName          string  `json:"userName" gorm:"null;default null;comment:用户帐号;"`
	TotalAmount       float32 `json:"totalAmount" gorm:"null;default null;comment:订单总金额;"`
	PayAmount         float32 `json:"payAmount" gorm:"null;default null;comment:应付金额（实际支付金额）;"`
	FreightAmount     float32 `json:"freightAmount" gorm:"null;default null;comment:运费金额;"`
	PromotionAmount   float32 `json:"promotionAmount" gorm:"null;default null;comment:促销优化金额（促销价、满减、阶梯价）"`
	IntegrationAmount float32 `json:"integrationAmount" gorm:"null;default null;comment:积分抵扣金额;"`
	CouponAmount      float32 `json:"couponAmount" gorm:"null;default null;comment:优惠券抵扣金额;"`
	DiscountAmount    float32 `json:"discountAmount" gorm:"null;default null;comment:管理员后台调整订单使用的折扣金额;"`
	PayType           int     `json:"payType" gorm:"null;default null;comment:支付方式：0->未支付；1->支付宝；2->微信;size:1;"`
	SourceType        int     `json:"sourceType" gorm:"null;default null;comment:订单来源：0->PC订单；1->app订单;"`
	Status            int     `json:"status" gorm:"null;default null;comment:订单状态：0->待付款；1->待发货；2->已发货；3->已完成；4->已关闭；5->无效订单;"`
	OrderType         int     `json:"orderType" gorm:"null;default null;comment:订单类型：0->正常订单；1->秒杀订单;"`
	LogisticsCompany  string  `json:"logisticsCompany" gorm:"null;default null;comment:物流公司(配送方式);"`
	LogisticsSn       string  `json:"logisticsSn" gorm:"null;default null;comment:物流单号;"`
	AutoConfirmDay    int     `json:"autoConfirmDay" gorm:"null;default null;comment:自动确认时间（天）;size:11;"`
	Integration       int     `json:"integration" gorm:"null;default null;comment:可以获得的积分;size:1;"`
	Growth            int     `json:"growth" gorm:"null;default null;comment:可以活动的成长值;size:1;"`
	PromotionInfo     string  `json:"promotionInfo" gorm:"null;default '';comment:活动信息;"`
	BillType          int     `json:"billType" gorm:"null;default null;size:1;comment:发票类型：0->不开发票；1->电子发票；2->纸质发票;"`
	BillHeader        string  `json:"billHeader" gorm:"null;default null;comment:发票抬头"`
	BillContent       string  `json:"billContent" gorm:"null;default null;comment:发票内容"`
	BillReceiverPhone string  `json:"billReceiverPhone" gorm:"null;default null;comment:收票人电话"`
	BillReceiverEmail string  `json:"billReceiverEmail" gorm:"null;default null;comment:收票人邮箱"`
	//ReceiverPhone         string       `json:"receiverPhone" gorm:"null;default null;comment:收货人电话"`
	//ReceiverName          string       `json:"receiverName" gorm:"null;default null;comment:收货人姓名"`
	//ReceiverPostCode      string       `json:"receiverPostCode" gorm:"null;default null;comment:收货人邮编"`
	//ReceiverProvince      string       `json:"receiverProvince" gorm:"null;default null;comment:省份/直辖市"`
	//ReceiverCity          string       `json:"receiverCity" gorm:"null;default null;comment:城市"`
	//ReceiverRegion        string       `json:"receiverRegion" gorm:"null;default null;comment:区"`
	//ReceiverDetailAddress string       `json:"receiverDetailAddress" gorm:"null;default null;comment:详细地址"`
	ReceiveAddressId int       `json:"memberReceiveAddressId" gorm:"not null;"`
	Note             string    `json:"note" gorm:"null;default null;comment:订单备注"`
	ConfirmStatus    int       `json:"confirmStatus" gorm:"null;default null;comment:确认收货状态：0->未确认；1->已确认;size:1;"`
	DeleteStatus     int       `json:"deleteStatus" gorm:"null;default null;comment:删除状态：0->未删除；1->已删除;size:1;"`
	UseIntegration   int       `json:"useIntegration" gorm:"null;default null;comment:下单时使用的积分;size:11;"`
	PaymentTime      time.Time `json:"paymentTime" gorm:"null;default null;comment:支付时间"`
	LogisticsTime    time.Time `json:"logisticsTime" gorm:"null;default 1000-01-01 00:00:00;comment:发货时间"`
	ReceiveTime      time.Time `json:"receiveTime" gorm:"null;default 1000-01-01 00:00:00;comment:确认收货时间"`
	//CommentTime      time.Time `json:"commentTime" gorm:"null;default 1000-01-01 00:00:00;comment:评价时间"`
	//ModifyTime       time.Time    `json:"modifyTime" gorm:"null;default 1000-01-01 00:00:00;comment:修改时间"`
	OrderItemList []*OrderItem `json:"orderItemList" gorm:"foreignKey:OrderId"`
	PrepayId      string       `json:"prepayId" gorm:"null;default null;comment:预支付交易会话标识"`
}

func (Order) TableName() string {
	return "oms_order"
}

// OrderItem 订单中所包含的商品
type OrderItem struct {
	global.GVA_MODEL
	OrderId           int     `json:"orderId" gorm:"not null;default null;comment:订单id;references:ID"`
	OrderSn           string  `json:"orderSn" gorm:"null;default null;comment:订单编号;"`
	ProductId         int     `json:"productId" gorm:"null;default null"`
	ProductSkuId      string  `json:"productSkuId" gorm:"null;default null;comment:商品sku编号;"`
	UserId            int     `json:"user_id" gorm:" not null;"`
	Quantity          int     `json:"quantity" gorm:"null;default null;comment:购买数量;"`
	Price             float32 `json:"price" gorm:"null;default null;comment:销售价格;"`
	ProductPic        string  `json:"productPic" gorm:"null;default null;comment:商品主图;"`
	ProductName       string  `json:"productName" gorm:"null;default null;comment:商品名称;"`
	ProductSubTitle   string  `json:"productSubTitle" gorm:"null;default null;comment:商品副标题（卖点）;"`
	ProductSkuCode    string  `json:"productSkuCode" gorm:"null;default null;comment:商品sku条码;"`
	MemberNickname    string  `json:"memberNickname" gorm:"null;default null;comment:会员昵称;"`
	DeleteStatus      int     `json:"deleteStatus" gorm:"null;default null;comment:是否删除;"`
	ProductCategoryId int     `json:"productCategoryId" gorm:"null;default null;comment:商品分类;"`
	ProductBrand      string  `json:"productBrand" gorm:"null;default null;comment:品牌;"`
	ProductSN         string  `json:"productSN" gorm:"null;default null;"`
	ProductAttr       string  `json:"productAttr" gorm:"null;default null;comment:商品销售属性:[{\"key\":\"颜色\",\"value\":\"颜色\"},{\"key\":\"容量\",\"value\":\"4G\"}];"`
	PromotionName     string  `json:"promotionName" gorm:"null;default null;comment:商品促销名称;"`
	PromotionAmount   float32 `json:"promotionAmount" gorm:"null;default null;comment:商品促销分解金额;"`
	PromotionMessage  string  `json:"promotionMessage" gorm:"-"` // 促销活动信息
	CouponAmount      float32 `json:"couponAmount" gorm:"null;default null;comment:优惠券优惠分解金额;"`
	IntegrationAmount float32 `json:"integrationAmount" gorm:"null;default null;comment:积分优惠分解金额;"`
	ReduceAmount      float32 `json:"reduceAmount" gorm:"-"` // 促销活动减去的金额，针对每个商品
	RealAmount        float32 `json:"realAmount" gorm:"null;default null;comment:促销活动减去的金额，针对每个商品;"`
	RealStock         int     `json:"realStock" gorm:"-"` // 剩余库存-锁定库存
	GiftIntegration   int     `json:"giftIntegration" gorm:"null;default null;comment:购买商品赠送积分;"`
	GiftGrowth        int     `json:"giftGrowth" gorm:"null;default null;comment:购买商品赠送成长值;"`
}

func (OrderItem) TableName() string {
	return "oms_order_item"
}

type OrderSetting struct {
	global.GVA_MODEL
	FlashOrderOvertime  int `json:"flashOrderOvertime" gorm:"null;default null;comment:秒杀订单超时关闭时间(分);"`
	NormalOrderOvertime int `json:"normalOrderOvertime" gorm:"null;default null;comment:常订单超时时间(分);"`
}

func (OrderSetting) TableName() string {
	return "oms_order_setting"
}

// CartItem 购物车表
type CartItem struct {
	global.GVA_MODEL
	ProductId  int      `json:"productId" gorm:"null;default null"`
	SkuStockId int      `json:"skuStockId" gorm:"null;default null;"`
	UserId     int      `json:"user_id" gorm:" not null;"`
	Quantity   int      `json:"quantity" gorm:"null;default null;comment:购买数量;"`
	Product    Product  `json:"product" gorm:"foreignKey:ProductId;"`
	SkuStock   SkuStock `json:"skuStock" gorm:"foreignKey:SkuStockId;"`
	Price      float32  `json:"price" gorm:"null;default null;comment:添加到购物车的价格;"`
}

func (CartItem) TableName() string {
	return "oms_cart_item"
}

// CartTmpItem 直接购买的虚拟购物车表
type CartTmpItem struct {
	global.GVA_MODEL
	ProductId  int      `json:"productId" gorm:"null;default null"`
	SkuStockId int      `json:"skuStockId" gorm:"null;default null;"`
	UserId     int      `json:"user_id" gorm:" not null;"`
	Quantity   int      `json:"quantity" gorm:"null;default null;comment:购买数量;"`
	Product    Product  `json:"product" gorm:"foreignKey:ProductId;"`
	SkuStock   SkuStock `json:"skuStock" gorm:"foreignKey:SkuStockId;"`
	Price      float32  `json:"price" gorm:"null;default null;comment:添加到购物车的价格;"`
}

func (CartTmpItem) TableName() string {
	return "oms_cart_tmp_item"
}

type CartCommonItem struct {
	global.GVA_MODEL
	ProductId  int      `json:"productId" gorm:"null;default null"`
	SkuStockId int      `json:"skuStockId" gorm:"null;default null;"`
	UserId     int      `json:"user_id" gorm:" not null;"`
	Quantity   int      `json:"quantity" gorm:"null;default null;comment:购买数量;"`
	Product    Product  `json:"product" gorm:"foreignKey:ProductId;"`
	SkuStock   SkuStock `json:"skuStock" gorm:"foreignKey:SkuStockId;"`
	Price      float32  `json:"price" gorm:"null;default null;comment:添加到购物车的价格;"`
	//CouponAmount float32  `json:"couponAmount"` // 优惠券优惠分解金额
}

//type CartPromotionItem struct {
//	OrderItem        `json:",inline"`
//	PromotionMessage string  `json:"promotionMessage"` // 促销活动信息
//	ReduceAmount     float32 `json:"reduceAmount"`     // 促销活动减去的金额，针对每个商品
//	RealStock        int     `json:"realStock"`        // 剩余库存-锁定库存
//	Integration      int     `json:"integration"`      // 购买商品赠送积分
//	Growth           int     `json:"growth"`           // 购买商品赠送成长值
//}
