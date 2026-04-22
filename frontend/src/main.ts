import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { bootstrapWailsCompat } from '@/lib/wailsCompat'
import { setLocale } from '@/i18n'
import { useConfigStore } from '@/stores/configStore'

async function bootstrap() {
  try {
    await bootstrapWailsCompat()
  } catch (error) {
    console.error('[wails] failed to bootstrap compatibility bridge', error)
  }

  const app = createApp(App)
  const pinia = createPinia()

  app.use(pinia)
  app.use(router)

  const configStore = useConfigStore(pinia)
  await configStore.loadConfig()
  setLocale(configStore.config?.language ?? navigator.language)

  app.mount('#app')
}

void bootstrap()
