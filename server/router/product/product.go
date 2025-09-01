package product

import (
	v1 "cooller/server/api/v1"
	"cooller/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (e *ProductRouter) InitProductRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {

	homePublicRouterWithoutRecord := RouterPub.Group("product")

	productRouter := Router.Group("product").Use(middleware.OperationRecord())
	productRouterWithoutRecord := Router.Group("product")
	homeApi := v1.ApiGroupApp.ProductApiGroup.ProductApi
	{

		productRouter.POST("create", homeApi.CreateProduct)
		productRouter.PUT("update", homeApi.UpdateProduct)
		productRouter.PUT("updateKeyword", homeApi.UpdateProductForKeyword)

		productRouter.POST("brand", homeApi.CreateProductBrand)
		productRouter.PUT("brand", homeApi.UpdateProductBrand)
		productRouter.DELETE("brand", homeApi.DeleteHomeProductBrand)
		productRouter.PUT("brandState", homeApi.UpdateProductBrandOnlineStat)

		productRouter.POST("attributeCategory", homeApi.CreateProductAttributeCategory)
		productRouter.PUT("attributeCategory", homeApi.UpdateProductAttributeCategory)
		productRouter.DELETE("attributeCategory", homeApi.DeleteProductAttributeCategory)
		productRouter.POST("attribute", homeApi.CreateProductAttribute)
		productRouter.PUT("attribute", homeApi.UpdateProductAttribute)
		productRouter.DELETE("attribute", homeApi.DeleteProductAttribute)
		productRouter.POST("productCategory", homeApi.CreateProductCategory)
		productRouter.PUT("productCategory", homeApi.UpdateProductCategory)
		productRouter.DELETE("productCategory", homeApi.DeleteProductCategory)
		productRouter.PUT("sku", homeApi.UpdateSKUStock)

	}
	{
		productRouterWithoutRecord.DELETE("deletes", homeApi.DeleteProducts)
	}
	{
		homePublicRouterWithoutRecord.GET("list", homeApi.GetProductList) // 获取商品列表
		homePublicRouterWithoutRecord.GET("simpleList", homeApi.GetBySimpleProductList)
		homePublicRouterWithoutRecord.GET("brand", homeApi.GetProductBrandList)

		homePublicRouterWithoutRecord.GET("brandDetail", homeApi.GetProductBrandByID)
		homePublicRouterWithoutRecord.GET("productDetail", homeApi.GetProductByID)
		homePublicRouterWithoutRecord.GET("attributeCategory", homeApi.GetProductAttributeCategoryList)
		homePublicRouterWithoutRecord.GET("attribute", homeApi.GetProductAttributeListByCategoryId)
		homePublicRouterWithoutRecord.GET("productCategory", homeApi.GetProductCategoryList)
		homePublicRouterWithoutRecord.GET("productList", homeApi.GetProductListByOnlyIDWithSort) // 获取商品列表
		homePublicRouterWithoutRecord.GET("allCategory", homeApi.GetProductAllCategory)          // 获取商品分类列表
		homePublicRouterWithoutRecord.GET("sku", homeApi.GetSkuStockByProductID)
	}
}
