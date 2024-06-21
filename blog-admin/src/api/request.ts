import { useLocalStorage } from '@/hooks/local-storage'
import { useArcoMessage } from '@/hooks/message'

const { errorMessage } = useArcoMessage();

export interface Result<T> {
  code: number;
  message: string;
  timestamp: number;
  data: T
}

export declare type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';

export const TOKEN_KEY: string = 'Authorization_Bearer_Token';

const { get } = useLocalStorage();

const baseUrl = import.meta.env.VITE_REQUEST_BASE_URL;

export function request<T>(url: string, method: HttpMethod, params?: any, data?: object, headers?: Record<string, string>): Promise<Result<T>> {
  return new Promise((resolve, reject) => {
    const token = get<string>(TOKEN_KEY);
    const httpHeaders: Record<string, string> = {};
    if (token && token.trim().length > 0) {
      httpHeaders['Authorization'] = `Bearer ${token}`;
    }
    headers && (Object.assign(httpHeaders, headers));
    if (params) {
      url += '?' + new URLSearchParams(params).toString();
    }
    fetch(baseUrl + url, {
      method: method,
      headers: httpHeaders,
      body: data ? JSON.stringify(data) : undefined
    }).then(async response => {
      const responseBody = await response.json() as Result<T>;
      const { code, message } = responseBody;
      if (code === 200) {
        resolve(responseBody);
        return;
      } else if (code === 401) {
        errorMessage('当前登录状态已失效，请先登录!');
      } else if (code === 403) {
        errorMessage('当前操作无权限');
      } else if (code === 400) {
        errorMessage(message);
      } else if (code === 500) {
        errorMessage('服务器错误，请稍后再试');
      } else {
        errorMessage('未知错误，请联系管理员');
      }
      reject(undefined);
    }).catch(err => {
      console.log(err)
      errorMessage('请求失败')
      reject(undefined);
    })
  })
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
export function fileUpload(url: string, formData: FormData, onProgress?: (event: ProgressEvent) => {}) {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    if (onProgress) {
      xhr.addEventListener('progress', onProgress);
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
    const token = get<string>(TOKEN_KEY);
    xhr.setRequestHeader('Authorization', `Bearer ${token}`);
    xhr.open('POST', baseUrl + url);
    xhr.send(formData);
  })
}