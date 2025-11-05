package product

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/product"
	productReq "cooller/server/model/product/request"
	date_conversion "cooller/server/utils/timer"
	"cooller/server/utils/upload"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderService struct{}

func (o *OrderService) GetProductCartById(userId int, id int) (cartItem product.CartItem, err error) {
	db := global.GVA_DB.Model(&product.CartItem{})
	db.Where("user_id = ? and id = ?", userId, id).First(&cartItem)
	return cartItem, err
}

func (o *OrderService) GetProductCartByIds(userId int, ids []int) (cartItem []*product.CartCommonItem, err error) {
	db := global.GVA_DB.Model(&product.CartItem{})
	var list []product.CartItem
	err = db.Debug().Preload("Product").Preload("SkuStock").Where("user_id = ? and id in ?", userId, ids).First(&list).Error
	result := make([]*product.CartCommonItem, len(list))
	for i, item := range list {
		result[i] = &product.CartCommonItem{
			GVA_MODEL: global.GVA_MODEL{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			},
			ProductId:  item.ProductId,
			SkuStockId: item.SkuStockId,
			UserId:     item.UserId,
			Quantity:   item.Quantity,
			Product:    item.Product,
			SkuStock:   item.SkuStock,
			Price:      item.Price,
		}
	}
	return result, err
}

func (o *OrderService) GetProductTmpCartByIds(userId int, ids []int) (cartItem []*product.CartCommonItem, err error) {
	db := global.GVA_DB.Model(&product.CartTmpItem{})
	var list []product.CartTmpItem
	err = db.Debug().Preload("Product").Preload("SkuStock").Where("user_id = ? and id in ?", userId, ids).First(&list).Error
	result := make([]*product.CartCommonItem, len(list))
	for i, item := range list {
		result[i] = &product.CartCommonItem{
			GVA_MODEL: global.GVA_MODEL{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			},
			ProductId:  item.ProductId,
			SkuStockId: item.SkuStockId,
			UserId:     item.UserId,
			Quantity:   item.Quantity,
			Product:    item.Product,
			SkuStock:   item.SkuStock,
			Price:      item.Price,
		}
	}
	return result, err
}

func (o *OrderService) CreateOrder(e *product.Order) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Debug().Create(&e).Error
	return err
}

func (o *OrderService) CreateOrderItem(e product.OrderItem) (err error) {
	db := global.GVA_DB.Model(&product.OrderItem{})
	err = db.Create(&e).Error

	return err
}

func (o *OrderService) CreateOrderItemByBatch(e []*product.OrderItem) (err error) {
	db := global.GVA_DB.Model(&product.OrderItem{})
	err = db.CreateInBatches(&e, len(e)).Error
	return err
}

func (o *OrderService) GetProductOrderById(id int) (order product.Order, err error) {
	db := global.GVA_DB.Model(&product.Order{})
	db.Debug().Where("id = ?", id).Preload("OrderItemList").First(&order)
	return order, err
}

func (o *OrderService) GetProductOrderItemById(id int) (orderItem product.OrderItem, err error) {
	db := global.GVA_DB.Model(&product.OrderItem{})
	db.Debug().Where("id = ?", id).First(&orderItem)
	return orderItem, err
}

func (o *OrderService) GetProductOrderList(info request.PageInfo) (list []product.Order, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Preload("OrderItemList").Find(&list).Error
	}
	return list, total, err
}

func (o *OrderService) GetProductOrderListByStatus(searchInfo productReq.SearchInfo) (list []product.Order, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)
	state := searchInfo.State
	db := global.GVA_DB.Model(&product.Order{})

	if searchInfo.OrderId > 0 {
		db = db.Where("id = ?", searchInfo.OrderId)
	}
	if len(searchInfo.OrderSn) > 0 {
		db = db.Where("order_sn = ?", strings.TrimSpace(searchInfo.OrderSn))
	}
	if len(searchInfo.ReceiverKeyword) > 0 {
		// 模糊查询需手动拼接%，Gorm会自动处理参数绑定
		db = db.Where("receiver_name like ? or receiver_phone like ?", "%"+searchInfo.ReceiverKeyword+"%")
	}
	if searchInfo.OrderType > 0 {
		db = db.Where("order_type = ?", searchInfo.OrderType-100)
	}
	if len(searchInfo.CreateTime) > 0 {
		thatDay := date_conversion.ParseStringDate(searchInfo.CreateTime)
		nextDay := thatDay.AddDate(0, 0, 1)

		db = db.Where("created_at between ? and ?", thatDay, nextDay)
	}
	if state >= 0 {
		db = db.Where("status = ?", state)
	}

	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Debug().Preload("OrderItemList").Order("id desc").Find(&list).Error
	}
	return list, total, err
}

