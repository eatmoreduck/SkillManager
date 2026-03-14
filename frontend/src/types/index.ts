export * from './config'
export * from './skill'
export * from './agent'
export * from './registry'

import type { Skill } from './skill'
import type { Agent } from './agent'
import type { RegistrySkill } from './registry'
import type { Config, ProxyConfig, Registry } from './config'

// Wails 绑定类型声明
declare global {
  interface Window {
    go: {
      main: {
        App: {
          SkillBinding: {
            ListInstalled: () => Promise<Skill[]>
            GetDetail: (id: string) => Promise<Skill>
            Install: (sourceURL: string, agents: string[]) => Promise<void>
            Uninstall: (id: string) => Promise<void>
            Update: (id: string) => Promise<Skill>
            UpdateContent: (id: string, content: string) => Promise<void>
            AssignAgents: (id: string, agents: string[]) => Promise<void>
          }
          AgentBinding: {
            ListAgents: () => Promise<Agent[]>
            DetectInstalled: () => Promise<Agent[]>
            AddCustomAgent: (agent: Agent) => Promise<void>
            RemoveAgent: (id: string) => Promise<void>
            ToggleAgent: (id: string, enabled: boolean) => Promise<void>
          }
          RegistryBinding: {
            ListRegistries: () => Promise<Registry[]>
            AddRegistry: (registry: Registry) => Promise<void>
            RemoveRegistry: (id: string) => Promise<void>
            Browse: (registryID: string, category: string) => Promise<RegistrySkill[]>
            Search: (query: string) => Promise<RegistrySkill[]>
          }
          ConfigBinding: {
            GetConfig: () => Promise<Config>
            UpdateProxy: (proxy: ProxyConfig) => Promise<void>
          }
        }
      }
    }
    wails?: {
      Call: {
        ByName: <T = unknown>(methodName: string, ...args: unknown[]) => Promise<T>
      }
    }
  }
}
