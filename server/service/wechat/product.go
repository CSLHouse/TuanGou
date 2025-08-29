package wechat

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/wechat"
	wechatRequest "cooller/server/model/wechat/request"
	wechatRes "cooller/server/model/wechat/response"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
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
	db := global.GVA_DB.Model(&wechat.Advertise{})
	err = db.Where("state = 1").Order("sort desc").Find(&list).Error
	return list, err
}

func (exa *HomeService) GetOnlineHomeBrandInfoList() (list []wechat.Brand, err error) {
	db := global.GVA_DB.Model(&wechat.Brand{})
	err = db.Where("show_status = 1").Order("sort desc").Find(&list).Error
	return list, err
}

func (exa *HomeService) GetOnlineHomeFlashPromotionInfoList() (list []*wechat.FlashPromotion, err error) {
	db := global.GVA_DB.Model(&wechat.FlashPromotion{})
	err = db.Where("status = 1").Find(&list).Error
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

func (exa *HomeService) GetNewProductInfoList(info request.NameAndStateSearchInfo) (list []wechatRes.NewProductRes, total int64, err error) {
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
		err = db.Debug().Where(cmdSearch).Omit("Product").Limit(limit).Offset(offset).Find(&list).Error
	}
	return list, total, err
}
func (exa *HomeService) GetNewProductIdsList() (list []int, err error) {
	db := global.GVA_DB.Model(&wechat.NewProduct{})
	err = db.Select("product_id").Scan(&list).Error
	return list, err
}

func (exa *HomeService) GetOnlineNewProductInfoList() (list []wechat.NewProduct, err error) {
	db := global.GVA_DB.Model(&wechat.NewProduct{})
	err = db.Where("recommend_status = 1").Preload("Product").Order("sort desc").Find(&list).Error
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

	db := global.GVA_DB.Model(&wechat.RecommendProduct{})

	err = db.Limit(limit).Offset(offset).Where("recommend_status = 1").Preload("Product").Order("sort desc").Find(&recommendProductList).Error
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

func (exa *HomeService) CreateHomeProduct(e *wechat.Product) (err error) {
	// TODO: 会员价格
	// 创建产品参数信
	for _, attribute := range e.ProductAttributeValueList {
		if len(attribute.Value) > 0 {
			err = exa.CreateProductAttributeValue(attribute)
			if err != nil {
				return err
			}
		}
	}
	//产品满减
	for _, fullReduction := range e.ProductFullReductionList {
		if fullReduction.FullPrice > 0 && fullReduction.ReducePrice > 0 {
			err = exa.CreateProductFullReduction(fullReduction)
			if err != nil {
				return err
			}
		}
	}
	// 产品阶梯价格
	for _, ladder := range e.ProductLadderList {
		if ladder.Count > 0 && ladder.Discount > 0 && ladder.Price > 0 {
			err = exa.CreateProductLadder(ladder)
			if err != nil {
				return err
			}
		}
	}
	// sku的库存
	for _, stock := range e.SkuStockList {
		if stock.Price > 0 && stock.PromotionPrice > 0 {
			stock.SkuCode = fmt.Sprintf("%d", time.Now().UnixNano())
			err = exa.CreateProductSKUStock(stock)
			if err != nil {
				return err
			}
		}
	}
	err = global.GVA_DB.Create(e).Error
	return err
}

func (exa *HomeService) UpdateHomeProductSynchronous(e *wechat.Product) (err error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return err
	}

	// 更新Product同时更新ProductLadder
	err = exa.UpdateProductLadderWithProduct(e.ID, e.ProductLadderList, tx)
	if err != nil {
		return err
	}
	// 更新Product同时更新满减
	err = exa.UpdateProductFullReductionWithProduct(e.ID, e.ProductFullReductionList, tx)
	if err != nil {
		return err
	}

	// 更新Product同时更新sku
	err = exa.UpdateSKUStockWithProduct(e.ID, e.SkuStockList, tx)
	if err != nil {
		return err
	}

	// 2. 显式更新关联的子模型（关键步骤）
	err = exa.UpdateProductAttributeValueWithProduct(e.ID, e.ProductAttributeValueList, tx)
	if err != nil {
		return err
	}

	return tx.Commit().Error
}

