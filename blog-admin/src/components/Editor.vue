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
  minHeight?: number,
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '请输入内容',
  mode: 'wysiwyg',
  fixedToolbar: true,
})

let vditor!: Vditor;
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
  getMarkdownValue, getHtmlValue, setValue,clear
})

onMounted(() => {
  vditor = new Vditor(templateId, {
    theme: useAppConfigStore().state.pageTheme === 'light' ? 'classic' : 'dark',
    toolbarConfig: {
      pin: props.fixedToolbar
    },
    minHeight: props.minHeight,
    mode: props.mode,
    cache: {
      enable: false
    },
    placeholder: props.placeholder,
    value: props.defaultValue,
    icon: 'ant',
    upload: {
      accept: 'image/png, image/jpeg, image/webp, image/gif',
      multiple: true,
      handler: handleImageUpload,
    }
  })
})
</script>

<template>
  <div :id="templateId" class="editor-container">

  </div>
</template>

<style scoped lang="scss">
</style>