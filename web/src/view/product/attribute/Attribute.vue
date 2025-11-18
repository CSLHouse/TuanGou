<template>
    <div>
      <div class="gva-table-box">
        <div class="gva-btn-list">
          <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        </div>
        <el-table
          ref="multipleTable"
          :data="tableData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="Id"
        >
          <el-table-column align="left" label="编号" prop="id" width="60"/>
          <el-table-column align="left" label="属性名称" prop="name" width="100"/>
          <el-table-column align="left" label="商品类型" prop="productType" width="100"/>
          <el-table-column align="left" label="属性是否可选" prop="selectType" width="120" />
          <el-table-column align="left" label="属性值的录入方式" prop="inputType" min-width="100"/>
          <el-table-column align="left" label="可选值列表" prop="inputList" width="220" />
          <el-table-column align="left" label="排序" prop="sort" width="80" />
          <el-table-column align="left" label="操作" min-width="160">
            <template #default="scope">
              <el-button type="primary" link icon="edit" @click="updateAttribute(scope.row)">编辑</el-button>
              <el-popover v-model="scope.row.visible" placement="top" width="160">
                <p>确定要删除吗？</p>
                <div style="text-align: right; margin-top: 8px;">
                  <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                  <el-button type="primary" @click="deleteAttribute(scope.row)">确定</el-button>
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
            :page-sizes="[10, 30, 50, 100]"
            :total.number="+total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
          />
        </div>
      </div>
      <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="商品属性">
        <el-form :model="attributeForm" label-width="180px">
          <el-form-item label="属性名称">
            <el-input v-model="attributeForm.name" autocomplete="off" />
          </el-form-item>
          <el-form-item label="商品类型">
            <el-input v-model="attributeForm.productType" autocomplete="off" disabled/>
          </el-form-item>
          <el-form-item label="分类筛选样式">
            <el-radio-group v-model="filterTypeState" class="ml-4" @change="HandleFilterTypeStateRadioChanged">
                <el-radio label="0" size="large">普通</el-radio>
                <el-radio label="1" size="large">颜色</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="能否进行检索">
            <el-radio-group v-model="searchTypeState" class="ml-4" @change="HandleSearchTypeStateRadioChanged">
                <el-radio label="0" size="large">不需要检索</el-radio>
                <el-radio label="1" size="large">关键字检索</el-radio>
                <el-radio label="2" size="large">范围检索</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="商品属性关联">
            <el-radio-group v-model="relatedState" class="ml-4" @change="HandleRelatedStateRadioChanged">
                <el-radio label="1" size="large">是</el-radio>
                <el-radio label="0" size="large">否</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="属性是否可选">
            <el-radio-group v-model="selectTypeState" class="ml-4" @change="HandleSelectTypeStateRadioChanged">
                <el-radio label="0" size="large">唯一</el-radio>
                <el-radio label="1" size="large">单选</el-radio>
                <el-radio label="2" size="large">复选</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="属性值的录入方式">
            <el-radio-group v-model="inputTypeState" class="ml-4" @change="HandleInputTypeStateRadioChanged">
                <el-radio label="0" size="large">手工录入</el-radio>
                <el-radio label="1" size="large">从下面列表种选择</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="属性值可选值列表">
            <el-input v-model="attributeForm.inputList" :autosize="true" type="textarea"/>
          </el-form-item>
          <el-form-item label="是否支持手动新增">
            <el-radio-group v-model="handAddState" class="ml-4" @change="HandleHandAddStateRadioChanged">
                <el-radio label="1" size="large">是</el-radio>
                <el-radio label="0" size="large">否</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="排序属性" >
            <el-input v-model="attributeForm.sort" autocomplete="off" type="number"/>
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="closeDialog">取消</el-button>
            <el-button type="primary" @click="enterDialog">提交</el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script setup>
  import {
    createProductAttribute,
    updateProductAttribute,
    deleteProductAttribute,
    getProductAttributeList
  } from '@/api/product'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { ref, onBeforeMount } from 'vue'
  import { ElMessage } from 'element-plus'
  import { formatDate } from '@/utils/format'
  import { useRoute } from 'vue-router'
