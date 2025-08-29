<template> 
  <el-card class="form-container" shadow="never">
    <el-form :model="coupon"
             :rules="rules"
             ref="couponFrom"
             label-width="150px"
             size="small">
      <el-form-item label="优惠券类型：">
        <el-select v-model="coupon.type">
          <el-option
            v-for="type in typeOptions"
            :key="type.value"
            :label="type.label"
            :value="type.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="优惠券名称：" prop="name">
        <el-input v-model="coupon.name" class="input-width"></el-input>
      </el-form-item>
      <el-form-item label="适用平台：">
        <el-select v-model="coupon.platform">
          <el-option
            v-for="item in platformOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="总发行量：" prop="publishCount">
        <el-input v-model.number="coupon.publishCount" placeholder="只能输入正整数" class="input-width"></el-input>
      </el-form-item>
      <el-form-item label="面额：" prop="amount">
        <el-input v-model.number="coupon.amount" placeholder="面值只能是数值，限2位小数" class="input-width">
          <template slot="append">元</template>
        </el-input>
      </el-form-item>
      <el-form-item label="每人限领：">
        <el-input v-model="coupon.perLimit" placeholder="只能输入正整数" class="input-width">
          <template slot="append">张</template>
        </el-input>
      </el-form-item>
      <el-form-item label="使用门槛：" prop="minPoint">
        <el-input v-model.number="coupon.minPoint" placeholder="只能输入正整数" class="input-width">
          <template slot="prepend">满</template>
          <template slot="append">元可用</template>
        </el-input>
      </el-form-item>
      <el-form-item label="领取日期：" prop="enableTime">
        <el-date-picker type="date" placeholder="选择日期" v-model="coupon.enableTime" class="input-width"></el-date-picker>
      </el-form-item>
      <el-form-item label="有效期：">
        <el-date-picker type="date" placeholder="选择日期" v-model="coupon.startTime" style="width: 150px"></el-date-picker>
        <span style="margin-left: 20px;margin-right: 20px">至</span>
        <el-date-picker type="date" placeholder="选择日期" v-model="coupon.endTime" style="width: 150px"></el-date-picker>
      </el-form-item>
      <el-form-item label="可使用商品：">
        <el-radio-group v-model="coupon.useType">
          <el-radio-button :label="0">全场通用</el-radio-button>
          <el-radio-button :label="1">指定分类</el-radio-button>
          <el-radio-button :label="2">指定商品</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item v-show="coupon.useType===1">
        <el-cascader
          clearable
          placeholder="请选择分类名称"
          v-model="selectProductCate"
          :options="productCateOptions">
        </el-cascader>
        <el-button @click="handleAddProductCategoryRelation()">添加</el-button>
        <el-table ref="productCateRelationTable"
                  :data="coupon.productCategoryRelationList"
                  style="width: 100%;margin-top: 20px"
                  border>
          <el-table-column label="分类名称" align="center">
            <template #default="scope">{{scope.row.parentCategoryName}}>{{scope.row.productCategoryName}}</template>
          </el-table-column>
          <el-table-column label="操作" align="center" width="100">
            <template #default="scope">
              <el-button size="mini"
                         type="text"
                         @click="handleDeleteProductCateRelation(scope.$index, scope.row)">删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <el-form-item v-show="coupon.useType===2">
        <el-select
          v-model="selectProduct"
          filterable
          remote
          reserve-keyword
          placeholder="商品名称/商品货号"
          :remote-method="searchProductMethod"
          :loading="selectProductLoading">
          <el-option
            v-for="item in selectProductOptions"
            :key="item.productId"
            :label="item.productName"
            :value="item.productId">
            <span style="float: left">{{ item.productName }}</span>
            <span style="float: right; color: #8492a6; font-size: 13px">NO.{{ item.productSn }}</span>
          </el-option>
        </el-select>
        <el-button @click="handleAddProductRelation()">添加</el-button>
        <el-table ref="productRelationTable"
                  :data="coupon.productRelationList"
                  style="width: 100%;margin-top: 20px"
                  border>
          <el-table-column label="商品名称" align="center">
            <template #default="scope">{{scope.row.productName}}</template>
          </el-table-column>
          <el-table-column label="货号" align="center"  width="120" >
            <template #default="scope">NO.{{scope.row.productSn}}</template>
          </el-table-column>
          <el-table-column label="操作" align="center" width="100">
            <template #default="scope">
              <el-button size="mini"
                         type="text"
                         @click="handleDeleteProductRelation(scope.$index, scope.row)">删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <el-form-item label="备注：">
        <el-input
          class="input-width"
          type="textarea"
          :rows="5"
          placeholder="请输入内容"
          v-model="coupon.note">
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit('couponFrom')">提交</el-button>
        <el-button v-if="!isEdit" @click="resetForm('couponFrom')">重置</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>
<script setup>
import { ref, reactive, getCurrentInstance, onBeforeMount } from 'vue';
import { createCoupon, getCoupon, updateCoupon } from '@/api/coupon';
import { getSimpleList } from '@/api/product';
import { ProductStore } from '@/pinia/modules/product';
import { ElMessageBox } from 'element-plus';

// Props定义
const props = defineProps({
  isEdit: {
    type: Boolean,
    default: false
  }
});

// 获取组件实例以访问路由
const instance = getCurrentInstance();
const route = instance.proxy.$route;
const router = instance.proxy.$router;

