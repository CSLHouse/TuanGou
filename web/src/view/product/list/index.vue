<template>
  <div class="app-container">
    <el-card class="filter-container" shadow="never">
      <div>
        <!-- <i class="el-icon-search"></i> -->
        <span>筛选搜索</span>
        <el-button
          style="float: right"
          @click="handleSearchList()"
          type="primary"
          size="small">
          查询结果
        </el-button>
        <el-button
          style="float: right;margin-right: 15px"
          @click="handleResetSearch()"
          size="small">
          重置
        </el-button>
      </div>
      <div style="margin-top: 15px">
        <el-form :inline="true" :model="listQuery" size="small" label-width="140px">
          <el-form-item label="输入搜索：">
            <el-input style="width: 203px" v-model="listQuery.keyword" placeholder="商品名称"></el-input>
          </el-form-item>
          <el-form-item label="商品货号：">
            <el-input style="width: 203px" v-model="listQuery.productSN" placeholder="商品货号"></el-input>
          </el-form-item>
          <el-form-item label="商品分类：">
            <el-cascader
              clearable
              v-model="listQuery.productCategoryId"
              :options="productCateOptions"
              placeholder="请选择" @change="handleCategoryIdChange">
            </el-cascader>
          </el-form-item>
          <el-form-item label="商品品牌：">
            <el-select v-model="listQuery.brandId" placeholder="请选择品牌" clearable>
              <el-option
                v-for="item in brandOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="上架状态：">
            <el-select v-model="listQuery.publishStatus" placeholder="全部" clearable>
              <el-option
                v-for="item in publishStatusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="审核状态：">
            <el-select v-model="listQuery.verifyStatus" placeholder="全部" clearable>
              <el-option
                v-for="item in verifyStatusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
    <el-card class="operate-container" shadow="never">
      <i class="el-icon-tickets"></i>
      <span>数据列表</span>
      <el-button
        class="btn-add"
        icon="plus"
        @click="handleAddProduct()"
        size="small">
        添加
      </el-button>
    </el-card>
    <div class="table-container">
      <el-table ref="productTable"
                :data="productList"
                style="width: 100%"
                @selection-change="handleSelectionChange"
                v-loading="listLoading"
                border>
        <el-table-column type="selection" width="60" align="center"></el-table-column>
        <el-table-column label="编号" width="100" align="center">
          <template #default="scope">{{scope.row.id}}</template>
        </el-table-column>
        <el-table-column label="商品图片" width="120" align="center">
          <template #default="scope"><img style="height: 80px" :src="scope.row.pic"></template>
        </el-table-column>
        <el-table-column label="商品名称" align="center">
          <template #default="scope">
            <p>{{scope.row.name}}</p>
            <p>品牌：{{scope.row.brandName}}</p>
          </template>
        </el-table-column>
        <el-table-column label="价格/货号" width="120" align="center">
          <template #default="scope">
            <p>价格：￥{{scope.row.price}}</p>
            <p>货号：{{scope.row.productSN}}</p>
          </template>
        </el-table-column>
        <el-table-column label="标签" width="140" align="center">
          <template #default="scope">
            <p>上架：
              <el-switch
                @change="handlePublishStatusChange(scope.$index, scope.row)"
                :active-value="1"
                :inactive-value="0"
                v-model="scope.row.publishStatus">
              </el-switch>
            </p>
            <p>新品：
              <el-switch
                @change="handleNewStatusChange(scope.$index, scope.row)"
                :active-value="1"
                :inactive-value="0"
                v-model="scope.row.newStatus">
              </el-switch>
            </p>
            <p>推荐：
              <el-switch
                @change="handleRecommendStatusChange(scope.$index, scope.row)"
                :active-value="1"
                :inactive-value="0"
                v-model="scope.row.recommandStatus">
              </el-switch>
            </p>
          </template>
        </el-table-column>
        <el-table-column label="排序" width="100" align="center">
          <template #default="scope">{{scope.row.sort}}</template>
        </el-table-column>
        <el-table-column label="SKU库存" width="100" align="center">
          <template #default="scope">
            <el-button type="primary" :icon="Edit" circle @click="handleShowSkuEditDialog(scope.row)" />
          </template>
        </el-table-column>
        <el-table-column label="销量" width="100" align="center">
          <template #default="scope">{{scope.row.sale}}</template>
        </el-table-column>
        <el-table-column label="审核状态" width="100" align="center">
          <template #default="scope">
            <p>{{ verifyStatusFilter(scope.row.verifyStatus)}}</p>
            <p>
              <el-button
                type="primary" link
                @click="handleShowVerifyDetail(scope.$index, scope.row)">审核详情
              </el-button>
            </p>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" align="center">
          <template #default="scope">
            <p>
              <el-button
                size="small"
                @click="handleShowProduct(scope.$index, scope.row)">查看
              </el-button>
              <el-button
                size="small"
                @click="handleUpdateProduct(scope.$index, scope.row)">编辑
              </el-button>
            </p>
            <p>
              <el-button
                size="small"
                @click="handleShowLog(scope.$index, scope.row)">日志
              </el-button>
              <el-popover v-model="scope.row.visible" placement="top" width="160">
                <p>确定要删除吗？</p>
                <div style="text-align: right; margin-top: 8px;">
                  <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                  <el-button type="primary" @click="handleDelete(scope.row)">确定</el-button>
                </div>
                <template #reference>
                  <el-button type="danger" link icon="delete" @click="scope.row.visible = true">删除</el-button>
                </template>
              </el-popover>
            </p>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div class="batch-operate-container">
      <el-select
        size="small"
        v-model="operateType" placeholder="批量操作">
        <el-option
          v-for="item in operates"
          :key="item.key"
          :label="item.label"
          :value="item.key">
        </el-option>
      </el-select>
      <el-button
        style="margin-left: 20px"
        class="search-button"
        @click="handleBatchOperate()"
        type="primary"
        size="small">
        确定
      </el-button>
    </div>
    <div class="pagination-container">
      <el-pagination
        background
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        layout="total, sizes,prev, pager, next,jumper"
        :page-size="listQuery.pageSize"
        :page-sizes="[5,10,15]"
        :current-page.sync="listQuery.page"
        :total.number="+total">
      </el-pagination>
    </div>
    <el-dialog title="编辑货品信息" v-model="editSkuInfo.dialogVisible" width="40%">
      <span>商品货号：</span>
      <span>{{editSkuInfo.productSN}}</span>
      <el-input placeholder="按sku编号搜索" v-model="editSkuInfo.keyword" size="small" 
        style="width: 50%;margin-left: 20px"
        clearable
        @clear="handleClearSkuCode"
      >
        <template #append>
          <el-button :icon="Search" @click="handleSearchEditSku"/>
        </template>
      </el-input>
      <el-table style="width: 100%;margin-top: 20px"
                :data="editSkuInfo.stockList"
                border>
        <el-table-column
          label="SKU编号"
          align="center">
          <template #default="scope">
            <el-input v-model="scope.row.skuCode"></el-input>
          </template>
        </el-table-column>
        <el-table-column
          v-for="(item,index) in editSkuInfo.productAttr"
          :label="item.name"
          :key="item.id"
          align="center">
          <template #default="scope">
            {{getProductSkuSp(scope.row,index)}}
          </template>
        </el-table-column>
        <el-table-column
          label="销售价格"
          width="80"
          align="center">
          <template #default="scope">
            <el-input v-model="scope.row.price"></el-input>
          </template>
        </el-table-column>
        <el-table-column
          label="商品库存"
          width="80"
          align="center">
          <template #default="scope">
            <el-input v-model="scope.row.stock"></el-input>
          </template>
        </el-table-column>
        <el-table-column
          label="库存预警值"
          width="100"
          align="center">
          <template #default="scope">
            <el-input v-model="scope.row.lowStock"></el-input>
          </template>
        </el-table-column>
      </el-table>
      <span slot="footer" class="dialog-footer">
        <el-button @click="editSkuInfo.dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleEditSkuConfirm">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import {
  getProductList,
  updateProductKeyword,
  getProductSKUStockByProductID,
  getProductAttributeList,
  updateProductSKUStock,
  deleteProducts,
} from '@/api/product'
import { ref, computed, onBeforeMount, watch } from 'vue'
import { ElMessage, ElMessageBox, ElTable } from 'element-plus'
import { Edit, Search } from '@element-plus/icons-vue'
import { ProductStore } from '@/pinia/modules/product'
import { useRouter } from "vue-router";
const router = useRouter()

