package wechat

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	"cooller/server/model/example"
	"cooller/server/model/product"
	"cooller/server/model/wechat"
	wechatRequest "cooller/server/model/wechat/request"
	wechatRes "cooller/server/model/wechat/response"
	"cooller/server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type HomeApi struct{}

// CreateRecommendProducts 创建推荐商品
func (e *HomeApi) CreateRecommendProducts(c *gin.Context) {
	var recommendProducts wechatRequest.AddRecommendProductRequest
	err := c.ShouldBindJSON(&recommendProducts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = wechatService.CreateRecommendProduct(&recommendProducts.Products)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateRecommendProducts 更新人气推荐商品
func (e *HomeApi) UpdateRecommendProducts(c *gin.Context) {
	var product wechatRequest.UpdateIdsKeywordRequest
	err := c.ShouldBindJSON(&product)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err3 := wechatService.UpdateRecommendProducts(&product)
	if err3 != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err3))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// UpdateRecommendProductSortById 更新人气推荐商品
func (e *HomeApi) UpdateRecommendProductSortById(c *gin.Context) {
	var info request.SortUpdateInfo
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err3 := wechatService.UpdateRecommendProductSortById(&info)
	if err3 != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err3))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteRecommendProducts 删除猜你喜欢商品
func (e *HomeApi) DeleteRecommendProducts(c *gin.Context) {
	var idsReq request.IdsReq
	err := c.ShouldBindJSON(&idsReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err3 := wechatService.DeleteRecommendProducts(idsReq.Ids)
	if err3 != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err3))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *HomeApi) CreateHomeAdvertise(c *gin.Context) {
	var home wechat.Advertise
	err := c.ShouldBindJSON(&home)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//home.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	err = wechatService.CreateHomeAdvertise(home)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteExaCustomer
// @Tags      ExaCustomer
// @Summary   删除客户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      example.ExaCustomer            true  "客户ID"
// @Success   200   {object}  response.Response{msg=string}  "删除客户"
// @Router    /customer/customer [delete]
func (e *HomeApi) DeleteHomeAdvertise(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	advertise, err := wechatService.GetHomeAdvertiseById(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	var file example.ExaFileUploadAndDownload
	file.Url = advertise.Pic
	if err := fileUploadAndDownloadService.DeleteFile(file); err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("阿里云图片删除失败:%s", file.Url), zap.Error(err))
	}

	err = wechatService.DeleteHomeAdvertise(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateExaCustomer
// @Tags      ExaCustomer
// @Summary   更新客户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      example.ExaCustomer            true  "客户ID, 客户信息"
// @Success   200   {object}  response.Response{msg=string}  "更新客户信息"
// @Router    /customer/customer [put]
func (e *HomeApi) UpdateHomeAdvertise(c *gin.Context) {
	var home wechat.Advertise
	err := c.ShouldBindJSON(&home)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//home.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	err = wechatService.UpdateHomeAdvertise(&home)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *HomeApi) UpdateHomeAdvertiseOnlineState(c *gin.Context) {
	var advertise wechatRequest.UpdateIdsKeywordRequest
	err := c.ShouldBindJSON(&advertise)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = wechatService.UpdateHomeAdvertiseByIdForKeyword(&advertise)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetHomeAdvertiseList
// @Tags      ExaCustomer
// @Summary   分页获取权限客户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取权限客户列表,返回包括列表,总数,页码,每页数量"
// @Router    /customer/customerList [get]
func (e *HomeApi) GetHomeAdvertiseList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	homeList, total, err := wechatService.GetHomeAdvertiseInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     homeList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// FindValidFlashPromotion 获取有效的秒杀活动
func FindValidFlashPromotion() (promotionValidList []*wechat.FlashPromotion) {
	promotionList, err := wechatService.GetOnlineHomeFlashPromotionInfoList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		return promotionValidList
	}
	for _, promotion := range promotionList {
		if promotion.Status > 0 {
			now := time.Now()
			if now.After(promotion.StartDate) && now.Before(promotion.EndDate) {
				promotionValidList = append(promotionValidList, promotion)
				return promotionValidList
			}
		}
	}
	return promotionValidList
}

func FindCurrentFlashPromotionSession() (sessionRes wechat.FlashPromotionSession, err error) {
	flashSessionList, err := wechatService.GetFlashSessionList()
	if err != nil || len(flashSessionList) < 1 {
		global.GVA_LOG.Error("GetFlashSessionList获取失败!", zap.Error(err))
		return sessionRes, fmt.Errorf("Not found FlashPromotionSession!")
	} else {
		now := time.Now()
		year, month, day := now.Date()
		for _, session := range flashSessionList {
			if session.Status > 0 {
				startTime := time.Date(year, month, day, session.StartTime.Hour(), session.StartTime.Minute(), session.StartTime.Second(), 0, time.Local)
				EndTime := time.Date(year, month, day, session.EndTime.Hour(), session.EndTime.Minute(), session.EndTime.Second(), 0, time.Local)
				if now.After(startTime) && now.Before(EndTime) {
					return session, nil
				}
			}
		}
		return sessionRes, fmt.Errorf("Not found!")
	}
}

func (e *HomeApi) FindNextFlashPromotionSession() (sessionRes wechat.FlashPromotionSession, err error) {
	flashSessionList, err := wechatService.GetFlashSessionList()
	if err != nil || len(flashSessionList) < 1 {
		return sessionRes, fmt.Errorf("Not found FlashPromotionSession!")
	} else {
		now := time.Now()
		year, month, day := now.Date()
		nextIndex := 0
		mini := 0
		isFind := false
		nextTime := now.AddDate(0, 0, 1)
		miniTime := now.AddDate(0, 0, 1)
		for i, session := range flashSessionList {
			if session.Status > 0 {
				startTime := time.Date(year, month, day, session.StartTime.Hour(), session.StartTime.Minute(), session.StartTime.Second(), 0, time.Local)
				if now.Before(startTime) && nextTime.After(now) && nextTime.After(startTime) {
					nextIndex = i
					nextTime = startTime
					isFind = true
				}
				if miniTime.After(startTime) {
					mini = i
					miniTime = startTime
				}
			}
		}
		if !isFind {
			return flashSessionList[mini], nil
		} else {
			return flashSessionList[nextIndex], nil
		}
	}
}

func (e *HomeApi) GetAllWechatContent(c *gin.Context) {
	homeList, err := wechatService.GetOnlineHomeAdvertiseInfoList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	brandList, err := wechatService.GetOnlineHomeBrandInfoList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	newProductList, err := wechatService.GetOnlineNewProductInfoList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	var homeFlashPromotion wechatRes.HomeFlashResponse
	promotionValidList := FindValidFlashPromotion()
	if len(promotionValidList) > 0 {
		currentFlash, err := FindCurrentFlashPromotionSession()
		if err == nil {
			homeFlashPromotion.StartTime = currentFlash.StartTime
			homeFlashPromotion.EndTime = currentFlash.EndTime
			for _, flashPromotion := range promotionValidList {
				sessionProductList, err := wechatService.GetFlashPromotionProductRelationListById(flashPromotion.ID, currentFlash.ID)
				if err != nil {
					global.GVA_LOG.Error("GetFlashPromotionProductRelationListById获取失败!", zap.Error(err))
				} else {
					productList := make([]product.Product, 0)
					for _, sessionProduct := range sessionProductList {
						productData, _, _ := CalculateProductPromotionPrice(sessionProduct.Product, sessionProductList)
						productList = append(productList, productData)
					}
					homeFlashPromotion.ProductList = productList
				}
			}
		}
		nextFlash, err1 := e.FindNextFlashPromotionSession()
		if err1 == nil {
			homeFlashPromotion.NextStartTime = nextFlash.StartTime
			homeFlashPromotion.NextEndTime = nextFlash.EndTime
		}
	}
	var pageInfo request.PageInfo
	pageInfo.Page = 1
	pageInfo.PageSize = 10
	hotProductList, err := wechatService.GetOnlineRecommendProductListInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	groupBuyList, err := wechatService.GetGroupBuyProductList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	var groupBuy wechatRes.GroupBuyResp
	for _, item := range groupBuyList {
		var groupBuyProduct wechatRes.SampleProductInfo
		groupBuyProduct.ID = item.ID
		groupBuyProduct.ProductId = item.ProductId
		groupBuyProduct.Name = item.Product.Name
		groupBuyProduct.Price = item.Price
		groupBuyProduct.OriginalPrice = item.Product.OriginalPrice
		groupBuyProduct.AlbumPics = item.Product.AlbumPics
		groupBuyProduct.Pic = item.Product.Pic
		groupBuyProduct.Sales = item.Product.Sale
		groupBuyProduct.Percent = item.Percent
		groupBuy.Groups = append(groupBuy.Groups, groupBuyProduct)
	}
	response.OkWithDetailed(wechatRes.HomeContentResponse{
		AdvertiseList:      homeList,
		BrandList:          brandList,
		NewProductList:     newProductList,
		HotProductList:     hotProductList,
		HomeFlashPromotion: homeFlashPromotion,
		GroupBuy:           groupBuy,
	}, "获取成功", c)
}

func (e *HomeApi) GetRecommendProductList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	recommendProductList, err := wechatService.GetOnlineRecommendProductListInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	var list []product.Product
	for _, recommend := range recommendProductList {
		if recommend.Product.ID != 0 {
			list = append(list, recommend.Product)
		}
	}
	response.OkWithDetailed(response.AllResult{
		List: list,
	}, "获取成功", c)
}

func (e *HomeApi) GetRecommendProductListByCondition(c *gin.Context) {
	var pageInfo wechatRequest.RecommendProductSearchInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	recommendProductList, total, err := wechatService.GetRecommendProductListByCondition(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	var list []product.Product
	for _, recommend := range recommendProductList {
		if recommend.Product.ID != 0 {
			list = append(list, recommend.Product)
		}
	}
	response.OkWithDetailed(response.PageResult{
		List:     recommendProductList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (e *HomeApi) GetNewProductList(c *gin.Context) {
	var info request.NameAndStateSearchInfo
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	newList, total, err := wechatService.GetNewProductInfoList(info)
	if err != nil {
		global.GVA_LOG.Error("获取新品推荐库失败!", zap.Error(err))
		response.FailWithMessage("获取新品推荐失败", c)
		return
	}

	productIds, err := wechatService.GetNewProductIdsList()
	if err != nil {
		global.GVA_LOG.Error("获取新品推荐ids失败!", zap.Error(err))
		response.FailWithMessage("获取新品推荐ids失败", c)
		return
	}

	response.OkWithDetailed(wechatRes.NewProductIdsList{
		List:       newList,
		Total:      total,
		Page:       info.Page,
		PageSize:   info.PageSize,
		ProductIds: productIds,
	}, "获取成功", c)
}

func (e *HomeApi) CreateNewProduct(c *gin.Context) {
	var newProducts wechatRequest.NewProductsRequest
	err := c.ShouldBindJSON(&newProducts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//category.SysUserId = utils.GetUserAuthorityId(c)
	for _, newProduct := range newProducts.NewProducts {
		newProduct.RecommendStatus = 1
		newProduct.Sort = 0
	}

	err = wechatService.CreateNewProductByBatch(newProducts.NewProducts)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateProductCategory 更新商品分类
func (e *HomeApi) UpdateNewProductSort(c *gin.Context) {
	var sortInfo request.SortUpdateInfo
	err := c.ShouldBindJSON(&sortInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//home.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	err = wechatService.UpdateNewProductForSort(&sortInfo)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *HomeApi) UpdateNewProductState(c *gin.Context) {
	var idsKeyword wechatRequest.UpdateIdsKeywordRequest
	err := c.ShouldBindJSON(&idsKeyword)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if idsKeyword.Value != 100 && idsKeyword.Value != 101 {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
	}
	//home.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	err = wechatService.UpdateNewProductForRecommendStatus(&idsKeyword)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteProductCategory 删除商品分类
func (e *HomeApi) DeleteNewProduct(c *gin.Context) {
	var idsReq request.IdsReq
	err := c.ShouldBindQuery(&idsReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = wechatService.DeleteNewProduct(idsReq.Ids)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetFlashPromotionProductList 获取当前活动中的商品关联表
func GetFlashPromotionProductList() (list []wechat.FlashPromotionProductRelation) {
	promotionValidList := FindValidFlashPromotion()
	if len(promotionValidList) > 0 {
		currentFlash, err := FindCurrentFlashPromotionSession()
		if err == nil {
			for _, flashPromotion := range promotionValidList {
				sessionProductList, err := wechatService.GetFlashPromotionProductRelationListById(flashPromotion.ID, currentFlash.ID)
				if err != nil {
					global.GVA_LOG.Error("GetFlashPromotionProductRelationListById获取失败!", zap.Error(err))
				} else {
					list = append(list, sessionProductList...)
					return list
				}
			}
		}
	}
	return list
}

func CalculateProductPromotionPrice(product product.Product, list []wechat.FlashPromotionProductRelation) (p product.Product, promotionMessage string, reduceAmount float32) {
	promotion := product.PromotionType
	// TODO: 各个优惠计算
	if promotion == 0 {
		promotionMessage = "无优惠"
	} else if promotion == 1 {
		promotionMessage = "特惠促销"
		now := time.Now().Unix()
		startTime := product.PromotionStartTime.Unix()
		endTime := product.PromotionEndTime.Unix()
		if now >= startTime && now <= endTime {
			product.Price = product.PromotionPrice
		}
	} else if promotion == 2 {
		promotionMessage = "会员优惠"
	} else if promotion == 3 {
		promotionMessage = "阶梯价格"
	} else if promotion == 4 {
		fullReductionList, err := productService.GetProductFullReductionByProductId(product.ID)
		if err != nil {
			global.GVA_LOG.Error("获取购物车物品失败!", zap.Error(err))
		}
		reductionStr := ""
		for _, reduction := range fullReductionList {
			if product.Price >= reduction.FullPrice && reduceAmount < reduction.ReducePrice {
				reduceAmount = reduction.ReducePrice
			}
			reductionStr += fmt.Sprintf("满%.2f元，减%.2f元,", reduction.FullPrice, reduction.ReducePrice)
		}
		promotionMessage = fmt.Sprintf("满减优惠：%s", reductionStr)
	} else if promotion == 5 { // 秒杀
		promotionMessage = "限时优惠"
		if list != nil && len(list) > 0 {
			for _, flashProduct := range list {
				if flashProduct.ProductId == product.ID {
					product.Price = flashProduct.FlashPromotionPrice
				}
			}
		} else {
			sessionProductList := GetFlashPromotionProductList()
			for _, relation := range sessionProductList {
				if product.ID == relation.ProductId {
					product.Price = relation.FlashPromotionPrice
				}
			}
		}
	}

	return product, promotionMessage, reduceAmount
}
