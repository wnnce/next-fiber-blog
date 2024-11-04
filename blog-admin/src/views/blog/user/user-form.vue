<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import DictSelect from '@/components/DictSelect.vue'
import type { User, UserUpdateForm } from '@/api/blog/user/types'
import { userApi } from '@/api/blog/user'

const { successMessage } = useArcoMessage();

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const modalShow = ref<boolean>(false);
const show = (record: User) => {
  const { userId, nickname, summary, email, link, labels, status } = record;
  Object.assign(formData, { userId, nickname, summary, email, link, labels, status: status || 0 });
  modalShow.value = true;
}
const onClose = () => {
  Object.assign(formData, defaultForm);
}

const defaultForm: UserUpdateForm = {
  userId: 0,
  nickname: undefined,
  summary: undefined,
  email: undefined,
  link: undefined,
  labels: [],
  status: 0,
}
const formData = reactive<UserUpdateForm>({ ...defaultForm })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  nickname: { required: true, message: '用户昵称不能为空' },
  email: { required: true, message: '用户邮箱不能为空', type: 'email' },
  link: { required: true, message: '用户站点链接不能为空', type: 'url' },
  status: { required: true, message: '用户状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    const result = await userApi.updateUser(formData);
    if (result.code === 200) {
      successMessage('修改成功');
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
  <a-modal title="修改用户信息" v-model:visible="modalShow" @close="onClose" :footer="false" unmount-on-close>
    <a-form layout="vertical" :model="formData" @submit-success="formSubmit" auto-label-width :rules="formRules">
      <a-form-item label="昵称" field="nickname">
        <a-input v-model="formData.nickname" placeholder="请输入用户昵称" />
      </a-form-item>
      <a-form-item label="简介" field="summary">
        <a-textarea v-model="formData.summary" placeholder="请输入用户简介" />
      </a-form-item>
      <a-form-item label="邮箱" field="email">
        <a-input v-model="formData.email" placeholder="请输入用户邮箱" />
      </a-form-item>
      <a-form-item label="链接" field="link">
        <a-input v-model="formData.link" placeholder="请输入用户站点链接" />
      </a-form-item>
      <a-form-item label="标签" field="labels">
        <a-input-tag v-model="formData.labels" allow-clear placeholder="用户个性化标签" />
      </a-form-item>
      <a-form-item label="状态" field="status">
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

</style>