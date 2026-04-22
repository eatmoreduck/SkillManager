<template>
  <div class="page-shell agents-view">
    <section class="page-hero">
      <div class="hero-copy">
        <p class="hero-kicker">{{ t('agents.heroKicker') }}</p>
        <h2 class="hero-title">{{ t('agents.title') }}</h2>
        <p class="hero-subtitle">{{ t('agents.heroSubtitle') }}</p>
      </div>

      <div class="hero-stats">
        <div class="hero-stat">
          <span class="hero-stat-value">{{ agentStore.enabledAgents.length }}</span>
          <span class="hero-stat-label">{{ t('agents.statsEnabled') }}</span>
        </div>
        <div class="hero-stat">
          <span class="hero-stat-value">{{ agentStore.installedAgents.length }}</span>
          <span class="hero-stat-label">{{ t('agents.statsInstalled') }}</span>
        </div>
        <div class="hero-stat">
          <span class="hero-stat-value">{{ agentStore.customAgents.length }}</span>
          <span class="hero-stat-label">{{ t('agents.statsCustom') }}</span>
        </div>
      </div>
    </section>

    <section class="glass-panel toolbar-panel">
      <div class="filter-toolbar">
        <div class="filter-toolbar-left">
          <n-button @click="handleRefresh" :loading="agentStore.loading">
            <template #icon><n-icon><RefreshOutline /></n-icon></template>
            {{ t('common.refresh') }}
          </n-button>
          <n-button @click="handleDetectInstalled" :loading="detecting">
            <template #icon><n-icon><SearchOutline /></n-icon></template>
            {{ t('agents.detectInstalled') }}
          </n-button>
        </div>

        <div class="filter-toolbar-right">
          <div class="badge-chip">
            <strong>{{ skillStore.skillCount }}</strong>
            <span>{{ t('agents.connectedSkills') }}</span>
          </div>
          <n-button type="primary" @click="showAddModal = true">
            <template #icon><n-icon><AddOutline /></n-icon></template>
            {{ t('agents.addCustom') }}
          </n-button>
        </div>
      </div>
    </section>

    <div v-if="agentStore.loading" class="state-surface">
      <n-spin size="large" :description="t('common.loading')" />
    </div>

    <div v-else-if="agentStore.agents.length === 0" class="state-surface">
      <n-empty :description="t('agents.empty')">
        <template #extra>
          <n-space>
            <n-button size="small" @click="handleDetectInstalled">
              {{ t('agents.detectInstalled') }}
            </n-button>
            <n-button size="small" type="primary" @click="showAddModal = true">
              {{ t('common.add') }}
            </n-button>
          </n-space>
        </template>
      </n-empty>
    </div>

    <section v-else class="glass-panel content-surface">
      <div class="section-bar">
        <div>
          <h3 class="section-title">{{ t('agents.sectionTitle') }}</h3>
          <p class="section-subtitle">{{ t('agents.sectionSubtitle') }}</p>
        </div>
        <div class="section-meta">{{ agentStore.agents.length }} {{ t('agents.resultsSuffix') }}</div>
      </div>

      <div class="agents-grid">
        <n-card
          v-for="agent in agentStore.agents"
          :key="agent.id"
          class="apple-card agent-card"
          hoverable
        >
          <template #header>
            <div class="agent-header">
              <div class="agent-header-copy">
                <span class="agent-overline">{{ agent.binaryName }}</span>
                <span class="agent-title">{{ agent.name }}</span>
              </div>
            </div>
          </template>

          <template #header-extra>
            <n-switch
              :value="agent.isEnabled"
              @update:value="(val: boolean) => handleToggle(agent.id, val)"
              :loading="togglingId === agent.id"
            />
          </template>

          <div class="agent-body">
            <div class="agent-line">
              <span class="agent-label">{{ t('agents.agentId') }}</span>
              <span class="agent-value">{{ agent.id }}</span>
            </div>
            <div class="agent-line">
              <span class="agent-label">{{ t('agents.binaryName') }}</span>
              <span class="agent-value">{{ agent.binaryName }}</span>
            </div>
            <div class="agent-line agent-line-path">
              <span class="agent-label">{{ t('agents.skillsDir') }}</span>
              <span class="agent-value">{{ agent.skillsDir }}</span>
            </div>
          </div>

          <template #footer>
            <div class="agent-tags">
              <n-tag :type="agent.isInstalled ? 'success' : 'default'" size="small">
                {{ agent.isInstalled ? t('common.installed') : t('common.notInstalled') }}
              </n-tag>
              <n-tag v-if="agent.isCustom" type="warning" size="small">
                {{ t('common.custom') }}
              </n-tag>
              <n-tag type="info" size="small">
                {{ t('agents.skillCount', { count: getSkillCount(agent.id) }) }}
              </n-tag>
            </div>
          </template>

          <template #action>
            <div class="agent-actions">
              <n-popconfirm
                v-if="agent.isCustom"
                @positive-click="handleRemove(agent.id)"
              >
                <template #trigger>
                  <n-button type="error" size="small" :loading="removingId === agent.id">
                    <template #icon><n-icon><TrashOutline /></n-icon></template>
                    {{ t('common.delete') }}
                  </n-button>
                </template>
                {{ t('agents.removeConfirm') }}
              </n-popconfirm>
            </div>
          </template>
        </n-card>
      </div>
    </section>

    <n-modal
      v-model:show="showAddModal"
      preset="dialog"
      :title="t('agents.addCustom')"
      :positive-text="t('common.add')"
      :negative-text="t('common.cancel')"
      :positive-button-props="{ disabled: !isFormValid }"
      :loading="adding"
      @positive-click="handleAdd"
      @close="resetForm"
      @negative-click="resetForm"
    >
      <n-form ref="formRef" :model="formData" label-placement="left" label-width="120">
        <n-form-item :label="t('agents.agentId')" required>
          <n-input
            v-model:value="formData.id"
            :placeholder="t('agents.exampleId')"
          />
        </n-form-item>
        <n-form-item :label="t('common.name')" required>
          <n-input
            v-model:value="formData.name"
            :placeholder="t('agents.exampleName')"
          />
        </n-form-item>
        <n-form-item :label="t('agents.skillsDir')" required>
          <n-input
            v-model:value="formData.skillsDir"
            :placeholder="t('agents.exampleSkillsDir')"
          />
        </n-form-item>
        <n-form-item :label="t('agents.binaryName')" required>
          <n-input
            v-model:value="formData.binaryName"
            :placeholder="t('agents.exampleBinaryName')"
          />
        </n-form-item>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  NButton,
  NIcon,
  NSpin,
  NEmpty,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NCard,
  NSpace,
  NTag,
  NSwitch,
  NPopconfirm,
  useMessage
} from 'naive-ui'
import { RefreshOutline, AddOutline, TrashOutline, SearchOutline } from '@vicons/ionicons5'
import { useAgentStore } from '@/stores/agentStore'
import { useSkillStore } from '@/stores/skillStore'
import type { Agent } from '@/types'
import { useI18n } from '@/i18n'