// UpdateSKUStockWithProduct 更新Product同时更新sku
func (exa *HomeService) UpdateSKUStockWithProduct(productId int, skus []*wechat.SkuStock, tx *gorm.DB) (err error) {
	// 3. 处理关联的产品SKU列表（使用ON DUPLICATE KEY UPDATE优化）
	if len(skus) > 0 {
		// 3.1 确保外键和时间戳正确
		now := time.Now()
		for _, sku := range skus {
			sku.ProductId = productId // 强制关联当前产品
			sku.UpdatedAt = now       // 更新时间戳
			if sku.CreatedAt.IsZero() {
				sku.CreatedAt = now // 新增记录时设置创建时间
			}
		}

		// 3.2 批量插入或更新关联属性
		// 使用ON DUPLICATE KEY UPDATE确保所有需要更新的字段被处理
		if err = tx.Debug().Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}}, // 按主键ID冲突判断
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at":      clause.Expr{SQL: "VALUES(updated_at)"},
				"product_id":      clause.Expr{SQL: "VALUES(product_id)"},
				"sku_code":        clause.Expr{SQL: "VALUES(sku_code)"},
				"price":           clause.Expr{SQL: "VALUES(price)"},
				"stock":           clause.Expr{SQL: "VALUES(stock)"},
				"low_stock":       clause.Expr{SQL: "VALUES(low_stock)"},
				"pic":             clause.Expr{SQL: "VALUES(pic)"},
				"promotion_price": clause.Expr{SQL: "VALUES(promotion_price)"},
				"sp_data":         clause.Expr{SQL: "VALUES(sp_data)"},
				// 如有其他需要更新的字段，在此处添加
			}),
		}).Create(&skus).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 3.3 清理数据库中存在但不在当前列表中的属性值
		// 收集当前列表中的所有属性ID
		var skuStockIDs []int
		for _, attr := range skus {
			if attr.ID > 0 {
				skuStockIDs = append(skuStockIDs, attr.ID)
			}
		}

		// 删除不在当前列表中的属性值（保持数据一致性）
		if len(skuStockIDs) > 0 {
			if err = tx.Where("product_id = ? AND id NOT IN (?)", productId, skuStockIDs).
				Delete(&wechat.SkuStock{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 如果当前列表没有带ID的属性，删除该产品所有属性
			if err = tx.Where("product_id = ?", productId).
				Delete(&wechat.SkuStock{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

// UpdateProductAttributeValueWithProduct 更新Product同时更新属性值
func (exa *HomeService) UpdateProductAttributeValueWithProduct(productId int, attrs []*wechat.ProductAttributeValue, tx *gorm.DB) (err error) {
	// 3. 处理关联的产品属性值列表（使用ON DUPLICATE KEY UPDATE优化）
	if len(attrs) > 0 {
		// 3.1 确保外键和时间戳正确
		now := time.Now()
		for _, attr := range attrs {
			attr.ProductId = productId // 强制关联当前产品
			attr.UpdatedAt = now       // 更新时间戳
			if attr.CreatedAt.IsZero() {
				attr.CreatedAt = now // 新增记录时设置创建时间
			}
		}

		// 3.2 批量插入或更新关联属性
		// 使用ON DUPLICATE KEY UPDATE确保所有需要更新的字段被处理
		if err = tx.Debug().Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}}, // 按主键ID冲突判断
			DoUpdates: clause.Assignments(map[string]interface{}{
				"product_id":           clause.Expr{SQL: "VALUES(product_id)"},
				"product_attribute_id": clause.Expr{SQL: "VALUES(product_attribute_id)"},
				"value":                clause.Expr{SQL: "VALUES(value)"},
				"updated_at":           clause.Expr{SQL: "VALUES(updated_at)"},
				// 如有其他需要更新的字段，在此处添加
			}),
		}).Create(&attrs).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 3.3 清理数据库中存在但不在当前列表中的属性值
		// 收集当前列表中的所有属性ID
		var currentAttrIDs []int
		for _, attr := range attrs {
			if attr.ID > 0 {
				currentAttrIDs = append(currentAttrIDs, attr.ID)
			}
		}

		// 删除不在当前列表中的属性值（保持数据一致性）
		if len(currentAttrIDs) > 0 {
			if err = tx.Where("product_id = ? AND id NOT IN (?)", productId, currentAttrIDs).
				Delete(&wechat.ProductAttributeValue{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 如果当前列表没有带ID的属性，删除该产品所有属性
			if err = tx.Where("product_id = ?", productId).
				Delete(&wechat.ProductAttributeValue{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

// UpdateProductLadderWithProduct 更新Product同时更新ProductLadder
func (exa *HomeService) UpdateProductLadderWithProduct(productId int, ladders []*wechat.ProductLadder, tx *gorm.DB) (err error) {
	// 3. 处理关联的ProductLadder列表（使用ON DUPLICATE KEY UPDATE优化）
	if len(ladders) > 0 {
		// 3.1 确保外键和时间戳正确
		now := time.Now()
		for _, attr := range ladders {
			attr.ProductId = productId // 强制关联当前产品
			attr.UpdatedAt = now       // 更新时间戳
			if attr.CreatedAt.IsZero() {
				attr.CreatedAt = now // 新增记录时设置创建时间
			}
		}

		// 3.2 批量插入或更新关联属性
		// 使用ON DUPLICATE KEY UPDATE确保所有需要更新的字段被处理
		if err = tx.Debug().Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}}, // 按主键ID冲突判断
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": clause.Expr{SQL: "VALUES(updated_at)"},
				"product_id": clause.Expr{SQL: "VALUES(product_id)"},
				"count":      clause.Expr{SQL: "VALUES(count)"},
				"discount":   clause.Expr{SQL: "VALUES(discount)"},
				"price":      clause.Expr{SQL: "VALUES(price)"},
				// 如有其他需要更新的字段，在此处添加
			}),
		}).Create(&ladders).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 3.3 清理数据库中存在但不在当前列表中的属性值
		// 收集当前列表中的所有属性ID
		var laddersIDs []int
		for _, attr := range ladders {
			if attr.ID > 0 {
				laddersIDs = append(laddersIDs, attr.ID)
			}
		}

		// 删除不在当前列表中的属性值（保持数据一致性）
		if len(laddersIDs) > 0 {
			if err = tx.Where("product_id = ? AND id NOT IN (?)", productId, laddersIDs).
				Delete(&wechat.ProductLadder{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 如果当前列表没有带ID的属性，删除该产品所有属性
			if err = tx.Where("product_id = ?", productId).
				Delete(&wechat.ProductLadder{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

// UpdateProductLadderWithProduct 更新Product同时更新满减
func (exa *HomeService) UpdateProductFullReductionWithProduct(productId int, fullReductions []*wechat.ProductFullReduction, tx *gorm.DB) (err error) {
	// 3. 处理关联的FullReduction列表（使用ON DUPLICATE KEY UPDATE优化）
	if len(fullReductions) > 0 {
		// 3.1 确保外键和时间戳正确
		now := time.Now()
		for _, attr := range fullReductions {
			attr.ProductId = productId // 强制关联当前产品
			attr.UpdatedAt = now       // 更新时间戳
			if attr.CreatedAt.IsZero() {
				attr.CreatedAt = now // 新增记录时设置创建时间
			}
		}

		// 3.2 批量插入或更新关联属性
		// 使用ON DUPLICATE KEY UPDATE确保所有需要更新的字段被处理
		if err = tx.Debug().Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}}, // 按主键ID冲突判断
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at":   clause.Expr{SQL: "VALUES(updated_at)"},
				"product_id":   clause.Expr{SQL: "VALUES(product_id)"},
				"full_price":   clause.Expr{SQL: "VALUES(full_price)"},
				"reduce_price": clause.Expr{SQL: "VALUES(reduce_price)"},
				// 如有其他需要更新的字段，在此处添加
			}),
		}).Create(&fullReductions).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 3.3 清理数据库中存在但不在当前列表中的属性值
		// 收集当前列表中的所有属性ID
		var fullsIDs []int
		for _, attr := range fullReductions {
			if attr.ID > 0 {
				fullsIDs = append(fullsIDs, attr.ID)
			}
		}

		// 删除不在当前列表中的属性值（保持数据一致性）
		if len(fullsIDs) > 0 {
			if err = tx.Where("product_id = ? AND id NOT IN (?)", productId, fullsIDs).
				Delete(&wechat.ProductFullReduction{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 如果当前列表没有带ID的属性，删除该产品所有属性
			if err = tx.Where("product_id = ?", productId).
				Delete(&wechat.ProductFullReduction{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

func (exa *HomeService) UpdateProductForKeyword(e *wechatRequest.UpdateIdsKeywordRequest) (err error) {
	db := global.GVA_DB.Model(&wechat.Product{})
	db.Where("id IN ?", e.Ids).Updates(map[string]interface{}{e.Key: e.Value})
	return err
}

// 获取查询数据库的命令
func (exa *HomeService) getProductSearchCmd(searchInfo wechatRequest.ProductSearchInfo) string {
	var cmdName string
	var cmdSN string
	var cmdCategoryName string
	var cmdBrand string
	var cmdPublishStatus string
	var cmdVerifyStatus string
	if len(searchInfo.Keyword) > 0 {
		cmdName += fmt.Sprintf("name like '%%%s%%'", searchInfo.Keyword)
	}
	if searchInfo.BrandId > 0 {
		cmdSN += fmt.Sprintf("brand_id = %d", searchInfo.BrandId)
	}
	if len(searchInfo.ProductSN) > 0 {
		cmdSN += fmt.Sprintf("product_sn like '%%%s%%'", strings.TrimSpace(searchInfo.ProductSN))
	}
	if len(searchInfo.ProductCategoryId) > 0 {
		cmdCategoryName += fmt.Sprintf("product_category_id like '%%%s%%'", strings.TrimSpace(searchInfo.ProductCategoryId))
	}
	if len(searchInfo.BrandName) > 0 {
		cmdBrand += fmt.Sprintf("brand_name like '%%%s%%'", strings.TrimSpace(searchInfo.BrandName))
	}
	if searchInfo.PublishStatus > 0 {
		cmdPublishStatus += fmt.Sprintf("publish_status = %d", searchInfo.PublishStatus-100)
	}
	if searchInfo.VerifyStatus > 0 {
		cmdVerifyStatus += fmt.Sprintf("verify_status = %d", searchInfo.VerifyStatus-100)
	}

	cmdSearch := ""
	cmds := [6]string{cmdName, cmdSN, cmdCategoryName, cmdBrand, cmdPublishStatus, cmdVerifyStatus}
	isFirst := true
	for _, cmd := range cmds {
		if len(cmd) > 0 {
			if isFirst {
				cmdSearch += cmd
				isFirst = false
			} else {
				cmdSearch += " and " + cmd
			}
		}
	}
	return cmdSearch
}

func (exa *HomeService) GetProductByID(id int) (product wechat.Product, err error) {
	err = global.GVA_DB.Where("id = ?", id).
		Preload("Brand").
		Preload("ProductLadderList").
		Preload("ProductFullReductionList").
		Preload("SkuStockList").
		Preload("ProductAttributeValueList").
		Preload("ProductAttributeList").First(&product).Error
	return product, err
}

// GetBuysAvatarsList 根据买家id，获取买家头像
func (exa *HomeService) GetBuysAvatarsList(buyersId []int) (avatars wechatRequest.BuyersInfo, err error) {
	db := global.GVA_DB.Model(&wechat.WXUser{})
	err = db.Select("avatarUrl").Where("id In ?", buyersId).Find(&avatars).Error
	return avatars, err
}
func (exa *HomeService) GetProductList(searchInfo wechatRequest.ProductSearchInfo) (list []wechat.Product, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	cmd := exa.getProductSearchCmd(searchInfo)

	db := global.GVA_DB.Model(&wechat.Product{})
	err = db.Where(cmd).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Debug().Limit(limit).Offset(offset).Debug().Where(cmd).Find(&list).Error
	}

	return list, total, err
}

func (exa *HomeService) GetBySimpleProductList(searchInfo request.KeyWordInfo) (list []wechat.Product, err error) {
	db := global.GVA_DB.Model(&wechat.Product{})
	err = db.Debug().Debug().Where("name like ? or product_sn like ?", "%"+searchInfo.Keyword+"%", "%"+searchInfo.Keyword+"%").Find(&list).Error

	return list, err
}

func (exa *HomeService) GetProductListByOnlyID(searchInfo *request.KeySearchInfo) (productList []wechat.Product, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&wechat.Product{})
	cmd := fmt.Sprintf("%s = %d", searchInfo.Key, searchInfo.ID)
	err = db.Where(cmd).Count(&total).Error
	if err != nil {
		return productList, total, err
	} else {
		err = db.Debug().Limit(limit).Offset(offset).Debug().Where(cmd).Find(&productList).Error
	}
	return productList, total, err
}

func (exa *HomeService) DeleteProducts(ids []int) (err error) {
	var products []wechat.Product
	db := global.GVA_DB
	//if err = db.Where("id in ?", ids).Preload("ProductLadderList").Preload("ProductFullReductionList").Preload("SkuStockList").Find(&products).Error; err != nil {
	//	return err
	//}

	db.Where("id in ?", ids).Select("ProductLadderList", "ProductFullReductionList", "SkuStockList").Delete(&products)
	return err
}

func (exa *HomeService) GetProductListByOnlyIDWithSort(searchInfo *request.SortSearchInfo) (productList []wechat.Product, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&wechat.Product{})

	cmd := fmt.Sprintf("%s = %d", searchInfo.Key, searchInfo.ID)
	var orderCmd string
	switch searchInfo.Sort {
	case 0:
		orderCmd = fmt.Sprintf("`sale` desc, `price` asc")
	case 2:
		orderCmd = fmt.Sprintf("`sale` desc")
	case 3:
		orderCmd = fmt.Sprintf("`price` asc")
	case 4:
		orderCmd = fmt.Sprintf("`price` desc")
	default:
		global.GVA_LOG.Error(fmt.Sprintf("sort参数错误,sort:%d", searchInfo.Sort), zap.Error(err))
		return productList, total, err
	}
	err = db.Where(cmd).Count(&total).Error
	if err != nil {
		return productList, total, err
	} else {
		err = db.Debug().Limit(limit).Offset(offset).Where(cmd).Order(orderCmd).Find(&productList).Error
	}
	return productList, total, err
}

func (exa *HomeService) GetProductBrand(id int) (brand wechat.Brand, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&brand).Error
	return
}

// 获取查询数据库的命令
func (exa *HomeService) getBrandSearchCmd(searchInfo wechatRequest.BrandSearchInfo) string {
	var cmdName string
	var cmdStatus string

	if len(searchInfo.Name) > 0 {
		cmdName += fmt.Sprintf("name like '%%%s%%'", searchInfo.Name)
	}
	if searchInfo.ShowStatus > 0 {
		cmdStatus += fmt.Sprintf("show_status = %d", searchInfo.ShowStatus-100)
	}

	cmdSearch := ""
	cmds := [2]string{cmdName, cmdStatus}
	isFirst := true
	for _, cmd := range cmds {
		if len(cmd) > 0 {
			if isFirst {
				cmdSearch += cmd
				isFirst = false
			} else {
				cmdSearch += " and " + cmd
			}
		}
	}
	return cmdSearch
}

func (exa *HomeService) GetProductBrandList(searchInfo wechatRequest.BrandSearchInfo) (list []wechat.Brand, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	cmd := exa.getBrandSearchCmd(searchInfo)

	db := global.GVA_DB.Model(&wechat.Brand{})
	err = db.Where(cmd).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Debug().Limit(limit).Offset(offset).Debug().Where(cmd).Find(&list).Error
	}

	return list, total, err
}

func (exa *HomeService) CreateHomeProductBrand(e *wechat.Brand) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) UpdateHomeBrand(e *wechat.Brand) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *HomeService) UpdateHomeBrandByIdForKeyword(e *wechatRequest.UpdateIdsKeywordRequest) (err error) {
	db := global.GVA_DB.Model(&wechat.Brand{})
	db.Where("id IN ?", e.Ids).Updates(map[string]interface{}{e.Key: e.Value})
	return err
}

func (exa *HomeService) DeleteHomeProductBrand(id int) (err error) {
	var brand wechat.Brand
	err = global.GVA_DB.Where("id = ?", id).Delete(&brand).Error
	return err
}

func (exa *HomeService) CreateProductAttributeCategory(e *wechat.ProductAttributeCategory) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) UpdateProductAttributeCategory(e *wechat.ProductAttributeCategory) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *HomeService) UpdateProductAttributeCategoryCount(id int, keyword string) (err error) {
	db := global.GVA_DB.Model(&wechat.ProductAttributeCategory{})
	//cmd := fmt.Sprintf("%s + %d", keyword, 1)
	db.Where("id = ?", id).Update(keyword, gorm.Expr(keyword+" + ?", 1))
	return err
}

func (exa *HomeService) DeleteProductAttributeCategory(id int) (err error) {
	var attribute wechat.ProductAttributeCategory
	err = global.GVA_DB.Where("id = ?", id).Delete(&attribute).Error
	return err
}

func (exa *HomeService) GetProductAttributeCategoryList(searchInfo request.PageInfo) (list interface{}, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	var productList []wechat.ProductAttributeCategory

	db := global.GVA_DB.Model(&wechat.ProductAttributeCategory{})
	err = db.Count(&total).Error
	if err != nil {
		return productList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Find(&productList).Error
	}

	return productList, total, err
}

func (exa *HomeService) GetProductAttributeListById(id int) (list []wechat.ProductAttribute, err error) {
	db := global.GVA_DB.Model(&wechat.ProductAttribute{})
	err = db.Where("product_attribute_category_id = ?", id).Find(&list).Error
	return list, err
}

func (exa *HomeService) CreateProductAttributeSynchronous(e *wechat.ProductAttribute) (err error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return err
	}
	// 修改商品类型属性和参数统计
	if e.Type == 0 {
		err = exa.UpdateProductAttributeCategoryCount(e.ProductAttributeCategoryId, "attribute_count")
	} else if e.Type == 1 {
		err = exa.UpdateProductAttributeCategoryCount(e.ProductAttributeCategoryId, "param_count")
	} else {
		tx.Rollback()
		return errors.New("属性类型错误")
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	// 存储属性
	err = global.GVA_DB.Create(&e).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (exa *HomeService) CreateProductAttribute(e *wechat.ProductAttribute) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) UpdateProductAttribute(e *wechat.ProductAttribute) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *HomeService) DeleteProductAttribute(id int) (err error) {
	var attribute wechat.ProductAttribute
	err = global.GVA_DB.Where("id = ?", id).Delete(&attribute).Error
	return err
}

