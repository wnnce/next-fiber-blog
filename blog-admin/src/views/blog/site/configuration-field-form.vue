<script setup lang="ts">
import { reactive, ref } from 'vue'
import DictSelect from '@/components/DictSelect.vue'

const emits = defineEmits<{
  (e: 'addField', field: string, name: string, type: string): void
}>()

const modalShow = ref<boolean>(false);
const show = () => {
  modalShow.value = true;
}
const onClose = () => {
  formData.field = '';
  formData.name = '';
  formData.type = '';
}
const formData = reactive({
  field: '',
  name: '',
  type: '',
})

const formSubmit = () => {
  emits('addField', formData.field, formData.name, formData.type);
  modalShow.value = false;
}

defineExpose({
  show
})

</script>

<template>
  <a-modal title="添加配置项" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" auto-label-width @submit-success="formSubmit" >
      <a-form-item label="字段名称" field="field" :rules="[{ required: true, message: '字段名称不能为空' }]">
        <a-input v-model="formData.field" placeholder="请输入字段名称" />
      </a-form-item>
      <a-form-item label="配置名称" field="name" :rules="[{ required: true, message: '配置名称不能为空' }]">
        <a-input v-model="formData.name" placeholder="请输入配置名称" />
      </a-form-item>
      <a-form-item label="类型" field="type" :rules="[{ required: true, message: '配置类型不能为空' }]">
        <dict-select dict-key="site_configuration_type" v-model="formData.type"
                     type="string" placeholder="请选择配置类型"
        />
      </a-form-item>
      <div class="flex justify-between" style="width: 100%; column-gap: 24px; margin-top: 16px">
        <a-button html-type="submit" type="primary" size="large" long>
          <template #icon><icon-save /></template>
          确定
        </a-button>
        <a-button size="large" long @click="modalShow = false">
          <template #icon><icon-close /></template>
          取消
        </a-button>
      </div>
    </a-form>
  </a-modal>
</template>

<style scoped lang="scss">

</style>