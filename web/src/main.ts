import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import permissionDirective from './directives/permission'
import { setOnUnauthorized } from './api'
import './assets/styles/global.css'

setOnUnauthorized(() => {
  router.push('/login')
})

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.directive('permission', permissionDirective)

app.mount('#app')