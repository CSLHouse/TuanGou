<template>
  <div style="margin-top: 50px">
    <el-form :model="modelValue" ref="productAttrForm" label-width="120px" class="form-inner-container" size="small">
      <el-form-item label="属性类型：">
        <el-select v-model="modelValue.productAttributeCategoryId"
                   placeholder="请选择属性类型"
                   @change="handleProductAttrChange">
          <el-option
            v-for="item in productAttributeCategoryOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="商品规格：">
        <el-card shadow="never" class="cardBg">
          <div v-for="(productAttr,idx) in selectProductAttr">
            {{productAttr.name}}：
            <el-checkbox-group v-if="productAttr.handAddStatus===0" v-model="selectProductAttr[idx].values" @change="changeSelectProductAttr">
              <el-checkbox v-for="item in getInputListArr(productAttr.inputList)" :label="item" :key="item"
                           class="littleMarginLeft"></el-checkbox>
            </el-checkbox-group>
            <div v-else>
              <el-checkbox-group v-model="selectProductAttr[idx].values" @change="changeSelectProductAttr">
                <div v-for="(item,index) in selectProductAttr[idx].options" style="display: inline-block"
                     class="littleMarginLeft">
                  <el-checkbox :label="item" :key="item"></el-checkbox>
                  <el-button type="text" class="littleMarginLeft" @click="handleRemoveProductAttrValue(idx,index)">删除
                  </el-button>
                </div>
              </el-checkbox-group>
              <el-input v-model="addProductAttrValue" style="width: 160px;margin-left: 10px" clearable></el-input>
              <el-button class="littleMarginLeft" @click="handleAddProductAttrValue(idx)">增加</el-button>
            </div>
          </div>
        </el-card>
        <el-table style="width: 100%;margin-top: 20px"
                  :data="skuStockList"
                  border>
          <el-table-column
            v-for="(item,index) in selectProductAttr"
            :label="item.name"
            :key="item.id"
            align="center">
            <template #default="scope">
              {{getProductSkuSp(scope.row, index)}}
            </template>
          </el-table-column>
          <el-table-column
            label="销售价格"
            width="80"
            align="center">
            <template #default="scope">
              <el-input-number v-model="scope.row.price"  size="small" :controls="false" style="width: 70px;"/>
              <!-- <el-input v-model="scope.row.price" type="number"></el-input> -->
            </template>
          </el-table-column>
          <el-table-column
            label="促销价格"
            width="80"
            align="center">
            <template #default="scope">
              <!-- <el-input v-model="scope.row.promotionPrice" type="number"></el-input> -->
              <el-input-number v-model="scope.row.promotionPrice"  size="small" :controls="false" style="width: 66px;"/>
            </template>
          </el-table-column>
          <el-table-column
            label="商品库存"
            width="80"
            align="center">
            <template #default="scope">
              <el-input-number v-model="scope.row.stock"  size="small" :controls="false" style="width: 66px;"/>
              <!-- <el-input v-model="scope.row.stock" type="number"></el-input> -->
            </template>
          </el-table-column>
          <el-table-column
            label="库存预警值"
            width="100"
            align="center">
            <template #default="scope">
              <!-- <el-input v-model="scope.row.lowStock" type="number"></el-input> -->
              <el-input-number v-model="scope.row.lowStock"  size="small" :controls="false" style="width: 86px;"/>
            </template>
          </el-table-column>
          <el-table-column
            label="SKU编号"
            width="140"
            align="center">
            <template #default="scope">
              <el-input v-model="scope.row.skuCode"></el-input>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            width="80"
            align="center">
            <template #default="scope">
              <el-button
                type="text"
                @click="handleRemoveProductSku(scope.$index, scope.row)">删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-button
          type="primary"
          style="margin-top: 20px"
          @click="handleRefreshProductSkuList">刷新列表
        </el-button>
        <el-button
          type="primary"
          style="margin-top: 20px"
          @click="handleSyncProductSku">一键同步
        </el-button>
      </el-form-item>
      <el-form-item label="属性图片：" v-if="hasAttrPic">
        <el-card shadow="never" class="cardBg">
          <div v-for="(item,index) in selectProductAttrPics">
            <span>{{item.name}}:</span>
            <!-- <single-upload v-model="item.pic"
                           style="width: 300px;display: inline-block;margin-left: 10px"></single-upload> -->
            <SelectImage v-model="item.pic" />

          </div>
        </el-card>
      </el-form-item>
      <el-form-item label="商品参数：">
        <el-card shadow="never" class="cardBg">
          <div v-for="(item,index) in selectProductParam" :class="{littleMarginTop:index!==0}">
            <div class="paramInputLabel">{{item.name}}:</div>
            <el-select v-if="item.inputType===1" class="paramInput" v-model="selectProductParam[index].value">
              <el-option
                v-for="item1 in getParamInputList(item.inputList)"
                :key="item1"
                :label="item1"
                :value="item1">
              </el-option>
            </el-select>
            <el-input v-else class="paramInput" v-model="selectProductParam[index].value"></el-input>
          </div>
        </el-card>
      </el-form-item>
      <el-form-item label="商品相册：">
        <!-- <multi-upload v-if="selectProductPics.length > 0" v-model="selectProductPics" ></multi-upload> -->
        <!-- <cooller-image  v-model="selectProductPics"></cooller-image> -->
        <SelectImage v-model="selectProductPics" :multiple="true" />

      </el-form-item>
      <el-form-item label="商品详情：">
        <el-tabs v-model="activeHtmlName" type="card">
          <el-tab-pane label="电脑端详情" name="pc">
            <TEditor ref="editor" :width="595" :height="300" v-model="modelValue.detailHTML"></TEditor>
            <!-- <tinymce :width="595" :height="300" v-model="value.detailHtml"></tinymce> -->
          </el-tab-pane>
          <el-tab-pane label="移动端详情" name="mobile">
            <TEditor :width="595" :height="300" v-model="modelValue.detailMobileHTML"></TEditor>
            <!-- <tinymce :width="595" :height="300" v-model="value.detailMobileHtml"></tinymce> -->
          </el-tab-pane>
        </el-tabs>
      </el-form-item>
      <el-form-item style="text-align: center">
        <el-button size="medium" @click="handlePrev">上一步，填写商品促销</el-button>
        <el-button type="primary" size="medium" @click="handleNext">下一步，选择商品关联</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, toRefs, computed, watch, onMounted, getCurrentInstance } from 'vue';