func (exa *HomeService) GetProductAttributeList(searchInfo request.TagSearchInfo) (list []wechat.ProductAttribute, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&wechat.ProductAttribute{})
	err = db.Where("product_attribute_category_id=? and type=?", searchInfo.Tag, searchInfo.State).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Debug().Where("product_attribute_category_id=? and type=?", searchInfo.Tag, searchInfo.State).Order("sort desc").Find(&list).Error
	}

	return list, total, err
}

func (exa *HomeService) CreateProductCategory(e *wechat.ProductCategory) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) UpdateProductCategory(e *wechat.ProductCategory) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *HomeService) DeleteProductCategory(id int) (err error) {
	var category wechat.ProductCategory
	err = global.GVA_DB.Where("id = ?", id).Delete(&category).Error
	return err
}

func (exa *HomeService) GetProductCategoryList(searchInfo request.TagSearchInfo) (productList []wechat.ProductCategory, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&wechat.ProductCategory{})
	err = db.Where("parent_id = ?", searchInfo.Tag).Count(&total).Error
	if err != nil {
		return productList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Debug().Where("parent_id = ?", searchInfo.Tag).Find(&productList).Error
	}

	return productList, total, err
}

func (exa *HomeService) GetAllProductCategoryList() (productList []*wechat.ProductCategory, err error) {
	db := global.GVA_DB.Model(&wechat.ProductCategory{})
	err = db.Find(&productList).Error
	return productList, err
}

