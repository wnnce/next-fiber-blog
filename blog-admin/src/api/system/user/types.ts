import type { PageQuery } from '@/assets/script/types'

export interface User {
  userId: number;
  username: string;
  nickname: string;
  password: string
  email: string;
  phone: string;
  avatar: string;
  lastLoginIp: string;
  lastLoginTime: string;
  roles: number[];
  remark: string;
  sort: number;
  status: number;
  createTime: string;
}

// 用户表单
export interface UserForm {
  userId?: number;
  username: string;
  nickname?: string;
  password: string;
  email?: string;
  phone?: string;
  avatar: string;
  roles: number[];
  sort: number;
  status: number;
  remark?: string
}

// 用户快捷更新表单
export interface UserUpdateForm {
  userId: number;
  status: number;
}

// 用户查询表单
export interface UserQueryForm extends PageQuery {
  username?: string;
  nickname?: string;
  email?: string;
  phone?: string;
  roleId?: number;
  createTimeBegin?: string;
  createTimeEnd?: string;
}

// 登录表单
export interface LoginForm {
  username?: string;
  password?: string;
}

// 更新密码表单
export interface ResetPasswordForm {
  oldPassword: string;
  newPassword: string;
}