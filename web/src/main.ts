import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import permissionDirective from './directives/permission'
import http, { setOnUnauthorized } from './api'
import './assets/styles/global.css'

async function bootstrap() {
  if (import.meta.env.DEV) {
    const { setupMockAdapter } = await import('./api/mock')
    setupMockAdapter(http)
  }

  setOnUnauthorized(() => {
    router.push('/login')
  })

  const app = createApp(App)

  app.use(createPinia())
  app.use(router)
  app.directive('permission', permissionDirective)

  app.mount('#app')
}

bootstrap()