func (exa *HomeService) GetProductAttributeValueByProductId(id int) (list []wechat.ProductAttributeValue, err error) {
	db := global.GVA_DB.Model(&wechat.ProductAttributeValue{})
	err = db.Where("product_id = ?", id).Find(&list).Error
	return list, err
}

func (exa *HomeService) CreateProductMemberPrice(e *wechat.MemberPrice) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) CreateProductAttributeValue(e *wechat.ProductAttributeValue) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) GetProductFullReductionByProductId(id int) (list []wechat.ProductFullReduction, err error) {
	db := global.GVA_DB.Model(&wechat.ProductFullReduction{})
	err = db.Where("product_id = ?", id).Find(&list).Error
	return list, err
}

func (exa *HomeService) CreateProductFullReduction(e *wechat.ProductFullReduction) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) CreateProductLadder(e *wechat.ProductLadder) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) CreateProductSKUStock(e *wechat.SkuStock) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *HomeService) GetProductSKUStockByProductId(id int, keyword string) (list []wechat.SkuStock, err error) {
	db := global.GVA_DB.Model(&wechat.SkuStock{})
	cmd := fmt.Sprintf("product_id = %d", id)
	if len(keyword) > 0 {
		cmd += fmt.Sprintf("and sku_code = %s", keyword)
	}
	err = db.Where(cmd).Find(&list).Error
	return list, err
}

