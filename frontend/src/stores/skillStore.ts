import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Skill } from '@/types'

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

function normalizeStringArray(value: unknown): string[] {
  if (!Array.isArray(value)) return []
  return value
    .map(item => normalizeString(item))
    .filter(Boolean)
}

function normalizeDateString(value: unknown): string {
  if (typeof value === 'string') return value
  if (value instanceof Date) return value.toISOString()
  if (value && typeof value === 'object') {
    const obj = value as Record<string, unknown>
    const nested = pickFirst(obj, ['time', 'Time', 'value', 'Value'])
    if (nested) return normalizeDateString(nested)
  }
  return ''
}

function normalizeSkill(item: unknown): Skill | null {
  if (!item || typeof item !== 'object') {
    return null
  }

  const obj = item as Record<string, unknown>
  const id = normalizeString(pickFirst(obj, ['id', 'ID']))
  const name = normalizeString(pickFirst(obj, ['name', 'Name']), id || 'Untitled Skill')

  return {
    id: id || name,
    name,
    description: normalizeString(pickFirst(obj, ['description', 'Description'])),
    author: normalizeString(pickFirst(obj, ['author', 'Author'])),
    version: normalizeString(pickFirst(obj, ['version', 'Version'])),
    tags: normalizeStringArray(pickFirst(obj, ['tags', 'Tags'])),
    agents: normalizeStringArray(pickFirst(obj, ['agents', 'Agents'])),
    content: normalizeString(pickFirst(obj, ['content', 'Content'])),
    localPath: normalizeString(pickFirst(obj, ['localPath', 'LocalPath'])),
    sourceUrl: normalizeString(pickFirst(obj, ['sourceUrl', 'SourceURL', 'sourceURL'])),
    installedAt: normalizeDateString(pickFirst(obj, ['installedAt', 'InstalledAt'])),
    updatedAt: normalizeDateString(pickFirst(obj, ['updatedAt', 'UpdatedAt'])),
    isManaged: Boolean(pickFirst(obj, ['isManaged', 'IsManaged']))
  }
}

/**
 * Wails v3 Call.ByName 返回值可能不是纯数组（如包装对象或 nil/undefined）。
 * 此函数确保结果始终为 Skill[]。
 */
function normalizeSkillArray(data: unknown): Skill[] {
  if (Array.isArray(data)) {
    return data
      .map(item => normalizeSkill(item))
      .filter((item): item is Skill => item !== null)
  }
  // 兜底：某些场景下 Wails 可能将 []model.Skill 序列化为对象包装
  if (data && typeof data === 'object') {
    const obj = data as Record<string, unknown>
    // 尝试常见的包装字段: data, skills, result, items
    for (const key of ['data', 'skills', 'result', 'items']) {
      if (Array.isArray(obj[key])) {
        return normalizeSkillArray(obj[key])
      }
    }
    // 如果对象本身看起来像单个 skill，包装为数组
    const normalized = normalizeSkill(obj)
    if (normalized) {
      return [normalized]
    }
  }
  return []
}

export const useSkillStore = defineStore('skill', () => {
  const skills = ref<Skill[]>([])
  const currentSkill = ref<Skill | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const skillCount = computed(() => skills.value.length)
  const installedSkillIds = computed(() => skills.value.map(s => s.id))

  // 加载已安装的 skills
  async function loadSkills() {
    loading.value = true
    error.value = null
    try {
      const result = await window.go.main.App.SkillBinding.ListInstalled()
      console.info('[skillStore] ListInstalled raw result:', typeof result, Array.isArray(result), result)
      skills.value = normalizeSkillArray(result)
      console.info('[skillStore] normalized skills count:', skills.value.length)
    } catch (e) {
      console.error('[skillStore] ListInstalled failed:', e)
      error.value = (e as Error).message
    } finally {
      loading.value = false
    }
  }

  // 获取 skill 详情
  async function getSkillDetail(id: string) {
    loading.value = true
    error.value = null
    try {
      const result = await window.go.main.App.SkillBinding.GetDetail(id)
      const normalized = normalizeSkill(result)
      if (!normalized) {
        throw new Error('Invalid skill payload returned from detail request')
      }
      currentSkill.value = normalized
      return normalized
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 安装 skill
  async function installSkill(sourceURL: string, agents: string[]) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.SkillBinding.Install(sourceURL, agents)
      await loadSkills()
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 卸载 skill
  async function uninstallSkill(id: string) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.SkillBinding.Uninstall(id)
      skills.value = skills.value.filter(s => s.id !== id)
      if (currentSkill.value?.id === id) {
        currentSkill.value = null
      }
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 更新 skill (git pull)
  async function updateSkill(id: string) {
    loading.value = true
    error.value = null
    try {
      const result = await window.go.main.App.SkillBinding.Update(id)
      const normalized = normalizeSkill(result)
      if (!normalized) {
        throw new Error('Invalid skill payload returned from update')
      }
      const index = skills.value.findIndex(s => s.id === id)
      if (index !== -1) {
        skills.value[index] = normalized
      }
      if (currentSkill.value?.id === id) {
        currentSkill.value = normalized
      }
      return normalized
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 更新 skill 内容
  async function updateSkillContent(id: string, content: string) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.SkillBinding.UpdateContent(id, content)
      const skill = skills.value.find(s => s.id === id)
      if (skill) {
        skill.content = content
        skill.updatedAt = new Date().toISOString()
      }
      if (currentSkill.value?.id === id) {
        currentSkill.value.content = content
      }
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 分配 agents
  async function assignAgents(id: string, agents: string[]) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.SkillBinding.AssignAgents(id, agents)
      const skill = skills.value.find(s => s.id === id)
      if (skill) {
        skill.agents = agents
      }
      if (currentSkill.value?.id === id) {
        currentSkill.value.agents = agents
      }
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 根据 ID 获取 skill
  function getSkillById(id: string): Skill | undefined {
    return skills.value.find(s => s.id === id)
  }

  // 检查 skill 是否已安装
  function isInstalled(id: string): boolean {
    return skills.value.some(s => s.id === id)
  }

  return {
    skills,
    currentSkill,
    loading,
    error,
    skillCount,
    installedSkillIds,
    loadSkills,
    getSkillDetail,
    installSkill,
    uninstallSkill,
    updateSkill,
    updateSkillContent,
    assignAgents,
    getSkillById,
    isInstalled
  }
})
