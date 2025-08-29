package business

import "cooller/server/global"

type Goods struct {
	global.GVA_MODEL
	CardId              int     `json:"cardId" gorm:"comment:原主人会员ID"`
	Name                string  `json:"name" gorm:"not null;comment:名称;size:100"`
	Price               float32 `json:"price" form:"price" gorm:"comment:价格"`
	Pic                 string  `json:"pic" form:"pic" gorm:"comment:图片"`
	Description         string  `json:"description" form:"description" gorm:"comment:商品描述"`
	OriginalPrice       float32 `json:"originalPrice" form:"originalPrice" gorm:"comment:市场价"`
	AlbumPics           string  `json:"albumPics" form:"albumPics" gorm:"size:500;comment:画册图片(头图)，连产品图片限制为5张，以逗号分割"`
	DetailMobileHTML    string  `json:"detailMobileHTML" form:"detailMobileHTML" gorm:"type:text;comment:移动端网页详情"`
	Date                string  `json:"date" form:"date" gorm:"type:text;comment:上架日期"`
	PromotionStartDate  string  `json:"promotionStartDate" form:"promotionStartDate" gorm:"type:text;comment:置换开始日期"`
	PromotionEndDate    string  `json:"promotionEndDate" form:"promotionEndDate" gorm:"comment:置换结束日期"`
	ProductCategoryId   int     `json:"productCategoryId" `
	ProductCategoryName string  `json:"productCategoryName" form:"productCategoryName" gorm:"comment:商品分类名称"`
	Sort                int     `json:"sort" form:"sort" gorm:"comment:排序"`
	Keywords            string  `json:"keywords" form:"keywords" gorm:"comment:关键字"`
	Note                string  `json:"note" form:"note" gorm:"comment:备注"`
	DisplaceCardId      int     `json:"displaceCardId" gorm:"comment:置换者会员ID"`
}

func (Goods) TableName() string {
	return "bus_goods"
}
