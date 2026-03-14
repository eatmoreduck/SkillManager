import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { bootstrapWailsCompat } from '@/lib/wailsCompat'

async function bootstrap() {
  try {
    await bootstrapWailsCompat()
  } catch (error) {
    console.error('[wails] failed to bootstrap compatibility bridge', error)
  }

  const app = createApp(App)

  app.use(createPinia())
  app.use(router)

  app.mount('#app')
}

void bootstrap()
