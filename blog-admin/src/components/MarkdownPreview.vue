<script setup lang="ts">
import Vditor from 'vditor/dist/method.min';
import { onMounted, watch } from 'vue'

interface Props {
  markdown: string
}
const props = defineProps<Props>();

// 及时预览
watch(props, () => {
  renderMarkdown();
})

const templateId: string = (Math.random() * 100).toFixed(0);

const renderMarkdown = () => {
  Vditor.preview(<HTMLDivElement>document.getElementById(templateId), props.markdown, {
    mode: 'dark',
    hljs: {
      lineNumber: true,
      style: 'native'
    },
  })
}
onMounted(() => {
  renderMarkdown();
})

</script>

<template>
  <div class="preview-container" :id="templateId">

  </div>
</template>

<style scoped lang="scss">
.preview-container {
  color: var(--text-color) !important;
}
</style>