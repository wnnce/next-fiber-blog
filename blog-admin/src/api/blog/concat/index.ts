import type { Concat, ConcatForm, ConcatQueryForm } from '@/api/blog/concat/types'
import { sendDelete, sendPost, sendPut } from '@/api/request'

// 联系方式api
export const concatApi = {
  /**
   * 管理端查询所有联系方式
   * @param query 查询参数
   */
  manageList: (query: ConcatQueryForm) => {
    return sendPost<Concat[]>('/concat/manage/list', query);
  },
  /**
   * 保存联系方式
   * @param form 表单参数
   */
  saveConcat: (form: ConcatForm) => {
    return sendPost<null>('/concat', form);
  },
  /**
   * 更新联系方式
   * @param form 表单参数
   */
  updateConcat: (form: ConcatForm) => {
    return sendPut<null>('/concat', form);
  },
  /**
   * 删除联系方式
   * @param concatId 联系方式Id
   */
  deleteConcat: (concatId: number) => {
    return sendDelete<null>(`/concat/${concatId}`)
  }
}