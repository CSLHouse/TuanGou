package wechat

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/product"
	"cooller/server/model/wechat"
	wechatRequest "cooller/server/model/wechat/request"
	"fmt"
	"strings"
)

type HomeService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateHomeAdvertise
//@description: 创建套餐
//@param: e model.Advertise
//@return: err error

func (exa *HomeService) CreateHomeAdvertise(e wechat.Advertise) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteHomeAdvertise
//@description: 删除套餐
//@param: e model.Advertise
//@return: err error

func (exa *HomeService) DeleteHomeAdvertise(id int) (err error) {
	var brand wechat.Advertise
	err = global.GVA_DB.Where("id = ?", id).Delete(&brand).Error

	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateHomeAdvertise
//@description: 更新套餐
//@param: e *model.Advertise
//@return: err error

func (exa *HomeService) UpdateHomeAdvertise(e *wechat.Advertise) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *HomeService) UpdateHomeAdvertiseByIdForKeyword(e *wechatRequest.UpdateIdsKeywordRequest) (err error) {
	db := global.GVA_DB.Model(&wechat.Advertise{})
	db.Where("id IN ?", e.Ids).Updates(map[string]interface{}{e.Key: e.Value})
	return err
}

func (exa *HomeService) GetHomeAdvertiseById(id int) (advertise wechat.Advertise, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&advertise).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetHomeAdvertiseInfoList
//@description: 获取套餐列表
//@param: sysUserAuthorityID string, info request.PageInfo
//@return: list interface{}, total int64, err error

func (exa *HomeService) GetHomeAdvertiseInfoList(info request.PageInfo) (list []wechat.Advertise, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&wechat.Advertise{})
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Find(&list).Error
	}
	return list, total, err
}

func (exa *HomeService) GetOnlineHomeAdvertiseInfoList() (list []wechat.Advertise, err error) {
	err = global.GVA_DB.Model(&wechat.Advertise{}).Where("state = 1").Order("sort desc").Find(&list).Error
	return list, err
}

func (exa *HomeService) GetOnlineHomeBrandInfoList() (list []product.Brand, err error) {
	err = global.GVA_DB.Where("show_status = 1").Order("sort desc").Find(&list).Error
	return list, err
}

func (exa *HomeService) GetOnlineHomeFlashPromotionInfoList() (list []*wechat.FlashPromotion, err error) {
	err = global.GVA_DB.Where("status = 1").Find(&list).Error
	return list, err
}

func (exa *HomeService) CreateNewProduct(e *wechat.NewProduct) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) CreateNewProductByBatch(e []*wechat.NewProduct) (err error) {
	err = global.GVA_DB.CreateInBatches(e, len(e)).Error
	return err
}

func (exa *HomeService) UpdateNewProductForSort(e *request.SortUpdateInfo) (err error) {
	db := global.GVA_DB.Model(&wechat.NewProduct{})
	db.Where("id = ?", e.ID).UpdateColumn("sort", e.Sort)
	return err
}

func (exa *HomeService) UpdateNewProductForRecommendStatus(e *wechatRequest.UpdateIdsKeywordRequest) (err error) {
	db := global.GVA_DB.Model(&wechat.NewProduct{})
	db.Where("id IN ?", e.Ids).Updates(map[string]interface{}{e.Key: e.Value - 100})
	return err
}

func (exa *HomeService) DeleteNewProduct(ids []int) (err error) {
	var newProduct wechat.NewProduct
	err = global.GVA_DB.Where("id in ?", ids).Delete(&newProduct).Error
	return err
}

