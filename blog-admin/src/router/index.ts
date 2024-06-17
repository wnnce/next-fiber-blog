import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import * as NProgress from "nprogress";
import "nprogress/nprogress.css";

const publicRoute: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'layout',
    component: () => import('@/layout/index.vue'),
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/system/login.vue')
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: publicRoute,
})

router.beforeEach((to, form, next) => {
  NProgress.start();
  next()
})

router.afterEach(() => {
  NProgress.done();
})

router.getRoutes()

export default router
