import type { Topic, TopicForm, TopicQueryForm, TopicUpdateForm } from '@/api/blog/topic/types'
import { sendDelete, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

export const topicApi = {
  /**
   * 保存博客动态
   * @param form 动态参数
   */
  saveTopic: (form: TopicForm) => {
    return sendPost<null>('/topic', form)
  },
  /**
   * 更新博客动态
   * @param form 动态参数
   */
  updateTopic: (form: TopicForm) => {
    return sendPut<null>('/topic', form)
  },
  /**
   * 快捷更新博客动态
   * @param form 快捷更新参数
   */
  updateSelective: (form: TopicUpdateForm) => {
    return sendPut<null>('/topic/status', form)
  },
  /**
   * 分页查询博客动态
   * @param query 动态查询参数
   */
  pageTopic: (query: TopicQueryForm) => {
    return sendPost<Page<Topic>>('/topic/page', query)
  },
  /**
   * 删除博客动态
   * @param id 动态Id
   */
  deleteTopic: (id: number) => {
    return sendDelete<null>(`/topic/${id}`)
  }
}