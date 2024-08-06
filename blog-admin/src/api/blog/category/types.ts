// 博客分类
export interface Category {
  // 分类Id
  categoryId: number;
  // 分类名称
  categoryName: string;
  // 类型描述
  description: string;
  // 封面
  coverUrl: string;
  // 查看次数
  viewNum: number;
  // 上级Id
  parentId: number;
  // 是否热门分类
  isHot: boolean;
  // 分类是否置顶
  isTop: boolean;
  // 排序
  sort: number;
  // 状态
  status: number;
  // 创建时间
  createTime: string;
  // 子分类
  children: Category[]
}

export interface ArticleCategory {
  categoryId: number;
  categoryName: string;
}

// 分类表单
export interface CategoryForm {
  categoryId?: number;
  categoryName: string;
  description?: string;
  coverUrl: string;
  parentId?: number;
  isHot: boolean;
  isTop: boolean;
  sort: number;
  status: number;
}

// 分类快捷更新表单
export interface CategoryUpdateForm {
  categoryId: number;
  isHot?: boolean;
  isTop?: boolean;
  status?: number
}