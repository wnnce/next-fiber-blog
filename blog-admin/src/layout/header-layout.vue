<script setup lang="ts">
import type { OptionItem } from '@/assets/script/types'
import { ref } from 'vue'
import { useAppConfigStore } from '@/stores/app-config'
import type { PageTheme } from '@/config/app-config'
import { changeTheme } from '@/assets/script/util'
import { IconSunFill } from '@arco-design/web-vue/es/icon'
import TextDropdown from '@/layout/components/TextDropdown.vue'
import Breadcrumd from '@/layout/components/Breadcrumd.vue'
import { useI18nLanguage } from '@/i18n'
import SettingDrawer from '@/layout/components/SettingDrawer.vue'
import { useLocalUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import LogoutModal from '@/layout/components/LogoutModal.vue'

const configStore = useAppConfigStore();
const i18nLanguage = useI18nLanguage();
const userStore = useLocalUserStore();

const iconSize = 18;

const handleLanguageSelect = (value: string | number) => {
  i18nLanguage.changeLanguage(value.toString());
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
const handleUserOptionSelect = (value: string | number) => {
  value = value.toString();
  if (value === 'logout') {
    logoutRef.value.show();
  }
}

const logoutRef= ref();

const settingRef = ref();
const settingShow = () => {
  settingRef.value.show();
}

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
          :options="i18nLanguage.languageOption"
          :text="i18nLanguage.currentLanguage.value.label"
          @select="handleLanguageSelect"
        />
      </div>
      <div class="theme pointer" @click="handleThemeSwitch">
        <transition name="hidden" mode="out-in">
          <icon-sun-fill class="icon-button" :size="iconSize" v-if="configStore.state.pageTheme === 'light'" style="color: #FFAB40"/>
          <icon-moon-fill class="icon-button" :size="iconSize" v-else style="color: #8A2BE2"/>
        </transition>
      </div>
      <div class="pointer" @click="settingShow">
        <icon-settings :size="iconSize"/>
      </div>
      <div class="user flex item-center">
        <a-avatar :size="32" :image-url="userStore.userInfo ? userStore.userInfo.avatar : ''">User</a-avatar>
        <text-dropdown :options="userOption"
                       :text="userStore.userInfo ? userStore.userInfo.nickname : '登录用户'"
                       @select="handleUserOptionSelect"
        />
      </div>
    </div>
    <setting-drawer ref="settingRef" />
    <logout-modal ref="logoutRef" />
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