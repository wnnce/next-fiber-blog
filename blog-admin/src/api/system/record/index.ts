import type {
  AccessRecord,
  AccessRecordQueryForm,
  ApplicationMonitor,
  LoginRecord,
  LoginRecordQueryForm
} from '@/api/system/record/types'
import { sendGet, sendPost } from '@/api/request'
import type { Page } from '@/assets/script/types'

// 日志/记录 接口
export const recordApi = {
  /**
   * 分页查询登录记录
   * @param query 查询参数
   */
  pageLoginRecord: (query: LoginRecordQueryForm) => {
    return sendPost<Page<LoginRecord>>('/system/record/login', query)
  },
  /**
   * 分页查询访问记录
   * @param query 查询参数
   */
  pageAccessRecord: (query: AccessRecordQueryForm) => {
    return sendPost<Page<AccessRecord>>('/system/record/access', query)
  },
  /**
   * 获取指标监控数据
   */
  applicationMonitor: () => {
    return sendGet<ApplicationMonitor>('/base/monitor');
  }
}