const productStore = ProductStore()
const defaultListQuery = {
  keyword: null,
  page: 1,
  pageSize: 5,
  publishStatus: null,
  verifyStatus: null,
  productSN: null,
  productCategoryId: null,
  brandId: null
};
const listQuery = ref(defaultListQuery)
const productCateOptions = ref()
const brandOptions = ref()

onBeforeMount(() => {
  getList()
  getProductCategoryData()
  getProductBrandData()
})
const getProductCategoryData = async() => {
  await productStore.BuildProductCategoryData()
  productCateOptions.value = productStore.ProductCategoryOptions
}

const getProductBrandData = async() => {
  await productStore.BuildBrandData()
  let productBrandList = productStore.RandData['list']
  
  brandOptions.value = productBrandList.map((item) => {
    return {value: item.id, label: item.name}
  })
}

const handleCategoryIdChange = () => {
  listQuery.value.productCategoryId = listQuery.value.productCategoryId.at(-1)
}

const handleSearchList = async() => {
  listQuery.value.page = 1;
  getList();
}
const handleResetSearch = async() => {
  listQuery.value = defaultListQuery
}

// const page = ref(1)
const total = ref(0)
// const pageSize = ref(5)
const productList = ref([])
const listLoading = ref(true)

// 分页
const handleSizeChange = (val) => {
  listQuery.value.pageSize = val
  getList()
}

