// 下拉选择数据项
export interface OptionItem {
  label: string | number,
  value: string | number
}

// 缓存项
export interface KeepaliveItem {
  menuId: string,
  name: string,
  path: string,
  isCache: boolean
}

// 分页查询参数
export interface PageQuery {
  page: number;
  size: number
}

/**
 * 分页返回对象
 */
export interface Page<T> {
  // 当前页码
  current: number;
  // 总页数
  pages: number;
  // 总记录数
  total: number;
  // 每页记录数
  size: number;
  // 记录列表
  records: T[];
}