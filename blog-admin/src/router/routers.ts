import type { RouteRecordRaw } from 'vue-router'

export const publicRouter: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'layout',
    redirect: '/index',
    component: () => import('@/layout/index.vue'),
    children: [
      {
        path: '/index',
        name: 'index',
        component: () => import('@/views/system/home.vue')
      }
    ],
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/system/login.vue')
  }
]