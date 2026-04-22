<template>
  <div class="page-shell skill-detail-view">
    <div v-if="loading" class="state-surface">
      <n-spin size="large" :description="t('common.loading')" />
    </div>

    <template v-else-if="skill">
      <section class="page-hero detail-hero">
        <div class="hero-copy">
          <n-button quaternary class="back-button" @click="handleBack">
            <template #icon><n-icon><ArrowBackOutline /></n-icon></template>
            {{ t('detail.backToList') }}
          </n-button>
          <p class="hero-kicker">{{ t('detail.heroKicker') }}</p>
          <h2 class="hero-title">{{ skill.name }}</h2>
          <p class="hero-subtitle">{{ skill.description || t('detail.heroSubtitle') }}</p>
          <div class="detail-hero-tags">
            <n-tag v-if="skill.version" type="info" size="small">v{{ skill.version }}</n-tag>
            <n-tag
              v-for="tag in normalizedTags(skill.tags).slice(0, 4)"
              :key="tag"
              size="small"
            >
              {{ tag }}
            </n-tag>
          </div>
        </div>

        <div class="detail-hero-actions">
          <n-space>
            <n-button @click="isEditing = true">
              <template #icon><n-icon><CreateOutline /></n-icon></template>
              {{ t('common.edit') }}
            </n-button>
            <n-button @click="handleUpdate" :loading="updating">
              <template #icon><n-icon><SyncOutline /></n-icon></template>
              {{ t('common.update') }}
            </n-button>
            <n-button type="error" @click="handleUninstall">
              <template #icon><n-icon><TrashOutline /></n-icon></template>
              {{ t('common.uninstall') }}
            </n-button>
          </n-space>
        </div>
      </section>

      <section v-if="isEditing" class="glass-panel editor-shell">
        <div class="section-bar">
          <div>
            <h3 class="section-title">{{ t('detail.editingTitle') }}</h3>
            <p class="section-subtitle">{{ t('detail.editingSubtitle') }}</p>
          </div>
        </div>
        <div class="editor-toolbar">
          <n-space>
            <n-button size="small" @click="insertMarkdown('**', '**')">
              <template #icon><n-icon><TextOutline /></n-icon></template>
              B
            </n-button>
            <n-button size="small" @click="insertMarkdown('*', '*')">
              <template #icon><n-icon><TextOutline /></n-icon></template>
              I
            </n-button>
            <n-button size="small" @click="insertMarkdown('`', '`')">
              <template #icon><n-icon><CodeSlash /></n-icon></template>
            </n-button>
          </n-space>
          <n-space>
            <n-button @click="isEditing = false">{{ t('common.cancel') }}</n-button>
            <n-button type="primary" :loading="saving" @click="handleSave">{{ t('common.save') }}</n-button>
          </n-space>
        </div>
        <textarea
          ref="editorRef"
          v-model="editContent"
          class="markdown-editor"
        />
      </section>

      <div v-else class="detail-layout">
        <section class="glass-panel detail-meta-panel">
          <div class="section-bar">
            <div>
              <h3 class="section-title">{{ t('detail.overviewTitle') }}</h3>
              <p class="section-subtitle">{{ t('detail.overviewSubtitle') }}</p>
            </div>
          </div>

          <n-descriptions label-placement="top" :column="2">
            <n-descriptions-item :label="t('common.author')">{{ skill.author }}</n-descriptions-item>
            <n-descriptions-item :label="t('common.version')">{{ skill.version }}</n-descriptions-item>
            <n-descriptions-item :label="t('common.tags')" :span="2">
              <n-space>
                <n-tag v-for="tag in normalizedTags(skill.tags)" :key="tag" size="small">{{ tag }}</n-tag>
              </n-space>
            </n-descriptions-item>
            <n-descriptions-item :label="t('detail.installedAt')">{{ formatDate(skill.installedAt) }}</n-descriptions-item>
            <n-descriptions-item :label="t('detail.updatedAt')">{{ formatDate(skill.updatedAt) }}</n-descriptions-item>
            <n-descriptions-item :label="t('detail.sourceUrl')" :span="2">{{ skill.sourceUrl || '-' }}</n-descriptions-item>
            <n-descriptions-item :label="t('detail.localPath')" :span="2">{{ skill.localPath || '-' }}</n-descriptions-item>
          </n-descriptions>

          <div class="detail-panel-actions">
            <div class="agent-stack">
              <div class="agent-stack-copy">
                <span class="agent-stack-title">{{ t('detail.assignedAgents') }}</span>
                <div class="agent-chip-list">
                  <n-tag v-for="agentId in skill.agents" :key="agentId" type="info">
                    {{ getAgentName(agentId) }}
                  </n-tag>
                </div>
              </div>
              <n-button size="small" @click="showAgentModal = true">
                {{ t('common.manage') }}
              </n-button>
            </div>
          </div>
        </section>

        <section class="glass-panel detail-content-panel">
          <div class="section-bar">
            <div>
              <h3 class="section-title">{{ t('detail.contentTitle') }}</h3>
              <p class="section-subtitle">{{ t('detail.contentSubtitle') }}</p>
            </div>
          </div>
          <div class="markdown-body" v-html="renderedContent" />
        </section>
      </div>

      <n-modal
        v-model:show="showAgentModal"
        preset="card"
        :title="t('detail.manageAgents')"
        style="width: 460px; max-width: 92vw"
      >
        <n-checkbox-group v-model:value="editAgents">
          <n-space vertical>
            <n-checkbox
              v-for="agent in agentStore.agents"
              :key="agent.id"
              :value="agent.id"
              :label="agent.name"
              :disabled="!agent.isEnabled"
            />
          </n-space>
        </n-checkbox-group>
        <template #footer>
          <n-space justify="end">
            <n-button @click="showAgentModal = false">{{ t('common.cancel') }}</n-button>
            <n-button type="primary" :loading="savingAgents" @click="handleSaveAgents">
              {{ t('common.save') }}
            </n-button>
          </n-space>
        </template>
      </n-modal>
    </template>

    <div v-else class="state-surface">
      <n-empty :description="t('detail.notFound')">
        <template #extra>
          <n-button @click="$router.push('/skills')">{{ t('detail.backToList') }}</n-button>
        </template>
      </n-empty>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NDescriptions,
  NDescriptionsItem,
  NSpace,
  NButton,
  NIcon,
  NTag,
  NSpin,
  NEmpty,
  NModal,
  NCheckboxGroup,
  NCheckbox,
  useMessage,
  useDialog
} from 'naive-ui'
import {
  ArrowBackOutline,
  CreateOutline,
  SyncOutline,
  TrashOutline,
  TextOutline,
  CodeSlash
} from '@vicons/ionicons5'
import { marked } from 'marked'
import { useSkillStore } from '@/stores/skillStore'
import { useAgentStore } from '@/stores/agentStore'
import type { Skill } from '@/types'
import { formatDateTime, useI18n } from '@/i18n'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const { t } = useI18n()

