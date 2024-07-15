<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import LoadImage from '@/components/LoadImage.vue'
import type { Concat, ConcatQueryForm } from '@/api/blog/concat/types'
import { concatApi } from '@/api/blog/concat'
import ConcatForm from '@/views/blog/concat/concat-form.vue'

const { successMessage, loading } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const tableData = ref<Concat[]>([]);
const defaultQueryForm: ConcatQueryForm = {
  name: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<ConcatQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await concatApi.manageList(queryForm);
    if (result && result.code === 200) {
      tableData.value = result.data;
    }
  } finally {
    tableLoading.value = false;
  }
}

const queryLoading = ref<boolean>(false);
const handleQuery = async () => {
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

const handleDelete = async (record: Concat) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await concatApi.deleteConcat(record.concatId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const formRef = ref();
const showForm = (record?: Concat) => {
  formRef.value.show(record);
}

onMounted(() => {
  queryTableData();
})
</script>

<template>
  <div class="table-card">
    <div class="search-div">
      <div class="search-item">
        <label>联系方式名称</label>
        <a-input v-model="queryForm.name" placeholder="请输入联系方式名称" />
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
        <a-button type="primary" @click="showForm(undefined)">
          <template #icon><icon-plus /></template>
          新增
        </a-button>
      </div>
      <right-operate @refresh="handleRefresh"/>
    </div>
    <a-table :data="tableData" :loading="tableLoading" :pagination="false">
      <template #columns>
        <a-table-column title="联系方式ID" data-index="concatId" />
        <a-table-column title="联系方式名称" data-index="name" />
        <a-table-column title="Logo" data-index="coverUrl">
          <template #cell="{ record }">
            <load-image :src="record.logoUrl" :local="false" thumbnail width="48px" height="48px" radius="8px" />
          </template>
        </a-table-column>
        <a-table-column title="源链接">
          <template #cell="{ record }">
            <a :href="record.targetUrl" target="_blank" class="link-text">{{ record.targetUrl }}</a>
          </template>
        </a-table-column>
        <a-table-column title="排序" data-index="sort" />
        <a-table-column title="创建时间" data-index="createTime" align="center" :width="280" />
        <a-table-column title="主要联系方式" align="center">
          <template #cell="{ record }">
            <a-switch :model-value="record.isMain" />
          </template>
        </a-table-column>
        <a-table-column title="状态" :width="60">
          <template #cell="{ record }">
            <a-switch :checked-value="0" :unchecked-value="1" :model-value="record.status" />
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
    <concat-form ref="formRef" @reload="handleQuery" />
  </div>
</template>

<style scoped lang="scss">

</style>