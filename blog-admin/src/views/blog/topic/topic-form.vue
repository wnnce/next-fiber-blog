<script setup lang="ts">

import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { FileItem } from '@arco-design/web-vue'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import ImageUpload from '@/components/ImageUpload.vue'
import DictSelect from '@/components/DictSelect.vue'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import type { Topic, TopicForm } from '@/api/blog/topic/types'
import { topicApi } from '@/api/blog/topic'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Topic) => {
  if (record) {
    const { topicId, content, imageUrls, location, mode, isHot, isTop, sort, status } = record;
    Object.assign(formData, { topicId, content, imageUrls: imageUrls || [], location, mode, isHot, isTop, sort, status })
    editorRef.value.setValue(content);
    formatImageToFileList();
  }
  modalShow.value = true;
}

const onClose = () => {
  fileList.value = [];
  Object.assign(formData, defaultFormData);
  formRef.value.clearValidate();
  editorRef.value.clear();
}

const formRef = ref();
const defaultFormData: TopicForm = {
  topicId: undefined,
  content: '',
  imageUrls: [],
  location: undefined,
  isHot: false,
  isTop: false,
  mode: undefined,
  sort: 1,
  status: 0,
}
const formData = reactive<TopicForm>({ ...defaultFormData })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  content: {
    validator: (value, callback) => {
      const markdownValue = editorRef.value.getMarkdownValue();
      if (!markdownValue || markdownValue.trim().length === 0) {
        callback("动态内容不能为空");
      } else {
        formData.content = markdownValue;
      }
    }
  },
  imageUrls: {
    validator: (value, callback) => {
      if (formData.mode === 2 && (!formData.imageUrls || formData.imageUrls.length === 0)) {
        callback('照片墙模式图片不能为空');
      }
    }
  },
  mode: { required: true, message: '动态类型不能为空' },
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.topicId) {
      result = await topicApi.updateTopic(formData);
    } else {
      result = await topicApi.saveTopic(formData);
    }
    if (result.code === 200) {
      successMessage(formData.topicId ? '修改成功' : '保存成功');
      emits('reload');
      modalShow.value = false;
    }
  } finally {
    submitButtonLoading.value = false;
  }
}

const fileList = ref<FileItem[]>([]);
const formatImageToFileList = () => {
  if (!formData.imageUrls || formData.imageUrls.length === 0) {
    return;
  }
  fileList.value = formData.imageUrls.map(item => {
    return {
      uid: (1e7 + ((Math.random() * 10000).toFixed(0) << 4)) % 1e6.toString(),
      status: 'done',
      percent: 1,
      url: item
    }
  })
}

const editorRef = ref();

defineExpose({
  show
})

</script>

<template>
  <a-modal :title="formData.topicId ? '修改动态' : '发布动态'" v-model:visible="modalShow" width="auto"
           @close="onClose" :footer="false"
  >
    <div class="topic-content">
      <a-form ref="formRef" :model="formData" auto-label-width :rules="formRules" @submit-success="formSubmit">
        <div class="topic-publish-form flex flex-column radius-md">
          <a-form-item field="content" hide-label>
            <MarkdownEditor ref="editorRef" :min-height="120" placeholder="说点什么吧..." mode="ir"
                    hide-code-preview
                    style="flex: 1"
            />
          </a-form-item>
          <a-form-item field="imageUrls" hide-label>
            <image-upload :file-list="fileList" :file-url="formData.imageUrls"
                          :limit="16"
            />
          </a-form-item>
          <div class="flex justify-between" style="column-gap: 16px">
            <a-form-item label="地点" field="location" hide-label>
              <a-select v-model="formData.location" placeholder="请选择发布地点" />
            </a-form-item>
            <a-form-item label="显示顺序" field="sort">
              <a-input-number v-model="formData.sort" placeholder="请输入显示顺序" />
            </a-form-item>
          </div>
          <div class="topic-form-option flex justify-between">
            <div class="option-left flex">
              <a-form-item field="mode" hide-label>
                <dict-select dict-key="topic_mode" v-model="formData.mode" type="number"
                             placeholder="动态类型"
                             width="120px"
                />
              </a-form-item>
              <a-form-item label="置顶" field="is_top">
                <a-switch v-model="formData.isTop" />
              </a-form-item>
              <a-form-item label="热门" field="is_hot">
                <a-switch v-model="formData.isHot" />
              </a-form-item>
              <a-form-item label="状态" field="status">
                <a-switch v-model="formData.status" :checked-value="0" :unchecked-value="1" />
              </a-form-item>
            </div>
            <div class="option-right">
              <a-button type="primary" size="large"
                        :loading="submitButtonLoading"
                        style="padding: 12px 48px !important; border-radius: 10px"
                        html-type="submit"
              >
                <template #icon>
                  <icon-send />
                </template>
              </a-button>
            </div>
          </div>
        </div>
      </a-form>
    </div>
  </a-modal>
</template>

<style scoped lang="scss">
.topic-content {
  width: 700px;
  margin: 0 auto;
}
.topic-form-option {
  margin-bottom: -16px;
  .option-left {
    column-gap: var(--space-mm);
  }
}
.topic-publish-form {
  width: 100%;
  padding: var(--space-mm);
  border: 1px solid #e8e8e8;
  :deep(.editor-container) {
    border: none !important;
  }
  :deep(.vditor-toolbar) {
    display: none;
  }
  :deep(.vditor-reset) {
    background-color: transparent !important;
    color: var(--text-color) !important;
    padding: 0 !important;
  }
}
</style>