func (exa *HomeService) GetProductSKUStockById(id int) (stock wechat.SkuStock, err error) {
	db := global.GVA_DB.Model(&wechat.SkuStock{})
	err = db.Where("id = ?", id).First(&stock).Error
	return stock, err
}

func (exa *HomeService) UpdateSKUStock(e *wechat.SkuStock) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *HomeService) DeleteSKUStockWithBatch(productId int) (err error) {
	var skuStock wechat.SkuStock
	err = global.GVA_DB.Where("product_id = ?", productId).Delete(&skuStock).Unscoped().Error
	return err
}

func (exa *HomeService) UpdateProductSkuStockForStock(id int, quantity int) (int64, error) {
	cmd := fmt.Sprintf("UPDATE pms_sku_stock SET lock_stock = lock_stock - %d, stock = stock - %d "+
		"WHERE id = %d AND stock - %d >= 0 AND lock_stock - %d >= 0", quantity, quantity, id, quantity, quantity)
	ret := global.GVA_DB.Debug().Exec(cmd)

	return ret.RowsAffected, ret.Error
}

func (exa *HomeService) GetProductCartList() (list []wechat.CartItem, err error) {
	db := global.GVA_DB.Preload("Product").Preload("SkuStock").Model(&wechat.CartItem{})
	err = db.Find(&list).Error
	return list, err
}

