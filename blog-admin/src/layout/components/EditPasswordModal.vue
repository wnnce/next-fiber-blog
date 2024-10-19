<script setup lang="ts">

import { reactive, ref } from 'vue'
import type { ResetPasswordForm } from '@/api/system/user/types'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import { userApi } from '@/api/system/user'
import { useLocalUserStore } from '@/stores/user'
import { useLocalStorage } from '@/hooks/local-storage'
import { LOCAl_USER_KEY, TOKEN_KEY } from '@/assets/script/constant'
import { useArcoMessage } from '@/hooks/message'
import { useRouter } from 'vue-router'

interface EditPasswordForm extends ResetPasswordForm {
  repeatPassword: string;
}

const router = useRouter();

const modalShow = ref<boolean>(false);
const show = () => {
  modalShow.value = true;
}
const onClose = () => {
  Object.assign(formData, defaultFormData);
}

const defaultFormData: EditPasswordForm = {
  oldPassword: '',
  newPassword: '',
  repeatPassword: '',
}
const formData = reactive<EditPasswordForm>({ ...defaultFormData });

const passwordFieldRule = (message: string): FieldRule<any>[] => {
  return [
    { required: true, message: message, type: 'string' },
    {
      validator: (value, callback) => {
        const reg = /^(?=.*[a-zA-Z])(?=.*\d).{8,}$/;
        if (!reg.test(value)) {
          callback('密码最少8位且包含字母和数字');
        } else {
          callback();
        }
      }
    }
  ]
}
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  oldPassword: passwordFieldRule('原密码不能为空'),
  newPassword: passwordFieldRule('新密码不能为空'),
  repeatPassword: [
    { required: true, type: 'string', message: '确认密码不能为空' },
    {
      validator: (value, callback) => {
        if (value !== formData.newPassword) {
          callback('确认密码和新密码不一致');
        } else {
          callback();
        }
      }
    }
  ]
}

const submitButtonLoading = ref<boolean>(false);

const formSubmit = async () => {
  submitButtonLoading.value = true;
  const requestForm: ResetPasswordForm = {
    oldPassword: btoa(formData.oldPassword),
    newPassword: btoa(formData.newPassword)
  }
  try {
    const result = await userApi.resetPassword(requestForm);
    if (result.code === 200) {
      useLocalUserStore().clear()
      useLocalStorage().remove(TOKEN_KEY, LOCAl_USER_KEY);
      useArcoMessage().successMessage('修改成功，请重新登录!');
      modalShow.value = false;
      setTimeout(() => {
        router.push({ path: '/login' })
      }, 300)
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
  <a-modal title="修改密码" v-model:visible="modalShow" unmountOnClose :footer="false"
           width="400px" @close="onClose"
  >
    <a-form layout="vertical" :model="formData" @submit-success="formSubmit" auto-label-width :rules="formRules">
      <a-form-item label="原密码" field="oldPassword">
        <a-input-password v-model="formData.oldPassword" placeholder="请输入原密码" autocomplete />
      </a-form-item>
      <a-form-item label="新密码" field="newPassword">
        <a-input-password v-model="formData.newPassword" placeholder="请输入新密码" autocomplete />
      </a-form-item>
      <a-form-item label="确认密码" field="repeatPassword">
        <a-input-password v-model="formData.repeatPassword" placeholder="请输入确认密码" autocomplete />
      </a-form-item>
      <div class="flex justify-between" style="width: 100%; column-gap: 24px">
        <a-button html-type="submit" type="primary" long :loading="submitButtonLoading">
          <template #icon><icon-save /></template>
          确认
        </a-button>
        <a-button long @click="modalShow = false" :disabled="submitButtonLoading">
          <template #icon><icon-close /></template>
          取消
        </a-button>
      </div>
    </a-form>
  </a-modal>
</template>

<style scoped lang="scss">

</style>