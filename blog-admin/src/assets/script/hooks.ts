import type { PageTheme } from '@/config/app-config'

interface StorageItem<T> {
  key: string;
  value: T;
  expireTime: number | undefined;
  lasting: boolean
}

export const useLocalStorage = () => {
  // 获取缓存对象
  function get<T>(key: string): T | undefined {
    const stringItem = window.localStorage.getItem(key)
    if (!stringItem) {
      return void 0;
    }
    const item = JSON.parse(stringItem) as StorageItem<T>;
    if (item.lasting) {
      return item.value;
    }
    // 判断是否过期
    if (item.expireTime && Date.now() >= item.expireTime) {
      setTimeout(() => {
        remove(key);
      }, 0)
      return void 0;
    }
    return item.value;
  }

  // 添加缓存对象
  function set<T>(key: string, value: T, expire: number | undefined) {
    const item: StorageItem<T> = {
      key: key,
      value: value,
      expireTime: expire ? Date.now() + expire : undefined,
      lasting: !expire
    }
    const stringValue = JSON.stringify(item);
    window.localStorage.setItem(key, stringValue);
  }

  // 删除缓存对象
  function remove(key: string) {
    window.localStorage.removeItem(key);
  }

  return { get, set, remove }
}