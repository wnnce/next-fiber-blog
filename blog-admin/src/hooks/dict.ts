import type { DictValue } from '@/api/system/dict/types'
import { dictApi } from '@/api/system/dict'

interface CacheItem {
  expireTime: number;
  value: DictValue[];
}

/**
 * 数据字典Hook 提供统一的查询字典数据方法 同时使用两级缓存来提升性能
 */
export const useDict = () => {
  // 数据字典请求缓存
  const promiseCache = new Map<string, Promise<DictValue[]>>();
  // 数据字典缓存
  const dictCache = new Map<string, CacheItem>();

  /**
   * 通过字典key查询字典数据列表
   * 现在数据字典缓存中查询 如果没有则发送请求向后端查询
   * 在请求过程中如果还有新的当前key查询 则会使用同一个Promise
   * @param dictKey 数据字典key
   */
  async function queryDict(dictKey: string): Promise<DictValue[]> {
    const cacheItem = dictCache.get(dictKey)
    if (cacheItem && cacheItem.expireTime > new Date().getTime()) {
      return cacheItem.value;
    } else if (cacheItem) {
      dictCache.delete(dictKey);
    }
    const dictPromise = promiseCache.get(dictKey);
    if (dictPromise) {
      return dictPromise;
    }
    const promise: Promise<DictValue[]> = new Promise(resolve => {
      dictApi.listDictValueByKey(dictKey).then(res => {
        const { code, data }  = res;
        if (code === 200 && data.length > 0) {
          // 添加字典缓存
          dictCache.set(dictKey, { expireTime: new Date().getTime() + 30000, value: data })
          resolve(data);
          return;
        }
        resolve([])
      }).catch(err => {
        console.log(err)
        resolve([])
      }).finally(() => {
        promiseCache.delete(dictKey);
      })
    })
    // 添加promise缓存
    promiseCache.set(dictKey, promise);
    return promise;
  }

  /**
   * 清除promise缓存和数据字典缓存
   */
  function clear() {
    promiseCache.clear();
    dictCache.clear();
  }

  return { queryDict, clear }
}