import { number } from 'echarts'

  const route = useRoute()

  const attributeForm = ref({
    name: '',
    selectType: 0,
    inputType: 0,
    inputList: '',
    sort: 0,
    filterType: 0,
    searchType: 0,
    relatedStatus: 0,
    handAddStatus: 0,
    type: 0,
  })
  
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  
  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }
  
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }
  
  const categoriesId = route.query.cid
  const attributeTypeName = route.query.cname
  const attributeType = route.query.ctype
  // onBeforeMount(() => {

  // })

  // 查询
  const getTableData = async() => {
    // let id = route.query.cid
    // let typeName = route.query.cname
    // let type = route.query.ctype
    const table = await getProductAttributeList({ tag: categoriesId, page: page.value, pageSize: pageSize.value, state: attributeType })
    if (table.code === 0) {
      tableData.value = table.data.list

      tableData.value.forEach((element)=>{
        element.productType = attributeTypeName
      })
    }
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
  }

  getTableData()
  
  const dialogFormVisible = ref(false)
  const type = ref('')
  const updateAttribute = async(row) => {
    type.value = 'update'
    attributeForm.value = row
    dialogFormVisible.value = true

    filterTypeState.value = String(row.filterType)
    searchTypeState.value = String(row.searchType)
    relatedState.value = String(row.relatedStatus)
    selectTypeState.value = String(row.selectType)
    inputTypeState.value = String(row.inputType)
    handAddState.value = String(row.handAddStatus)
  }
  const closeDialog = () => {
    dialogFormVisible.value = false
    attributeForm.value = {
        name: '',
        selectType: 0,
        inputType: 0,
        inputList: '',
        sort: 0,
        filterType: 0,
        searchType: 0,
        relatedStatus: 0,
        handAddStatus: 0,
        type: 0,
        productAttributeCategoryId: 0,
    }
  }
  const deleteAttribute = async(row) => {
    row.visible = false
    const res = await deleteProductAttribute({Id: row.id })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  }
  const enterDialog = async() => {
    let res
    attributeForm.value.type = Number(attributeType)
    attributeForm.value.sort = Number(attributeForm.value.sort)

    switch (type.value) {
      case 'create':
        res = await createProductAttribute(attributeForm.value)
        break
      case 'update':
        res = await updateProductAttribute(attributeForm.value)
        break
      default:
        res = await createProductAttribute(attributeForm.value)
        break
    }
  
    if (res.code === 0) {
      closeDialog()
      getTableData()
    }
  }
  const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
    attributeForm.value.productType = attributeTypeName
    attributeForm.value.productAttributeCategoryId = Number(categoriesId)
  }
  
  const filterTypeState = ref("0")
  const searchTypeState = ref("0")
  const relatedState = ref("0")
  const selectTypeState = ref("0")
  const inputTypeState = ref("0")
  const handAddState = ref("0")

  const HandleFilterTypeStateRadioChanged = () => {
    attributeForm.value.filterType = Number(filterTypeState.value)
  }
  const HandleSearchTypeStateRadioChanged = () => {
    attributeForm.value.searchType = Number(searchTypeState.value)
  }
  const HandleRelatedStateRadioChanged = () => {
    attributeForm.value.relatedStatus = Number(relatedState.value)
  }
  const HandleSelectTypeStateRadioChanged = () => {
    attributeForm.value.selectType = Number(selectTypeState.value)
  }
  const HandleInputTypeStateRadioChanged = () => {
    attributeForm.value.inputType = Number(inputTypeState.value)
  }
  const HandleHandAddStateRadioChanged = () => {
    attributeForm.value.handAddStatus = Number(handAddState.value)
  }

  </script>
  
  <script>
  
  export default {
    name: 'Attribute'
  }
  </script>
  
  <style></style>
  