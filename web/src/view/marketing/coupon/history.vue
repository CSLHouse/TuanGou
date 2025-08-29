<template>
  <div class="app-container">
    <div class="table-layout">
      <el-row>
        <el-col :span="4" class="table-cell-title">名称</el-col>
        <el-col :span="4" class="table-cell-title">优惠券类型</el-col>
        <el-col :span="4" class="table-cell-title">可使用商品</el-col>
        <el-col :span="4" class="table-cell-title">使用门槛</el-col>
        <el-col :span="4" class="table-cell-title">面值</el-col>
        <el-col :span="4" class="table-cell-title">状态</el-col>
      </el-row>
      <el-row>
        <el-col :span="4" class="table-cell">{{ coupon.name }}</el-col>
        <el-col :span="4" class="table-cell">{{ formatType(coupon.type) }}</el-col>
        <el-col :span="4" class="table-cell">{{ formatUseType(coupon.useType) }}</el-col>
        <el-col :span="4" class="table-cell">满{{ coupon.minPoint }}元可用</el-col>
        <el-col :span="4" class="table-cell">{{ coupon.amount }}元</el-col>
        <el-col :span="4" class="table-cell">{{ formatStatus(coupon.endTime) }}</el-col>
      </el-row>
      <el-row>
        <el-col :span="4" class="table-cell-title">有效期</el-col>
        <el-col :span="4" class="table-cell-title">总发行量</el-col>
        <el-col :span="4" class="table-cell-title">已领取</el-col>
        <el-col :span="4" class="table-cell-title">待领取</el-col>
        <el-col :span="4" class="table-cell-title">已使用</el-col>
        <el-col :span="4" class="table-cell-title">未使用</el-col>
      </el-row>
      <el-row>
        <el-col :span="4" class="table-cell" style="font-size: 13px">
          {{ formatDateFilter(coupon.startTime) }}至{{ formatDateFilter(coupon.endTime) }}
        </el-col>
        <el-col :span="4" class="table-cell">{{ coupon.publishCount }}</el-col>
        <el-col :span="4" class="table-cell">{{ coupon.receiveCount }}</el-col>
        <el-col :span="4" class="table-cell">{{ coupon.publishCount - coupon.receiveCount }}</el-col>
        <el-col :span="4" class="table-cell">{{ coupon.useCount }}</el-col>
        <el-col :span="4" class="table-cell">{{ coupon.publishCount - coupon.useCount }}</el-col>
      </el-row>
    </div>
    <el-card class="filter-container" shadow="never">
      <div>
        <i class="el-icon-search"></i>
        <span>筛选搜索</span>
        <el-button
          style="float:right"
          type="primary"
          @click="handleSearchList()"
          size="small">
          查询搜索
        </el-button>
        <el-button
          style="float:right;margin-right: 15px"
          @click="handleResetSearch()"
          size="small">
          重置
        </el-button>
      </div>
      <div style="margin-top: 15px">
        <el-form :inline="true" :model="listQuery" size="small" label-width="140px">
          <el-form-item label="使用状态：">
            <el-select v-model="listQuery.useStatus" placeholder="全部" clearable class="input-width">
              <el-option v-for="item in useTypeOptions"
                         :key="item.value"
                         :label="item.label"
                         :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="订单编号：">
            <el-input v-model="listQuery.orderSn" class="input-width" placeholder="订单编号"></el-input>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
    <div class="table-container">
      <el-table ref="couponHistoryTable"
                :data="list"
                style="width: 100%;"
                v-loading="listLoading" border>
        <el-table-column label="优惠码" width="160" align="center">
          <template #default="scope">{{ scope.row.couponCode }}</template>
        </el-table-column>
        <el-table-column label="领取会员" width="140" align="center">
          <template #default="scope">{{ scope.row.memberNickname }}</template>
        </el-table-column>
        <el-table-column label="领取方式" width="100" align="center">
          <template #default="scope">{{ formatGetType(scope.row.getType) }}</template>
        </el-table-column>
        <el-table-column label="领取时间" width="160" align="center">
          <template #default="scope">{{ formatTime(scope.row.createTime) }}</template>
        </el-table-column>
        <el-table-column label="当前状态" width="140" align="center">
          <template #default="scope">{{ formatCouponHistoryUseType(scope.row.useStatus) }}</template>
        </el-table-column>
        <el-table-column label="使用时间" width="160" align="center">
          <template #default="scope">{{ formatTime(scope.row.useTime) }}</template>
        </el-table-column>
        <el-table-column label="订单编号" align="center">
          <template #default="scope">{{ scope.row.orderSn === null ? 'N/A' : scope.row.orderSn }}</template>
        </el-table-column>
      </el-table>
    </div>
    <div class="pagination-container">
      <el-pagination
        background
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        layout="total, sizes,prev, pager, next,jumper"
        :current-page.sync="listQuery.page"
        :page-size="listQuery.pageSize"
        :page-sizes="[5,10,15]"
        :total="total">
      </el-pagination>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, getCurrentInstance, onMounted } from 'vue';
