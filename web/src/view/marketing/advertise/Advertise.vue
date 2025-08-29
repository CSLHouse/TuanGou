<template>
    <div>
      <div>
        <el-card class="box-card" title="筛选搜索">
          <el-form :inline="true" :model="searchData" class="demo-form-inline">
            <el-form-item label="广告名称：" class="form-item">
              <el-input v-model.number="searchData.name" placeholder="广告名称" clearable />
            </el-form-item>
            <el-form-item label="广告位置：" class="form-item">
              <el-input v-model="searchData.type" placeholder="广告位置" clearable />
            </el-form-item>
            <el-form-item label="到期时间：" class="form-item">
              <el-input v-model="searchData.endTime" placeholder="到期时间：" clearable />
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
          <el-table-column align="left" label="广告名称" prop="name" width="120" />
          <el-table-column align="left" label="广告位置" prop="type" width="100" />
          <el-table-column align="left" label="广告图片" prop="pic" width="100" >
            <template #default="scope">
                <img  :src="scope.row.pic" alt="" style="width: 50px;height: 50px">
            </template>
          </el-table-column>
          <el-table-column align="left" label="时间" prop="dateTime" width="240" />
          <el-table-column align="left" label="上线/下线" prop="state" width="100" >
            <template #default="scope">
                <el-switch
                    v-model="scope.row.state"
                    :active-value="1"
                    :inactive-value="0"
                    @change="handleChangeOnlineState(scope.row)">
                </el-switch>
            </template>
          </el-table-column>
          <el-table-column align="left" label="点击次数" prop="clickCount" width="80" />
          <el-table-column align="left" label="生成订单" prop="orderCount" width="80" />
          <el-table-column align="left" label="操作" min-width="160">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="updateProduct(scope.row)">变更</el-button>
              <el-popover v-model="scope.row.visible" placement="top" width="160">
                <p>确定要删除吗？</p>
                <div style="text-align: right; margin-top: 8px;">
                  <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                  <el-button type="primary" @click="deleteProduct(scope.row)">确定</el-button>
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
            <el-form-item label="广告名称">
                <el-input v-model="productForm.name" autocomplete="off" />
            </el-form-item>
            <el-form-item label="广告位置">
                <el-input v-model.number="productForm.type" autocomplete="off" />
            </el-form-item>
            <el-form-item label="开始时间">
                <div class="demo-date-picker">
                    <el-date-picker
                        v-model="productForm.startTime"
                        type="datetime"
                        placeholder="选择时间"
                    >
                    </el-date-picker>
                </div>
            </el-form-item>
            <el-form-item label="结束时间">
                <div class="demo-date-picker">
                    <el-date-picker
                        v-model="productForm.endTime"
                        type="datetime"
                        placeholder="选择时间"
                    >
                    </el-date-picker>
                </div>
            </el-form-item>
            <el-form-item label="上线/下线">
                <el-radio-group v-model="onlineState" class="ml-4" @change="HandleonlineRadioChanged">
                    <el-radio :label="0" size="large">下线</el-radio>
                    <el-radio :label="1" size="large">上线</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="广告图片">
              <cooller-single-upload :listType="'picture'" v-model="productForm.pic" ></cooller-single-upload>
            </el-form-item>
            <el-form-item label="排序">
                <el-input v-model.number="productForm.sort" autocomplete="off" />
            </el-form-item>
            <el-alert type="info" show-icon :closable="false">
              <p>
                指向品牌详情：/subpages/brand/brandDetail?id=品牌ID，品牌ID到【商品】-【品牌管理】里面找；
              </p>
              <p>
                指向商品详情：/subpages/product/product?id=商品ID，商品ID到【商品】-【商品列表】里面找；
              </p>
            </el-alert>
            <el-form-item label="广告链接">
                <el-input v-model="productForm.url" autocomplete="off" />
            </el-form-item>
            <el-form-item label="广告备注">
                <el-input v-model="productForm.note" autocomplete="off" />
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
    getAdvertiseList,
    createAdvertise,
    updateAdvertise,
    deletedvertise,
    updateAdvertiseByIdForState
  } from '@/api/product'
  import { ref, reactive, onBeforeMount, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { ProductStore } from '@/pinia/modules/product'  
  import coollerSingleUpload from '@/components/upload/coollerSingleUpload.vue'

  const searchData = reactive({
    name: null,
    type: null,
    endTime: null,
  })
  
  const onSearch = async() => {
    // TODO: 根据广告名称或到期时间搜索
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
  
  const productStore = ProductStore()
  // 查询
  const getTableData = async() => {
    const res = await getAdvertiseList({ name: searchData.name, type:searchData.type,
        endTime:searchData.endTime,
        page: page.value, pageSize: pageSize.value })
    if ('code' in res && res.code === 0) {
        productData.value = res.data.list
        productData.value.forEach(element => {
            element.dateTime = "开始时间：" + element.startTime + "\n 结束时间：" + element.endTime
        });
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.pageSize
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
    { id: 0, label: '显示品牌', key: "state", dbKey: "state", value: 1 },
    { id: 1, label: '隐藏品牌', key: "state", dbKey: "state",  value: 0 },
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
    const res = await updateAdvertiseByIdForState({ids: updateList, key: stateOption.value.dbKey, value: stateOption.value.value })
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
  const productForm = ref({
    name: '',
    type: 1,
    pic: '',
    startTime: '',
    endTime: '',
    state: 0,
    clickCount: 0,
    orderCount: 0,
    url: '',
    note: '',
    sort: 0,
  })
  const updateProduct = async(row) => {
    dialogFormVisible.value = true
    type.value = 'update'
    productForm.value = row
    onlineState = row.state
    // let picName = getImageUrlName(row.pic)
    // fileList.value = []
    // fileList.value.push({name: picName, url: row.pic})
  }
  const openDialog = () => {
    dialogFormVisible.value = true
    type.value = 'create'
  }

  const enterDialog = async() => {
    let res
    console.log("----enterDialog---", productForm.value)
    
    switch (type.value) {
        case 'create':
            res = await createAdvertise(productForm.value)
            break
        case 'update':
            res = await updateAdvertise(productForm.value)
            break
        default:
            res = await createAdvertise(productForm.value)
            break
    }

    if (res.code === 0) {
        ElMessage({
            type: 'success',
            message: '操作成功'
          })
        closeDialog()
        await productStore.BuildBrandData(true)
        getTableData()
    }
  }
  const closeDialog = () => {
    dialogFormVisible.value = false
  }

  const deleteProduct = async(row) => {
    console.log("=====row:", row)
    const res = await deletedvertise({id: row.id})
    if ('code' in res && res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功!',
      })
      getTableData()
    }
  }
  let onlineState = 1
  const HandleonlineRadioChanged = () => {
    productForm.value.state = onlineState
  }
  
  const handleChangeOnlineState = async(row) => {
    console.log("--------row:", row)
    const res = await updateAdvertiseByIdForState({ids: [row.id], key: 'state', value: row.state})
    if ('code' in res && res.code === 0) {
      ElMessage({
        type: 'success',
        message: '更新成功!',
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
  