<template>
  <div class="page-shell registry-view">
    <section class="page-hero registry-hero">
      <div class="hero-copy">
        <p class="hero-kicker">{{ t('registry.heroKicker') }}</p>
        <h2 class="hero-title">{{ t('registry.title') }}</h2>
        <p class="hero-subtitle">{{ t('registry.heroSubtitle') }}</p>
      </div>

      <div class="hero-stats">
        <div class="hero-stat">
          <span class="hero-stat-value">{{ registryStore.registries.length }}</span>
          <span class="hero-stat-label">{{ t('registry.statsSources') }}</span>
        </div>
        <div class="hero-stat">
          <span class="hero-stat-value">{{ displaySkills.length }}</span>
          <span class="hero-stat-label">{{ t('registry.statsPackages') }}</span>
        </div>
        <div class="hero-stat">
          <span class="hero-stat-value">{{ totalStars }}</span>
          <span class="hero-stat-label">{{ t('registry.statsPopularity') }}</span>
        </div>
      </div>
    </section>

    <section class="glass-panel toolbar-panel registry-toolbar-panel">
      <div class="filter-toolbar">
        <div class="filter-toolbar-left registry-toolbar-left">
          <n-select
            v-model:value="selectedRegistryId"
            :options="registryOptions"
            :placeholder="t('registry.selectPlaceholder')"
            class="registry-select"
            :loading="registryStore.loading"
            @update:value="handleRegistryChange"
          />
          <n-input
            v-model:value="searchQuery"
            :placeholder="t('registry.searchPlaceholder')"
            clearable
            class="search-input"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <n-icon><SearchOutline /></n-icon>
            </template>
          </n-input>
          <n-button @click="handleSearch" :loading="searching">
            {{ t('common.search') }}
          </n-button>
        </div>

        <div class="filter-toolbar-right">
          <div class="badge-chip">
            <strong>{{ currentRegistryName }}</strong>
            <span>{{ t('registry.currentSource') }}</span>
          </div>
          <n-button @click="handleRefresh" :loading="browsing">
            <template #icon><n-icon><RefreshOutline /></n-icon></template>
            {{ t('common.refresh') }}
          </n-button>
        </div>
      </div>
    </section>

    <n-alert
      v-if="registryStore.error"
      type="error"
      :title="t('common.loadFailed', { error: registryStore.error })"
      class="error-alert"
    />

    <div v-if="registryStore.loading || browsing" class="state-surface">
      <n-spin size="large" :description="t('common.loading')" />
    </div>

    <div v-else-if="!registryStore.hasRegistries" class="state-surface">
      <n-empty :description="t('registry.emptyRegistries')">
        <template #extra>
          <n-button size="small" @click="$router.push('/settings')">
            {{ t('registry.goToSettings') }}
          </n-button>
        </template>
      </n-empty>
    </div>

    <div v-else-if="displaySkills.length === 0 && !isSearchApplied" class="state-surface">
      <n-empty :description="t('registry.emptyCurrent')" />
    </div>

    <div v-else-if="displaySkills.length === 0 && isSearchApplied" class="state-surface">
      <n-empty :description="t('registry.emptySearch')">
        <template #extra>
          <n-button size="small" @click="clearSearch">
            {{ t('registry.clearSearch') }}
          </n-button>
        </template>
      </n-empty>
    </div>

    <section v-else class="glass-panel content-surface">
      <div class="section-bar">
        <div>
          <h3 class="section-title">{{ t('registry.sectionTitle') }}</h3>
          <p class="section-subtitle">{{ t('registry.sectionSubtitle') }}</p>
        </div>
        <div class="section-meta">{{ displaySkills.length }} {{ t('registry.resultsSuffix') }}</div>
      </div>

      <div class="skills-grid">
        <article
          v-for="skill in displaySkills"
          :key="skill.id"
          class="apple-card registry-skill-card"
          @click="handleSkillClick(skill)"
        >
          <div class="card-header">
            <div class="card-copy">
              <span class="skill-overline">{{ skill.author }}</span>
              <span class="skill-name">{{ skill.name }}</span>
            </div>
            <div class="card-badges">
              <n-tag v-if="skillStore.isInstalled(skill.id)" size="small" type="success">
                {{ t('common.installed') }}
              </n-tag>
              <n-tag v-if="skill.category" size="small" type="info">
                {{ skill.category }}
              </n-tag>
            </div>
          </div>

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

          <div class="registry-card-footer">
            <div class="skill-meta">
              <span class="meta-pill">
                <n-icon><PersonOutline /></n-icon>
                {{ skill.author }}
              </span>
              <span v-if="skill.stars" class="meta-pill">
                <n-icon><StarOutline /></n-icon>
                {{ skill.stars }}
              </span>
            </div>

            <n-button
              type="primary"
              size="small"
              @click.stop="handleInstallClick(skill)"
            >
              <template #icon><n-icon><DownloadOutline /></n-icon></template>
              {{ t('common.install') }}
            </n-button>
          </div>
        </article>
      </div>
    </section>

    <n-modal
      v-model:show="showInstallModal"
      preset="dialog"
      :title="t('registry.installSkill')"
      :positive-text="t('common.install')"
      :negative-text="t('common.cancel')"
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
        <n-form-item :label="t('registry.installToAgents')">
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

    <n-modal
      v-model:show="showDetailModal"
      preset="card"
      :title="selectedSkill?.name || t('registry.skillDetail')"
      style="width: 680px; max-width: 92vw"
    >
      <div v-if="selectedSkill" class="skill-detail">
        <n-descriptions label-placement="left" :column="1" bordered>
          <n-descriptions-item :label="t('common.name')">
            {{ selectedSkill.name }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('common.author')">
            {{ selectedSkill.author }}
          </n-descriptions-item>
          <n-descriptions-item v-if="selectedSkill.category" :label="t('common.category')">
            {{ selectedSkill.category }}
          </n-descriptions-item>
          <n-descriptions-item v-if="selectedSkill.stars" :label="t('registry.popularity')">
            {{ selectedSkill.stars }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('common.description')">
            {{ selectedSkill.description }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('common.tags')">
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
          <n-descriptions-item :label="t('registry.installUrl')">
            <n-ellipsis style="max-width: 420px">
              {{ selectedSkill.installUrl }}
            </n-ellipsis>
          </n-descriptions-item>
        </n-descriptions>

        <div class="detail-actions">
          <n-button type="primary" @click="handleInstallClick(selectedSkill)">
            <template #icon><n-icon><DownloadOutline /></n-icon></template>
            {{ t('registry.installThisSkill') }}
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
  NAlert,
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
import { useSkillStore } from '@/stores/skillStore'
import type { RegistrySkill } from '@/types'
import { useI18n } from '@/i18n'

const message = useMessage()
const { t } = useI18n()

const registryStore = useRegistryStore()
const agentStore = useAgentStore()
const skillStore = useSkillStore()

const selectedRegistryId = ref<string | null>(null)
const searchQuery = ref('')
const browsing = ref(false)
const searching = ref(false)
const installing = ref(false)
const showInstallModal = ref(false)
const showDetailModal = ref(false)
const selectedSkill = ref<RegistrySkill | null>(null)
const selectedAgents = ref<string[]>([])

const registryOptions = computed(() =>
  registryStore.registries.map(r => ({
    label: r.name + (r.isDefault ? t('registry.defaultSuffix') : ''),
    value: r.id
  }))
)

const isSearchApplied = computed(() => Boolean(registryStore.searchQuery))

const displaySkills = computed(() => {
  if (registryStore.searchQuery && registryStore.searchResults.length > 0) {
    return registryStore.searchResults
  }
  if (registryStore.searchQuery && registryStore.searchResults.length === 0) {
    return []
  }
  return registryStore.browseResults
})

const currentRegistryName = computed(() => {
  return registryStore.currentRegistry?.name || t('registry.emptyRegistries')
})

const totalStars = computed(() => {
  return displaySkills.value.reduce((sum, skill) => sum + (skill.stars || 0), 0)
})

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
    await skillStore.loadSkills()
    message.success(t('registry.installSuccess', { name: selectedSkill.value.name }))
    closeInstallModal()
  } catch (error) {
    message.error(t('common.installFailed', { error: (error as Error).message }))
  } finally {
    installing.value = false
  }
}

onMounted(async () => {
  await Promise.all([
    registryStore.loadRegistries(),
    agentStore.loadAgents(),
    skillStore.loadSkills()
  ])

  if (registryStore.currentRegistry) {
    selectedRegistryId.value = registryStore.currentRegistry.id
    await loadBrowseResults()
  } else if (registryStore.registries.length > 0) {
    selectedRegistryId.value = registryStore.registries[0].id
    await loadBrowseResults()
  }
})

watch(() => registryStore.currentRegistry, (newRegistry) => {
  if (newRegistry) {
    selectedRegistryId.value = newRegistry.id
  }
})

watch(searchQuery, (value) => {
  if (!value.trim() && registryStore.searchQuery) {
    registryStore.clearSearch()
  }
})
</script>

<style scoped>
.registry-view {
  padding: 20px;
}

.registry-hero {
  align-items: center;
  padding: 18px 22px;
}

.registry-hero .hero-copy {
  max-width: 620px;
}

.registry-hero .hero-kicker {
  margin-bottom: 8px;
}

.registry-hero .hero-subtitle {
  margin-top: 6px;
}

.registry-hero .hero-stats {
  min-width: min(430px, 100%);
  gap: 10px;
}

.registry-hero .hero-stat {
  padding: 14px 16px;
}

.registry-hero .hero-stat-value {
  font-size: 24px;
}

.registry-toolbar-panel {
  padding: 16px 18px;
}

.registry-toolbar-left {
  flex: 1 1 520px;
  min-width: 0;
}

.registry-select {
  width: min(220px, 100%);
}

.search-input {
  width: min(280px, 100%);
}

.error-alert {
  border-radius: 24px;
  overflow: hidden;
}

.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 14px;
  overflow: auto;
  flex: 1;
  min-height: 0;
  align-content: start;
  padding-right: 6px;
}

.registry-skill-card {
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 196px;
  padding: 16px 18px 14px;
}

.card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.card-copy {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.skill-overline {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--text-tertiary);
}

.skill-name {
  font-size: 16px;
  line-height: 1.2;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.card-badges {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 8px;
}

.skill-description {
  color: var(--text-secondary);
  font-size: 13px;
  line-height: 1.6;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.skill-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 6px;
}

.skill-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.registry-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: auto;
}

.meta-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.7);
  color: var(--text-secondary);
  font-size: 11px;
  font-weight: 600;
}

.install-modal-content {
  padding: 8px 0;
}

.skill-preview h4 {
  margin-bottom: 8px;
  font-size: 18px;
  line-height: 1.2;
}

.skill-preview p {
  color: var(--text-secondary);
  line-height: 1.6;
}

.skill-detail {
  padding: 8px 0;
  max-height: min(72vh, 680px);
  overflow: auto;
}

.detail-actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .registry-view {
    padding: 16px;
  }

  .registry-select,
  .search-input,
  .registry-toolbar-left {
    width: 100%;
  }

  .card-header {
    flex-direction: column;
  }

  .card-badges {
    justify-content: flex-start;
  }
}
</style>