func (exa *HomeService) CreateProductCart(e *wechat.CartItem) (err error) {
	db := global.GVA_DB.Model(&wechat.CartItem{})
	var cart wechat.CartItem
	result := db.Where("user_id = ? and product_id = ? and sku_stock_id = ?", e.UserId, e.ProductId, e.SkuStockId).First(&cart)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			err = global.GVA_DB.Create(&e).Error
			return err
		}
		err = result.Error
	} else {
		err = db.Debug().Where("id = ? and user_id = ?", cart.ID, e.UserId).Update("quantity", cart.Quantity+1).Error
		return err
	}
	return err
}

func (exa *HomeService) UpdateProductCartQuantity(userId int, id int, quantity int) (err error) {
	db := global.GVA_DB.Model(&wechat.CartItem{})
	db.Debug().Where("user_id = ? and id = ?", userId, id).UpdateColumn("quantity", quantity)
	return err
}

func (exa *HomeService) DeleteProductCartById(userId int, id int) (err error) {
	var cart wechat.CartItem
	err = global.GVA_DB.Where("user_id = ? and id = ?", userId, id).Delete(&cart).Error
	return err
}

func (exa *HomeService) DeleteProductCartByIds(userId int, ids []int) (err error) {
	var cart wechat.CartItem
	err = global.GVA_DB.Where("user_id = ? and id in ?", userId, ids).Delete(&cart).Error
	return err
}

