import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import permissionDirective from './directives/permission'
import './assets/styles/global.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.directive('permission', permissionDirective)

app.mount('#app')