import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Skill } from '@/types'

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
      skills.value = result || []
    } catch (e) {
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
      currentSkill.value = result
      return result
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
      const index = skills.value.findIndex(s => s.id === id)
      if (index !== -1) {
        skills.value[index] = result
      }
      if (currentSkill.value?.id === id) {
        currentSkill.value = result
      }
      return result
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
