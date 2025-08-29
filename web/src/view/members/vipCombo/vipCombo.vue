<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
      </div>
      <el-table
        ref="multipleTable"
        :data="comboData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="Id"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="套餐编号" prop="id"></el-table-column>
        <el-table-column align="left" label="套餐名称" prop="comboName" width="120" />
        <el-table-column align="left" label="套餐类型" prop="comboType" width="120" />
        <el-table-column align="left" label="套餐价格" prop="comboPrice" width="120" />
        <el-table-column align="left" label="天数/次数/金额" prop="times" width="120" />
        <el-table-column align="left" label="套餐状态" prop="stateStr" width="120" />
        <el-table-column align="left" label="操作" min-width="160">
          <template #default="scope">
            <el-button type="primary" link icon="edit" @click="updateCombo(scope.row)">变更</el-button>
            <el-popover v-model="scope.row.visible" placement="top" width="160">
              <p>确定要删除吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                <el-button type="primary" @click="deleteCombo(scope.row)">确定</el-button>
              </div>
              <template #reference>
                <el-button type="primary" link icon="delete" @click="scope.row.visible = true">删除</el-button>
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
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="Vip套餐">
      <el-form  :model="comboForm" label-width="140px" >
        <el-form-item label="会员套餐名称">
          <el-input v-model="comboForm.comboName" autocomplete="off" />
        </el-form-item>
        <el-form-item label="会员套餐类型">
          <el-select v-model="comboTypeValue" value-key="id" class="m-2" placeholder="请选择会员卡" size="large">
              <el-option
                v-for="item in config.comboTypeOptions"
                :key="item.id"
                :label="item.label"
                :value="item"
              />
          </el-select>
        </el-form-item>
        <el-form-item label="价格">
          <el-input v-model.number="comboForm.comboPrice" autocomplete="off" />
        </el-form-item>
        <el-form-item label="天数/次数/金额">
          <el-input v-model.number="comboForm.times" autocomplete="off" />
        </el-form-item>
        <el-form-item label="套餐状态">
          <el-select v-model="stateValue" value-key="id" class="m-2" placeholder="请选择会员卡" size="large">
            <el-option
            v-for="item in stateOptions"
            :key="item.id"
            :label="item.label"
            :value="item"
            />
          </el-select>
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
  createExaVIPCombo,
  updateExaVIPCombo,
  deleteExaVIPCombo,
  getExaVIPComboList
} from '@/api/combo'
import { ref, watch, onBeforeMount } from 'vue'
import { ElMessage } from 'element-plus'
import { formatDate } from '@/utils/format'
import config from '@/core/config'
import { comboStore } from '@/pinia/modules/combo'

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const comboData = ref([])

const comboForm = ref({
  Id: 0,
  comboName: '',
  comboType: null,
  comboPrice: 0,
  times: 0,
  state: 0,
  stateStr: '',
})

type Option = {
  id: number
  label: string
}
const comboTypeValue = ref<Option>()
// const comboTypeOptions = ref([
//     { id: 1, label: 'Vip次卡' },
//     { id: 2, label: 'Vip周期卡' },
//     { id: 3, label: 'Vip充值卡' },
//   ])
watch(() => comboTypeValue.value, () => {
  comboForm.value.comboType = comboTypeValue.value.id
})

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getComboData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getComboData()
}

type stateOption = {
  id: number
  label: string
}
const stateValue = ref<stateOption>()
const stateOptions = ref([])
watch(() => stateValue.value, () => {
  comboForm.value.state = Number(stateValue.value.id)
})

// 查询
const getComboData = async() => {
  await vipComboStore.GetAllVIPCombos()
  comboData.value = vipComboStore.comboList
  let total_num = 0
  if (comboData.value['length'] > 0) {
    comboData.value.forEach(element => {
      total_num += 1
      if (config.stateOptions[element.state]) {
        element.stateStr = config.stateOptions[element.state]
      }
    });
  }
  
  
  total.value = total_num
  page.value = 1
  pageSize.value = total_num
}

onBeforeMount(() => {
  for (let key in config.stateOptions) {
    let stateMap = {id: key, label: config.stateOptions[key]}
    stateOptions.value.push(stateMap)
  }
  stateValue.value = {id: 0, label: "正常"}
  comboForm.value.state = 0
  getComboData()
})


const dialogFormVisible = ref(false)
const operateType = ref('')
const updateCombo = async(row) => {
  dialogFormVisible.value = true
  operateType.value = 'update'
  comboForm.value = row
  // config.comboTypeOptions.forEach((element, _) => {
  //   if (element.id == comboForm.value.comboType) {
  //     comboTypeValue.value = element
  //   }
  // })
}
const closeDialog = () => {
  dialogFormVisible.value = false
  comboForm.value = {
    Id: 0,
    comboName: '',
    comboType: null,
    comboPrice: 0,
    times: 0,
    state: 0,
    stateStr: '',
  }
  comboTypeValue.value = {
    id: 0,
    label: ""
  }
}
const deleteCombo = async(row) => {
  row.visible = false
  console.log("---row:", row)
  const res = await deleteExaVIPCombo({ id: row.id })
  if ("code" in res && res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (comboData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getComboData()
  }
}
const vipComboStore = comboStore()

const enterDialog = async() => {
  if (comboForm.value.comboName == '' || !comboForm.value.comboType 
    || comboForm.value.comboPrice < 1 || comboForm.value.times < 1 || comboForm.value.state < 0) {
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
      res = await createExaVIPCombo(comboForm.value)
      break
    case 'update':
      res = await updateExaVIPCombo(comboForm.value)
      if (res.code === 0) {
        let comboList = await vipComboStore.comboList
        if (comboList['length'] < 1) {
          await vipComboStore.GetAllVIPCombos()
        }
      }
      break
    default:
      res = await createExaVIPCombo(comboForm.value)
      break
  }

  if (res.code === 0) {
    closeDialog()
    getComboData()
  }
}
const openDialog = () => {
  operateType.value = 'create'
  dialogFormVisible.value = true
}

</script>

<style></style>
