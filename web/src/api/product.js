import service from '@/utils/request'

export const getProductList = (params) => {
    return service({
        url: '/product/list',
        method: 'get',
        params
    })
}

export const getSimpleList = (params) => {
    return service({
        url: '/product/simpleList',
        method: 'get',
        params
    })
}

export const deleteProducts = (data) => {
    return service({
        url: '/product/deletes',
        method: 'delete',
        data,
    })
}

export const getProductDetail = (params) => {
    return service({
        url: '/product/productDetail',
        method: 'get',
        params
    })
}

export const updateProductKeyword = (data) => {
    return service({
        url: '/product/updateKeyword',
        method: 'put',
        data
    })
}

export const updateProduct = (data) => {
    return service({
        url: '/product/update',
        method: 'put',
        data
    })
}

export const getBrandList = (params) => {
    return service({
        url: '/product/brand',
        method: 'get',
        params
    })
}

export const createProduct = (data) => {
    return service({
        url: '/product/create',
        method: 'post',
        data
    })
}

export const createProductBrand = (data) => {
    return service({
        url: '/product/brand',
        method: 'post',
        data
    })
}

export const updateProductBrand = (data) => {
    return service({
        url: '/product/brand',
        method: 'put',
        data
    })
}

export const updateBrandByIdForState = (data) => {
    return service({
        url: '/product/brandState',
        method: 'put',
        data
    })
}

export const deleteProductBrand = (params) => {
    return service({
        url: '/product/brand',
        method: 'delete',
        params: params,
    })
}

// 获取首页轮播广告表
export const getAdvertiseList = (params) => {
    return service({
        url: '/product/advertiseList',
        method: 'get',
        params
    })
}

// 创建首页轮播广告
export const createAdvertise = (data) => {
    return service({
      url: '/product/advertise',
      method: 'post',
      data
    })
}

// 更新首页轮播广告
export const updateAdvertise = (data) => {
    return service({
        url: '/product/advertise',
        method: 'put',
        data
    })
}

export const deletedvertise = (params) => {
    return service({
        url: '/product/advertise',
        method: 'delete',
        params: params,
    })
}

// 更新首页轮播广告状态
export const updateAdvertiseByIdForState = (data) => {
    return service({
        url: '/product/advertiseState',
        method: 'put',
        data
    })
}


// 获取首页推荐专题表 猜你喜欢
export const getRecommendProductList = (params) => {
    return service({
        url: '/product/recommendProduct',
        method: 'get',
        params
    })
}

// 获取首页推荐专题表排序 猜你喜欢
export const updateRecommendProductByIdForSort = (data) => {
    return service({
        url: '/product/updateRecommendSort',
        method: 'post',
        data
    })
}

// 获取首页推荐专题表 猜你喜欢
export const addRecommendProductList = (data) => {
    return service({
        url: '/product/recommendProduct',
        method: 'post',
        data
    })
}
// 更新首页推荐专题表
export const updateRecommendProduct = (data) => {
    return service({
        url: '/product/recommendProduct',
        method: 'put',
        data
    })
}

// 删除首页推荐专题表
export const deleteRecommendProduct = (data) => {
    return service({
        url: '/product/recommendProduct',
        method: 'delete',
        data
    })
}

// 获取商品属性分类列表
export const getProductAttributeCategoryList = (params) => {
    return service({
        url: '/product/attributeCategory',
        method: 'get',
        params
    })
}

// 创建商品属性分类列表
export const createProductAttributeCategory = (data) => {
    return service({
      url: '/product/attributeCategory',
      method: 'post',
      data
    })
}

// 更新商品属性分类列表
export const updateProductAttributeCategory = (data) => {
    return service({
        url: '/product/attributeCategory',
        method: 'put',
        data
    })
}

// 删除商品属性分类列表
export const deleteProductAttributeCategory = (data) => {
    return service({
      url: '/product/attributeCategory',
      method: 'delete',
      data
    })
}

// 获取商品属性参数列表
export const getProductAttributeList = (params) => {
    return service({
        url: '/product/attribute',
        method: 'get',
        params
    })
}

// 创建商品属性参数列表
export const createProductAttribute = (data) => {
    return service({
      url: '/product/attribute',
      method: 'post',
      data
    })
}

// 更新商品属性参数列表
export const updateProductAttribute = (data) => {
    return service({
        url: '/product/attribute',
        method: 'put',
        data
    })
}

// 删除商品属性参数列表
export const deleteProductAttribute = (data) => {
    return service({
      url: '/product/attribute',
      method: 'delete',
      data
    })
}

// 获取商品分类列表
export const getProductCategoryList = (params) => {
    return service({
        url: '/product/productCategory',
        method: 'get',
        params
    })
}

// 获取商品分类列表
export const getProductAllCategory = (params) => {
    return service({
        url: '/product/allCategory',
        method: 'get',
        params
    })
}

// 创建商品分类列表
export const createProductCategory= (data) => {
    return service({
      url: '/product/productCategory',
      method: 'post',
      data
    })
}

// 更新商品分类列表
export const updateProductCategory = (data) => {
    return service({
        url: '/product/productCategory',
        method: 'put',
        data
    })
}

// 删除商品分类列表
export const deleteProductCategory = (data) => {
    return service({
      url: '/product/productCategory',
      method: 'delete',
      data
    })
}

// 获取商品sku库存
export const getProductSKUStockByProductID = (params) => {
    return service({
        url: '/product/sku',
        method: 'get',
        params
    })
}

// 更新商品sku库存
export const updateProductSKUStock = (data) => {
    return service({
        url: '/product/sku',
        method: 'put',
        data
    })
}

export function updateDeleteStatus(params) {
    return request({
      url:'/product/update/deleteStatus',
      method:'post',
      params:params
    })
  }

export function updateNewStatus(params) {
    return request({
      url:'/product/update/newStatus',
      method:'post',
      params:params
    })
  }

export function updateRecommendStatus(params) {
    return request({
      url:'/product/update/recommendStatus',
      method:'post',
      params:params
    })
}

export function updatePublishStatus(params) {
    return request({
      url:'/product/update/publishStatus',
      method:'post',
      params:params
    })
  }

// 更新商品属性参数值
export const updateProductAttributeValue = (data) => {
    return service({
        url: '/product/attributeValue',
        method: 'put',
        data
    })
}