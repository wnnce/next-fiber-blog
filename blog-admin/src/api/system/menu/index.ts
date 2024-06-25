import { sendDelete, sendGet, sendPost, sendPut } from '@/api/request'
import type { Menu, MenuForm } from '@/api/system/menu/types'

/**
 * 系统菜单api接口
 */
export const menuApi = {
  /**
   * 获取菜单树形列表
   */
  listTreeMenu: () => {
    return sendGet<Menu[]>('/system/menu/tree')
  },
  /**
   * 获取菜单树管理列表
   */
  manageListTree: () => {
    return sendGet<Menu[]>('/system/menu/manage/tree')
  },
  /**
   * 保存系统菜单
   * @param menu 菜单参数
   */
  saveSysMenu: (menu: MenuForm) => {
    return sendPost<null>('/system/menu', menu)
  },
  /**
   * 更新系统菜单
   * @param menu 菜单参数
   */
  updateSysMenu: (menu: MenuForm) => {
    return sendPut<null>('/system/menu', menu)
  },
  /**
   * 删除系统菜单
   * @param id 菜单id
   */
  deleteSysMenu: (id: number) => {
    return sendDelete<null>(`/system/menu/${id}`)
  }
}