func (exa *HomeService) GetNewProductInfoList(info request.NameAndStateSearchInfo) (list []wechat.NewProduct, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&wechat.NewProduct{})
	var cmdName string
	var cmdStatus string

	if len(info.Keyword) > 0 {
		cmdName += fmt.Sprintf("product_name like '%%%s%%'", info.Keyword)
	}
	if info.State > 0 {
		cmdStatus += fmt.Sprintf("recommend_status = %d", info.State-100)
	}

	cmdSearch := ""
	cmd := [2]string{cmdName, cmdStatus}
	isFirst := true
	for _, cmdItem := range cmd {
		if len(cmdItem) > 0 {
			if isFirst {
				cmdSearch += cmdItem
				isFirst = false
			} else {
				cmdSearch += " and " + cmdItem
			}
		}
	}
	err = db.Where(cmdSearch).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Debug().Where(cmdSearch).Preload("Product").Limit(limit).Offset(offset).Find(&list).Error
	}
	return list, total, err
}
func (exa *HomeService) GetNewProductIdsList() (list []int, err error) {
	db := global.GVA_DB.Model(&wechat.NewProduct{})
	err = db.Select("product_id").Scan(&list).Error
	return list, err
}

func (exa *HomeService) GetOnlineNewProductInfoList() (list []wechat.NewProduct, err error) {
	err = global.GVA_DB.Where("recommend_status = 1").Preload("Product").Order("sort desc").Find(&list).Error
	return list, err
}

func (exa *HomeService) CreateRecommendProduct(e *[]wechat.RecommendProduct) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}
func (exa *HomeService) UpdateRecommendProducts(e *wechatRequest.UpdateIdsKeywordRequest) (err error) {
	db := global.GVA_DB.Model(&wechat.RecommendProduct{})
	db.Where("id IN ?", e.Ids).Updates(map[string]interface{}{e.Key: e.Value})
	return err
}

func (exa *HomeService) UpdateRecommendProductSortById(e *request.SortUpdateInfo) (err error) {
	db := global.GVA_DB.Model(&wechat.RecommendProduct{})
	db.Where("id = ?", e.ID).UpdateColumn("sort", e.Sort)
	return err
}

func (exa *HomeService) DeleteRecommendProducts(ids []int) (err error) {
	var product wechat.RecommendProduct
	err = global.GVA_DB.Where("id in ?", ids).Delete(&product).Error
	return err
}

func (exa *HomeService) GetOnlineRecommendProductListInfoList(pageInfo request.PageInfo) (recommendProductList []wechat.RecommendProduct, err error) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)

	err = global.GVA_DB.Limit(limit).Offset(offset).Where("recommend_status = 1").Preload("Product").Order("sort desc").Find(&recommendProductList).Error
	return recommendProductList, err
}

func (exa *HomeService) GetRecommendProductListByCondition(searchInfo wechatRequest.RecommendProductSearchInfo) (list []wechat.RecommendProduct, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	var cmdList []interface{}
	var cmdString string
	if len(searchInfo.ProductName) > 0 {
		cmdString = "product_name = ?"
		cmdList = append(cmdList, strings.TrimSpace(searchInfo.ProductName))
	}
	if searchInfo.RecommendStatus > 0 {
		if len(cmdList) >= 1 {
			cmdString += " and recommend_status = ?"
			cmdList = append(cmdList, searchInfo.RecommendStatus-100)
		} else {
			cmdString += "recommend_status = ?"
			cmdList = append(cmdList, searchInfo.RecommendStatus-100)
		}
	}

	db := global.GVA_DB.Model(&wechat.RecommendProduct{})
	err = db.Where(cmdString, cmdList...).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Where(cmdString, cmdList...).Order("sort desc").Find(&list).Error
	}

	return list, total, err
}

func (exa *HomeService) GetFlashPromotionList(searchInfo request.PageInfo) (list []wechat.FlashPromotion, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	var cmdName string

	if len(searchInfo.Keyword) > 0 {
		cmdName += fmt.Sprintf("title like '%%%s%%'", searchInfo.Keyword)
		db := global.GVA_DB.Model(&wechat.FlashPromotion{})
		err = db.Where(cmdName).Count(&total).Error
		if err != nil {
			return list, total, err
		} else {
			err = global.GVA_DB.Limit(limit).Offset(offset).Where(cmdName).Find(&list).Error
		}
	} else {
		db := global.GVA_DB.Model(&wechat.FlashPromotion{})
		err = db.Count(&total).Error
		if err != nil {
			return list, total, err
		} else {
			err = db.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, total, err
}

func (exa *HomeService) CreateFlashPromotion(e *wechat.FlashPromotion) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) UpdateFlashPromotion(e *wechat.FlashPromotion) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *HomeService) DeleteFlashPromotionById(id int) (err error) {
	var flash wechat.FlashPromotion
	err = global.GVA_DB.Where("id = ?", id).Delete(&flash).Error
	return err
}