// Pinia Store
const productStore = ProductStore();

// 常量定义
const defaultCoupon = {
  type: 0,
  name: null,
  platform: 0,
  amount: null,
  perLimit: 1,
  minPoint: null,
  startTime: null,
  endTime: null,
  useType: 0,
  note: null,
  publishCount: null,
  productRelationList: [],
  productCategoryRelationList: []
};

const defaultTypeOptions = [
  { label: '全场赠券', value: 0 },
  { label: '会员赠券', value: 1 },
  { label: '购物赠券', value: 2 },
  { label: '注册赠券', value: 3 }
];

const defaultPlatformOptions = [
  { label: '全平台', value: 0 },
  { label: '移动平台', value: 1 },
  { label: 'PC平台', value: 2 }
];

// 响应式数据
const coupon = reactive({ ...defaultCoupon });
const typeOptions = ref([...defaultTypeOptions]);
const platformOptions = ref([...defaultPlatformOptions]);
const rules = ref({
  name: [
    { required: true, message: '请输入优惠券名称', trigger: 'blur' },
    { min: 2, max: 140, message: '长度在 2 到 140 个字符', trigger: 'blur' }
  ],
  publishCount: [
    { type: 'number', required: true, message: '只能输入正整数', trigger: 'blur' }
  ],
  amount: [
    { type: 'number', required: true, message: '面值只能是数值，0.01-10000，限2位小数', trigger: 'blur' }
  ],
  minPoint: [
    { type: 'number', required: true, message: '只能输入正整数', trigger: 'blur' }
  ]
});

// 其他响应式变量
const couponFrom = ref(null);
const selectProduct = ref(null);
const selectProductLoading = ref(false);
const selectProductOptions = ref([]);
const selectProductCate = ref(null);
const productCateOptions = ref([]);

onBeforeMount(() => {
  if (props.isEdit) {
    let param = { id: route.query.id }
    getCoupon(param).then(response => {
      Object.assign(coupon, response.data);
    });
  }
  getProductCateList();
});

// 方法定义
const onSubmit = (formName) => {
  couponFrom.value.validate((valid) => {
    if (valid) {
      ElMessageBox.confirm('是否提交数据', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        if (props.isEdit) {
          updateCoupon(coupon).then(() => {
            couponFrom.value.resetFields();
            ElMessageBox.confirm({
              message: '修改成功',
              type: 'success',
              duration: 1000
            });
            router.back();
          });
        } else {
          createCoupon(coupon).then(() => {
            couponFrom.value.resetFields();
            ElMessageBox.confirm({
              message: '提交成功',
              type: 'success',
              duration: 1000
            });
            router.back();
          });
        }
      });
    } else {
      ElMessageBox.confirm({
        message: '验证失败',
        type: 'error',
        duration: 1000
      });
      return false;
    }
  });
};

const resetForm = (formName) => {
  couponFrom.value.resetFields();
  Object.assign(coupon, defaultCoupon);
};

const searchProductMethod = async(query) => {
  if (query !== '') {
    selectProductLoading.value = true;
    await getSimpleList({ keyword: query }).then(response => {

      selectProductLoading.value = false;
      const productList = response.data.list;
      console.log("----[searchProductMethod]------productList:", productList)
      selectProductOptions.value = productList.map(item => ({
        productId: item.id,
        productName: item.name,
        productSn: item.productSN
      }));
    });
    console.log("----[searchProductMethod]------selectProductOptions:", selectProductOptions.value)
  } else {
    // selectProductOptions.value = [];
  }
  console.log("----[searchProductMethod]------selectProductOptions:", selectProductOptions.value)
};

const handleAddProductRelation = () => {
  if (selectProduct.value === null) {
    ElMessageBox.confirm({
      message: '请先选择商品',
      type: 'warning'
    });
    return;
  }
  coupon.productRelationList.push(getProductById(selectProduct.value));
  selectProduct.value = null;
};

const handleDeleteProductRelation = (index) => {
  coupon.productRelationList.splice(index, 1);
};

const handleAddProductCategoryRelation = () => {
  if (selectProductCate.value === null || selectProductCate.value.length === 0) {
    ElMessageBox.confirm({
      message: '请先选择商品分类',
      type: 'warning'
    });
    return;
  }
  coupon.productCategoryRelationList.push(getProductCateByIds(selectProductCate.value));
  selectProductCate.value = [];
};

const handleDeleteProductCateRelation = (index) => {
  coupon.productCategoryRelationList.splice(index, 1);
};

const getProductById = (id) => {

  return selectProductOptions.value.find(item => item.productId === id) || null;
};

const getProductCateList = () => {
  productStore.BuildProductCategoryData();
  productCateOptions.value = productStore.ProductCategoryOptions;
};

const getProductCateByIds = (ids) => {
  let name, parentName;
  for (let i = 0; i < productCateOptions.value.length; i++) {
    if (productCateOptions.value[i].value === ids[0]) {
      parentName = productCateOptions.value[i].label;
      for (let j = 0; j < productCateOptions.value[i].children.length; j++) {
        if (productCateOptions.value[i].children[j].value === ids[1]) {
          name = productCateOptions.value[i].children[j].label;
        }
      }
    }
  }
  return { productCategoryId: ids[1], productCategoryName: name, parentCategoryName: parentName };
};
</script>
<style scoped>
  .input-width {
    width: 60%;
  }
</style>


