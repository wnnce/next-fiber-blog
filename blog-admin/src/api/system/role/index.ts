import type { Role, RoleForm, RoleQueryForm } from '@/api/system/role/types'
import { sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

export const roleApi = {
  /**
   * 分页查询系统角色
   * @param query 查询参数
   */
  pageSysRole: (query: RoleQueryForm) => {
    return sendPost<Page<Role>>('/system/role/page', query)
  },
  saveSysRole: (form: RoleForm) => {
    return sendPost<null>('/system/role', form);
  },
  updateSysRole: (form: RoleForm) => {
    return sendPut<null>('/system/role', form);
  }
}