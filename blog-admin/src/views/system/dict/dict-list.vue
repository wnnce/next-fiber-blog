<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import { constant } from '@/assets/script/constant'
import type { Dict, DictQueryForm } from '@/api/system/dict/types'
import { dictApi } from '@/api/system/dict'
import DictForm from '@/views/system/dict/dict-form.vue'
import DictValueDrawer from '@/views/system/dict/dict-value-drawer.vue'

const { successMessage, loading, errorMessage } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<Dict[]>([]);
const defaultQueryForm: DictQueryForm = {
  page: 1,
  size: 10,
  dictName: undefined,
  dictKey: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<DictQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await dictApi.pageDict(queryForm);
    if (result && result.code === 200) {
      const { total, records } = result.data;
      recordTotal.value = total;
      tableData.value = records;
    }
  } finally {
    tableLoading.value = false;
  }
}

const queryLoading = ref<boolean>(false);
const handleQuery = async () => {
  queryForm.page = 1;
  queryLoading.value = true;
  try {
    await queryTableData();
  } finally {
    queryLoading.value = false;
  }
}

const handleReset = () => {
  Object.assign(queryForm, defaultQueryForm);
  dateRange.value = [];
}

const handleDateChange = () => {
  if (!dateRange.value) {
    queryForm.createTimeBegin = undefined;
    queryForm.createTimeEnd = undefined;
  } else {
    const [ begin, end ] = dateRange.value;
    queryForm.createTimeBegin = begin;
    queryForm.createTimeEnd = end;
  }
}

const handleDelete = async (record: Dict) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await dictApi.deleteDict(record.dictId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const handleUpdateStatus = async (newValue: number, id: number) => {
  const result = await dictApi.updateDictStatus({ dictId: id, status: newValue })
  if (result.code === 200) {
    successMessage('更新成功');
    return true;
  }
  return false;
}

const formRef = ref();
const showForm = (record?: Dict) => {
  formRef.value.show(record);
}

const dictValueRef = ref();
const showDictValue = (record: Dict) => {
  const { dictId, dictKey, status } = record;
  if (status === 1) {
    errorMessage('禁用状态下无法查看字典详情');
    return;
  }
  dictValueRef.value.show(dictId, dictKey);
}

onMounted(() => {
  queryTableData();
})
</script>

<template>
  <div class="table-card">
    <div class="search-div">
      <div class="search-item">
        <label>字典名称</label>
        <a-input v-model="queryForm.dictName" placeholder="请输入字典名称" />
      </div>
      <div class="search-item">
        <label>字典KEY</label>
        <a-input v-model="queryForm.dictKey" placeholder="请输入字典key" />
      </div>
      <div class="search-item">
        <label>创建时间</label>
        <a-range-picker allow-clear v-model="dateRange" @change="handleDateChange" />
      </div>
      <div class="search-buttons">
        <a-button type="primary" @click="handleQuery" :loading="queryLoading">
          <template #icon><icon-search /></template>
          搜索
        </a-button>
        <a-button @click="handleReset">
          <template #icon><icon-refresh /></template>
          重置
        </a-button>
      </div>
    </div>
    <div class="flex justify-between">
      <div class="flex" style="column-gap: 12px">
        <a-button type="primary" @click="showForm(undefined)">
          <template #icon><icon-plus /></template>
          新增
        </a-button>
      </div>
      <right-operate @refresh="queryTableData"/>
    </div>
    <a-table :data="tableData" :loading="tableLoading" :pagination="false">
      <template #columns>
        <a-table-column title="字典ID" data-index="dictId" />
        <a-table-column title="字典名称" data-index="dictName" />
        <a-table-column title="字典KEY">
          <template #cell="{ record }">
            <span class="link-text" @click="showDictValue(record)">
              {{ record.dictKey }}
            </span>
          </template>
        </a-table-column>
        <a-table-column title="排序" data-index="sort" />
        <a-table-column title="创建时间" data-index="createTime" />
        <a-table-column title="备注" data-index="remark" />
        <a-table-column title="状态">
          <template #cell="{ record }">
            <a-switch :checked-value="0" :unchecked-value="1" v-model="record.status"
                      :before-change="(newValue: number) => handleUpdateStatus(newValue, record.dictId)"
            />
          </template>
        </a-table-column>
        <a-table-column title="操作" align="center">
          <template #cell="{ record }">
            <a-button type="text" shape="circle" @click="showForm(record)">
              <template #icon><icon-edit /></template>
            </a-button>
            <a-popconfirm content="是否确认删除数据？" type="error" position="lt"
                          :ok-button-props="{ status: 'danger' }"
                          @ok="handleDelete(record)"
            >
              <a-button type="text" shape="circle" status="danger">
                <template #icon><icon-delete /></template>
              </a-button>
            </a-popconfirm>
          </template>
        </a-table-column>
      </template>
    </a-table>
    <div class="flex justify-end">
      <a-pagination :total="recordTotal" size="medium"
                    v-model:current="queryForm.page"
                    v-model:page-size="queryForm.size"
                    :page-size-options="constant.pageSizeOption"
                    show-page-size
                    show-total
                    @page-size-change="handleQuery"
                    @change="queryTableData"
      />
    </div>
    <dict-form ref="formRef" @reload="queryTableData"/>
    <dict-value-drawer ref="dictValueRef" />
  </div>
</template>

<style scoped lang="scss">

</style>