const skillStore = useSkillStore()
const agentStore = useAgentStore()

const skill = ref<Skill | null>(null)
const loading = ref(true)
const isEditing = ref(false)
const editContent = ref('')
const editAgents = ref<string[]>([])
const saving = ref(false)
const savingAgents = ref(false)
const updating = ref(false)
const showAgentModal = ref(false)
const editorRef = ref<HTMLTextAreaElement | null>(null)

const skillId = computed(() => decodeURIComponent(route.params.id as string))

const renderedContent = computed(() => {
  return String(marked.parse(skill.value?.content || ''))
})

function formatDate(date: string): string {
  if (!date) return '-'
  return formatDateTime(date)
}

function normalizedTags(tags: string[] | null | undefined): string[] {
  return Array.isArray(tags) ? tags : []
}

function getAgentName(agentId: string): string {
  return agentStore.getAgentById(agentId)?.name || agentId
}

function handleBack() {
  router.push('/skills')
}

async function loadSkill() {
  loading.value = true
  try {
    skill.value = await skillStore.getSkillDetail(skillId.value)
    editAgents.value = [...(skill.value?.agents || [])]
  } catch (error) {
    message.error(t('common.loadFailed', { error: (error as Error).message }))
  } finally {
    loading.value = false
  }
}

function insertMarkdown(before: string, after: string) {
  if (!editorRef.value) return
  const textarea = editorRef.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = editContent.value.substring(start, end)
  const newText = before + selectedText + after
  editContent.value = editContent.value.substring(0, start) + newText + editContent.value.substring(end)
  nextTick(() => {
    textarea.focus()
    textarea.setSelectionRange(start + before.length, start + before.length + selectedText.length)
  })
}

