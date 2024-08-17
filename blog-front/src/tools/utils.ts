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

/**
 * 睡眠函数 强制等待指定时间
 * @param interval 需要等待的时间
 */
export const sleep = (interval: number): Promise<null> => {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(null);
    }, interval)
  })
}

/**
 * 节流函数 借助window.requestAnimationFrame实现节流
 * @param func 需要执行的函数
 * @param callback 事件被节流后的回调函数
 */
export const throttle = (func: (...args: any[]) => void, callback?: () => void): () => void => {
  let lock = false;
  return (...args: any[]) => {
    if (lock) {
      callback && callback();
      return;
    }
    lock = true;
    window.requestAnimationFrame(() => {
      func(...args);
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
export const throttleTimer = (func: (...args: any[]) => void, delay: number, callback?: () => void): () => void => {
  let timer: NodeJS.Timeout | number | undefined = undefined;
  return (...args: any[]) => {
    if (timer) {
      callback && callback()
      return;
    }
    func(...args);
    timer = setTimeout(() => {
      timer = undefined;
    }, delay)
  }
}