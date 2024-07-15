<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import { constant } from '@/assets/script/constant'
import { tagApi } from '@/api/blog/tags'
import type { Tag } from '@/api/blog/tags/types'
import LoadImage from '@/components/LoadImage.vue'
import TagsForm from '@/views/blog/tags/tags-form.vue'
import type { Link, LinkQueryForm } from '@/api/blog/link/types'
import { linkApi } from '@/api/blog/link'
import LinksForm from '@/views/blog/links/links-form.vue'

const { successMessage, loading } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<Link[]>([]);
const defaultQueryForm: LinkQueryForm = {
  page: 1,
  size: 10,
  name: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<LinkQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await linkApi.pageLink(queryForm);
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

const handleDelete = async (record: Link) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await linkApi.deleteLink(record.linkId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const formRef = ref();
const showForm = (record?: Link) => {
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
        <label>友情链接名称</label>
        <a-input v-model="queryForm.name" placeholder="请输入友情链接名称" />
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
        <a-table-column title="ID" data-index="linkId" />
        <a-table-column title="友情链接名称" data-index="name" />
        <a-table-column title="封面" data-index="coverUrl">
          <template #cell="{ record }">
            <load-image :src="record.coverUrl" thumbnail width="48px" height="48px" radius="8px" />
          </template>
        </a-table-column>
        <a-table-column title="源链接">
          <template #cell="{ record }">
            <a :href="record.targetUrl" target="_blank" class="link-text">{{ record.targetUrl }}</a>
          </template>
        </a-table-column>
        <a-table-column title="点击次数" data-index="clickNum" align="center" />
        <a-table-column title="排序" data-index="sort" />
        <a-table-column title="创建时间" data-index="createTime" align="center" :width="280" />
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
    <links-form ref="formRef" @reload="handleQuery" />
  </div>
</template>

<style scoped lang="scss">

</style>