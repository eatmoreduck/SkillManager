<template>
  <div class="page-shell settings-view">
    <section class="page-hero settings-hero">
      <div class="hero-copy">
        <p class="hero-kicker">{{ t('settings.heroKicker') }}</p>
        <h2 class="hero-title">{{ t('settings.title') }}</h2>
        <p class="hero-subtitle">{{ t('settings.heroSubtitle') }}</p>
      </div>

      <div class="hero-stats">
        <div class="hero-stat">
          <span class="hero-stat-value">{{ registryStore.registries.length }}</span>
          <span class="hero-stat-label">{{ t('settings.statsRegistries') }}</span>
        </div>
        <div class="hero-stat">
          <span class="hero-stat-value">{{ proxyForm.enabled ? t('common.enabled') : t('common.disabled') }}</span>
          <span class="hero-stat-label">{{ t('settings.statsProxy') }}</span>
        </div>
        <div class="hero-stat">
          <span class="hero-stat-value hero-stat-value-text">{{ currentLanguageLabel }}</span>
          <span class="hero-stat-label">{{ t('settings.statsLanguage') }}</span>
        </div>
      </div>
    </section>

    <div class="settings-grid">
      <section class="glass-panel settings-card">
        <div class="settings-card-head">
          <div>
            <h3 class="section-title">{{ t('settings.proxySettings') }}</h3>
            <p class="section-subtitle">{{ t('settings.proxyHint') }}</p>
          </div>
        </div>

        <n-form
          ref="proxyFormRef"
          :model="proxyForm"
          label-placement="left"
          label-width="112"
          class="settings-form"
        >
          <n-form-item :label="t('settings.enableProxy')">
            <n-switch v-model:value="proxyForm.enabled" />
          </n-form-item>

          <n-form-item :label="t('settings.proxyType')">
            <n-select
              v-model:value="proxyForm.type"
              :options="proxyTypeOptions"
              :disabled="!proxyForm.enabled"
              style="width: 220px"
            />
          </n-form-item>

          <n-form-item :label="t('common.language')">
            <n-select
              v-model:value="selectedLanguage"
              :options="languageOptions"
              style="width: 240px"
            />
          </n-form-item>

          <n-form-item :label="t('settings.host')">
            <n-input
              v-model:value="proxyForm.host"
              placeholder="127.0.0.1"
              :disabled="!proxyForm.enabled"
            />
          </n-form-item>

          <n-form-item :label="t('settings.port')">
            <n-input-number
              v-model:value="proxyForm.port"
              :min="1"
              :max="65535"
              placeholder="7890"
              :disabled="!proxyForm.enabled"
              style="width: 180px"
            />
          </n-form-item>

          <n-form-item :label="t('settings.username')">
            <n-input
              v-model:value="proxyForm.username"
              :placeholder="t('settings.optional')"
              :disabled="!proxyForm.enabled"
            />
          </n-form-item>

          <n-form-item :label="t('settings.password')">
            <n-input
              v-model:value="proxyForm.password"
              type="password"
              :placeholder="t('settings.optional')"
              :disabled="!proxyForm.enabled"
              show-password-on="click"
            />
          </n-form-item>
        </n-form>
      </section>

      <section class="glass-panel settings-card registry-card">
        <div class="settings-card-head">
          <div>
            <h3 class="section-title">{{ t('settings.registryManagement') }}</h3>
            <p class="section-subtitle">{{ t('settings.registryHint') }}</p>
          </div>
          <n-button type="primary" size="small" @click="showAddRegistryModal = true">
            <template #icon><n-icon><AddOutline /></n-icon></template>
            {{ t('settings.addRegistry') }}
          </n-button>
        </div>

        <div v-if="registryStore.loading" class="settings-loading">
          <n-spin size="medium" :description="t('common.loading')" />
        </div>

        <n-empty
          v-else-if="registryStore.registries.length === 0"
          :description="t('settings.noRegistries')"
        />

        <div v-else class="registry-stack">
          <div v-for="registry in registryStore.registries" :key="registry.id" class="registry-row">
            <div class="registry-copy">
              <div class="registry-name-row">
                <span class="registry-name">{{ registry.name }}</span>
                <n-tag v-if="registry.isDefault" type="success" size="small">
                  {{ t('common.default') }}
                </n-tag>
              </div>
              <span class="registry-url">{{ registry.url }}</span>
            </div>
            <n-button
              quaternary
              type="error"
              size="small"
              :disabled="registry.isDefault"
              @click="handleRemoveRegistry(registry.id)"
            >
              <template #icon><n-icon><TrashOutline /></n-icon></template>
            </n-button>
          </div>
        </div>
      </section>
    </div>

    <section class="glass-panel action-bar">
      <div class="action-copy">
        <h3 class="section-title">{{ t('settings.saveSettings') }}</h3>
        <p class="section-subtitle">{{ t('settings.saveHint') }}</p>
      </div>
      <n-button
        type="primary"
        size="large"
        :loading="saving"
        :disabled="!hasChanges"
        @click="handleSave"
      >
        <template #icon><n-icon><SaveOutline /></n-icon></template>
        {{ t('settings.saveSettings') }}
      </n-button>
    </section>

    <n-modal
      v-model:show="showAddRegistryModal"
      preset="dialog"
      :title="t('settings.addRegistryTitle')"
      :positive-text="t('common.add')"
      :negative-text="t('common.cancel')"
      :positive-button-props="{ disabled: !newRegistryName || !newRegistryUrl }"
      :loading="addingRegistry"
      @positive-click="handleAddRegistry"
    >
      <n-form label-placement="left" label-width="80">
        <n-form-item :label="t('common.name')">
          <n-input
            v-model:value="newRegistryName"
            :placeholder="t('settings.officialRegistryExample')"
          />
        </n-form-item>
        <n-form-item label="URL">
          <n-input
            v-model:value="newRegistryUrl"
            :placeholder="t('settings.registryUrlExample')"
          />
        </n-form-item>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
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
  NTag,
  NModal,
  useMessage,
  useDialog
} from 'naive-ui'
import { AddOutline, TrashOutline, SaveOutline } from '@vicons/ionicons5'
import { useConfigStore } from '@/stores/configStore'
import { useRegistryStore } from '@/stores/registryStore'
import type { ProxyConfig } from '@/types'
import { languageOptions, normalizeLocale, setLocale, useI18n } from '@/i18n'

