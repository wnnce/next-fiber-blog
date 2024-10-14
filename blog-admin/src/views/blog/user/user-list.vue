<script setup lang="ts">
import { onMounted, reactive, ref, shallowRef } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import type { User, UserQueryForm } from '@/api/blog/user/types'
import { userApi } from '@/api/blog/user'
import { constant } from '@/assets/script/constant'
import LoadImage from '@/components/LoadImage.vue'
import DictLabel from '@/components/DictLabel.vue'
import UserInfo from '@/views/blog/user/user-info.vue'
import UserForm from '@/views/blog/user/user-form.vue'

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = shallowRef<User[]>([]);
const defaultQueryForm: UserQueryForm = {
  page: 1,
  size: 10,
  username: undefined,
  nickname: undefined,
  email: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<UserQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await userApi.pageUser(queryForm);
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

const userInfoRef = ref();
const showUserInfo = (record: User) => {
  userInfoRef.value.show(record);
}

const formRef = ref();
const showForm = (record: User) => {
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
        <label>用户名</label>
        <a-input v-model="queryForm.username" placeholder="请输入用户名" />
      </div>
      <div class="search-item">
        <label>昵称</label>
        <a-input v-model="queryForm.nickname" placeholder="请输入昵称" />
      </div>
      <div class="search-item">
        <label>邮箱</label>
        <a-input v-model="queryForm.email" placeholder="请输入邮箱" />
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
        <a-table-column title="头像">
          <template #cell="{ record }">
            <LoadImage :src="record.avatar" :width="48" :height="48" radius="50%" />
          </template>
        </a-table-column>
        <a-table-column title="用户名">
          <template #cell="{ record }">
            <span class="link-text" @click="showUserInfo(record)">{{ record.username }}</span>
          </template>
        </a-table-column>
        <a-table-column title="昵称" data-index="nickname" />
        <a-table-column title="邮箱" data-index="email" />
        <a-table-column title="链接">
          <template #cell="{ record }">
            <a class="link-text" :href="record.link" target="_blank">{{ record.link }}</a>
          </template>
        </a-table-column>
        <a-table-column title="等级" data-index="level" />
        <a-table-column title="创建时间" data-index="createTime" align="center" />
        <a-table-column title="用户类型">
          <template #cell="{ record }">
            <span>{{ record.userType === 1 ? '管理员' : '普通用户' }}</span>
          </template>
        </a-table-column>
        <a-table-column title="状态">
          <template #cell="{ record }">
            <a-tag :color="record.status && record.status === 1 ? 'red' : 'green'">
              <DictLabel dict-key="dict_status" :value="record.status || 0" />
            </a-tag>
          </template>
        </a-table-column>
        <a-table-column title="操作" align="center">
          <template #cell="{ record }">
            <a-button type="text" shape="circle" @click="showForm(record)">
              <template #icon><icon-edit /></template>
            </a-button>
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
    <user-info ref="userInfoRef" />
    <user-form ref="formRef" @reload="queryTableData" />
  </div>
</template>

<style scoped lang="scss">

</style>