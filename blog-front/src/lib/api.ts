import { request } from '@/lib/request'
import { Category, Concat, Notice, Result, Tag } from '@/lib/types'
import { SiteConfigurationItem } from '@/hooks/site-configuration'

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

/**
 * 配置博客站点配置
 */
export const querySiteConfiguration = (): Promise<Result<Record<string, SiteConfigurationItem>>> => {
  return request<Record<string, SiteConfigurationItem>>('/open/site/configuration', 'GET')
}

/**
 * 获取联系方式列表
 */
export const listConcat = (): Promise<Result<Concat[]>> => {
  return request<Concat[]>('/open/concat/list', 'GET');
}

/**
 * 通过通知类型查询通知公告列表
 * @param noticeType 1: 首页弹窗通知 2:公告板通知
 */
export const listNoticeByType = (noticeType: 1 | 2): Promise<Result<Notice[]>> => {
  const requestUrl = noticeType === 1 ? '/open/notice/index' : '/open/notice/public';
  return request<Notice[]>(requestUrl, 'GET');
}