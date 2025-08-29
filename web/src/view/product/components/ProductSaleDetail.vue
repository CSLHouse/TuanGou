<template>
  <div style="margin-top: 50px">
    <el-form :model="modelValue" ref="productSaleForm" label-width="120px" class="form-inner-container" size="small">
      <el-form-item label="赠送积分：">
        <el-input v-model.number="modelValue.giftPoint"></el-input>
      </el-form-item>
      <el-form-item label="赠送成长值：">
        <el-input v-model.number="modelValue.giftGrowth"></el-input>
      </el-form-item>
      <el-form-item label="积分购买限制：">
        <el-input v-model.number="modelValue.usePointLimit"></el-input>
      </el-form-item>
      <el-form-item label="预告商品：">
        <el-switch
          v-model="modelValue.previewStatus"
          :active-value="1"
          :inactive-value="0">
        </el-switch>
      </el-form-item>
      <el-form-item label="商品上架：">
        <el-switch
          v-model="modelValue.publishStatus"
          :active-value="1"
          :inactive-value="0">
        </el-switch>
      </el-form-item>
      <el-form-item label="商品推荐：">
        <span style="margin-right: 10px">新品</span>
        <el-switch
          v-model="modelValue.newStatus"
          :active-value="1"
          :inactive-value="0">
        </el-switch>
        <span style="margin-left: 10px;margin-right: 10px">推荐</span>
        <el-switch
          v-model="modelValue.recommandStatus"
          :active-value="1"
          :inactive-value="0">
        </el-switch>
      </el-form-item>
      <el-form-item label="服务保证：">
        <el-checkbox-group v-model="selectServiceList">
          <el-checkbox :label="1">无忧退货</el-checkbox>
          <el-checkbox :label="2">快速退款</el-checkbox>
          <el-checkbox :label="3">免费包邮</el-checkbox>
        </el-checkbox-group>
      </el-form-item>
      <el-form-item label="详细页标题：">
        <el-input v-model="modelValue.detailTitle"></el-input>
      </el-form-item>
      <el-form-item label="详细页描述：">
        <el-input v-model="modelValue.detailDesc"></el-input>
      </el-form-item>
      <el-form-item label="商品关键字：">
        <el-input v-model="modelValue.keywords"></el-input>
      </el-form-item>
      <el-form-item label="商品备注：">
        <el-input v-model="modelValue.note" type="textarea" :autoSize="true"></el-input>
      </el-form-item>
      <el-form-item label="选择优惠方式：">
        <el-radio-group v-model="modelValue.promotionType" size="small">
          <el-radio-button :label="0">无优惠</el-radio-button>
          <el-radio-button :label="1">特惠促销</el-radio-button>
          <el-radio-button :label="2">会员价格</el-radio-button>
          <el-radio-button :label="3">阶梯价格</el-radio-button>
          <el-radio-button :label="4">满减价格</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item v-show="modelValue.promotionType===1">
        <div>
          开始时间：
          <el-date-picker
            v-model="modelValue.promotionStartTime"
            type="date"
            placeholder="选择开始时间">
          </el-date-picker>
        </div>
        <div class="littleMargin">
          结束时间：
          <el-date-picker
            v-model="modelValue.promotionEndTime"
            type="date"
            placeholder="选择结束时间">
          </el-date-picker>
        </div>
        <div class="littleMargin">
          促销价格：
          <el-input-number style="width: 220px" v-model="modelValue.promotionPrice" :precision="2" placeholder="输入促销价格"></el-input-number>
        </div>

      </el-form-item>
      <el-form-item v-show="modelValue.promotionType===2">
        <div v-for="(item, index) in modelValue.memberPriceList" :class="{littleMargin:index!==0}">
          {{item.memberLevelName}}：
          <el-input-number v-model="item.memberPrice" :precision="2" tyle="width: 200px"></el-input-number>
        </div>
      </el-form-item>
      <el-form-item v-show="modelValue.promotionType===3">
        <el-table :data="modelValue.productLadderList"
                  style="width: 80%" border>
          <el-table-column
            label="数量"
            align="center"
            width="120">
            <template #default="scope">
              <el-input-number v-model="scope.row.count" :precision="2"></el-input-number>
            </template>
          </el-table-column>
          <el-table-column
            label="折扣"
            align="center"
            width="120">
            <template #default="scope">
              <el-input-number v-model="scope.row.discount" :precision="2"></el-input-number>
            </template>
          </el-table-column>
          <el-table-column
            align="center"
            label="操作">
            <template #default="scope">
              <el-button type="danger" link @click="handleRemoveProductLadder(scope.$index, scope.row)">删除</el-button>
              <el-button type="primary" link @click="handleAddProductLadder(scope.$index, scope.row)">添加</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <el-form-item v-show="modelValue.promotionType===4">
        <el-table :data="modelValue.productFullReductionList"
                  style="width: 80%" border>
          <el-table-column
            label="满"
            align="center"
            width="120">
            <template #default="scope">
              <el-input-number v-model="scope.row.fullPrice" :precision="2"></el-input-number>
            </template>
          </el-table-column>
          <el-table-column
            label="立减"
            align="center"
            width="120">
            <template #default="scope">
              <el-input-number v-model="scope.row.reducePrice" :precision="2"></el-input-number>
            </template>
          </el-table-column>
          <el-table-column
            align="center"
            label="操作">
            <template #default="scope">
              <el-button type="danger" link @click="handleRemoveFullReduction(scope.$index, scope.row)">删除</el-button>
              <el-button type="primary" link @click="handleAddFullReduction(scope.$index, scope.row)">添加</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <el-form-item style="text-align: center">
        <el-button size="small" @click="handlePrev">上一步，填写商品信息</el-button>
        <el-button type="primary" size="small" @click="handleNext">下一步，填写商品属性</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
