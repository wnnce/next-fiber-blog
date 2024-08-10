// 站点配置
import { sendGet, sendPut } from '@/api/request'

export declare type SiteConfigurationType = 'image' | 'text' | 'number' | 'markdown' | 'html' | 'color';

export interface SiteConfigurationItem {
  name: string;
  type: SiteConfigurationType;
  value: string | number;
  extend: boolean;
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
  }
}