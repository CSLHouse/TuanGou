package response

import (
	"cooller/server/global"
	"cooller/server/model/wechat"
	wechatReq "cooller/server/model/wechat/request"
	"time"
)

type FlashPromotionResModel struct {
	NextStartTime int         `json:"nextStartTime"`
	EndTime       int         `json:"endTime"`
	ProductList   interface{} `json:"productList"`
}

type HomeContentResponse struct {
	AdvertiseList      interface{} `json:"advertiseList"`
	BrandList          interface{} `json:"brandList"`
	HomeFlashPromotion interface{} `json:"homeFlashPromotion"`
	NewProductList     interface{} `json:"newProductList"`
	HotProductList     interface{} `json:"hotProductList"`
	GroupBuy           interface{} `json:"groupBuy"`
}

type HomeRecommendResponse struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Pic      string  `json:"pic"`
	SubTitle string  `json:"subTitle"`
	Price    float32 `json:"price"`
}

type HomeFlashResponse struct {
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	NextStartTime time.Time `json:"nextStartTime"`
	NextEndTime   time.Time `json:"nextEndTime"`
	//HotProductList interface{} `json:"hotProductList"`
	ProductList []wechat.Product `json:"productList"`
}

// SampleProductInfo 简单商品信息
type SampleProductInfo struct {
	ID            int     `json:"id"`
	Name          string  `json:"name" form:"name" gorm:"comment:物品名称"`
	Price         float32 `json:"price" form:"price" gorm:"comment:商品价格"` // 状态
	OriginalPrice float32 `json:"originalPrice" form:"originalPrice" gorm:"comment:市场价"`
	Sales         int     `json:"sales" form:"sales" gorm:"comment:销量"`
	AlbumPics     string  `json:"albumPics" form:"albumPics" gorm:"size:500;comment:画册图片(头图)，连产品图片限制为5张，以逗号分割"`
	Pic           string  `json:"pic" form:"pic" gorm:"comment:图片"`
	Percent       int     `json:"percent" form:"percent" gorm:"comment:团购进度条"`
	ProductId     int     `json:"productId" gorm:"not null;comment:物品序号"`
}

type GroupBuyResp struct {
	Groups []SampleProductInfo `json:"groups" form:"groups" gorm:"comment:团购值"`
}

//type GroupBuyInfo struct {
//	MainProduct SampleProductInfo `json:"mainProduct" form:"mainProduct" gorm:"comment:主推荐"`
//	SubProduct  SampleProductInfo `json:"subProduct" form:"subProduct" gorm:"comment:附加推荐"`
//}

type ProductDetails struct {
	Product    wechat.Product       `json:"product" form:"product" gorm:"comment:商品信息"`
	BuyersList wechatReq.BuyersInfo `json:"buyersList" form:"buyersList" gorm:"comment:买家信息"`
}

type NewProductRes struct {
	global.GVA_MODEL
	ProductId       int    `json:"productId" gorm:"not null;comment:物品序号"`
	ProductName     string `json:"productName" form:"productName" gorm:"comment:物品名称"`                     // 套餐价格
	RecommendStatus int    `json:"recommendStatus" form:"recommendStatus" gorm:"comment:推荐状态：0->下线；1->上线"` // 状态
	Sort            int    `json:"sort" form:"sort" gorm:"comment:排序"`
}

type NewProductIdsList struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	ProductIds []int       `json:"productIds"`
}
