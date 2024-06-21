// 菜单对象
export interface Menu {
  // 菜单id
  menuId: number;
  // 菜单名词
  menuName: string;
  // 菜单类型 1：目录 2：菜单
  menuType: number;
  // 上级菜单id
  parentId: number;
  // 菜单路由地址
  path: string;
  // 菜单组件地址
  component: string;
  // 菜单图标名称
  icon: string;
  // 是否frame窗口
  isFrame: boolean;
  // frame窗口地址
  frameUrl: string;
  // 是否缓存
  isCache: boolean;
  // 是否可见
  isVisible: boolean;
  // 是否关闭
  isDisable: boolean;
  // 排序
  sort: number;
  // 子菜单
  children?: Menu[]
}