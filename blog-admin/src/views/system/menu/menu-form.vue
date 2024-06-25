<script setup lang="ts">

import { computed, reactive, ref } from 'vue'
import type { Menu, MenuForm } from '@/api/system/menu/types'
import type { TreeNodeData } from '@arco-design/web-vue'
import * as ArcoIcons from '@arco-design/web-vue/es/icon';
import type { Result } from '@/api/request'
import { menuApi } from '@/api/system/menu'
import { useArcoMessage } from '@/hooks/message'

const { successMessage } = useArcoMessage();

interface Props {
  treeMenu: Menu[]
}

const emits = defineEmits<{
  (e: 'reload'): void
}>()
const props = defineProps<Props>();

const modalShow = ref<boolean>(false);
const show = (record?: Menu, parentId?: number) => {
  if (!record && parentId) {
    formData.parentId = parentId;
  } else if (record) {
    Object.assign(formData, record);
  }
  modalShow.value = true;
}
const onClose = () => {
  Object.assign(formData, defaultMenuForm);
}

const defaultMenuForm: MenuForm = {
  menuId: undefined,
  menuName: '',
  menuType: undefined,
  parentId: undefined,
  path: undefined,
  component: undefined,
  icon: undefined,
  isFrame: false,
  isCache: true,
  isVisible: true,
  isDisable: false,
  sort: 1
}
const formData = reactive<MenuForm>({ ...defaultMenuForm })

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.menuId) {
      result = await menuApi.updateSysMenu(formData);
    } else {
      result = await menuApi.saveSysMenu(formData);
    }
    if (result.code === 200) {
      successMessage(formData.menuId ? '修改成功' : '保存成功');
      emits('reload');
      modalShow.value = false;
    }
  } finally {
    submitButtonLoading.value = false;
  }
}

const treeSelectData = computed((): TreeNodeData[] => {
  const treeOptions = parseMenuToSelectOption(props.treeMenu);
  return [{
    key: 0,
    title: '主类目',
    disabled: false,
    children: treeOptions
  }]
})
const parseMenuToSelectOption = (menus: Menu[]): TreeNodeData[] => {
  if (!menus || menus.length === 0) {
    return [];
  }
  return menus.map(item => {
    return {
      key: item.menuId,
      title: item.menuName,
      disabled: item.menuType === 2,
      children: item.children && item.children.length > 0 ? parseMenuToSelectOption(item.children) : undefined,
    }
  })
}

defineExpose({
  show
})

</script>

<template>
  <a-modal :title="formData.menuName ? '修改参数配置' : '添加参数配置'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" @submit="formSubmit" auto-label-width>
      <a-form-item label="上级菜单" field="parentId" :rules="[ {required: true, message: '上级菜单不能为空'} ]">
        <a-tree-select v-model="formData.parentId" :data="treeSelectData" placeholder="请选择上级菜单"/>
      </a-form-item>
      <a-form-item label="菜单名称" field="menuName" :rules="[ {required: true, message: '菜单名称不能为空'} ]">
        <a-input v-model="formData.menuName" placeholder="菜单名称不能为空" />
      </a-form-item>
      <a-form-item label="菜单图标" field="icon" :rules="[ {required: true, message: '菜单图标不能为空'} ]">
        <a-popover trigger="click" show-arrow position="bl">
          <a-input :model-value="formData.icon" placeholder="请选择菜单图标" readonly>
            <template #prefix v-if="formData.icon">
              <component :is="ArcoIcons[formData.icon as keyof typeof ArcoIcons]" />
            </template>
          </a-input>
          <template #content>
            <div class="icon-container flex">
              <div v-for="(item, index) in ArcoIcons" :key="index" class="pointer" @click.stop="formData.icon = item.name">
                <component :is="ArcoIcons[item.name]"/>
              </div>
            </div>
          </template>
        </a-popover>
      </a-form-item>
      <div class="flex justify-between">
        <a-form-item label="菜单类型" field="menuType" :rules="[ {required: true, message: '菜单类型不能为空'} ]">
          <a-radio-group v-model="formData.menuType">
            <a-radio :value="1">目录</a-radio>
            <a-radio :value="2">菜单</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="显示排序" field="sort" :rules="[ {required: true, message: '序号不能为空'} ]">
          <a-input-number v-model="formData.sort" placeholder="请输入序号" :min="0" />
        </a-form-item>
      </div>
      <a-form-item label="路由地址" field="path" :rules="[ {required: true, message: '路由地址不能为空'} ]">
        <a-input v-model="formData.path" placeholder="请输入路由地址" />
      </a-form-item>
      <template v-if="formData.menuType === 2">
        <div class="flex">
          <a-form-item label="外链" field="isFrame" required>
            <a-radio-group v-model="formData.isFrame">
              <a-radio :value="true">是</a-radio>
              <a-radio :value="false">否</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="是否缓存" field="isCache" required>
            <a-radio-group v-model="formData.isCache">
              <a-radio :value="true">是</a-radio>
              <a-radio :value="false">否</a-radio>
            </a-radio-group>
          </a-form-item>
        </div>
        <transition name="switch" mode="out-in">
          <a-form-item label="外链地址" field="frameUrl" v-if="formData.isFrame">
            <a-input v-model="formData.frameUrl" placeholder="请输入外链地址" />
          </a-form-item>
          <a-form-item label="组件地址" field="component" v-else>
            <a-input v-model="formData.component" placeholder="请输入组件地址" />
          </a-form-item>
        </transition>
      </template>
      <div class="flex">
        <a-form-item label="显示状态" field="isVisible" :rules="[ {required: true, message: '显示状态不能为空'} ]">
          <a-radio-group v-model="formData.isVisible">
            <a-radio :value="true">显示</a-radio>
            <a-radio :value="false">隐藏</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="菜单状态" field="isDisable" :rules="[ {required: true, message: '菜单状态不能为空'} ]">
          <a-radio-group v-model="formData.isDisable">
            <a-radio :value="false">启用</a-radio>
            <a-radio :value="true">停用</a-radio>
          </a-radio-group>
        </a-form-item>
      </div>
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