func (exa *HomeService) ClearProductCartUserId(id int) (err error) {
	var cart wechat.CartItem
	err = global.GVA_DB.Where("user_id = ?", id).Delete(&cart).Error
	return err
}

func (exa *HomeService) CreateProductTmpCart(e *wechat.CartTmpItem) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
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
			err = db.Limit(limit).Offset(offset).Where(cmdName).Find(&list).Error
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
	db := global.GVA_DB.Model(&wechat.FlashPromotionProductRelation{})
	err = db.Debug().Where("flash_promotion_id = ? and flash_promotion_session_id = ?", flashPromotionId, flashPromotionSessionId).Preload("Product").Find(&list).Error
	return list, err
}

func (exa *HomeService) CreateFlashPromotionProductRelation(e []wechat.FlashPromotionProductRelation) (err error) {
	err = global.GVA_DB.CreateInBatches(e, len(e)).Error
	return err
}

func (exa *HomeService) UpdateProductPromotionType(relation *wechat.FlashPromotionProductRelation, session *wechat.FlashPromotionSession) (err error) {
	db := global.GVA_DB.Model(&wechat.Product{})
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

func (exa *HomeService) GetFlashSessionList() (list []wechat.FlashPromotionSession, err error) {
	db := global.GVA_DB.Model(&wechat.FlashPromotionSession{})
	err = db.Find(&list).Error
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
	db := global.GVA_DB.Model(&wechat.FlashPromotionSession{})
	err = db.Preload("ProductRelation").Find(&list).Error
	return list, err
}

func (exa *HomeService) GetGroupBuyProductList() (list []*wechat.GroupBuyProduct, err error) {
	db := global.GVA_DB.Model(&wechat.GroupBuyProduct{})
	err = db.Preload("Product").Find(&list).Error
	return list, err
}
