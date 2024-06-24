import type { OptionItem } from '@/assets/script/types'
import { computed } from 'vue'
import { useAppConfigStore } from '@/stores/app-config'

const appConfig = useAppConfigStore();

// 国际化hooks
export const useI18nLanguage = () => {
  // 支持的语言选项
  const languageOption: OptionItem[] = [
    { label: '中文', value: 'zh-CN' },
    { label: 'English', value: 'en-US'}
  ];
  // 当前使用的语言
  const currentLanguage = computed((): OptionItem => {
    const find = languageOption.find(item => item.value === appConfig.state.language);
    if (find) {
      return find;
    }
    appConfig.state.language = languageOption[0].value.toString();
    return languageOption[0];
  })

  // 更改语言
  function changeLanguage(value: string){
    appConfig.state.language = value;
  }

  return { languageOption, currentLanguage, changeLanguage };
}