func (exa *HomeService) UpdateFlashPromotionStatus(id int, status int) (err error) {
	db := global.GVA_DB.Model(&wechat.FlashPromotion{})
	db.Debug().Where("id = ?", id).UpdateColumn("status", status)
	return err
}

func (exa *HomeService) GetFlashPromotionProductRelationList(searchInfo wechatRequest.FlashProductRelationInfo) (list []wechat.FlashPromotionProductRelation, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&wechat.FlashPromotionProductRelation{})
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Where("flash_promotion_id = ? and flash_Promotion_session_id = ?", searchInfo.FlashPromotionId, searchInfo.FlashPromotionSessionId).Preload("Product").Find(&list).Error
	}

	return list, total, err
}

func (exa *HomeService) GetFlashPromotionProductRelationListById(flashPromotionId int, flashPromotionSessionId int) (list []wechat.FlashPromotionProductRelation, err error) {
	err = global.GVA_DB.Debug().Where("flash_promotion_id = ? and flash_promotion_session_id = ?", flashPromotionId, flashPromotionSessionId).Preload("Product").Find(&list).Error
	return list, err
}

func (exa *HomeService) CreateFlashPromotionProductRelation(e []wechat.FlashPromotionProductRelation) (err error) {
	err = global.GVA_DB.CreateInBatches(e, len(e)).Error
	return err
}

func (exa *HomeService) UpdateProductPromotionType(relation *wechat.FlashPromotionProductRelation, session *wechat.FlashPromotionSession) (err error) {
	db := global.GVA_DB.Model(&product.Product{})
	err = db.Where("id=?", relation.ProductId).UpdateColumn("promotion_type", 5).Error
	return err
}

func (exa *HomeService) UpdateFlashPromotionProductRelation(e *wechat.FlashPromotionProductRelation) (err error) {
	err = global.GVA_DB.Debug().Save(e).Error
	return err
}

func (exa *HomeService) DeleteFlashPromotionProductRelationById(id int) (err error) {
	var flash wechat.FlashPromotionProductRelation
	err = global.GVA_DB.Where("id = ?", id).Delete(&flash).Error
	return err
}

func (exa *HomeService) GetFlashSessionList() (list []*wechat.FlashPromotionSession, err error) {
	err = global.GVA_DB.Find(&list).Error
	return list, err
}

func (exa *HomeService) GetFlashSessionById(id int) (flashSession wechat.FlashPromotionSession, err error) {
	db := global.GVA_DB.Model(&wechat.FlashPromotionSession{})
	err = db.Where("id = ?", id).First(&flashSession).Error
	return flashSession, err
}

func (exa *HomeService) CreateFlashSession(e *wechat.FlashPromotionSession) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) UpdateFlashSession(e *wechat.FlashPromotionSession) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *HomeService) DeleteFlashSessionById(id int) (err error) {
	var flash wechat.FlashPromotionSession
	err = global.GVA_DB.Where("id = ?", id).Delete(&flash).Error
	return err
}

func (exa *HomeService) UpdateFlashSessionStatus(id int, status int) (err error) {
	db := global.GVA_DB.Model(&wechat.FlashPromotionSession{})
	db.Where("id = ?", id).UpdateColumn("status", status)
	return err
}

func (exa *HomeService) GetFlashSessionSelectList(id int) (list []*wechat.FlashPromotionSession, err error) {
	err = global.GVA_DB.Preload("ProductRelation").Find(&list).Error
	return list, err
}

func (exa *HomeService) GetGroupBuyProductList() (list []*wechat.GroupBuyProduct, err error) {
	err = global.GVA_DB.Preload("Product").Find(&list).Error
	return list, err
}
