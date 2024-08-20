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