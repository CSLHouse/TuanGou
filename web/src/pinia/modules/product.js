import { getProductAttributeCategoryList, getBrandList, getProductAllCategory } from '@/api/product'
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

export const ProductStore = defineStore('product', () => {

  const ProductAttributeCategoryList = ref()
  const RandData = ref()
  const ProductAttributeCategoryOptions = ref([])
  const ProductCategoryList = ref()
  const ProductCategoryOptions = ref([])

  const setProductAttributeCategoryList = (val) => {
    ProductAttributeCategoryList.value = val
  }

  /* 获取产品属性分类*/
  const GetProductAttributeCategoryList = async() => {
    const res = await getProductAttributeCategoryList()
    if (res.code === 0) {
      ProductAttributeCategoryList.value = []
      setProductAttributeCategoryList(res.data.list)
    }
    return res
  }
  
  const setBrandData = (val) => {
    RandData.value = val
  }

  /* 获取产品品牌分类*/
  const BuildBrandData = async(isRefresh) => {
    if (!RandData.value || isRefresh) {
      const res = await getBrandList()
      if (res.code === 0) {
        setBrandData(res.data)
      }
    }
  }
  
  const parseProductAttributeCategory = () => {
    let productAttributeMap = {}
    ProductAttributeCategoryOptions.value = []
    // console.log("--[parseProductAttributeCategory]ProductAttributeCategoryList:", ProductAttributeCategoryList.value)
    ProductAttributeCategoryList.value.forEach((item) => {
        let splitted = item.name.split("-")
        if ( !productAttributeMap[splitted[0]]) {
            if (splitted.length > 1) {
                productAttributeMap[splitted[0]] = [{"id": item.id, "data": splitted[1]}]
            } else {
                productAttributeMap[splitted[0]] = {"id": item.id, "data": splitted[0]}
            }
        } else {
            if (splitted.length > 1) {
                productAttributeMap[splitted[0]].push({"id": item.id, "data": splitted[1]})
            }
        }
    })
    // console.log("---productAttributeMap-", productAttributeMap)
    let count = 0
    for (let key in productAttributeMap) {
        let productAttribute = {}
        productAttribute["label"] = key
        if (Array.isArray(productAttributeMap[key])) {
          productAttribute["value"] = count
          productAttribute["children"] = []
          productAttributeMap[key].forEach((item) => {
            let productAttributeItem = {}
            productAttributeItem["label"] = item.data
            productAttributeItem["value"] = item.id
            productAttribute["children"].push(productAttributeItem)
          })
            count += 1
        } else {
          productAttribute["value"] = productAttributeMap[key].id
        }
        ProductAttributeCategoryOptions.value.push(productAttribute)
    }
    // console.log("--[parseProductAttributeCategory]ProductAttributeCategoryOptions:", ProductAttributeCategoryOptions.value)
  }
  const BuildProductAttributeData = async(isRefresh) => {
    if (!ProductAttributeCategoryList.value || isRefresh) {
      await GetProductAttributeCategoryList()
    }
    parseProductAttributeCategory()
  }

  /* 获取产品分类*/
  const GetProductCategoryList = async() => {
    const res = await getProductAllCategory()
    if (res.code === 0) {
      ProductCategoryList.value = []
      setProductCategoryList(res.data)
    }
    return res
  }
  const setProductCategoryList = (val) => {
    ProductCategoryList.value = val
  }
  const parseProductCategory = () => {
    ProductCategoryOptions.value = []
    ProductCategoryList.value.forEach((element) => {
      if (element.parentId == 0) {
        let category = {label: element.name, value: element.id, children: []}
        ProductCategoryList.value.forEach((item) => {
          if (item.parentId === element.id) {
            let children = {label: item.name, value: item.id}
            category.children.push(children)
          }
        })
        ProductCategoryOptions.value.push(category)
      }
    })
  }
  const BuildProductCategoryData = async(isRefresh) => {
    if (!ProductCategoryList.value || isRefresh) {
      await GetProductCategoryList()
    }
    parseProductCategory()
  }
  
  return {
    ProductAttributeCategoryList,
    GetProductAttributeCategoryList,
    RandData,
    BuildBrandData,
    ProductAttributeCategoryOptions,
    BuildProductAttributeData,
    ProductCategoryOptions,
    BuildProductCategoryData,
  }
})
