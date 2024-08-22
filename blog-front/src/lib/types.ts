/**
 * 统一数据返回对象
 */
export interface Result<T> {
  code: number;
  message: string;
  timestamp: number;
  data: T,
}

/**
 * 博客标签数据
 */
export interface Tag {
  tagId: number;
  tagName: string;
  coverUrl: string;
  viewNum: number;
  color: string;
  articleNum: number;
}

/**
 * 博客分类数据
 */
export interface Category {
  categoryId: number;
  categoryName: string;
  coverUrl: string;
  viewNum: number;
  parentId: number;
  isHot: boolean;
  isTop: boolean;
  articleNum: number;
  children?: Category[]
}

/**
 * 联系方式
 */
export interface Concat {
  concatId: number;
  name: string;
  iconSvg: string;
  targetUrl: string;
  isMain: boolean;
}

/**
 * 通知公告数据
 */
export interface Notice {
  noticeId: number;
  title: string;
  message: string;
  level: number; // 通知级别 1:info 2:warn 3:error
  noticeType: number;
}