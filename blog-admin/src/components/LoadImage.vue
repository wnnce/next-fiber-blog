<script setup lang="ts">
import { computed, ref } from 'vue'
import { sliceThumbnailImageUrl } from '@/assets/script/util'

declare type ImageMode = 'contain' | 'cover' | 'fill' | 'none' | 'scale-down';

declare type ImageStatus = 'loading' | 'done' | 'error';

interface Props {
  // 图片链接
  src: string,
  // 是否只显示缩略图
  thumbnail?: boolean;
  // 是否开启过渡效果
  animation?: boolean;
  // 宽度
  width?: number | string;
  // 高度
  height?: number | string;
  // 是否懒加载
  lazy?: boolean;
  // 图片描述
  alt?: string;
  // 圆角
  radius?: number | string;
  // 图片显示模式
  mode?: ImageMode;
}
const props = withDefaults(defineProps<Props>(), {
  local: false,
  thumbnail: true,
  animation: true,
  width: '60px',
  height: '60px',
  lazy: true,
  alt: 'image',
  radius: '0',
  mode: 'cover'
})
const emits = defineEmits<{
  // 图片加载完成事件
  (e: 'done'): void,
  // 图片加载错误事件
  (e: 'error'): void
}>()

const _width = computed(() => {
  if (typeof props.width === 'string') {
    return props.width;
  }
  return `${props.width}px`;
})
const _height = computed(() => {
  if (typeof  props.height === 'string') {
    return props.height;
  }
  return `${props.height}px`;
})
const _radius = computed(() => {
  if (typeof props.radius === 'string') {
    return props.radius;
  }
  return `${props.radius}px`;
})
const thumbnailUrl = computed(() => {
  const min = Math.min(parseFloat(_width.value), parseFloat(_height.value));
  return sliceThumbnailImageUrl(props.src, min);
})
const maskBlur = computed(() => {
  if (loadingStatus.value === 'loading') {
    return 'blur(4px)';
  }
  return 'blur(0)';
})

const loadingStatus = ref<ImageStatus>('loading');
const blurMaskShow = ref<boolean>(true);

const handleLoadError = () => {
  loadingStatus.value = 'error';
  emits('error');
}
const handleLoadDone = () => {
  loadingStatus.value = 'done';
  if (!props.thumbnail && props.animation) {
    setTimeout(() => {
      blurMaskShow.value = false;
    }, 300)
  }
  emits('done');
}
</script>

<template>
  <div class="load-image">
    <div class="mask thumb-loading-mask absolute-center" v-if="loadingStatus === 'loading'">
      <icon-loading spin />
    </div>
    <div class="mask thumb-error-mask absolute-center" v-else-if="loadingStatus === 'error'">
      <icon-close-circle />
    </div>
    <!-- 缩略图模式之加载缩略图 -->
    <template v-if="thumbnail">
      <img :src="thumbnailUrl" :loading="lazy ? 'lazy' : undefined"
           :alt="alt"
           @load="handleLoadDone"
           @error="handleLoadError"
      />
    </template>
    <!-- 如果不是缩略图模式也不开启过渡动画的话，直接加载原图 -->
    <template v-else-if="!animation">
      <img :src="props.src" :loading="lazy ? 'lazy' : undefined"
           :alt="alt"
           @load="handleLoadDone"
           @error="handleLoadError"
      />
    </template>
    <!-- 不是缩略图但启用原图过渡动画模式 -->
    <template v-else>
      <div class="mask blur-mask" v-if="blurMaskShow" />
      <div class="mask thumb-image-mask" v-if="loadingStatus === 'loading'"
           :style="{ backgroundImage: `url(${thumbnailUrl})` }"
      />
      <img :src="props.src" :loading="lazy ? 'lazy' : undefined"
           :alt="alt"
           @load="handleLoadDone"
           @error="handleLoadError"
      />
    </template>
  </div>
</template>

<style scoped lang="scss">
.load-image {
  height: v-bind(_height);
  width: v-bind(_width);
  border-radius: v-bind(_radius);
  overflow: hidden;
  position: relative;
  img {
    height: 100%;
    width: 100%;
    object-fit: v-bind(mode);
  }
  .mask {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
  }
  .thumb-loading-mask {
    background-color: var(--color-fill-1);
    color: var(--primary-color);
  }
  .thumb-error-mask {
    background-color: rgb(var(--danger-1));
    color: rgb(var(--danger-5));
  }
  .blur-mask {
    transition: all 300ms ease;
    z-index: 2;
    backdrop-filter: v-bind(maskBlur);
  }
  .thumb-image-mask {
    background-repeat: no-repeat;
    background-size: cover;
    background-position: center;
  }
}
</style>