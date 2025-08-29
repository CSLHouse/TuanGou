package product

import (
	"cooller/server/global"
	"time"
)

type Coupon struct {
	global.GVA_MODEL
	Name         string    `json:"name" gorm:"not null;comment:名称;size:100"`                              // 套餐ID
	Type         int       `json:"type" form:"type" gorm:"comment:优惠券类型；0->全场赠券；1->会员赠券；2->购物赠券；3->注册赠券"` // 套餐名称
	Platform     int       `json:"platform" form:"platform" gorm:"comment:使用平台：0->全部；1->移动；2->PC;size:1"` // 套餐类型
	Count        int       `json:"count" form:"count" gorm:"comment:数量;size:11"`                          // 套餐价格
	Amount       int       `json:"amount" form:"amount" gorm:"comment:金额"`                                // 天数/次数/金额
	PerLimit     int       `json:"perLimit" form:"perLimit" gorm:"comment:每人限领张数;size:11"`                // 状态
	MinPoint     int       `json:"minPoint" form:"minPoint" gorm:"comment:使用门槛；0表示无门槛;size:1"`
	StartTime    time.Time `json:"startTime" form:"startTime" gorm:"comment:开始时间"` // 套餐价格
	EndTime      time.Time `json:"endTime" form:"endTime" gorm:"comment:结束时间"`
	UseType      int       `json:"useType" form:"useType" gorm:"comment:使用类型：0->全场通用；1->指定分类；2->指定商品;size:1"`
	Note         string    `json:"note" form:"note" gorm:"comment:备注;size:500"`
	PublishCount string    `json:"publishCount" form:"publishCount" gorm:"comment:发行数量;size:11"`
	UseCount     int       `json:"useCount" form:"useCount" gorm:"comment:已使用数量;size:11"`
	ReceiveCount int       `json:"receiveCount" form:"receiveCount" gorm:"comment:领取数量"`
	EnableTime   time.Time `json:"enableTime" form:"enableTime" gorm:"comment:可以领取的日期"`
	Code         string    `json:"code" form:"code" gorm:"comment:优惠码;"`
	MemberLevel  int       `json:"memberLevel" form:"memberLevel" gorm:"comment:可领取的会员类型：0->无限时;size:1"`
}

func (Coupon) TableName() string {
	return "sms_coupon"
}

type CouponHistory struct {
	global.GVA_MODEL
	CouponId       int       `json:"couponId" form:"couponId" gorm:"not null;comment:名称;"`               // 套餐ID
	MemberId       int       `json:"memberId" form:"memberId" gorm:"not null;"`                          // 套餐名称
	CouponCode     string    `json:"couponCode" form:"couponCode" gorm:"not null"`                       // 套餐类型
	MemberNickname string    `json:"memberNickname" form:"memberNickname" gorm:"comment:领取人昵称;not null"` // 套餐价格
	GetType        int       `json:"getType" form:"getType" gorm:"comment:获取类型：0->后台赠送；1->主动获取;size:1"`  // 天数/次数/金额
	UseStatus      int       `json:"useStatus" form:"useStatus" gorm:"comment:使用状态：0->未使用；1->已使用；2->已过期;size:1"`
	UseTime        time.Time `json:"useTime" form:"useTime" gorm:"comment:使用时间"` // 套餐价格
	OrderId        int       `json:"orderId" form:"orderId" gorm:"comment:订单编号"`
	OrderSn        string    `json:"orderSn" form:"orderSn" gorm:"comment:订单号码;"`
}

func (CouponHistory) TableName() string {
	return "sms_coupon_history"
}
