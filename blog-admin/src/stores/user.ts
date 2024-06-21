import { defineStore } from 'pinia'
import type { Menu } from '@/api/system/menu/types'
import { computed, reactive, ref } from 'vue'
import type { RouteRecordRaw } from 'vue-router'

export const useLocalUserStore = defineStore('user', () => {
  // 用户详情
  const userInfo = reactive({});
  // 树形菜单
  const treeMenu = ref<Menu[]>([]);
  // 菜单列表
  const menuList = ref<Menu[]>([]);
  // 带路由的菜单列表
  const menuRouteList = ref<RouteRecordRaw[]>([]);
  // 菜单父节点id集合Map
  const menuParentIdListMap = new Map<string, string[]>();

  const getTreeMenu = computed(() => {
    return treeMenu.value;
  })

  const getMenuList = computed(() => {
    return menuList.value;
  })

  const getMenuParentIdListMap = () => {
    return menuParentIdListMap;
  }

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
    menuList.value = tempList;
    tempList.forEach(item => {
      const parentIdList: string[] = [];
      searchMenuParent(item, parentIdList)
      menuParentIdListMap.set(item.menuId.toString(), parentIdList);
    })
  }

  const searchMenuParent = (menu: Menu, parentIdList: string[]) => {
    if (!menu.parentId || menu.parentId === 0) {
      return;
    }
    const parent = menuList.value.find(item => item.menuId === menu.parentId);
    if (parent) {
      parentIdList.unshift(menu.parentId.toString());
      searchMenuParent(parent, parentIdList);
    }
  }

  return { userInfo, setTreeMenu, treeMenu: getTreeMenu, menuList: getMenuList, getMenuParentIdListMap, menuRouteList }
})