func (o *OrderService) GetProductOrderItemListByStatus(searchInfo productReq.SearchInfo) (list []product.OrderItem, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)
	state := searchInfo.State
	db := global.GVA_DB.Model(&product.OrderItem{})

	if searchInfo.OrderId > 0 {
		db = db.Where("order_id = ?", searchInfo.OrderId)
	}
	if len(searchInfo.OrderSn) > 0 {
		db = db.Where("order_sn = ?", strings.TrimSpace(searchInfo.OrderSn))
	}

	if state >= 0 {
		db = db.Where("status = ?", state)
	}

	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Debug().Order("id desc").Find(&list).Error
	}
	return list, total, err
}

func (o *OrderService) UpdateOrderStatus(e *productReq.PaySuccessRequest, status int) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Debug().Where("id = ?", e.OrderId).Updates(map[string]interface{}{"pay_type": e.PayType, "status": status}).Error
	return err
}

func (o *OrderService) UpdateOrderStatusById(orderId int, status int) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Debug().Where("id = ?", orderId).UpdateColumn("status", status).Error
	return err
}

func (o *OrderService) UpdateOrderAfterSalesStatusById(orderId int, status int) (err error) {
	db := global.GVA_DB.Model(&product.OrderItem{})
	err = db.Debug().Where("id = ?", orderId).UpdateColumn("status", status).Error
	return err
}

func (o *OrderService) UpdateOrderItemStatusById(orderId int, status int) (err error) {
	db := global.GVA_DB.Model(&product.OrderItem{})
	err = db.Debug().Where("id = ?", orderId).UpdateColumn("status", status).Error
	return err
}

func (o *OrderService) UpdateOrderStatusByOrderSn(orderSn *string, status int) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Debug().Where("order_sn = ?", orderSn).UpdateColumn("status", status).Error
	return err
}

func (o *OrderService) UpdateOrderPaySuccess(id int, prepayId string) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Debug().Where("id = ?", id).Updates(map[string]interface{}{"prepay_id": prepayId}).Error
	return err
}

func (o *OrderService) CancelOrder(id int) (outTrade string, err error) {
	var order product.Order
	db := global.GVA_DB.Model(&product.Order{})
	if err := db.Preload("OrderItemList").First(&order, id).Error; err != nil {
		return "", err
	}
	// 执行关联删除操作
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&order).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return order.OrderSn, err
	}
	return order.OrderSn, nil
}

func (o *OrderService) DeleteManyOrder(ids []int) (orderList []product.Order, err error) {
	db := global.GVA_DB.Model(&product.Order{})
	if err := db.Where("id in ?", ids).Preload("OrderItemList").Find(&orderList).Error; err != nil {
		return orderList, err
	}
	// 执行关联删除操作
	err = db.Select("OrderItemList").Delete(&orderList).Error
	if err != nil {
		return orderList, err
	}
	return orderList, nil
}

func (o *OrderService) UpdateManyOrderStatus(ids []int, status int) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Where("id in ?", ids).UpdateColumn("status", status).Error
	return err
}

func (o *OrderService) UpdateOrderReceiverInfo(e *productReq.OrderReceiveAddress) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Select("receiver_name", "receiver_phone", "receiver_post_code", "receiver_detail_address", "receiver_province", "receiver_city", "receiver_region").
		Where("id=?", e.OrderId).
		Updates(map[string]interface{}{
			"receiver_name":           e.ReceiverName,
			"receiver_phone":          e.ReceiverPhone,
			"receiver_post_code":      e.ReceiverPostCode,
			"receiver_detail_address": e.ReceiverDetailAddress,
			"receiver_province":       e.ReceiverProvince,
			"receiver_city":           e.ReceiverCity,
			"receiver_region":         e.ReceiverRegion,
		}).Error
	return err
}

func (o *OrderService) UpdateOrderMoneyInfo(info *productReq.OrderMoneyInfo) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Where("id = ?", info.OrderId).UpdateColumn("discount_amount", info.DiscountAmount).Error
	return err
}

func (o *OrderService) UpdateOrderNoteInfo(info *productReq.OrderNoteInfo) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Where("id = ?", info.OrderId).UpdateColumn("note", info.Note).Error
	return err
}

func (o *OrderService) UpdateOrdersStatus(ids []int, status int) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Where("id in ?", ids).UpdateColumn("status", status).Error
	return err
}

func (o *OrderService) GetOrderSetting(id int) (order product.OrderSetting, err error) {
	db := global.GVA_DB.Model(&product.OrderSetting{})
	db.Where("id = ?", id).First(&order)
	return order, err
}

func (o *OrderService) UpdateOrderSetting(e *product.OrderSetting) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

