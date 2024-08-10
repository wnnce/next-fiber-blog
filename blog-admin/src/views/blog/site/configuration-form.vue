<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { FileItem } from '@arco-design/web-vue'
import ImageUpload from '@/components/ImageUpload.vue'
import { siteApi, type SiteConfigurationItem, type SiteConfigurationType } from '@/api/blog/site'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import ConfigurationFieldForm from '@/views/blog/site/configuration-field-form.vue'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Record<string, SiteConfigurationItem>) => {
  if (record) {
    Object.keys(record).forEach(key => {
      const { type, value } = record[key];
      if (type === 'image') {
        imageRecord.value[key] = [{
          uid: new Date().getTime().toString(),
          status: 'done',
          percent: 1,
          url: value.toString(),
        }]
      }
    })
    Object.assign(formData, record)
  }
  modalShow.value = true;
}
const onClose = () => {
  imageRecord.value = {};
  Object.keys(formData).forEach(key => delete formData[key]);
}
const formData = reactive<Record<string, SiteConfigurationItem>>({})

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    const result = await siteApi.updateConfiguration(formData);
    if (result.code === 200) {
      successMessage('修改成功');
      emits('reload');
      modalShow.value = false;
    }
  } finally {
    submitButtonLoading.value = false;
  }
}

const imageRecord = ref<Record<string, FileItem[]>>({});

const fieldRef = ref();
const showFieldForm = () => {
  fieldRef.value.show();
}
const handlerAddField = (field: string, name: string, type: string) => {
  formData[field] = {
    name: name,
    type: type as SiteConfigurationType,
    extend: true,
    value: type === 'image' ? 0 : '',
  }
}
const handleDeleteField = (field: string) => {
  delete formData[field];
}

defineExpose({
  show
})

</script>

<template>
  <a-modal title="修改站点配置" v-model:visible="modalShow" @close="onClose" :footer="false" width="800px">
    <a-form :model="formData" auto-label-width @submit-success="formSubmit" >
      <div style="max-height: 600px; overflow-y: auto; padding: 8px">
        <div class="flex" v-for="(item, key) in formData" :key="key" style="column-gap: 12px">
          <a-form-item :label="item.name || key" :field="key" style="flex: 1"
                       :rules="[{ required: !item.extend, message: `${item.name || key} 不能为空` }]"
          >
            <template v-if="item.type === 'text'">
              <a-input v-model="item.value as string" placeholder="请输入内容" />
            </template>
            <template v-else-if="item.type === 'number'">
              <a-input-number v-model="item.value as number" placeholder="请输入内容" />
            </template>
            <template v-else-if="item.type === 'html'">
              <a-textarea v-model="item.value as string" placeholder="请输入HTML字符" />
            </template>
            <template v-else-if="item.type === 'markdown'">
              <MarkdownEditor v-model="item.value as string" mode="ir" :min-height="200" :fixed-toolbar="false" />
            </template>
            <template v-else-if="item.type === 'image'">
              <image-upload :file-list="imageRecord[key]" v-model:file-url="item.value as string" />
            </template>
            <template v-else-if="item.type === 'color'">
              <a-color-picker v-model="item.value as string" />
            </template>
          </a-form-item>
          <div style="flex-shrink: 0" v-if="item.extend">
            <a-popconfirm content="是否确认删除配置项?" type="error" position="lb"
                          :ok-button-props="{ status: 'danger' }"
                          @ok="handleDeleteField(key)"
            >
              <a-button type="text" shape="round" status="danger">
                <template #icon><icon-delete /></template>
              </a-button>
            </a-popconfirm>
          </div>
        </div>

      </div>
      <div class="flex justify-between" style="width: 100%; column-gap: 24px; margin-top: 16px">
        <a-button size="large" long @click="showFieldForm">
          <template #icon><icon-plus /></template>
          添加配置项
        </a-button>
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
    <configuration-field-form ref="fieldRef" @add-field="handlerAddField" />
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