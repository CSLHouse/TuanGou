
<template>
    <div>
      <el-upload
        ref="uploadRef"
        :list-type="listType"
        :action="`${path}/fileUploadAndDownload/upload`"
        :headers="{ 'x-token': userStore.token }"
        :show-file-list="isShow"
        :on-success="handleImageSuccess"
        :before-upload="beforeImageUpload"
        :on-remove="handleRemove"
        :on-preview="handlePreview"
        :limit="1"
        :file-list="fileDetailList"
        :on-exceed="handleExceed"
      >
        <el-button size="small" type="primary">
          <template #default="scope">
              <div v-if="isRefresh">
                重新上传
              </div>
              <div v-else>
                点击上传
              </div>
          </template>
        </el-button>
      </el-upload>
    </div>
  </template>
  
  <script setup>
  import ImageCompress from '@/utils/image'
  import { ref, toRef, nextTick, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { useUserStore } from '@/pinia/modules/user'
  import { deleteFile } from '@/api/fileUploadAndDownload'
  import { genFileId } from 'element-plus'
  // import { genFileId } from 'element-plus'
  // import type { UploadInstance, UploadProps, UploadRawFile } from 'element-plus'

  const emits = defineEmits(["update:modelValue"])
  const props = defineProps({
    imageUrl: {
      type: String,
      default: ''
    },
    fileSize: {
      type: Number,
      default: 2048 // 2M 超出后执行压缩
    },
    maxWH: {
      type: Number,
      default: 1920 // 图片长宽上限
    },
    isShow: {
      type: Boolean,
      default: true
    },
    listType: {
      type: String,
      default: "text"
    },
    modelValue: String,

  })
  const modelValue = toRef(props.modelValue)
  const path = ref(import.meta.env.VITE_BASE_API)
  const userStore = useUserStore()
  let fileDetailList = ref([])  // 只能上传一个
  const isRefresh = ref(false)
  computed: {
    if (props.modelValue) {
      var index = props.modelValue.lastIndexOf("\/")
      var name = props.modelValue.substring(index + 1, props.modelValue.length)
      fileDetailList.value[0] = {url: props.modelValue, name: name}

      isRefresh.value = true
    }
  }

  const beforeImageUpload = (file) => {
    const isJPG = file.type === 'image/jpeg'
    const isPng = file.type === 'image/png'
    if (!isJPG && !isPng) {
      ElMessage.error('上传头像图片只能是 jpg或png 格式!')
      return false
    }
  
    const isRightSize = file.size / 1024 < props.fileSize
    if (!isRightSize) {
      // 压缩
      const compress = new ImageCompress(file, props.fileSize, props.maxWH)
      return compress.compress()
    }
    return isRightSize
  }
  
  const handleImageSuccess = (res) => {
    const { data } = res
    if (data.file) {
        fileDetailList.value = []
        fileDetailList.value[0] = {url: data.file.url, name: data.file.name}
        emits('update:modelValue', data.file.url)
    }
  }
  
  const handleRemove = async(file) => {
    if (file.status == "success") {
      const res = await deleteFile(file)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!',
        })
      }
    }
  }
  const handlePreview = (url) => {
    ElMessage.warning('已上传过同名文件!')
    return
  }

  const uploadRef = ref()

  const handleExceed = async(files) => {

    const res = await deleteFile(fileDetailList.value[0])
    if (res.code === 0) {
      fileDetailList.value = []
      uploadRef.value.clearFiles()
      nextTick(() => {
        files[0].uid = genFileId()
        uploadRef.value.handleStart(files[0])
        uploadRef.value.submit()
        fileDetailList.value[0] = files[0]
      })

      ElMessage({
        type: 'success',
        message: '更新成功!',
      })
    }
  }
  </script>
  
  <script>
  
  export default {
    name: 'CoollerSingeUpload',
    methods: {
  
    }
  }
  </script>
  
  <style lang="scss" scoped>
  .image-uploader {
    border: 1px dashed #d9d9d9;
    width: 180px;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
  }
  .image-uploader {
    border-color: #409eff;
  }
  .image-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 178px;
    height: 178px;
    line-height: 178px;
    text-align: center;
  }
  .image {
    width: 178px;
    height: 178px;
    display: block;
  }

  .upload-queue{
  .el-upload-list__item:hover{
    .el-upload-list__item-actions{
      display: block;
      opacity: 1;
    }
  }

  .el-upload--picture-card .el-upload--text{
    width: 100%;
  }
}
  </style>
  