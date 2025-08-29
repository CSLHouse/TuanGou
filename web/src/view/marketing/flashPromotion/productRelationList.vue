<template>
  <div class="app-container">
    <el-card class="operate-container" shadow="never">
      <i class="el-icon-tickets"></i>
      <span>数据列表</span>
      <el-button size="small" class="btn-add" @click="handleSelectProduct()" style="margin-left: 20px">添加</el-button>
    </el-card>
    <div class="table-container">
      <el-table ref="productRelationTable"
                :data="list"
                style="width: 100%;"
                v-loading="listLoading" border>
        <el-table-column label="编号" width="100" align="center" prop="id">
          <!-- <template #default="scope">{{scope.row.id}}</template> -->
        </el-table-column>
        <el-table-column label="商品名称" align="center">
          <template #default="scope">{{scope.row.product.name}}</template>
        </el-table-column>
        <el-table-column label="货号" width="140" align="center">
          <template #default="scope">NO.{{scope.row.product.productSN}}</template>
        </el-table-column>
        <el-table-column label="商品价格" width="100" align="center">
          <template #default="scope">￥{{scope.row.product.price}}</template>
        </el-table-column>
        <el-table-column label="剩余数量" width="100" align="center">
          <template #default="scope">{{scope.row.product.stock}}</template>
        </el-table-column>
        <el-table-column label="秒杀价格" width="100" align="center">
          <template #default="scope">
            <p v-if="scope.row.flashPromotionPrice!==null">
              ￥{{scope.row.flashPromotionPrice}}
            </p>
          </template>
        </el-table-column>
        <el-table-column label="秒杀数量" width="100" align="center" prop="flashPromotionCount">
          <!-- <template #default="scope">{{scope.row.flashPromotionCount}}</template> -->
        </el-table-column>
        <el-table-column label="限购数量" width="100" align="center" prop="flashPromotionLimit">
          <!-- <template #default="scope">{{scope.row.flashPromotionLimit}}</template> -->
        </el-table-column>
        <el-table-column label="排序" width="100" align="center" prop="sort">
          <!-- <template #default="scope">{{scope.row.sort}}</template> -->
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template #default="scope">
            <el-button size="small"
                      key="primary"
                      type="primary"
                      link
                      @click="handleUpdate(scope.$index, scope.row)">编辑
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
        background
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        layout="total, sizes,prev, pager, next,jumper"
        :current-page.sync="listQuery.page"
        :page-size="listQuery.pageSize"
        :page-sizes="[5,10,15]"
        :total.number="+total">
      </el-pagination>
    </div>
    <el-dialog title="选择商品" v-model="selectDialogVisible" width="50%">
      <el-input v-model="dialogData.listQuery.keyword"
                style="width: 250px;margin-bottom: 20px"
                size="small"
                placeholder="商品名称搜索">
        <el-button slot="append" icon="el-icon-search" @click="handleSelectSearch()"></el-button>
      </el-input>
      <el-table :data="dialogData.list"
                @selection-change="handleDialogSelectionChange" border>
        <el-table-column type="selection" width="60" align="center" :selectable="checkSelectedStatus"/>
        <el-table-column label="商品名称" align="center">
          <template #default="scope">{{scope.row.name}}</template>
        </el-table-column>
        <el-table-column label="货号" width="160" align="center">
          <template #default="scope">NO.{{scope.row.productSN}}</template>
        </el-table-column>
        <el-table-column label="价格" width="120" align="center">
          <template #default="scope">￥{{scope.row.price}}</template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination
          background
          @size-change="handleDialogSizeChange"
          @current-change="handleDialogCurrentChange"
          layout="prev, pager, next"
          :current-page.sync="dialogData.listQuery.page"
          :page-size="dialogData.listQuery.pageSize"
          :page-sizes="[5,10,15]"
          :total.number="+dialogData.total">
        </el-pagination>
      </div>
      <div style="clear: both;"></div>
      <div slot="footer">
        <el-button  size="small" @click="selectDialogVisible = false">取 消</el-button>
        <el-button  size="small" type="primary" @click="handleSelectDialogConfirm()">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="编辑秒杀商品信息"
      v-model="editDialogVisible"
      width="40%">
      <el-form :model="flashProductRelation"
               ref="flashProductRelationForm"
               label-width="150px" size="small">
        <el-form-item label="商品名称：">
          <span>{{flashProductRelation.product.name}}</span>
        </el-form-item>
        <el-form-item label="货号：">
          <span>NO.{{flashProductRelation.product.productSN}}</span>
        </el-form-item>
        <el-form-item label="商品价格：">
          <span>￥{{flashProductRelation.product.price}}</span>
        </el-form-item>
        <el-form-item label="秒杀价格：">
          <el-input-number v-model="flashProductRelation.flashPromotionPrice" :precision="2" :controls="false" class="input-width"/>
          <!-- <el-input v-model="flashProductRelation.flashPromotionPrice" class="input-width">
            <template slot="prepend">￥</template>
          </el-input> -->
        </el-form-item>
        <el-form-item label="剩余数量：">
          <span>{{flashProductRelation.product.stock}}</span>
        </el-form-item>
        <el-form-item label="秒杀数量：">
          <el-input v-model.number="flashProductRelation.flashPromotionCount" class="input-width"></el-input>
        </el-form-item>
        <el-form-item label="限购数量：">
          <el-input v-model.number="flashProductRelation.flashPromotionLimit" class="input-width"></el-input>
        </el-form-item>
        <el-form-item label="排序：">
          <el-input v-model.number="flashProductRelation.sort" class="input-width"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="editDialogVisible = false" size="small">取 消</el-button>
        <el-button type="primary" @click="handleEditDialogConfirm()" size="small">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
  import {fetchList,createFlashProductRelation,deleteFlashProductRelation,updateFlashProductRelation} from '@/api/flashProductRelation';
  import {getProductList as fetchProductList} from '@/api/product';
  const defaultListQuery = {
    page: 1,
    pageSize: 5,
    flashPromotionId: null,
    flashPromotionSessionId: null
  };
  export default {
    name:'flashPromotionProductRelationList',
    data() {
      return {
        listQuery: Object.assign({}, defaultListQuery),
        list: null,
        total: 0,
        listLoading: false,
        dialogVisible: false,
        selectDialogVisible:false,
        dialogData:{
          list: null,
          total: 0,
          multipleSelection:[],
          listQuery:{
            keyword: null,
            page: 1,
            pageSize: 5
          }
        },
        editDialogVisible:false,
        flashProductRelation:{
          product:{}
        }
      }
    },
    created(){
      this.listQuery.flashPromotionId = parseInt(this.$route.query.flashPromotionId);
      this.listQuery.flashPromotionSessionId = parseInt(this.$route.query.flashPromotionSessionId);
      this.getList();
    },
    methods:{
      handleSizeChange(val) {
        this.listQuery.page = 1;
        this.listQuery.pageSize = val;
        this.getList();
      },
      handleCurrentChange(val) {
        this.listQuery.page = val;
        this.getList();
      },
      handleSelectProduct(){
        this.selectDialogVisible=true;
        this.getDialogList();
      },
      handleUpdate(index,row){
        this.editDialogVisible = true;
        this.flashProductRelation = Object.assign({},row);
      },
      handleDelete(index,row){
        this.$confirm('是否要删除该商品?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          deleteFlashProductRelation({id: row.id}).then(response => {
            this.$message({
              type: 'success',
              message: '删除成功!'
            });
            this.getList();
          });
        });
      },
      handleSelectSearch(){
        this.getDialogList();
      },
      handleDialogSizeChange(val) {
        this.dialogData.listQuery.page = 1;
        this.dialogData.listQuery.pageSize = val;
        this.getDialogList();
      },
      handleDialogCurrentChange(val) {
        this.dialogData.listQuery.page = val;
        this.getDialogList();
      },
      handleDialogSelectionChange(val){
        this.dialogData.multipleSelection = val;
      },
      handleSelectDialogConfirm(){
        if (this.dialogData.multipleSelection < 1) {
          this.$message({
            message: '请选择一条记录',
            type: 'warning',
            duration: 1000
          });
          return;
        }
        let selectProducts = [];
        for (let i = 0; i < this.dialogData.multipleSelection.length; i++) {
          selectProducts.push({
            productId: this.dialogData.multipleSelection[i].id,
            flashPromotionId: this.listQuery.flashPromotionId,
            flashPromotionSessionId: this.listQuery.flashPromotionSessionId
          });
        }
        this.$confirm('使用要进行添加操作?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          createFlashProductRelation({flashes: selectProducts}).then(response=>{
            this.selectDialogVisible = false;
            this.dialogData.multipleSelection = [];
            this.getList();
            this.$message({
              type: 'success',
              message: '添加成功!'
            });
          });
        });
      },
      handleEditDialogConfirm(){
        this.$confirm('是否要确认?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
            updateFlashProductRelation(this.flashProductRelation).then(response => {
              this.$message({
                message: '修改成功！',
                type: 'success'
              });
              this.editDialogVisible =false;
              this.getList();
            })
        })
      },
      getList() {
        this.listLoading = true;
        fetchList(this.listQuery).then(response => {
          this.listLoading = false;
          this.list = response.data.list;
          this.total = response.data.total;
        });
      },
      getDialogList(){
        fetchProductList(this.dialogData.listQuery).then(response=>{
          // for (var i = response.data.list.length - 1; i >= 0; i--) {
          //   this.list.forEach(element => {
          //     if (element.product.id == response.data.list[i].id) {
          //       response.data.list.splice(i, 1)
          //     }
          //   });
          // }
          this.dialogData.list = response.data.list;
          this.dialogData.total = response.data.total;
        })
      },
      checkSelectedStatus(row) {
        if (this.list.length < 1) {
          return true
        }
        for (var i = 0; i < this.list.length; i++) {
          if (this.list[i].product.id == row.id) {
            return false
          }
        }
        return true
      }
    }
  }
</script>
<style scoped>
  .operate-container{
    margin-top: 0;
  }
  .input-width{
    width: 200px;
  }
</style>


