<template>
  <div class="registry-view">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <n-select
          v-model:value="selectedRegistryId"
          :options="registryOptions"
          placeholder="选择 Registry"
          style="width: 200px"
          :loading="registryStore.loading"
          @update:value="handleRegistryChange"
        />
        <n-input
          v-model:value="searchQuery"
          placeholder="搜索 Skills..."
          clearable
          style="width: 300px"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <n-icon><SearchOutline /></n-icon>
          </template>
        </n-input>
        <n-button @click="handleSearch" :loading="searching">
          搜索
        </n-button>
      </div>
      <div class="toolbar-right">
        <n-button @click="handleRefresh" :loading="browsing">
          <template #icon><n-icon><RefreshOutline /></n-icon></template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="registryStore.loading || browsing" class="loading-container">
      <n-spin size="large" description="加载中..." />
    </div>

    <!-- 空状态 -->
    <div v-else-if="!registryStore.hasRegistries" class="empty-container">
      <n-empty description="暂无可用的 Registry">
        <template #extra>
          <n-button size="small" @click="$router.push('/settings')">
            前往设置添加 Registry
          </n-button>
        </template>
      </n-empty>
    </div>

    <div v-else-if="displaySkills.length === 0 && !searching" class="empty-container">
      <n-empty description="当前 Registry 暂无 Skills" />
    </div>

    <div v-else-if="displaySkills.length === 0 && searching" class="empty-container">
      <n-empty description="未找到匹配的 Skills">
        <template #extra>
          <n-button size="small" @click="clearSearch">
            清除搜索
          </n-button>
        </template>
      </n-empty>
    </div>

    <!-- Skills 卡片网格 -->
    <div v-else class="skills-grid">
      <n-card
        v-for="skill in displaySkills"
        :key="skill.id"
        class="registry-skill-card"
        hoverable
        @click="handleSkillClick(skill)"
      >
        <template #header>
          <div class="card-header">
            <span class="skill-name">{{ skill.name }}</span>
            <n-tag v-if="skill.category" size="small" type="info">
              {{ skill.category }}
            </n-tag>
          </div>
        </template>

        <div class="skill-description">{{ skill.description }}</div>

        <div class="skill-tags">
          <n-tag
            v-for="tag in skill.tags.slice(0, 3)"
            :key="tag"
            size="small"
            :bordered="false"
          >
            {{ tag }}
          </n-tag>
          <n-tag v-if="skill.tags.length > 3" size="small" :bordered="false">
            +{{ skill.tags.length - 3 }}
          </n-tag>
        </div>

        <template #footer>
          <div class="skill-meta">
            <span class="author">
              <n-icon><PersonOutline /></n-icon>
              {{ skill.author }}
            </span>
            <span v-if="skill.stars" class="stars">
              <n-icon><StarOutline /></n-icon>
              {{ skill.stars }}
            </span>
          </div>
        </template>

        <template #action>
          <n-button
            type="primary"
            size="small"
            @click.stop="handleInstallClick(skill)"
          >
            <template #icon><n-icon><DownloadOutline /></n-icon></template>
            安装
          </n-button>
        </template>
      </n-card>
    </div>

    <!-- 安装弹窗：选择 Agent -->
    <n-modal
      v-model:show="showInstallModal"
      preset="dialog"
      title="安装 Skill"
      positive-text="安装"
      negative-text="取消"
      :positive-button-props="{ disabled: selectedAgents.length === 0 }"
      :loading="installing"
      @positive-click="handleInstallConfirm"
      @negative-click="closeInstallModal"
    >
      <div v-if="selectedSkill" class="install-modal-content">
        <div class="skill-preview">
          <h4>{{ selectedSkill.name }}</h4>
          <p>{{ selectedSkill.description }}</p>
        </div>
        <n-divider />
        <n-form-item label="选择要分配的 Agent">
          <n-checkbox-group v-model:value="selectedAgents">
            <n-space vertical>
              <n-checkbox
                v-for="agent in agentStore.enabledAgents"
                :key="agent.id"
                :value="agent.id"
                :label="agent.name"
              />
            </n-space>
          </n-checkbox-group>
        </n-form-item>
      </div>
    </n-modal>

    <!-- Skill 详情弹窗 -->
    <n-modal
      v-model:show="showDetailModal"
      preset="card"
      :title="selectedSkill?.name || 'Skill 详情'"
      style="width: 600px; max-width: 90vw"
    >
      <div v-if="selectedSkill" class="skill-detail">
        <n-descriptions label-placement="left" :column="1" bordered>
          <n-descriptions-item label="名称">
            {{ selectedSkill.name }}
          </n-descriptions-item>
          <n-descriptions-item label="作者">
            {{ selectedSkill.author }}
          </n-descriptions-item>
          <n-descriptions-item v-if="selectedSkill.category" label="分类">
            {{ selectedSkill.category }}
          </n-descriptions-item>
          <n-descriptions-item v-if="selectedSkill.stars" label="热度">
            {{ selectedSkill.stars }}
          </n-descriptions-item>
          <n-descriptions-item label="描述">
            {{ selectedSkill.description }}
          </n-descriptions-item>
          <n-descriptions-item label="标签">
            <n-space>
              <n-tag
                v-for="tag in selectedSkill.tags"
                :key="tag"
                size="small"
              >
                {{ tag }}
              </n-tag>
            </n-space>
          </n-descriptions-item>
          <n-descriptions-item label="安装 URL">
            <n-ellipsis style="max-width: 400px">
              {{ selectedSkill.installUrl }}
            </n-ellipsis>
          </n-descriptions-item>
        </n-descriptions>

        <div class="detail-actions">
          <n-button type="primary" @click="handleInstallClick(selectedSkill)">
            <template #icon><n-icon><DownloadOutline /></n-icon></template>
            安装此 Skill
          </n-button>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import {
  NSelect,
  NInput,
  NButton,
  NIcon,
  NSpin,
  NEmpty,
  NCard,
  NTag,
  NModal,
  NCheckboxGroup,
  NCheckbox,
  NSpace,
  NDescriptions,
  NDescriptionsItem,
  NEllipsis,
  NDivider,
  NFormItem,
  useMessage
} from 'naive-ui'
import {
  SearchOutline,
  RefreshOutline,
  DownloadOutline,
  PersonOutline,
  StarOutline
} from '@vicons/ionicons5'
import { useRegistryStore } from '@/stores/registryStore'
import { useAgentStore } from '@/stores/agentStore'
import type { RegistrySkill } from '@/types'

