<template>
  <div class="app-container">
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
          <el-form-item label="优惠券名称：">
            <el-input v-model="listQuery.name" class="input-width" placeholder="优惠券名称"></el-input>
          </el-form-item>
          <el-form-item label="优惠券类型：">
            <el-select v-model="listQuery.type" placeholder="全部" clearable class="input-width">
              <el-option v-for="item in typeOptions"
                         :key="item.value"
                         :label="item.label"
                         :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
    <el-card class="operate-container" shadow="never">
      <i class="el-icon-tickets"></i>
      <span>数据列表</span>
      <el-button size="mini" class="btn-add" @click="handleAdd()">添加</el-button>
    </el-card>
    <div class="table-container">
      <el-table ref="couponTable"
                :data="list"
                style="width: 100%;"
                @selection-change="handleSelectionChange"
                v-loading="listLoading" border>
        <el-table-column type="selection" width="60" align="center"></el-table-column>
        <el-table-column label="编号" width="80" align="center">
          <template #default="scope">{{scope.row.id}}</template>
        </el-table-column>
        <el-table-column label="优惠劵名称" width="120" align="center">
          <template #default="scope">{{scope.row.name}}</template>
        </el-table-column>
        <el-table-column label="优惠券类型" width="100" align="center">
          <template #default="scope">{{formatType(scope.row.type)}}</template>
        </el-table-column>
        <el-table-column label="可使用商品" width="100" align="center">
          <template #default="scope">{{formatUseType(scope.row.useType)}}</template>
        </el-table-column>
        <el-table-column label="使用门槛" width="120" align="center">
          <template #default="scope">满{{scope.row.minPoint}}元可用</template>
        </el-table-column>
        <el-table-column label="面值" width="100" align="center">
          <template #default="scope">{{scope.row.amount}}元</template>
        </el-table-column>
        <el-table-column label="适用平台" width="100" align="center">
          <template #default="scope">{{formatPlatform(scope.row.platform)}}</template>
        </el-table-column>
        <el-table-column label="有效期" width="190" align="center">
          <template #default="scope">{{formatDate(scope.row.startTime)}}至{{formatDate(scope.row.endTime)}}</template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="scope">{{formatStatus(scope.row.endTime)}}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center">
          <template #default="scope">
            <el-button size="mini"
                       type="text"
                       @click="handleView(scope.$index, scope.row)">查看</el-button>
            <el-button size="mini"
                       type="text"
                       @click="handleUpdate(scope.$index, scope.row)">
              编辑</el-button>
            <el-button size="mini"
                       type="text"
                       @click="handleDelete(scope.$index, scope.row)">删除</el-button>
          </template>
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
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { fetchList, deleteCoupon } from '@/api/coupon';
import { formatDate as formatDateUtil } from '@/utils/date';
import { ElMessage, ElMessageBox } from 'element-plus';

// 定义默认查询参数
const defaultListQuery = {
  page: 1,
  pageSize: 10,
  name: null,
  type: null
};

// 优惠券类型选项
const typeOptions = [
  { label: '全场赠券', value: 0 },
  { label: '会员赠券', value: 1 },
  { label: '购物赠券', value: 2 },
  { label: '注册赠券', value: 3 }
];

// 响应式数据
const listQuery = ref({ ...defaultListQuery });
const list = ref(null);
const total = ref(null);
const listLoading = ref(false);
const multipleSelection = ref([]);
const couponTable = ref(null);
const router = useRouter();

// 生命周期 - 组件挂载时获取列表
onMounted(() => {
  getList();
});

// 过滤器转换为普通函数
const formatType = (type) => {
  for (let i = 0; i < typeOptions.length; i++) {
    if (type === typeOptions[i].value) {
      return typeOptions[i].label;
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

const formatDate = (time) => {
  if (time == null || time === '') {
    return 'N/A';
  }
  let date = new Date(time);
  return formatDateUtil(date, 'yyyy-MM-dd');
};

const formatStatus = (endTime) => {
  let now = new Date().getTime();
  let endDate = new Date(endTime);
  return endDate > now ? '未过期' : '已过期';
};

// 方法
const handleResetSearch = () => {
  listQuery.value = { ...defaultListQuery };
};

const handleSearchList = () => {
  listQuery.value.page = 1;
  getList();
};

const handleSelectionChange = (val) => {
  multipleSelection.value = val;
};

const handleSizeChange = (val) => {
  listQuery.value.page = 1;
  listQuery.value.pageSize = val;
  getList();
};

const handleCurrentChange = (val) => {
  listQuery.value.page = val;
  getList();
};

const handleAdd = () => {
  router.push({ path: '/layout/marketing/addCoupon' });
};

const handleView = (index, row) => {
  router.push({ path: '/layout/marketing/couponHistory', query: { id: row.id } });
};

const handleUpdate = (index, row) => {
  router.push({ path: '/layout/marketing/updateCoupon', query: { id: row.id } });
};

const handleDelete = (index, row) => {
  ElMessageBox.confirm('是否进行删除操作?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    let params = { id: row.id };
    deleteCoupon(params).then(response => {
      ElMessageBox.confirm({
        type: 'success',
        message: '删除成功!'
      });
      getList();
    });
  });
};

const getList = () => {
  listLoading.value = true;
  fetchList(listQuery.value).then(response => {
    listLoading.value = false;
    list.value = response.data.list;
    total.value = response.data.total;
  });
};
</script>

<style scoped>
.input-width {
  width: 203px;
}
</style>