<template>
    <div>
      <div>
        <el-card class="box-card">
          <el-form :inline="true" :model="searchData" class="demo-form-inline">
            <el-form-item label="订单编号：" class="form-item">
              <el-input v-model.number="searchData.orderId" placeholder="按订单编号搜索" clearable 
                @input="onTriggerSearch" @clear="onSearch"/>
            </el-form-item>
            <el-form-item label="联系电话：" class="form-item">
              <el-input v-model.number="searchData.telephone" placeholder="按联系电话搜索" clearable 
                @input="onTriggerSearch" @clear="onSearch"/>
            </el-form-item>
            <el-form-item label="订单类型：" class="form-item">
              <el-input v-model.number="searchData.comboType" placeholder="按订单类型搜索" clearable 
                @input="onTriggerSearch" @clear="onSearch"/>
            </el-form-item>
            <el-form-item label="购买日期：" class="form-item">
              <el-date-picker
                v-model="searchData.buyDate"
                placeholder="按购买日期搜索"
                clearable
                type="date"
                format="YYYY/MM/DD"
                value-format="YYYY-MM-DD"
              />
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
        >
          <el-table-column align="left" label="序号" prop="Id" width="60"></el-table-column>
          <el-table-column align="left" label="订单编号" prop="orderId" width="180" />
          <el-table-column align="left" label="用户名称" prop="nameStr" width="180" />
          <el-table-column align="left" label="客户经理" prop="manager" width="120" />
          <el-table-column align="left" label="订单类型" prop="comboType" width="100" />
          <el-table-column align="left" label="商品名称" prop="comboName" width="160" />
          <el-table-column align="left" label="商品价格" prop="comboPrice" width="100" />
          <el-table-column align="left" label="实收金额" prop="collection" width="100" />
          <el-table-column align="left" label="购买日期" prop="buyDate" width="100" />
          <el-table-column align="left" label="订单状态" prop="stateStr" width="80" />
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
    getVIPOrderList,
  } from '@/api/member'
  import { ref, reactive, onBeforeMount, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { comboStore } from '@/pinia/modules/combo'
  import config from '@/core/config'
  import { trim } from '@/utils/stringFun'
  const searchData = reactive({
    orderId: null,
    telephone: null,
    buyDate: '',
    comboType: null,
  })
  
  const onTriggerSearch = () => {
    if (searchData.orderId > 10|| searchData.telephone > 1000 ||  trim(searchData.buyDate).length == 10 || searchData.comboType > 0) {
        getOrderList()
    }
  }
  
  const getOrderList = async() => {
    const res = await getVIPOrderList({orderId: searchData.orderId, telephone: searchData.telephone, buyDate: searchData.buyDate, 
        comboType: searchData.comboType, page: page.value, pageSize: pageSize.value})
    if ('code' in res && res.code === 0) {
      if (res.data.list) {
        tableData.value = res.data.list
        tableData.value.forEach(memberElement => {
            memberElement.nameStr = memberElement.memberName + "(" + memberElement.telephone + ")"
            let productName = ""
            if (memberElement.isNew) {
                productName = productName + '新'
            }

            config.memberStateOptions.forEach(stateElement => {
                if (memberElement.state == stateElement.id) {
                    productName = productName + stateElement.label
                    memberElement.stateStr = stateElement.label
                }
            })
            productName = productName + ":" + memberElement.comboType
            memberElement.comboName = productName
            
        });
        total.value = res.data.total
        page.value = res.data.page
        pageSize.value = res.data.pageSize
      }
    }
  }
  
  const onCancel = () => {
    searchData.orderId = null
    searchData.telephone = null
    searchData.buyDate = ''
    searchData.comboType = null
    getOrderList()
  }
  
  const onExport = () => {
  
  }
  
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  
  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getOrderList()
  }
  
  const handleCurrentChange = (val) => {
    page.value = val
    getOrderList()
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
    getOrderList()
  }

  getComboData()
  
  watch(() => searchData.buyDate, () => {
    getOrderList()
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
  </style>
  