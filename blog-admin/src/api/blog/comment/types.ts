// 评论数据
import type { PageQuery } from '@/assets/script/types'

export interface Comment {
  commentId: number;
  content: string;
  userId: number;
  articleId?: number;
  topicId?: number;
  fid: number;
  rid: number;
  location: string;
  commentIp: string;
  commentUa: string;
  voteUp: number;
  voteDown: number;
  commentType: number;
  isHot: boolean;
  isTop: boolean;
  isColl: boolean;
  sort: number;
  createTime: string;
  updateTime: string;
  status: number;
  username: string;
  articleTile: string;
}

// CommentQueryForm 评论查询参数
export interface CommentQueryForm extends PageQuery {
  articleId?: number;
  topicId?: number;
  fid?: number;
  commentType?: number;
  createTimeBegin?: string;
  createTimeEnd?: string;
}

// CommentQueryForm 评论快捷更新表单
export interface CommentUpdateForm {
  commentId: number;
  isHot?: boolean;
  isTop?: boolean;
  isColl?: boolean;
  status?: number;
}