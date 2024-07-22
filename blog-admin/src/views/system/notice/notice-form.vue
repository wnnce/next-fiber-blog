<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import DictSelect from '@/components/DictSelect.vue'
import type { Notice, NoticeForm } from '@/api/system/notice/types'
import { noticeApi } from '@/api/system/notice'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Notice) => {
  if (record) {
    const {  noticeId, title, message, level, noticeType, sort, status } = record;
    Object.assign(formData, { noticeId, title, message, level, noticeType, sort, status })
  }
  modalShow.value = true;
}
const onClose = () => {
  Object.assign(formData, defaultConfigForm);
}

const defaultConfigForm: NoticeForm = {
  noticeId: undefined,
  title: '',
  message: '',
  level: undefined,
  noticeType: undefined,
  sort: 1,
  status: 0,
}
const formData = reactive<NoticeForm>({ ...defaultConfigForm })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  title: { required: true, message: '通知标题不能为空' },
  message: { required: true, message: '通知内容不能为空' },
  level: { required: true, message: '通知级别不能为空' },
  noticeType: { required: true, message: '通知类型不能为空' },
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.noticeId) {
      result = await noticeApi.updateNotice(formData);
    } else {
      result = await noticeApi.saveNotice(formData);
    }
    if (result.code === 200) {
      successMessage(formData.noticeId ? '修改成功' : '保存成功');
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
  <a-modal :title="formData.noticeId ? '修改通知' : '添加通知'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" layout="vertical" @submit="formSubmit" :rules="formRules">
      <a-form-item label="通知标题" field="title">
        <a-input v-model="formData.title" placeholder="请输入通知标题" />
      </a-form-item>
      <a-form-item label="通知内容" field="message">
        <a-textarea v-model="formData.message" placeholder="请输入通知或公告内容" />
      </a-form-item>
      <a-form-item label="级别" field="level">
        <dict-select dict-key="notice_level" v-model="formData.level" placeholder="请选择通知级别" />
      </a-form-item>
      <a-form-item label="类型" field="noticeType">
        <dict-select dict-key="notice_type" v-model="formData.noticeType" placeholder="请选择通知类型" />
      </a-form-item>
      <a-form-item label="显示顺序" field="sort">
        <a-input-number v-model="formData.sort" placeholder="请输入显示顺序" />
      </a-form-item>
      <a-form-item label="状态" field="status" required>
        <!--        <a-switch :checked-value="0" :unchecked-value="1" v-model="formData.status" />-->
        <dict-select dict-key="dict_status" v-model="formData.status" type="number" />
      </a-form-item>
      <a-form-item>
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
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<style scoped lang="scss">

</style>