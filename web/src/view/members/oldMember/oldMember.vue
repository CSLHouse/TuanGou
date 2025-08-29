<template>
  <div>
    <warning-bar
      title="支持模糊查找，左键选中自动填充！"
    />
  </div>
  
  <el-menu class="el-menu-demo" mode="horizontal">
    <el-menu-item index="1">添加订单</el-menu-item>
  </el-menu>
  <div class="h-6" />
  <el-form
    ref="ruleFormRef"
    :model="memberForm"
    status-icon
    label-width="120px"
    class="demo-ruleForm"
  >
    <el-form-item label="手机/会员卡" prop="cardId" >
      <el-col :span="8">
        <el-input v-model="memberForm.onlyId" type="text" autocomplete="off"
        clearable @input="onTriggerSearch" @clear="onSearch"/>
      </el-col>
    </el-form-item>
    <el-form-item label="" prop="cardId">
      <div v-if="tableVisible">
        <el-table :data="tableData" style="width: 100%"
        highlight-current-row
        @row-click="OnConfirm">
          <el-table-column prop="cardId" label="会员卡号" width="120" height="40" />
          <el-table-column prop="telephone" label="手机号" width="120" height="40" />
          <el-table-column prop="userName" label="姓名" width="80" />
          <el-table-column prop="comboType" label="套餐类型" width="80"/>
          <el-table-column prop="remainTimes" label="当前剩余" width="80"/>
          <el-table-column prop="deadline" label="截止日期" width="100"/>
        </el-table>
      </div>
      
    </el-form-item>
    <el-form-item label="会员姓名" prop="telephone">
      <el-col :span="8">
        <el-input v-model="memberForm.userName" type="tel" autocomplete="off" disabled/>
      </el-col>
    </el-form-item>
    <el-form-item label="VIP卡" prop="comboId">
      <el-select v-model.number="comboOption" class="m-2" placeholder="请选择会员卡" size="large">
          <el-option
          v-for="item in comboOptions"
          :key="item.key"
          :label="item.value"
          :value="item"
          />
      </el-select>
    </el-form-item>
    <el-form-item label="赠送" prop="gift">
      <el-col :span="8">
        <el-input v-model.number="memberForm.times" type="number" autocomplete="off" />
      </el-col>
      <a style="color: red; margin-left: 10px;">整数，填写赠送的天数/次数/金额</a>
    </el-form-item>
    <el-form-item label="实收金额" prop="collection">
      <el-col :span="8">
        <el-input v-model.number="memberForm.collection" type="number" autocomplete="off" />
      </el-col>
      <a style="color: red; margin-left: 10px;">整数，必须与选择的Vip卡价格一致</a>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="submitForm(ruleFormRef)"
        >确认</el-button
      >
      <el-button @click="resetForm(ruleFormRef)">清空</el-button>
    </el-form-item>
  </el-form>
</template>
  
<script lang="ts" setup>
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { searchVIPCards, renewVIPCards } from '@/api/member'
  import { reactive, ref, onBeforeMount, watch } from 'vue'
  import { FormInstance, ElMessage } from 'element-plus'
  import { comboStore } from '@/pinia/modules/combo'
  import { trim } from '@/utils/stringFun'
  import config from '@/core/config'

  const ruleFormRef = ref<FormInstance>()
  const tableVisible = ref(false)
  const tableData = ref([])
  interface comboItem {
    key: number,
    value: string,
    price: number,
  }
  const comboOption = ref<comboItem>()
  const comboOptions = ref<comboItem[]>([])
  // 查询会员卡类型
  const vipComboStore = comboStore()

  const memberForm = reactive({
    onlyId: null,
    id: null,
    cardId: null,
    userName: "",
    times: null,
    comboId: null,
    collection: null,
  })

  const onTriggerSearch = () => {
    if (memberForm.onlyId > 1000 ) {
      onSearch()
    }
  }

  const onSearch = async() => {
    const res = await searchVIPCards({ onlyId: memberForm.onlyId})
    if ('code' in res && res.code === 0) {
      if (res.data.list) {
        tableVisible.value = true
        tableData.value = res.data.list
        tableData.value.forEach(memberElement => {
          memberElement.comboType = memberElement.combo.comboName

          config.memberStateOptions.forEach(stateElement => {
            if (memberElement.state == stateElement.id) {
              memberElement.stateStr = stateElement.label
            }
          })
        });
      }
    }
  }
  
  const OnConfirm = async(row) => {
    memberForm.id = row.id
    memberForm.cardId = row.cardId
    memberForm.userName = row.userName
    memberForm.comboId = row.comboId
    memberForm.times = 0
    memberForm.collection = 0
    let comboList = vipComboStore.comboList
    if (comboList['length'] < 1) {
      await vipComboStore.GetAllVIPCombos()
      comboList = vipComboStore.comboList
    }
    comboList.forEach(element => {
      if (element.Id == memberForm.comboId) {
        comboOption.value = {key: element.id, value: element.comboName, price: element.comboPrice}
      }
    });
  }

  const getComboData = async() => {
    let comboList = vipComboStore.comboList
    if (comboList['length'] < 1) {
      await vipComboStore.GetAllVIPCombos()
      comboList = vipComboStore.comboList
    }
    comboList.forEach(element => {
      comboOptions.value.push({key: element.id, value: element.comboName, price: element.comboPrice})
    });
  }

  onBeforeMount(() => {
    getComboData()
  })

  watch(() => comboOption.value, () => {
    memberForm.comboId = comboOption.value.key
  })

  const submitForm = async(formEl: FormInstance | undefined) => {
    if (memberForm.collection != comboOption.value.price) {
      ElMessage({
        type: 'error',
        message: '信息填写错误'
      })
      return
    }
    const memberTable = {
      id: memberForm.id,
      cardId: memberForm.cardId,
      times: memberForm.times,
      comboId: memberForm.comboId,
      collection: memberForm.collection
    }
    let res = await renewVIPCards(memberTable)
    if ('code' in res && res.code === 0) {
      formEl.resetFields()
      tableVisible.value = false
      memberForm.userName = ""
      memberForm.onlyId = null
      memberForm.comboId = 0
      memberForm.times = null
      memberForm.collection = null
      tableData.value = []
      comboOption.value = {key: 0, value: '', price: 0}
      ElMessage({
        type: 'success',
        message: '修改成功'
      })
    }
  }
  
  const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
    tableVisible.value = false
    memberForm.userName = ""
    memberForm.onlyId = null
    memberForm.comboId = 0
    memberForm.times = null
    memberForm.collection = null
    tableData.value = []
    comboOption.value = {key: 0, value: '', price: 0}
  }
</script>