import SingleUpload from '@/components/upload/singleUpload.vue';
import CoollerImage from '@/components/upload/coollerImage.vue'
import TEditor from '@/components/TEditor.vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getProductAttributeList, updateProductAttributeValue } from '@/api/product'
import { ProductStore } from '@/pinia/modules/product'
import SelectImage from '@/components/selectImage/selectImage.vue'
import { is } from 'core-js/core/object';
import { update } from '@/api/skuStock';
const productStore = ProductStore()

// 类型定义
interface ProductAttributeCategoryOption {
  label: string;
  value: number;
}

interface SelectProductAttr {
  id: number;
  name: string;
  handAddStatus: number;
  inputList: string;
  values: string[];
  options: string[];
}

interface SelectProductParam {
  id: number;
  name: string;
  value: string | null;
  inputType: number;
  inputList: string;
}

interface SelectProductAttrPic {
  name: string;
  pic: string | null;
}

interface SkuStock {
  spData: string;
  price?: number;
  promotionPrice?: number;
  stock?: number;
  lowStock?: number;
  skuCode?: string;
  pic?: string;
}

interface ProductAttributeValue {
  id?: number;
  productAttributeId: number;
  value: string;
}

interface SpDataItem {
  key: string;
  value: string;
}

interface ValueProps {
  id?: number;
  productAttributeCategoryId?: number;
  skuStockList: SkuStock[];
  productAttributeValueList: ProductAttributeValue[];
  pic?: string;
  albumPics?: string;
  detailHTML?: string;
  detailMobileHTML?: string;
}

// 组件Props
const props = defineProps<{
  modelValue: ValueProps;
  isEdit?: boolean;
}>();

const { modelValue, isEdit } = toRefs(props);

// 组件Emits
const emit = defineEmits<{
  (e: 'prevStep'): void;
  (e: 'nextStep'): void;
}>();

