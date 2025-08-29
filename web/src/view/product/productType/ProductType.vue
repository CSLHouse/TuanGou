<template>
    <div>
      <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">添加</el-button>
        </div>
        <el-table
          ref="multipleTable"
          :data="attributeData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="Id"
        >
          <el-table-column align="left" label="编号" prop="id" width="60"></el-table-column>
          <el-table-column align="left" label="类型名称" prop="name" width="180" />
          <el-table-column align="left" label="属性数量" prop="attributeCount" width="160" />
          <el-table-column align="left" label="参数数量" prop="paramCount" width="160" />
          <el-table-column align="left" label="设置" prop="ID" width="260">
            <template #default="scope">
                <el-button type="primary" link icon="edit" @click="editAttributes(scope.row)">属性列表</el-button>
                <el-button type="primary" link icon="edit" @click="editParames(scope.row)">参数列表</el-button>
            </template>
          </el-table-column>
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
            :page-sizes="[5, 5, 5, 5]"
            :total.number="+total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
          />
        </div>
      </div>
      <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="类型名称">
        <el-form :inline="true" :model="attributeForm" label-width="80px">
            <el-form-item label="类型名称">
              <el-input v-model="attributeForm.name" autocomplete="off" />
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
    getProductAttributeCategoryList,
    createProductAttributeCategory,
    updateProductAttributeCategory,
    deleteProductAttributeCategory,
  } from '@/api/product'
  import { ref, reactive, onBeforeMount, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { comboStore } from '@/pinia/modules/combo'
  import config from '@/core/config'
  import { trim } from '@/utils/stringFun'
  import { number } from 'echarts'
  import { useRouter } from "vue-router";
  const router = useRouter()

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(5)
  const attributeData = ref([])
  
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
    const res = await getProductAttributeCategoryList({ page: page.value, pageSize: pageSize.value })
    if ('code' in res && res.code === 0) {
      console.log(res)

        attributeData.value = res.data.list
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.pageSize
    }
  }

  getTableData()
  
  const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
  }

  const type = ref('')
  const dialogFormVisible = ref(false)
  const attributeForm = ref({
    name: '',
  })
  const updateAttribute = async(row) => {
    attributeForm.value = row
    dialogFormVisible.value = true
    type.value = 'update'
  }

  const closeDialog = () => {
    dialogFormVisible.value = false
    attributeForm.value = {
        name: '',
    }
  }
  const enterDialog = async() => {
    let res
    switch (type.value) {
        case 'create':
            res = await createProductAttributeCategory(attributeForm.value)
            break
        case 'update':
            res = await updateProductAttributeCategory(attributeForm.value)
            break
        default:
            res = await createProductAttributeCategory(attributeForm.value)
            break
    }

    if (res.code === 0) {
        closeDialog()
        getTableData()
    }
  }

  const deleteAttribute = async(row) => {
    let res = await deleteProductAttributeCategory(row)
    if ("code" in res && res.code === 0) {
        getTableData()
    }
  }

  const editAttributes = (row) => {
    router.push({ path: 'Attribute', query: { cid: row.id, cname: row.name, ctype: 0 }})
  }

  const editParames = (row) => {
    router.push({ path: 'Attribute', query: { cid: row.id, cname: row.name, ctype: 1 }})
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
  