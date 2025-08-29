<template>
    <div>
      <div>
        <el-card class="box-card">

          <el-form :inline="true" :model="searchData" class="demo-form-inline">
            <el-form-item label="开始日期：" class="form-item">
              <el-date-picker
                v-model="searchData.startDate"
                placeholder="按开始日期搜索"
                clearable
                type="date"
                format="YYYY/MM/DD"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item label="结束日期：" class="form-item">
              <el-date-picker
                v-model="searchData.endDate"
                placeholder="按结束日期搜索"
                clearable
                type="date"
                format="YYYY/MM/DD"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
            <el-form-item class="form-item">
              <el-button type="primary" @click="getsStatementList">搜索</el-button>
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
        <div style="margin-bottom: 10px;">
            <el-text class="mx-1">会员卡总流水:</el-text>
            <el-text class="mx-1" type="danger">{{totalStream}}</el-text>
            <el-text class="mx-1">会员卡总单量:</el-text>
            <el-text class="mx-1" type="danger">{{totalOrder}}</el-text>
            <el-text class="mx-1">总会员:</el-text>
            <el-text class="mx-1" type="danger">{{totalMember}}</el-text>
            <el-text class="mx-1">入店总人数:</el-text>
            <el-text class="mx-1" type="danger">{{totalConsumer}}</el-text>
        </div>
        
        <el-table
          ref="multipleTable"
          :data="tableData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="Id"
        >
          <el-table-column align="left" label="日期" prop="date" width="100"></el-table-column>
          <el-table-column align="left" label="会员卡流水" prop="recharge" width="180" />
          <el-table-column align="left" label="会员卡单量" prop="cardNumber" width="160" />
          <el-table-column align="left" label="新增会员" prop="newMember" width="120" />
          <el-table-column align="left" label="入店统计" prop="consumeNumber" width="100" />
        </el-table>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import {
    getStatementList,
    getStatisticsList,
  } from '@/api/member'
  import { ref, reactive, onBeforeMount, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  const searchData = reactive({
    startDate: null,
    endDate: null,
  })
  const totalStream = ref(0)    // 总流水
  const totalOrder = ref(0)     // 总单量
  const totalMember = ref(0)     // 总会员
  const totalConsumer = ref(0)     // 入店总人数

  const getsStatisticsList = async() => {
    const res = await getStatisticsList()
    if ('code' in res && res.code === 0) {
      let total = res.data
      totalStream.value = total.totalStream
      totalOrder.value = total.totalOrder
      totalMember.value = total.totalMember
      totalConsumer.value = total.totalConsumer
    }
  }
  
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])

  const onTriggerSearch = () => {
    if (searchData.startDate >= 10|| searchData.endDate > 1000 ) {
        getsStatementList()
    }
  }
  
  const getsStatementList = async() => {
    const res = await getStatementList({startDate: searchData.startDate, endDate: searchData.endDate})
    if ('code' in res && res.code === 0) {
      if (res.data.list) {
        tableData.value = res.data.list
        let rechargeTotal = 0
        let cardTotal = 0
        let newMemberTotal = 0
        let consumeNumberTotal = 0
        tableData.value.forEach(item => {
            rechargeTotal += item.recharge
            cardTotal += item.cardNumber
            newMemberTotal += item.newMember
            consumeNumberTotal += item.consumeNumber
        })
        tableData.value.push({date: "合计", recharge: rechargeTotal, cardNumber: cardTotal, newMember: newMemberTotal, consumeNumber: consumeNumberTotal})
      }
    }
  }
  
  const onCancel = () => {
    searchData.startDate = null
    searchData.endDate = null
    getsStatementList()
  }
  
  const onExport = () => {
  
  }
  
  getsStatementList()
  getsStatisticsList()
  
  // watch(() => searchData.startDate, () => {
  //   getsStatementList()
  // })
  
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
  