package product

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/common/response"
	"cooller/server/model/product"
	"cooller/server/model/wechat"
	wechatRequest "cooller/server/model/wechat/request"
	wechatRes "cooller/server/model/wechat/response"
	"cooller/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type ProductApi struct{}

func (e *ProductApi) CreateProduct(c *gin.Context) {
	var productInfo product.Product
	err := c.ShouldBindJSON(&productInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = productService.CreateHomeProduct(&productInfo)
	if err != nil {
		global.GVA_LOG.Error("创建产品失败!", zap.Error(err))
		response.FailWithMessage("创建产品失败", c)
		return
	}
	// 新品
	if productInfo.NewStatus == 1 {
		var newProduct wechat.NewProduct
		newProduct.ProductId = productInfo.ID
		newProduct.ProductName = productInfo.Name
		newProduct.RecommendStatus = productInfo.PublishStatus
		err = wechatService.CreateNewProduct(&newProduct)
		if err != nil {
			global.GVA_LOG.Error("创建新品失败!", zap.Error(err))
			response.FailWithMessage("创建新品失败", c)
			return
		}
	}
	// 推荐
	if productInfo.RecommandStatus == 1 {
		var recommendProduct []wechat.RecommendProduct
		var recommend wechat.RecommendProduct
		recommend.ProductId = productInfo.ID
		recommend.ProductName = productInfo.Name
		recommend.RecommendStatus = productInfo.PublishStatus
		recommendProduct = append(recommendProduct, recommend)
		err = wechatService.CreateRecommendProduct(&recommendProduct)
		if err != nil {
			global.GVA_LOG.Error("创建推荐商品失败!", zap.Error(err))
			response.FailWithMessage("创建推荐商品失败", c)
			return
		}
	}

	response.OkWithMessage("创建成功", c)
}

// GetProductByID 获取商品详情
func (e *ProductApi) GetProductByID(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	productData, err := productService.GetProductByID(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("获取商品失败!", zap.Error(err))
		response.FailWithMessage("获取商品失败", c)
		return
	}
	// 买家id是以字符串方式存储：12,46,88,需要分割且转换一下
	var buyersIdList []int
	var avatarsList wechatRequest.BuyersInfo
	if len(productData.Buyers) > 0 {
		buyersIdStrList := strings.Split(productData.Buyers, ",")
		for _, buyersIdStr := range buyersIdStrList {
			buyersId, err := strconv.Atoi(buyersIdStr)
			if err != nil {
				global.GVA_LOG.Error("获取买家头像失败!", zap.Error(err))
				continue
			}
			buyersIdList = append(buyersIdList, buyersId)
		}
		if len(buyersIdList) > 0 {
			avatarsList, err = productService.GetBuysAvatarsList(buyersIdList)
			if err != nil {
				global.GVA_LOG.Error("获取买家头像失败!", zap.Error(err))
				response.FailWithMessage("获取买家头像失败", c)
				return
			}
		}
	}

	var productDetails wechatRes.ProductDetails
	productDetails.Product = productData
	productDetails.BuyersList = avatarsList
	response.OkWithData(productDetails, c)
}

func (e *ProductApi) GetProductList(c *gin.Context) {
	var pageInfo wechatRequest.ProductSearchInfo
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
	productList, total, err := productService.GetProductList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     productList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetBySimpleProductList 通过关键字模糊查找商品
func (e *ProductApi) GetBySimpleProductList(c *gin.Context) {
	var keyword request.KeyWordInfo
	err := c.ShouldBindQuery(&keyword)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	productList, err := productService.GetBySimpleProductList(keyword)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.AllResult{
		List: productList,
	}, "获取成功", c)
}

func (e *ProductApi) GetProductListByOnlyID(c *gin.Context) {
	var pageInfo request.KeySearchInfo
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
	productList, total, err := productService.GetProductListByOnlyID(&pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     productList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
func (e *ProductApi) DeleteProducts(c *gin.Context) {
	var idsInfo request.IdsReq
	err := c.ShouldBindJSON(&idsInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = productService.DeleteProducts(idsInfo.Ids)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (e *ProductApi) GetProductListByOnlyIDWithSort(c *gin.Context) {
	var pageInfo request.SortSearchInfo
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
	productList, total, err := productService.GetProductListByOnlyIDWithSort(&pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     productList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
func (e *ProductApi) UpdateProduct(c *gin.Context) {
	var product product.Product
	var err = c.ShouldBindJSON(&product)
	if err != nil {
		response.FailWithMessage("Json Parse Error!!", c)
		return
	}
	// 修改优惠方式
	if product.PromotionType == 0 { // 无优惠
		product.PromotionPerLimit = 0
		product.PromotionStartTime = nil
		product.PromotionEndTime = nil
		product.ProductLadderList = nil
	} else if product.PromotionType == 1 { //特惠促销
		product.ProductLadderList = nil
	} else if product.PromotionType == 2 { // 会员价格
		product.PromotionPerLimit = 0
		product.PromotionStartTime = nil
		product.PromotionEndTime = nil
		product.ProductLadderList = nil
	} else if product.PromotionType == 3 { // 阶梯价格
		product.PromotionPerLimit = 0
		product.PromotionStartTime = nil
		product.PromotionEndTime = nil
	} else if product.PromotionType == 4 { // 满减价格

	}
	err3 := productService.UpdateHomeProductSynchronous(&product)
	if err3 != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err3))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// UpdateProductForKeyword 更新商品
func (e *ProductApi) UpdateProductForKeyword(c *gin.Context) {
	var product wechatRequest.UpdateIdsKeywordRequest
	var err = c.ShouldBindJSON(&product)
	if err != nil {
		response.FailWithMessage("Json Parse Error!!", c)
		return
	}

	err3 := productService.UpdateProductForKeyword(&product)
	if err3 != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err3))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetProductAttribute 获取商品分类
//func (e *ProductApi) GetProductAttributeCategoryList(c *gin.Context) {
//	var pageInfo request.PageInfo
//	err := c.ShouldBindQuery(&pageInfo)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	if pageInfo.Page == 0 {
//		pageInfo.Page = 1
//	}
//	if pageInfo.PageSize == 0 {
//		pageInfo.PageSize = 100
//	}
//	categoryList, total, err := wechatService.GetProductAttributeCategoryList(pageInfo)
//	if err != nil {
//		global.GVA_LOG.Error("获取失败!", zap.Error(err))
//		response.FailWithMessage("获取失败"+err.Error(), c)
//		return
//	}
//
//	response.OkWithDetailed(response.PageResult{
//		List:     categoryList,
//		Total:    total,
//		Page:     pageInfo.Page,
//		PageSize: pageInfo.PageSize,
//	}, "获取成功", c)
//}

// GetProductBrand 获取品牌详情
func (e *ProductApi) GetProductBrandByID(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := productService.GetProductBrand(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetProductBrandList 获取商品品牌列表
func (e *ProductApi) GetProductBrandList(c *gin.Context) {
	var pageInfo wechatRequest.BrandSearchInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 100
	}
	brandList, total, err := productService.GetProductBrandList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     brandList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (e *ProductApi) CreateProductBrand(c *gin.Context) {
	var brand product.Brand

	err := c.ShouldBindJSON(&brand)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = productService.CreateHomeProductBrand(&brand)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (e *ProductApi) UpdateProductBrand(c *gin.Context) {
	var brand product.Brand
	err := c.ShouldBindJSON(&brand)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//home.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	err = productService.UpdateHomeBrand(&brand)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *ProductApi) UpdateProductBrandOnlineStat(c *gin.Context) {
	var brand wechatRequest.UpdateIdsKeywordRequest
	err := c.ShouldBindJSON(&brand)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = productService.UpdateHomeBrandByIdForKeyword(&brand)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (e *ProductApi) DeleteHomeProductBrand(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindQuery(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = productService.DeleteHomeProductBrand(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetProductAttributeCategoryList 获取商品属性列表
func (e *ProductApi) GetProductAttributeCategoryList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 100
	}
	attributeList, total, err := productService.GetProductAttributeCategoryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     attributeList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// CreateProductAttributeCategory 创建商品属性分类
func (e *ProductApi) CreateProductAttributeCategory(c *gin.Context) {
	var attribute product.ProductAttributeCategory

	err := c.ShouldBindJSON(&attribute)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = productService.CreateProductAttributeCategory(&attribute)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateProductAttributeCategory 更新商品属性分类
func (e *ProductApi) UpdateProductAttributeCategory(c *gin.Context) {
	var attribute product.ProductAttributeCategory
	err := c.ShouldBindJSON(&attribute)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//home.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	err = productService.UpdateProductAttributeCategory(&attribute)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteProductAttributeCategory 删除商品属性分类
func (e *ProductApi) DeleteProductAttributeCategory(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = productService.DeleteProductAttributeCategory(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// CreateProductAttribute 创建商品属性参数
func (e *ProductApi) CreateProductAttribute(c *gin.Context) {
	var attribute product.ProductAttribute

	err := c.ShouldBindJSON(&attribute)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = productService.CreateProductAttributeSynchronous(&attribute)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

// UpdateProductAttribute 更新商品属性参数
func (e *ProductApi) UpdateProductAttribute(c *gin.Context) {
	var attribute product.ProductAttribute
	err := c.ShouldBindJSON(&attribute)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//home.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	err = productService.UpdateProductAttribute(&attribute)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteProductAttribute 删除商品属性参数
func (e *ProductApi) DeleteProductAttribute(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = productService.DeleteProductAttribute(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetProductAttributeListByCategoryId 根据商品分类id获取商品属性参数列表
func (e *ProductApi) GetProductAttributeListByCategoryId(c *gin.Context) {
	var pageInfo request.TagSearchInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if pageInfo.Tag < 1 {
		response.FailWithMessage("商品属性分类id不可为空", c)
		return
	}
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 100
	}
	attributeList, total, err := productService.GetProductAttributeList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     attributeList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// CreateProductCategory 创建商品分类
func (e *ProductApi) CreateProductCategory(c *gin.Context) {
	var category product.ProductCategory
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	category.SysUserId = utils.GetUserAuthorityId(c)
	err = productService.CreateProductCategory(&category)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateProductCategory 更新商品分类
func (e *ProductApi) UpdateProductCategory(c *gin.Context) {
	var category product.ProductCategory
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//home.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	err = productService.UpdateProductCategory(&category)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteProductCategory 删除商品分类
func (e *ProductApi) DeleteProductCategory(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = productService.DeleteProductCategory(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetProductCategoryList 根据商品分类id获取商品分类
func (e *ProductApi) GetProductCategoryList(c *gin.Context) {
	var pageInfo request.TagSearchInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 100
	}
	attributeList, total, err := productService.GetProductCategoryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     attributeList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetProductCategoryTreeList 获取商品分类树
func (e *ProductApi) GetProductAllCategory(c *gin.Context) {
	categoryList, err := productService.GetAllProductCategoryList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithData(categoryList, c)
}

func (e *ProductApi) GetSkuStockByProductID(c *gin.Context) {
	var info request.KeyWordInfo
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(info, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	skuStockList, err := productService.GetProductSKUStockByProductId(info.ID, info.Keyword)
	if err != nil {
		global.GVA_LOG.Error("获取商品sku库存失败!", zap.Error(err))
		response.FailWithMessage("获取sku库存失败", c)
		return
	}
	response.OkWithData(skuStockList, c)
}

func (e *ProductApi) UpdateSKUStock(c *gin.Context) {
	var stockList []product.SkuStock
	err := c.ShouldBindJSON(&stockList)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//home.SysUserAuthorityID = utils.GetUserAuthorityId(c)
	for _, stock := range stockList {
		err = productService.UpdateSKUStock(&stock)
		if err != nil {
			global.GVA_LOG.Error("更新失败!", zap.Error(err))
			response.FailWithMessage("更新失败", c)
			return
		}
	}
	response.OkWithMessage("更新成功", c)
}
