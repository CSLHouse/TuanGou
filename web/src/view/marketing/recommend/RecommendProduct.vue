<template>
    <div>
      <div>
        <el-card class="box-card" title="筛选搜索">
          <el-form :inline="true" :model="searchData" class="demo-form-inline">
            <el-form-item label="商品名称：" class="form-item">
              <el-input v-model="searchData.productName" placeholder="商品名称" clearable />
            </el-form-item>
            <el-form-item label="推荐状态：" class="form-item">
              <el-select v-model="searchData.recommendStatus" placeholder="全部" clearable class="input-width">
              <el-option v-for="item in recommendOptions"
                         :key="item.value"
                         :label="item.label"
                         :value="item.value">
              </el-option>
            </el-select>
            </el-form-item>
          </el-form>
          <el-button
            style="float:right; margin-right: 25px"
            @click="onReset()"
            size="small">
            重置
          </el-button>
          <el-button
            style="float:right"
            type="primary"
            @click="onSearch()"
            size="small">
            查询搜索
          </el-button>
        </el-card>
      </div>
      <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">选择商品</el-button>
        </div>
        <el-table
          ref="multipleTable"
          :data="hotProductData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="Id"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="40" />
          <el-table-column align="left" label="编号" prop="id" width="60"></el-table-column>
          <el-table-column align="left" label="商品编号" prop="productId" width="120" />
          <el-table-column align="left" label="商品名称" prop="productName" width="460" />
          <el-table-column align="left" label="是否推荐" prop="state" width="120" >
            <template #default="scope">
                <el-switch
                    @change="handleRecommendStatusStatusChange(scope.$index, scope.row)"
                    v-model="scope.row.recommendStatus"
                    :active-value="1"
                    :inactive-value="0">
                </el-switch>
            </template>
          </el-table-column>
          <el-table-column align="left" label="排序" prop="sort" width="120" />
          <el-table-column align="left" label="操作" min-width="160">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="handleEditSort(scope.$index, scope.row)">设置排序</el-button>
              <el-popover v-model="scope.row.visible" placement="top" width="160">
                <p>确定要删除吗？</p>
                <div style="text-align: right; margin-top: 8px;">
                  <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                  <el-button type="primary" @click="handleDelete(scope.row)">确定</el-button>
                </div>
                <template #reference>
                  <el-button type="danger" link icon="delete" @click="scope.row.visible = true">删除</el-button>
                </template>
              </el-popover>
            </template>
          </el-table-column>
        </el-table>
        <div class="gva-pagination">
          <div style="margin-top: 15px;float: left;">
            <el-select v-model="stateOption" value-key="id" class="m-2" placeholder="批量操作" size="large">
                <el-option
                    v-for="item in stateOptions"
                    :key="item.id"
                    :label="item.label"
                    :value="item"
                />
            </el-select>
            <el-button type="primary" @click="toggleSelection()">确定</el-button>
          </div>
          <el-pagination
            :current-page="pageHot"
            :page-size="pageSizeHot"
            :page-sizes="[5, 5, 5, 5]"
            :total.number="+totalHot"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handleCurrentChangeHot"
            @size-change="handleSizeChangeHot"
          />
        </div>
      </div>
      <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="选择商品">
        <el-input
            v-model="productSearchData.name"
            placeholder="商品名称搜索"
            class="input-with-select"
            >
            <template #append>
                <el-button :icon="Search" @click="onSearchProduct" />
            </template>
        </el-input>
        <el-table
          ref="multipleTable"
          :data="productData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="Id"
          @selection-change="handleSelectionProductChange"
        >
          <el-table-column type="selection" width="40" />
          <el-table-column align="left" label="商品编号" prop="id" width="200" />
          <el-table-column align="left" label="商品名称" prop="name" width="200" />
          <el-table-column align="left" label="货号" prop="productSN" width="200" />
          <el-table-column align="left" label="价格" prop="price" width="80" />
        </el-table>
        <el-pagination
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[5, 5, 5, 5]"
            :total.number="+total"
            layout=" prev, pager, next"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
          />
        <el-button  @click="closeDialog()">取消</el-button>
        <el-button type="primary" @click="toggleSelectionProduct()">确定</el-button>
      </el-dialog>
      <el-dialog title="设置排序"
              v-model="sortDialogVisible"
               width="40%">
        <el-form :model="sortDialogData"
                label-width="150px">
          <el-form-item label="排序：">
            <el-input v-model.number="sortDialogData.sort" style="width: 200px"></el-input>
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
  import {
    getProductList,
    getRecommendProductList,
    addRecommendProductList,
    updateRecommendProductByIdForSort,
    updateRecommendProduct,
    deleteRecommendProduct,
  } from '@/api/product'
  import { ref, reactive, onBeforeMount, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ProductStore } from '@/pinia/modules/product'  
  import { Search } from '@element-plus/icons-vue'

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

  const recommendOptions = ref(defaultRecommendOptions)
  const searchData = reactive({
    productName: null,
    recommendStatus: null,
  })
  
  const onSearch = async() => {
    // TODO: 根据广告名称或到期时间搜索
    getTableData()
  }
  
  const onReset = () => {
    searchData.productName = null
    searchData.recommendStatus = null
    getTableData()
  }
  
  const pageHot = ref(1)
  const totalHot = ref(0)
  const pageSizeHot = ref(5)
  const hotProductData = ref([])
  
  // 分页
  const handleSizeChangeHot = (val) => {
    pageSizeHot.value = val
    getTableData()
  }
  
  const handleCurrentChangeHot = (val) => {
    pageHot.value = val
    getTableData()
  }
  
  const productStore = ProductStore()
  // 查询
  const getTableData = async() => {
    const res = await getRecommendProductList({ productName: searchData.productName,
         recommendStatus:searchData.recommendStatus,
        page: pageHot.value, pageSize: pageSizeHot.value })
    if ('code' in res && res.code === 0) {
        hotProductData.value = res.data.list
        totalHot.value = res.data.total
        pageHot.value = res.data.page
        pageSizeHot.value = res.data.pageSize
    }
  }

  getTableData()
  interface stateItem {
    id: number,
    label: string,
    key: string,
    dbKey: string,
    value: number,
  }
  const stateOption = ref<stateItem>()
  const stateOptions = ref<stateItem[]>([])
  stateOptions.value = [
    { id: 0, label: '设为推荐', key: "recommendStatus", dbKey: "recommend_status", value: 1 },
    { id: 1, label: '取消推荐', key: "recommendStatus", dbKey: "recommend_status",  value: 0 },
    { id: 2, label: '删除', key: "delete", dbKey: "delete",  value: 0 },
  ]

  const multipleSelection = ref()
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }
  var updateList:number[] = new Array() 
  const toggleSelection = async() => {
    if (!stateOption.value || !multipleSelection.value || multipleSelection.value.length < 1) {
      ElMessage.warning('请选择一条记录')
      return
    }
    
    multipleSelection.value.forEach((item) => {
        item[stateOption.value.key] = stateOption.value.value
        updateList.push(item.id)
    })
    if (stateOption.value.id < 2) {
        const res = await updateRecommendProduct({ids: updateList, key: stateOption.value.dbKey, value: stateOption.value.value })
        if ('code' in res && res.code == 0) {
            hotProductData.value.forEach(element => {
                element[stateOption.value.key] = stateOption.value.value
            });
        }
    } else if (stateOption.value.id == 2) {
        const res = await deleteRecommendProduct({ids: updateList })
        if ('code' in res && res.code == 0) {
            getTableData()
        }
    } else {
        ElMessage.error('没有此操作')
    }
    
    updateList = []
  }
  
  const type = ref('')
  const dialogFormVisible = ref(false)
  const productForm = ref({
    productId: 0,
    productName: '',
    recommendStatus: 0,
    sort: 0,
  })
  const updateProduct = async(row) => {
    dialogFormVisible.value = true
    type.value = 'update'
    productForm.value = row
  }

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(5)
  const productData = ref([])
  const getProductListTableData = async() => {
    const res = await getProductList({ keyword: productSearchData.name, page: page.value, pageSize: pageSize.value })
    if ('code' in res && res.code === 0) {
        productData.value = []
        res.data.list.forEach(element => {
          let isExit = false
          hotProductData.value.forEach(item => {
            if (item.productId == element.id) {
              isExit = true
            }
          })
          if (!isExit) {
            element.productSN = "NO." + element.productSN
            element.price = "￥" + element.price
            productData.value.push(element)
          }
            
        });
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.pageSize
    }
  }
   // 分页
   const handleSizeChange = (val) => {
    pageSize.value = val
    getProductListTableData()
  }
  
  const handleCurrentChange = (val) => {
    page.value = val
    getProductListTableData()
  }

  const openDialog = () => {
    dialogFormVisible.value = true
    type.value = 'create'
    getProductListTableData()
  }

  const closeDialog = () => {
    dialogFormVisible.value = false
  }

  const productSearchData = reactive({
    name: null,
  })

  const onSearchProduct = async() => {
    getProductListTableData()
  }

  const multipleProductSelection = ref()
  const handleSelectionProductChange = (val) => {
    multipleProductSelection.value = val
  }
  var updateProductList:object[] = new Array() 
  const toggleSelectionProduct = async() => {
    if (!multipleProductSelection.value || multipleProductSelection.value.length < 1) {
      ElMessage.warning('请选择一条记录')
      return
    }
    multipleProductSelection.value.forEach((item) => {
        let product = {productId: item.id, productName: item.name}
        updateProductList.push(product)
    })
    
    const res = await addRecommendProductList({products: updateProductList })
    if ('code' in res && res.code == 0) {
        getTableData()
        closeDialog()
        ElMessage.success("添加成功！")
    }
    updateProductList = []
  }

  const sortDialogVisible = ref(false)
  const sortDialogData = ref({sort: 0, id: null})
  const handleEditSort = async (index, row) => {
    sortDialogVisible.value = true;
    sortDialogData.value.id = row.id
    sortDialogData.value.sort = row.sort
  }

  const handleUpdateSort = async () => {
      ElMessageBox.confirm('是否要修改排序?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async() => {
        updateRecommendProductByIdForSort(sortDialogData.value).then(response=>{
                sortDialogVisible.value =false;
                getTableData();
                ElMessage({
                  type: 'success',
                  message: '删除成功!',
                })
              });
        }).catch(() => {
          ElMessage({
            type: 'info',
            message: '已取消删除',
          })
      })
  }

  const handleDelete = async (row) => {
    deleteProduct([row.id])
  }
  const deleteProduct = async(ids) => {
      ElMessageBox.confirm('是否要删除该推荐?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async() => {
        deleteRecommendProduct({"ids": ids}).then(response=>{
            getTableData();
            ElMessage({
              type: 'success',
              message: '删除成功!',
            })
          });
        }).catch(() => {
          ElMessage({
            type: 'info',
            message: '已取消删除',
          })
      })  
  }

  const handleRecommendStatusStatusChange = async (inex, row) => {
    const res = await updateRecommendProduct({ids: [row.id], key: "recommend_status", value: row.recommendStatus })
    if ('code' in res && res.code == 0) {
      ElMessage({
        type: 'success',
        message: '推荐成功!',
      })
    }
  }
  </script>
  
  <style scoped>
  .form-item {
    margin: 2px 5px 0 2px;
  }
  .box-card {
    margin: 10px 0 10px 0;
    border: 1px 0 1px 0;
  }

  :deep(.gva-table-box) .el-table .cell {
    white-space: pre-line !important;
  }
  </style>
  