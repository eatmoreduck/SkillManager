import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Config, ProxyConfig } from '@/types'

export const useConfigStore = defineStore('config', () => {
  const config = ref<Config | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function loadConfig() {
    loading.value = true
    error.value = null
    try {
      const result = await window.go.main.App.ConfigBinding.GetConfig()
      config.value = result
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      loading.value = false
    }
  }

  async function updateProxy(proxy: ProxyConfig) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.ConfigBinding.UpdateProxy(proxy)
      if (config.value) {
        config.value.proxy = proxy
      }
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    config,
    loading,
    error,
    loadConfig,
    updateProxy
  }
})
