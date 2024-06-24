import type { Config, ConfigForm, ConfigQueryForm } from '@/api/system/config/types'
import { sendDelete, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

// 系统参数配置接口
export const configApi = {
  /**
   * 分页查询系统配置
   * @param query 查询参数
   */
  pageSysConfig: (query: ConfigQueryForm) => {
    return sendPost<Page<Config>>('/system/config/page', query)
  },
  /**
   * 保存系统配置
   * @param config 系统配置参数
   */
  saveSysConfig: (config: ConfigForm) => {
    return sendPost<null>('/system/config', config);
  },
  /**
   * 更新系统配置
   * @param config 系统配置参数
   */
  updateSysConfig: (config: ConfigForm) => {
    return sendPut<null>('/system/config', config);
  },
  /**
   * 删除系统配置
   * @param id 系统配置id
   */
  deleteSysConfig: (id: number) => {
    return sendDelete<null>(`/system/config/${id}`);
  }
}