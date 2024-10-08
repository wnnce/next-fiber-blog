import type { Role, RoleForm, RoleQueryForm, RoleUpdateForm } from '@/api/system/role/types'
import { sendDelete, sendGet, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

export const roleApi = {
  /**
   * 分页查询系统角色
   * @param query 查询参数
   */
  pageSysRole: (query: RoleQueryForm) => {
    return sendPost<Page<Role>>('/system/role/page', query)
  },
  /**
   * 获取所有角色列表
   */
  listAllSysROle: () => {
    return sendGet<Role[]>('/system/role/list');
  },
  /**
   * 保存系统角色
   * @param form
   */
  saveSysRole: (form: RoleForm) => {
    return sendPost<null>('/system/role', form);
  },
  /**
   * 更新系统角色
   * @param form
   */
  updateSysRole: (form: RoleForm) => {
    return sendPut<null>('/system/role', form);
  },
  /**
   * 快捷更新系统角色
   * @param form
   */
  updateSelective: (form: RoleUpdateForm) => {
    return sendPut<null>('/system/role/status', form)
  },
  /**
   * 删除系统角色
   * @param id 角色id
   */
  deleteSysRole: (id: number) => {
    return sendDelete<null>(`/system/role/${id}`);
  }
}