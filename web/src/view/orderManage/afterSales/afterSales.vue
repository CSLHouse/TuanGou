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
            <el-form-item label="输入联系方式：">
              <el-input v-model="listQuery.contact" class="input-width" placeholder="手机/QQ/邮箱" clearable></el-input>
            </el-form-item>
            <el-form-item label="申请状态：">
              <el-select v-model="listQuery.state" class="input-width" placeholder="全部" clearable>
                <el-option v-for="item in statusOptions"
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
                  v-loading="listLoading" border>
          <el-table-column type="selection" width="60" align="center"></el-table-column>
          <el-table-column label="编号" width="60" align="center">
            <template #default="scope">{{scope.row.id}}</template>
          </el-table-column>
          <el-table-column label="用户ID" align="center">
            <template #default="scope">{{scope.row.userId}}</template>
          </el-table-column>
          <el-table-column label="订单ID"  align="center">
            <template #default="scope">{{scope.row.orderItemId}}</template>
          </el-table-column>
          <el-table-column label="金额"  align="center">
            <template #default="scope">{{scope.row.realAmount}}</template>
          </el-table-column>
          <el-table-column label="反馈内容" width="140" align="center">
            <template #default="scope">{{scope.row.content}}</template>
          </el-table-column>
          <el-table-column label="联系方式" width="120" align="center">
            <template #default="scope">{{scope.row.contact}}</template>
          </el-table-column>
          <el-table-column label="图片" min-width="192" align="center">
            <template #default="scope">
              <div class="image-container">
                <img 
                    v-for="(item, index) in scope.row.imagesList" 
                    :key="item"  
                    :src="item" 
                    alt="图片" 
                    class="preview-img"
                    @click="handlePreview(item)"
                >
                </div>
            </template>
          </el-table-column>
          <el-table-column label="申请状态" width="120" align="center">
            <template #default="scope">{{formatStatus(scope.row.status)}}</template>
          </el-table-column>
          
          <el-table-column label="操作" width="200" align="center">
            <template #default="scope">
              <el-button
                size="small"
                @click="handleUpdateOrder(scope.$index, scope.row)"
                >关闭申请</el-button>
            </template>
          </el-table-column>
        </el-table>
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
    </div>
    
    <el-dialog 
        v-model="dialogVisible" 
        title="图片预览" 
        :close-on-click-modal="true" 
        :width="800"
    >
        <img 
        :src="previewImage" 
        alt="预览图" 
        class="preview-large-img"
        >
    </el-dialog>

  </template>
  <script>
    import {getOrderDealList,updateOrderDealStatus} from '@/api/order'
    import {formatDate} from '@/utils/date';
    const defaultListQuery = {
      page: 1,
      pageSize: 10,
      contact: null,
      state: -1,
    };
    export default {
      name: "afterSales",
      data() {
        return {
          listQuery: Object.assign({}, defaultListQuery),
          listLoading: true,
          list: null,
          total: null,
          statusOptions: [
            {
              label: '未处理',
              value: 0
            },
            {
              label: '已处理',
              value: 1
            },
          ],
          dialogVisible: false,
          previewImage: '',
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
        formatStatus() {
            return (value) => {
                if (value === 1) {
                    return '已处理';
                }  else {
                    return '未处理';
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
        handleUpdateOrder(index, row){
            if(row.status === 1){
                this.$message({
                    message: '该申请已处理，无需重复操作',
                    type: 'warning',
                    duration: 1000
                });
                return;
            }
            updateOrderDealStatus({"dealId": row.id, "orderItemId": row.orderItemId, "status": 1}).then(response=>{
                if(response.code === 0){
                this.$message({
                    message: '修改成功',
                    type: 'success',
                    duration: 1000
                });
                this.getList();
                } else {
                this.$message({
                    message: response.message,
                    type: 'warning',
                    duration: 1000
                });
                }
            });
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
          getOrderDealList(this.listQuery).then(response => {
            this.listLoading = false;
            this.list = response.data.list;
            this.total = response.data.total;
          });
        },
        handlePreview(imgUrl) {
          this.previewImage = imgUrl;
          this.dialogVisible = true;
        }
      }
    }
  </script>
  <style scoped>
    .input-width {
      width: 203px;
    }
    .edit-icon {
      color: #409eff;
      cursor: pointer;
      font-size: 16px;
      transition: color 0.2s;
    }

    .edit-icon:hover {
      color: #66b1ff;
    }

    .edit-container {
      position: relative;
    }

    /* 图片容器：自适应排列 */
    .image-container {
    display: flex;        /* 启用Flex布局 */
    flex-wrap: wrap;      /* 超出换行 */
    gap: 5px;             /* 图片间距（替代margin-right） */
    width: 100%;          /* 占满表格列宽度 */
    padding: 5px 0;       /* 上下内边距 */
    }

    /* 缩略图样式 */
    .preview-img {
    width: 50px;
    height: 50px;
    object-fit: cover;    /* 保持比例，裁剪多余部分 */
    cursor: pointer;      /* 鼠标悬停显示手型 */
    border: 1px solid #eee; /* 轻微边框 */
    }

    /* 弹窗内大图样式 */
    .preview-large-img {
    display: block;
    margin: 0 auto;
    max-width: 100%;      /* 不超过弹窗宽度 */
    max-height: 70vh;     /* 不超过屏幕高度的70% */
    object-fit: contain;  /* 保持完整显示 */
    }
  </style>
  
  
  