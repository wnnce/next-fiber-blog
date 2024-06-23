import { defineStore } from 'pinia'
import type { Menu } from '@/api/system/menu/types'
import { computed, reactive, ref } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import type { KeepaliveItem } from '@/assets/script/types'
import { useAppConfigStore } from '@/stores/app-config'

export const useLocalUserStore = defineStore('user', () => {
  // 用户详情
  const userInfo = reactive({});

  // 带路由的菜单Map
  const menuRouteMap = new Map<string, RouteRecordRaw>();
  const _menuRouteMap = computed(() => {
    return menuRouteMap;
  })
  const setMenuRoute = (menuRoutes: RouteRecordRaw[]) => {
    menuRoutes.forEach(item => menuRouteMap.set(item.name ? item.name.toString() : '', item));
  }
  // 树形菜单
  const treeMenu = ref<Menu[]>([]);
  const _treeMenu = computed(() => {
    return treeMenu.value;
  })
  /**
   * 保存后端返回的属性菜单列表
   * @param menus 树形菜单列表
   * 保存的同时将树形菜单遍历一遍，将所有的菜单展开保存到map中
   * 然后再拿到子菜单的所有上级菜单id 保存到map中
   */
  const setTreeMenu = (menus: Menu[]) => {
    treeMenu.value = menus;
    const tempList: Menu[] = [];
    function dfs(menu: Menu) {
      if (!menu) {
        return;
      }
      tempList.push(menu);
      if (menu.children && menu.children.length > 0) {
        menu.children.forEach(item => dfs(item));
      }
    }
    menus.forEach(item => dfs(item));
    tempList.forEach(item => menuListMap.set(item.menuId.toString(), item));
    tempList.forEach(item => {
      const parentIdList: string[] = [];
      searchMenuParent(item, parentIdList)
      menuParentIdListMap.set(item.menuId.toString(), parentIdList);
    })
  }


  // 菜单列表
  const menuListMap = new Map<string, Menu>();
  const _menuListMap = computed(() => {
    return menuListMap;
  })

  const getMenuParentIdListMap = () => {
    return menuParentIdListMap;
  }

  // 菜单父节点id集合Map
  const menuParentIdListMap = new Map<string, string[]>();
  /**
   * 查询菜单的上级id列表 递归操作 如果上级id不存在或者 === 0的话 就直接返回
   * @param menu 需要查询的菜单
   * @param parentIdList 上级菜单列表
   */
  const searchMenuParent = (menu: Menu, parentIdList: string[]) => {
    if (!menu.parentId || menu.parentId === 0) {
      return;
    }
    const parent = menuListMap.get(menu.parentId.toString());
    if (parent) {
      parentIdList.unshift(menu.parentId.toString());
      searchMenuParent(parent, parentIdList);
    }
  }


  // 缓存的组件id列表
  const keepaliveList = ref<KeepaliveItem[]>([]);
  const _keepaliveList = computed(() => {
    return keepaliveList.value;
  })
  // keepalive组件需要缓存的组件name列表
  const keepaliveInclude = computed((): string[] => {
    return keepaliveList.value.filter(item => item.isCache === true)
      .map(item => item.menuId);
  })
  const addKeepaliveComponent = (menuId: string, isCache: boolean) => {
    let item: KeepaliveItem;
    if (menuId === '-1') {
      item = {
        menuId: menuId,
        path: '/index',
        name: '首页',
        isCache: isCache
      }
    } else {
      const route = menuRouteMap.get(menuId);
      if (!route) {
        return;
      }
      item = {
        menuId: menuId,
        path: route.path,
        name: route.meta ? route.meta.componentName as string : '',
        isCache: isCache
      }
    }
    if (keepaliveList.value.length >= useAppConfigStore().state.maxCachePage) {
      keepaliveList.value.shift();
    }
    keepaliveList.value.push(item);
  }
  // 通过菜单id查询该组件是否被缓存
  const queryKeepaliveComponent = (menuId: string): KeepaliveItem | undefined => {
    return keepaliveList.value.find(item => item.menuId === menuId);
  }
  // 删除某一个被缓存的组件
  const removeKeepaliveComponent = (menuId: string) => {
    keepaliveList.value = keepaliveList.value.filter(item => item.menuId !== menuId);
  }
  // 删除所有被缓存的组件
  const removeAllKeepaliveComponent = () => {
    keepaliveList.value = [];
  }

  return { userInfo, setTreeMenu, treeMenu: _treeMenu, menuListMap: _menuListMap, getMenuParentIdListMap, menuRouteMap: _menuRouteMap,
    setMenuRoute, keepaliveInclude, addKeepaliveComponent, queryKeepaliveComponent, removeKeepaliveComponent, removeAllKeepaliveComponent,
    keepaliveList: _keepaliveList }
})