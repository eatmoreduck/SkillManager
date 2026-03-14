<template>
  <div class="settings-view">
    <!-- 代理设置区域 -->
    <n-card title="代理设置" class="settings-card">
      <n-form
        ref="proxyFormRef"
        :model="proxyForm"
        label-placement="left"
        label-width="100"
      >
        <n-form-item label="启用代理">
          <n-switch v-model:value="proxyForm.enabled" />
        </n-form-item>

        <n-form-item label="代理类型">
          <n-select
            v-model:value="proxyForm.type"
            :options="proxyTypeOptions"
            :disabled="!proxyForm.enabled"
            style="width: 200px"
          />
        </n-form-item>

        <n-form-item label="主机">
          <n-input
            v-model:value="proxyForm.host"
            placeholder="127.0.0.1"
            :disabled="!proxyForm.enabled"
            style="width: 300px"
          />
        </n-form-item>

        <n-form-item label="端口">
          <n-input-number
            v-model:value="proxyForm.port"
            :min="1"
            :max="65535"
            placeholder="7890"
            :disabled="!proxyForm.enabled"
            style="width: 150px"
          />
        </n-form-item>

        <n-form-item label="用户名">
          <n-input
            v-model:value="proxyForm.username"
            placeholder="可选"
            :disabled="!proxyForm.enabled"
            style="width: 300px"
          />
        </n-form-item>

        <n-form-item label="密码">
          <n-input
            v-model:value="proxyForm.password"
            type="password"
            placeholder="可选"
            :disabled="!proxyForm.enabled"
            show-password-on="click"
            style="width: 300px"
          />
        </n-form-item>
      </n-form>
    </n-card>

    <!-- Registry 管理区域 -->
    <n-card title="Registry 管理" class="settings-card">
      <template #header-extra>
        <n-button type="primary" size="small" @click="showAddRegistryModal = true">
          <template #icon><n-icon><AddOutline /></n-icon></template>
          添加 Registry
        </n-button>
      </template>

      <!-- 加载状态 -->
      <div v-if="registryStore.loading" class="loading-container">
        <n-spin size="medium" description="加载中..." />
      </div>

      <!-- 空状态 -->
      <n-empty
        v-else-if="registryStore.registries.length === 0"
        description="暂无配置的 Registry"
      />

      <!-- Registry 列表 -->
      <n-list v-else bordered>
        <n-list-item v-for="registry in registryStore.registries" :key="registry.id">
          <template #prefix>
            <n-tag v-if="registry.isDefault" type="success" size="small">
              默认
            </n-tag>
          </template>
          <n-thing :title="registry.name" :description="registry.url" />
          <template #suffix>
            <n-button
              quaternary
              type="error"
              size="small"
              :disabled="registry.isDefault"
              @click="handleRemoveRegistry(registry.id)"
            >
              <template #icon><n-icon><TrashOutline /></n-icon></template>
            </n-button>
          </template>
        </n-list-item>
      </n-list>
    </n-card>

    <!-- 保存按钮 -->
    <div class="action-bar">
      <n-button
        type="primary"
        size="large"
        :loading="saving"
        :disabled="!hasChanges"
        @click="handleSave"
      >
        <template #icon><n-icon><SaveOutline /></n-icon></template>
        保存设置
      </n-button>
    </div>

    <!-- 添加 Registry 弹窗 -->
    <n-modal
      v-model:show="showAddRegistryModal"
      preset="dialog"
      title="添加 Registry"
      positive-text="添加"
      negative-text="取消"
      :positive-button-props="{ disabled: !newRegistryName || !newRegistryUrl }"
      :loading="addingRegistry"
      @positive-click="handleAddRegistry"
    >
      <n-form label-placement="left" label-width="80">
        <n-form-item label="名称">
          <n-input
            v-model:value="newRegistryName"
            placeholder="例如: Official Registry"
          />
        </n-form-item>
        <n-form-item label="URL">
          <n-input
            v-model:value="newRegistryUrl"
            placeholder="例如: https://registry.example.com"
          />
        </n-form-item>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  NCard,
  NForm,
  NFormItem,
  NSwitch,
  NSelect,
  NInput,
  NInputNumber,
  NButton,
  NIcon,
  NSpin,
  NEmpty,
  NList,
  NListItem,
  NThing,
  NTag,
  NModal,
  useMessage,
  useDialog
} from 'naive-ui'
import { AddOutline, TrashOutline, SaveOutline } from '@vicons/ionicons5'
import { useConfigStore } from '@/stores/configStore'
import { useRegistryStore } from '@/stores/registryStore'
import type { ProxyConfig } from '@/types'

const message = useMessage()
const dialog = useDialog()

const configStore = useConfigStore()
const registryStore = useRegistryStore()

// 代理表单
const proxyFormRef = ref()
const proxyForm = ref<ProxyConfig>({
  enabled: false,
  type: 'http',
  host: '',
  port: 7890,
  username: '',
  password: ''
})

// 原始代理配置（用于检测变更）
const originalProxy = ref<ProxyConfig>({
  enabled: false,
  type: 'http',
  host: '',
  port: 7890,
  username: '',
  password: ''
})

// 代理类型选项
const proxyTypeOptions = [
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'SOCKS5', value: 'socks5' }
]

// 添加 Registry 相关
const showAddRegistryModal = ref(false)
const newRegistryName = ref('')
const newRegistryUrl = ref('')
const addingRegistry = ref(false)

// 保存状态
const saving = ref(false)

// 检测是否有变更
const hasChanges = computed(() => {
  return JSON.stringify(proxyForm.value) !== JSON.stringify(originalProxy.value)
})

// 加载配置
async function loadSettings() {
  try {
    await configStore.loadConfig()
    if (configStore.config?.proxy) {
      proxyForm.value = { ...configStore.config.proxy }
      originalProxy.value = { ...configStore.config.proxy }
    }
    await registryStore.loadRegistries()
  } catch (error) {
    message.error('加载配置失败: ' + (error as Error).message)
  }
}

// 保存设置
async function handleSave() {
  saving.value = true
  try {
    await configStore.updateProxy(proxyForm.value)
    originalProxy.value = { ...proxyForm.value }
    message.success('保存成功')
  } catch (error) {
    message.error('保存失败: ' + (error as Error).message)
  } finally {
    saving.value = false
  }
}

// 添加 Registry
async function handleAddRegistry() {
  if (!newRegistryName.value || !newRegistryUrl.value) return

  addingRegistry.value = true
  try {
    await registryStore.addRegistry({
      id: '',
      name: newRegistryName.value,
      url: newRegistryUrl.value,
      isDefault: false
    })
    message.success('添加成功')
    showAddRegistryModal.value = false
    newRegistryName.value = ''
    newRegistryUrl.value = ''
  } catch (error) {
    message.error('添加失败: ' + (error as Error).message)
  } finally {
    addingRegistry.value = false
  }
}

// 删除 Registry
function handleRemoveRegistry(id: string) {
  const registry = registryStore.registries.find(r => r.id === id)
  if (registry?.isDefault) {
    message.warning('无法删除默认 Registry')
    return
  }

  dialog.warning({
    title: '确认删除',
    content: `确定要删除 Registry "${registry?.name}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await registryStore.removeRegistry(id)
        message.success('删除成功')
      } catch (error) {
        message.error('删除失败: ' + (error as Error).message)
      }
    }
  })
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings-view {
  padding: 20px;
  height: 100%;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-card {
  flex-shrink: 0;
}

.loading-container {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.action-bar {
  display: flex;
  justify-content: flex-end;
  padding: 10px 0;
  flex-shrink: 0;
}
</style>