// 响应式变量
const hasEditCreated = ref(false);
const productAttributeCategoryOptions = ref<ProductAttributeCategoryOption[]>([]);
const selectProductAttr = ref<SelectProductAttr[]>([]);
const selectProductParam = ref<SelectProductParam[]>([]);
const selectProductAttrPics = ref<SelectProductAttrPic[]>([]);
const addProductAttrValue = ref('');
const activeHtmlName = ref('pc');
const skuStockList = ref<SkuStock[]>([]);


// 计算属性
const hasAttrPic = computed(() => {
  return selectProductAttrPics.value.length > 0;
});

const productId = computed(() => {
  return modelValue.value.id;
});

const selectProductPics = computed({
  get(): string[] {
    const pics: string[] = [];
    if (!modelValue.value.pic) return pics;
    
    pics.push(modelValue.value.pic);
    
    if (!modelValue.value.albumPics) return pics;
    
    const albumPics = modelValue.value.albumPics.split(',');
    pics.push(...albumPics);
    // console.log("----[selectProductPics]------pics:", pics);
    return pics;
  },
  set(newValue: any) {

    if (!newValue || newValue.length === 0) {
      modelValue.value.pic = '';
      modelValue.value.albumPics = '';
    } else {
      modelValue.value.pic = newValue[0];
      modelValue.value.albumPics = '';
      
      if (newValue.length > 1) {
         for (let i = 1; i < newValue.length; i++) {
          modelValue.value.albumPics += newValue[i];
          if (i !== newValue.length - 1) {
            modelValue.value.albumPics += ',';
          }
        }
        // modelValue.value.albumPics = newValue.slice(1).join(',');
      }
    }
  }
});

// 生命周期
onMounted(() => {
  getProductAttrCateList();
  // console.log("----[ProductAttrDetail]------modelValue:", modelValue.value);
});

// 监听商品ID变化
watch(productId, (newValue) => {
  if (!isEdit.value) return;
  if (hasEditCreated.value) return;
  if (!newValue) return;
  
  handleEditCreated();
  skuStockList.value = modelValue.value.skuStockList;

});

// 方法定义
const handleEditCreated = () => {
  if (modelValue.value.productAttributeCategoryId) {
    handleProductAttrChange(modelValue.value.productAttributeCategoryId);

  }
  hasEditCreated.value = true;
};

const getProductAttrCateList = async() => {
  await productStore.BuildProductAttributeData()
  let productAttributeCategoryList = productStore.ProductAttributeCategoryList
  for (let i = 0; i < productAttributeCategoryList.length; i++) {
    productAttributeCategoryOptions.value.push({
      label: productAttributeCategoryList[i].name, 
      value: productAttributeCategoryList[i].id
    });
  }

  // const param = { pageNum: 1, pageSize: 100 };
  // fetchProductAttrCateList(param).then(response => {
  //   productAttributeCategoryOptions.value = [];
  //   const list = response.data.list;
  //   for (let i = 0; i < list.length; i++) {
  //     productAttributeCategoryOptions.value.push({ 
  //       label: list[i].name, 
  //       value: list[i].id 
  //     });
  //   }
  // });
};

const changeSelectProductAttr = () => {
  // console.log("----[changeSelectProductAttr]------selectProductAttr:", selectProductAttr.value);
  refreshProductSkuList();
  refreshProductAttrPics();
};

const getProductAttrList = (type: number, cid: number) => {
  // const param = { pageNum: 1, pageSize: 100, type };
  getProductAttributeList({ tag: cid, page: 1, pageSize: 100, state: type }).then(response => {
    const list = response.data.list;
    if (type === 0) {
      selectProductAttr.value = [];
      for (let i = 0; i < list.length; i++) {
        let options: string[] = [];
        let values: string[] = [];
        
        options = getInputListArr(list[i].inputList);
        if (isEdit.value) {
          if (list[i].handAddStatus === 1) {
             let optionAdd = getEditAttrOptions(list[i].id);
             options = [...options, ...optionAdd];
            //  去重
            options = options.filter((element, i) => i === options.indexOf(element))
          }
          values = getEditAttrValues(i);
        } else {
          options = getInputListArr(list[i].inputList);
        }
        
        selectProductAttr.value.push({
          id: list[i].id,
          name: list[i].name,
          handAddStatus: list[i].handAddStatus,
          inputList: list[i].inputList,
          values,
          options
        });
      }
      if (isEdit.value) {
        refreshProductAttrPics();
      }

    } else {
      selectProductParam.value = [];
      for (let i = 0; i < list.length; i++) {
        let value: string | null = null;
        if (isEdit.value) {
          value = getEditParamValue(list[i].id);
        }
        
        selectProductParam.value.push({
          id: list[i].id,
          name: list[i].name,
          value,
          inputType: list[i].inputType,
          inputList: list[i].inputList
        });
      }
      // console.log("----[getProductAttrList]------selectProductParam:", selectProductParam.value);

    }
  });
};

