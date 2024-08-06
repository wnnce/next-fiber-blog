// 博客标签
import type { PageQuery } from '@/assets/script/types'

export interface Tag {
  tagId: number;
  tagName: string;
  coverUrl: string;
  viewNum: number;
  color: string;
  sort: number;
  status: number;
  createTime: string;
  articleNum: number;
}

export interface ArticleTag {
  tagId: number;
  tagName: string;
  color: string;
}

// 博客标签表单
export interface TagForm {
  tagId?: number;
  tagName: string;
  coverUrl: string;
  color: string;
  sort: number;
  status: number;
}

// TagUpdateForm 博客标签快捷更新表单
export interface TagUpdateForm {
  tagId: number;
  status: number
}

// 博客标签查询表单
export interface TagQueryForm extends PageQuery{
  tagName?: string;
  createTimeBegin?: string;
  createTimeEnd?: string;
}