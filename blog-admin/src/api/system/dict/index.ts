import type {
  Dict,
  DictForm,
  DictQueryForm,
  DictValue,
  DictValueForm,
  DictValueQueryForm
} from '@/api/system/dict/types'
import { sendDelete, sendGet, sendPost, sendPut } from '@/api/request'
import type { Page } from '@/assets/script/types'

export const dictApi = {
  /**
   * 分页查询系统字典
   * @param query 查询参数
   */
  pageDict: (query: DictQueryForm) => {
    return sendPost<Page<Dict>>('/system/dict/page', query)
  },
  /**
   * 保存系统字典
   * @param form 字典参数
   */
  saveDict: (form: DictForm) => {
    return sendPost<null>('/system/dict', form)
  },
  /**
   * 更新系统字典
   * @param form 字典参数
   */
  updateDict: (form: DictForm) => {
    return sendPut<null>('/system/dict', form)
  },
  /**
   * 更新系统字典状态
   * @param form 状态参数
   */
  updateDictStatus: (form: DictForm) => {
    return sendPut<null>('/system/dict/status', form)
  },
  /**
   * 删除系统字典
   * @param dictId 字典Id
   */
  deleteDict: (dictId: number) => {
    return sendDelete<null>(`/system/dict/${dictId}`)
  },
  /**
   * 分页查询字典数据
   * @param query 查询参数
   */
  pageDictValue: (query: DictValueQueryForm) => {
    return sendPost<Page<DictValue>>('/system/dict/value/page', query)
  },
  /**
   * 保存字典数据
   * @param form 字典数据参数
   */
  saveDictValue: (form: DictValueForm) => {
    return sendPost<null>('/system/dict/value', form)
  },
  /**
   * 更新字典数据
   * @param form 字典数据参数
   */
  updateDictValue: (form: DictValueForm) => {
    return sendPut<null>('/system/dict/value', form)
  },
  /**
   * 更新字典数据状态
   * @param form 状态参数
   */
  updateDictValueStatus: (form: DictValueForm) => {
    return sendPut<null>('/system/dict/value/status', form)
  },
  /**
   * 删除字典数据
   * @param valueId 字典数据Id
   */
  deleteDictValue: (valueId: number) => {
    return sendDelete<null>(`/system/dict/value/${valueId}`)
  },
  /**
   * 通过字典Key查询字典数据列表
   * @param dictKey 字典数据key
   */
  listDictValueByKey: (dictKey: string) => {
    return sendGet<DictValue[]>(`/open/dict/${dictKey}`)
  }
}