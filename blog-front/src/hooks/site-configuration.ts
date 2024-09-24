import { querySiteConfiguration } from '@/lib/api'

export declare type SiteConfigurationType = 'image' | 'text' | 'number' | 'markdown' | 'html' | 'color';

export interface SiteConfigurationItem {
  name: string;
  type: SiteConfigurationType;
  value: string | number;
  extend: boolean;
}

const useSiteConfiguration = () => {
  async function queryConfigs(...keys: string[]): Promise<SiteConfigurationItem[]> {
    const result = await querySiteConfiguration();
    if (result.code === 200 && result.data) {
      return keys.map(key => result.data[key]);
    }
    return []
  }
  return { queryConfigs }
}

export default useSiteConfiguration;