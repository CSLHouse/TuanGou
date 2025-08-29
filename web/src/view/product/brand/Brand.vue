<template>
    <div>
      <div>
        <el-card class="box-card" title="筛选搜索">
          <el-form :inline="true" :model="searchData" class="demo-form-inline">
            <el-form-item label="输入搜索：" class="form-item">
              <el-input v-model.number="searchData.name" placeholder="品牌名称/关键字" clearable />
            </el-form-item>
          </el-form>
          <el-button type="primary" @click="onSearch">搜索结果</el-button>
        </el-card>
      </div>
      <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        </div>
        <el-table
          ref="multipleTable"
          :data="productData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="Id"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="40" />
          <el-table-column align="left" label="编号" prop="id" width="60"></el-table-column>
          <el-table-column align="left" label="商品名称" prop="name" width="120" />
          <el-table-column align="left" label="品牌首字母" prop="firstLetter" width="100" />
          <el-table-column align="left" label="排序" prop="sort" width="60" />
          <el-table-column align="left" label="品牌制造商" prop="status" width="100" >
            <template #default="scope">
                <el-switch
                    v-model="scope.row.factoryStatus"
                    :active-value="1"
                    :inactive-value="0">
                </el-switch>
            </template>
          </el-table-column>
          <el-table-column align="left" label="是否显示" prop="status" width="100" >
            <template #default="scope">
                <el-switch
                    v-model="scope.row.showStatus"
                    :active-value="1"
                    :inactive-value="0">
                </el-switch>
            </template>
          </el-table-column>
          
          <el-table-column align="left" label="相关" prop="content" width="180" />
          <el-table-column align="left" label="操作" min-width="160">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="updateProduct(scope.row)">变更</el-button>
              <el-popover v-model="scope.row.visible" placement="top" width="160">
                <p>确定要删除吗？</p>
                <div style="text-align: right; margin-top: 8px;">
                  <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                  <el-button type="primary" @click="handleDeleteBrand(scope.row)">确定</el-button>
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
        <el-form  :model="productForm" label-width="140px" >
            <el-form-item label="品牌名称">
                <el-input v-model="productForm.name" autocomplete="off" />
            </el-form-item>
            <el-form-item label="品牌首字母">
                <el-input v-model.number="productForm.firstLetter" autocomplete="off" />
            </el-form-item>
            <el-form-item label="品牌LOGO">
              <!-- <upload-image
                v-model:imageUrl="productForm.logo"
                :file-size="512"
                :max-w-h="1080"
                class="upload-btn-media-library"
                @on-success="uploadLogoSuccess"
              /> -->
              <SelectImage v-model="productForm.logo" />
            </el-form-item>
            <el-form-item label="品牌专区大图">
              <SelectImage v-model="productForm.bigPic" />
              <!-- <upload-image
                v-model:imageUrl="productForm.bigPic"
                :file-size="512"
                :max-w-h="1080"
                class="upload-btn-media-library"
                @on-success="uploadBigPicSuccess"
              /> -->
            </el-form-item>
            <el-form-item label="品牌故事">
                <el-input v-model="productForm.brandStory" type="textarea" autocomplete="off" />
            </el-form-item>
            <el-form-item label="排序">
                <el-input v-model.number="productForm.sort" autocomplete="off" />
            </el-form-item>
            <el-form-item label="是否显示">
                <el-radio-group v-model="showState" class="ml-4">
                    <el-radio label="1" size="large">是</el-radio>
                    <el-radio label="0" size="large">否</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="品牌制造商">
                <el-radio-group v-model="brandState" class="ml-4">
                    <el-radio label="1" size="large">是</el-radio>
                    <el-radio label="0" size="large">否</el-radio>
                </el-radio-group>
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
    createProductBrand,
    updateProductBrand,
  } from '@/api/product'
  import { ref, reactive, onBeforeMount, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { ProductStore } from '@/pinia/modules/product'  
  import { useUserStore } from '@/pinia/modules/user'
  import SelectImage from '@/components/selectImage/selectImage.vue'
  import { deleteProductBrand } from '@/api/product'

  const userStore = useUserStore()

  const searchData = reactive({
    name: null,
    
  })
  
  const onSearch = async() => {
    // TODO: 根据品牌名或关键词搜索
    getTableData(true)
  }
  
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(5)
  const productData = ref([])
  
  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData(true)
  }
  
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData(true)
  }
  
  const productStore = ProductStore()
  // 查询
  const getTableData = async(refresh) => {
    await productStore.BuildBrandData(refresh)
    productData.value = productStore.RandData['list']
    productData.value.forEach(element => {
      element.content = "商品：" + element.productCount + "  评价：" + element.productCommentCount
    });
    total.value = productStore.RandData['total']
    page.value = productStore.RandData['page']
    pageSize.value = productStore.RandData['pageSize']
  }

  getTableData(true)
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
    { id: 0, label: '显示品牌', key: "showStatus", dbKey: "show_status", value: 1 },
    { id: 1, label: '隐藏品牌', key: "showStatus", dbKey: "show_status",  value: 0 },
  ]

  const multipleSelection = ref()
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }
  var updateList:number[] = new Array() 
  const toggleSelection = async() => {
    // if (!multipleSelection.value || multipleSelection.value.length < 1) return
    // multipleSelection.value.forEach((item) => {
    //     item[stateOption.value.key] = stateOption.value.value
    //     updateList.push(item.id)
    // })
    // const res = await updateProductKeyword({ids: updateList, key: stateOption.value.dbKey, value: stateOption.value.value })
    // if ('code' in res && res.code !== 0) {
    //     productData.value.forEach(element => {
    //         element[stateOption.value.key] = stateOption.value.value
    //     });
    // }
    // updateList = []
  }
  
  const type = ref('')
  const dialogFormVisible = ref(false)
  const productForm = ref({
    name: '',
    firstLetter: '',
    sort: 0,
    factoryStatus: 0,
    showStatus: 0,
    logo: '',
    bigPic: '',
    brandStory: '',
  })
  const updateProduct = async(row) => {
    dialogFormVisible.value = true
    type.value = 'update'
    productForm.value = row
  }
  const openDialog = () => {
    dialogFormVisible.value = true
    type.value = 'create'
    resetForm()
  }

  const enterDialog = async() => {
    let res
    switch (type.value) {
        case 'create':
            productForm.value.showStatus = Number(showState.value)
            productForm.value.factoryStatus = Number(brandState.value)
            res = await createProductBrand(productForm.value)
            break
        case 'update':
            res = await updateProductBrand(productForm.value)
            break
        default:
            res = await createProductBrand(productForm.value)
            break
    }

    if (res.code === 0) {
        closeDialog()
        await productStore.BuildBrandData(true)
        getTableData(true)
    }
  }
  const closeDialog = () => {
    dialogFormVisible.value = false
  }

  const showState = ref("1")
  const brandState = ref("1")

  const handleDeleteBrand = async(row) => {
    const res = await deleteProductBrand({id: row.id})
    if ('code' in res && res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功!',
      })
      getTableData(true)
    }
  }
  const resetForm = () => {
    productForm.value = {
      name: '',
      firstLetter: '',
      sort: 0,
      factoryStatus: 0,
      showStatus: 0,
      logo: '',
      bigPic: '',
      brandStory: '',
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
  