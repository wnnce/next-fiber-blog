<script setup lang="ts">

import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { FileItem } from '@arco-design/web-vue'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import ImageUpload from '@/components/ImageUpload.vue'
import type { Tag, TagForm } from '@/api/blog/tags/types'
import { tagApi } from '@/api/blog/tags'
import DictSelect from '@/components/DictSelect.vue'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Tag) => {
  if (record) {
    const { tagId, tagName, color, coverUrl, sort, status } = record;
    Object.assign(formData, { tagId, tagName, color, coverUrl: coverUrl || '', sort, status })
    formatAvatarToFileList();
  }
  modalShow.value = true;
}
const onClose = () => {
  fileList.value = [];
  Object.assign(formData, defaultFormData);
}

const defaultFormData: TagForm = {
  tagId: undefined,
  tagName: '',
  color: '#fff',
  coverUrl: '',
  sort: 1,
  status: 0,
}
const formData = reactive<TagForm>({ ...defaultFormData })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  tagName: { required: true, message: '标签名称不能为空' },
  color: { required: true, message: '标签颜色不能为空' },
  coverUrl: { required: true, message: '标签封面不能为空' },
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.tagId) {
      result = await tagApi.updateTag(formData);
    } else {
      result = await tagApi.saveTag(formData);
    }
    if (result.code === 200) {
      successMessage(formData.tagId ? '修改成功' : '保存成功');
      emits('reload');
      modalShow.value = false;
    }
  } finally {
    submitButtonLoading.value = false;
  }
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
  <a-modal :title="formData.tagId ? '修改标签' : '添加标签'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" auto-label-width @submit-success="formSubmit" :rules="formRules">
      <a-form-item label="标签名称" field="tagName">
        <a-input v-model="formData.tagName" placeholder="请输入标签名称" />
      </a-form-item>
      <a-form-item label="标签颜色" field="color">
        <a-color-picker v-model="formData.color" />
      </a-form-item>
      <a-form-item label="封面" field="coverUrl">
        <image-upload v-model:file-list="fileList" v-model:file-url="formData.coverUrl" width="100px" height="100px" />
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
.tree-select {
  width: 100%;
  max-height: 160px;
  overflow-y: auto;
  border: 1px solid var(--border-color);
}
</style>