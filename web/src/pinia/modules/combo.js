import { getAllVIPComboList } from '@/api/combo'
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

export const comboStore = defineStore('combo', () => {
  const comboInfo = ref({
    Id: 0,
    storeName: '',
    comboName: '',
    comboType: '',
    comboPrice: 0,
    amount: 0,
  })
  const comboList = ref([])
  const setComboList = (val) => {
    comboList.value = val
  }

  /* 获取用户信息*/
  const GetAllVIPCombos = async() => {
    const res = await getAllVIPComboList()
    if (res.code === 0) {
        setComboList(res.data.list)
    }
    return res
  }
  

  return {
    comboInfo,
    comboList,
    GetAllVIPCombos,
  }
})
