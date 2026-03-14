// Agent 运行时结构
export interface Agent {
  id: string
  name: string
  skillsDir: string
  binaryName: string
  priorityPaths: string[]
  isInstalled: boolean
  isEnabled: boolean
  isCustom: boolean
}
