<script setup lang="ts">

import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { TreeNodeData } from '@arco-design/web-vue'
import { menuApi } from '@/api/system/menu'
import type { Menu } from '@/api/system/menu/types'
import type { User, UserForm } from '@/api/system/user/types'
import { userApi } from '@/api/system/user'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import type { OptionItem } from '@/assets/script/types'
import ImageUpload from '@/components/ImageUpload.vue'

interface Props {
  roleSelectOption: OptionItem[]
}

defineProps<Props>();

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: User) => {
  if (record) {
    const { userId, username, nickname, email, phone, avatar, remark, roles, sort, status } = record;
    Object.assign(formData, { userId, username, nickname, email, phone, avatar, remark, roles, sort, status })
  }
  queryTreeSelectData();
  modalShow.value = true;
}
const onClose = () => {
  treeSelectData.value = [];
  Object.assign(formData, defaultFormData);
}

const defaultFormData: UserForm = {
  userId: undefined,
  username: '',
  nickname: undefined,
  password: '',
  email: undefined,
  phone: undefined,
  avatar: '',
  sort: 1,
  status: 0,
  roles: [],
  remark: undefined,
}
const formData = reactive<UserForm>({ ...defaultFormData })

const formRules: Record<string, FieldRule<any> | FieldRule<any>[] | undefined> = {
  username: { required: true, message: '用户名不能为空' },
  password: [
    { required: !formData.userId, message: '密码不能为空' },
    {
      validator(value, callback) {
        if (!formData.userId) {
          if (!value || value.trim().length === 0) {
            callback('密码不能为空');
            return;
          }
          const regExp = new RegExp('^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d|.*[!@#$%^&*?])[A-Za-z\\d!@#$%^&*?]{8,}$')
          if (!regExp.test(value)) {
            callback('密码最低8为，且为数字，字母组合')
            return;
          }
        }
        callback();
      },
    },
  ],
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
  roles: { required: true, message: '所属角色不能为空' },
  email: { type: 'email', required: false },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  console.log(formData)
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.userId) {
      result = await userApi.updateSysUser(formData);
    } else {
      result = await userApi.saveSysUser(formData);
    }
    if (result.code === 200) {
      successMessage(formData.userId ? '修改成功' : '保存成功');
      emits('reload');
      modalShow.value = false;
    }
  } finally {
    submitButtonLoading.value = false;
  }
}

const treeSelectData = ref<TreeNodeData[]>([]);
const queryTreeSelectData = async () => {
  const result = await menuApi.manageListTree();
  const { code, data } = result;
  if (code === 200) {
    treeSelectData.value = parseMenuToSelectOption(data);
  }
}

const parseMenuToSelectOption = (menus: Menu[]): TreeNodeData[] => {
  if (!menus || menus.length === 0) {
    return [];
  }
  return menus.map(item => {
    return {
      key: item.menuId,
      title: item.menuName,
      children: item.children && item.children.length > 0 ? parseMenuToSelectOption(item.children) : undefined,
    }
  })
}

defineExpose({
  show
})

</script>

<template>
  <a-modal :title="formData.userId ? '修改用户' : '添加用户'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" auto-label-width @submit="formSubmit" :rules="formRules">
      <a-form-item label="用户名" field="username">
        <a-input v-model="formData.username" placeholder="请输入角色名" />
      </a-form-item>
      <a-form-item label="昵称" field="nickname">
        <a-input v-model="formData.nickname" placeholder="请输入昵称" />
      </a-form-item>
      <a-form-item label="密码" field="password" v-if="!formData.userId">
        <a-input-password v-model="formData.password" placeholder="请输入密码" />
      </a-form-item>
      <a-form-item label="邮箱" field="email">
        <a-input v-model="formData.email" placeholder="请输入邮箱" />
      </a-form-item>
      <a-form-item label="手机号" field="phone">
        <a-input v-model="formData.phone" placeholder="请输入手机号" />
      </a-form-item>
      <a-form-item label="头像" field="avatar">
        <image-upload v-model:file-list="formData.avatar"/>
      </a-form-item>
      <a-form-item label="显示顺序" field="sort">
        <a-input-number v-model="formData.sort" placeholder="请输入显示顺序" />
      </a-form-item>
      <a-form-item label="所属角色" field="roles">
        <a-select multiple v-model="formData.roles" placeholder="请选择所属角色">
          <a-option v-for="item in roleSelectOption" :key="item.value" :value="item.value" :label="item.label" />
        </a-select>
      </a-form-item>
      <a-form-item label="备注" field="remark">
        <a-textarea v-model="formData.remark" placeholder="备注" />
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