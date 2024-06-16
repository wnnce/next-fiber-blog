import type { PageTheme } from '@/config/app-config'

/**
 * 切换白天 黑暗模式
 * @param theme 选择的主题模式
 */
export const changeTheme = (theme: PageTheme) => {
  if (theme === 'light') {
    document.body.classList.remove('dark');
    document.body.classList.add('light');
    document.body.removeAttribute('arco-theme');
  } else {
    document.body.classList.remove('light');
    document.body.classList.add('dark');
    document.body.setAttribute('arco-theme', 'dark');
  }
}