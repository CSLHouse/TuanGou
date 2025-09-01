package product

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/product"
	wechatReq "cooller/server/model/wechat/request"
	date_conversion "cooller/server/utils/timer"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type OrderService struct{}

func (o *OrderService) GetProductCartById(userId int, id int) (cartItem product.CartItem, err error) {
	db := global.GVA_DB.Model(&product.CartItem{})
	db.Where("user_id = ? and id = ?", userId, id).First(&cartItem)
	return cartItem, err
}

func (o *OrderService) GetProductCartByIds(userId int, ids []int) (cartItem []product.CartItem, err error) {
	db := global.GVA_DB.Model(&product.CartItem{})
	db.Debug().Preload("Product").Preload("SkuStock").Where("user_id = ? and id in ?", userId, ids).First(&cartItem)
	return cartItem, err
}

func (o *OrderService) GetProductTmpCartByIds(userId int, ids []int) (cartItem []product.CartTmpItem, err error) {
	db := global.GVA_DB.Model(&product.CartTmpItem{})
	db.Debug().Preload("Product").Preload("SkuStock").Where("user_id = ? and id in ?", userId, ids).First(&cartItem)
	return cartItem, err
}

func (o *OrderService) CreateOrder(e *product.Order) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Create(&e).Error
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

func (o *OrderService) GetProductOrderListByStatus(searchInfo request.StateInfo) (list []product.Order, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)
	state := searchInfo.State
	db := global.GVA_DB.Model(&product.Order{})

	var cmdList []interface{}
	var cmdString string
	if len(searchInfo.OrderSn) > 0 {
		cmdList = append(cmdList, strings.TrimSpace(searchInfo.OrderSn))
		cmdString = "order_sn = ?"
	}
	if len(searchInfo.ReceiverKeyword) > 0 {
		keyword := "%" + strings.TrimSpace(searchInfo.ReceiverKeyword) + "%"
		if len(cmdList) >= 1 {
			cmdList = append(cmdList, keyword)
			cmdList = append(cmdList, keyword)
			cmdString += " and receiver_name like ? or receiver_phone like ?"
		} else {
			cmdList = append(cmdList, keyword)
			cmdList = append(cmdList, keyword)
			cmdString += "receiver_name like ? or receiver_phone like ?"
		}
	}

	if searchInfo.OrderType > 0 {
		if len(cmdList) >= 1 {
			cmdList = append(cmdList, searchInfo.OrderType-100)
			cmdString += " and order_type = ?"
		} else {
			cmdList = append(cmdList, searchInfo.OrderType-100)
			cmdString += "order_type = ?"
		}
	}
	if len(searchInfo.CreateTime) > 0 {
		thatDay := date_conversion.ParseStringDate(searchInfo.CreateTime)
		nextDay := thatDay.AddDate(0, 0, 1)

		if len(cmdList) >= 1 {
			cmdList = append(cmdList, thatDay)
			cmdList = append(cmdList, nextDay)
			cmdString += " and created_at between ? and ?"
		} else {
			cmdList = append(cmdList, thatDay)
			cmdList = append(cmdList, nextDay)
			cmdString += "created_at between ? and ?"
		}
	}

	switch state {
	case -1:
		{
			if len(cmdList) > 0 {
				err = db.Where(cmdString, cmdList...).Count(&total).Error
				if err != nil {
					return list, total, err
				} else {
					err = db.Where(cmdString, cmdList...).Limit(limit).Offset(offset).Preload("OrderItemList").Order("id desc").Find(&list).Error
				}
			} else {
				err = db.Count(&total).Error
				if err != nil {
					return list, total, err
				} else {
					err = db.Limit(limit).Offset(offset).Preload("OrderItemList").Order("id desc").Find(&list).Error
				}
			}
			return list, total, err
		}
	case 0, 3, 4:
		{
			if len(cmdString) > 0 {
				cmdString = fmt.Sprintf(" %s and status = %d", cmdString, state)
			} else {
				cmdString = fmt.Sprintf("status = %d", state)
			}
			err = db.Debug().Where(cmdString, cmdList...).Count(&total).Error
			if err != nil {
				return list, total, err
			} else {

				err = db.Debug().Where(cmdString, cmdList...).Limit(limit).Offset(offset).Preload("OrderItemList").Order("id desc").Find(&list).Error
			}
			return list, total, err
		}
	case 1, 2:
		{
			if len(cmdString) > 0 {
				cmdString = fmt.Sprintf("%s and status = 1 or status = 2", cmdString)
			} else {
				cmdString = "status = 1 or status = 2"
			}
			err = db.Where(cmdString, cmdList...).Count(&total).Error
			if err != nil {
				return list, total, err
			} else {
				err = db.Debug().Where(cmdString, cmdList...).Limit(limit).Offset(offset).Preload("OrderItemList").Order("id desc").Find(&list).Error
			}
			return list, total, err
		}
	default:
		return list, total, err
	}
}

func (o *OrderService) UpdateOrderStatus(e *wechatReq.PaySuccessRequest, status int) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Debug().Where("id = ?", e.OrderId).Updates(map[string]interface{}{"pay_type": e.PayType, "status": status}).Error
	return err
}

func (o *OrderService) UpdateOrderStatusById(orderId int, status int) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Debug().Where("id = ?", orderId).UpdateColumn("status", status).Error
	return err
}

func (o *OrderService) UpdateOrderStatusByOrderSn(orderSn *string, status int) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Debug().Where("order_sn = ?", orderSn).UpdateColumn("status", status).Error
	return err
}

func (o *OrderService) UpdateOrderPrepayId(id int, prepayId string) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Debug().Where("id = ?", id).UpdateColumn("prepay_id", prepayId).Error
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

func (o *OrderService) UpdateOrderReceiverInfo(e *wechatReq.OrderReceiveAddress) (err error) {
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

func (o *OrderService) UpdateOrderMoneyInfo(info *wechatReq.OrderMoneyInfo) (err error) {
	db := global.GVA_DB.Model(&product.Order{})
	err = db.Where("id = ?", info.OrderId).UpdateColumn("discount_amount", info.DiscountAmount).Error
	return err
}

func (o *OrderService) UpdateOrderNoteInfo(info *wechatReq.OrderNoteInfo) (err error) {
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
