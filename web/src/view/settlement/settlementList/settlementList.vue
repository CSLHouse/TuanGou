<template>
  <div class="app-container">
    <el-card class="filter-container" shadow="never">
      <div>
        <!-- <i class="el-icon-search"></i> -->
        <span>筛选搜索</span>
        <el-button
          style="float: right"
          @click="handleSearchList()"
          type="primary"
          size="small">
          查询结果
        </el-button>
        <el-button
          style="float: right;margin-right: 15px"
          @click="handleResetSearch()"
          size="small">
          重置
        </el-button>
      </div>
      <div style="margin-top: 15px">
        <el-form :inline="true" :model="listQuery" size="small" label-width="140px">
          <el-form-item label="输入搜索：">
            <el-input style="width: 203px" v-model="listQuery.userId" placeholder="用户ID"></el-input>
          </el-form-item>
          <el-form-item label="结算单号：">
            <el-input style="width: 203px" v-model="listQuery.settlementNo" placeholder="结算单号"></el-input>
          </el-form-item>
          <el-form-item label="用户昵称：">
            <el-input style="width: 203px" v-model="listQuery.userName" placeholder="用户昵称"></el-input>
          </el-form-item>
          <el-form-item label="创建时间：">
              <el-date-picker
                class="input-width"
                v-model="listQuery.CreatedAt"
                value-format="YYYY-MM-DD"
                type="date"
                placeholder="请选择时间">
              </el-date-picker>
            </el-form-item>
          <el-form-item label="结算时间：">
              <el-date-picker
                class="input-width"
                v-model="listQuery.settlementTime"
                value-format="YYYY-MM-DD"
                type="date"
                placeholder="请选择时间">
              </el-date-picker>
            </el-form-item>
          <el-form-item label="结算状态：">
            <el-select v-model="listQuery.status" placeholder="全部" clearable>
              <el-option
                v-for="item in statusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
          
        </el-form>
      </div>
    </el-card>
    <div class="table-container">
      <el-table ref="productTable"
                :data="productList"
                style="width: 100%"
                @selection-change="handleSelectionChange"
                v-loading="listLoading"
                :default-sort="{ prop: 'id', order: 'descending' }"
                border>
        <!-- <el-table-column type="selection" width="60" align="center"></el-table-column> -->
        <!-- <el-table-column label="编号" width="80" align="center">
          <template #default="scope">{{scope.row.id}}</template>
        </el-table-column> -->
        <el-table-column prop="userId" label="用户ID" sortable  width="100" align="center">
          <template #default="scope">{{scope.row.userId}}</template>
        </el-table-column>
        <el-table-column label="用户昵称" min-width="120" align="center">
          <template #default="scope">
            <p>{{scope.row.userName}}</p>
          </template>
        </el-table-column>
        <el-table-column label="结算单号" min-width="200" align="center">
          <template #default="scope">
            <p>{{scope.row.settlementNo}}</p>
          </template>
        </el-table-column>
        <el-table-column label="结算金额" width="100" align="center">
          <template #default="scope">
            <p>{{scope.row.totalAmount}}</p>
          </template>
        </el-table-column>  
        <el-table-column prop="CreatedAt" label="账单生成时间" sortable  width="160" align="center">
          <template #default="scope">
            <p>{{formatDate(scope.row.CreatedAt)}}</p>
          </template>
        </el-table-column>  
       <el-table-column prop="settlementTime" label="结算时间" sortable  width="160" align="center">
          <template #default="scope">
            <p>{{formatDate(scope.row.settlementTime)}}</p>
          </template>
        </el-table-column> 
        <el-table-column prop="status" label="结算状态" sortable width="120" align="center">
          <template #default="scope">
            <p v-if="scope.row.status === 0" style="color: #ff4949;">{{verifyStatusFilter(scope.row.status)}}</p>
            <p v-else>{{ verifyStatusFilter(scope.row.status)}}</p>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" width="100">
            <template #default="scope">
              <el-button size="small"
                         type="text"
                         @click="handleSettlement(scope.$index, scope.row)">结算
              </el-button>
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
        :page-size="listQuery.pageSize"
        :page-sizes="[5,10,15]"
        :current-page.sync="listQuery.page"
        :total.number="+total">
      </el-pagination>
    </div>
  </div>
</template>
<script setup lang="ts">
import {
  getSettlementList,
  updateSettlementState
} from '@/api/team'
import { ref, computed, onBeforeMount, watch } from 'vue'
import { ElMessage, ElMessageBox, ElTable } from 'element-plus'
import { Edit, Search } from '@element-plus/icons-vue'
import { useRouter } from "vue-router";
import { formatDate as formatDateUtil } from '@/utils/date';
const router = useRouter()

const defaultListQuery = {
  page: 1,
  pageSize: 5,
  CreatedAt: null,
  userId: null,
  userName: null,
  settlementNo: null,
  settlementTime: null,
  status: null,
};
const listQuery = ref(defaultListQuery)

onBeforeMount(() => {
  getList()
})

const handleSearchList = async() => {
  listQuery.value.page = 1;
  getList();
}
const handleResetSearch = async() => {
  listQuery.value = defaultListQuery
}

// const page = ref(1)
const total = ref(0)
// const pageSize = ref(5)
const productList = ref([])
const listLoading = ref(true)

// 分页
const handleSizeChange = (val) => {
  listQuery.value.pageSize = val
  getList()
}

const handleCurrentChange = (val) => {
  console.log("handleCurrentChange", val)
  listQuery.value.page = val
  getList()
}

// 查询
const getList = async() => {
  listLoading.value = true;
  const res = await getSettlementList(listQuery.value)
  if ('code' in res && res.code === 0) {
    listLoading.value = false;
    productList.value = res.data.list
    total.value = res.data.total;
  }
}

const statusOptions = [
  {
      value: 100,
      label: '未结算'
  },
  {
      value: 101,
      label: '已结算'
  },
]

const formatDate = (time) => {
  if (time == null || time === '') {
    return 'N/A';
  }
  let date = new Date(time);
  return formatDateUtil(date, 'yyyy-MM-dd HH:mm:ss');
};

const multipleSelection = ref()
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 结算状态过滤器
const verifyStatusFilter = computed(() => {
  return (value: any) => {
      if (value === 0) {
        return '未结算';
      } else {
        return '已结算';
      }
  }
})

const handleSettlement = (index, row) => {
  ElMessageBox.confirm(
    `是否确认对结算单号为 <strong>${row.settlementNo}</strong> 的结算？`,
    '操作确认',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      dangerouslyUseHTMLString: true,
      type: 'warning',
    }
  )
    .then(async() => {
      // 这里调用结算接口
      // 假设结算接口调用成功
        await updateSettlementState({ id: row.id, status: 1 }).then((res) => {
          if ('code' in res && res.code === 0) {
            ElMessage({
                type: 'success',
                message: '结算成功!',
            });
            // 刷新列表
            getList();
          } else {
            ElMessage({
              type: 'error',
              message:  res.msg || '结算失败，请重试！',
            });
            return;
          }
        });
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '已取消结算',
      });
    });
};
</script>
<style></style>


