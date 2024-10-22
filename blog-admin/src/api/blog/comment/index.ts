import type { Comment, CommentQueryForm, CommentUpdateForm } from '@/api/blog/comment/types'
import { sendDelete, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

/**
 * 评论接口
 */
export const commentApi = {
  /**
   * 分页查询评论数据
   * @param query 查询参数
   */
  pageComment: (query: CommentQueryForm) => {
    return sendPost<Page<Comment>>('/comment/manage/page', query)
  },
  /**
   * 快捷更新评论
   * @param form 需要更新的参数
   */
  updateSelective: (form: CommentUpdateForm) => {
    return sendPut<null>('/comment/manage/status', form)
  },
  /**
   * 通过id删除评论
   * @param id 待删除的评论id
   */
  deleteComment: (id: number) => {
    return sendDelete<null>(`/comment/manage/${id}`)
  }
}