func (o *OrderService) UpdateOrderLogistics(infos *productReq.UpdateLogisticsRequest, ids []int) (err error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return err
	}

	// 3. 构造批量更新SQL（使用CASE WHEN）
	companyCase := "CASE id " // 物流公司的CASE语句
	snCase := "CASE id "      // 物流单号的CASE语句
	args := []interface{}{}   // SQL参数（防止注入）

	for _, info := range infos.LogisticsInfos {
		// 拼接CASE条件（使用参数化占位符?）
		companyCase += "WHEN ? THEN ? "
		snCase += "WHEN ? THEN ? "
		// 绑定参数（ID、公司名、ID、单号）
		args = append(args, info.ID, info.LogisticsCompany, info.ID, info.LogisticsSn)
	}
	companyCase += "END" // 结束CASE语句
	snCase += "END"

	// 构造WHERE IN条件（筛选需要更新的ID）
	inPlaceholders := "(" + strings.Repeat("?,", len(ids)-1) + "?)" // 如 (?,?)
	for _, id := range ids {
		args = append(args, id) // 这里int会被隐式转换为interface{}
	}

	fmt.Println(append([]interface{}{companyCase, snCase}, args...), len(args))
	// 4. 执行批量更新SQL
	sql := fmt.Sprintf(`UPDATE oms_order SET logistics_company = %s, logistics_sn = %s WHERE id IN %s`, companyCase, snCase, inPlaceholders)
	fmt.Println(sql)
	if err := tx.Debug().Exec(sql, args...).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("批量更新失败: %v", err)
	}

	return tx.Commit().Error
}

func (o *OrderService) UploadFile(header *multipart.FileHeader, noSave string, userId int) (file product.AfterSalesUpload, err error) {
	n, err := snowflake.NewNode(1)
	if err != nil {
		global.GVA_LOG.Error("创建id失败!", zap.Error(err))
	}
	s := strings.Split(header.Filename, ".")
	tag := s[len(s)-1]

	uuid := n.Generate()
	fileName := uuid.String()
	header.Filename = fmt.Sprintf("%s.%s", fileName, tag)
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header, userId, 1)
	if uploadErr != nil {
		panic(err)
	}

	f := product.AfterSalesUpload{
		Url:       filePath,
		Name:      header.Filename,
		Tag:       tag,
		Key:       key,
		SysUserId: userId,
		FileId:    uuid.Int64(),
	}
	if noSave == "0" {
		err = o.Upload(&f)
		fmt.Println(f)
		return f, err
	}
	return f, nil
}

func (o *OrderService) Upload(file *product.AfterSalesUpload) error {
	return global.GVA_DB.Create(file).Error
}

func (o *OrderService) CreateDealOrderApply(apply *product.AfterSalesApply) error {
	return global.GVA_DB.Debug().Create(apply).Error
}

func (o *OrderService) GetDealOrderList(searchInfo productReq.OrderDealSearchRequest) (list []*product.AfterSalesApply, total int64, err error) {
	// 处理分页参数默认值，避免页码或页大小为负数/零
	page := searchInfo.Page
	if page <= 0 {
		page = 1
	}
	pageSize := searchInfo.PageSize
	if pageSize <= 0 || pageSize > 100 { // 限制最大页大小，防止恶意请求
		pageSize = 10
	}
	limit := pageSize
	offset := pageSize * (page - 1)

	db := global.GVA_DB.Model(&product.AfterSalesApply{})
	if len(searchInfo.Contact) > 0 {
		db = db.Where("contact = ?", searchInfo.Contact)
	}
	if searchInfo.Status > 0 {
		db = db.Where("status = ?", searchInfo.Status-100)
	}

	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	} else {

		err = db.Limit(limit).Offset(offset).Debug().Find(&list).Error
	}

	return list, total, err
}

func (o *OrderService) UpdateDealOrderSynchronous(info productReq.UpdateDealOrderRequest) (err error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return err
	}

	err = tx.Model(&product.AfterSalesApply{}).Debug().Where("id = ?", info.DealId).UpdateColumn("status", info.Status).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新状态失败: %w", err)
	}

	err = tx.Model(&product.OrderItem{}).Debug().Where("id = ?", info.OrderItemId).UpdateColumn("status", 2).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新状态失败: %w", err)
	}

	// 更新Product同时更新ProductLadder
	return tx.Commit().Error

}

func (o *OrderService) GetOrderDealUploadImages(ids []int) (imagesList []product.AfterSalesUpload, err error) {
	db := global.GVA_DB.Model(&product.AfterSalesUpload{})
	db.Where("id in ?", ids).Find(&imagesList)
	return imagesList, err
}
