<script setup lang="ts">
import type { OptionItem } from '@/assets/script/types'
import { onMounted, ref } from 'vue'
import { useAppConfigStore } from '@/stores/app-config'
import type { PageTheme } from '@/config/app-config'
import { changeTheme } from '@/assets/script/util'

const configStore = useAppConfigStore();

const languageActiveIndex = ref<number>(0);
const languageOption: OptionItem[] = [
  { label: '中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US'}
]
const handleLanguageSelect = (value: string) => {
  console.log(value);
  changeLanguageActiveIndex(value);
  configStore.state.language = value;
}
const changeLanguageActiveIndex = (value: string) => {
  languageActiveIndex.value = languageOption.findIndex(item => item.value === value)
}

const handleThemeSwitch = () => {
  const newTheme: PageTheme = configStore.state.pageTheme === 'light' ? 'dark' : 'light';
  configStore.state.pageTheme = newTheme;
  changeTheme(newTheme);
}

onMounted(() => {
  changeLanguageActiveIndex(configStore.state.language);
})

</script>

<template>
  <header class="header flex justify-between shadow-md">
    <div class="logo flex item-center">
      <div class="logo-image">
        <img src="/images/logo.png" alt="logo">
      </div>
      <div class="logo-title">
        博客后台管理
      </div>
    </div>
    <div class="option flex">
      <div class="notice pointer">
        <a-badge dot :count="1">
          <icon-notification size="18" />
        </a-badge>
      </div>
      <div class="language">
        <a-dropdown :popup-max-height="false" @select="handleLanguageSelect">
          <span class="pointer">{{ languageOption[languageActiveIndex].label }} <icon-down/></span>
          <template #content>
            <a-doption v-for="(item, index) in languageOption" :key="index" :value="item.value">
              {{ item.label }}
            </a-doption>
          </template>
        </a-dropdown>
      </div>
      <div class="theme pointer" @click="handleThemeSwitch">
        <icon-sun-fill size="18" v-if="configStore.state.pageTheme === 'light'" style="color: red"/>
        <icon-moon-fill size="18" v-else style="color: #1E90FF"/>
      </div>
    </div>
  </header>
</template>

<style scoped lang="scss">
header {
  background-color: var(--card-color);
  padding: var(--space-sm) var(--space-xl);
  border: 1px solid var(--border-color);
  width: 100%;
  .logo {
    column-gap: var(--space-sm);
    font-size: 18px;
    font-weight: 500;
    .logo-image {
      height: 28px;
      width: 28px;
      overflow: hidden;
      img {
        height: 100%;
        width: 100%;
      }
    }
  }
  .option {
    align-items: center;
    column-gap: var(--space-xl);
  }
}
</style>