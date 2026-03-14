// Skill 表示已安装的技能
export interface Skill {
  id: string
  name: string
  description: string
  author: string
  version: string
  tags: string[]
  agents: string[]
  content: string
  localPath: string
  sourceUrl: string
  installedAt: string
  updatedAt: string
}
