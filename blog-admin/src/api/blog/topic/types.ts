// 博客动态
import type { PageQuery } from '@/assets/script/types'

export interface Topic {
  topicId: number;
  content: string;
  imageUrls: string[];
  location: string;
  isHot: boolean;
  isTop: boolean;
  voteUp: number;
  mode: number;
  sort: number;
  status: number;
  createTime: string;
}

// 博客动态表单
export interface TopicForm {
  topicId?: number;
  content: string;
  imageUrls: string[];
  location?: string;
  isHot: boolean;
  isTop: boolean;
  mode?: number;
  sort?: number;
  status?: number;
}

// 博客动态快捷更新表单
export interface TopicUpdateForm {
  topicId: number
  isHot?: boolean;
  isTop?: boolean;
  status?: number;
}

// 博客动态查询表单
export interface TopicQueryForm extends PageQuery {
  location?: string;
  status?: number;
  createTimeBegin?: string;
  createTimeEnd?: string;
}