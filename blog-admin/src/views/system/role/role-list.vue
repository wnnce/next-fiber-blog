<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import type { Role, RoleQueryForm } from '@/api/system/role/types'
import { roleApi } from '@/api/system/role'
import RoleForm from '@/views/system/role/role-form.vue'

const { successMessage, loading } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<Role[]>([]);
const defaultQueryForm: RoleQueryForm = {
  page: 1,
  size: 10,
  name: undefined,
  key: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<RoleQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await roleApi.pageSysRole(queryForm);
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

const handleDelete = async (record: Role) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await roleApi.deleteSysRole(record.roleId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const handleUpdateStatus = async (newStatus: number, roleId: number) => {
  const result = await roleApi.updateSelective({ roleId: roleId, status: newStatus })
  if (result.code === 200) {
    successMessage('更新成功')
    return true;
  }
  return false;
}

const formRef = ref();
const showForm = (record?: Role) => {
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
        <label>角色名称</label>
        <a-input v-model="queryForm.name" placeholder="请输入角色名称" />
      </div>
      <div class="search-item">
        <label>角色标识</label>
        <a-input v-model="queryForm.key" placeholder="请输入角色标识" />
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
    <a-table :data="tableData" :loading="tableLoading" >
      <template #columns>
        <a-table-column title="角色ID" data-index="roleId" />
        <a-table-column title="角色名称" data-index="roleName" />
        <a-table-column title="角色标识" data-index="roleKey" />
        <a-table-column title="排序" data-index="sort" />
        <a-table-column title="创建时间" data-index="createTime" align="center" />
        <a-table-column title="备注" data-index="remark" />
        <a-table-column title="状态">
          <template #cell="{ record }">
            <a-switch :checked-value="0" :unchecked-value="1" v-model="record.status"
                      :before-change="(newValue: number) => handleUpdateStatus(newValue, record.roleId)"
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
    <role-form ref="formRef" @reload="queryTableData" />
  </div>
</template>

<style scoped lang="scss">

</style>