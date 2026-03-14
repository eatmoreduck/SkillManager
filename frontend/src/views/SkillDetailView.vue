<template>
  <div class="skill-detail-view">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <n-spin size="large" />
    </div>

    <!-- Skill 详情 -->
    <template v-else-if="skill">
      <div class="detail-header">
        <n-page-header @back="handleBack">
          <template #title>{{ skill.name }}</template>
          <template #subtitle>{{ skill.description }}</template>
          <template #extra>
            <n-space>
              <n-button @click="isEditing = true">
                <template #icon><n-icon><CreateOutline /></n-icon></template>
                编辑
              </n-button>
              <n-button @click="handleUpdate" :loading="updating">
                <template #icon><n-icon><SyncOutline /></n-icon></template>
                更新
              </n-button>
              <n-button type="error" @click="handleUninstall">
                <template #icon><n-icon><TrashOutline /></n-icon></template>
                卸载
              </n-button>
            </n-space>
          </template>
        </n-page-header>
      </div>

      <!-- 编辑模式 -->
      <div v-if="isEditing" class="editor-container">
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
            <n-button @click="isEditing = false">取消</n-button>
            <n-button type="primary" :loading="saving" @click="handleSave">保存</n-button>
          </n-space>
        </div>
        <textarea
          ref="editorRef"
          v-model="editContent"
          class="markdown-editor"
        />
      </div>

      <!-- 查看模式 -->
      <div v-else class="detail-content">
        <n-card class="info-card">
          <n-descriptions label-placement="left" :column="2">
            <n-descriptions-item label="作者">{{ skill.author }}</n-descriptions-item>
            <n-descriptions-item label="版本">{{ skill.version }}</n-descriptions-item>
            <n-descriptions-item label="标签">
              <n-space>
                <n-tag v-for="tag in normalizedTags(skill.tags)" :key="tag" size="small">{{ tag }}</n-tag>
              </n-space>
            </n-descriptions-item>
            <n-descriptions-item label="安装时间">{{ formatDate(skill.installedAt) }}</n-descriptions-item>
            <n-descriptions-item label="已分配 Agent" :span="2">
              <n-space>
                <n-tag v-for="agentId in skill.agents" :key="agentId" type="info">
                  {{ getAgentName(agentId) }}
                </n-tag>
                <n-button size="small" @click="showAgentModal = true">
                  管理
                </n-button>
              </n-space>
            </n-descriptions-item>
          </n-descriptions>
        </n-card>

        <n-card title="SKILL.md 内容" class="content-card">
          <div class="markdown-body" v-html="renderedContent" />
        </n-card>
      </div>

      <!-- Agent 管理弹窗 -->
      <n-modal
        v-model:show="showAgentModal"
        preset="card"
        title="管理 Agent 分配"
        style="width: 400px"
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
            <n-button @click="showAgentModal = false">取消</n-button>
            <n-button type="primary" :loading="savingAgents" @click="handleSaveAgents">
              保存
            </n-button>
          </n-space>
        </template>
      </n-modal>
    </template>

    <!-- 错误状态 -->
    <div v-else class="error-container">
      <n-empty description="Skill 不存在">
        <template #extra>
          <n-button @click="$router.push('/skills')">返回列表</n-button>
        </template>
      </n-empty>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NPageHeader,
  NCard,
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

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()

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
  return marked.parse(skill.value?.content || '')
})

function formatDate(date: string): string {
  if (!date) return '-'
  return new Date(date).toLocaleString()
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
    message.error('加载失败')
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
    message.success('保存成功')
  } catch (error) {
    message.error('保存失败: ' + (error as Error).message)
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
    message.success('更新成功')
  } catch (error) {
    message.error('更新失败: ' + (error as Error).message)
  } finally {
    updating.value = false
  }
}

function handleUninstall() {
  if (!skill.value) return
  dialog.warning({
    title: '确认卸载',
    content: `确定要卸载 "${skill.value.name}" 吗？此操作不可撤销。`,
    positiveText: '卸载',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await skillStore.uninstallSkill(skill.value!.id)
        message.success('卸载成功')
        router.push('/skills')
      } catch (error) {
        message.error('卸载失败')
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
    message.success('保存成功')
  } catch (error) {
    message.error('保存失败')
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
  padding: 16px;
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.loading-container,
.error-container {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.detail-header {
  margin-bottom: 16px;
  flex-shrink: 0;
}

.detail-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.info-card {
  margin-bottom: 16px;
}

.content-card {
  flex: 1;
}

.markdown-body {
  line-height: 1.6;
}

.editor-container {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.editor-toolbar {
  display: flex;
  justify-content: space-between;
  padding: 8px;
  border: 1px solid var(--n-border-color);
  border-bottom: none;
  border-radius: 4px 4px 0 0;
  background: var(--n-color-hover);
}

.markdown-editor {
  flex: 1 1 auto;
  width: 100%;
  min-height: 0;
  height: 100%;
  padding: 16px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 14px;
  line-height: 1.6;
  border: 1px solid var(--n-border-color);
  border-radius: 0 0 4px 4px;
  resize: none;
  outline: none;
  box-sizing: border-box;
}
</style>
