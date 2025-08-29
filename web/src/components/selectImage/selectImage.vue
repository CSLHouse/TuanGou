<template>
  <div
    v-if="!multiple"
    :class="size === 'default' ? 'update-image' : 'single-img-small'"
    :style="{
      'background-image': `url(${getUrl(modelValue)})`,
    }"
  >
    <span class="update" @click="openChooseImg">
      <el-icon>
        <edit />
      </el-icon>
      选择图片</span>
  </div>
  <div v-else class="multiple-img">
    <div
      v-for="(item, index) in multipleValue"
      :key="index"
      class="update-image"
      :style="{
        'background-image': `url(${getUrl(item)})`,
      }"
    >
      <span class="update" @click="deleteImg(index)">
        <el-icon>
          <delete />
        </el-icon>
        删除图片</span>
    </div>
    <div class="add-image">
      <span class="update" @click="openChooseImg">
        <el-icon>
          <folder-add />
        </el-icon>
        上传图片</span>
    </div>
  </div>

  <el-drawer v-model="drawer" title="媒体库" size="650px">
    <warning-bar
      title="点击“文件名/备注”可以编辑文件名或者备注内容。"
    />
    <div class="gva-btn-list">
      <upload-common
        v-model:imageCommon="imageCommon"
        class="upload-btn-media-library"
        @on-success="getImageList"
      />
      <upload-image
        v-model:imageUrl="imageUrl"
        :file-size="512"
        :max-w-h="1080"
        class="upload-btn-media-library"
        @on-error="handleUploadError"
        @on-success="getImageList"
      />
      <el-form ref="searchForm" :inline="true" :model="search">
        <el-form-item label="">
          <el-input v-model="search.keyword" class="keyword" placeholder="请输入文件名或备注" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="getImageList">查询</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="media">
      <div v-for="(item,key) in picList" :key="key" class="media-box">
        <div class="header-img-box-list">
          <!-- 选择按钮 -->
          <div 
            class="select-btn"
            :class="{ 'selected': isSelected(item) }"
            @click.stop="toggleSelect(item)"
          >
            <el-icon v-if="isSelected(item)">
              <check />
            </el-icon>
          </div>

          <el-image
            :key="key"
            :src="getUrl(item.url)"
            @click="chooseImg(item.url)"
          >
            <template #error>
              <div class="header-img-box-list">
                <el-icon>
                  <picture />
                </el-icon>
              </div>
            </template>
          </el-image>
        </div>
        <div class="img-title" @click="editFileNameFunc(item)">{{ item.name }}</div>
      </div>
    </div>
    <div class="bottom-btn-list">
      <el-dropdown @command="handleCommand">
        <el-button type="primary">
          批量操作<el-icon class="el-icon--right"><arrow-down /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="selectAll">选择</el-dropdown-item>
            <el-dropdown-item command="deleteSelected">删除</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :total.number="+total"
        :style="{'justify-content':'center'}"
        layout="total, prev, pager, next"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
      </div>
  </el-drawer>
</template>

<script setup>

import { getUrl } from '@/utils/image'
import { onMounted, ref, watch } from 'vue'
import { getFileList, editFileName, deleteFiles } from '@/api/fileUploadAndDownload'
import UploadImage from '@/components/upload/image.vue'
import UploadCommon from '@/components/upload/common.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check, Delete, FolderAdd, Picture } from '@element-plus/icons-vue'

const imageUrl = ref('')
const imageCommon = ref('')

const search = ref({})
const page = ref(1)
const total = ref(0)
const pageSize = ref(20)

const props = defineProps({
  modelValue: {
    type: [String, Array],
    default: ''
  },
  multiple: {
    type: Boolean,
    default: false
  },
  size: {
    type: String,
    default: 'default'
  }
})

const multipleValue = ref([])

// onMounted(() => {
//   if (props.multiple) {
//     multipleValue.value = props.modelValue
//   }
//   console.log("----[onMounted]------multipleValue:", multipleValue.value)
// })

watch(() => props.modelValue, (newVal) => {
  if (props.multiple) {
    if (Array.isArray(newVal) && newVal.length > 0) {
      multipleValue.value = newVal
    } 
  } else {
    imageUrl.value = newVal
  }
  // console.log("----[watch]------multipleValue:", multipleValue.value)
})

const emits = defineEmits(['update:modelValue'])

const deleteImg = (index) => {
  multipleValue.value.splice(index, 1)
  emits('update:modelValue', multipleValue.value)
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getImageList()
}

const handleCurrentChange = (val) => {
  page.value = val
  getImageList()
}
const editFileNameFunc = async(row) => {
  ElMessageBox.prompt('请输入文件名或者备注', '编辑', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /\S/,
    inputErrorMessage: '不能为空',
    inputValue: row.name
  }).then(async({ value }) => {
    row.name = value
    // console.log(row)
    const res = await editFileName(row)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '编辑成功!',
      })
      getImageList()
    }
  }).catch(() => {
    ElMessage({
      type: 'info',
      message: '取消修改'
    })
  })
}

const drawer = ref(false)
const picList = ref([])

