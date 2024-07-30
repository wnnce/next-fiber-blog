
// 博客文章
import type { PageQuery } from '@/assets/script/types'

export interface Article {
  articleId: number;
  title: string;
  summary: string;
  coverUrl: string;
  categoryIds: number[];
  tagIds: number[];
  viewNum: number;
  shareNum: number;
  voteUp: number;
  content: string;
  protocol: string;
  tips: string;
  password: string;
  isHot: boolean;
  isTop: boolean;
  isComment: boolean;
  isPrivate: boolean;
  createTime: string;
  sort: number;
  status: number;
}

// 博客表单
export interface ArticleForm {
  articleId?: number;
  title?: string;
  summary?: string;
  coverUrl: string;
  categoryIds: number[];
  tagIds: number[];
  content?: string;
  protocol?: string;
  tips?: string;
  password?: string;
  isHot: boolean;
  isTop: boolean;
  isComment: boolean;
  isPrivate: boolean;
  sort?: number;
  status?: number;
}

// 博客文章快捷更新表单
export interface ArticleUpdateFrom {
  articleId: number;
  isTop?: boolean;
  isHot?: boolean;
  isComment?: boolean;
  status?: number;
}

// 博客文章查询表单
export interface ArticleQueryForm extends PageQuery {
  title?: string;
  categoryId?: number;
  tagId?: number;
  status?: number;
  createTimeBegin?: string;
  createTimeEnd?: string;
}