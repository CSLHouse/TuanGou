<template> 
  <div>
    <el-upload
      :action="useOss?ossUploadUrl:minioUploadUrl"
      :headers="{ 'x-token': userStore.token }"
      :data="useOss?dataObj:null"
      list-type="picture-card"
      :file-list="fileList"
      :before-upload="beforeUpload"
      :on-remove="handleRemove"
      :on-success="handleUploadSuccess"
      :on-preview="handlePreview"
      :limit="maxCount"
      :on-exceed="handleExceed"
    >
      <i class="el-icon-plus"></i>
    </el-upload>
    <el-dialog :visible.sync="dialogVisible">
      <img width="100%" :src="dialogImageUrl" alt="">
    </el-dialog>
  </div>
</template>
<script setup>
import { ref, toRef,computed, defineProps, defineEmits } from 'vue'
import { policy } from '@/api/oss'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'

// 定义props
const props = defineProps({
  // 图片属性数组
  modelValue: {
    type: Array,
    default: () => []
  },
  // 最大上传图片数量
  maxCount: {
    type: Number,
    default: 5
  }
})
const modelValue = toRef(props.modelValue)
const maxCount = toRef(props.maxCount)
// 定义emit - 使用setup语法糖的defineEmits
const emit = defineEmits(['update:modelValue'])

// 响应式数据
const dataObj = ref({
  policy: '',
  signature: '',
  key: '',
  ossaccessKeyId: '',
  dir: '',
  host: ''
})

const userStore = useUserStore()
const dialogVisible = ref(false)
const dialogImageUrl = ref(null)
const useOss = ref(false) // 使用oss->true;使用MinIO->false
const ossUploadUrl = ref('http://macro-oss.oss-cn-shenzhen.aliyuncs.com')
const path = ref(import.meta.env.VITE_BASE_API)
const minioUploadUrl = ref(path.value + '/fileUploadAndDownload/upload')

// 计算属性fileList
const fileList = computed(() => {
  let fileList = []
  if (modelValue.value && modelValue.value.length > 0) {
    // 将value数组转换为fileList数组
    for (let i = 0; i < modelValue.value.length; i++) {
      fileList.push({ url: modelValue.value[i] })
    }
  }
  console.log('fileList', fileList)
  return fileList
})

// 方法定义
const emitInput = (fileList) => {
  let value = []
  for (let i = 0; i < fileList.length; i++) {
    value.push(fileList[i].url)
  }
  // 直接使用defineEmits返回的emit函数触发事件
  emit('update:modelValue', value)
}

const handleRemove = (file, fileList) => {
  emitInput(fileList)
}

const handlePreview = (file) => {
  dialogVisible.value = true
  dialogImageUrl.value = file.url
}

const beforeUpload = (file) => {
  if (!useOss.value) {
    // 不使用oss不需要获取策略
    return true
  }
  return new Promise((resolve, reject) => {
    policy()
      .then(response => {
        dataObj.value.policy = response.data.policy
        dataObj.value.signature = response.data.signature
        dataObj.value.ossaccessKeyId = response.data.accessKeyId
        dataObj.value.key = response.data.dir + '/${filename}'
        dataObj.value.dir = response.data.dir
        dataObj.value.host = response.data.host
        resolve(true)
      })
      .catch(err => {
        console.log(err)
        reject(false)
      })
  })
}

const handleUploadSuccess = (res, file) => {
  let url = dataObj.value.host + '/' + dataObj.value.dir + '/' + file.name
  if (!useOss.value) {
    // 不使用oss直接获取图片路径
    url = res.data.url
  }
  // 创建新的fileList数组，避免直接修改计算属性
  const newFileList = [...fileList.value, { name: file.name, url: url }]
  emitInput(newFileList)
}

const handleExceed = (files, fileList) => {
  ElMessage({
    message: `最多只能上传${maxCount.value}张图片`,
    type: 'warning',
    duration: 1000
  })
}
</script>
<style>

</style>


