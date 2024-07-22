import type { Notice, NoticeForm, NoticeQueryForm } from '@/api/system/notice/types'
import { sendDelete, sendGet, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

// 通知公告接口
export const noticeApi = {
  /**
   * 保存通知
   * @param form 通知表单参数
   */
  saveNotice: (form: NoticeForm) => {
    return sendPost<null>('/system/notice', form);
  },
  /**
   * 更新通知
   * @param form
   */
  updateNotice: (form: NoticeForm) => {
    return sendPut<null>('/system/notice', form)
  },
  /**
   * 分页查询通知
   * @param query
   */
  pageNotice: (query: NoticeQueryForm) => {
    return sendPost<Page<Notice>>('/system/notice/page', query)
  },
  /**
   * 删除通知
   * @param id
   */
  deleteNotice: (id: number) => {
    return sendDelete<null>(`/system/notice/${id}`)
  },
  /**
   * 获取管理端通知列表
   */
  listAdminNotice: () => {
    return sendGet<Notice[]>('/base/notice/admin')
  }
}