const handleCurrentChange = (val) => {
  console.log("handleCurrentChange", val)
  listQuery.value.page = val
  getList()
}

// 查询
const getList = async() => {
  listLoading.value = true;
  const res = await getProductList(listQuery.value)
  if ('code' in res && res.code === 0) {
    listLoading.value = false;
    productList.value = res.data.list
    total.value = res.data.total;
  }
}

const publishStatusOptions = [
  {
      value: 101,
      label: '上架'
  },
  {
      value: 100,
      label: '下架'
  },
]
const verifyStatusOptions = [
  {
      value: 101,
      label: '审核通过'
  },
  {
      value: 100,
      label: '未审核'
  },
]

const handleAddProduct = () => {
  router.push({path: '/layout/product/add'});
}

const multipleSelection = ref()
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const handlePublishStatusChange = async(_, row) => {
  let ids = [];
  ids.push(row.id);
  updateProductKeyword({ids: ids, key: 'publish_status', value: row.publishStatus })
}

const handleNewStatusChange = async(_, row) => {
  let ids = [];
  ids.push(row.id);
  updateProductKeyword({ids: ids, key: 'new_status', value: row.newStatus })
}

const handleRecommendStatusChange = async(_, row) => {
  let ids = [];
  ids.push(row.id);
  updateProductKeyword({ids: ids, key: 'recommand_status', value: row.recommandStatus })
}

const editSkuInfo = ref({
  dialogVisible:false,
  productId:null,
  productSN:'',
  productAttributeCategoryId:null,
  stockList:[],
  productAttr:[],
  keyword:null
})

const handleShowSkuEditDialog = async(row) => {
  editSkuInfo.value.dialogVisible = true;
  editSkuInfo.value.productId = row.id;
  editSkuInfo.value.productSN = row.productSN;
  editSkuInfo.value.productAttributeCategoryId = row.productAttributeCategoryId;
  editSkuInfo.value.keyword = null;
  const res = await getProductSKUStockByProductID({id: row.id})
  if ('code' in res && res.code === 0) {
    editSkuInfo.value.stockList = res.data
  } 
  
  if(row.productAttributeCategoryId != null){
    const attributeRes = await getProductAttributeList({tag: row.productAttributeCategoryId, state: 0})
    if ('code' in attributeRes && attributeRes.code === 0) {
      editSkuInfo.value.productAttr = attributeRes.data.list
    }
  }
}

