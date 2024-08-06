<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { DictValue, DictValueForm } from '@/api/system/dict/types'
import { dictApi } from '@/api/system/dict'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import DictSelect from '@/components/DictSelect.vue'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (dictId: number, dictKey: string, record?: DictValue, ) => {
  if (record) {
    const { id, label, value, sort, status, remark } = record;
    Object.assign(formData, { id, label, value, sort, status, remark })
  }
  formData.dictId = dictId;
  formData.dictKey = dictKey;
  modalShow.value = true;
}
const onClose = () => {
  Object.assign(formData, defaultConfigForm);
}

const defaultConfigForm: DictValueForm = {
  id: undefined,
  dictId: undefined,
  dictKey: undefined,
  label: undefined,
  value: undefined,
  sort: 1,
  status: 0,
  remark: undefined,
}
const formData = reactive<DictValueForm>({ ...defaultConfigForm })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  dictId: { required: true, message: '字典ID不能为空' },
  dictKey: { required: true, message: '字典KEY不能为空' },
  label: { required: true, message: '数据标签不能为空' },
  value: { required: true, message: '数据值不能为空' },
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.id) {
      result = await dictApi.updateDictValue(formData);
    } else {
      result = await dictApi.saveDictValue(formData);
    }
    if (result.code === 200) {
      successMessage(formData.id ? '修改成功' : '保存成功');
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
  <a-modal :title="formData.id ? '修改字典数据' : '添加字典数据'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" auto-label-width @submit-success="formSubmit" :rules="formRules">
      <a-form-item label="字典KEY" field="dictKey">
        <a-input :model-value="formData.dictKey" disabled />
      </a-form-item>
      <a-form-item label="标签" field="label">
        <a-input v-model="formData.label" placeholder="请输入数据标签" />
      </a-form-item>
      <a-form-item label="数据值" field="value">
        <a-input v-model="formData.value" placeholder="请输入数据值" />
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