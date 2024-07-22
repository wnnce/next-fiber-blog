<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import { constant } from '@/assets/script/constant'
import type { Dict, DictQueryForm } from '@/api/system/dict/types'
import { dictApi } from '@/api/system/dict'
import DictForm from '@/views/system/dict/dict-form.vue'
import DictValueDrawer from '@/views/system/dict/dict-value-drawer.vue'
import type { Notice, NoticeQueryForm } from '@/api/system/notice/types'
import { noticeApi } from '@/api/system/notice'
import DictLabel from '@/components/DictLabel.vue'
import NoticeForm from '@/views/system/notice/notice-form.vue'

const { successMessage, loading, errorMessage } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<Notice[]>([]);
const defaultQueryForm: NoticeQueryForm = {
  page: 1,
  size: 10,
  title: undefined,
  level: undefined,
  noticeType: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<NoticeQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await noticeApi.pageNotice(queryForm);
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

const handleDelete = async (record: Notice) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await noticeApi.deleteNotice(record.noticeId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const formRef = ref();
const showForm = (record?: Notice) => {
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
        <label>字典名称</label>
        <a-input v-model="queryForm.title" placeholder="请输入通知标题" />
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
        <a-table-column title="ID" data-index="noticeId" />
        <a-table-column title="标题" data-index="title" />
        <a-table-column title="级别">
          <template #cell="{ record }">
            <dict-label dict-key="notice_level" :value="record.level" />
          </template>
        </a-table-column>
        <a-table-column title="类型">
          <template #cell="{ record }">
            <dict-label dict-key="notice_type" :value="record.noticeType" />
          </template>
        </a-table-column>
        <a-table-column title="排序" data-index="sort" />
        <a-table-column title="创建时间" data-index="createTime" />
        <a-table-column title="备注" data-index="remark" />
        <a-table-column title="状态">
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
    <notice-form ref="formRef" @reload="queryTableData" />
  </div>
</template>

<style scoped lang="scss">

</style>