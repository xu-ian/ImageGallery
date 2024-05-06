import { createApp } from 'vue'
import App from './App.vue'
import { createWebHistory, createRouter } from 'vue-router'
import './assets/style.scss';
import VueCookies from 'vue3-cookies'

import HelloWorld from './components/HelloWorld.vue'
import AddFile from './components/AddFile.vue'
import ViewImages from './components/ViewImages.vue'
import AuthenticateUser from './components/Authentication.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
  { path: '/', component: HelloWorld },
  { path: '/files', component: AddFile},
  { path: '/gallery', component: ViewImages},
  { path: '/authentication', component: AuthenticateUser}
  ],
})

createApp(App).use(router).use(VueCookies).mount('#app')
