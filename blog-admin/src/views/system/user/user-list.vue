<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import { roleApi } from '@/api/system/role'
import type { User, UserQueryForm } from '@/api/system/user/types'
import { userApi } from '@/api/system/user'
import type { OptionItem } from '@/assets/script/types'
import { constant } from '@/assets/script/constant'
import UserForm from '@/views/system/user/user-form.vue'

const { successMessage, loading } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<User[]>([]);
const defaultQueryForm: UserQueryForm = {
  page: 1,
  size: 10,
  username: undefined,
  nickname: undefined,
  email: undefined,
  phone: undefined,
  roleId: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<UserQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await userApi.pageSysUser(queryForm);
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

const handleDelete = async (record: User) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await userApi.deleteSysUser(record.userId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const handleUpdateStatus = async (newStatus: number, userId: number) => {
  const result = await userApi.updateSelective({ userId: userId, status: Number(newStatus) })
  if (result.code === 200) {
    successMessage('更新成功')
    return true;
  }
  return false;
}

const roleSelectOption = ref<OptionItem[]>([]);
const queryRoleSelectData = async () => {
  const result = await roleApi.listAllSysROle();
  const { code, data } = result;
  if (code === 200 && data) {
    roleSelectOption.value = data.map(item => {
      return {
        label: item.roleName,
        value: item.roleId
      }
    })
  }
}

const formRef = ref();
const showForm = (record?: User) => {
  formRef.value.show(record);
}

onMounted(() => {
  queryTableData();
  queryRoleSelectData();
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
        <label>手机号</label>
        <a-input v-model="queryForm.phone" placeholder="请输入手机号" />
      </div>
      <div class="search-item">
        <label>角色</label>
        <a-select v-model="queryForm.roleId" placeholder="请选择角色"
                  allow-clear
                  allow-search
                  style="width: 200px"
        >
          <a-option v-for="item in roleSelectOption" :key="item.value" :value="item.value" :label="item.label.toString()" />
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
        <a-button type="primary" @click="showForm(undefined)">
          <template #icon><icon-plus /></template>
          新增
        </a-button>
      </div>
      <right-operate @refresh="handleRefresh"/>
    </div>
    <a-table :data="tableData" :loading="tableLoading" :pagination="false">
      <template #columns>
        <a-table-column title="用户名">
          <template #cell="{ record }">
            <span class="link-text">{{ record.username }}</span>
          </template>
        </a-table-column>
        <a-table-column title="昵称" data-index="nickname" />
        <a-table-column title="邮箱" data-index="email" />
        <a-table-column title="手机号" data-index="phone" />
        <a-table-column title="排序" data-index="sort" />
        <a-table-column title="创建时间" data-index="createTime" align="center" />
        <a-table-column title="最后登录IP" data-index="lastLoginIp" />
        <a-table-column title="最后登录时间" data-index="lastLoginTime" />
        <a-table-column title="状态">
          <template #cell="{ record }">
            <a-switch :checked-value="0" :unchecked-value="1" v-model="record.status"
                      :before-change="(newValue: number) => handleUpdateStatus(newValue, record.userId)"
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
    <user-form :role-select-option="roleSelectOption" ref="formRef" @reload="handleQuery" />
  </div>
</template>

<style scoped lang="scss">

</style>