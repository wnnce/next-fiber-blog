import './assets/style/main.scss'

import { createApp } from 'vue'
import { Message } from '@arco-design/web-vue';
import App from './App.vue'
import router, { buildRoute } from './router'
import { stores } from '@/stores'
import { menuApi } from '@/api/system/menu'
import { useLocalUserStore } from '@/stores/user'

const mountVue = () => {
  const app = createApp(App)
  Message._context = app._context;
  app.use(stores)
  app.use(router)
  app.mount('#app')
}
const queryMenu = async () => {
  try {
    const result = await menuApi.listTreeMenu();
    const { code, data } = result;
    if (code === 200) {
      const routeList = buildRoute(data)
      mountVue();
      useLocalUserStore().setTreeMenu(data)
      useLocalUserStore().menuRouteList = routeList;
    }
  } catch (e) {
    console.log(e);
    mountVue();
  }

}
queryMenu();