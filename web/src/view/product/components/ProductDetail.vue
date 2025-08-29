<template>
  <el-card class="form-container" shadow="never">
    <el-steps :active="active" finish-status="success" align-center>
      <el-step title="填写商品信息"></el-step>
      <el-step title="填写商品促销"></el-step>
      <el-step title="填写商品属性"></el-step>
      <el-step title="选择商品关联"></el-step>
    </el-steps>
    <product-info-detail
      v-show="showStatus[0]"
      v-model="productParam"
      :is-edit="isEdit"
      @nextStep="nextStep">
    </product-info-detail>
    <product-sale-detail
      v-show="showStatus[1]"
      v-model="productParam"
      :is-edit="isEdit"
      @nextStep="nextStep"
      @prevStep="prevStep">
    </product-sale-detail>
    <product-attr-detail
      v-show="showStatus[2]"
      v-model="productParam"
      :is-edit="isEdit"
      @nextStep="nextStep"
      @prevStep="prevStep">
    </product-attr-detail>
    <product-relation-detail
      v-show="showStatus[3]"
      v-model="productParam"
      :is-edit="isEdit"
      @prevStep="prevStep"
      @finishCommit="finishCommit">
    </product-relation-detail>
  </el-card>
</template>

<script setup lang="ts">
  import ProductInfoDetail from './ProductInfoDetail.vue';
  import ProductSaleDetail from './ProductSaleDetail.vue';
  import ProductAttrDetail from './ProductAttrDetail.vue';
  import ProductRelationDetail from './ProductRelationDetail.vue';
  import {createProduct,getProductDetail,updateProduct} from '@/api/product';
  import { ref, computed, onBeforeMount, watch, toRefs, PropType } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useRouter } from 'vue-router'
  const router = useRouter()

  const defaultProductParam = {
    id: null,
    albumPics: '',
    brandId: null,
    brandName: '',
    deleteStatus: 0,
    description: '商品介绍',
    detailDesc: '详细页描述',
    detailHTML: '',
    detailMobileHTML: '',
    detailTitle: '详细页标题',
    feightTemplateId: 0,
    flashPromotionCount: 0,
    flashPromotionId: 0,
    flashPromotionPrice: 0,
    flashPromotionSort: 0,
    giftPoint: 100,
    giftGrowth: 100,
    keywords: '木耳袜',
    lowStock: 0,
    name: '木耳袜',
    newStatus: 1,
    note: '商品备注',
    originalPrice: 1,
    pic: '',
    //会员价格{memberLevelId: 0,memberPrice: 0,memberLevelName: null}
    memberPriceList: [],
    //商品满减
    productFullReductionList: [{fullPrice: 0, reducePrice: 0}],
    //商品阶梯价格
    productLadderList: [{count: 0,discount: 0,price: 0}],
    previewStatus: 1,
    price: 0.02,
    productAttributeCategoryId: null,
    //商品属性相关{productAttributeId: 0, value: ''}
    productAttributeValueList: [],
    //商品sku库存信息{lowStock: 0, pic: '', price: 0, sale: 0, skuCode: '', spData: '', stock: 0}
    skuStockList: [],
    //商品相关专题{subjectId: 0}
    subjectProductRelationList: [],
    //商品相关优选{prefrenceAreaId: 0}
    prefrenceAreaProductRelationList: [],
    productCategoryId: null,
    productCategoryName: '',
    productSN: 'No.0001',
    promotionEndTime: null,
    promotionPerLimit: 0,
    promotionPrice: null,
    promotionStartTime: null,
    promotionType: 0,
    publishStatus: 1,
    recommandStatus: 1,
    sale: 0,
    serviceIds: '',
    sort: 0,
    stock: 100,
    subTitle: '副标题',
    unit: '双',
    usePointLimit: 0,
    verifyStatus: 0,
    weight: 10
  }

let prop = defineProps({
  isEdit: Boolean
})
// 要响应式
const {isEdit} = toRefs(prop);

const active = ref(0)
const productParam = ref(defaultProductParam)
const showStatus = [true, false, false, false]

onBeforeMount(() => {
  if(isEdit.value){
    console.log('编辑商品', router.currentRoute.value.query.id);
    getProductDetail({id: router.currentRoute.value.query.id}).then(response=>{
      productParam.value = response.data;
      if ("code" in response && response.code === 0) {
        productParam.value =  response.data.product;
        // let serviceList = productParam.value.serviceIds.split(',')
        // checkList.value = []
        // for (let i = 0; i < serviceList.length; i++) {
        //   let element = serviceList[i].trim()
        //   if (element) {
        //       checkList.value.push(Number(element))
        //   }
        // }
        // promotionTypeState.value = productForm.value.promotionType

        // // 初始化商品规格 库存数据
        // skuTableData.value = productForm.value.skuStockList
        // for (let i = 0; i < skuTableData.value.length; i++) {
        //   let element = skuTableData.value[i]
        //   let spData = JSON.parse(element.spData)
        //   for (let j = 0; j < spData.length; j++) {
        //       let spItem = spData[j]
        //       element[spItem.key] = spItem.value
        //   }
        // }
      }
    });
  }
})

const  hideAll = () => {
  for (let i = 0; i < showStatus.length; i++) {
    showStatus[i] = false;
  }
}
const prevStep = () => {
  if (active.value > 0 && active.value < showStatus.length) {
    active.value--;
    hideAll();
    showStatus[active.value] = true;
  }
}
const nextStep = () => {
  if (active.value < showStatus.length - 1) {
    active.value++;
    hideAll();
    showStatus[active.value] = true;
  }
}
const finishCommit = (isEdit) => {
  ElMessageBox.confirm('是否要提交该产品', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    if (productParam.value.promotionEndTime && productParam.value.promotionEndTime.length === 0) {
      productParam.value.promotionEndTime = null;
    }
    if (productParam.value.promotionStartTime && productParam.value.promotionStartTime.length === 0) {
      productParam.value.promotionStartTime = null;
    }
    if(isEdit){
      updateProduct(productParam.value).then(response=>{
        ElMessage({
          message: '提交成功',
          type: 'success',
          duration: 1000
        })
        router.back();
      });
    }else{
      createProduct(productParam.value).then(response=>{
        if ("code" in response && response.code !== 0) {
          ElMessage({
            message: '提交失败',
            type: 'error',
            duration: 1000
          })
          return;
        }
      
        location.reload();
      });
    }
  });
}
</script>


<style>
  .form-container {
    width: 960px;
  }
  .form-inner-container {
    width: 800px;
  }
</style>


