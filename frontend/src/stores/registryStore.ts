import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Registry, RegistrySkill } from '@/types'

function pickFirst(obj: Record<string, unknown>, keys: string[]): unknown {
  for (const key of keys) {
    if (obj[key] !== undefined && obj[key] !== null) {
      return obj[key]
    }
  }
  return undefined
}

function normalizeString(value: unknown, fallback = ''): string {
  if (typeof value === 'string') return value
  if (typeof value === 'number' || typeof value === 'boolean') return String(value)
  return fallback
}

function normalizeBoolean(value: unknown): boolean {
  if (typeof value === 'boolean') return value
  if (typeof value === 'string') {
    return value.toLowerCase() === 'true'
  }
  if (typeof value === 'number') return value !== 0
  return false
}

function normalizeNumber(value: unknown): number {
  if (typeof value === 'number' && Number.isFinite(value)) return value
  if (typeof value === 'string') {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : 0
  }
  return 0
}

function normalizeStringArray(value: unknown): string[] {
  if (!Array.isArray(value)) return []
  return value
    .map(item => normalizeString(item))
    .filter(Boolean)
}

function normalizeRegistry(item: unknown): Registry | null {
  if (!item || typeof item !== 'object') {
    return null
  }

  const obj = item as Record<string, unknown>
  const id = normalizeString(pickFirst(obj, ['id', 'ID']))
  const name = normalizeString(pickFirst(obj, ['name', 'Name']), id)

  return {
    id: id || name,
    name,
    url: normalizeString(pickFirst(obj, ['url', 'URL'])),
    isDefault: normalizeBoolean(pickFirst(obj, ['isDefault', 'IsDefault']))
  }
}

function normalizeRegistrySkill(item: unknown): RegistrySkill | null {
  if (!item || typeof item !== 'object') {
    return null
  }

  const obj = item as Record<string, unknown>
  const id = normalizeString(pickFirst(obj, ['id', 'ID']))
  const name = normalizeString(pickFirst(obj, ['name', 'Name']), id || 'Untitled Skill')
  const author = normalizeString(pickFirst(obj, ['author', 'Author', 'source', 'Source']))

  return {
    id: id || name,
    name,
    description: normalizeString(pickFirst(obj, ['description', 'Description'])),
    author,
    stars: normalizeNumber(pickFirst(obj, ['stars', 'Stars', 'installs', 'Installs'])),
    tags: normalizeStringArray(pickFirst(obj, ['tags', 'Tags'])),
    installUrl: normalizeString(pickFirst(obj, ['installUrl', 'InstallURL', 'installURL'])),
    category: normalizeString(pickFirst(obj, ['category', 'Category']))
  }
}

function normalizeRegistryArray(data: unknown): Registry[] {
  if (!Array.isArray(data)) return []
  return data
    .map(item => normalizeRegistry(item))
    .filter((item): item is Registry => item !== null)
}

function normalizeRegistrySkillArray(data: unknown): RegistrySkill[] {
  if (!Array.isArray(data)) return []
  return data
    .map(item => normalizeRegistrySkill(item))
    .filter((item): item is RegistrySkill => item !== null)
}

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
      console.info('[registryStore] ListRegistries raw result:', typeof result, Array.isArray(result), result)
      registries.value = normalizeRegistryArray(result)
      if (registries.value.length > 0 && !currentRegistry.value) {
        currentRegistry.value = registries.value.find(r => r.isDefault) || registries.value[0]
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
      console.info('[registryStore] Browse raw result:', typeof result, Array.isArray(result), result)
      browseResults.value = normalizeRegistrySkillArray(result)
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
      console.info('[registryStore] Search raw result:', typeof result, Array.isArray(result), result)
      searchResults.value = normalizeRegistrySkillArray(result)
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
