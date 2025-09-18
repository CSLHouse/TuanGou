package product

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/product"
	"cooller/server/model/wechat"
	wechatRequest "cooller/server/model/wechat/request"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type ProductService struct{}

func (exa *ProductService) CreateHomeProduct(e *product.Product) (err error) {
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

func (exa *ProductService) UpdateHomeProductSynchronous(e *product.Product) (err error) {
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
func (exa *ProductService) UpdateSKUStockWithProduct(productId int, skus []*product.SkuStock, tx *gorm.DB) (err error) {
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
				Delete(&product.SkuStock{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 如果当前列表没有带ID的属性，删除该产品所有属性
			if err = tx.Where("product_id = ?", productId).
				Delete(&product.SkuStock{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

// UpdateProductAttributeValueWithProduct 更新Product同时更新属性值
func (exa *ProductService) UpdateProductAttributeValueWithProduct(productId int, attrs []*product.ProductAttributeValue, tx *gorm.DB) (err error) {
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
				Delete(&product.ProductAttributeValue{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 如果当前列表没有带ID的属性，删除该产品所有属性
			if err = tx.Where("product_id = ?", productId).
				Delete(&product.ProductAttributeValue{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

// UpdateProductLadderWithProduct 更新Product同时更新ProductLadder
func (exa *ProductService) UpdateProductLadderWithProduct(productId int, ladders []*product.ProductLadder, tx *gorm.DB) (err error) {
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
				Delete(&product.ProductLadder{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 如果当前列表没有带ID的属性，删除该产品所有属性
			if err = tx.Where("product_id = ?", productId).
				Delete(&product.ProductLadder{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

// UpdateProductLadderWithProduct 更新Product同时更新满减
func (exa *ProductService) UpdateProductFullReductionWithProduct(productId int, fullReductions []*product.ProductFullReduction, tx *gorm.DB) (err error) {
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
				Delete(&product.ProductFullReduction{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 如果当前列表没有带ID的属性，删除该产品所有属性
			if err = tx.Where("product_id = ?", productId).
				Delete(&product.ProductFullReduction{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

func (exa *ProductService) UpdateProductForKeyword(e *wechatRequest.UpdateIdsKeywordRequest) (err error) {
	db := global.GVA_DB.Model(&product.Product{})
	db.Where("id IN ?", e.Ids).Updates(map[string]interface{}{e.Key: e.Value})
	return err
}

// 获取查询数据库的命令
func (exa *ProductService) getProductSearchCmd(searchInfo wechatRequest.ProductSearchInfo) string {
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

func (exa *ProductService) GetProductByID(id int) (product product.Product, err error) {
	err = global.GVA_DB.Where("id = ?", id).
		Preload("Brand").
		Preload("ProductLadderList").
		Preload("ProductFullReductionList").
		Preload("SkuStockList").
		Preload("ProductAttributeValueList").
		Preload("ProductAttributeList").First(&product).Error
	return product, err
}

func (exa *ProductService) GetProductCouponByID(id int) (productData product.Product, err error) {
	err = global.GVA_DB.Debug().Where("id = ?", id).Preload("CouponList").First(&productData).Error
	return productData, err
}

// GetBuysAvatarsList 根据买家id，获取买家头像
func (exa *ProductService) GetBuysAvatarsList(buyersId []int) (avatars wechatRequest.BuyersInfo, err error) {
	db := global.GVA_DB.Model(&wechat.WXUser{})
	err = db.Select("avatarUrl").Where("id In ?", buyersId).Find(&avatars).Error
	return avatars, err
}
func (exa *ProductService) GetProductList(searchInfo wechatRequest.ProductSearchInfo) (list []product.Product, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	cmd := exa.getProductSearchCmd(searchInfo)

	db := global.GVA_DB.Model(&product.Product{})
	err = db.Where(cmd).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Debug().Limit(limit).Offset(offset).Debug().Where(cmd).Find(&list).Error
	}

	return list, total, err
}

func (exa *ProductService) GetBySimpleProductList(searchInfo request.KeyWordInfo) (list []product.Product, err error) {
	db := global.GVA_DB.Model(&product.Product{})
	err = db.Debug().Debug().Where("name like ? or product_sn like ?", "%"+searchInfo.Keyword+"%", "%"+searchInfo.Keyword+"%").Find(&list).Error

	return list, err
}

func (exa *ProductService) GetProductListByOnlyID(searchInfo *request.KeySearchInfo) (productList []product.Product, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&product.Product{})
	cmd := fmt.Sprintf("%s = %d", searchInfo.Key, searchInfo.ID)
	err = db.Where(cmd).Count(&total).Error
	if err != nil {
		return productList, total, err
	} else {
		err = db.Debug().Limit(limit).Offset(offset).Debug().Where(cmd).Find(&productList).Error
	}
	return productList, total, err
}

func (exa *ProductService) DeleteProducts(ids []int) (err error) {
	var products []product.Product
	db := global.GVA_DB
	//if err = db.Where("id in ?", ids).Preload("ProductLadderList").Preload("ProductFullReductionList").Preload("SkuStockList").Find(&products).Error; err != nil {
	//	return err
	//}

	db.Where("id in ?", ids).Select("ProductLadderList", "ProductFullReductionList", "SkuStockList").Delete(&products)
	return err
}

func (exa *ProductService) GetProductListByOnlyIDWithSort(searchInfo *request.SortSearchInfo) (productList []product.Product, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&product.Product{})

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

func (exa *ProductService) GetProductBrand(id int) (brand product.Brand, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&brand).Error
	return
}

// 获取查询数据库的命令
func (exa *ProductService) getBrandSearchCmd(searchInfo wechatRequest.BrandSearchInfo) string {
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

func (exa *ProductService) GetProductBrandList(searchInfo wechatRequest.BrandSearchInfo) (list []product.Brand, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	cmd := exa.getBrandSearchCmd(searchInfo)

	db := global.GVA_DB.Model(&product.Brand{})
	err = db.Where(cmd).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Debug().Limit(limit).Offset(offset).Debug().Where(cmd).Find(&list).Error
	}

	return list, total, err
}

func (exa *ProductService) CreateHomeProductBrand(e *product.Brand) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *ProductService) UpdateHomeBrand(e *product.Brand) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *ProductService) UpdateHomeBrandByIdForKeyword(e *wechatRequest.UpdateIdsKeywordRequest) (err error) {
	db := global.GVA_DB.Model(&product.Brand{})
	db.Where("id IN ?", e.Ids).Updates(map[string]interface{}{e.Key: e.Value})
	return err
}

func (exa *ProductService) DeleteHomeProductBrand(id int) (err error) {
	var brand product.Brand
	err = global.GVA_DB.Where("id = ?", id).Delete(&brand).Error
	return err
}

func (exa *ProductService) CreateProductAttributeCategory(e *product.ProductAttributeCategory) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *ProductService) UpdateProductAttributeCategory(e *product.ProductAttributeCategory) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *ProductService) UpdateProductAttributeCategoryCount(id int, keyword string) (err error) {
	db := global.GVA_DB.Model(&product.ProductAttributeCategory{})
	//cmd := fmt.Sprintf("%s + %d", keyword, 1)
	db.Where("id = ?", id).Update(keyword, gorm.Expr(keyword+" + ?", 1))
	return err
}

func (exa *ProductService) DeleteProductAttributeCategory(id int) (err error) {
	var attribute product.ProductAttributeCategory
	err = global.GVA_DB.Where("id = ?", id).Delete(&attribute).Error
	return err
}

func (exa *ProductService) GetProductAttributeCategoryList(searchInfo request.PageInfo) (list interface{}, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	var productList []product.ProductAttributeCategory

	db := global.GVA_DB.Model(&product.ProductAttributeCategory{})
	err = db.Count(&total).Error
	if err != nil {
		return productList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Find(&productList).Error
	}

	return productList, total, err
}

func (exa *ProductService) GetProductAttributeListById(id int) (list []product.ProductAttribute, err error) {
	db := global.GVA_DB.Model(&product.ProductAttribute{})
	err = db.Where("product_attribute_category_id = ?", id).Find(&list).Error
	return list, err
}

func (exa *ProductService) CreateProductAttributeSynchronous(e *product.ProductAttribute) (err error) {
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

func (exa *ProductService) CreateProductAttribute(e *product.ProductAttribute) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *ProductService) UpdateProductAttribute(e *product.ProductAttribute) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *ProductService) DeleteProductAttribute(id int) (err error) {
	var attribute product.ProductAttribute
	err = global.GVA_DB.Where("id = ?", id).Delete(&attribute).Error
	return err
}

func (exa *ProductService) GetProductAttributeList(searchInfo request.TagSearchInfo) (list []product.ProductAttribute, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&product.ProductAttribute{})
	err = db.Where("product_attribute_category_id=? and type=?", searchInfo.Tag, searchInfo.State).Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Debug().Where("product_attribute_category_id=? and type=?", searchInfo.Tag, searchInfo.State).Order("sort desc").Find(&list).Error
	}

	return list, total, err
}

func (exa *ProductService) CreateProductCategory(e *product.ProductCategory) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *ProductService) UpdateProductCategory(e *product.ProductCategory) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *ProductService) DeleteProductCategory(id int) (err error) {
	var category product.ProductCategory
	err = global.GVA_DB.Where("id = ?", id).Delete(&category).Error
	return err
}

func (exa *ProductService) GetProductCategoryList(searchInfo request.TagSearchInfo) (productList []product.ProductCategory, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&product.ProductCategory{})
	err = db.Where("parent_id = ?", searchInfo.Tag).Count(&total).Error
	if err != nil {
		return productList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Debug().Where("parent_id = ?", searchInfo.Tag).Find(&productList).Error
	}

	return productList, total, err
}

func (exa *ProductService) GetAllProductCategoryList() (productList []*product.ProductCategory, err error) {
	db := global.GVA_DB.Model(&product.ProductCategory{})
	err = db.Find(&productList).Error
	return productList, err
}

func (exa *ProductService) GetProductCategoryById(id int) (productCategory product.ProductCategory, err error) {
	db := global.GVA_DB.Model(&product.ProductCategory{})
	err = db.Where("id = ?", id).Preload("CouponList").First(&productCategory).Error
	return productCategory, err
}

func (exa *ProductService) GetProductAttributeValueByProductId(id int) (list []product.ProductAttributeValue, err error) {
	db := global.GVA_DB.Model(&product.ProductAttributeValue{})
	err = db.Where("product_id = ?", id).Find(&list).Error
	return list, err
}

func (exa *ProductService) CreateProductAttributeValue(e *product.ProductAttributeValue) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *ProductService) GetProductFullReductionByProductId(id int) (list []product.ProductFullReduction, err error) {
	db := global.GVA_DB.Model(&product.ProductFullReduction{})
	err = db.Where("product_id = ?", id).Find(&list).Error
	return list, err
}

func (exa *ProductService) CreateProductFullReduction(e *product.ProductFullReduction) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *ProductService) CreateProductLadder(e *product.ProductLadder) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *ProductService) CreateProductSKUStock(e *product.SkuStock) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (exa *ProductService) GetProductSKUStockByProductId(id int, keyword string) (list []product.SkuStock, err error) {
	db := global.GVA_DB.Model(&product.SkuStock{})
	cmd := fmt.Sprintf("product_id = %d", id)
	if len(keyword) > 0 {
		cmd += fmt.Sprintf("and sku_code = %s", keyword)
	}
	err = db.Where(cmd).Find(&list).Error
	return list, err
}

func (exa *ProductService) GetProductSKUStockById(id int) (stock product.SkuStock, err error) {
	db := global.GVA_DB.Model(&product.SkuStock{})
	err = db.Where("id = ?", id).First(&stock).Error
	return stock, err
}

func (exa *ProductService) UpdateSKUStock(e *product.SkuStock) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (exa *ProductService) DeleteSKUStockWithBatch(productId int) (err error) {
	var skuStock product.SkuStock
	err = global.GVA_DB.Where("product_id = ?", productId).Delete(&skuStock).Unscoped().Error
	return err
}

func (exa *ProductService) UpdateProductSkuStockForStock(id int, quantity int) (int64, error) {
	cmd := fmt.Sprintf("UPDATE pms_sku_stock SET lock_stock = lock_stock - %d, stock = stock - %d "+
		"WHERE id = %d AND stock - %d >= 0 AND lock_stock - %d >= 0", quantity, quantity, id, quantity, quantity)
	ret := global.GVA_DB.Debug().Exec(cmd)

	return ret.RowsAffected, ret.Error
}

func (exa *ProductService) GetProductCartList() (list []product.CartItem, err error) {
	db := global.GVA_DB.Preload("Product").Preload("SkuStock").Model(&product.CartItem{})
	err = db.Find(&list).Error
	return list, err
}

func (exa *ProductService) CreateProductCart(e *product.CartItem) (err error) {
	db := global.GVA_DB.Model(&product.CartItem{})
	var cart product.CartItem
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

func (exa *ProductService) UpdateProductCartQuantity(userId int, id int, quantity int) (err error) {
	db := global.GVA_DB.Model(&product.CartItem{})
	db.Debug().Where("user_id = ? and id = ?", userId, id).UpdateColumn("quantity", quantity)
	return err
}

func (exa *ProductService) DeleteProductCartById(userId int, id int) (err error) {
	var cart product.CartItem
	err = global.GVA_DB.Where("user_id = ? and id = ?", userId, id).Delete(&cart).Error
	return err
}

func (exa *ProductService) DeleteProductCartByIds(userId int, ids []int) (err error) {
	var cart product.CartItem
	err = global.GVA_DB.Where("user_id = ? and id in ?", userId, ids).Delete(&cart).Error
	return err
}

func (exa *ProductService) ClearProductCartUserId(id int) (err error) {
	var cart product.CartItem
	err = global.GVA_DB.Where("user_id = ?", id).Delete(&cart).Error
	return err
}

func (exa *ProductService) CreateProductTmpCart(e *product.CartTmpItem) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}