async function handleSave() {
  if (!skill.value) return
  saving.value = true
  try {
    await skillStore.updateSkillContent(skill.value.id, editContent.value)
    skill.value.content = editContent.value
    isEditing.value = false
    message.success(t('common.saveSuccess'))
  } catch (error) {
    message.error(t('common.saveFailed', { error: (error as Error).message }))
  } finally {
    saving.value = false
  }
}

async function handleUpdate() {
  if (!skill.value) return
  updating.value = true
  try {
    const updated = await skillStore.updateSkill(skill.value.id)
    skill.value = updated
    message.success(t('common.updateSuccess'))
  } catch (error) {
    message.error(t('common.updateFailed', { error: (error as Error).message }))
  } finally {
    updating.value = false
  }
}

function handleUninstall() {
  if (!skill.value) return
  dialog.warning({
    title: t('detail.confirmUninstallTitle'),
    content: t('detail.confirmUninstallContent', { name: skill.value.name }),
    positiveText: t('common.uninstall'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await skillStore.uninstallSkill(skill.value!.id)
        message.success(t('common.uninstallSuccess'))
        router.push('/skills')
      } catch (error) {
        message.error(t('common.uninstallFailed', { error: (error as Error).message }))
      }
    }
  })
}

async function handleSaveAgents() {
  if (!skill.value) return
  savingAgents.value = true
  try {
    await skillStore.assignAgents(skill.value.id, editAgents.value)
    skill.value.agents = [...editAgents.value]
    showAgentModal.value = false
    message.success(t('common.saveSuccess'))
  } catch (error) {
    message.error(t('common.saveFailed', { error: (error as Error).message }))
  } finally {
    savingAgents.value = false
  }
}

onMounted(async () => {
  await agentStore.loadAgents()
  await loadSkill()
  if (skill.value) {
    editContent.value = skill.value.content
  }
})
</script>

<style scoped>
.skill-detail-view {
  padding: 22px;
}

.back-button {
  margin-bottom: 14px;
  border-radius: 999px;
}

.detail-hero {
  align-items: stretch;
}

.detail-hero-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 16px;
}

.detail-hero-actions {
  display: flex;
  align-items: flex-start;
}

.detail-layout {
  display: grid;
  grid-template-columns: minmax(320px, 0.95fr) minmax(0, 1.25fr);
  gap: 18px;
  min-height: 0;
  flex: 1;
}

.detail-meta-panel,
.detail-content-panel,
.editor-shell {
  min-height: 0;
}

.detail-meta-panel {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.detail-content-panel {
  display: flex;
  flex-direction: column;
  gap: 18px;
  overflow: hidden;
}

.detail-panel-actions {
  margin-top: auto;
}

.agent-stack {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 16px 18px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.62);
  border: 1px solid rgba(255, 255, 255, 0.7);
}

.agent-stack-copy {
  min-width: 0;
}

.agent-stack-title {
  display: block;
  margin-bottom: 10px;
  font-size: 13px;
  font-weight: 700;
  color: var(--text-primary);
}

.agent-chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.markdown-body {
  overflow: auto;
  flex: 1;
  min-height: 0;
  padding-right: 4px;
  color: var(--text-secondary);
  line-height: 1.8;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3) {
  color: var(--text-primary);
  letter-spacing: -0.03em;
  margin: 1.2em 0 0.6em;
}

.markdown-body :deep(p),
.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  margin-bottom: 1em;
}

.markdown-body :deep(pre) {
  padding: 16px;
  border-radius: 20px;
  background: rgba(19, 44, 72, 0.92);
  color: #eef5ff;
  overflow: auto;
}

.markdown-body :deep(code) {
  font-family: 'JetBrains Mono', 'SFMono-Regular', Consolas, monospace;
}

.editor-shell {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.editor-toolbar {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  padding: 14px 16px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.68);
  border: 1px solid rgba(255, 255, 255, 0.72);
}

.markdown-editor {
  flex: 1 1 auto;
  width: 100%;
  min-height: 420px;
  padding: 18px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  background: rgba(255, 255, 255, 0.74);
  color: var(--text-primary);
  font-family: 'JetBrains Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: 14px;
  line-height: 1.7;
  resize: vertical;
  outline: none;
}

@media (max-width: 1024px) {
  .detail-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .skill-detail-view {
    padding: 16px;
  }

  .detail-hero-actions,
  .agent-stack {
    width: 100%;
  }

  .detail-hero-actions :deep(.n-space) {
    width: 100%;
    justify-content: flex-start;
    flex-wrap: wrap;
  }

  .agent-stack {
    flex-direction: column;
  }

  .editor-toolbar {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
