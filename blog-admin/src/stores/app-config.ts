import { defineStore } from 'pinia'
import { computed, reactive, watch } from 'vue'
import { type AppConfig, defaultAppConfig } from '@/config/app-config'
import { useLocalStorage } from '@/hooks/local-storage'
const { set } = useLocalStorage();

export const useAppConfigStore = defineStore('app-config', () => {
  const state = reactive<AppConfig>({ ...defaultAppConfig })

  const sideWidth = computed(() => {
    return typeof state.sideWidth === 'number' ? `${state.sideWidth}px` : state.sideWidth;
  })

  watch(state, (newValue) => {
    set<AppConfig>('app-config', newValue, undefined);
  })

  return { state, sideWidth }
})