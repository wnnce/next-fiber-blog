<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { constant } from '@/assets/script/constant'
import TagsForm from '@/views/blog/tags/tags-form.vue'
import type { AccessRecord, AccessRecordQueryForm } from '@/api/system/record/types'
import { recordApi } from '@/api/system/record'

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<AccessRecord[]>([]);
const defaultQueryForm: AccessRecordQueryForm = {
  page: 1,
  size: 10,
  ip: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<AccessRecordQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await recordApi.pageAccessRecord(queryForm);
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

const handleRefresh = () => {
  queryTableData();
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

onMounted(() => {
  queryTableData();
})
</script>

<template>
  <div class="table-card">
    <div class="search-div">
      <div class="search-item">
        <label>IP地址</label>
        <a-input v-model="queryForm.ip" placeholder="通过IP地址搜索" />
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
        <a-button status="danger" disabled>
          <template #icon>
            <icon-delete />
          </template>
          删除
        </a-button>
        <a-button status="danger" disabled>
          <template #icon>
            <icon-refresh />
          </template>
          清空
        </a-button>
      </div>
      <right-operate @refresh="handleRefresh"/>
    </div>
    <a-table :data="tableData" :loading="tableLoading" :pagination="false">
      <template #columns>
        <a-table-column title="记录ID" data-index="id" :width="120"/>
        <a-table-column title="来源" data-index="referee" align="center" />
        <a-table-column title="访问IP" data-index="accessIp" align="center" />
        <a-table-column title="地址" data-index="location" align="center" />
        <a-table-column title="UA" ellipsis align="center">
          <template #cell="{ record }">
            <a-tooltip :content="record.accessUa">
              <p>{{ record.accessUa }}</p>
            </a-tooltip>
          </template>
        </a-table-column>
        <a-table-column title="访问时间" data-index="createTime" :width="260" align="center" />
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
    <tags-form ref="formRef" @reload="handleQuery" />
  </div>
</template>

<style scoped lang="scss">

</style>