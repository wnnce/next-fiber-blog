import { sendGet } from '@/api/request'
import type { Menu } from '@/api/system/menu/types'

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
  }
}