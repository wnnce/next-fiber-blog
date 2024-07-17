<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { constant } from '@/assets/script/constant'
import TagsForm from '@/views/blog/tags/tags-form.vue'
import type { LoginRecord, LoginRecordQueryForm } from '@/api/system/record/types'
import { recordApi } from '@/api/system/record'

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<LoginRecord[]>([]);
const defaultQueryForm: LoginRecordQueryForm = {
  page: 1,
  size: 10,
  username: undefined,
  loginType: undefined,
  result: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<LoginRecordQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await recordApi.pageLoginRecord(queryForm);
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
        <label>用户名</label>
        <a-input v-model="queryForm.username" placeholder="通过用户名搜索" />
      </div>
      <div class="search-item">
        <label>登录类型</label>
        <a-select v-model="queryForm.loginType" placeholder="通过登录类型搜索">
          <a-option :value="1" label="博客登录" />
          <a-option :value="2" label="后台登录" />
        </a-select>
      </div>
      <div class="search-item">
        <label>状态</label>
        <a-select v-model="queryForm.result" placeholder="通过状态搜索">
          <a-option :value="0" label="成功" />
          <a-option :value="1" label="失败" />
        </a-select>
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
        <a-table-column title="用户名" data-index="username" />
        <a-table-column title="用户类型" data-index="userType" align="center" />
        <a-table-column title="登录IP" data-index="loginIp" align="center" />
        <a-table-column title="地址" data-index="location" align="center" />
        <a-table-column title="UA" :width="300" ellipsis align="center">
          <template #cell="{ record }">
            <a-tooltip :content="record.loginUa">
              <p>{{ record.loginUa }}</p>
            </a-tooltip>
          </template>
        </a-table-column>
        <a-table-column title="登录时间" data-index="createTime" :width="260" align="center" />
        <a-table-column title="登录类型" data-index="loginType" :width="120">
          <template #cell="{ record }">
            {{ record.loginType === 1 ? '博客登录' : '后台登录' }}
          </template>
        </a-table-column>
        <a-table-column title="状态" align="center" :width="120">
          <template #cell="{ record }">
            <a-tag v-if="record.result == 0" color="green">成功</a-tag>
            <a-tag v-else color="red">失败</a-tag>
          </template>
        </a-table-column>
        <a-table-column title="备注" data-index="remark" align="center" />
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