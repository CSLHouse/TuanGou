<template>
  <div id="app">
    <el-config-provider :locale="zhCn">
      <router-view :key="route.path"/>
    </el-config-provider>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<script setup>
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
// element 2.3.8之前使用下面的语句
// import zhCn from 'element-plus/lib/locale/lang/zh-cn'
// TODO:解决动态路由参数变化页面不刷新的问题
import { useRoute } from 'vue-router';
const route = useRoute()
</script>
<style lang="scss">
@tailwind base;
@tailwind components;
@tailwind utilities;
// 引入初始化样式
#app {
  background: #eee;
  height: 100vh;
  overflow: hidden;
  font-weight: 400 !important;
}
.el-button{
  font-weight: 400 !important;
}
</style>
