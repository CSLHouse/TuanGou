<template>
    <div>
      <div>
        <el-card class="box-card" title="筛选搜索">
          <el-form :inline="true" :model="searchData" class="demo-form-inline">
            <el-form-item label="品牌名称：" class="form-item">
              <el-input v-model="searchData.name" placeholder="品牌名称" clearable />
            </el-form-item>
            <el-form-item label="推荐状态：" class="form-item">
                <el-select v-model="searchData.showStatus" value-key="id"
                    placeholder="请选择" size="default">
                    <el-option
                        v-for="item in recommendStatusOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                    />
                </el-select>
            </el-form-item>
          </el-form>
          <el-button type="default" @click="handleReset">重置</el-button>
          <el-button type="primary" @click="onSearch">搜索结果</el-button>
        </el-card>
      </div>
      <div class="gva-table-box">
        <el-table
          ref="multipleTable"
          :data="productData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="Id"
          @selection-change="handleSelectionChange"
          border
        >
          <el-table-column type="selection" width="40" />
          <el-table-column align="left" label="编号" prop="id" width="60"></el-table-column>
          <el-table-column align="left" label="品牌名称" prop="name" width="140" />
          <el-table-column align="left" label="是否推荐" prop="type" width="100" >
            <template #default="scope">
                <el-switch
                    v-model="scope.row.showStatus"
                    :active-value="1"
                    :inactive-value="0"
                    @change="handleChangeOnlineState(scope.row)">
                </el-switch>
            </template>
          </el-table-column>
          <el-table-column align="left" label="排序" prop="sort" width="100" />
          <el-table-column align="left" label="相关" prop="about" width="200" >
            <template #default="scope">
                商品：
                <span style="color: #409EFF;">{{ scope.row.productCount}}</span>
                评价：
                <span style="color: #409EFF;">{{ scope.row.productCommentCount}}</span>
            </template>
          </el-table-column>  
       
          <el-table-column align="left" label="状态" prop="status" width="80" />
          <el-table-column align="left" label="操作" min-width="160">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="updateProduct(scope.row)">设置排序</el-button>
              <el-popover v-model="scope.row.visible" placement="top" width="160">
                <p>确定要删除吗？</p>
                <div style="text-align: right; margin-top: 8px;">
                  <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                  <el-button type="primary" @click="deleteBrand(scope.row)">确定</el-button>
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
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[5, 5, 5, 5]"
            :total.number="+total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
          />
        </div>
      </div>

      <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="编辑品牌">
        <el-form :model="brandForm" label-width="140px" >
            <el-form-item label="排序">
                <el-input v-model.number="brandForm.sort" autocomplete="off" />
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
            <el-button @click="closeDialog">取 消</el-button>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            </div>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script setup lang="ts">
  import {
    getBrandList,
    updateProductBrand,
    deleteProductBrand,
    updateBrandByIdForState,
  } from '@/api/product'
  import { ref, reactive, onBeforeMount, watch } from 'vue'
  import { ElMessage } from 'element-plus'

  const searchData = reactive({
    name: null,
    showStatus: null,
  })
  
  const onSearch = async() => {
    getTableData()
  }
  
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(5)
  const productData = ref([])
  
  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }
  
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }
  
  // 查询
  const getTableData = async() => {
    const res = await getBrandList({ name: searchData.name, 
        showStatus: searchData.showStatus,
        page: page.value, pageSize: pageSize.value })
    if ('code' in res && res.code === 0) {
        productData.value = res.data.list
        productData.value.forEach(element => {
            if (element.showStatus == 0) {
                element.status = '未推荐'
            } else {
                element.status = '推荐中'
            }
        });
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.pageSize
    }
  }

  getTableData()

  const recommendStatusOptions = ref([
    {label: "未推荐", value: 100},
    {label: "推荐中", value: 101}
  ])

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
    { id: 0, label: '设为推荐', key: "showStatus", dbKey: "show_status", value: 1 },
    { id: 1, label: '取消推荐', key: "showStatus", dbKey: "show_status",  value: 0 },
  ]

  const multipleSelection = ref()
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }
  var updateList:number[] = new Array() 
  const toggleSelection = async() => {
    if (!multipleSelection.value || multipleSelection.value.length < 1) return
    multipleSelection.value.forEach((item) => {
        item[stateOption.value.key] = stateOption.value.value
        updateList.push(item.id)
    })
    const res = await updateBrandByIdForState({ids: updateList, key: stateOption.value.dbKey, value: stateOption.value.value })
    if ('code' in res && res.code !== 0) {
        productData.value.forEach(element => {
            element[stateOption.value.key] = stateOption.value.value
        });
        ElMessage({
          type: 'success',
          message: '更新成功!',
        })
    }
    updateList = []
  }
  
  const type = ref('')
  const dialogFormVisible = ref(false)
  const brandForm = ref({
    id: null,
    brandId: 1,
    name: '',
    showStatus: '',
    productCount: 0,
    productCommentCount: 0,
    sort: '',
    status: '',
  })

  const updateProduct = async(row) => {
    dialogFormVisible.value = true
    brandForm.value = row
  }


  const enterDialog = async() => {
    let res = await updateProductBrand(brandForm.value)
    console.log("----enterDialog---", brandForm.value)
    
    if ("code" in res && res.code === 0) {
        ElMessage({
            type: 'success',
            message: '添加成功'
          })
        closeDialog()
    }
  }
  const closeDialog = () => {
    dialogFormVisible.value = false
  }

  const deleteBrand = async(row) => {
    const res = await deleteProductBrand({id: row.id})
    if ('code' in res && res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功!',
      })
      getTableData()
    }
  }

  const handleChangeOnlineState = async(row) => {
    const res = await updateBrandByIdForState({ids: [row.id], key: 'show_status', value: row.showStatus})
    if ('code' in res && res.code === 0) {
        row.status = row.showStatus == 0 ? "未推荐" : "推荐中"
      ElMessage({
        type: 'success',
        message: '更新成功!',
      })
    }
  }

  const handleReset = () => {
    searchData.name = null
    searchData.showStatus = null
    getTableData()
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
  