const getEditAttrOptions = (id: number): string[] => {
  const options: string[] = [];
  for (let i = 0; i < modelValue.value.productAttributeValueList.length; i++) {
    const attrValue = modelValue.value.productAttributeValueList[i];
    if (attrValue.productAttributeId === id) {
      const strArr = attrValue.value.split(',');
      options.push(...strArr);
      break;
    }
  }
  return options;
};

const getEditAttrValues = (index: number): string[] => {
  const values = new Set<string>();
  
  for (let i = 0; i < modelValue.value.skuStockList.length; i++) {
    const sku = modelValue.value.skuStockList[i];
    const spData: SpDataItem[] = JSON.parse(sku.spData);
    
    if (spData && spData.length > index) {
      values.add(spData[index].value);
    }
  }
  
  return Array.from(values);
};

const getEditParamValue = (id: number): string | null => {
  for (let i = 0; i < modelValue.value.productAttributeValueList.length; i++) {
    if (id === modelValue.value.productAttributeValueList[i].productAttributeId) {
      return modelValue.value.productAttributeValueList[i].value;
    }
  }
  return null;
};

const handleProductAttrChange = (value: number) => {
  selectProductAttr.value = [];
  selectProductParam.value = [];
  selectProductAttrPics.value = [];
  skuStockList.value = [];

  getProductAttrList(0, value);
  getProductAttrList(1, value);
};

const getInputListArr = (inputList: string): string[] => {
  return inputList.split('\n');
};

const handleAddProductAttrValue = async(idx: number) => {
  const options = selectProductAttr.value[idx].options;
  
  if (!addProductAttrValue.value) {
    ElMessage({
      message: '属性值不能为空',
      type: 'warning',
      duration: 1000
    });
    return;
  }
  
  if (options.indexOf(addProductAttrValue.value) !== -1) {
    ElMessage({
      message: '属性值不能重复',
      type: 'warning',
      duration: 1000
    });
    return;
  }
  
  options.push(addProductAttrValue.value);
  addProductAttrValue.value = '';

  // console.log("----[handleAddProductAttrValue]------selectProductAttr:", selectProductAttr.value);
  // console.log("----[handleAddProductAttrValue]------modelValue:", modelValue.value);

  // let productAttribute = modelValue.value.productAttributeValueList.find(attr => attr.productAttributeId === selectProductAttr.value[idx].id);
  // if (productAttribute) {
  //   productAttribute.value += ',' + addProductAttrValue.value;
  //   console.log("----[handleAddProductAttrValue]------productAttribute:", productAttribute);
    
  //   let res = await updateProductAttributeValue(productAttribute);
  //   if ('code' in res && res.code === 0) {
  //     addProductAttrValue.value = '';
  //     ElMessage({
  //       message: '添加成功',
  //       type: 'success',
  //       duration: 1000
  //     });
  //   } else {
  //     ElMessage({
  //       message: '添加失败',
  //       type: 'error',
  //       duration: 1000
  //     });
  //   }
  // }
};

const handleRemoveProductAttrValue = (idx: number, index: number) => {
  selectProductAttr.value[idx].options.splice(index, 1);
};

const getProductSkuSp = (row: SkuStock, index: number): string | null => {
  const spData: SpDataItem[] = JSON.parse(row.spData);
  if (spData && index < spData.length) {
    return spData[index].value;
  }
  return null;
};

const handleRefreshProductSkuList = () => {
  ElMessageBox.confirm('刷新列表将导致sku信息重新生成，是否要刷新', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    refreshProductAttrPics();
    refreshProductSkuList();
  });
};

