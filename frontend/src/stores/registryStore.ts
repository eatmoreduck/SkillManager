import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Registry, RegistrySkill } from '@/types'

export const useRegistryStore = defineStore('registry', () => {
  const registries = ref<Registry[]>([])
  const currentRegistry = ref<Registry | null>(null)
  const browseResults = ref<RegistrySkill[]>([])
  const searchResults = ref<RegistrySkill[]>([])
  const searchQuery = ref('')
  const loading = ref(false)
  const browsing = ref(false)
  const searching = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const hasRegistries = computed(() => registries.value.length > 0)
  const defaultRegistry = computed(() => registries.value.find(r => r.isDefault))

  // 加载 registries
  async function loadRegistries() {
    loading.value = true
    error.value = null
    try {
      const result = await window.go.main.App.RegistryBinding.ListRegistries()
      registries.value = result || []
      if (result && result.length > 0 && !currentRegistry.value) {
        currentRegistry.value = result.find(r => r.isDefault) || result[0]
      }
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      loading.value = false
    }
  }

  // 切换当前 registry
  function switchRegistry(id: string) {
    const registry = registries.value.find(r => r.id === id)
    if (registry) {
      currentRegistry.value = registry
    }
  }

  // 添加 registry
  async function addRegistry(registry: Registry) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.RegistryBinding.AddRegistry(registry)
      await loadRegistries()
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 删除 registry
  async function removeRegistry(id: string) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.RegistryBinding.RemoveRegistry(id)
      registries.value = registries.value.filter(r => r.id !== id)
      if (currentRegistry.value?.id === id) {
        currentRegistry.value = registries.value[0] || null
      }
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 浏览 registry
  async function browse(registryId?: string, category?: string) {
    browsing.value = true
    error.value = null
    try {
      const rid = registryId || currentRegistry.value?.id || ''
      const cat = category || ''
      const result = await window.go.main.App.RegistryBinding.Browse(rid, cat)
      browseResults.value = result || []
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      browsing.value = false
    }
  }

  // 搜索
  async function search(query: string) {
    if (!query.trim()) {
      searchResults.value = []
      searchQuery.value = ''
      return
    }

    searching.value = true
    error.value = null
    try {
      const result = await window.go.main.App.RegistryBinding.Search(query)
      searchResults.value = result || []
      searchQuery.value = query
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      searching.value = false
    }
  }

  // 清除搜索
  function clearSearch() {
    searchResults.value = []
    searchQuery.value = ''
  }

  // 清除错误
  function clearError() {
    error.value = null
  }

  return {
    registries,
    currentRegistry,
    browseResults,
    searchResults,
    searchQuery,
    loading,
    browsing,
    searching,
    error,
    hasRegistries,
    defaultRegistry,
    loadRegistries,
    switchRegistry,
    addRegistry,
    removeRegistry,
    browse,
    search,
    clearSearch,
    clearError
  }
})
