// 系统配置
import type { PageQuery } from '@/assets/script/types'

export interface Config {
  configId: number;
  configName: string;
  configKey: string;
  configValue: string;
  createTime: string;
  updateTime: string;
  remark: string;
}

// 系统配置表单
export interface ConfigForm {
  configId?: number;
  configName: string;
  configKey: string;
  configValue: string;
  remark: string;
}

// 系统配置查询表单
export interface ConfigQueryForm extends PageQuery {
  name?: string;
  key?: string;
  createTimeBegin?: string;
  createTimeEnd?: string;
}