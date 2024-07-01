<script setup lang="ts">

import { ref } from 'vue'
import type { FileItem } from '@arco-design/web-vue'

const uploadFileList = ref<FileItem[]>([]);
</script>

<template>

  <a-upload :file-list="uploadFileList" draggable :auto-upload="false" multiple
            tip="" :show-file-list="false"
            accept="image/png, image/jpeg, image/webp, image/gif"
            style="height: 100px; width: 100px"
  >
    <template #upload-button>
      <div class="button-container radius-md">
        <icon-plus />
        <span>点击或者拖拽上传文件</span>
      </div>
    </template>
  </a-upload>
</template>

<style scoped lang="scss">
.button-container {
  width: 100px;
  height: 100px;
  border: 1px dashed var(--color-border-2);
  background-color: var(--color-fill-2);
  transition: all 500ms ease;
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  justify-content: center;
  align-items: center;
  padding: var(--space-sm);
  &:hover {
    border-color: var(--color-border-3);
    background-color: var(--color-fill-3);
  }
}
</style>