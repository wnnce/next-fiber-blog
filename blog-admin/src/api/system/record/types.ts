
// 登录日志
import type { PageQuery } from '@/assets/script/types'

export interface LoginRecord {
  id: number;
  userId: number;
  userType: number;
  username: string;
  loginIp: string;
  location: string;
  loginUa: string;
  createTime: string;
  remark: string;
  result: number;
  loginType: number;
}

// 登录日志查询表单
export interface LoginRecordQueryForm extends PageQuery {
  username?: string;
  loginType?: number;
  result?: number;
  createTimeBegin?: string;
  createTimeEnd?: string;
}

// 访问记录
export interface AccessRecord {
  id: number;
  location: string;
  referee: string;
  accessIp: string;
  accessUa: string;
  createTime: string;
}

// 访问记录查询表单
export interface AccessRecordQueryForm extends PageQuery {
  ip?: string;
  createTimeBegin?: string;
  createTimeEnd?: string;
}
