<template>
  <div class="app-container">
    <el-card shadow="never" class="operate-container">
      <i class="el-icon-tickets"></i>
      <span>数据列表</span>
    </el-card>
    <div class="table-container">
      <el-table ref="selectSessionTable"
                :data="list"
                style="width: 100%;"
                v-loading="listLoading" border>
        <el-table-column label="编号" width="100" align="center" prop="id">
        </el-table-column>
        <el-table-column label="秒杀时间段名称" align="center" prop="name">
        </el-table-column>
        <el-table-column label="每日开始时间" align="center">
          <template #default="scope">{{ formatDate(scope.row.startTime) }}</template>
        </el-table-column>
        <el-table-column label="每日结束时间" align="center">
          <template #default="scope">{{ formatDate(scope.row.endTime) }}</template>
        </el-table-column>
        <el-table-column label="商品数量" align="center" prop="productCount">
        </el-table-column>
        <el-table-column label="操作" align="center">
          <template #default="scope">
            <el-button size="small"
                      key="primary"
                      type="primary"
                      link
                      @click="handleShowRelation(scope.$index, scope.row)">商品列表
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>
<script>
  import {fetchSelectList} from '@/api/flashSession';
  import {formatDate} from '@/utils/date';
  export default {
    name: 'selectSessionList',
    data() {
      return {
        list: null,
        listLoading: false
      }
    },
    created() {
      this.getList();
    },
    computed: {
      formatDate() {
        return (time) => {
          if (time == null || time === '') {
            return 'N/A';
          }
          let date = new Date(time);
          return formatDate(date, 'yyyy-MM-dd HH:mm:ss')
        }
      }
      
    },
    methods: {
      handleShowRelation(index, row){
        this.$router.push({path: '/layout/marketing/productRelation', query: {
          flashPromotionId: this.$route.query.flashPromotionId, flashPromotionSessionId: row.id}})
      },
      getList() {
        this.listLoading = true;
        fetchSelectList({id: this.$route.query.flashPromotionId}).then(response => {
          this.listLoading = false;
          this.list = response.data.list;
        });
      }
    }
  }
</script>
<style scoped>
  .operate-container {
    margin-top: 0;
  }
</style>


