<script setup lang="ts">

import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import type { Concat, ConcatForm } from '@/api/blog/concat/types'
import { concatApi } from '@/api/blog/concat'
import DictSelect from '@/components/DictSelect.vue'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Concat) => {
  if (record) {
    const { concatId, name, iconSvg, targetUrl, isMain, sort, status } = record;
    Object.assign(formData, { concatId, name, iconSvg, targetUrl, isMain, sort, status })
  }
  modalShow.value = true;
}
const onClose = () => {
  Object.assign(formData, defaultFormData);
}

const defaultFormData: ConcatForm = {
  concatId: undefined,
  name: '',
  iconSvg: '',
  targetUrl: '',
  isMain: false,
  sort: 1,
  status: 0,
}
const formData = reactive<ConcatForm>({ ...defaultFormData })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  name: { required: true, message: '联系方式名称不能为空' },
  iconSvg: { required: true, message: 'Icon不能为空' },
  targetUrl: { required: true, message: '源链接不能为空', type: 'url' },
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.concatId) {
      result = await concatApi.updateConcat(formData);
    } else {
      result = await concatApi.saveConcat(formData);
    }
    if (result.code === 200) {
      successMessage(formData.concatId ? '修改成功' : '保存成功');
      emits('reload');
      modalShow.value = false;
    }
  } finally {
    submitButtonLoading.value = false;
  }
}

defineExpose({
  show
})

</script>

<template>
  <a-modal :title="formData.concatId ? '修改联系方式' : '添加联系方式'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" auto-label-width @submit-success="formSubmit" :rules="formRules">
      <a-form-item label="名称" field="name">
        <a-input v-model="formData.name" placeholder="请输入联系方式名称" />
      </a-form-item>
      <a-form-item label="源链接" field="targetUrl">
        <a-input v-model="formData.targetUrl" placeholder="请输入源链接" />
      </a-form-item>
      <a-form-item label="Icon" field="Icon">
        <a-textarea v-model="formData.iconSvg" placeholder="请输入Icon SVG图片内容" />
      </a-form-item>
      <a-form-item label="是否主要" field="isMain">
        <a-switch v-model="formData.isMain" />
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