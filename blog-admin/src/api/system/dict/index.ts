import type {
  Dict,
  DictForm,
  DictQueryForm,
  DictValue,
  DictValueForm,
  DictValueQueryForm
} from '@/api/system/dict/types'
import { sendDelete, sendPost } from '@/api/request'
import type { Page } from '@/assets/script/types'

export const dictApi = {
  pageDict: (query: DictQueryForm) => {
    return sendPost<Page<Dict>>('/system/dict/page', query)
  },
  saveDict: (form: DictForm) => {
    return sendPost<null>('/system/dict', form)
  },
  updateDict: (form: DictForm) => {
    return sendPost<null>('/system/dict', form)
  },
  updateDictStatus: (form: DictForm) => {
    return sendPost<null>('/system/dict/status', form)
  },
  deleteDict: (dictId: number) => {
    return sendDelete<null>(`/system/dict/${dictId}`)
  },
  pageDictValue: (query: DictValueQueryForm) => {
    return sendPost<Page<DictValue>>('/system/dict/value/page', query)
  },
  saveDictValue: (form: DictValueForm) => {
    return sendPost<null>('/system/dict/value', form)
  },
  updateDictValue: (form: DictValueForm) => {
    return sendPost<null>('/system/dict/value', form)
  },
  updateDictValueStatus: (form: DictValueForm) => {
    return sendPost<null>('/system/dict/value/status', form)
  },
  deleteDictValue: (valueId: number) => {
    return sendDelete<null>(`/system/dict/value/${valueId}`)
  },
}