// import {fetchList as fetchMemberLevelList} from '@/api/memberLevel'
import { ref, computed, onBeforeMount, watch, toRefs, PropType } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

interface MemberPriceValue {
  memberLevelId: number, memberLevelName: string, memberPrice: number
}

interface ProductLadderValue {
  count: number, discount: number, price: number
}

interface ProductFullReductionValue {
  fullPrice: number, reducePrice: number
}

let prop = defineProps({
  modelValue: Object as () => ({
    id: number,
    memberPriceList: MemberPriceValue[],
    serviceIds: string,
    productLadderList: ProductLadderValue[],
    productFullReductionList: ProductFullReductionValue[],
    giftPoint: number,
    giftGrowth: number,
    usePointLimit: number,
    previewStatus: number,
    publishStatus: number,
    newStatus: number,
    recommandStatus: number,
    detailTitle: string,
    detailDesc: string,
    keywords: string,
    note: string,
    promotionType: number,
    promotionStartTime: string,
    promotionEndTime: string,
    promotionPrice: number,
    

  }),
  isEdit: Boolean
})
// 要响应式
const {modelValue, isEdit} = toRefs(prop);

const pickerOptions1 = ref({
  disabledDate(time) {
    return time.getTime() < Date.now();
  }
})

onBeforeMount(() => {
  if (isEdit.value) {
  } else {
    modelValue.value.memberPriceList = [
        {memberLevelId: 1, memberLevelName: "黄金会员", memberPrice: null},
        {memberLevelId: 2, memberLevelName: "白金会员", memberPrice: null},
        {memberLevelId: 3, memberLevelName: "钻石会员", memberPrice: null}
    ]
    // fetchMemberLevelList({defaultStatus: 0}).then(response => {
    //   let memberPriceList = [];
    //   for (let i = 0; i < response.data.length; i++) {
    //     let item = response.data[i];
    //     memberPriceList.push({memberLevelId: item.id, memberLevelName: item.name})
    //   }
    //   modelValue.value.memberPriceList = memberPriceList;
    // });
  }
})

const selectServiceList = computed({
  get() {
    let list = [];
    if (modelValue.value.serviceIds === undefined || modelValue.value.serviceIds == null || modelValue.value.serviceIds === '') return list;
    let ids = modelValue.value.serviceIds.split(',');
    for (let i = 0; i < ids.length; i++) {
      list.push(Number(ids[i]));
    }
    return list;
  },
  set(newValue) {
    let serviceIds = '';
    if (newValue != null && newValue.length > 0) {
      for (let i = 0; i < newValue.length; i++) {
        serviceIds += newValue[i] + ',';
      }
      if (serviceIds.endsWith(',')) {
        serviceIds = serviceIds.substr(0, serviceIds.length - 1)
      }
      modelValue.value.serviceIds = serviceIds;
    } else {
      modelValue.value.serviceIds = null;
    }
  }
})

const handleEditCreated = () => {
  let ids = modelValue.value.serviceIds.split(',');
  console.log('handleEditCreated', ids);
  for (let i = 0; i < ids.length; i++) {
    selectServiceList.value.push(Number(ids[i]));
  }
}
const handleRemoveProductLadder = (index, row) =>  {
  let productLadderList = modelValue.value.productLadderList;
  if (productLadderList.length === 1) {
    productLadderList.pop();
    productLadderList.push({
      count: 0,
      discount: 0,
      price: 0
    })
  } else {
    productLadderList.splice(index, 1);
  }
}
const handleAddProductLadder = (index, row) => {
  let productLadderList = modelValue.value.productLadderList;
  if (productLadderList.length < 3) {
    productLadderList.push({
      count: 0,
      discount: 0,
      price: 0
    })
  } else {
    ElMessage({
      message: '最多只能添加三条',
      type: 'warning'
    })
  }
}
const handleRemoveFullReduction = (index, row) => {
  let fullReductionList = modelValue.value.productFullReductionList;
  if (fullReductionList.length === 1) {
    fullReductionList.pop();
    fullReductionList.push({
      fullPrice: 0,
      reducePrice: 0
    });
  } else {
    fullReductionList.splice(index, 1);
  }
}
const handleAddFullReduction = (index, row) => {
  let fullReductionList = modelValue.value.productFullReductionList;
  if (fullReductionList.length < 3) {
    fullReductionList.push({
      fullPrice: 0,
      reducePrice: 0
    });
  } else {
    ElMessage({
      message: '最多只能添加三条',
      type: 'warning'
    })
  }
}

const emits = defineEmits(["prevStep", "nextStep"]);
const handlePrev = () => {
  emits('prevStep')
}
const handleNext = () => {
  emits('nextStep')
}

</script>


<style scoped>
  .littleMargin {
    margin-top: 10px;
  }
</style>