const message = useMessage()
const { t } = useI18n()

const agentStore = useAgentStore()
const skillStore = useSkillStore()

const showAddModal = ref(false)
const detecting = ref(false)
const adding = ref(false)
const togglingId = ref<string | null>(null)
const removingId = ref<string | null>(null)

const formRef = ref()
const formData = ref({
  id: '',
  name: '',
  skillsDir: '',
  binaryName: ''
})

const isFormValid = computed(() => {
  return formData.value.id.trim() &&
    formData.value.name.trim() &&
    formData.value.skillsDir.trim() &&
    formData.value.binaryName.trim()
})

function getSkillCount(agentId: string): number {
  return skillStore.skills.filter(skill => skill.agents.includes(agentId)).length
}

async function handleRefresh() {
  await Promise.all([
    agentStore.loadAgents(),
    skillStore.loadSkills()
  ])
}

async function handleDetectInstalled() {
  detecting.value = true
  try {
    await agentStore.detectInstalled()
    message.success(t('agents.detectSuccess'))
  } catch (error) {
    message.error(t('agents.detectFailed', { error: (error as Error).message }))
  } finally {
    detecting.value = false
  }
}

async function handleToggle(id: string, enabled: boolean) {
  togglingId.value = id
  try {
    await agentStore.toggleAgent(id, enabled)
    message.success(enabled ? t('common.enabled') : t('common.disabled'))
  } catch (error) {
    message.error(t('agents.toggleFailed', { error: (error as Error).message }))
  } finally {
    togglingId.value = null
  }
}

async function handleRemove(id: string) {
  removingId.value = id
  try {
    await agentStore.removeAgent(id)
    message.success(t('common.deleteSuccess'))
  } catch (error) {
    message.error(t('common.deleteFailed', { error: (error as Error).message }))
  } finally {
    removingId.value = null
  }
}

async function handleAdd() {
  if (!isFormValid.value) return false

  adding.value = true
  try {
    const agent: Agent = {
      id: formData.value.id.trim(),
      name: formData.value.name.trim(),
      skillsDir: formData.value.skillsDir.trim(),
      binaryName: formData.value.binaryName.trim(),
      priorityPaths: [],
      isInstalled: false,
      isEnabled: true,
      isCustom: true
    }
    await agentStore.addCustomAgent(agent)
    message.success(t('common.addSuccess'))
    resetForm()
    return true
  } catch (error) {
    message.error(t('common.addFailed', { error: (error as Error).message }))
    return false
  } finally {
    adding.value = false
  }
}

function resetForm() {
  formData.value = {
    id: '',
    name: '',
    skillsDir: '',
    binaryName: ''
  }
}

onMounted(async () => {
  await Promise.all([
    agentStore.loadAgents(),
    skillStore.loadSkills()
  ])
})
</script>

<style scoped>
.agents-view {
  padding: 22px;
}

.agents-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 18px;
  overflow: auto;
  flex: 1;
  min-height: 0;
  align-content: start;
  padding-right: 6px;
}

.agent-header-copy {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.agent-overline {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--text-tertiary);
}

.agent-title {
  font-size: 18px;
  line-height: 1.2;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.agent-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.agent-line {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.agent-label {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  color: var(--text-tertiary);
}

.agent-value {
  color: var(--text-secondary);
  line-height: 1.6;
  word-break: break-word;
}

.agent-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.agent-actions {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .agents-view {
    padding: 16px;
  }
}
</style>
