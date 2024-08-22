import { Result } from '@/lib/types'

export declare type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';

const baseUrl = process.env.NEXT_PUBLIC_BASE_URL

export function request<T> (url: string, method: HttpMethod, params?: Record<string, any>, body?: Record<string, any>, headers?: Record<string, string>): Promise<Result<T>> {
  return new Promise((resolve, reject) => {
    if (params) {
      const requestParams = new URLSearchParams(params).toString();
      url += `?${requestParams}`;
    }
    fetch(baseUrl + url, {
      method: method,
      headers: headers,
      body: body ? JSON.stringify(body) : undefined,
    }).then(async response => {
      const result = await response.json() as Result<T>;
      resolve(result);
    }).catch(err => {
      reject(err)
    })
  })
}