import { formatDate } from '@/utils/date';
import { getCoupon, fetchCouponHistoryList } from '@/api/coupon';

// 获取组件实例，用于访问Vue 2的路由
const instance = getCurrentInstance();
const $route = instance.proxy.$route;

// 常量定义
const defaultTypeOptions = [
  { label: '全场赠券', value: 0 },
  { label: '会员赠券', value: 1 },
  { label: '购物赠券', value: 2 },
  { label: '注册赠券', value: 3 }
];

const defaultListQuery = {
  page: 1,
  pageSize: 10,
  useStatus: null,
  orderSn: null,
  couponId: null
};

const defaultUseTypeOptions = [
  { label: "未使用", value: 0 },
  { label: "已使用", value: 1 },
  { label: "已过期", value: 2 }
];

// 响应式变量
const coupon = reactive({});
const listQuery = reactive({ ...defaultListQuery });
const useTypeOptions = defaultUseTypeOptions;
const list = ref(null);
const total = ref(null);
const listLoading = ref(false);

// 初始化逻辑
onMounted(() => {
  // 使用Vue 2的路由获取方式
  let couponId = $route.query.id;
  if (!couponId) {
    ElMessage.error("未获取到优惠券ID");
    return;
  }
  let param = {id: couponId};
  getCoupon(param).then(response => {
    Object.assign(coupon, response.data);
  });

  listQuery.couponId = couponId;
  getList();
});

// 方法定义
const getList = () => {
  listLoading.value = true;
  fetchCouponHistoryList(listQuery).then(response => {
    listLoading.value = false;
    list.value = response.data.list;
    total.value = response.data.total;
  });
};

// 过滤器转换为普通函数
const formatType = (type) => {
  for (let i = 0; i < defaultTypeOptions.length; i++) {
    if (type === defaultTypeOptions[i].value) {
      return defaultTypeOptions[i].label;
    }
  }
  return '';
};

const formatUseType = (useType) => {
  if (useType === 0) {
    return '全场通用';
  } else if (useType === 1) {
    return '指定分类';
  } else {
    return '指定商品';
  }
};

const formatPlatform = (platform) => {
  if (platform === 1) {
    return '移动平台';
  } else if (platform === 2) {
    return 'PC平台';
  } else {
    return '全平台';
  }
};

const formatDateFilter = (time) => {
  if (time == null || time === '') {
    return 'N/A';
  }
  let date = new Date(time);
  return formatDate(date, 'yyyy-MM-dd');
};

const formatStatus = (endTime) => {
  let now = new Date().getTime();
  if (endTime > now) {
    return '未过期';
  } else {
    return '已过期';
  }
};

const formatGetType = (type) => {
  if (type === 1) {
    return '主动获取';
  } else {
    return '后台赠送';
  }
};

const formatCouponHistoryUseType = (useType) => {
  if (useType === 0) {
    return '未使用';
  } else if (useType === 1) {
    return '已使用';
  } else {
    return '已过期';
  }
};

const formatTime = (time) => {
  if (time == null || time === '') {
    return 'N/A';
  }
  let date = new Date(time);
  return formatDate(date, 'yyyy-MM-dd hh:mm:ss');
};

// 事件处理方法
const handleResetSearch = () => {
  Object.assign(listQuery, defaultListQuery);
  listQuery.couponId = $route.query.id;
};

const handleSearchList = () => {
  listQuery.page = 1;
  getList();
};

const handleSizeChange = (val) => {
  listQuery.page = 1;
  listQuery.pageSize = val;
  getList();
};

const handleCurrentChange = (val) => {
  listQuery.page = val;
  getList();
};
</script>

<style scoped>
/* 样式部分保持不变 */
.app-container {
  width: 80%;
  margin: 20px auto;
}

.filter-container {
  margin-top: 20px;
}

.table-layout {
  margin-top: 20px;
  border-left: 1px solid #DCDFE6;
  border-top: 1px solid #DCDFE6;
}

.table-cell {
  height: 60px;
  line-height: 40px;
  border-right: 1px solid #DCDFE6;
  border-bottom: 1px solid #DCDFE6;
  padding: 10px;
  font-size: 14px;
  color: #606266;
  text-align: center;
  overflow: hidden;
}

.table-cell-title {
  border-right: 1px solid #DCDFE6;
  border-bottom: 1px solid #DCDFE6;
  padding: 10px;
  background: #F2F6FC;
  text-align: center;
  font-size: 14px;
  color: #303133;
}
</style>
