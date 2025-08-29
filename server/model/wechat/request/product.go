package request

import (
	"cooller/server/model/wechat"
)

type UpdateIdsKeywordRequest struct {
	Ids   []int  `json:"ids" gorm:"not null;comment:物品序号"`
	Key   string `json:"key" gorm:"not null；comment:各种状态：上架、新品、推荐"`
	Value int    `json:"value" gorm:"not null；comment:状态值:0->不是；1->是"`
}

type AddRecommendProductRequest struct {
	Products []wechat.RecommendProduct `json:"products" gorm:"not null"`
}

type MemberLevel struct {
	MemberLevelId   int     `json:"memberLevelId" gorm:"not null"`
	MemberLevelName string  `json:"memberLevelName" gorm:"not null"`
	MemberPrice     float32 `json:"memberPrice" `
}

type ProductSearchInfo struct {
	Keyword           string `json:"keyword" form:"keyword" gorm:"comment:商品名称"`
	BrandId           int    `json:"brandId" form:"brandId" gorm:"comment:物品序号"`
	ProductSN         string `json:"productSN" form:"productSN" gorm:"comment:货号"`
	ProductCategoryId string `json:"productCategoryId" form:"productCategoryId" gorm:"comment:商品分类"`
	BrandName         string `json:"brandName" form:"brandName" gorm:"comment:品牌"`
	PublishStatus     int    `json:"publishStatus" form:"publishStatus" gorm:"comment:上架状态;"`
	VerifyStatus      int    `json:"verifyStatus" form:"verifyStatus" gorm:"comment:审核状态;"`
	Page              int    `json:"page" form:"page"`         // 页码
	PageSize          int    `json:"pageSize" form:"pageSize"` // 每页大小
}

type BrandSearchInfo struct {
	Name       string `json:"name" form:"name" gorm:"comment:品牌名称"`
	ShowStatus int    `json:"showStatus" form:"showStatus" gorm:"comment:推荐状态"`
	Page       int    `json:"page" form:"page"`         // 页码
	PageSize   int    `json:"pageSize" form:"pageSize"` // 每页大小
}

type FlashProductRelationInfo struct {
	FlashPromotionId        int `json:"flashPromotionId" form:"flashPromotionId" gorm:"comment:品牌名称"`
	FlashPromotionSessionId int `json:"flashPromotionSessionId" form:"flashPromotionSessionId" gorm:"comment:推荐状态"`
	Page                    int `json:"page" form:"page"`         // 页码
	PageSize                int `json:"pageSize" form:"pageSize"` // 每页大小
}

type RecommendProductSearchInfo struct {
	ProductName     string `json:"productName" form:"productName" gorm:"comment:商品名称"`
	RecommendStatus int    `json:"recommendStatus" form:"recommendStatus" gorm:"comment:推荐状态"`
	Page            int    `json:"page" form:"page"`         // 页码
	PageSize        int    `json:"pageSize" form:"pageSize"` // 每页大小
}

type BuyersInfo struct {
	Avatars []string `json:"avatars" form:"avatars" gorm:"comment:微信头像url"`
}

type NewProductsRequest struct {
	NewProducts []*wechat.NewProduct `json:"newProducts"`
}
