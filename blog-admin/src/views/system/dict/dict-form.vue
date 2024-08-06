<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { Dict, DictForm } from '@/api/system/dict/types'
import { dictApi } from '@/api/system/dict'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import DictSelect from '@/components/DictSelect.vue'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Dict) => {
  if (record) {
    const { dictId, dictName, dictKey, sort, status, remark } = record;
    Object.assign(formData, { dictId, dictName, dictKey, sort, status, remark })
  }
  modalShow.value = true;
}
const onClose = () => {
  Object.assign(formData, defaultConfigForm);
}

const defaultConfigForm: DictForm = {
  dictId: undefined,
  dictName: '',
  dictKey: '',
  sort: 1,
  status: 0,
  remark: '',
}
const formData = reactive<DictForm>({ ...defaultConfigForm })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  dictName: { required: true, message: '字典名称不能为空' },
  dictKey: { required: true, message: '字典KEY不能为空' },
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.dictId) {
      result = await dictApi.updateDict(formData);
    } else {
      result = await dictApi.saveDict(formData);
    }
    if (result.code === 200) {
      successMessage(formData.dictId ? '修改成功' : '保存成功');
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
  <a-modal :title="formData.dictId ? '修改字典' : '添加字典'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" layout="vertical" @submit="formSubmit" :rules="formRules">
      <a-form-item label="字典名称" field="dictName">
        <a-input v-model="formData.dictName" placeholder="请输入字典名称" />
      </a-form-item>
      <a-form-item label="字典KEY" field="dictKey">
        <a-input v-model="formData.dictKey" placeholder="请输入字典KEY" />
      </a-form-item>
      <a-form-item label="备注" field="remark">
        <a-textarea v-model="formData.remark" placeholder="备注" />
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