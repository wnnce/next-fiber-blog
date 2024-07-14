import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import * as NProgress from "nprogress";
import "nprogress/nprogress.css";
import { publicRouter } from '@/router/routers'
import type { Menu } from '@/api/system/menu/types'
import { useArcoMessage } from '@/hooks/message'

const { errorMessage } = useArcoMessage();

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: publicRouter,
  strict: true,
  scrollBehavior: () => ({left: 0, top: 0})
})

router.beforeEach((to, form, next) => {
  NProgress.start();
  const isDisable = to.meta ? to.meta.isDisable as boolean : true;
  if (isDisable) {
    errorMessage('路由地址被禁用');
    return;
  }
  next()
})

router.afterEach(() => {
  NProgress.done();
})

router.getRoutes()

export default router

export const buildRoute = (treeMenu: Menu[]): RouteRecordRaw[] => {
  const menuList: Menu[] = parseTreeMenuToMenuList(treeMenu);
  const views = import.meta.glob('../views/**/**.vue')
  const routeList: RouteRecordRaw[] = menuList.map(item => formatMenuToRoute(item, views))
  router.addRoute({
    path: '/',
    component: () => import('@/layout/index.vue'),
    children: routeList,
  });
  return routeList;
}

const parseTreeMenuToMenuList = (treeMenu: Menu[], prefixPath?: string): Menu[] => {
  let menuList: Menu[] = []
  treeMenu.forEach(item => {
    if (item.menuType === 2) {
      item.path = item.path.startsWith('/') ? item.path : `/${item.path}`;
      if (prefixPath) {
        item.path = prefixPath + item.path
      }
      menuList.push(item);
    } else if (item.children && item.children.length > 0) {
      const childrenList = parseTreeMenuToMenuList(item.children, item.path);
      menuList = menuList.concat(childrenList);
    }
  })
  return menuList;
}

const formatMenuToRoute = (menu: Menu, views: Record<string, () => Promise<unknown>>): RouteRecordRaw => {
  let componentPath;
  if (menu.isFrame) {
    componentPath = ''
  } else {
    const prefixPath = menu.component.startsWith('/') ? menu.component : `/${menu.component}`;
    componentPath = prefixPath.endsWith('.vue') ? prefixPath : `${prefixPath}.vue`;
  }
  const menuIdStr = menu.menuId.toString();
  const findComponent = views[`../views${componentPath}`] || undefined;
  findComponent && ((findComponent() as Promise<{ default: { name: string } }>)
    .then(val => val && (val.default.name = menuIdStr)))
  return {
    path: menu.path.startsWith('/') ? menu.path : `/${menu.path}`,
    name: menuIdStr,
    meta: {
      id: menuIdStr,
      componentName: menu.menuName,
      icon: menu.icon,
      isVisible: menu.isVisible,
      isDisable: menu.isDisable,
      keepAlive: menu.isCache,
      isFrame: menu.isFrame,
      frameUrl: menu.frameUrl
    },
    component: findComponent,
  }
}