<script setup lang="ts">
import { onMounted, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import type { Category, CategoryUpdateForm } from '@/api/blog/category/types'
import LoadImage from '@/components/LoadImage.vue'
import CategoryForm from '@/views/blog/category/category-form.vue'
import { categoryApi } from '@/api/blog/category'

const { successMessage, loading } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const tableData = ref<Category[]>([]);
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await categoryApi.manageTree();
    if (result && result.code === 200) {
      tableData.value = result.data;
    }
  } finally {
    tableLoading.value = false;
  }
}

const handleDelete = async (record: Category) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await categoryApi.deleteCategory(record.categoryId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const handleUpdateStatus = async (form: CategoryUpdateForm) => {
  const result = await categoryApi.updateCategoryStatus(form);
  if (result.code === 200) {
    successMessage('更新成功');
    return true;
  }
  return false;
}

const formRef = ref();
const showForm = (record?: Category, parentId?: number) => {
  formRef.value.show(record, parentId);
}

onMounted(() => {
  queryTableData();
})
</script>

<template>
  <div class="table-card">
    <div class="flex justify-between">
      <div class="flex" style="column-gap: 12px">
        <a-button type="primary" @click="showForm(undefined, 0)">
          <template #icon><icon-plus /></template>
          新增
        </a-button>
      </div>
      <right-operate @refresh="queryTableData"/>
    </div>
    <a-table :data="tableData" :loading="tableLoading" row-key="categoryId" :pagination="false">
      <template #columns>
        <a-table-column title="分类名称" data-index="categoryName" />
        <a-table-column title="封面" data-index="coverUrl">
          <template #cell="{ record }">
            <load-image :src="record.coverUrl" :local="false" thumbnail width="48px" height="48px" radius="8px" />
          </template>
        </a-table-column>
        <a-table-column title="查看次数" data-index="viewNum" />
        <a-table-column title="文章数量" data-index="articleNum" />
        <a-table-column title="排序" data-index="sort" align="center"/>
        <a-table-column title="创建时间" data-index="createTime" align="center"/>
        <a-table-column title="热门" :width="60">
          <template #cell="{ record }">
            <a-switch v-model="record.isHot"
                      :before-change="(newValue: boolean) => handleUpdateStatus({ categoryId: record.categoryId, isHot: newValue })"
            />
          </template>
        </a-table-column>
        <a-table-column title="置顶" :width="60">
          <template #cell="{ record }">
            <a-switch v-model="record.isTop"
                      :before-change="(newValue: boolean) => handleUpdateStatus({ categoryId: record.categoryId, isTop: newValue })"
            />
          </template>
        </a-table-column>
        <a-table-column title="状态" :width="60">
          <template #cell="{ record }">
            <a-switch :checked-value="0" :unchecked-value="1"
                      v-model="record.status"
                      :before-change="(newValue: number) => handleUpdateStatus({ categoryId: record.categoryId, status: newValue })"
            />
          </template>
        </a-table-column>
        <a-table-column title="操作" align="center">
          <template #cell="{ record }">
            <a-button type="text" shape="circle" @click="showForm(undefined, record.categoryId)">
              <template #icon><icon-plus /></template>
            </a-button>
            <a-button type="text" shape="circle" @click="showForm(record)">
              <template #icon><icon-edit /></template>
            </a-button>
            <a-popconfirm content="是否确认删除数据？" type="error" position="lt"
                          :ok-button-props="{ status: 'danger' }"
                          @ok="handleDelete(record)"
            >
              <a-button type="text" shape="circle" status="danger" :disabled="record.articleNum > 0">
                <template #icon><icon-delete /></template>
              </a-button>
            </a-popconfirm>
          </template>
        </a-table-column>
      </template>
    </a-table>
    <category-form :tree-category="tableData" ref="formRef" @reload="queryTableData" />
  </div>
</template>

<style scoped lang="scss">

</style>