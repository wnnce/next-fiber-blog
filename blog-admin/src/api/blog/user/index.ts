import type { ExpertiseDetail, ExpertiseQueryForm, User, UserQueryForm, UserUpdateForm } from '@/api/blog/user/types'
import { sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

/**
 * 博客端用户接口
 */
export const userApi = {
  /**
   * 分页查询用户信息
   * @param query 查询参数
   */
  pageUser: (query: UserQueryForm) => {
    return sendPost<Page<User>>('/user/page', query)
  },
  /**
   * 更新用户信息
   * @param form
   */
  updateUser: (form: UserUpdateForm) => {
    return sendPut<null>('/user', form)
  },
  /**
   * 分页查询用户经验值明细
   * @param query 查询参数
   */
  pageExpertise: (query: ExpertiseQueryForm) => {
    return sendPost<Page<ExpertiseDetail>>('/user/expertise/page', query)
  }
}