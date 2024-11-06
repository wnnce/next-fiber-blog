import { querySiteConfiguration } from '@/lib/api'

export declare type SiteConfigurationType = 'image' | 'text' | 'number' | 'markdown' | 'html' | 'color';

export interface SiteConfigurationItem {
  name: string;
  type: SiteConfigurationType;
  value: string | number;
  extend: boolean;
}

/**
 * 查询指定的站点配置项
 * @param keys 需要查询的配置key
 */
export async function querySiteConfigs(...keys: string[]): Promise<SiteConfigurationItem[]> {
  const result = await querySiteConfiguration();
  if (result.code === 200 && result.data) {
    return keys.map(key => result.data[key]);
  }
  return []
}