const chooseImg = (url) => {
  if (props.multiple) {
    multipleValue.value.push(url)
    emits('update:modelValue', multipleValue.value)
  } else {
    emits('update:modelValue', url)
  }
  drawer.value = false
}
const openChooseImg = async() => {
  await getImageList()
  drawer.value = true
}

const getImageList = async() => {
  const res = await getFileList({ page: page.value, pageSize: pageSize.value, ...search.value })
  if (res.code === 0) {
    picList.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

const handleUploadError = async(res) => {
  console.log("-[handleUploadError]--error upload", res)
}

// 新增：存储选中的图片URL
const selectedUrls = ref([])

// 新增：切换图片选中状态
const toggleSelect = (file) => {
  if (selectedUrls.value.includes(file)) {
    // 取消选中
    selectedUrls.value = selectedUrls.value.filter(item => item.url !== file.url)
  } else {
    // 选中
    selectedUrls.value.push(file)
  }
}

// 新增：判断图片是否被选中
const isSelected = (file) => {
  return selectedUrls.value.includes(file)
}

const handleCommand = (command) => {
  if (command === 'selectAll') {
    // 全选
    multipleValue.value = selectedUrls.value.map(item => item.url)
    emits('update:modelValue', multipleValue.value)
    // selectedUrls.value = [] // 清空已选列表
  } else if (command === 'deleteSelected') {
    // 删除选中
    if (selectedUrls.value.length === 0) {
      ElMessage.warning('请先选择要删除的图片!');
      return;
    }
    ElMessageBox.confirm(`确定删除选中的 ${selectedUrls.value.length} 张图片吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      // 发送删除请求
      let selectedUrlsToDelete = selectedUrls.value.map(item => ({ id: item.id, key: item.key}));
      // 调用删除接口
      deleteFiles({filesInfo:selectedUrlsToDelete}).then(res => {
        if (res.code === 0) {
          ElMessage.success('删除成功!');
          // 更新图片列表
          getImageList();
          // 清空已选列表
          selectedUrls.value = [];
        } else {
          ElMessage.error('删除失败，请稍后重试。');
        }
      }).catch(() => {
        ElMessage.error('删除操作失败，请稍后重试。');
      }); 
    }).catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除'
      })
    });
  }
}
</script>
<script>

export default {
  name: 'SelectImage',
  methods: {

  }
}
</script>

<style scoped lang="scss">

.multiple-img{
  display: flex;
  gap:8px;
}

.add-image{
  width: 120px;
  height: 120px;
  line-height: 120px;
  display: flex;
  justify-content: center;
  border-radius: 20px;
  border: 1px dashed #ccc;
  background-size: cover;
  cursor: pointer;
}

.single-img-small {
  cursor: pointer;
  height: 100px;
  width: 100px;
  margin-left: 10px;
  background-repeat: no-repeat;
  background-size: cover;
  background-blend-mode: multiply, multiply;
}

.update-image {
  cursor: pointer;
  width: 120px;
  height: 120px;
  line-height: 120px;
  display: flex;
  justify-content: center;
  border-radius: 20px;
  border: 1px dashed #ccc;
   background-repeat: no-repeat;
   background-size: cover;
  &:hover {
    color: #fff;
    background: linear-gradient(
            to bottom,
            rgba(255, 255, 255, 0.15) 0%,
            rgba(0, 0, 0, 0.15) 100%
    ),
    radial-gradient(
            at top center,
            rgba(255, 255, 255, 0.4) 0%,
            rgba(0, 0, 0, 0.4) 120%
    )
    #989898;
    background-blend-mode: multiply, multiply;
    background-size: cover;
    .update {
      color: #fff;
    }
  }
  .update {
    height: 120px;
    width: 120px;
    text-align: center;
    color: transparent;
  }
}

.upload-btn-media-library {
  margin-left: 20px;
}

.media {
  display: flex;
  flex-wrap: wrap;
  

  .media-box {
    width: 120px;
    margin-left: 20px;
    margin-bottom: 20px; // 增加底部间距，避免内容重叠

    .img-title {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      line-height: 36px;
      text-align: center;
      cursor: pointer;
    }

    .header-img-box-list {
      width: 120px;
      height: 120px;
      border: 1px dashed #ccc;
      border-radius: 8px;
      text-align: center;
      line-height: 120px;
      position: relative;
      cursor: pointer;
      overflow: hidden;
      .el-image__inner {
        max-width: 120px;
        max-height: 120px;
        vertical-align: middle;
        width: unset;
        height: unset;
      }
    }
  }
}

// 选择按钮样式
.select-btn {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 2px solid #409eff;
  background-color: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0; // 默认隐藏
  transition: all 0.3s ease;
  color: #409eff;
  z-index: 10;
  opacity: 1;
  
  // 鼠标悬停在图片上时显示按钮
  //.header-img-box-list:hover & {
   // opacity: 1;
  //}
  
  // 选中状态样式
  &.selected {
    background-color: #409eff;
    color: white;
  }
}

.bottom-btn-list {
  display: flex;
  justify-content: space-between; /* 两端对齐 */
  align-items: center; /* 垂直居中 */
  padding: 16px; /* 可选：添加内边距 */
  width: 100%; /* 确保占满父容器宽度 */
  box-sizing: border-box; /* 避免padding影响总宽度 */
}

</style>
