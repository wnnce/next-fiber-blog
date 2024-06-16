
export declare type PageTheme = 'light' | 'dark';

export interface AppConfig {
  // 语言
  language: string;
  // 侧边栏宽度
  sideWidth: string | number;
  // 主题
  pageTheme: PageTheme;
  // 是否显示标签页
  showPageTag: boolean;
}

// 网站默认配置
export const defaultAppConfig: AppConfig = {
  language: 'zh-CN',
  sideWidth: 200,
  pageTheme: 'light',
  showPageTag: true
}