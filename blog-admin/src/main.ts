import './assets/style/main.scss'

import { createApp } from 'vue'

import App from './App.vue'
import router from './router'
import { stores } from '@/stores'

const app = createApp(App)

app.use(stores)
app.use(router)

app.mount('#app')
