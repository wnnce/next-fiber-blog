<script setup lang="ts">

import { reactive, ref } from 'vue'
import type { Config, ConfigQueryForm } from '@/api/system/config/types'
import { configApi } from '@/api/system/config'
import RightOperate from '@/components/RightOperate.vue'

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<Config[]>([]);
const defaultQueryForm: ConfigQueryForm = {
  page: 1,
  size: 10,
  name: undefined,
  key: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<ConfigQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await configApi.pageSysConfig(queryForm);
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
const handleQuery = () => {
  queryForm.page = 1;
  queryLoading.value = true;
  try {
    queryTableData();
  } finally {
    queryLoading.value = false;
  }
}

const handleReset = () => {
  Object.assign(queryForm, defaultQueryForm);
  dateRange.value = [];
}

const handleDateChange = () => {
  const [ begin, end ] = dateRange.value;
  queryForm.createTimeBegin = begin;
  queryForm.createTimeEnd = end;
}

</script>

<template>
  <div class="card flex flex-column" style="row-gap: 12px">
    <div class="search-div">
      <div class="search-item">
        <label>参数名称</label>
        <a-input v-model="queryForm.name" placeholder="请输入参数名称" />
      </div>
      <div class="search-item">
        <label>参数KEY</label>
        <a-input v-model="queryForm.key" placeholder="请输入参数key" />
      </div>
      <div class="search-item">
        <label>创建时间</label>
        <a-range-picker v-model="dateRange" @change="handleDateChange" />
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
        <a-button type="primary">
          <template #icon><icon-plus /></template>
          新增
        </a-button>
      </div>
      <right-operate @refresh="queryTableData"/>
    </div>
    <a-table :data="tableData" >
      <template #columns>
        <a-table-column title="参数名称" data-index="configName" />
        <a-table-column title="参数键" data-index="configKey" />
        <a-table-column title="参数值" data-index="configValue" />
        <a-table-column title="创建时间" data-index="createTime" />
        <a-table-column title="备注" data-index="remark" />
        <a-table-column title="操作">
          <template #cell="{ record }">
            {{ record }}
          </template>
        </a-table-column>
      </template>
    </a-table>
  </div>
</template>

<style scoped lang="scss">

</style>