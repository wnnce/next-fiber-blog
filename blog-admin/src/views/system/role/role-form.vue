<script setup lang="ts">

import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { Role, RoleForm } from '@/api/system/role/types'
import { roleApi } from '@/api/system/role'
import type { TreeNodeData } from '@arco-design/web-vue'
import { menuApi } from '@/api/system/menu'
import type { Menu } from '@/api/system/menu/types'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Role) => {
  if (record) {
    const { roleId, roleName, roleKey, sort, status, remark, menus } = record;
    Object.assign(formData, { roleId, roleName, roleKey, sort, status, remark, menus })
  }
  queryTreeSelectData();
  modalShow.value = true;
}
const onClose = () => {
  treeSelectData.value = [];
  Object.assign(formData, defaultFormData);
}

const defaultFormData: RoleForm = {
  roleId: undefined,
  roleName: '',
  roleKey: '',
  sort: 1,
  status: 0,
  remark: '',
  menus: []
}
const formData = reactive<RoleForm>({ ...defaultFormData })

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  console.log(formData)
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.roleId) {
      result = await roleApi.updateSysRole(formData);
    } else {
      result = await roleApi.saveSysRole(formData);
    }
    if (result.code === 200) {
      successMessage(formData.roleId ? '修改成功' : '保存成功');
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
  <a-modal :title="formData.roleId ? '修改角色' : '添加角色'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" auto-label-width @submit-success="formSubmit">
      <a-form-item label="角色名称" field="roleName" :rules="[ {required: true, message: '角色名称不能为空'} ]">
        <a-input v-model="formData.roleName" placeholder="请输入角色名称" />
      </a-form-item>
      <a-form-item label="角色标识" field="roleKey" :rules="[ {required: true, message: '角色标识不能为空'} ]">
        <a-input v-model="formData.roleKey" placeholder="请输入角色标识" />
      </a-form-item>
      <a-form-item label="显示顺序" field="sort" :rules="[ {required: true, message: '显示顺序不能为空'} ]">
        <a-input-number v-model="formData.sort" placeholder="请输入显示顺序" />
      </a-form-item>
      <a-form-item label="授权菜单" field="menus" :rules="[ {required: true, message: '授权菜单不能为空'} ]">
        <a-tree v-model:checked-keys="formData.menus" :data="treeSelectData"
                checkable
                default-expand-all
                default-expand-checked
                class="tree-select"
        />
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