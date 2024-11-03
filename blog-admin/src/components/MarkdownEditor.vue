<script setup lang="ts">
import Vditor from 'vditor'
import "vditor/src/assets/less/index.less"
import { onMounted, watch } from 'vue'
import { useAppConfigStore } from '@/stores/app-config'
import { fileUpload } from '@/api/request'

interface Props {
  placeholder?: string,
  mode?: 'wysiwyg' | 'ir' | 'sv',
  defaultValue?: string,
  fixedToolbar?: boolean,
  height?: number,
  hideCodePreview?: boolean,
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '请输入内容',
  mode: 'wysiwyg',
  fixedToolbar: true,
  hideCodePreview: false,
})

const modeValue = defineModel<string>({
  required: false
})
const updateModelValue = () => {
  // 如果等于 undefined 说明没有传递 modelValue 不触发更新
  if (modeValue.value !== undefined) {
    modeValue.value = getMarkdownValue();
  }
}

let vditor!: Vditor;
let wordCount: number = 0;
const templateId: string = new Date().getTime().toString();

// 监听主题变化
watch(useAppConfigStore().state, (newValue) => {
  vditor.setTheme(newValue.pageTheme === 'light' ? 'classic' : 'dark')
})

const handleImageUpload = async (files: File[]) : Promise<string> => {
  for (let file of files) {
    const fileName = file.name;
    vditor.tip(`${fileName} 上传中...`)
    const formData = new FormData();
    formData.append("image", file);
    try {
      const result = await fileUpload('/base/upload/image', formData)
      const { code, data } = result;
      if (code === 200 && data) {
        vditor.tip(`${fileName} 上传完成`)
        vditor.insertValue(`\n![${fileName}](${data})`)
      }
    } catch (err) {
      vditor.tip(`${fileName} 上传失败`)
      console.log(err)
    }
  }
  return '文件上传结束';
}

const getMarkdownValue = (): string => {
  return vditor.getValue();
}

const getHtmlValue = (): string => {
  return vditor.getHTML();
}

const getWordCount = (): number => {
  return wordCount;
}

const getVditor = (): Vditor => {
  return vditor;
}

const setValue = (value: string) => {
  vditor.setValue(value, false)
}

const clear = () => {
  vditor.setValue("", true)
}

defineExpose({
  getMarkdownValue, getHtmlValue, getWordCount, getVditor, setValue, clear
})

onMounted(() => {
  vditor = new Vditor(templateId, {
    counter: {
      enable: true,
      type: 'text',
      after: (length) => {
        wordCount = length;
      }
    },
    preview: {
      hljs: {
        lineNumber: true,
      },
      markdown: {
        codeBlockPreview: !props.hideCodePreview,
      }
    },
    theme: useAppConfigStore().state.pageTheme === 'light' ? 'classic' : 'dark',
    toolbarConfig: {
      pin: props.fixedToolbar
    },
    height: props.height,
    mode: props.mode,
    cache: {
      enable: false
    },
    placeholder: props.placeholder,
    value: modeValue.value || props.defaultValue,
    icon: 'ant',
    upload: {
      accept: 'image/png, image/jpeg, image/webp, image/gif',
      multiple: true,
      handler: handleImageUpload,
    },
    input: updateModelValue,
  })
})
</script>

<template>
  <div :id="templateId" class="editor-container">

  </div>
</template>

<style scoped lang="scss">
.editor-container {
  width: 100%;
}
:deep(.vditor-reset) {
  color: var(--text-color) !important;
}
:deep(.vditor-content) {
  max-height: 100%;
  overflow-y: auto;
}
.vditor {
  --border-color: var(--color-border-2);
  --second-color: rgba(88, 96, 105, 0.36);

  --panel-background-color: #ffffff;
  --panel-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);

  --toolbar-background-color: var(--card-color);
  --toolbar-icon-color: var(--color-text-2);
  --toolbar-icon-hover-color: #4285f4;
  --toolbar-height: 35px;
  --toolbar-divider-margin-top: 8px;

  --textarea-background-color: var(--card-color);

  &--dark {
    --panel-background-color: var(--card-color);
    --panel-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);

    --textarea-background-color: var(--card-color);
  }
}
</style>