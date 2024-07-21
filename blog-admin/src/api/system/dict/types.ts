import type { PageQuery } from '@/assets/script/types'

export interface Dict {
  dictId: number;
  dictName: string;
  dictKey: string;
  remark?: string;
  sort: number;
  status: number;
  createTime: string
}

export interface DictForm {
  dictId?: number;
  dictName?: string;
  dictKey?: string;
  remark?: string;
  sort?: number;
  status?: number;
}

export interface DictQueryForm extends PageQuery {
  dictName?: string;
  dictKey?: string;
  createTimeBegin?: string;
  createTimeEnd?: string;
}

export interface DictValue {
  id: number;
  dictId: number;
  dictKey: string;
  label: string;
  value: string;
  remark?: string;
  sort: number;
  status: number;
  createTime: string
}

export interface DictValueForm {
  id?: number;
  dictId?: number;
  dictKey?: string;
  label?: string;
  value?: string;
  remark?: string;
  sort?: number;
  status?: number;
}

export interface DictValueQueryForm extends PageQuery {
  dictId?: number;
  dictKey?: string;
  label?: string;
  createTimeBegin?: string;
  createTimeEnd?: string;
}