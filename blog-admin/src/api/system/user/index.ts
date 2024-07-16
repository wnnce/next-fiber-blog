import type { LoginForm, ResetPasswordForm, User, UserForm, UserQueryForm } from '@/api/system/user/types'
import { sendDelete, sendGet, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

export const userApi = {
  /**
   * 保存系统用户
   * @param form 用户参数
   */
  saveSysUser: (form: UserForm) => {
    return sendPost<null>('/system/user', form);
  },
  /**
   * 更新系统用户
   * @param form 用户参数
   */
  updateSysUser: (form: UserForm) => {
    return sendPut<null>('/system/user', form);
  },
  /**
   * 分页查询系统用户
   * @param query 查询参数
   */
  pageSysUser: (query: UserQueryForm) => {
    return sendPost<Page<User>>('/system/user/page', query);
  },
  /**
   * 获取用户详情
   */
  queryUserInfo: () => {
    return sendGet<User>('/system/user/info');
  },
  /**
   * 用户登录
   * @param form 登录参数
   */
  login: (form: LoginForm) => {
    return sendPost<string>('/system/user/login', form);
  },
  /**
   * 删除用户
   * @param id 用户id
   */
  deleteSysUser: (id: number) => {
    return sendDelete<null>(`/system/user/${id}`)
  },
  /**
   * 重置密码
   * @param form 重置密码参数
   */
  resetPassword: (form: ResetPasswordForm) => {
    return sendPut<null>('/system/user/password', form)
  },
  /**
   * 退出登录
   */
  logout: () => {
    return sendPost<null>('/system/user/logout')
  }
}