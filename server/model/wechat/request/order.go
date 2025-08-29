package request

type OrderCreateRequest struct {
	AppId                  string `json:"appId" gorm:"null;default null"`
	OpenId                 string `json:"openId" gorm:"not null;" binding:"required"`
	IP                     string `json:"ip" gorm:"not null;" `
	Ids                    []int  `json:"ids" gorm:"not null;"`
	CouponId               int    `json:"couponId" gorm:"null;default null"`
	MemberReceiveAddressId int    `json:"memberReceiveAddressId" gorm:"not null;"`
	PayType                int    `json:"payType" gorm:"null;default null;comment:支付方式"`
	UseIntegration         int    `json:"useIntegration" gorm:"null;default null"`
	Note                   string `json:"note" gorm:"null;default null"`
	BuyType                int    `json:"buyType" gorm:"null;default null;comment:购买类型，1直接购买，2购物车购买"`
}

type PaySuccessRequest struct {
	OrderId int `json:"orderId" gorm:"not null;" binding:"required"`
	PayType int `json:"payType" gorm:"not null;"`
}

// OrderReceiveAddress 订单收货地址表
type OrderReceiveAddress struct {
	OrderId               int    `json:"orderId" `
	ReceiverName          string `json:"receiverName" gorm:"not null;comment:收货人名称;size:100"`
	ReceiverPhone         string `json:"receiverPhone" gorm:"not null；comment:手机号;size:11"`
	Status                int    `json:"status" gorm:"comment:是否默认;0->是；1->否"`
	ReceiverPostCode      string `json:"receiverPostCode" gorm:"comment:邮政编码"`
	ReceiverProvince      string `json:"receiverProvince" gorm:"comment:省份/直辖市"`
	ReceiverCity          string `json:"receiverCity" gorm:"not null；comment:城市"`
	ReceiverRegion        string `json:"receiverRegion" gorm:"not null；comment:区"`
	ReceiverDetailAddress string `json:"receiverDetailAddress" gorm:"comment:详细地址(楼层、门牌号)"`

	//SysUserAuthorityID uint   `json:"sys_user_authority_id" form:"sys_user_authority_id" gorm:"comment:管理角色ID"`
}

// OrderMoneyInfo 订单费用信息
type OrderMoneyInfo struct {
	OrderId        int     `json:"orderId" `
	DiscountAmount float32 `json:"discountAmount" gorm:"not null;comment:折扣金额;"`
}

// OrderNoteInfo 订单费用信息
type OrderNoteInfo struct {
	OrderId int    `json:"orderId" `
	Note    string `json:"note" gorm:"not null;comment:备注;"`
}
