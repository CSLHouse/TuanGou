<template>
  <div>
    <div>
      <el-card class="box-card">
        <el-form :inline="true" :model="searchData" class="demo-form-inline">
          <el-form-item label="联系电话：" class="form-item">
            <el-input v-model="searchData.telephone" placeholder="按联系电话搜索" clearable 
              @input="onTriggerSearch" @clear="onSearch"/>
          </el-form-item>
          <el-form-item label="姓名：" class="form-item">
            <el-input v-model="searchData.memberName" placeholder="按姓名搜索" clearable 
              @input="onTriggerSearch" @clear="onSearch"/>
          </el-form-item>
          <el-form-item label="会员到期：" class="form-item">
            <el-date-picker
              v-model="searchData.deadline"
              placeholder="按到期日期搜索"
              clearable
              type="date"
              format="YYYY/MM/DD"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item label="状态：" class="form-item">
            <el-select v-model="searchStateValue" value-key="id" class="m-2" 
              placeholder="请选择会员卡" size="large" clearable @clear="onSearch">
              <el-option
              v-for="item in config.memberStateOptions"
              :key="item.id"
              :label="item.label"
              :value="item"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="会员形式：" class="form-item">
            <el-select v-model="searchData.tmp" value-key="id" class="m-2" 
              placeholder="请选择会员形式" size="large" clearable @clear="onSearch">
              <el-option
              v-for="item in memberType"
              :key="item.id"
              :label="item.label"
              :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item class="form-item">
            <el-button type="primary" @click="onSearch">搜索</el-button>
          </el-form-item>
          <el-form-item class="form-item">
            <el-button  @click="onCancel">清空</el-button>
          </el-form-item>
          <el-form-item class="form-item">
            <el-button type="primary" @click="onExport">导出
              <el-icon class="el-icon--right"><Download /></el-icon>
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
      
    </div>
    <div class="gva-table-box">
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="Id"
        @cell-dblclick="updateMember"
      >
        <!-- <el-table-column type="selection" width="55" /> -->
        <el-table-column align="left" label="序号" prop="id" width="60"></el-table-column>
        <el-table-column align="left" label="会员卡号" prop="cardId" width="120" />
        <el-table-column align="left" label="姓名" prop="userName" width="120" />
        <el-table-column align="left" label="联系方式" prop="telephone" width="120" />
        <el-table-column align="left" label="会员类型" prop="comboType" width="100" />
        <el-table-column align="left" label="剩余次数/金额" prop="remainTimes" width="90" />
        <el-table-column align="left" label="办卡日期" prop="startDate" width="100" />
        <el-table-column align="left" label="截止日期" prop="deadline" width="100" />
        <el-table-column align="left" label="状态" prop="stateStr" width="80" />
        <el-table-column align="left" label="操作" min-width="160">
          <template #default="scope">
            <el-button type="primary" link icon="edit" @click="updateMember(scope.row)">变更</el-button>
            <el-popover v-model="scope.row.visible" placement="top" width="160">
              <p>确定要删除吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                <el-button type="primary" @click="deleteMember(scope.row)">确定</el-button>
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
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="修改会员">
      <el-form  :model="memberForm" label-width="140px" >
        <el-form-item label="会员卡号">
          <el-input v-model="memberForm.cardId" autocomplete="off" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="memberForm.telephone" autocomplete="off" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="memberForm.userName" autocomplete="off" />
        </el-form-item>
        <el-form-item label="会员类型">
          <el-select v-model="comboOption" value-key="key" class="m-2" placeholder="请选择会员卡" size="large">
            <el-option
            v-for="item in comboOptions"
            :key="item.key"
            :label="item.value"
            :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="当前次数/金额">
          <el-input v-model.number="memberForm.remainTimes" autocomplete="off" />
        </el-form-item>
        <el-form-item label="会员起始日期">
          <el-date-picker
            v-model="memberForm.startDate"
            type="date"
            format="YYYY/MM/DD"
            value-format="YYYY-MM-DD">
          </el-date-picker>
        </el-form-item>
        <el-form-item label="会员截止日期">
          <el-date-picker
            v-model="memberForm.deadline"
            type="date"
            format="YYYY/MM/DD"
            value-format="YYYY-MM-DD">
          </el-date-picker>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="stateValue" value-key="id" class="m-2" placeholder="请选择会员卡" size="large">
            <el-option
            v-for="item in config.memberStateOptions"
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
  createVIPMember,
  updateVIPMember,
  deleteVIPMember,
  getVIPMemberList,
  searchVIPMembers,
} from '@/api/member'
import { ref, reactive, onBeforeMount, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { comboStore } from '@/pinia/modules/combo'
import config from '@/core/config'
import { trim } from '@/utils/stringFun'
const searchData = reactive({
  telephone: null,
  memberName: '',
  deadline: '',
  state: null,
  stateStr: "",
  tmp: null,
})

const onTriggerSearch = () => {
  if (searchData.telephone > 1000 || trim(searchData.memberName).length > 1 || (searchData.deadline && trim(searchData.deadline).length == 10) || searchData.state > 0) {
    onSearch()
  }
}

const onSearch = async() => {
  const res = await searchVIPMembers({ telephone: searchData.telephone, memberName: searchData.memberName, 
      deadline: searchData.deadline, state: searchData.state, tmp: searchData.tmp, page: page.value, pageSize: pageSize.value})
  if ('code' in res && res.code === 0) {
    if (res.data.list) {
      tableData.value = res.data.list
      tableData.value.forEach(memberElement => {
        memberElement.comboType = memberElement.combo.comboName
        config.memberStateOptions.forEach(stateElement => {
          if (memberElement.state == stateElement.id) {
            memberElement.stateStr = stateElement.label
          }
        })
      });
      total.value = res.data.total
      page.value = res.data.page
      pageSize.value = res.data.pageSize
    }
  }
}

const onCancel = () => {
  searchData.telephone = null
  searchData.memberName = ''
  searchData.deadline = ''
  searchData.state = null
  searchData.stateStr = ''
  getTableData()
}

const onExport = () => {

}

const memberForm = ref({
  Id: 0,
  cardId: null,
  telephone: null,
  userName: '',
  comboType: '',
  remainTimes: null,
  startDate: '',
  deadline: '',
  state: null,
  comboId: null,
  stateStr: ''
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

// 查询
const getTableData = async() => {
  const table = await getVIPMemberList({ page: page.value, pageSize: pageSize.value })
  if ('code' in table && table.code === 0) {
    tableData.value = table.data.list
    tableData.value.forEach(memberElement => {
      memberElement.comboType = memberElement.combo.comboName
      config.memberStateOptions.forEach(stateElement => {
        if (memberElement.state == stateElement.id) {
          memberElement.stateStr = stateElement.label
        }
      })
    });
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

// 查询会员卡类型
const vipComboStore = comboStore()
interface comboItem {
  key: number,
  value: string,
}
const comboOption = ref<comboItem>()
const comboOptions = ref<comboItem[]>([])

const getComboData = async() => {
  let comboList = vipComboStore.comboList
  if (comboList['length'] < 1) {
    await vipComboStore.GetAllVIPCombos()
    comboList = vipComboStore.comboList
  }

  comboOptions.value = comboList.map((item) => {
    return {key: item.Id, value: item.comboName}
  })
}
watch(() => comboOption.value, () => {
  memberForm.value.comboType = comboOption.value.value
  memberForm.value.comboId = comboOption.value.key
})

getTableData()
// 
type Option = {
  id: number
  label: string
}
const stateValue = ref<Option>()
watch(() => stateValue.value, () => {
  memberForm.value.state = stateValue.value.id
})

const searchStateValue = ref<Option>()
watch(() => searchStateValue.value, () => {
  searchData.state = searchStateValue.value.id
  onSearch()
})

watch(() => searchData.deadline, () => {
  onSearch()
})

const type = ref('')
const updateMember = async(row) => {
  dialogFormVisible.value = true
  getComboData()
  type.value = 'update'
  memberForm.value = row
  comboOption.value = {key:row.comboId, value: row.comboType}
  stateValue.value = {id: row.state, label:row.stateStr}
}

const deleteMember = async(row) => {
  row.visible = false
  const res = await deleteVIPMember({ id: row.id })
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

const dialogFormVisible = ref(false)
const enterDialog = async() => {
  let res
  switch (type.value) {
    case 'create':
      res = await createVIPMember(memberForm.value)
      break
    case 'update':
      res = await updateVIPMember(memberForm.value)
      break
    default:
      res = await createVIPMember(memberForm.value)
      break
  }

  if (res.code === 0) {
    closeDialog()
    getTableData()
  }
}

const closeDialog = () => {
  dialogFormVisible.value = false
}

const memberType = ref([
    {
      id: 1,
      label: '会员',
      value: 100
    },
    {
      id: 2,
      label: '临时会员',
      value: 101
    }
  ])
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
