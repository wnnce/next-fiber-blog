import './assets/style/main.scss'
import { createApp } from 'vue'
import { Message } from '@arco-design/web-vue';
import App from './App.vue'
import router, { buildRoute } from './router'
import { stores } from '@/stores'
import { menuApi } from '@/api/system/menu'
import { useLocalUserStore } from '@/stores/user'
import { useLocalStorage } from '@/hooks/local-storage'
import { LOCAl_USER_KEY, TOKEN_KEY } from '@/assets/script/constant'
import type { User } from '@/api/system/user/types'

const mountVue = () => {
  const app = createApp(App)
  Message._context = app._context;
  app.use(stores)
  app.use(router)
  app.mount('#app')
}

// 用户已登录时启动应用 查询菜单 获取本地保存的用户信息
const queryMenu = async () => {
  try {
    const result = await menuApi.listTreeMenu();
    const { code, data } = result;
    if (code === 200) {
      const routeList = buildRoute(data)
      mountVue();
      useLocalUserStore().setTreeMenu(data);
      useLocalUserStore().setMenuRoute(routeList);
      const localUser = useLocalStorage().get<User>(LOCAl_USER_KEY);
      localUser && (useLocalUserStore().userInfo = localUser);
    }
  } catch (e) {
    console.log(e);
    mountVue();
  }
}

const token = useLocalStorage().get<string>(TOKEN_KEY);

token && token.trim().length > 0 ? queryMenu() : mountVue();