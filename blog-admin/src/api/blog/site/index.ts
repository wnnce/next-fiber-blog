// 站点配置
import { sendGet, sendPut } from '@/api/request'

export declare type SiteConfigurationType = 'image' | 'text' | 'number' | 'markdown' | 'html' | 'color';

export interface SiteConfigurationItem {
  name: string;
  type: SiteConfigurationType;
  value: string | number;
  extend: boolean;
}

// 首页统计信息
export interface IndexStats {
  toDayAccess: number;
  toDayComment: number;
  totalAccess: number;
  totalComment: number;
  totalTopic: number;
  totalArticle: number;
  totalUser: number;
  articleTotalView: number;
  accessArray: DayStats[];
  commentArray: DayStats[];
  userArray: DayStats[];
  articleArray: DayStats[];
}

export interface DayStats {
  dateItem: string;
  countItem: number;
}

export const siteConfigurationRequiredField: string[] = ['tabTitle', 'logo', 'avatar', 'title', 'summary', 'about', 'powered', 'icp', 'articleSize', 'topicSize', 'commentSize', 'primaryColor'];

export const siteApi = {
  /**
   * 查询站点配置信息
   */
  configuration: () => {
    return sendGet<Record<string, SiteConfigurationItem>>('/open/site/configuration')
  },
  /**
   * 更新站点配置信息
   * @param form 配置信息参数
   */
  updateConfiguration: (form: Record<string, SiteConfigurationItem>) => {
    return sendPut<null>('/base/site/configuration', form)
  },
  /**
   * 查询首页统计信息
   */
  indexStats: () => {
    return sendGet<IndexStats>('/base/index/stats')
  }
}