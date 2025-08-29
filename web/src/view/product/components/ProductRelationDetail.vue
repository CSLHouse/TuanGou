<template>
  <div style="margin-top: 50px">
    <el-form :model="modelValue"
             ref="productRelationForm"
             label-width="120px"
             class="form-inner-container"
             size="small">
      <!-- <el-form-item label="关联专题：">
        <el-transfer
          style="display: inline-block"
          filterable
          :filter-method="filterMethod"
          filter-placeholder="请输入专题名称"
          v-model="selectSubject"
          :titles="subjectTitles"
          :data="subjectList">
        </el-transfer>
      </el-form-item>
      <el-form-item label="关联优选：">
        <el-transfer
          style="display: inline-block"
          filterable
          :filter-method="filterMethod"
          filter-placeholder="请输入优选名称"
          v-model="selectPrefrenceArea"
          :titles="prefrenceAreaTitles"
          :data="prefrenceAreaList">
        </el-transfer>
      </el-form-item> -->
      <el-form-item style="text-align: center">
        <el-button size="small" @click="handlePrev">上一步，填写商品属性</el-button>
        <el-button type="primary" size="small" @click="handleFinishCommit">完成，提交商品</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
// import {fetchListAll as fetchSubjectList} from '@/api/subject'
// import {fetchList as fetchPrefrenceAreaList} from '@/api/prefrenceArea'
import { ref, computed, onBeforeMount, watch, toRefs, PropType } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

let prop = defineProps({
  modelValue: Object as () => ({
    id: number,
    productAttributeCategoryId: number,
    pic: string,
    albumPics: string,
    detailHTML: string,
    detailMobileHTML: string,
  }),
  isEdit: Boolean
})
// 要响应式
const {modelValue, isEdit} = toRefs(prop);


//所有专题列表
const subjectList = ref([])
//专题左右标题
const subjectTitles = ['待选择', '已选择']
//所有专题列表
const prefrenceAreaList = ref([])
//专题左右标题
const prefrenceAreaTitles = ['待选择', '已选择']

onBeforeMount(() => {
  getSubjectList();
   getPrefrenceAreaList();
})

//选中的专题
const selectSubject = computed({
  get() {
    let subjects =[];
    if(this.value.subjectProductRelationList==null||this.value.subjectProductRelationList.length<=0){
      return subjects;
    }
    for(let i=0;i<this.value.subjectProductRelationList.length;i++){
      subjects.push(this.value.subjectProductRelationList[i].subjectId);
    }
    return subjects;
  },
  set(newValue) {
    this.value.subjectProductRelationList=[];
    for(let i=0;i<newValue.length;i++){
      this.value.subjectProductRelationList.push({subjectId:newValue[i]});
    }
  }
})

//选中的优选
const selectPrefrenceArea = computed({
  get() {
    let prefrenceAreas =[];
    if(this.value.prefrenceAreaProductRelationList==null||this.value.prefrenceAreaProductRelationList.length<=0){
      return prefrenceAreas;
    }
    for(let i=0;i<this.value.prefrenceAreaProductRelationList.length;i++){
      prefrenceAreas.push(this.value.prefrenceAreaProductRelationList[i].prefrenceAreaId);
    }
    return prefrenceAreas;
  },
  set(newValue) {
    this.value.prefrenceAreaProductRelationList=[];
    for(let i=0;i<newValue.length;i++){
      this.value.prefrenceAreaProductRelationList.push({prefrenceAreaId:newValue[i]});
    }
  }
})

const filterMethod = (query, item) => {
  return item.label.indexOf(query) > -1;
}
const getSubjectList = () => {
  // fetchSubjectList().then(response => {
  //   let list = response.data;
  //   for (let i = 0; i < list.length; i++) {
  //     subjectList.value.push({
  //       label: list[i].title,
  //       key: list[i].id
  //     });
  //   }
  // });
}
const getPrefrenceAreaList = () => {
  // fetchPrefrenceAreaList().then(response=>{
  //   let list = response.data;
  //   for (let i = 0; i < list.length; i++) {
  //     prefrenceAreaList.value.push({
  //       label: list[i].name,
  //       key: list[i].id
  //     });
  //   }
  // });
}

const emits = defineEmits(["prevStep", "finishCommit"]);
const handlePrev = () => {
  emits('prevStep')
}
const handleFinishCommit = () => {
  emits('finishCommit', isEdit.value);
}

</script>


<style scoped>

</style>
