<script setup lang="ts">
import type { OptionItem } from '@/assets/script/types'
import { onMounted, ref } from 'vue'
import { useAppConfigStore } from '@/stores/app-config'
import type { PageTheme } from '@/config/app-config'
import { changeTheme } from '@/assets/script/util'
import { IconSunFill } from '@arco-design/web-vue/es/icon'
import TextDropdown from '@/layout/components/TextDropdown.vue'
import Breadcrumd from '@/layout/components/Breadcrumd.vue'

const configStore = useAppConfigStore();

const iconSize = 18;

const languageActiveIndex = ref<number>(0);
const languageOption: OptionItem[] = [
  { label: '中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US'}
]
const handleLanguageSelect = (value: string | number) => {
  value = value.toString();
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

const userOption: OptionItem[] = [
  { label: '个人中心', value: 'user-info' },
  { label: '修改密码', value: 'reset-password' },
  { label: '退出登录', value: 'logout' }
]
const handleUserOptionClick = (value: string | number) => {
  value = value.toString();
  console.log(value);
}

onMounted(() => {
  changeLanguageActiveIndex(configStore.state.language);
})

</script>

<template>
  <header class="header flex justify-between">
    <div class="left flex item-center">
      <div class="logo flex item-center">
        <div class="logo-image">
          <img src="/images/logo.png" alt="logo">
        </div>
        <div class="logo-title">
          博客后台管理
        </div>
      </div>
      <div class="bread">
        <breadcrumd />
      </div>
    </div>
    <div class="option flex item-center">
      <div class="notice pointer">
        <a-badge dot :count="1">
          <icon-notification size="18" />
        </a-badge>
      </div>
      <div class="language">
        <text-dropdown
          :options="languageOption"
          :text="languageOption[languageActiveIndex].label"
          @select="handleLanguageSelect"
        />
      </div>
      <div class="theme pointer" @click="handleThemeSwitch">
        <transition name="hidden" mode="out-in">
          <icon-sun-fill class="icon-button" :size="iconSize" v-if="configStore.state.pageTheme === 'light'" style="color: #FFAB40"/>
          <icon-moon-fill class="icon-button" :size="iconSize" v-else style="color: #8A2BE2"/>
        </transition>
      </div>
      <div class="pointer">
        <icon-settings :size="iconSize"/>
      </div>
      <div class="user flex item-center">
        <a-avatar :size="32">A</a-avatar>
        <text-dropdown :options="userOption" text="Admin" @select="handleUserOptionClick" />
      </div>
    </div>
  </header>
</template>

<style scoped lang="scss">
header {
  background-color: var(--card-color);
  padding: var(--space-sm) var(--space-xl);
  width: 100%;
  .left {
    column-gap: var(--space-xxl);
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
  }
  .option {
    column-gap: var(--space-xl);
    .user {
      column-gap: var(--space-sm);
    }
  }
}
</style>