const message = useMessage()

const registryStore = useRegistryStore()
const agentStore = useAgentStore()

// 状态
const selectedRegistryId = ref<string | null>(null)
const searchQuery = ref('')
const browsing = ref(false)
const searching = ref(false)
const installing = ref(false)
const showInstallModal = ref(false)
const showDetailModal = ref(false)
const selectedSkill = ref<RegistrySkill | null>(null)
const selectedAgents = ref<string[]>([])

// 计算属性
const registryOptions = computed(() =>
  registryStore.registries.map(r => ({
    label: r.name + (r.isDefault ? ' (默认)' : ''),
    value: r.id
  }))
)

const displaySkills = computed(() => {
  if (searchQuery.value && registryStore.searchResults.length > 0) {
    return registryStore.searchResults
  }
  if (searchQuery.value && registryStore.searchResults.length === 0) {
    return []
  }
  return registryStore.browseResults
})

// 方法
async function handleRegistryChange(registryId: string) {
  registryStore.switchRegistry(registryId)
  await loadBrowseResults()
}

async function handleSearch() {
  if (!searchQuery.value.trim()) {
    clearSearch()
    return
  }
  searching.value = true
  try {
    await registryStore.search(searchQuery.value)
  } finally {
    searching.value = false
  }
}

function clearSearch() {
  searchQuery.value = ''
  registryStore.clearSearch()
}

async function handleRefresh() {
  await loadBrowseResults()
}

async function loadBrowseResults() {
  if (!selectedRegistryId.value) return
  browsing.value = true
  try {
    await registryStore.browse(selectedRegistryId.value)
  } finally {
    browsing.value = false
  }
}

function handleSkillClick(skill: RegistrySkill) {
  selectedSkill.value = skill
  showDetailModal.value = true
}

function handleInstallClick(skill: RegistrySkill) {
  selectedSkill.value = skill
  showDetailModal.value = false
  showInstallModal.value = true
  // 默认选中所有启用的 agent
  selectedAgents.value = agentStore.enabledAgents.map(a => a.id)
}

function closeInstallModal() {
  showInstallModal.value = false
  selectedSkill.value = null
  selectedAgents.value = []
}

async function handleInstallConfirm() {
  if (!selectedSkill.value || selectedAgents.value.length === 0) return

  installing.value = true
  try {
    await window.go.main.App.SkillBinding.Install(
      selectedSkill.value.installUrl,
      selectedAgents.value
    )
    message.success(`Skill "${selectedSkill.value.name}" 安装成功`)
    closeInstallModal()
  } catch (error) {
    message.error('安装失败: ' + (error as Error).message)
  } finally {
    installing.value = false
  }
}

// 初始化
onMounted(async () => {
  await Promise.all([
    registryStore.loadRegistries(),
    agentStore.loadAgents()
  ])

  // 设置默认选中的 registry
  if (registryStore.currentRegistry) {
    selectedRegistryId.value = registryStore.currentRegistry.id
    await loadBrowseResults()
  } else if (registryStore.registries.length > 0) {
    selectedRegistryId.value = registryStore.registries[0].id
    await loadBrowseResults()
  }
})

// 监听 registry 变化
watch(() => registryStore.currentRegistry, (newRegistry) => {
  if (newRegistry) {
    selectedRegistryId.value = newRegistry.id
  }
})
</script>

<style scoped>
.registry-view {
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
  gap: 12px;
  align-items: center;
}

.toolbar-right {
  display: flex;
  gap: 8px;
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
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
  overflow-y: auto;
  flex: 1;
}

.registry-skill-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.registry-skill-card:hover {
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
  margin-bottom: 8px;
}

.skill-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--n-text-color-3);
}

.skill-meta .author,
.skill-meta .stars {
  display: flex;
  align-items: center;
  gap: 4px;
}

.install-modal-content {
  padding: 8px 0;
}

.skill-preview h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
}

.skill-preview p {
  margin: 0;
  color: var(--n-text-color-2);
  font-size: 13px;
}

.skill-detail {
  padding: 8px 0;
}

.detail-actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
