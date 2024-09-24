import parse from 'ua-parser-js';
import UAParser from 'ua-parser-js'

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

export const formatWordCount = (num: number): number | string => {
  if (num < 1000) {
    return num
  }
  return `${(num / 1000).toFixed(1)}k`
}

/**
 * 格式化时间为统一格式 一分钟内返回刚刚 一小时内返回xx分钟前 一天内返回xx小时前 否则返回yyyy-MM-dd格式时间
 * @param datetime 需要被格式化的日期时间字符串
 */
export const formatDateTime = (datetime: string): string => {
  const now = new Date();
  if (datetime.indexOf('T') > -1) {
    datetime = datetime.replace('T', ' ');
  }
  if (datetime.indexOf('Z') > -1) {
    datetime = datetime.replace('Z', '')
  }
  const date = new Date(datetime);
  const diffMs = now.getTime() - date.getTime();
  const diffMinutes = Math.floor(diffMs / 60000);
  if (diffMinutes <= 1) {
    return '刚刚';
  }
  if (diffMinutes < 60) {
    return `${diffMinutes} 分钟前`;
  }
  const diffHours = Math.floor(diffMinutes / 60);
  if (diffHours < 24) {
    return `${diffHours} 小时前`
  }
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDay() + 1).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

/**
 * 解析UserAgent头
 * @param ua 待解析的ua
 */
export const formatUa = (ua: string): UAParser.IResult => {
  return parse(ua);
}