const handleSyncProductSku = () => {
  ElMessageBox.confirm('将同步第一个sku到所有sku,是否继续', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    if (skuStockList.value && skuStockList.value.length > 0) {
      const tempSkuList = [...skuStockList.value];
      const price = skuStockList.value[0].price;
      const promotionPrice = skuStockList.value[0].promotionPrice;
      const stock = skuStockList.value[0].stock;
      const lowStock = skuStockList.value[0].lowStock;
      tempSkuList.forEach(sku => {
        if (price !== undefined) {
          sku.price = Number(price);
        }
        if (promotionPrice !== undefined) {
          sku.promotionPrice = Number(promotionPrice);
        }
        if (stock !== undefined) {
          sku.stock = stock;
        }
        if (lowStock !== undefined) {
          sku.lowStock = lowStock;
        }
      });
      
      skuStockList.value = [...tempSkuList];
    }
  });
};

const refreshProductSkuList = () => {
  skuStockList.value = [];
  if (selectProductAttr.value.length === 1) {
    const attr = selectProductAttr.value[0];
    attr.values.forEach(value => {
      skuStockList.value.push({
        spData: JSON.stringify([{ key: attr.name, value }])
      });
    });
  } else if (selectProductAttr.value.length === 2) {
    const attr0 = selectProductAttr.value[0];
    const attr1 = selectProductAttr.value[1];
    
    attr0.values.forEach(val0 => {
      if (attr1.values.length === 0) {
        skuStockList.value.push({
          spData: JSON.stringify([{ key: attr0.name, value: val0 }])
        });
        return;
      }
      
      attr1.values.forEach(val1 => {
        const spData: SpDataItem[] = [
          { key: attr0.name, value: val0 },
          { key: attr1.name, value: val1 }
        ];
        skuStockList.value.push({
          spData: JSON.stringify(spData)
        });
      });
    });
  } else if (selectProductAttr.value.length >= 3) {
    const attr0 = selectProductAttr.value[0];
    const attr1 = selectProductAttr.value[1];
    const attr2 = selectProductAttr.value[2];
    
    attr0.values.forEach(val0 => {
      if (attr1.values.length === 0) {
        skuStockList.value.push({
          spData: JSON.stringify([{ key: attr0.name, value: val0 }])
        });
        return;
      }
      
      attr1.values.forEach(val1 => {
        if (attr2.values.length === 0) {
          const spData: SpDataItem[] = [
            { key: attr0.name, value: val0 },
            { key: attr1.name, value: val1 }
          ];
          skuStockList.value.push({
            spData: JSON.stringify(spData)
          });
          return;
        }
        
        attr2.values.forEach(val2 => {
          const spData: SpDataItem[] = [
            { key: attr0.name, value: val0 },
            { key: attr1.name, value: val1 },
            { key: attr2.name, value: val2 }
          ];
          skuStockList.value.push({
            spData: JSON.stringify(spData)
          });
        });
      });
    });
  } else {
    ElMessage({
      message: '不得多于三个商品属性',
      type: 'warning',
      duration: 1000
    });
  }

  if (isEdit.value && modelValue.value.skuStockList.length > 0) {
    skuStockList.value.forEach(sku => {
      const existingSku = modelValue.value.skuStockList.find(existing => existing.spData === sku.spData);
      if (existingSku) {
        sku.price = existingSku.price;
        sku.promotionPrice = existingSku.promotionPrice;
        sku.stock = existingSku.stock;
        sku.lowStock = existingSku.lowStock;
        sku.skuCode = existingSku.skuCode;
        sku.pic = existingSku.pic;
      }
    });
  }
};

const refreshProductAttrPics = () => {
  selectProductAttrPics.value = [];
  // console.log("----[refreshProductAttrPics]------selectProductAttr:", selectProductAttr.value);
  if (selectProductAttr.value.length >= 1) {
    const values = selectProductAttr.value[0].values;
    values.forEach(value => {
      let pic: string | null = null;
      if (isEdit.value) {
        pic = getProductSkuPic(value);
      }
      selectProductAttrPics.value.push({ name: value, pic });
    });
  }
  // console.log("----[refreshProductAttrPics]------selectProductAttrPics:", selectProductAttrPics.value);
};

