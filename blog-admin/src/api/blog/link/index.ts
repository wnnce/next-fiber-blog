import type { Link, LinkForm, LinkQueryForm } from '@/api/blog/link/types'
import { sendDelete, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

// 友情链接接口
export const linkApi = {
  /**
   * 分页查询友情链接
   * @param query 查询参数
   */
  pageLink: (query: LinkQueryForm) => {
    return sendPost<Page<Link>>('/link/page', query);
  },
  /**
   * 保存友情链接
   * @param form 友情链接参数
   */
  saveLink: (form: LinkForm) => {
    return sendPost<null>('/link', form);
  },
  /**
   * 更新友情链接
   * @param form 友情链接参数
   */
  updateLink: (form: LinkForm) => {
    return sendPut<null>('/link', form)
  },
  /**
   * 删除友情链接
   * @param linkId 友情链接Id
   */
  deleteLink: (linkId: number) => {
    return sendDelete<null>(`/link/${linkId}`);
  }
}