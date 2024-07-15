<script setup lang="ts">

import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { FileItem } from '@arco-design/web-vue'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import ImageUpload from '@/components/ImageUpload.vue'
import type { Link, LinkForm } from '@/api/blog/link/types'
import { linkApi } from '@/api/blog/link'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Link) => {
  if (record) {
    const { linkId, name, summary, coverUrl, targetUrl, sort, status } = record;
    Object.assign(formData, { linkId, name, summary, coverUrl: coverUrl || '', targetUrl, sort, status })
    formatAvatarToFileList();
  }
  modalShow.value = true;
}
const onClose = () => {
  fileList.value = [];
  Object.assign(formData, defaultFormData);
}

const defaultFormData: LinkForm = {
  linkId: undefined,
  name: '',
  summary: undefined,
  coverUrl: '',
  targetUrl: '',
  sort: 1,
  status: 0,
}
const formData = reactive<LinkForm>({ ...defaultFormData })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  name: { required: true, message: '友情链接名称不能为空' },
  targetUrl: { required: true, message: '源链接不能为空', type: 'url' },
  coverUrl: { required: true, message: '标签封面不能为空' },
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.linkId) {
      result = await linkApi.updateLink(formData);
    } else {
      result = await linkApi.saveLink(formData);
    }
    if (result.code === 200) {
      successMessage(formData.linkId ? '修改成功' : '保存成功');
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
  <a-modal :title="formData.linkId ? '修改友情链接' : '添加友情链接'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" auto-label-width @submit="formSubmit" :rules="formRules">
      <a-form-item label="友情链接名称" field="name">
        <a-input v-model="formData.name" placeholder="请输入友情链接名称" />
      </a-form-item>
      <a-form-item label="简介" field="summary">
        <a-textarea v-model="formData.summary" placeholder="请输入友情链接简介" />
      </a-form-item>
      <a-form-item label="源链接" field="targetUrl">
        <a-input v-model="formData.targetUrl" placeholder="请输入源链接" />
      </a-form-item>
      <a-form-item label="封面" field="coverUrl">
        <image-upload v-model:file-list="fileList" v-model:file-url="formData.coverUrl" width="100px" height="100px" />
      </a-form-item>
      <a-form-item label="显示顺序" field="sort">
        <a-input-number v-model="formData.sort" placeholder="请输入显示顺序" />
      </a-form-item>
      <a-form-item label="状态" field="status" required>
        <a-switch :checked-value="0" :unchecked-value="1" v-model="formData.status" />
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