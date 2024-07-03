<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { FileItem } from '@arco-design/web-vue'
import { useArcoMessage } from '@/hooks/message'
import { fileUpload } from '@/api/request'
import { sliceImageUrl } from '@/assets/script/util'

const { errorMessage } = useArcoMessage();

interface ImageUploadProps {
  tip?: string;
  showTip?: boolean;
  width?: string | number;
  height?: string | number;
  fileList: string | string[];
  limit?: number;
  circle?: boolean;
}

const props = withDefaults(defineProps<ImageUploadProps>(), {
  tip: '点击或者拖拽上传文件',
  showTip: false,
  height: '100px',
  width: '100px',
  limit: 0,
  circle: false
})
const emits = defineEmits<{
  // 更新modelValue
  (e: 'update:fileList', value: string | string[]): void,
  // 文件状态发送改变时触发
  (e: 'change', value: FileItem): void,
  // 某个文件上传成功时触发
  (e: 'upload', value: FileItem): void,
  // 文件全部上传成功时触发
  (e: 'ok'): void,
}>()

const _tip = computed(() => {
  if (!isDragEnter.value) {
    return props.tip;
  }
  return '释放文件并开始上传';
})
const _width = computed(() => {
  if (typeof props.width === 'string') {
    return props.width;
  }
  return `${props.width}px`
})
const _height = computed(() => {
  if (typeof props.height === 'string') {
    return props.height;
  }
  return `${props.height}px`
})
const _limit = computed(() => {
  if (typeof props.fileList === 'string') {
    return 1;
  }
  return props.limit;
})
const _borderRadius = computed(() => {
  return props.circle ? '50%' : '12px';
})

const isDragEnter = ref<boolean>(false);
const fileUploading = ref<boolean>(false);
const uploadFileList = ref<FileItem[]>([]);

const handleUpload = async () => {
  if (fileUploading.value) {
    return;
  }
  if (!uploadFileList.value || uploadFileList.value.length === 0) {
    return;
  }
  const waitUploadList = uploadFileList.value.filter(item => item.status === 'init')
  if (!waitUploadList || waitUploadList.length === 0) {
    // 检查是否全部上传完成
    checkFileAllUploadDone();
    return;
  }
  fileUploading.value = true;
  for (const item of waitUploadList) {
    if (!item.file) {
      item.status = 'error';
      item.url && (window.URL.revokeObjectURL(item.url))
      continue
    }
    item.status = 'uploading';
    const formData = new FormData();
    formData.append('image', item.file);
    try {
      const result = await fileUpload('/other/upload/image', formData, (event: ProgressEvent) => {
        const { lengthComputable, loaded, total } = event;
        if (lengthComputable) {
          item.percent = parseFloat((loaded / total).toFixed(2));
        }
      })
      const { code, data } = result;
      if (code === 200 && data) {
        const blobUrl = item.url;
        (blobUrl && blobUrl.startsWith('blob:')) && (URL.revokeObjectURL(blobUrl));
        item.url = data;
        item.status = 'done';
        item.file = undefined;
        if (typeof props.fileList === 'string') {
          emits('update:fileList', data);
        } else {
          const doneFileList: string[] = []
          uploadFileList.value.forEach(item => {
            if (!item.status || item.status !== 'done' || !item.url) {
              return;
            }
            doneFileList.push(item.url);
          })
          emits('update:fileList', doneFileList);
        }
      } else {
        item.status = 'error';
      }
    } catch (err) {
      console.log(err);
      item.status = 'error';
    } finally {
      emits('change', item);
    }
  }
  fileUploading.value = false;
  // 继续递归调用
  handleUpload();
}

const onUploadChange = (_: FileItem[], fileItem: FileItem) => {
  if (uploadFileList.value.length >= _limit.value) {
    errorMessage('上传文件数量达到最大限制');
    return;
  }
  uploadFileList.value.push(fileItem);
  handleUpload();
}

const handleDeleteUploadFile = (uid: string) => {
  const findIndex = uploadFileList.value.findIndex(item => item.uid === uid)
  if (findIndex < 0) {
    return;
  }
  const { status, url } = uploadFileList.value[findIndex];
  if (status === 'uploading') {
    errorMessage('文件上传中无法删除');
    return;
  }
  // 删除当前图片
  uploadFileList.value.splice(findIndex, 1);
  // 如果url存在并且是blob资源 那么就释放资源
  (url && url.startsWith('blob:')) && (URL.revokeObjectURL(url));
  // 检查是否全部上传
  checkFileAllUploadDone();
}

// 检查文件是否全部上传完成
const checkFileAllUploadDone = () => {
  const doneCount = uploadFileList.value.filter(item => item.status && item.status === 'done').length;
  if (doneCount > 0 && doneCount === uploadFileList.value.length) {
    emits('ok');
  }
}

// 处理文件上传重试
const handleRetryUpload = (item: FileItem) => {
  if (item.status && item.status === 'error') {
    item.percent = 0;
    item.status = 'init';
    handleUpload();
  }
}

