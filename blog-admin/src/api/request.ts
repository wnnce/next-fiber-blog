import { useLocalStorage } from '@/hooks/local-storage'
import { useArcoMessage } from '@/hooks/message'
import { LOCAl_USER_KEY, TOKEN_KEY } from '@/assets/script/constant'
import { useLocalUserStore } from '@/stores/user'
import router from '@/router/index';

const { errorMessage } = useArcoMessage();

export interface Result<T> {
  code: number;
  message: string;
  timestamp: number;
  data: T
}

export declare type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';

const { get, remove } = useLocalStorage();

const baseUrl = import.meta.env.VITE_REQUEST_BASE_URL;

// 请求Promise缓存
const fetchCache = new Map<string, Promise<Result<any>>>;

export function request<T>(url: string, method: HttpMethod, params?: any, data?: object, headers?: Record<string, string>): Promise<Result<T>> {
  const token = get<string>(TOKEN_KEY);
  // 计算请求缓存的key
  const keyItems: string[] = [url, method.toString()];
  const httpHeaders: Record<string, string> = {};
  if (token && token.trim().length > 0) {
    httpHeaders['Authorization'] = `Bearer ${token}`;
  }
  headers && (Object.assign(httpHeaders, headers));
  if (params) {
    const queryParams = new URLSearchParams(params).toString()
    keyItems.push(`queryLength:${queryParams.length}`)
    url += '?' + queryParams;
  }
  const requestBody = data ? JSON.stringify(data) : undefined
  requestBody && (keyItems.push(`bodyLength:${requestBody.length}`))
  const promiseKey = keyItems.join('');
  // 判断请求是否在处理中 如果正在处理则直接返回缓存的Promise
  const cachePromise = fetchCache.get(promiseKey);
  if (cachePromise) {
    return cachePromise
  }
  const fetchPromise: Promise<Result<T>> = new Promise((resolve, reject) => {
    fetch(baseUrl + url, {
      method: method,
      headers: httpHeaders,
      body: requestBody
    }).then(async response => {
      const responseBody = await response.json() as Result<T>;
      const { code, message } = responseBody;
      if (code === 200) {
        resolve(responseBody);
        return;
      } else if (code === 401) {
        errorMessage('当前登录状态已失效，请重新登录!');
        // 清除登录用户信息 菜单信息 路由信息
        useLocalUserStore().clear();
        remove(TOKEN_KEY, LOCAl_USER_KEY);
        // 跳转到登录页
        router.push({ path: '/login' })
      } else if (code === 403) {
        errorMessage('当前操作无权限');
      } else if (code === 400) {
        errorMessage(message);
      } else if (code === 500) {
        errorMessage(import.meta.env.DEV ? message : '服务器错误，请稍后再试');
      } else {
        errorMessage('未知错误，请联系管理员');
      }
      reject(undefined);
    }).catch(err => {
      console.log(err)
      reject(undefined);
    }).finally(() => {
      // 请求结束后清理缓存
      fetchCache.delete(promiseKey)
    })
  })
  fetchCache.set(promiseKey, fetchPromise)
  return fetchPromise;
}

/**
 * 发送普通get请求
 * @param url 请求地址
 * @param params 请求url参数
 */
export function sendGet<T>(url: string, params?: any) {
  return request<T>(url, 'GET', params)
}

/**
 * 发送json数据get请求
 * @param url 请求地址
 * @param params 请求url参数
 * @param data 请求body json参数
 */
export function sendGetByPost<T>(url: string, params?: any, data?: object) {
  return request<T>(url, 'GET', params, data, { 'content-type': 'application/json' });
}

/**
 * 发送post请求
 * @param url 请求地址
 * @param data 请求json数据
 */
export function sendPost<T>(url: string, data?: object) {
  return request<T>(url, 'POST', undefined, data, { 'content-type': 'application/json' });
}

/**
 * 发送put请求
 * @param url 请求地址
 * @param data 请求json数据
 */
export function sendPut<T>(url: string, data?: object) {
  return request<T>(url, 'PUT', undefined, data, { 'content-type': 'application/json' });
}

/**
 * 发送delete请求
 * @param url 请求地址
 * @param params 请求url参数
 */
export function sendDelete<T>(url: string, params?: object) {
  return request<T>(url, 'DELETE', params);
}

/**
 * 文件上传请求
 * @param url 请求地址
 * @param formData 包含待上传文件的form表单
 * @param onProgress 文件上传进度回调
 */
export function fileUpload(url: string, formData: FormData, onProgress?: (event: ProgressEvent) => void): Promise<Result<string>> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    if (onProgress) {
      xhr.upload.onprogress = onProgress
    }
    xhr.addEventListener('error', err => {
      console.log(err);
      errorMessage('文件上传失败');
      reject(err);
    })
    xhr.addEventListener('load', () => {
      const { status, responseText } = xhr;
      const result = JSON.parse(responseText) as Result<string>
      if (status != 200) {
        errorMessage(result.message)
        reject(result)
        return;
      } else {
        resolve(result);
      }
    })
    xhr.open('POST', baseUrl + url);
    // 添加验证Token
    const token = get<string>(TOKEN_KEY);
    xhr.setRequestHeader('Authorization', `Bearer ${token}`);
    xhr.send(formData);
  })
}