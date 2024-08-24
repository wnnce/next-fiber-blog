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

/**
 * 拼接缩略图图片地址 （仅适用于七牛云）
 * @param imageUrl 后端返回的七牛图片地址
 * @param h 图片短边的长度 长边自适应
 */
export const sliceThumbnailImageUrl = (imageUrl: string, h: number = 100) => {
  if (imageUrl.startsWith('/b-oss/')) {
    return imageUrl + `?imageView2/0/h/${h}`
  }
  return imageUrl;
}


export const parseWordCount = (count: number): number | string => {
  if (count < 1000) {
    return count
  }
  return `${(count / 1000).toFixed(1)}k`
}

/**
 * 节流函数 借助window.requestAnimationFrame实现节流
 * @param func 需要执行的函数
 * @param callback 事件被节流后的回调函数
 */
export const throttle = (func: () => void, callback?: () => void): () => void => {
  let lock = false;
  return (...args) => {
    if (lock) {
      callback && callback();
      return;
    }
    lock = true;
    window.requestAnimationFrame(() => {
      func.apply(this, args)
      lock = false;
    })
  }
}

/**
 * 节流函数 基于间隔事件进行节流
 * @param func 需要执行的函数
 * @param delay 间隔事件
 * @param callback 事件被节流后的回调函数
 */
export const throttleTimer = (func: () => void, delay: number, callback?: () => void): () => void => {
  let timer: number | undefined = undefined;
  return () => {
    if (timer) {
      callback && callback()
      return;
    }
    timer = setTimeout(() => {
      func()
      timer = undefined;
    }, delay)
  }
}