const message = useMessage()
const dialog = useDialog()
const { t } = useI18n()

const configStore = useConfigStore()
const registryStore = useRegistryStore()

const proxyFormRef = ref()
const proxyForm = ref<ProxyConfig>({
  enabled: false,
  type: 'http',
  host: '',
  port: 7890,
  username: '',
  password: ''
})

const originalProxy = ref<ProxyConfig>({
  enabled: false,
  type: 'http',
  host: '',
  port: 7890,
  username: '',
  password: ''
})

const proxyTypeOptions = [
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'SOCKS5', value: 'socks5' }
]

const selectedLanguage = ref('zh-CN')
const originalLanguage = ref('zh-CN')
const showAddRegistryModal = ref(false)
const newRegistryName = ref('')
const newRegistryUrl = ref('')
const addingRegistry = ref(false)
const saving = ref(false)

const hasChanges = computed(() => {
  return JSON.stringify(proxyForm.value) !== JSON.stringify(originalProxy.value) ||
    selectedLanguage.value !== originalLanguage.value
})

const currentLanguageLabel = computed(() => {
  return languageOptions.find(option => option.value === selectedLanguage.value)?.label || selectedLanguage.value
})

async function loadSettings() {
  try {
    await configStore.loadConfig()
    if (configStore.config?.proxy) {
      proxyForm.value = { ...configStore.config.proxy }
      originalProxy.value = { ...configStore.config.proxy }
    }
    selectedLanguage.value = normalizeLocale(configStore.config?.language)
    originalLanguage.value = selectedLanguage.value
    await registryStore.loadRegistries()
  } catch (error) {
    message.error(t('settings.loadFailed', { error: (error as Error).message }))
  }
}

async function handleSave() {
  saving.value = true
  try {
    if (JSON.stringify(proxyForm.value) !== JSON.stringify(originalProxy.value)) {
      await configStore.updateProxy(proxyForm.value)
      originalProxy.value = { ...proxyForm.value }
    }
    if (selectedLanguage.value !== originalLanguage.value) {
      await configStore.updateLanguage(selectedLanguage.value)
      originalLanguage.value = selectedLanguage.value
      setLocale(selectedLanguage.value)
    }
    message.success(t('common.saveSuccess'))
  } catch (error) {
    message.error(t('common.saveFailed', { error: (error as Error).message }))
  } finally {
    saving.value = false
  }
}

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
    message.success(t('common.addSuccess'))
    showAddRegistryModal.value = false
    newRegistryName.value = ''
    newRegistryUrl.value = ''
  } catch (error) {
    message.error(t('common.addFailed', { error: (error as Error).message }))
  } finally {
    addingRegistry.value = false
  }
}

function handleRemoveRegistry(id: string) {
  const registry = registryStore.registries.find(r => r.id === id)
  if (registry?.isDefault) {
    message.warning(t('settings.cannotRemoveDefaultRegistry'))
    return
  }

  dialog.warning({
    title: t('settings.confirmDeleteRegistryTitle'),
    content: t('settings.confirmDeleteRegistryContent', { name: registry?.name || '' }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await registryStore.removeRegistry(id)
        message.success(t('common.deleteSuccess'))
      } catch (error) {
        message.error(t('common.deleteFailed', { error: (error as Error).message }))
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
  overflow: auto;
}

.settings-hero {
  align-items: center;
  padding: 18px 22px;
}

.settings-hero .hero-copy {
  max-width: 620px;
}

.settings-hero .hero-kicker {
  margin-bottom: 8px;
}

.settings-hero .hero-subtitle {
  margin-top: 6px;
  max-width: 560px;
  line-height: 1.5;
}

.settings-hero .hero-stats {
  gap: 10px;
  min-width: min(420px, 100%);
}

.settings-hero .hero-stat {
  padding: 12px 14px;
  border-radius: 20px;
}

.settings-hero .hero-stat-value {
  font-size: 22px;
}

.settings-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.04fr) minmax(300px, 0.86fr);
  gap: 14px;
  align-items: start;
}

.settings-card {
  padding: 18px 20px;
}

.settings-card-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}

.settings-form {
  margin-top: 2px;
}

.settings-form :deep(.n-form-item) {
  margin-bottom: 10px;
}

.settings-form :deep(.n-form-item-feedback-wrapper) {
  min-height: 0;
}

.settings-loading {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.registry-stack {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.registry-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 14px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.62);
  border: 1px solid rgba(255, 255, 255, 0.7);
}

.registry-copy {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.registry-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.registry-name {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

.registry-url {
  color: var(--text-secondary);
  line-height: 1.5;
  word-break: break-all;
}

.action-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  position: sticky;
  bottom: 0;
  z-index: 2;
  padding: 16px 18px;
}

.action-copy {
  min-width: 0;
}

.hero-stat-value-text {
  font-size: 18px;
  line-height: 1.2;
}

@media (max-width: 960px) {
  .settings-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .settings-view {
    padding: 16px;
  }

  .settings-card-head,
  .action-bar,
  .registry-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .settings-card {
    padding: 16px;
  }

  .action-bar :deep(.n-button) {
    width: 100%;
  }
}
</style>
