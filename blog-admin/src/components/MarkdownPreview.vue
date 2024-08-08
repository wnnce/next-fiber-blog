<script setup lang="ts">
import Vditor from 'vditor/dist/method.min'
import { onMounted } from 'vue'

interface Props {
  markdown: string
}

const props = defineProps<Props>();

const templateId: string = (Math.random() * 100).toFixed(0);

onMounted(() => {
  Vditor.preview(<HTMLDivElement>document.getElementById(templateId), props.markdown, {
    mode: 'dark',
    hljs: {
      lineNumber: true,
      style: 'native'
    },
  })
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