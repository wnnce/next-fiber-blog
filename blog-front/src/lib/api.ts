import { request } from '@/lib/request'
import { Category, Result, Tag } from '@/lib/types'

/**
 * 获取博客分类列表
 */
export const listCategory = (): Promise<Result<Category[]>> => {
  return request<Category[]>('/open/category/list', 'GET');
}

/**
 * 获取博客标签列表
 */
export const listTag = (): Promise<Result<Tag[]>> => {
  return request<Tag[]>('/open/tag/list', 'GET');
}