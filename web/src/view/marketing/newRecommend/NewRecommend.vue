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
          <el-form-item label="商品名称：">
            <el-input v-model="listQuery.keyword" class="input-width" clearable placeholder="商品名称"></el-input>
          </el-form-item>
          <el-form-item label="推荐状态：">
            <el-select v-model="listQuery.state" placeholder="全部" clearable class="input-width">
              <el-option v-for="item in recommendOptions"
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
      <el-button size="mini" class="btn-add" @click="handleSelectProduct()">选择商品</el-button>
    </el-card>
    <div class="table-container">
      <el-table ref="newProductTable"
                :data="list"
                style="width: 100%;"
                @selection-change="handleSelectionChange"
                v-loading="listLoading" border>
        <el-table-column type="selection" width="60" align="center"></el-table-column>
        <el-table-column label="编号" width="120" align="center" prop="id">
          <template #default="scope">{{scope.row.id}}</template>
        </el-table-column>
        <el-table-column label="商品名称" align="center">
          <template #default="scope">{{scope.row.productName}}</template>
        </el-table-column>
        <el-table-column label="是否推荐" width="200" align="center">
          <template #default="scope">
            <el-switch
              @change="handleRecommendStatusStatusChange(scope.$index, scope.row)"
              :active-value="1"
              :inactive-value="0"
              v-model="scope.row.recommendStatus">
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column label="排序" width="160" align="center">
          <template #default="scope">{{scope.row.sort}}</template>
        </el-table-column>
        <el-table-column label="状态" width="160" align="center">
          <template #default="scope">{{ formatRecommendStatus(scope.row.recommendStatus) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center">
          <template #default="scope">
            <el-button size="mini"
                       type="text"
                       @click="handleEditSort(scope.$index, scope.row)">设置排序
            </el-button>
            <el-button size="mini"
                       type="text"
                       @click="handleDelete(scope.$index, scope.row)">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div class="batch-operate-container">
      <el-select
        size="small"
        v-model="operateType" placeholder="批量操作">
        <el-option
          v-for="item in operates"
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
        layout="total, sizes,prev, pager, next,jumper"
        :page-size="listQuery.pageSize"
        :page-sizes="[5,10,15]"
        :current-page.sync="listQuery.page"
        :total="total">
      </el-pagination>
    </div>

    <el-dialog title="选择商品" v-model="selectDialogVisible"  width="50%">
      <el-input v-model="dialogData.listQuery.keyword"
                style="width: 250px;margin-bottom: 20px"
                size="small"
                placeholder="商品名称搜索">
        <template #append>
            <el-button  :icon="Search" @click="handleSelectSearch()"></el-button>
        </template>
        
      </el-input>
      <el-table :data="dialogData.list"
                @selection-change="handleDialogSelectionChange" border>
        <el-table-column type="selection" :selectable="selectable" width="60" align="center"></el-table-column>
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
          :total="dialogData.total">
        </el-pagination>
      </div>
      <div style="clear: both;"></div>
      <div slot="footer">
        <el-button  size="small" @click="selectDialogVisible = false">取 消</el-button>
        <el-button  size="small" type="primary" @click="handleSelectDialogConfirm()">确 定</el-button>
      </div>
    </el-dialog>
    
    <el-dialog title="设置排序"
               v-model="sortDialogVisible"
               width="40%">
      <el-form :model="sortDialogData"
               label-width="150px">
        <el-form-item label="排序：">
          <el-input v-model="sortDialogData.sort" type="number" style="width: 200px"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="sortDialogVisible = false" size="small">取 消</el-button>
        <el-button type="primary" @click="handleUpdateSort" size="small">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
  import {fetchList,updateRecommendStatus,deleteNewProduct,createNewProduct,updateNewProductSort} from '@/api/newProduct';
  import {getProductList as fetchProductList} from '@/api/product';
  import { ref, reactive, onBeforeMount, computed } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { Search } from '@element-plus/icons-vue'
  const defaultListQuery = {
    page: 1,
    pageSize: 5,
    keyword: null,
    state: null
  };
  const defaultRecommendOptions = [
    {
      label: '未推荐',
      value: 100
    },
    {
      label: '推荐中',
      value: 101
    }
  ];

    const listQuery = reactive(defaultListQuery)
    
    const recommendOptions = ref(defaultRecommendOptions);
    const total = ref(0)
    const list = ref([])
    const listLoading = ref(false)
    const multipleSelection = ref([]);
    const operates = ref([
        { label: "设为推荐", value: 100 },
        { label: "取消推荐", value: 101 },
        { label: "删除", value: 2 }
    ]);
    
    const operateType = ref(null);
    const selectDialogVisible = ref(false);
    const dialogData = reactive({
      list: null,
      total: null,
      multipleSelection: [],
      listQuery: {
        keyword: null,
        page: 1,
        pageSize: 5
      }
    });
    const sortDialogVisible = ref(false);
    const sortDialogData = reactive({sort: 0, id: null});

    const recommendStatusParams = ref({
        ids: [],
        key: 'recommend_status',
        value: 0
    });

    const newProductIds = ref([]);
    const selectable = ref((row) => {
      return !newProductIds.value.includes(row.id);
    });


    onBeforeMount(() => {
        getList()
    })

    const formatRecommendStatus = computed((status) => {
      return (status) => {
        if (status === 1) {
          return '推荐中';
        } else {
          return '未推荐';
        }
      };
    });
    const handleResetSearch = () => {
      Object.assign(listQuery, defaultListQuery);
    };

    const handleSearchList = () => {
      listQuery.page = 1;
      getList();
    };

    const handleSelectionChange = (val) => {
      multipleSelection.value = val;
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

    const handleRecommendStatusStatusChange = (index, row) => {
      updateRecommendStatusStatus([row.id], row.recommendStatus + 100);
    };

    const handleDelete = (index, row) => {

      deleteProduct([row.id]);
    };

    const handleBatchOperate = () => {
      if (multipleSelection.value.length < 1) {
        ElMessage({
          message: '请选择一条记录',
          type: 'warning',
          duration: 1000
        });
        return;
      }
      let ids = [];
      for (let i = 0; i < multipleSelection.value.length; i++) {
        ids.push(multipleSelection.value[i].id);
      }
      if (operateType.value === 0) {
        //设为推荐
        updateRecommendStatusStatus(ids, 101);
      } else if (operateType.value === 1) {
        //取消推荐
        updateRecommendStatusStatus(ids, 100);
      } else if(operateType.value===2){
        //删除
        deleteProduct(ids);
      }else {
        ElMessage({
          message: '请选择批量操作类型',
          type: 'warning',
          duration: 1000
        });
      }
    };

    const handleSelectProduct = () => {
      selectDialogVisible.value = true;
      dialogData.listQuery.page = 1;
      dialogData.listQuery.pageSize = 5;
      getDialogList();
    };

    const handleSelectSearch = () => {
      dialogData.listQuery.page = 1;
      getDialogList();
    };

    const handleDialogSizeChange = (val) => {
      dialogData.listQuery.page = 1;
      dialogData.listQuery.pageSize = val;
      getDialogList();
    };

    const handleDialogCurrentChange = (val) => {
      dialogData.listQuery.page = val;
      getDialogList();
    };

    const handleDialogSelectionChange = (val) => {
      dialogData.multipleSelection = val;
    };

    const handleSelectDialogConfirm = () => {
      if (dialogData.multipleSelection.length < 1) {
        ElMessage({
          message: '请选择一条记录',
          type: 'warning',
          duration: 1000
        });
        return;
      }
      let selectProducts = {newProducts: []};
      for (let i = 0; i < dialogData.multipleSelection.length; i++) {
        selectProducts.newProducts.push({
          productId: dialogData.multipleSelection[i].id,
          productName: dialogData.multipleSelection[i].name
        });
      }

      console.log("选择的商品", selectProducts);
      ElMessageBox.confirm('使用要进行添加操作?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        createNewProduct(selectProducts).then(response => {
          selectDialogVisible.value = false;
          dialogData.multipleSelection = [];
          getList();
          ElMessage({
            type: 'success',
            message: '添加成功!'
          });
        });
      });
    };

    const handleEditSort = (index, row) => {
      sortDialogVisible.value = true;
      sortDialogData.sort = Number(row.sort);
      sortDialogData.id = row.id;
    };

    const handleUpdateSort = () => {
      ElMessageBox.confirm('是否要修改排序?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        sortDialogData.sort = Number(sortDialogData.sort);
        updateNewProductSort(sortDialogData).then(response => {
          sortDialogVisible.value = false;
          getList();
          ElMessage({
            type: 'success',
            message: '修改成功!'
          });
        });
      });
    };

    const getList = () => {
      listLoading.value = true;
      fetchList(listQuery).then(response => {
        listLoading.value = false;
        list.value = response.data.list;
        total.value = response.data.total;
        newProductIds.value = response.data.productIds
        console.log("获取推荐列表", list.value);
      });
    };

    const updateRecommendStatusStatus = (ids, status) => {
      ElMessageBox.confirm('是否要修改推荐状态?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        recommendStatusParams.value.ids = ids;
        recommendStatusParams.value.value = status;
        updateRecommendStatus(recommendStatusParams.value).then(response => {
          getList();
          ElMessage({
            type: 'success',
            message: '修改成功!'
          });
        });
      }).catch(() => {
        ElMessage({
          type: 'success',
          message: '已取消操作!'
        });
        getList();
      });
    };

    const deleteProduct = (ids) => {
      ElMessageBox.confirm('是否要删除推荐?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let params = new URLSearchParams();
        params.append("ids", ids);
        deleteNewProduct(params).then(response => {
          getList();
          ElMessage({
            type: 'success',
            message: '删除成功!'
          });
        });
      }).catch(() => {
        ElMessage({
          type: 'info',
          message: '已取消操作!'
        });
      });
    };

    const getDialogList = () => {
      dialogData.listQuery.page = 1;
      dialogData.listQuery.pageSize = 5;
      fetchProductList(dialogData.listQuery).then(response => {
        dialogData.list = response.data.list;
        dialogData.total = response.data.total;
        // console.log("获取商品列表 selectDialogVisible", selectDialogVisible.value);
      });
    };

    
  </script>
  
  <style scoped>

  </style>