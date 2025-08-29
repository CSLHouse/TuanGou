<template> 
  <div>
    <el-upload
      :headers="{ 'x-token': userStore.token }"
      :action="useOss?ossUploadUrl:minioUploadUrl"
      :data="useOss?dataObj:null"
      list-type="picture"
      :multiple="false" :show-file-list="showFileList"
      :file-list="fileList"
      :before-upload="beforeUpload"
      :on-remove="handleRemove"
      :on-success="handleUploadSuccess"
      :on-preview="handlePreview">
      <el-button size="small" type="primary">点击上传</el-button>
      <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过10MB</div>
    </el-upload>
    <el-dialog :visible.sync="dialogVisible">
      <img width="100%" :src="fileList[0].url" alt="">
    </el-dialog>
  </div>
</template>
<script setup>
import { ref, computed, toRef } from 'vue';
import { policy } from '@/api/oss';
import { useUserStore } from '@/pinia/modules/user'

// 定义 emits
const emit = defineEmits(['update:modelValue']);

// 接收 props
const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
});

// 响应式数据
const dataObj = ref({
  policy: '',
  signature: '',
  key: '',
  ossaccessKeyId: '',
  dir: '',
  host: '',
  // callback:'',
});
const userStore = useUserStore()
const modelValue = toRef(props.modelValue)

const dialogVisible = ref(false);
const useOss = ref(false); // 使用oss->true;使用MinIO->false
const ossUploadUrl = ref('https://pingogo.oss-cn-shanghai.aliyuncs.com');
const minioUploadUrl = ref(import.meta.env.VITE_BASE_API + '/fileUploadAndDownload/upload');

// 计算属性
const imageUrl = computed(() => {
  return modelValue.value;
});

const imageName = computed(() => {
  if (modelValue.value != null && modelValue.value !== '') {
    return modelValue.value.substr(modelValue.value.lastIndexOf("/") + 1);
  } else {
    return null;
  }
});

const fileList = computed(() => {
  return [{
    name: imageName.value,
    url: imageUrl.value
  }];
});

const showFileList = computed({
  get() {
    return modelValue.value !== null && modelValue.value !== '' && modelValue.value !== undefined;
  },
  set(newValue) {
    // 空实现，保持原有逻辑
  }
});

// 方法
const emitInput = (val) => {
  emit('update:modelValue', val);
};

const handleRemove = (file, fileList) => {
  emitInput('');
};

const handlePreview = (file) => {
  dialogVisible.value = true;
};

const beforeUpload = (file) => {
  if (!useOss.value) {
    // 不使用oss不需要获取策略
    return true;
  }
  return new Promise((resolve, reject) => {
    policy().then(response => {
      dataObj.value.policy = response.data.policy;
      dataObj.value.signature = response.data.signature;
      dataObj.value.ossaccessKeyId = response.data.accessKeyId;
      dataObj.value.key = response.data.dir + '/${filename}';
      dataObj.value.dir = response.data.dir;
      dataObj.value.host = response.data.host;
      // dataObj.value.callback = response.data.callback;
      resolve(true);
    }).catch(err => {
      console.log(err);
      reject(false);
    });
  });
};

const handleUploadSuccess = (res, file) => {

  showFileList.value = true;
  // 由于fileList是计算属性，不能直接pop，这里通过修改props.value间接更新
  // let url = dataObj.value.host + '/' + dataObj.value.dir + '/' + file.name;
  if (!useOss.value) {
    // 不使用oss直接获取图片路径
    // url = res.data.url;
    const { data } = res
    if (data.file) {
      fileList.value.push({url: data.file.url})
      modelValue.value = data.file.url
    }
    emitInput(data.file.url);
  }
};
</script>
<style>

</style>