const formatServerImageUrl = (imageUrl: string | undefined) => {
  if (!imageUrl) {
    return '';
  }
  if (imageUrl.startsWith('blob:')) {
    return imageUrl;
  }
  return sliceImageUrl(imageUrl);
}
</script>

<template>
  <div class="upload-container">
    <transition-group name="list">
      <div class="common-card image-card" v-for="item in uploadFileList" :key="item.uid">
        <div class="image-mask init-mask absolute-center" v-if="item.status === 'init'">
          <icon-loading spin />
          <span>待上传</span>
        </div>
        <div class="image-mask loading-mask absolute-center" v-else-if="item.status === 'uploading'">
          <a-progress :percent="item.percent" type="circle" size="small"  />
        </div>
        <div class="done-mask" v-else-if="item.status === 'done' && !circle">
          <div class="done-icon">
            <icon-check />
          </div>
        </div>
        <div class="image-mask fail-mask absolute-center flex-column" v-else-if="item.status === 'error'">
          <div class="flex" style="column-gap: 8px; align-items: flex-end">
            <icon-refresh :size="16" class="pointer" @click="handleRetryUpload(item)"/>
            <span>|</span>
            <icon-close-circle :size="16" class="danger-color" />
          </div>
          <span class="danger-color">上传失败</span>
        </div>
        <div class="delete-pop flex justify-center pointer danger-color" @click="handleDeleteUploadFile(item.uid)"
             v-if="item.status !== 'uploading'"
        >
          <icon-delete />
        </div>
        <img :src="formatServerImageUrl(item.url)" alt="upload">
      </div>
    </transition-group>
    <a-upload :file-list="uploadFileList" draggable :auto-upload="false" multiple
              tip="" :show-file-list="false" :limit="_limit"
              accept="image/png, image/jpeg, image/webp, image/gif"
              style="height: 100px; width: 100px"
              @change="onUploadChange" v-show="uploadFileList.length < _limit"
    >
      <template #upload-button>
        <div class="common-card button-card" :class="isDragEnter ? 'drag-enter' : ''">
          <div class="drag-mask"
               @dragenter="isDragEnter = true"
               @dragleave="isDragEnter = false"
               @drop="isDragEnter = false"
          />
          <icon-plus />
          <span v-if="showTip">{{ _tip }}</span>
        </div>
      </template>
    </a-upload>
  </div>
</template>

<style scoped lang="scss">
.upload-container {
  display: flex;
  gap: var(--space-sm);
  flex-wrap: wrap;
  .common-card {
    height: v-bind(_height);
    width: v-bind(_width);
    overflow: hidden;
    border-radius: v-bind(_borderRadius);
  }
  .image-card {
    --delete-mask-opacity: 0;
    position: relative;
    &:hover {
      --delete-mask-opacity: 1;
    }
    img {
      height: 100%;
      width: 100%;
      object-fit: cover;
    }
    .danger-color {
      color: rgb(var(--danger-5));
    }
    .image-mask {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      z-index: 2;
      color: rgb(240, 240, 240);
      font-size: 12px;
      border-radius: v-bind(_borderRadius);
    }
    .init-mask {
      background-color: rgba(0, 0, 0, 0.5);
      flex-direction: column;
      row-gap: var(--space-xs);
    }
    .loading-mask {
      background-color: rgba(0, 0, 0, 0.5);
    }
    .done-mask {
      position: absolute;
      top: 0;
      right: 0;
      height: 0;
      width: 0;
      border-bottom: 36px solid transparent;
      border-left: 40px solid transparent;
      border-top: 36px solid rgb(var(--success-4));
      .done-icon {
        position: absolute;
        top: -36px;
        right: 6px;
        color: white;
      }
    }
    .fail-mask {
      row-gap: var(--space-xs);
      background-color: rgba(0, 0, 0, 0.6);
    }
    .delete-pop {
      z-index: 20;
      background-color: rgba(0, 0, 0, 0.3);
      backdrop-filter: blur(2px);
      position: absolute;
      left: 0;
      right: 0;
      bottom: 0;
      padding: 4px 0;
      opacity: var(--delete-mask-opacity);
      transition: opacity 300ms ease;
    }
  }
  .button-card {
    border: 1px dashed var(--color-border-2);
    background-color: var(--color-fill-2);
    transition: all 300ms ease;
    display: flex;
    flex-direction: column;
    row-gap: var(--space-sm);
    justify-content: center;
    align-items: center;
    padding: var(--space-sm);
    color: var(--color-text-2);
    position: relative;
    &:hover {
      border-color: var(--color-border-3);
      background-color: var(--color-fill-3);
    }
    .drag-mask {
      position: absolute;
      left: 0;
      top: 0;
      right: 0;
      bottom: 0;
    }
  }
  .drag-enter {
    border-color: rgb(var(--primary-4)) !important;
    background-color: rgb(var(--primary-1)) !important;
  }
}
</style>