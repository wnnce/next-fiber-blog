import type { PageTheme } from '@/config/app-config'

const qiniuDomain = import.meta.env.VITE_QINIU_DOMAIN;

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

/**
 * 拼接七牛云图片地址
 * @param imageUrl 后端返回的七牛图片地址
 */
export const sliceImageUrl = (imageUrl: string) => {
  if (imageUrl.startsWith('/')) {
    return qiniuDomain + imageUrl.substring(1);
  }
  return qiniuDomain + imageUrl;
}