// 博客端用户信息
import type { PageQuery } from '@/assets/script/types'

export interface User {
  userId: number;
  avatar: string;
  nickname: string;
  summary: string;
  email: string;
  link: string;
  username: string;
  userType: number;
  labels: string[];
  createTime: string;
  status: number;
  level: number;
  expertise: number;
  registerIp: string;
  registerLocation: string
}

// 博客端用户更新表单
export interface UserUpdateForm {
  userId: number;
  nickname?: string;
  summary?: string;
  email?: string;
  link?: string;
  labels?: string[];
  status?: number;
}

// 博客用户查询参数
export interface UserQueryForm extends PageQuery {
  nickname?: string;
  email?: string;
  username?: string;
  level?: number;
  createTimeBegin?: string;
  createTimeEnd?: string;
}

// 经验值明细
export interface ExpertiseDetail {
  id: number;
  userId: number;
  username: string;
  nickname: string;
  detail: number;
  detailType: number;
  source: number;
  createTime: string;
  remark: string;
}

// 经验值明细查询参数
export interface ExpertiseQueryForm extends PageQuery {
  username?: string;
  detailType?: number;
  source?: number;
  createTimeBegin?: string;
  createTimeEnd?: string;
}