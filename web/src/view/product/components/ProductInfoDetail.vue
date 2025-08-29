<template>
  <div style="margin-top: 50px">
    <el-form :model="modelValue" :rules="rules" ref="productInfoForm" label-width="120px" class="form-inner-container" size="small">
      <el-form-item label="商品分类：" prop="productCategoryId">
        <el-cascader
          v-model="selectProductCateValue"
          :options="productCateOptions"
          placeholder="请选择" clearable @change="handleProductTypeChange">
        </el-cascader>
      </el-form-item>
      <el-form-item label="商品名称：" prop="name">
        <el-input v-model="modelValue.name"></el-input>
      </el-form-item>
      <el-form-item label="副标题：" prop="subTitle">
        <el-input v-model="modelValue.subTitle"></el-input>
      </el-form-item>
      <el-form-item label="商品品牌：" prop="brandId">
        <el-select
          v-model="modelValue.brandId"
          @change="handleBrandChange"
          placeholder="请选择品牌">
          <el-option
            v-for="item in brandOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="商品介绍：">
        <el-input
          :autoSize="true"
          v-model="modelValue.description"
          type="textarea"
          placeholder="请输入内容"></el-input>
      </el-form-item>
      <el-form-item label="商品货号：">
        <el-input v-model="modelValue.productSN"></el-input>
      </el-form-item>
      <el-form-item label="商品售价：">
        <el-input-number v-model="modelValue.price" :precision="2"></el-input-number>
      </el-form-item>
      <el-form-item label="市场价：">
        <el-input-number v-model="modelValue.originalPrice" :precision="2"></el-input-number>
      </el-form-item>
      <el-form-item label="商品库存：">
        <el-input v-model.number="modelValue.stock"></el-input>
      </el-form-item>
      <el-form-item label="计量单位：">
        <el-input v-model="modelValue.unit"></el-input>
      </el-form-item>
      <el-form-item label="商品重量：">
        <el-input-number v-model="modelValue.weight" :precision="2" ></el-input-number>
        <span style="margin-left: 20px">克</span>
      </el-form-item>
      <el-form-item label="排序">
        <el-input v-model.number="modelValue.sort"></el-input>
      </el-form-item>
      <el-form-item style="text-align: center">
        <el-button type="primary" size="small" @click="handleNext('productInfoForm')">下一步，填写商品促销</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
// import {fetchListWithChildren} from '@/api/productCate'
// import {fetchList as fetchBrandList} from '@/api/brand'
import { ref, computed, onBeforeMount, watch, toRefs, getCurrentInstance } from 'vue'
import { ElMessage } from 'element-plus';
import { ProductStore } from '@/pinia/modules/product'
const productStore = ProductStore()

let prop = defineProps({
  modelValue: Object as () => ({
    id: number,
    productCategoryId: number,
    name: string,
    subTitle: string,
    brandId: number,
    description: string,
    productSN: string,
    price: number,
    originalPrice: number,
    stock: number,
    unit: string,
    weight: number,
    sort: number,
    brandName: string,
    productCategoryName: string,
  }),
  isEdit: Boolean
})
// 要响应式
const {modelValue, isEdit} = toRefs(prop);

const hasEditCreated = ref(false)
//选中商品分类的值
const selectProductCateValue = ref([])
const productCateOptions = ref([])
const brandOptions = ref([])
const rules = {
  name: [
    {required: true, message: '请输入商品名称', trigger: 'blur'},
    {min: 2, max: 140, message: '长度在 2 到 140 个字符', trigger: 'blur'}
  ],
  subTitle: [{required: true, message: '请输入商品副标题', trigger: 'blur'}],
  productCategoryId: [{required: true, message: '请选择商品分类', trigger: 'blur'}],
  brandId: [{required: true, message: '请选择商品品牌', trigger: 'blur'}],
  description: [{required: true, message: '请输入商品介绍', trigger: 'blur'}],
  requiredProp: [{required: true, message: '该项为必填项', trigger: 'blur'}]
}

