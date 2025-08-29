<template>
    <div>
      <div>
        <el-card class="box-card">
          <el-form :inline="true" :model="searchData" class="demo-form-inline">
            <el-form-item label="会员电话：" class="form-item">
              <el-input v-model="searchData.telephone" placeholder="按会员电话搜索" clearable 
                @input="onTriggerSearch" @clear="onSearch"/>
            </el-form-item>
            <!-- <el-form-item label="状态：" class="form-item">
              <el-select v-model="searchStateValue" value-key="id" class="m-2" 
                placeholder="请选择会员卡" size="large" clearable @clear="onSearch">
                <el-option
                v-for="item in config.memberStateOptions"
                :key="item.id"
                :label="item.label"
                :value="item"
                />
              </el-select>
            </el-form-item> -->
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
          <div v-if="tableVisible">
            <el-table :data="memberData" style="width: 100%"
                highlight-current-row
                :show-header="false"
                @row-click="OnConfirm">
                <el-table-column prop="cardId" label="会员卡号" width="120" height="40" />
                <el-table-column prop="telephone" label="手机号" width="120" height="40" />
                <el-table-column prop="memberName" label="姓名" width="80" />
                <el-table-column prop="comboType" label="套餐类型" width="80"/>
                <el-table-column prop="remainTimes" label="当前剩余" width="80"/>
                <el-table-column prop="deadline" label="截止日期" width="100"/>
                <el-table-column prop="stateStr" label="状态" width="100"/>
            </el-table>
          </div>
        </el-card>
      </div>
      <div class="gva-table-box">
        <el-table
          ref="multipleTable"
          :data="tableData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="Id"
        >
          <el-table-column align="left" label="入店序号" prop="Id" width="100"></el-table-column>
          <el-table-column align="left" label="会员" prop="memberStr" width="180" />
          <el-table-column align="left" label="会员状态" prop="memberStateStr" width="180" />
          <el-table-column align="left" label="会员类型" prop="memberType" width="100" />
          <el-table-column align="left" label="会员到期时间" prop="deadline" width="120" />
          <el-table-column align="left" label="剩余次数/余额" prop="remainTimes" width="120" />
          <el-table-column align="left" label="消费次数/余额" prop="consumeTimes" width="120" />
          <el-table-column align="left" label="打卡时间" prop="punchDate" width="160" />
          <el-table-column align="left" label="状态" prop="stateStr" width="80" />
          <!-- <el-table-column align="left" label="操作" min-width="160">
            <template #default="scope">
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
          </el-table-column> -->
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
    </div>
  </template>
  
  <script setup lang="ts">
  import {
    getVIPConsumeList,
    searchVIPMembers,
  } from '@/api/member'
  import { ref, reactive, onBeforeMount, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { comboStore } from '@/pinia/modules/combo'
  import config from '@/core/config'
  import { trim } from '@/utils/stringFun'
  const searchData = reactive({
    telephone: null,
  })
  const tableVisible = ref(false)
  const memberData = ref([])
  const onTriggerSearch = () => {
    if (searchData.telephone >= 1000 ) {
      onSearch()
    }
  }
  
  const onSearch = async() => {
    const res = await searchVIPMembers({ telephone: searchData.telephone})
    if ('code' in res && res.code === 0) {
      if (res.data.list) {
        tableVisible.value = true
        memberData.value = res.data.list
        memberData.value.forEach(memberElement => {
            if (memberElement.state == 0) {
                memberElement.stateStr = "已确认"
            }
          config.memberStateOptions.forEach(stateElement => {
            if (memberElement.memberState == stateElement.id) {
              memberElement.memberStateStr = stateElement.label
            }
          })
        });
      }
    }
  }
  
  const onCancel = () => {
    searchData.telephone = null
    getTableData()
  }
  
  const onExport = () => {
  
  }
  
  const OnConfirm = async(row) => {
    getTableData()
  }
  
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
    const table = await getVIPConsumeList({ page: page.value, pageSize: pageSize.value })
    if ('code' in table && table.code === 0) {
      tableData.value = table.data.list
      tableData.value.forEach(memberElement => {
        memberElement.memberStr = memberElement.memberName + "(" + memberElement.telephone + ")"        
        if (memberElement.state == 0) {
            memberElement.stateStr = "已确认"
        }
        config.memberStateOptions.forEach(stateElement => {
          if (memberElement.memberState == stateElement.id) {
            memberElement.memberStateStr = stateElement.label
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
    // 一定放在GetAllVIPCombos返回之后
    getTableData()
  }
  getComboData()
  
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
  