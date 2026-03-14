<template>
  <n-card
    class="skill-card"
    hoverable
    @click="$emit('click')"
  >
    <template #header>
      <div class="card-header">
        <span class="skill-name">{{ skill.name }}</span>
        <n-tag v-if="skill.version" size="small" type="info">v{{ skill.version }}</n-tag>
      </div>
    </template>

    <template #header-extra>
      <n-dropdown
        trigger="click"
        :options="dropdownOptions"
        @select="handleDropdownSelect"
      >
        <n-button text @click.stop>
          <template #icon><n-icon><EllipsisVertical /></n-icon></template>
        </n-button>
      </n-dropdown>
    </template>

    <div class="skill-description">{{ skill.description }}</div>

    <div class="skill-tags">
      <n-tag
        v-for="tag in skillTags(skill).slice(0, 3)"
        :key="tag"
        size="small"
        :bordered="false"
      >
        {{ tag }}
      </n-tag>
      <n-tag v-if="skillTags(skill).length > 3" size="small" :bordered="false">
        +{{ skillTags(skill).length - 3 }}
      </n-tag>
    </div>

    <template #footer>
      <div class="skill-meta">
        <span class="install-time">{{ formatInstallTime(skill.installedAt) }}</span>
        <div class="agent-icons" v-if="skillAgents.length > 0">
          <div
            v-for="agent in skillAgents.slice(0, 4)"
            :key="agent.id"
            class="agent-icon"
            :title="agent.name"
            v-html="getAgentIcon(agent.id)"
          />
          <div v-if="skillAgents.length > 4" class="agent-more">
            +{{ skillAgents.length - 4 }}
          </div>
        </div>
      </div>
    </template>
  </n-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NCard, NTag, NButton, NIcon, NDropdown } from 'naive-ui'
import { EllipsisVertical } from '@vicons/ionicons5'
import type { Skill, Agent } from '@/types'
import { useAgentStore } from '@/stores/agentStore'

const props = defineProps<{
  skill: Skill
}>()

const emit = defineEmits<{
  click: []
  uninstall: [skillId: string]
}>()

const agentStore = useAgentStore()

const dropdownOptions = [
  { label: '查看详情', key: 'detail' },
  { label: '卸载', key: 'uninstall' }
]

const fallbackAgentIcon = `
<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
  <rect x="3" y="3" width="18" height="18" rx="5" fill="currentColor" opacity="0.16"></rect>
  <path d="M8 15.5V8.5H10.2L12 11.1L13.8 8.5H16V15.5H13.9V12.1L12.6 14H11.4L10.1 12.1V15.5H8Z" fill="currentColor"></path>
</svg>
`.trim()

const iconModules = import.meta.glob('../assets/icons/*.svg', {
  query: '?raw',
  import: 'default',
  eager: true
}) as Record<string, string>

const agentIcons = Object.fromEntries(
  Object.entries(iconModules).map(([path, svg]) => {
    const match = path.match(/\/([^/]+)\.svg$/)
    return [match?.[1] ?? path, svg]
  })
) as Record<string, string>

function getAgentIcon(agentId: string): string {
  return agentIcons[agentId] || fallbackAgentIcon
}

function formatInstallTime(time: string): string {
  if (!time) return ''
  try {
    const date = new Date(time)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))

    if (days === 0) return '今天安装'
    if (days === 1) return '昨天安装'
    if (days < 7) return `${days} 天前安装`
    if (days < 30) return `${Math.floor(days / 7)} 周前安装`
    return `${Math.floor(days / 30)} 个月前安装`
  } catch {
    return ''
  }
}

// 获取使用该 skill 的 agent 列表
const skillAgents = computed(() => {
  if (!Array.isArray(props.skill.agents)) return []
  return props.skill.agents
    .map(id => agentStore.getAgentById(id))
    .filter((a): a is Agent => a !== undefined)
})

function skillTags(skill: Skill): string[] {
  return Array.isArray(skill.tags) ? skill.tags : []
}

function handleDropdownSelect(key: string) {
  if (key === 'detail') {
    emit('click')
  } else if (key === 'uninstall') {
    emit('uninstall', props.skill.id)
  }
}
</script>

<style scoped>
.skill-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.skill-card:hover {
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.skill-name {
  font-weight: 600;
  font-size: 15px;
}

.skill-description {
  color: var(--n-text-color-2);
  font-size: 13px;
  line-height: 1.5;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.skill-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.skill-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--n-text-color-3);
}

.agent-icons {
  display: flex;
  align-items: center;
  gap: 4px;
}

.agent-icon {
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.agent-icon :deep(svg) {
  width: 100%;
  height: 100%;
}

.agent-more {
  font-size: 11px;
  color: var(--n-text-color-3);
  margin-left: 2px;
}
</style>
