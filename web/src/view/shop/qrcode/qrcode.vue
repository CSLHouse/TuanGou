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
        <el-table-column align="center" label="ID" prop="id" width="180" />
          <el-table-column align="center" label="标题" prop="title" width="160"></el-table-column>
          <el-table-column align="center" label="扫描次数" prop="count" width="180" />
          <el-table-column align="center" label="预览" prop="memberStateStr" width="180" >
            <template #default="scope">
                <img :src="scope.row.remoteUrl" style="height: 80px"/>
            </template>
          </el-table-column>
          <el-table-column align="center" label="更新日期" prop="UpdatedAt" width="160" >
            <template #default="scope">
                {{ formatTime(scope.row.UpdatedAt)  }}
            </template>
          </el-table-column>
          <el-table-column align="center" label="状态" prop="UpdatedAt" width="160" >
            <template #default="scope">
                <div v-if="scope.row.isExpired != 0" style="color: red;"> {{ formatIsExpired(scope.row.isExpired) }}</div>
                <div v-else style="color: green;"> {{ formatIsExpired(scope.row.isExpired) }}</div>
            </template>
          </el-table-column>
          <el-table-column align="center" label="操作" min-width="160">
            <template #default="scope">
              <el-popover v-model="scope.row.visible" placement="top" width="160">
                <p>确定要删除吗？</p>
                <div style="text-align: right; margin-top: 8px;">
                  <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                  <el-button type="primary" @click="deleteQrCodeInfo(scope.row)">确定</el-button>
                </div>
                <template #reference>
                  <el-button type="danger" link icon="delete" @click="scope.row.visible = true">删除</el-button>
                </template>
              </el-popover>
              <el-button
                size="small"
                @click="handleShowUpdateQrCode(scope.row)">更新</el-button>
              <el-button
                size="small"
                @click="handleDowloadQrCode(scope.row)">下载</el-button>
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
      <el-dialog title="二维码信息"
            v-model="qrcodeDialogVisible"
               width="40%">
        <el-form :model="qrcodeForm"
                ref="receiverInfoForm"
                label-width="150px">
            <el-form-item label="标题：">
                <el-input v-model="qrcodeForm.title" style="width: 200px"></el-input>
            </el-form-item>
            <el-form-item label="链接：">
                <el-input v-model="qrcodeForm.url" style="width: 200px">
                </el-input>
            </el-form-item>
        </el-form>
        <span slot="footer" class="dialog-footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button type="primary" @click="enterDialog">确 定</el-button>
        </span>
    </el-dialog>
    </div>
  </template>
  

  <script setup lang="ts">
  import {
    getQrCodeList,
    deleteQrCode,
    updateQrCode,
    createQrCode,
    downloadQrCode
  } from '@/api/qrcode.js'
  import { ref, reactive, onBeforeMount, watch, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  import {formatDate} from '@/utils/date';

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const qrcodeDialogVisible = ref(false)
  const operateType = ref('')
  const qrcodeForm = ref({
    title: '',
    url: null,
  })

  onBeforeMount(() => {
    getTableData()
  })

  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }
  
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }
  
  const formatTime = computed(() => {
    return (time: any) => {
        if (time == null || time === '') {
            return 'N/A';
        }
        let date = new Date(time);
        return formatDate(date, 'yyyy-MM-dd HH:mm:ss')
    }
  })

  const formatIsExpired = computed(() => {
    return (value: any) => {
        if (value == 0) {
            return '未过期';
        } else if ( value == 1 ) {
            return '过期';
        } else {
            return '错误';
        }
    }
  })

  // 查询
  const getTableData = async() => {
    const table = await getQrCodeList({ page: page.value, pageSize: pageSize.value })
    if ('code' in table && table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  const deleteQrCodeInfo = async(row) => {
    const res = await deleteQrCode({ id: row.id })
    if ("code" in res && res.code === 0) {
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
  
  const handleShowUpdateQrCode = async(row) => {
    operateType.value = 'update'
    qrcodeDialogVisible.value = true
    qrcodeForm.value = row
  }

  const openDialog = () => {
    operateType.value = 'create'
    qrcodeDialogVisible.value = true
  }

  const closeDialog = () => {
    qrcodeDialogVisible.value = false
    qrcodeForm.value = {
        title: '',
        url: null,
    }
  }

  const enterDialog = async() => {
    if (qrcodeForm.value.title == '' || qrcodeForm.value.url == "") {
        ElMessage({
            type: 'error',
            message: '填写信息有误'
        })
        closeDialog()
        return
        }
    let res
    switch (operateType.value) {
        case 'create':
        res = await createQrCode(qrcodeForm.value)
        break
        case 'update':
        res = await updateQrCode(qrcodeForm.value)
        break
        default:
        res = await createQrCode(qrcodeForm.value)
        break
    }

    if (res.code === 0) {
        closeDialog()
        getTableData()
    }
  }

  const handleDowloadQrCode = async(row) => {
    window.location.href = import.meta.env.VITE_BASE_API + "/qrcode/download?id=" + row.id
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
  </style>
  