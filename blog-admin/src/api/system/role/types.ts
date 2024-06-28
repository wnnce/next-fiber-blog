// 系统角色
import type { PageQuery } from '@/assets/script/types'

export interface Role {
  roleId: number;
  roleName: string;
  roleKey: string;
  sort: number;
  status: number;
  createTime: string;
  remark: string;
  menus: number[];
}

// 系统角色表单
export interface RoleForm {
  roleId?: number;
  roleName: string;
  roleKey: string;
  sort: number;
  status: number;
  remark: string;
  menus: number[];
}

// 系统角色查询参数
export interface RoleQueryForm extends PageQuery{
  name?: string;
  key?: string;
  createTimeBegin?: string;
  createTimeEnd?: string;
}