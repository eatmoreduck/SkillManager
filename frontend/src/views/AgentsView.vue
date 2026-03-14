<template>
  <div class="agents-view">
    <!-- 工具栏 -->
    <div class="toolbar">
      <n-space>
        <n-button @click="handleRefresh" :loading="agentStore.loading">
          <template #icon><n-icon><RefreshOutline /></n-icon></template>
          刷新
        </n-button>
        <n-button @click="handleDetectInstalled" :loading="detecting">
          <template #icon><n-icon><SearchOutline /></n-icon></template>
          检测已安装
        </n-button>
      </n-space>
      <n-button type="primary" @click="showAddModal = true">
        <template #icon><n-icon><AddOutline /></n-icon></template>
        添加自定义 Agent
      </n-button>
    </div>

    <!-- 加载状态 -->
    <div v-if="agentStore.loading" class="loading-container">
      <n-spin size="large" description="加载中..." />
    </div>

    <!-- 空状态 -->
    <div v-else-if="agentStore.agents.length === 0" class="empty-container">
      <n-empty description="暂无 Agent">
        <template #extra>
          <n-space>
            <n-button size="small" @click="handleDetectInstalled">
              检测已安装
            </n-button>
            <n-button size="small" type="primary" @click="showAddModal = true">
              添加自定义
            </n-button>
          </n-space>
        </template>
      </n-empty>
    </div>

    <!-- Agent 列表 -->
    <div v-else class="agents-grid">
      <n-card
        v-for="agent in agentStore.agents"
        :key="agent.id"
        :title="agent.name"
        hoverable
        class="agent-card"
      >
        <template #header-extra>
          <n-switch
            :value="agent.isEnabled"
            @update:value="(val: boolean) => handleToggle(agent.id, val)"
            :loading="togglingId === agent.id"
          />
        </template>

        <n-space vertical>
          <n-text depth="3">
            ID: {{ agent.id }}
          </n-text>
          <n-text depth="3">
            二进制: {{ agent.binaryName }}
          </n-text>
          <n-text depth="3">
            Skills 目录: {{ agent.skillsDir }}
          </n-text>
          <n-space align="center">
            <n-tag :type="agent.isInstalled ? 'success' : 'default'" size="small">
              {{ agent.isInstalled ? '已安装' : '未安装' }}
            </n-tag>
            <n-tag v-if="agent.isCustom" type="warning" size="small">
              自定义
            </n-tag>
            <n-tag type="info" size="small">
              {{ getSkillCount(agent.id) }} Skills
            </n-tag>
          </n-space>
        </n-space>

        <template #action>
          <n-space justify="end">
            <n-popconfirm
              v-if="agent.isCustom"
              @positive-click="handleRemove(agent.id)"
            >
              <template #trigger>
                <n-button type="error" size="small" :loading="removingId === agent.id">
                  <template #icon><n-icon><TrashOutline /></n-icon></template>
                  删除
                </n-button>
              </template>
              确定要删除此 Agent 吗？
            </n-popconfirm>
          </n-space>
        </template>
      </n-card>
    </div>

    <!-- 添加自定义 Agent 弹窗 -->
    <n-modal
      v-model:show="showAddModal"
      preset="dialog"
      title="添加自定义 Agent"
      positive-text="添加"
      negative-text="取消"
      :positive-button-props="{ disabled: !isFormValid }"
      :loading="adding"
      @positive-click="handleAdd"
      @close="resetForm"
      @negative-click="resetForm"
    >
      <n-form ref="formRef" :model="formData" label-placement="left" label-width="120">
        <n-form-item label="Agent ID" required>
          <n-input
            v-model:value="formData.id"
            placeholder="例如: my-agent"
          />
        </n-form-item>
        <n-form-item label="名称" required>
          <n-input
            v-model:value="formData.name"
            placeholder="例如: My Agent"
          />
        </n-form-item>
        <n-form-item label="Skills 目录" required>
          <n-input
            v-model:value="formData.skillsDir"
            placeholder="例如: /path/to/.my-agent/skills"
          />
        </n-form-item>
        <n-form-item label="二进制文件名" required>
          <n-input
            v-model:value="formData.binaryName"
            placeholder="例如: my-agent"
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
  NText,
  NTag,
  NSwitch,
  NPopconfirm,
  useMessage
} from 'naive-ui'
import { RefreshOutline, AddOutline, TrashOutline, SearchOutline } from '@vicons/ionicons5'
import { useAgentStore } from '@/stores/agentStore'
import { useSkillStore } from '@/stores/skillStore'
import type { Agent } from '@/types'

const message = useMessage()

const agentStore = useAgentStore()
const skillStore = useSkillStore()

// 状态
const showAddModal = ref(false)
const detecting = ref(false)
const adding = ref(false)
const togglingId = ref<string | null>(null)
const removingId = ref<string | null>(null)

// 表单数据
const formRef = ref()
const formData = ref({
  id: '',
  name: '',
  skillsDir: '',
  binaryName: ''
})

// 表单验证
const isFormValid = computed(() => {
  return formData.value.id.trim() &&
    formData.value.name.trim() &&
    formData.value.skillsDir.trim() &&
    formData.value.binaryName.trim()
})

// 获取 Agent 关联的 Skills 数量
function getSkillCount(agentId: string): number {
  return skillStore.skills.filter(skill => skill.agents.includes(agentId)).length
}

// 刷新列表
async function handleRefresh() {
  await Promise.all([
    agentStore.loadAgents(),
    skillStore.loadSkills()
  ])
}

// 检测已安装的 Agents
async function handleDetectInstalled() {
  detecting.value = true
  try {
    await agentStore.detectInstalled()
    message.success('检测完成')
  } catch (error) {
    message.error('检测失败: ' + (error as Error).message)
  } finally {
    detecting.value = false
  }
}

// 切换 Agent 启用状态
async function handleToggle(id: string, enabled: boolean) {
  togglingId.value = id
  try {
    await agentStore.toggleAgent(id, enabled)
    message.success(enabled ? '已启用' : '已禁用')
  } catch (error) {
    message.error('操作失败: ' + (error as Error).message)
  } finally {
    togglingId.value = null
  }
}

// 删除 Agent
async function handleRemove(id: string) {
  removingId.value = id
  try {
    await agentStore.removeAgent(id)
    message.success('删除成功')
  } catch (error) {
    message.error('删除失败: ' + (error as Error).message)
  } finally {
    removingId.value = null
  }
}

// 添加自定义 Agent
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
    message.success('添加成功')
    resetForm()
    return true
  } catch (error) {
    message.error('添加失败: ' + (error as Error).message)
    return false
  } finally {
    adding.value = false
  }
}

// 重置表单
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

.loading-container,
.empty-container {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.agents-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
  overflow-y: auto;
  flex: 1;
}

.agent-card {
  transition: transform 0.2s ease;
}

.agent-card:hover {
  transform: translateY(-2px);
}
</style>
