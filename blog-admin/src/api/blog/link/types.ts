// 友情链接
import type { PageQuery } from '@/assets/script/types'

export interface Link {
  linkId: number;
  name: string;
  summary: string;
  coverUrl: string;
  targetUrl: string;
  clickNum: number;
  sort: number;
  status: number;
  createTime: string;
}

// 友情链接表单
export interface LinkForm {
  linkId?: number;
  name: string;
  summary?: string;
  coverUrl: string;
  targetUrl: string;
  sort: number;
  status: number;
}

// 友情链接查询表单
export interface LinkQueryForm extends PageQuery {
  name?: string;
  createTimeBegin?: string;
  createTimeEnd?: string;
}