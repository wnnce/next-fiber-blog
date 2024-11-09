import { querySiteConfiguration } from '@/lib/api'
import { Result } from '@/lib/types'

export declare type SiteConfigurationType = 'image' | 'text' | 'number' | 'markdown' | 'html' | 'color';

export interface SiteConfigurationItem {
  name: string;
  type: SiteConfigurationType;
  value: string | number;
  extend: boolean;
}

let siteConfigurationCache: Record<string, SiteConfigurationItem>
let siteConfigurationPromiseCache: Promise<Result<Record<string, SiteConfigurationItem>>> | undefined = undefined;

/**
 * 查询指定的站点配置项
 * @param keys 需要查询的配置key
 */
export async function querySiteConfigs(...keys: string[]): Promise<SiteConfigurationItem[]> {
  if (siteConfigurationCache) {
    return keys.map(key => siteConfigurationCache[key]);
  }
  if (!siteConfigurationPromiseCache) {
    siteConfigurationPromiseCache = querySiteConfiguration();
  }
  const result = await siteConfigurationPromiseCache;
  siteConfigurationPromiseCache = undefined;
  if (result.code === 200) {
    siteConfigurationCache = result.data;
    return keys.map(key => siteConfigurationCache[key]);
  }
  return []
}