// 审核状态过滤器
const verifyStatusFilter = computed(() => {
  return (value: any) => {
      if (value === 1) {
        return '审核通过';
      } else {
        return '未审核';
      }
  }
})
// 审核详情点击事件
const handleShowVerifyDetail = (index,row) => {
  console.log("handleShowVerifyDetail",row);
}

// 查看按钮点击事件
const handleShowProduct = (index,row) => {
  console.log("handleShowProduct",row);
}

const handleUpdateProduct = (index,row) => {
  router.push({path: '/layout/product/update', query:{id: row.id}});
}

const handleShowLog = (index,row) => {
  console.log("handleShowLog",row);
}

const handleDelete = async( row) => {
  const res = await deleteProducts({ids: [row.id]})
  if ('code' in res && res.code === 0) {
    getList()
    ElMessage({
      message: '删除成功',
      type: 'success',
      duration: 1000
    });
  }
}

const operateType = ref()
const operates = [
  {
    label: "商品上架",
    key: "publishOn",
    dbKey: "publish_status", 
    value: 1
  },
  {
    label: "商品下架",
    key: "publishOff",
    dbKey: "publish_status",
    value: 0
  },
  {
    label: "设为推荐",
    key: "recommendOn",
    dbKey: "recommand_status",
    value: 1
  },
  {
    label: "取消推荐",
    key: "recommendOff",
    dbKey: "recommand_status",
    value: 0
  },
  {
    label: "设为新品",
    key: "newOn",
    dbKey: "new_status",
    value: 1
  },
  {
    label: "取消新品",
    key: "newOff",
    dbKey:"new_status",
    value: 0
  },
  {
    label: "移入回收站",
    key: "recycle",
    dbKey: "",
    value: 0
  }
]

var updateList:number[] = new Array() 
const handleBatchOperate = async() => {
  if (operateType.value == null) {
    ElMessage({
      message: '请选择操作类型',
      type: 'warning',
      duration: 1000
    });
    return
  }
  if (!multipleSelection.value || multipleSelection.value.length < 1){
    ElMessage({
      message: '请选择要操作的商品',
      type: 'warning',
      duration: 1000
    });
    return
  }
  
  multipleSelection.value.forEach((item) => {
      updateList.push(item.id)
  })
  if (operateType.value == 'recycle') {
    const res = await deleteProducts({ids: updateList})
    if ('code' in res && res.code === 0) {
      ElMessage({
        message: '删除成功',
        type: 'success',
        duration: 1000
      });
      getList()
    }
  } else {
    let dbkey
    let value
    for (let index = 0; index < operates.length; index++) {
      const element = operates[index];
      if (element.key == operateType.value) {
        dbkey = element.dbKey
        value = element.value
        break
      }
    }
    const res = await updateProductKeyword({ids: updateList, key: dbkey, value: value })
    if ('code' in res && res.code === 0) {
      ElMessage({
        message: '修改成功',
        type: 'success',
        duration: 1000
      });
      getList()
    }
  }
  updateList = []
}

const handleClearSkuCode = () => {
  editSkuInfo.value.keyword = null
}

const handleSearchEditSku = async (params) => {
  editSkuInfo.value.stockList = []
  const res = await getProductSKUStockByProductID({id: editSkuInfo.value.productId, keyword: editSkuInfo.value.keyword})
  if ('code' in res && res.code === 0) {
    editSkuInfo.value.stockList = res.data
  }
}

const getProductSkuSp = (row, index) => {
  let spData = JSON.parse(row.spData);
  if(spData !=null && index < spData.length){
    return spData[index].value;
  }else{
    return null;
  }
}

const handleEditSkuConfirm = async() => {
  if(editSkuInfo.value.stockList == null||editSkuInfo.value.stockList.length <= 0){
    ElMessage({
      message: '暂无sku信息',
      type: 'warning',
      duration: 1000
    });
    return
  }
  ElMessageBox.confirm('是否要进行修改', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    const res = updateProductSKUStock(editSkuInfo.value.stockList)
    if ('code' in res && res.code !== 0) {
      ElMessage({
        type: 'success',
        message: '更新成功'
      })
      editSkuInfo.value.dialogVisible = false;
    }
  })
}

</script>
<style></style>


