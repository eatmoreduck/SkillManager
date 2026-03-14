<template>
  <div class="skills-view">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <n-input
          v-model:value="searchQuery"
          placeholder="搜索 Skills..."
          clearable
          style="width: 300px"
        >
          <template #prefix>
            <n-icon><SearchOutline /></n-icon>
          </template>
        </n-input>
      </div>
      <n-button type="primary" @click="showInstallModal = true">
        <template #icon><n-icon><AddOutline /></n-icon></template>
        从 URL 安装
      </n-button>
    </div>

    <n-alert
      v-if="skillStore.error"
      type="error"
      title="加载 Skills 失败"
      class="error-alert"
    >
      {{ skillStore.error }}
    </n-alert>

    <!-- 加载状态 -->
    <div v-if="skillStore.loading" class="loading-container">
      <n-spin size="large" description="加载中..." />
    </div>

    <!-- 空状态 -->
    <div v-else-if="filteredSkills.length === 0" class="empty-container">
      <n-empty description="暂无已安装的 Skills">
        <template #extra>
          <n-button size="small" @click="$router.push('/registry')">
            浏览 Registry
          </n-button>
        </template>
      </n-empty>
    </div>

    <!-- Skill 列表 -->
    <div v-else class="skills-grid">
      <SkillCard
        v-for="skill in filteredSkills"
        :key="skill.id"
        :skill="skill"
        @click="handleSkillClick(skill)"
        @uninstall="handleUninstall"
      />
    </div>

    <!-- 从 URL 安装弹窗 -->
    <n-modal
      v-model:show="showInstallModal"
      preset="dialog"
      title="从 URL 安装 Skill"
      positive-text="安装"
      negative-text="取消"
      :positive-button-props="{ disabled: !installURL || selectedAgents.length === 0 }"
      :loading="installing"
      @positive-click="handleInstall"
    >
      <n-form label-placement="left" label-width="100">
        <n-form-item label="Git URL">
          <n-input
            v-model:value="installURL"
            placeholder="https://github.com/user/skill-name.git"
          />
        </n-form-item>
        <n-form-item label="选择 Agent">
          <n-checkbox-group v-model:value="selectedAgents">
            <n-space>
              <n-checkbox
                v-for="agent in agentStore.enabledAgents"
                :key="agent.id"
                :value="agent.id"
                :label="agent.name"
              />
            </n-space>
          </n-checkbox-group>
        </n-form-item>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NInput,
  NButton,
  NIcon,
  NSpin,
  NEmpty,
  NModal,
  NForm,
  NFormItem,
  NCheckboxGroup,
  NCheckbox,
  NSpace,
  NAlert,
  useMessage,
  useDialog
} from 'naive-ui'
import { SearchOutline, AddOutline } from '@vicons/ionicons5'
import { useSkillStore } from '@/stores/skillStore'
import { useAgentStore } from '@/stores/agentStore'
import SkillCard from '@/components/SkillCard.vue'
import type { Skill } from '@/types'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const skillStore = useSkillStore()
const agentStore = useAgentStore()

const searchQuery = ref('')
const showInstallModal = ref(false)
const installURL = ref('')
const selectedAgents = ref<string[]>([])
const installing = ref(false)

function normalizedTags(tags: string[] | null | undefined): string[] {
  return Array.isArray(tags) ? tags : []
}

const filteredSkills = computed(() => {
  if (!searchQuery.value) return skillStore.skills
  const query = searchQuery.value.toLowerCase()
  return skillStore.skills.filter(skill =>
    skill.name.toLowerCase().includes(query) ||
    skill.description.toLowerCase().includes(query) ||
    normalizedTags(skill.tags).some(tag => tag.toLowerCase().includes(query))
  )
})

function handleSkillClick(skill: Skill) {
  router.push(`/skills/${encodeURIComponent(skill.id)}`)
}

async function handleInstall() {
  if (!installURL.value || selectedAgents.value.length === 0) return

  installing.value = true
  try {
    await skillStore.installSkill(installURL.value, selectedAgents.value)
    message.success('安装成功')
    showInstallModal.value = false
    installURL.value = ''
    selectedAgents.value = []
  } catch (error) {
    message.error('安装失败: ' + (error as Error).message)
  } finally {
    installing.value = false
  }
}

function handleUninstall(skillId: string) {
  dialog.warning({
    title: '确认卸载',
    content: '确定要卸载此 Skill 吗？此操作不可撤销。',
    positiveText: '卸载',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await skillStore.uninstallSkill(skillId)
        message.success('卸载成功')
      } catch (error) {
        message.error('卸载失败: ' + (error as Error).message)
      }
    }
  })
}

onMounted(async () => {
  await Promise.all([
    skillStore.loadSkills(),
    agentStore.loadAgents()
  ])
  // 默认选中所有启用的 agent
  selectedAgents.value = agentStore.enabledAgents.map(a => a.id)
})
</script>

<style scoped>
.skills-view {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-shrink: 0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.error-alert {
  margin-bottom: 16px;
  flex-shrink: 0;
}

.loading-container,
.empty-container {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  overflow-y: auto;
  flex: 1;
}
</style>
