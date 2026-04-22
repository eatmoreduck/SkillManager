// 代理配置
export interface ProxyConfig {
  enabled: boolean
  type: 'http' | 'https' | 'socks5'
  host: string
  port: number
  username: string
  password: string
}

// Agent 配置（持久化）
export interface AgentConfig {
  id: string
  name: string
  skillsDir: string
  binaryName: string
  priorityPaths: string[]
  isEnabled: boolean
  isCustom: boolean
}

// Registry 配置
export interface Registry {
  id: string
  name: string
  url: string
  isDefault: boolean
}

// 应用配置
export interface Config {
  version: string
  language: string
  proxy: ProxyConfig
  registries: Registry[]
  agents: AgentConfig[]
}
