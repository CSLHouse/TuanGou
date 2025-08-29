<template>
    <div>
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
            <el-form-item label="输入搜索：">
              <el-input v-model="listQuery.orderSn" class="input-width" placeholder="订单编号" clearable></el-input>
            </el-form-item>
            <el-form-item label="收货人：">
              <el-input v-model="listQuery.receiverKeyword" class="input-width" placeholder="收货人姓名/手机号码" clearable></el-input>
            </el-form-item>
            <el-form-item label="提交时间：">
              <el-date-picker
                class="input-width"
                v-model="listQuery.createTime"
                value-format="YYYY-MM-DD"
                type="date"
                placeholder="请选择时间">
              </el-date-picker>
            </el-form-item>
            <el-form-item label="订单状态：">
              <el-select v-model="listQuery.state" class="input-width" placeholder="全部" clearable>
                <el-option v-for="item in statusOptions"
                           :key="item.value"
                           :label="item.label"
                           :value="item.value">
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="订单分类：">
              <el-select v-model="listQuery.orderType" class="input-width" placeholder="全部" clearable>
                <el-option v-for="item in orderTypeOptions"
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
      </el-card>
      <div class="table-container">
        <el-table ref="orderTable"
                  :data="list"
                  style="width: 100%;"
                  @selection-change="handleSelectionChange"
                  v-loading="listLoading" border>
          <el-table-column type="selection" width="60" align="center"></el-table-column>
          <el-table-column label="编号" width="60" align="center">
            <template #default="scope">{{scope.row.id}}</template>
          </el-table-column>
          <el-table-column label="订单编号" width="180" align="center">
            <template #default="scope">{{scope.row.orderSn}}</template>
          </el-table-column>
          <el-table-column label="提交时间" width="180" align="center">
            <template #default="scope">{{ formatTime(scope.row.CreatedAt) }}</template>
          </el-table-column>
          <el-table-column label="用户账号" width="140" align="center">
            <template #default="scope">{{scope.row.userName}}</template>
          </el-table-column>
          <el-table-column label="订单金额" width="120" align="center">
            <template #default="scope">￥{{scope.row.totalAmount}}</template>
          </el-table-column>
          <el-table-column label="支付方式" width="120" align="center">
            <template #default="scope">{{formatPayType(scope.row.payType)}}</template>
          </el-table-column>
          <el-table-column label="订单状态" width="120" align="center">
            <template #default="scope">{{formatStatus(scope.row.status)}}</template>
          </el-table-column>
          <el-table-column label="操作" width="200" align="center">
            <template #default="scope">
              <el-button
                size="small"
                @click="handleViewOrder(scope.$index, scope.row)"
              >查看订单</el-button>
              <el-button
                size="small"
                @click="handleCloseOrder(scope.$index, scope.row)"
                v-show="scope.row.status===0">关闭订单</el-button>
              <el-button
                size="small"
                @click="handleCompleteOrder(scope.$index, scope.row)"
                v-show="scope.row.status===1 || scope.row.status===2||scope.row.status===3">已完成</el-button>
              <el-button
                size="small"
                type="danger"
                @click="handleDeleteOrder(scope.$index, scope.row)"
                v-show="scope.row.status===4">删除订单</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="batch-operate-container">
        <el-select
          size="small"
          v-model="operateType" placeholder="批量操作">
          <el-option
            v-for="item in operateOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
        <el-button
          style="margin-left: 20px"
          class="search-button"
          @click="handleBatchOperate()"
          type="primary"
          size="small">
          确定
        </el-button>
      </div>
      <div class="pagination-container">
        <el-pagination
          background
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          layout="total,prev, pager, next,jumper"
          v-model:current-page="listQuery.page"
          v-model:page-size="listQuery.pageSize"
          :page-sizes="[5,10,15]"
          :total.number="+total">
        </el-pagination>
      </div>
      <el-dialog
        title="关闭订单"
        v-model="closeOrder.dialogVisible" width="30%">
        <span style="vertical-align: top">操作备注：</span>
        <el-input
          style="width: 80%"
          type="textarea"
          :rows="5"
          placeholder="请输入内容"
          v-model="closeOrder.content">
        </el-input>
        <span slot="footer" class="dialog-footer">
          <el-button @click="closeOrder.dialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="handleCloseOrderConfirm">确 定</el-button>
        </span>
      </el-dialog>
      <!-- <logistics-dialog v-model="logisticsDialogVisible"></logistics-dialog> -->
    </div>
  </template>
  <script>
    import {fetchList,closeOrder,deleteOrder,updateOrderCompletedStatus} from '@/api/order'
    import {formatDate} from '@/utils/date';
    import { useUserStore } from '@/pinia/modules/user'
    const userStore = useUserStore()
    const defaultListQuery = {
      page: 1,
      pageSize: 10,
      orderSn: null,
      receiverKeyword: null,
      state: null,
      orderType: null,
    //   sourceType: null,
      createTime: null,
    };
    export default {
      name: "orderList",
      data() {
        return {
          listQuery: Object.assign({}, defaultListQuery),
          listLoading: true,
          list: null,
          total: null,
          operateType: null,
          multipleSelection: [],
          closeOrder:{
            dialogVisible:false,
            content:null,
            orderIds:[]
          },
          statusOptions: [
            {
              label: '待付款',
              value: 0
            },
            {
              label: '待发货',
              value: 1
            },
            {
              label: '已发货',
              value: 2
            },
            {
              label: '已完成',
              value: 3
            },
            {
              label: '已关闭',
              value: 4
            }
          ],
          orderTypeOptions: [
            {
              label: '正常订单',
              value: 100
            },
            {
              label: '秒杀订单',
              value: 101
            }
          ],
          operateOptions: [
            // {
            //   label: "批量发货",
            //   value: 1
            // },
            {
              label: "关闭订单",
              value: 2
            },
            {
              label: "删除订单",
              value: 3
            }
          ],
          logisticsDialogVisible:false
        }
      },
      created() {
        this.getList();
      },
      computed: {
        formatTime() {
            return (time) => {
                if (time == null || time === '') {
                    return 'N/A';
                }
                let date = new Date(time);
                return formatDate(date, 'yyyy-MM-dd HH:mm:ss')
            }
        },
        formatPayType() {
            return (value) => {
                if (value === 1) {
                    return '支付宝';
                } else if (value === 2) {
                    return '微信';
                } else {
                    return '未支付';
                }
            }
        },
        formatStatus() {
            return (value) => {
                if (value === 1) {
                    return '待发货';
                } else if (value === 2) {
                    return '已发货';
                } else if (value === 3) {
                    return '已完成';
                } else if (value === 4) {
                    return '已关闭';
                } else if (value === 5) {
                    return '无效订单';
                } else {
                    return '待付款';
                }
            }
        },
      },
      methods: {
        handleResetSearch() {
          this.listQuery = Object.assign({}, defaultListQuery);
        },
        handleSearchList() {
          this.listQuery.page = 1;
          this.getList();
        },
        handleSelectionChange(val){
          this.multipleSelection = val;
        },
        handleViewOrder(index, row){
          this.$router.push({path:'/layout/orderManage/orderDetail',query:{id:row.id}})
        },
        handleCloseOrder(index, row){
          this.closeOrder.dialogVisible = true;
          this.closeOrder.orderIds = [row.id];
        },
        handleViewLogistics(index, row){
          this.logisticsDialogVisible=true;
        },
        handleDeleteOrder(index, row){
          let ids=[];
          ids.push(row.id);
          this.deleteOrder(ids);
        },
        handleBatchOperate(){
          if(this.multipleSelection==null||this.multipleSelection.length<1){
            this.$message({
              message: '请选择要操作的订单',
              type: 'warning',
              duration: 1000
            });
            return;
          }
          if(this.operateType === 1){
            //批量发货
          }else if(this.operateType === 2){
            //关闭订单
            this.closeOrder.orderIds = [];
            for(let i=0;i<this.multipleSelection.length;i++){
              this.closeOrder.orderIds.push(this.multipleSelection[i].id);
            }
            this.closeOrder.dialogVisible=true;
          }else if(this.operateType === 3){
            //删除订单
            let ids=[];
            for(let i=0;i<this.multipleSelection.length;i++){
              ids.push(this.multipleSelection[i].id);
            }
            this.deleteOrder(ids);
          }
        },
        handleSizeChange(val){
          this.listQuery.page = 1;
          this.listQuery.pageSize = val;
          this.getList();
        },
        handleCurrentChange(val){
          this.listQuery.page = val;
          this.getList();
        },
        handleCloseOrderConfirm() {
          if (this.closeOrder.content == null || this.closeOrder.content === '') {
            this.$message({
              message: '操作备注不能为空',
              type: 'warning',
              duration: 1000
            });
            return;
          }
          let params = {"ids": this.closeOrder.orderIds, "note": this.closeOrder.content}
          closeOrder(params).then(response=>{
              this.closeOrder.orderIds=[];
              this.closeOrder.dialogVisible=false;
              this.getList();
              this.$message({
                message: '修改成功',
                type: 'success',
                duration: 1000
              });
            });
        },
        getList() {
          this.listLoading = true;
          fetchList(this.listQuery).then(response => {
            this.listLoading = false;
            this.list = response.data.list;
            this.total = response.data.total;
          });
        },
        deleteOrder(ids){
          this.$confirm('是否要进行该删除操作?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            let data = {"ids": ids}
            deleteOrder(data).then(response=>{
              this.$message({
                message: '删除成功！',
                type: 'success',
                duration: 1000
              });
              this.getList();
            });
          })
        },
        handleCompleteOrder(index, row) {
          updateOrderCompletedStatus({"ids": [row.id]}).then(response=>{
            this.$message({
              message: '更新成功！',
              type: 'success',
              duration: 1000
            });
          })
        },
      }
    }
  </script>
  <style scoped>
    .input-width {
      width: 203px;
    }
  </style>
  
  
  