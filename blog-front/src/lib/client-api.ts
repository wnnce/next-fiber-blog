import { HttpMethod, request } from '@/lib/request'
import { Comment, Page, Result, SimpleArticle, User } from '@/lib/types'

export const clientAuthTokenKey = "client-auth-token";

function baseClientRequest<T>(url: string, method: HttpMethod, params?: Record<string, any>, body?: Record<string, any>, headers?: Record<string, string>): Promise<Result<T>> {
  const authToken = localStorage.getItem(clientAuthTokenKey) || '';
  if (!headers) {
    headers = {
      Authorization: `Bearer ${authToken}`
    }
  } else {
    headers['Authorization'] = authToken;
  }
  return request(url, method, params, body, headers);
}

/**
 * github登录接口
 * @param code github授权后返回的临时code
 */
export const loginWithGithub = (code: string) => {
  return request<string>('/open/classic/login/github', 'GET', { code: code });
}

/**
 * 获取用户详细信息
 */
export const userInfo = () => {
  return baseClientRequest<User>('/user/info', 'GET');
}

/**
 * 用户退出登录接口
 */
export const logout = () => {
  return baseClientRequest<null>('/user/logout', 'GET');
}

/**
 * 分页查询评论数据列表
 * @param data 查询参数
 */
export const pageComment = (data: Record<string, any>) => {
  return baseClientRequest<Page<Comment>>('/comment/page', 'POST', undefined, data);
}

/**
 * 保存评论
 * @param data 评论数据
 */
export const saveComment = (data: Record<string, any>) => {
  return baseClientRequest<null>('/comment', 'POST', undefined, data);
}

/**
 * 评论点赞接口
 * @param commentId 点赞的评论Id
 */
export const commentVoteUp = (commentId: number) => {
  return baseClientRequest<null>(`/comment/vote-up/${commentId}`, 'POST')
}

/**
 * 查询评论数量
 * @param data 查询参数
 */
export const totalComment = (data: Record<string, any>) => {
  return baseClientRequest<number>('/comment/total', 'POST', undefined, data);
}

/**
 * 动态点赞
 * @param topicId 点赞的动态id
 */
export const topicVoteUp = (topicId: number) => {
  return baseClientRequest<null>(`/open/topic/vote-up/${topicId}`, 'POST')
}

/**
 * 文章点赞
 * @param articleId 点赞的文章id
 */
export const articleVoteUp = (articleId: number) => {
  return baseClientRequest<null>(`/open/article/vote-up/${articleId}`, 'POST')
}

/**
 * 搜索文章
 * @param keyword 搜索关键字
 */
export const searchArticle = (keyword: string) => {
  return baseClientRequest<SimpleArticle[]>('/open/search/article', 'GET', { keyword })
}

/**
 * 记录访问请求
 */
export const traceAccess = () => {
  return baseClientRequest<null>('/open/trace/access', 'GET')
}