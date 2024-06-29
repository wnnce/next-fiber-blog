<script setup lang="ts">

import { onMounted, ref } from 'vue'
import type { Config } from '@/api/system/config/types'
import { configApi } from '@/api/system/config'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import type { Menu } from '@/api/system/menu/types'
import { menuApi } from '@/api/system/menu'
import * as ArcoIcons from '@arco-design/web-vue/es/icon';
import MenuForm from '@/views/system/menu/menu-form.vue'

const { successMessage, loading } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const tableData = ref<Menu[]>([]);
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await menuApi.manageListTree();
    if (result && result.code === 200) {
      tableData.value = result.data;
    }
  } finally {
    tableLoading.value = false;
  }
}

const handleDelete = async (record: Menu) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await menuApi.deleteSysMenu(record.menuId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }

}

const formRef = ref();
const showForm = (record?: Menu, parentId?: number) => {
  formRef.value.show(record, parentId);
}

onMounted(() => {
  queryTableData();
})
</script>

<template>
  <div class="card flex flex-column" style="row-gap: 12px">
    <div class="flex justify-between">
      <div class="flex" style="column-gap: 12px">
        <a-button type="primary" @click="showForm">
          <template #icon><icon-plus /></template>
          新增
        </a-button>
      </div>
      <right-operate @refresh="queryTableData"/>
    </div>
    <a-table :data="tableData" :loading="tableLoading">
      <template #columns>
        <a-table-column title="菜单名称" data-index="menuName" />
        <a-table-column title="图标" align="center">
          <template #cell="{ record }">
            <component :is="ArcoIcons[record.icon as keyof typeof ArcoIcons]" />
          </template>
        </a-table-column>
        <a-table-column title="排序" data-index="sort" align="center"/>
        <a-table-column title="路由地址" data-index="path" align="center"/>
        <a-table-column title="组件路径" data-index="component" align="center"/>
        <a-table-column title="创建时间" data-index="createTime" align="center"/>
        <a-table-column title="操作" align="center">
          <template #cell="{ record }">
            <a-button type="text" shape="circle" @click="showForm(undefined, record.menuId)">
              <template #icon><icon-plus /></template>
            </a-button>
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
    <menu-form ref="formRef" :tree-menu="tableData" @reload="queryTableData"/>
  </div>
</template>

<style scoped lang="scss">

</style>