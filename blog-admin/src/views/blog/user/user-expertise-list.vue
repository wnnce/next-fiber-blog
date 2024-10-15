<script setup lang="ts">
import { onMounted, reactive, ref, shallowRef } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import type { ExpertiseDetail, ExpertiseQueryForm } from '@/api/blog/user/types'
import { userApi } from '@/api/blog/user'
import { constant } from '@/assets/script/constant'
import DictLabel from '@/components/DictLabel.vue'
import DictSelect from '@/components/DictSelect.vue'

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = shallowRef<ExpertiseDetail[]>([]);
const defaultQueryForm: ExpertiseQueryForm = {
  page: 1,
  size: 10,
  username: undefined,
  source: undefined,
  detailType: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<ExpertiseQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await userApi.pageExpertise(queryForm);
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
  queryRoleSelectData();
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
        <a-input v-model="queryForm.username" placeholder="请输入用户名" />
      </div>
      <div class="search-item">
        <label>经验类型</label>
        <DictSelect dict-key="expertise_detail_type" v-model="queryForm.detailType" type="number" placeholder="请选择类型" width="140px" />
      </div>
      <div class="search-item">
        <label>来源</label>
        <DictSelect dict-key="expertise_source" v-model="queryForm.source" type="number" placeholder="请选择来源" width="140px" />
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

      </div>
      <right-operate @refresh="handleRefresh"/>
    </div>
    <a-table :data="tableData" :loading="tableLoading" :pagination="false">
      <template #columns>
        <a-table-column title="ID" data-index="id" />
        <a-table-column title="用户名">
          <template #cell="{ record }">
            <span class="link-text" @click="showUserInfo(record)">{{ record.username }}</span>
          </template>
        </a-table-column>
        <a-table-column title="昵称" data-index="nickname" />
        <a-table-column title="明细" data-index="detail" />
        <a-table-column title="类型">
          <template #cell="{ record }">
            <DictLabel dict-key="expertise_detail_type" :value="record.detailType" />
          </template>
        </a-table-column>
        <a-table-column title="来源">
          <template #cell="{ record }">
            <DictLabel dict-key="expertise_source" :value="record.source" />
          </template>
        </a-table-column>
        <a-table-column title="创建时间" data-index="createTime" align="center" />
        <a-table-column title="备注" data-index="remark" />
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
  </div>
</template>

<style scoped lang="scss">

</style>