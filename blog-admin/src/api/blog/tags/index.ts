import type { Tag, TagForm, TagQueryForm } from '@/api/blog/tags/types'
import { sendDelete, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

// 系统标签API
export const tagApi = {
  /**
   * 分页查询博客标签
   * @param query 查询参数
   */
  pageTag: (query: TagQueryForm) => {
    return sendPost<Page<Tag>>('/tag/page', query);
  },
  /**
   * 保存博客标签
   * @param form 表单参数
   */
  saveTag: (form: TagForm) => {
    return sendPost<null>('/tag', form);
  },
  /**
   * 更新博客标签
   * @param form 表单参数
   */
  updateTag: (form: TagForm) => {
    return sendPut<null>('/tag', form);
  },
  /**
   * 删除博客标签
   * @param id 标签Id
   */
  deleteTag: (id: number) => {
    return sendDelete<null>(`/tag/${id}`)
  }
}