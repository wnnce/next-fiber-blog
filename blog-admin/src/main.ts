import './assets/style/main.scss'

import { createApp } from 'vue'
import { Message } from '@arco-design/web-vue';

import App from './App.vue'
import router from './router'
import { stores } from '@/stores'

const app = createApp(App)

Message._context = app._context;

app.use(stores)
app.use(router)

app.mount('#app')