onBeforeMount(() => {
  getProductCateList();
  getBrandList();
})

const productId = computed(() => modelValue.value.id)

// watch(() => selectProductCateValue.value, (newValue) => {
//   if (newValue != null && newValue.length === 2) {
//     modelValue.value.productCategoryId = newValue[1];
//     modelValue.value.productCategoryName = getCateNameById(modelValue.value.productCategoryId);
//   } else {
//     modelValue.value.productCategoryId = null;
//     modelValue.value.productCategoryName = null;
//   }
// })

watch(() => productId.value, (newValue) => {
  if(!isEdit.value)return;
  if(hasEditCreated.value)return;
  if(newValue===undefined||newValue==null||newValue===0)return;
  handleEditCreated();
})


const handleProductTypeChange = () => {
  modelValue.value.productCategoryId = selectProductCateValue.value.at(-1)
  modelValue.value.productCategoryName = getCateNameById(modelValue.value.productCategoryId)
}


const handleEditCreated = () => {
  // console.log('商品分类数据', productCateOptions.value);
  // console.log('编辑商品，设置商品分类', modelValue.value);
  if (modelValue.value.productCategoryId != null) {
    for (let i = 0; i < productCateOptions.value.length; i++) {
      if (productCateOptions.value[i].value === modelValue.value.productCategoryId) {
        selectProductCateValue.value.push(productCateOptions.value[i].value);
        break;
      } else {
        for (let j = 0; j < productCateOptions.value[i].children.length; j++) {
          if (productCateOptions.value[i].children[j].value === modelValue.value.productCategoryId) {
            selectProductCateValue.value.push(productCateOptions.value[i].value);
            selectProductCateValue.value.push(productCateOptions.value[i].children[j].value);
            break;
          }
        }
      }
    }
  }
  // console.log('选中商品分类', selectProductCateValue.value);
  // if(modelValue.value.productCategoryId != null){
  //   for (let index = 0; index < productCateOptions.value.length; index++) {
  //       const element = productCateOptions.value[index];
  //       if ("children" in element) {
  //           for (let index = 0; index < element.children.length; index++) {
  //               const item = element.children[index];
  //               if (item.value == modelValue.value.productCategoryId) {
  //                   selectProductCateValue.value.push(element.value)
  //                   selectProductCateValue.value.push(item.value)
  //               }
  //               break
  //           }
  //       } else {
  //           if ("value" in element && element.value == modelValue.value.productCategoryId) {
  //               selectProductCateValue.value.push(element.value)
  //               break
  //           }
  //       }
  //   }
  // }
  hasEditCreated.value = true;
}

const getProductCateList = async() => {
  await productStore.BuildProductCategoryData()
  productCateOptions.value = productStore.ProductCategoryOptions
}
const getBrandList = async() => {
  await productStore.BuildBrandData()
  let brandList = productStore.RandData['list']
  for (let i = 0; i < brandList.length; i++) {
    brandOptions.value.push({label: brandList[i].name, value: brandList[i].id});
  }
}
const getCateNameById = (id) => {
  let name=null;
  for(let i=0; i<productCateOptions.value.length; i++){
    for(let j=0; j<productCateOptions.value[i].children.length;j++){
      if(productCateOptions.value[i].children[j].value===id){
        name = productCateOptions.value[i].children[j].label;
        return name;
      }
    }
  }
  return name;
}

const emits = defineEmits(["nextStep"]);
const productInfoForm = ref(null)
const handleNext = (formName) => {
  productInfoForm.value.validate((valid) => {
    if (valid) {
      emits('nextStep');
    } else {
      ElMessage({
        message: '验证失败',
        type: 'error',
        duration:1000
      });
      return false;
    }
  });
}
const handleBrandChange = (val) => {
  let brandName = '';
  for (let i = 0; i < brandOptions.value.length; i++) {
    if (brandOptions.value[i].value === val) {
      brandName = brandOptions.value[i].label;
      break;
    }
  }
  modelValue.value.brandName = brandName;
}

</script>


<style scoped>
</style>
