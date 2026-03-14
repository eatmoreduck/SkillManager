import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Agent } from '@/types'

export const useAgentStore = defineStore('agent', () => {
  const agents = ref<Agent[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const enabledAgents = computed(() => agents.value.filter(a => a.isEnabled))
  const installedAgents = computed(() => agents.value.filter(a => a.isInstalled))
  const customAgents = computed(() => agents.value.filter(a => a.isCustom))

  // 加载 agents
  async function loadAgents() {
    loading.value = true
    error.value = null
    try {
      const result = await window.go.main.App.AgentBinding.ListAgents()
      agents.value = result || []
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      loading.value = false
    }
  }

  // 检测已安装的 agents
  async function detectInstalled() {
    loading.value = true
    error.value = null
    try {
      const result = await window.go.main.App.AgentBinding.DetectInstalled()
      console.info('[agentStore] detectInstalled found', result?.length ?? 0, 'agents')
      await loadAgents()
      return result
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 添加自定义 agent
  async function addCustomAgent(agent: Agent) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.AgentBinding.AddCustomAgent(agent)
      await loadAgents()
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 删除 agent
  async function removeAgent(id: string) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.AgentBinding.RemoveAgent(id)
      agents.value = agents.value.filter(a => a.id !== id)
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 切换 agent 启用状态
  async function toggleAgent(id: string, enabled: boolean) {
    loading.value = true
    error.value = null
    try {
      await window.go.main.App.AgentBinding.ToggleAgent(id, enabled)
      const agent = agents.value.find(a => a.id === id)
      if (agent) {
        agent.isEnabled = enabled
      }
    } catch (e) {
      error.value = (e as Error).message
      throw e
    } finally {
      loading.value = false
    }
  }

  // 根据 ID 获取 agent
  function getAgentById(id: string): Agent | undefined {
    return agents.value.find(a => a.id === id)
  }

  return {
    agents,
    loading,
    error,
    enabledAgents,
    installedAgents,
    customAgents,
    loadAgents,
    detectInstalled,
    addCustomAgent,
    removeAgent,
    toggleAgent,
    getAgentById
  }
})
