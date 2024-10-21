
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

export const applicationMonitorRequiredField: string[] = [
  'hostname', 'platform', 'platformVersion', 'cpuNumber', 'cpuPercent', 'memoryTotal', 'memoryUsed',
  'memoryAvailable', 'memoryUsedPercent', 'Sys', 'HeapSys', 'HeapInuse', 'HeapIdle', 'StackSys',
  'StackInuse', 'PauseTotalNs', 'NumGC', 'GCCPUFraction'
]

export const applicationMonitorRequiredFieldName: string[] = [
  '主机名称', '平台', '平台版本', 'CPU核心数', 'CPU使用率', '总内存', '已用内存', '可用内存', '内存使用率',
  '程序占用内存', '堆内存', '堆已用内存', '堆可用内存', '栈内存', '栈已用内存', '暂停时间', '垃圾回收次数', '垃圾回收占用CPU比例'
]

// 应用程序监控数据
export interface ApplicationMonitor extends Record<string, number | string | number[] | string[] | Record<string, any>[]> {
  hostname: string; // 主机名称
  platform: string; // 主机平台
  platformVersion: string; // 系统版本
  cpuNumber: number; // cpu核心数
  cpuPercent: number; // cpu使用率
  memoryTotal: number; // 总内存
  memoryUsed: number; // 已用内存
  memoryAvailable: number; // 可用内存
  memoryUsedPercent: number; // 内存使用率
  Sys: number; // 程序从系统获取的内存大小
  HeapSys: number; // 堆内存大小
  HeapIdle: number; // 堆空闲内存
  HeapInuse: number; // 堆已用内存
  StackSys: number; // 栈内存大小
  StackInuse: number; // 栈已用内存
  PauseTotalNs: number; // 垃圾回收总占用时间 单位：纳秒
  PauseNs: number[]; // 每次垃圾回收程序暂停的时间 单位：纳秒 256长度的数组
  NumGC: number; // 垃圾回收的次数
  GCCPUFraction: number; // 垃圾回收占用的CPU时间比例
}