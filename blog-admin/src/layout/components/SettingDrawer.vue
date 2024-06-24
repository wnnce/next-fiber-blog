<script setup lang="ts">

import { ref } from 'vue'
import { useAppConfigStore } from '@/stores/app-config'
import { useI18nLanguage } from '@/i18n'
import type { PageTheme } from '@/config/app-config'
import { changeTheme } from '@/assets/script/util'

const configStore = useAppConfigStore();
const i18nLanguage = useI18nLanguage();

const drawerShow = ref<boolean>(false);
const show = () => {
  drawerShow.value = true;
}
const onClose = () => {
}

const handleThemeSwitchChange = (value: string) => {
  const newTheme = value as PageTheme;
  configStore.state.pageTheme = newTheme;
  changeTheme(newTheme);
}

defineExpose({
  show
})
</script>

<template>
  <a-drawer title="网站设置" v-model:visible="drawerShow"
            :mask-closable="false" unmount-on-close class="onClose"
            :footer="false" width="300px"
  >
    <div class="flex flex-column form">
      <div class="input-item">
        <label>语言设置</label>
        <a-select v-model="configStore.state.language">
          <a-option v-for="item in i18nLanguage.languageOption" :key="item.value" :value="item.value" :label="item.label" />
        </a-select>
      </div>
      <div class="input-item">
        <label>侧边栏宽度</label>
        <a-input-number v-model="configStore.state.sideWidth" placeholder="请输入侧边栏宽度" :min="100" :max="500" />
      </div>
      <div class="input-item">
        <label>页面缓存最大数量</label>
        <a-input-number v-model="configStore.state.maxCachePage" placeholder="请输入数量" :min="5" :max="30" />
      </div>
      <div class="flex justify-between">
        <label>
          主题设置
        </label>
        <a-switch :model-value="configStore.state.pageTheme" :checked-value="'light'" :unchecked-value="'dark'"
                  @change="handleThemeSwitchChange"
        >
          <template #checked>
            Light
          </template>
          <template #unchecked>
            Dark
          </template>
        </a-switch>
      </div>
      <div class="flex justify-between">
        <label>
          开启Tabs
        </label>
        <a-switch v-model="configStore.state.showPageTag" :checked-value="true" :unchecked-value="false" />
      </div>
    </div>
  </a-drawer>
</template>

<style scoped lang="scss">
.form {
  row-gap: 20px;
  .input-item {
    display: flex;
    flex-direction: column;
    row-gap: var(--space-xs);
  }
}

</style>