const getProductSkuPic = (name: string): string | null => {
  for (let i = 0; i < skuStockList.value.length; i++) {
    const spData: SpDataItem[] = JSON.parse(skuStockList.value[i].spData);
    if (name === spData[0].value) {
      return skuStockList.value[i].pic || null;
    }
  }
  return null;
};

const mergeProductAttrValue = () => {  
  console.log("----[mergeProductAttrValue]------selectProductAttr:", selectProductAttr.value);
  console.log("----[mergeProductAttrValue]------selectProductParam:", selectProductParam.value);
  console.log("----[mergeProductAttrValue]------modelValue:", modelValue.value);
  selectProductAttr.value.forEach(attr => {
    if (attr.handAddStatus === 1 && attr.options && attr.options.length > 0) {
      if (isEdit.value) {
        // 编辑状态下，更新属性值
        let existingAttr = modelValue.value.productAttributeValueList.find(item => item.productAttributeId === attr.id);
        if (existingAttr) {
          existingAttr.value = getOptionStr(attr.values);
        } else {
          modelValue.value.productAttributeValueList.push({
            productAttributeId: attr.id,
            value: getOptionStr(attr.values)
          });
        }
      } else {
        // 非编辑状态下，直接添加属性值
        modelValue.value.productAttributeValueList.push({
          productAttributeId: attr.id,
          value: getOptionStr(attr.values)
        });
      }
    } else if (attr.handAddStatus === 0 && attr.values && attr.values.length > 0) {
      if (isEdit.value) {
        // 编辑状态下，更新属性值
        let existingAttr = modelValue.value.productAttributeValueList.find(item => item.productAttributeId === attr.id);
        if (existingAttr) {
          existingAttr.value = getOptionStr(attr.values);
          return;
        } else {
          // 如果不存在，则继续添加
          modelValue.value.productAttributeValueList.push({
            productAttributeId: attr.id,
            value: getOptionStr(attr.values)
          });
        }

      } else {
        // 非编辑状态下，直接添加属性值
        modelValue.value.productAttributeValueList.push({
          productAttributeId: attr.id,
          value: getOptionStr(attr.values)
        });
      }
    }

  });
  
  selectProductParam.value.forEach(param => {
    if (isEdit.value) {
      // 编辑状态下，更新参数值
      let existingParam = modelValue.value.productAttributeValueList.find(item => item.productAttributeId === param.id);
      if (existingParam) {
        existingParam.value = param.value || '';
        return;
      } else {
        // 如果不存在，则继续添加
        modelValue.value.productAttributeValueList.push({
          productAttributeId: param.id,
          value: param.value || ''
        });
      }
    } else {
      // 非编辑状态下，直接添加参数值
      modelValue.value.productAttributeValueList.push({
        productAttributeId: param.id,
        value: param.value || ''
      });
    }
    
  });
};

const mergeProductAttrPics = () => {
  selectProductAttrPics.value.forEach(picItem => {
    skuStockList.value.forEach(sku => {
      const spData: SpDataItem[] = JSON.parse(sku.spData);
      if (spData[0].value === picItem.name) {
        sku.pic = picItem.pic || '';
      }
    });
  });
};

const getOptionStr = (arr: string[]): string => {
  return arr.join(',');
};

const handleRemoveProductSku = (index: number, row) => {
  const list = skuStockList.value;
  if (list.length === 1) {
    list.pop();
  } else {
    list.splice(index, 1);
  }
};

const getParamInputList = (inputList: string): string[] => {
  return inputList.split('\n');
};

const handlePrev = () => {
  emit('prevStep');
};

const handleNext = () => {
  mergeProductAttrValue();
  mergeProductAttrPics();
  modelValue.value.skuStockList = skuStockList.value;
  console.log("----[ProductAttrDetail]------modelValue:", modelValue.value);
  emit('nextStep');
};

// 暴露给模板使用的变量和方法
defineExpose({
  handlePrev,
  handleNext
});
</script>


<style scoped>
  .littleMarginLeft {
    margin-left: 10px;
  }

  .littleMarginTop {
    margin-top: 10px;
  }

  .paramInput {
    width: 250px;
  }

  .paramInputLabel {
    display: inline-block;
    width: 100px;
    text-align: right;
    padding-right: 10px
  }

  .cardBg {
    background: #F8F9FC;
  }
</style>
