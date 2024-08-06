<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import type { FileItem, TreeNodeData } from '@arco-design/web-vue'
import type { Result } from '@/api/request'
import { useArcoMessage } from '@/hooks/message'
import type { Category, CategoryForm } from '@/api/blog/category/types'
import { categoryApi } from '@/api/blog/category'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import ImageUpload from '@/components/ImageUpload.vue'
import DictSelect from '@/components/DictSelect.vue'

const { successMessage } = useArcoMessage();

interface Props {
  treeCategory: Category[]
}

const emits = defineEmits<{
  (e: 'reload'): void
}>()
const props = defineProps<Props>();

const modalShow = ref<boolean>(false);
const show = (record?: Category, parentId?: number) => {
  if (!record && (parentId || parentId === 0)) {
    formData.parentId = parentId;
  } else if (record) {
    const { categoryId, categoryName, description, coverUrl, isHot, isTop, parentId, sort, status } = record;
    Object.assign(formData, { categoryId, categoryName, description, parentId, coverUrl, sort, status, isHot, isTop });
  }
  formatAvatarToFileList();
  modalShow.value = true;
}
const onClose = () => {
  fileList.value = [];
  Object.assign(formData, defaultMenuForm);
}

const defaultMenuForm: CategoryForm = {
  categoryId: undefined,
  categoryName: '',
  parentId: undefined,
  description: undefined,
  coverUrl: '',
  isHot: false,
  isTop: false,
  sort: 1,
  status: 0,
}
const formData = reactive<CategoryForm>({ ...defaultMenuForm })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  categoryName: { required: true, message: '分类名称不能为空' },
  parentId: { required: true, message: '上级分类不能为空' },
  coverUrl: { required: true, message: '分类封面不能为空' },
  isHot: { required: true, message: '是否热门不能为空' },
  isTop: { required: true, message: '是否置顶不能为空' },
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.categoryId) {
      result = await categoryApi.updateCategory(formData);
    } else {
      result = await categoryApi.saveCategory(formData);
    }
    if (result.code === 200) {
      successMessage(formData.categoryId ? '修改成功' : '保存成功');
      emits('reload');
      modalShow.value = false;
    }
  } finally {
    submitButtonLoading.value = false;
  }
}

const treeSelectData = computed((): TreeNodeData[] => {
  const treeOptions = parseMenuToSelectOption(props.treeCategory);
  return [{
    key: 0,
    title: '主分类',
    disabled: false,
    children: treeOptions
  }]
})
const parseMenuToSelectOption = (categorys: Category[]): TreeNodeData[] => {
  if (!categorys || categorys.length === 0) {
    return [];
  }
  return categorys.map(item => {
    return {
      key: item.categoryId,
      title: item.categoryName,
      children: item.children && item.children.length > 0 ? parseMenuToSelectOption(item.children) : undefined,
    }
  })
}

const fileList = ref<FileItem[]>([]);
const formatAvatarToFileList = () => {
  if (!formData.coverUrl || formData.coverUrl.trim().length === 0) {
    return;
  }
  fileList.value = [{
    uid: new Date().getTime().toString(),
    status: 'done',
    percent: 1,
    url: formData.coverUrl
  }]
}

defineExpose({
  show
})

</script>

<template>
  <a-modal :title="formData.categoryId ? '修改分类' : '新增分类'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" @submit-success="formSubmit" auto-label-width :rules="formRules">
      <a-form-item label="上级分类" field="parentId">
        <a-tree-select v-model="formData.parentId" :data="treeSelectData" :tree-props="{ defaultExpandAll: false }"
                       placeholder="请选择上级菜单"
        />
      </a-form-item>
      <a-form-item label="分类名称" field="categoryName">
        <a-input v-model="formData.categoryName" placeholder="请输入分类名称" />
      </a-form-item>
      <a-form-item label="分类描述" field="description">
        <a-textarea v-model="formData.description" placeholder="请输入分类描述" />
      </a-form-item>
      <a-form-item label="封面" field="coverUrl">
        <image-upload v-model:file-list="fileList" v-model:file-url="formData.coverUrl" width="200px" height="100px" />
      </a-form-item>
      <a-form-item label="热门分类" field="isHot">
        <a-switch v-model="formData.isHot" />
      </a-form-item>
      <a-form-item label="分类置顶" field="isTop">
        <a-switch v-model="formData.isTop" />
      </a-form-item>
      <a-form-item label="显示顺序" field="sort">
        <a-input-number v-model="formData.sort" placeholder="请输入显示顺序" />
      </a-form-item>
      <a-form-item label="状态" field="status" required>
        <dict-select dict-key="dict_status" v-model="formData.status" type="number" />
      </a-form-item>
      <div class="flex justify-between" style="width: 100%; column-gap: 24px">
        <a-button html-type="submit" type="primary" size="large" long :loading="submitButtonLoading">
          <template #icon><icon-save /></template>
          提交
        </a-button>
        <a-button size="large" long @click="modalShow = false" :disabled="submitButtonLoading">
          <template #icon><icon-close /></template>
          取消
        </a-button>
      </div>
    </a-form>
  </a-modal>
</template>

<style scoped lang="scss">
.icon-container {
  width: 300px;
  height: 160px;
  gap: var(--space-sm);
  flex-wrap: wrap;
  overflow-y: auto;
  font-size: 20px;
  > div {
    text-align: center;
    width: calc((100% - (var(--space-sm) * 4)) / 5);
  }
}
</style>