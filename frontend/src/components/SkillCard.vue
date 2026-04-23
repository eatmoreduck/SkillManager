<template>
  <article
    class="apple-card skill-card"
    @click="$emit('click')"
  >
    <div class="card-topbar">
      <div class="card-header">
        <div class="skill-avatar">{{ skillInitial }}</div>
        <div class="card-copy">
          <span class="skill-overline">{{ skill.author || t('common.author') }}</span>
          <span class="skill-name">{{ displayName }}</span>
        </div>
        <n-tag v-if="skill.version" size="small" type="info">v{{ skill.version }}</n-tag>
      </div>
      <n-dropdown
        trigger="click"
        :options="dropdownOptions"
        @select="handleDropdownSelect"
      >
        <n-button text class="card-menu-button" @click.stop>
          <template #icon><n-icon><EllipsisVertical /></n-icon></template>
        </n-button>
      </n-dropdown>
    </div>

    <div class="skill-description">{{ skill.description }}</div>

    <div class="skill-metrics">
      <span class="metric-pill">
        <strong>{{ skillAgents.length }}</strong>
        {{ t('skills.statsAgents') }}
      </span>
      <span class="metric-pill">
        <strong>{{ skillTags(skill).length }}</strong>
        {{ t('skills.statsTags') }}
      </span>
    </div>

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

    <div class="skill-meta">
      <span class="install-time">{{ installTimeText }}</span>
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
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NTag, NButton, NIcon, NDropdown } from 'naive-ui'
import { EllipsisVertical } from '@vicons/ionicons5'
import type { Skill, Agent } from '@/types'
import { useAgentStore } from '@/stores/agentStore'
import { formatRelativeTime, useI18n } from '@/i18n'

const props = defineProps<{
  skill: Skill
}>()

const emit = defineEmits<{
  click: []
  uninstall: [skillId: string]
}>()

const agentStore = useAgentStore()
const { t } = useI18n()

const dropdownOptions = computed(() => [
  { label: t('common.viewDetails'), key: 'detail' },
  { label: t('common.uninstall'), key: 'uninstall' }
])

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

const displayName = computed(() => {
  const name = typeof props.skill.name === 'string' ? props.skill.name.trim() : ''
  if (name) return name
  const id = typeof props.skill.id === 'string' ? props.skill.id.trim() : ''
  return id || 'Untitled Skill'
})

const skillInitial = computed(() => {
  return displayName.value.slice(0, 1).toUpperCase()
})

function getAgentIcon(agentId: string): string {
  return agentIcons[agentId] || fallbackAgentIcon
}

const installTimeText = computed(() => {
  const relativeTime = formatRelativeTime(props.skill.installedAt)
  return relativeTime ? t('skills.installedRelative', { time: relativeTime }) : ''
})

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
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 176px;
  padding: 16px 18px 14px;
}

.card-topbar {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.skill-avatar {
  width: 40px;
  height: 40px;
  border-radius: 14px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 16px;
  font-weight: 800;
  color: var(--text-primary);
  background: linear-gradient(135deg, rgba(93, 161, 255, 0.34), rgba(131, 226, 186, 0.24));
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.72);
}

.card-copy {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.skill-overline {
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--text-tertiary);
}

.skill-name {
  font-weight: 700;
  font-size: 16px;
  line-height: 1.2;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.card-menu-button {
  border-radius: 999px;
}

.skill-description {
  color: var(--text-secondary);
  font-size: 13px;
  line-height: 1.55;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.skill-metrics {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}

.metric-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.72);
  color: var(--text-secondary);
  font-size: 11px;
  font-weight: 600;
}

.metric-pill strong {
  color: var(--text-primary);
}

.skill-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.skill-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: auto;
}

.install-time {
  line-height: 1.5;
}

.agent-icons {
  display: flex;
  align-items: center;
  gap: 4px;
}

.agent-icon {
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  padding: 3px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.72);
}

.agent-icon :deep(svg) {
  width: 100%;
  height: 100%;
}

.agent-more {
  font-size: 11px;
  color: var(--text-secondary);
  margin-left: 2px;
}
</style>
