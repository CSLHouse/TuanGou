<template>
    <div>
      <el-card shadow="never">
        <div>
          <el-icon><Search /></el-icon>
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
            <el-form-item label="活动名称：">
              <el-input v-model="listQuery.keyword" class="input-width" placeholder="活动名称" clearable></el-input>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
      <el-card class="operate-container" shadow="never">
        <el-icon><Grid /></el-icon>
        <span>数据列表</span>
        <el-button size="small" class="btn-add" @click="handleAdd()" style="margin-left: 20px">添加活动</el-button>
        <el-button size="small" class="btn-add" @click="handleShowSessionList()">秒杀时间段列表</el-button>
      </el-card>
      <div class="table-container">
        <el-table ref="flashTable"
                  :data="list"
                  style="width: 100%;"
                  v-loading="listLoading" border>
          <el-table-column type="selection" width="60" align="center"></el-table-column>
          <el-table-column label="编号" width="100" align="center" prop="id">
          </el-table-column>
          <el-table-column label="活动标题" align="center" prop="title">
          </el-table-column>
          <el-table-column label="活动状态" width="140" align="center">
            <template #default="scope">{{ formatActiveStatus(scope.row) }}</template>
          </el-table-column>
          <el-table-column label="开始时间" width="160" align="center">
            <template #default="scope">{{ formatDate(scope.row.startDate) }}</template>
          </el-table-column>
          <el-table-column label="结束时间" width="160" align="center">
            <template #default="scope">{{ formatDate(scope.row.endDate) }}</template>
          </el-table-column>
          <el-table-column label="上线/下线" width="100" align="center">
            <template #default="scope">
              <el-switch
                @change="handleStatusChange(scope.$index, scope.row)"
                :active-value="1"
                :inactive-value="0"
                v-model="scope.row.status">
              </el-switch>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" align="center">
            <template #default="scope">
              <el-button size="small"
                         key="primary"
                         type="primary"
                         link
                         @click="handleSelectSession(scope.$index, scope.row)">设置商品
              </el-button>
              <el-button size="small"
                         key="primary"
                         type="primary"
                         link
                         @click="handleUpdate(scope.$index, scope.row)">
                编辑
              </el-button>
              <el-button size="small"
                         key="danger"
                         type="danger"
                         link
                         @click="handleDelete(scope.$index, scope.row)">删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="listQuery.page"
          v-model:page-size="listQuery.pageSize"
          :page-sizes="[5, 10, 15]"
          small="small"
          :background="false"
          layout="total, sizes, prev, pager, next, jumper"
          :total.number="+total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
      <el-dialog
        title="添加活动"
        v-model="dialogVisible"
        width="40%">
        <el-form :model="flashPromotion"
                 ref="flashPromotionForm"
                 label-width="150px" size="small">
          <el-form-item label="活动标题：">
            <el-input v-model="flashPromotion.title" style="width: 250px"></el-input>
          </el-form-item>
          <el-form-item label="开始时间：">
            <el-date-picker
              v-model="flashPromotion.startDate"
              type="date"
              placeholder="请选择时间">
            </el-date-picker>
          </el-form-item>
          <el-form-item label="结束时间：">
            <el-date-picker
              v-model="flashPromotion.endDate"
              type="date"
              placeholder="请选择时间">
            </el-date-picker>
          </el-form-item>
          <el-form-item label="上线/下线">
            <el-radio-group v-model="flashPromotion.status">
              <el-radio :label="1">上线</el-radio>
              <el-radio :label="0">下线</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
        <span slot="footer" class="dialog-footer">
          <el-button @click="dialogVisible = false" size="small">取 消</el-button>
          <el-button type="primary" @click="handleDialogConfirm()" size="small">确 定</el-button>
        </span>
      </el-dialog>
    </div>
  </template>
  <script>
    import {fetchList, updateStatus, deleteFlash, createFlash, updateFlash} from '@/api/flash';
    import {formatDate} from '@/utils/date';
    import { computed } from 'vue';
  
    const defaultListQuery = {
      page: 1,
      pageSize: 5,
      keyword: null
    };
    const defaultFlashPromotion = {
      id: null,
      title: null,
      startDate: null,
      endDate: null,
      status: 0
    };
    export default {
      name: 'flashPromotion',
      data() {
        return {
          listQuery: Object.assign({}, defaultListQuery),
          list: null,
          total: 0,
          listLoading: false,
          dialogVisible: false,
          flashPromotion: Object.assign({}, defaultFlashPromotion),
          isEdit: false
        }
      },
      created() {
        this.getList();
      },
      computed: {
        formatActiveStatus() {
          return (row) => {
            let nowTime = new Date().getTime();
            let startDate = new Date(row.startDate).getTime()
            let endDate = new Date(row.endDate).getTime()
            if (nowTime >= startDate && nowTime <= endDate) {
              return '活动进行中';
            } else if (nowTime > endDate) {
              return '活动已结束';
            } else {
              return '活动未开始';
            }
          }
        },
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
        handleResetSearch() {
          this.listQuery = Object.assign({}, defaultListQuery);
          this.getList();
        },
        handleSearchList() {
          this.listQuery.page = 1;
          this.getList();
        },
        handleSizeChange(val) {
          this.listQuery.page = 1;
          this.listQuery.pageSize = val;
          this.getList();
        },
        handleCurrentChange(val) {
          this.listQuery.page = val;
          this.getList();
        },
        handleAdd() {
          this.dialogVisible = true;
          this.isEdit = false;
          this.flashPromotion = Object.assign({},defaultFlashPromotion);
        },
        handleShowSessionList() {
          this.$router.push({path: '/layout/marketing/session'})
        },
        handleStatusChange(index, row) {
          this.$confirm('是否要修改该状态?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            updateStatus(row.id, {status: row.status}).then(response => {
              this.$message({
                type: 'success',
                message: '修改成功!'
              });
            });
          }).catch(() => {
            this.$message({
              type: 'info',
              message: '取消修改'
            });
            this.getList();
          });
        },
        handleDelete(index, row) {
          this.$confirm('是否要删除该活动?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            deleteFlash({id: row.id}).then(response => {
              this.$message({
                type: 'success',
                message: '删除成功!'
              });
              this.getList();
            });
          });
        },
        handleUpdate(index, row) {
          this.dialogVisible = true;
          this.isEdit = true;
          this.flashPromotion = Object.assign({},row);
        },
        handleDialogConfirm() {
          this.$confirm('是否要确认?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            if (this.isEdit) {
              updateFlash(this.flashPromotion).then(response => {
                this.$message({
                  message: '修改成功！',
                  type: 'success'
                });
                this.dialogVisible =false;
                this.getList();
              })
            } else {
              createFlash(this.flashPromotion).then(response => {
                this.$message({
                  message: '添加成功！',
                  type: 'success'
                });
                this.dialogVisible =false;
                this.getList();
              })
            }
          })
        },
        handleSelectSession(index,row){
          this.$router.push({path:'/layout/marketing/selectSession', query: {flashPromotionId: row.id}})
        },
        getList() {
          this.listLoading = true;
          let that = this
          fetchList(this.listQuery).then(response => {
            that.listLoading = false;
            that.list = response.data.list;
            that.total = response.data.total;
          });
        }
      }
    }
  </script>
  <style scoped></style>
  
  
  