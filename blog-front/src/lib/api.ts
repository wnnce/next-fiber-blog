import { request } from '@/lib/request'
import {
  Article,
  ArticleArchive,
  Category,
  Concat, HotArticle,
  Link,
  Notice,
  Page,
  Result,
  SiteStats,
  Tag,
  Topic
} from '@/lib/types'
import { SiteConfigurationItem } from '@/hooks/site-configuration'

/**
 * 获取博客分类列表
 */
export const listCategory = (): Promise<Result<Category[]>> => {
  return request<Category[]>('/open/category/list', 'GET');
}
/**
 * 查询博客分类详情
 * @param id 分类Id
 */
export const queryCategory = (id: number): Promise<Result<Category>> => {
  return request<Category>(`/open/category/${id}`, 'GET')
}

/**
 * 获取博客标签列表
 */
export const listTag = (): Promise<Result<Tag[]>> => {
  return request<Tag[]>('/open/tag/list', 'GET');
}
/**
 * 查询博客标签详情
 * @param id 标签id
 */
export const queryTag = (id: number): Promise<Result<Tag>> => {
  return request<Tag>(`/open/tag/${id}`, 'GET');
}
/**
 * 配置博客站点配置
 */
export const querySiteConfiguration = (): Promise<Result<Record<string, SiteConfigurationItem>>> => {
  return request<Record<string, SiteConfigurationItem>>('/open/site/configuration', 'GET')
}

/**
 * 获取站点统计数据
 */
export const querySiteStats = (): Promise<Result<SiteStats>> => {
  return request<SiteStats>('/open/site/stats', 'GET')
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

/**
 * 分页查询文章列表
 * @param data 查询参数
 */
export const pageArticle = (data: Record<string, any>) => {
  return request<Page<Article>>('/open/article/page', 'POST', undefined, data);
}

/**
 * 查询分类 / 标签所关联的文章列表
 * @param data
 */
export const pageLabelArticle = (data: Record<string, any>) => {
  return request<Page<Article>>('/open/article/label/page', 'POST', undefined, data);
}

/**
 * 查询置顶的文章列表
 */
export const listTopArticle = () => {
  return request<Article[]>('/open/article/top', 'GET')
}

/**
 * 查询热门文章列表
 */
export const listHotArticle = () => {
  return request<HotArticle[]>('/open/article/hot', 'GET')
}

/**
 * 获取博客文章的归档信息
 */
export const articleArchives = () => {
  return request<ArticleArchive[]>('/open/article/archives', 'GET')
}

/**
 * 获取博客文章详情
 * @param id 文章id
 */
export const queryArticle = (id: number) => {
  return request<Article>(`/open/article/${id}`, 'GET');
}

/**
 * 获取友情链接列表
 */
export const listLinks = (): Promise<Result<Link[]>> => {
  return request<Link[]>('/open/link/list', 'GET');
}

/**
 * 获取动态列表
 * @param data 查询条件
 */
export const pageTopic = (data: Record<string, any>) => {
  return request<Page<Topic>>('/open/topic/page', 'POST', undefined, data);
}