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