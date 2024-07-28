import type { Article, ArticleForm, ArticleQueryForm, ArticleUpdateFrom } from '@/api/blog/article/types'
import { sendDelete, sendGet, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

export const articleApi = {
  /**
   * 保存博客文章
   * @param form 文章参数
   */
  saveArticle: (form: ArticleForm) => {
    return sendPost<null>('/article', form)
  },
  /**
   * 更新博客文章
   * @param form 文章参数
   */
  updateArticle: (form: ArticleForm) => {
    return sendPut<null>('/article', form)
  },
  /**
   * 快捷更新博客文章
   * @param form 表单参数
   */
  updateSelective: (form: ArticleUpdateFrom) => {
    return sendPut<null>('/article/status', form)
  },
  /**
   * 查询博客文章详情
   * @param articleId 文章id
   */
  queryArticleInfo: (articleId: number) => {
    return sendGet<Article>(`/article/info/${articleId}`)
  },
  /**
   * 分页查询文章列表
   * @param query 查询参数
   */
  pageArticle: (query: ArticleQueryForm) => {
    return sendPost<Page<Article>>('/article/manage/page', query)
  },
  /**
   * 删除博客文章
   * @param articleId 文章id
   */
  deleteArticle: (articleId: number) => {
    return sendDelete<null>(`/article/${articleId}`)
  }

}