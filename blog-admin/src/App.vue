<script setup lang="ts">
import { RouterView } from 'vue-router'
import { computed, onMounted } from 'vue'
import { useAppConfigStore } from '@/stores/app-config'
import enUS from '@arco-design/web-vue/es/locale/lang/en-us';
import zhCN from '@arco-design/web-vue/es/locale/lang/zh-cn';
import type { ArcoLang } from '@arco-design/web-vue/es/locale/interface'
import { useLocalStorage } from '@/assets/script/hooks'
import type { AppConfig, PageTheme } from '@/config/app-config'
import { changeTheme } from '@/assets/script/util'

const configStore = useAppConfigStore();

const { get } = useLocalStorage();

const localeMap = new Map<string, ArcoLang>([
  ['zh-CN', zhCN], ['en-US', enUS]
])

const language = computed(() => {
  return localeMap.get(configStore.state.language);
})

const initAppConfig = async () => {
  const localConfig = get<AppConfig>('app-config');
  localConfig && (Object.assign(configStore.state, localConfig));
  changeTheme(configStore.state.pageTheme);
}
const listenerColorScheme = () => {
  window.matchMedia('(prefers-color-scheme: dark)').addListener(e => {
    const newTheme: PageTheme = e.matches ? 'dark' : 'light';
    configStore.state.pageTheme = newTheme;
    changeTheme(newTheme);
  })
}

onMounted(() => {
  listenerColorScheme();
  initAppConfig();
})

</script>

<template>
  <a-config-provider :locale="language">
    <router-view />
  </a-config-provider>
</template>

<style scoped lang="scss">

</style>
