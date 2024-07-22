import type { PageQuery } from '@/assets/script/types'

export interface Notice {
  noticeId: number;
  title: string;
  message: string;
  level: number;
  noticeType: number;
  createTime: string;
  sort: number;
  status: number;
}

export interface NoticeForm {
  noticeId?: number;
  title?: string;
  message?: string;
  level?: number;
  noticeType?: number;
  sort?: number;
  status?: number;
}

export interface NoticeQueryForm extends PageQuery {
  title?: string;
  level?: number;
  noticeType?: number;
  createTimeBegin?: string;
  createTimeEnd?: string;
}