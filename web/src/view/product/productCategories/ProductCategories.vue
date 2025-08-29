<template>
    <div>
      <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">添加</el-button>
        </div>
        <el-table
          ref="multipleTable"
          :data="tableData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="Id"
        >
          <el-table-column align="left" label="编号" prop="id" width="60"></el-table-column>
          <el-table-column align="left" label="分类名称" prop="name" width="100" />
          <el-table-column align="left" label="级别" prop="level" width="60" >
            <template #default="scope">
                <i v-if="scope.row.level==0">一级</i>
                <i v-else>二级</i>
            </template>
          </el-table-column>
          <el-table-column align="left" label="商品数量" prop="productCount" width="100" />
          <el-table-column align="left" label="商品单位" prop="productUnit" width="80"/>
          <el-table-column align="left" label="导航栏" prop="navStatus" width="100" >
            <template #default="scope">
                <el-switch
                    v-model="scope.row.navStatus"
                    :active-value="1"
                    :inactive-value="0"
                    @change="HandleNavStatusRadio(scope.row)">
                </el-switch>
            </template>
          </el-table-column>
          <el-table-column align="left" label="是否显示" prop="showStatus" width="100" >
            <template #default="scope">
                <el-switch
                    v-model="scope.row.showStatus"
                    :active-value="1"
                    :inactive-value="0"
                    @change="HandleShowStatusRadio(scope.row)">
                </el-switch>
            </template>
          </el-table-column>
          <el-table-column align="left" label="排序" prop="sort" width="60"/>
          <el-table-column align="left" label="设置" width="200">
            <template #default="scope">
                <el-button type="primary" link icon="edit" @click="showNext(scope.row)">{{showNextText}}</el-button>
                <el-button type="primary" link icon="edit" @click="editParames(scope.row)">转移商品</el-button>
            </template>
          </el-table-column>
          <el-table-column align="left" label="操作" min-width="160">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="handleUpdateProductCategory(scope.row)">编辑</el-button>
              <el-popover v-model="scope.row.visible" placement="top" width="160">
                <p>确定要删除吗？</p>
                <div style="text-align: right; margin-top: 8px;">
                  <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                  <el-button type="primary" @click="handleDeleteProductCategory(scope.row)">确定</el-button>
                </div>
                <template #reference>
                  <el-button type="danger" link icon="delete" @click="scope.row.visible = true">删除</el-button>
                </template>
              </el-popover>
            </template>
          </el-table-column>
        </el-table>
        <div class="gva-pagination">
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
      <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="类型名称">
        <el-form :model="productForm" label-width="160px">
            <el-form-item label="分类名称">
              <el-input v-model="productForm.name" autocomplete="off" />
            </el-form-item>
            <el-form-item label="上级分类">
                <el-select v-model="productOption" class="m-2" placeholder="请选择品牌" size="large" >
                    <el-option
                    v-for="item in productOptions"
                    :key="item.key"
                    :label="item.value"
                    :value="item"
                    />
                </el-select>
            </el-form-item>
            <el-form-item label="数量单位">
              <el-input v-model="productForm.productUnit" autocomplete="off" />
            </el-form-item>
            <el-form-item label="排序">
              <!-- <el-input v-model="productForm.sort" autocomplete="off" /> -->
              <el-input-number v-model.number="productForm.sort" />
            </el-form-item>
            <el-form-item label="是否显示">
                <el-radio-group v-model="showStatus" class="ml-4" @change="HandleShowStatusRadioChanged">
                    <el-radio :value="1" label="1" size="large">是</el-radio>
                    <el-radio :value="0" label="0" size="large">否</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="是否显示在导航栏">
                <el-radio-group v-model="navStatus" class="ml-4" @change="HandleNavStatusRadioChanged">
                    <el-radio :value="1" label="1" size="large">是</el-radio>
                    <el-radio :value="0" label="0" size="large">否</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="分类图标">
                <el-upload
                    class="product-image-upload"
                    :limit=1
                    :action="`${path}/fileUploadAndDownload/upload`"
                    :headers="{ 'x-token': userStore.token }"
                    :on-preview="handlePreview"
                    :on-remove="handleRemove"
                    :on-error="uploadError"
                    :on-success="uploadIconSuccess"
                    >

                    <div v-if="productForm.icon && productForm.icon.length > 0" class="image-preview-container">
                        <img :src="productForm.icon" class="preview-image" />  
                        <el-icon class="delete-icon" @click.stop="handleRemove"><Delete /></el-icon>
                    </div>
                    <div v-else>
                        <el-icon class="upload-icon"><Plus /></el-icon>
                        <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过10MB</div>
                    </div>
                </el-upload>
            </el-form-item>

            <el-row :gutter="20" v-for="(item, index) in dynamicItem" :key="index">
                <el-col :span="20">
                    <el-form-item label="筛选属性" prop="'price' + index">
                        <el-cascader v-model="item.id" :options="productCategoryOptions" placeholder="请选择" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="4">
                    <el-button @click.prevent="deleteItem(item, index)">删除行</el-button>
                </el-col>
            </el-row>
            <el-form-item>
                <el-button type="primary" size="mini" @click="addItem(dynamicItem.length)">新增</el-button>
            </el-form-item>
            <el-form-item label="关键词">
              <el-input v-model="productForm.keywords" autocomplete="off" />
            </el-form-item>
            <el-form-item label="分类描述">
              <el-input v-model="productForm.description" autocomplete="off" />
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
    getProductCategoryList,
    createProductCategory,
    updateProductCategory,
    deleteProductCategory,
  } from '@/api/product'
  import { ref, reactive, onBeforeMount, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { ProductStore } from '@/pinia/modules/product'
  import config from '@/core/config'
  import { trim } from '@/utils/stringFun'
  import { number } from 'echarts'
  import type { UploadProps, UploadFile, UploadInstance  } from 'element-plus'
  import { useUserStore } from '@/pinia/modules/user'
  const userStore = useUserStore()
  const productStore = ProductStore()

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(5)
  const tableData = ref([])
  const currentLevel = ref(0)
  const isDisabled = ref(false)
  const showNextText = ref('查看下级')
  const maxLevel = ref(2) // 最大级别
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
    const res = await getProductCategoryList({ tag: currentLevel.value, page: page.value, pageSize: pageSize.value })
    if ('code' in res && res.code === 0) {
        tableData.value = res.data.list
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.pageSize
    }
  }

  getTableData()
  
  const productCategoryOptions = ref([])
  const openDialog = async() => {
    type.value = 'create'
    dialogFormVisible.value = true
    await getAllTableData()
    await productStore.BuildProductAttributeData()
    productCategoryOptions.value = productStore.ProductCategoryOptions
  }

  const type = ref('')
  const dialogFormVisible = ref(false)
  const productForm = ref({
    parentId: 0,
    name: '',
    level: 0,
    productCount: 0,
    productUnit: '',
    navStatus: 0,
    showStatus: 0,
    sort: 0,
    icon: "",
    keywords: '',
    description: '',
  })
  const handleUpdateProductCategory = async(row) => {
    productForm.value = row
    console.log("-----productForm.value-----", productForm.value)
    dialogFormVisible.value = true
    type.value = 'update'
    showStatus.value = String(row.showStatus)
    navStatus.value = String(row.navStatus)
  }

  const closeDialog = () => {
    dialogFormVisible.value = false
    productForm.value = {
        parentId: 0,
        name: '',
        level: 0,
        productCount: 0,
        productUnit: '',
        navStatus: 0,
        showStatus: 0,
        sort: 0,
        icon: '',
        keywords: '',
        description: '',
    }
    dynamicItem.value = [{
        id: 0
    }]
  }
  const enterDialog = async() => {
    
    let res
    switch (type.value) {
        case 'create':
            if (productForm.value.parentId != 0) {
              productForm.value.level = 1
            }
            res = await createProductCategory(productForm.value)
            break
        case 'update':
            res = await updateProductCategory(productForm.value)
            break
        default:
            res = await createProductCategory(productForm.value)
            break
    }

    if (res.code === 0) {
        closeDialog()
        getTableData()
    }
  }

  const handleDeleteProductCategory = async(row) => {
    let res = await deleteProductCategory(row)
    if ("code" in res && res.code === 0) {
        getTableData()
    }
  }

  const showNext = (row) => {
    maxLevel.value = maxLevel.value - 1
    if (maxLevel.value == 1) {
      currentLevel.value = row.id
      showNextText.value = '查看上级'
      // isDisabled.value = false
      getTableData()
    } else if (maxLevel.value == 0) {
      maxLevel.value = 2
      currentLevel.value = 0
      showNextText.value = '查看下级'
      // isDisabled.value = true
      getTableData()
    } else {
      console.log(" 级别选择错误，当前级别：", currentLevel.value)
    }

  }

  const editParames = (row) => {
  }

  interface productItem {
    key: number,
    value: string,
  }
  const productOption = ref<productItem>({key: 0, value: "无上级分类"})
  const productOptions = ref<productItem[]>([])

  const getAllTableData = async() => {
    const res = await getProductCategoryList({ tag: 0, page: 1, pageSize: 100 })
    if ('code' in res && res.code === 0) {
        productOptions.value.push({key: 0, value: "无上级分类"})
        res.data.list.forEach((item)=>{
            productOptions.value.push({key: item.id, value: item.name})
        })
    }
  }
  watch(() => productOption.value, () => {
    productForm.value.parentId = productOption.value.key
  })

  const showStatus = ref("0")
  const navStatus = ref("0")
   const HandleNavStatusRadio = async (event) => {
    console.log("NavStatus.value:", event)
    let res = await updateProductCategory(event)
    if ("code" in res && res.code === 0) {
        getTableData()
        ElMessage({
            type: 'success',
            message: '更新成功'
        })
    } else {
        ElMessage({
            type: 'error',
            message: '更新失败'
        })
    }
  }

  const HandleShowStatusRadio = async (event) => {
    console.log("ShowStatus.value:", event)
    let res = await updateProductCategory(event)
    if ("code" in res && res.code === 0) {
        getTableData()
        ElMessage({
            type: 'success',
            message: '更新成功'
        })
    } else {
        ElMessage({
            type: 'error',
            message: '更新失败'
        })
    }
  }

  const HandleShowStatusRadioChanged = () => {
    productForm.value.showStatus = Number(showStatus.value)
    // updateProductCategory(productForm.value)
  }
  const HandleNavStatusRadioChanged = async () => {
    productForm.value.navStatus = Number(navStatus.value)
  }

  const uploadError = () => {
    ElMessage({
        type: 'error',
        message: '图片上传失败'
    })
  }
  
  const path = ref(import.meta.env.VITE_BASE_API)
  const uploadIconSuccess = (res) => {
    const { data } = res
    if (data.file) {
        productForm.value.icon = data.file.url
        console.log("-----bigPic.url-----", data.file.url)
    }
  }

  const handleRemove = (uploadFile) => {
    console.log(uploadFile)
    productForm.value.icon = ""
  }
  const handlePreview: UploadProps['onPreview'] = (file) => {
    console.log(file)
  }

  const disabled = ref(false)

  const dynamicItem = ref([{
        id: 0
  }])
  const addItem = (length) =>{
    if (length >= 10) {

    } else {
        dynamicItem.value.push({
            id: 0
        })
    }
  }
  const deleteItem = (item, index) =>{
    // const index = dynamicItem.value.indexOf(item)
    if (index >= 0) {
        dynamicItem.value.splice(index, 1)
    } 
  }
  watch(() => dynamicItem.value, () => {
    console.log("--------dynamicItem:", dynamicItem)
  })
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

  .product-image-upload {
  display: inline-block;
}

.upload-icon {
  width: 100px;
  height: 100px;
  line-height: 100px;
  text-align: center;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  font-size: 24px;
  color: #8c8c8c;
  transition: all 0.3s;
}

.upload-icon:hover {
  color: #409eff;
  border-color: #409eff;
}

.image-preview-container {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: all 0.3s;
}

.delete-icon {
  position: absolute;
  top: 0;
  right: 0;
  width: 24px;
  height: 24px;
  background-color: rgba(0, 0, 0, 0.5);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: all 0.3s;
}

.image-preview-container:hover .delete-icon {
  opacity: 1;
}

.image-preview-container:hover .preview-image {
  filter: brightness(0.8);
}
  </style>
  