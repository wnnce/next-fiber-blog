<script setup lang="ts">

import { reactive, ref } from 'vue'
import type { Config, ConfigForm } from '@/api/system/config/types'
import { configApi } from '@/api/system/config'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Config) => {
  if (record) {
    const { configId, configName, configKey, configValue, remark } = record;
    Object.assign(configForm, { configId, configName, configKey, configValue, remark })
  }
  modalShow.value = true;
}
const onClose = () => {
  Object.assign(configForm, defaultConfigForm);
}

const defaultConfigForm: ConfigForm = {
  configId: undefined,
  configName: '',
  configKey: '',
  configValue: '',
  remark: '',
}
const configForm = reactive<ConfigForm>({ ...defaultConfigForm })

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (configForm.configId) {
      result = await configApi.updateSysConfig(configForm);
    } else {
      result = await configApi.saveSysConfig(configForm);
    }
    if (result.code === 200) {
      successMessage(configForm.configId ? '修改成功' : '保存成功');
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
  <a-modal :title="configForm.configId ? '修改参数配置' : '添加参数配置'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="configForm" layout="vertical" @submit-success="formSubmit">
      <a-form-item label="参数名称" field="configName" :rules="[ {required: true, message: '参数名称不能为空'} ]">
        <a-input v-model="configForm.configName" placeholder="请输入参数名称" />
      </a-form-item>
      <a-form-item label="参数KEY" field="configKey" :rules="[ {required: true, message: '参数KEY不能为空'} ]">
        <a-input v-model="configForm.configKey" placeholder="请输入参数KEY" />
      </a-form-item>
      <a-form-item label="参数值" field="configValue" :rules="[ {required: true, message: '参数值不能为空'} ]">
        <a-textarea v-model="configForm.configValue" placeholder="请输入参数值" />
      </a-form-item>
      <a-form-item label="备注" field="remark">
        <a-textarea v-model="configForm.remark" placeholder="备注" />
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