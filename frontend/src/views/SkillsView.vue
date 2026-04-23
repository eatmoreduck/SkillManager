<template>
  <div class="page-shell skills-view">
    <section class="page-hero skills-hero">
      <div class="hero-copy">
        <p class="hero-kicker">{{ t('skills.heroKicker') }}</p>
        <h2 class="hero-title">{{ t('skills.title') }}</h2>
        <p class="hero-subtitle">{{ t('skills.heroSubtitle') }}</p>
      </div>

      <div class="hero-stats">
        <div class="hero-stat">
          <span class="hero-stat-value">{{ skillStore.skillCount }}</span>
          <span class="hero-stat-label">{{ t('skills.statsInstalled') }}</span>
        </div>
        <div class="hero-stat">
          <span class="hero-stat-value">{{ assignedAgentCount }}</span>
          <span class="hero-stat-label">{{ t('skills.statsAgents') }}</span>
        </div>
        <div class="hero-stat">
          <span class="hero-stat-value">{{ tagCount }}</span>
          <span class="hero-stat-label">{{ t('skills.statsTags') }}</span>
        </div>
      </div>
    </section>

    <section class="glass-panel toolbar-panel skills-toolbar-panel">
      <div class="filter-toolbar skills-toolbar">
        <div class="filter-toolbar-left skills-toolbar-left">
          <n-input
            v-model:value="searchQuery"
            :placeholder="t('skills.searchPlaceholder')"
            clearable
            class="search-input"
          >
            <template #prefix>
              <n-icon><SearchOutline /></n-icon>
            </template>
          </n-input>
          <div class="badge-chip">
            <strong>{{ filteredSkills.length }}</strong>
            <span>{{ t('skills.filteredHint') }}</span>
          </div>
        </div>

        <div class="filter-toolbar-right">
          <n-button type="primary" @click="showInstallModal = true">
            <template #icon><n-icon><AddOutline /></n-icon></template>
            {{ t('skills.installFromUrl') }}
          </n-button>
        </div>
      </div>
    </section>

    <n-alert
      v-if="skillStore.error"
      type="error"
      :title="t('skills.loadFailedTitle')"
      class="error-alert"
    >
      {{ skillStore.error }}
    </n-alert>

    <div v-if="skillStore.loading" class="state-surface">
      <n-spin size="large" :description="t('common.loading')" />
    </div>

    <div v-else-if="filteredSkills.length === 0" class="state-surface">
      <n-empty :description="t('skills.emptyInstalled')">
        <template #extra>
          <n-button size="small" @click="$router.push('/registry')">
            {{ t('skills.browseRegistry') }}
          </n-button>
        </template>
      </n-empty>
    </div>

    <section v-else class="glass-panel content-surface">
      <div class="section-bar">
        <div>
          <h3 class="section-title">{{ t('skills.sectionTitle') }}</h3>
          <p class="section-subtitle">{{ t('skills.sectionSubtitle') }}</p>
        </div>
        <div class="section-meta">{{ filteredSkills.length }} {{ t('skills.resultsSuffix') }}</div>
      </div>

      <div class="skills-grid">
        <SkillCard
          v-for="skill in filteredSkills"
          :key="skill.id"
          :skill="skill"
          @click="handleSkillClick(skill)"
          @uninstall="handleUninstall"
        />
      </div>
    </section>

    <n-modal
      v-model:show="showInstallModal"
      preset="dialog"
      :title="t('skills.installFromUrlTitle')"
      :positive-text="t('common.install')"
      :negative-text="t('common.cancel')"
      :positive-button-props="{ disabled: !installURL || selectedAgents.length === 0 }"
      :loading="installing"
      @positive-click="handleInstall"
    >
      <n-form label-placement="left" label-width="100">
        <n-form-item :label="t('skills.gitUrl')">
          <n-input
            v-model:value="installURL"
            placeholder="https://github.com/user/skill-name.git"
          />
        </n-form-item>
        <n-form-item :label="t('common.selectAgent')">
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
import { useI18n } from '@/i18n'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const { t } = useI18n()

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
  const all = Array.isArray(skillStore.skills) ? skillStore.skills : []
  if (!searchQuery.value) return all
  const query = searchQuery.value.toLowerCase()
  return all.filter(skill =>
    (skill.name || '').toLowerCase().includes(query) ||
    (skill.description || '').toLowerCase().includes(query) ||
    normalizedTags(skill.tags).some(tag => tag.toLowerCase().includes(query))
  )
})

const assignedAgentCount = computed(() => {
  const all = Array.isArray(skillStore.skills) ? skillStore.skills : []
  return new Set(all.flatMap(skill => skill.agents || [])).size
})

const tagCount = computed(() => {
  const all = Array.isArray(skillStore.skills) ? skillStore.skills : []
  return new Set(all.flatMap(skill => normalizedTags(skill.tags))).size
})

function handleSkillClick(skill: Skill) {
  router.push(`/skills/${encodeURIComponent(skill.id)}`)
}

async function handleInstall() {
  if (!installURL.value || selectedAgents.value.length === 0) return

  installing.value = true
  try {
    await skillStore.installSkill(installURL.value, selectedAgents.value)
    message.success(t('common.installSuccess'))
    showInstallModal.value = false
    installURL.value = ''
    selectedAgents.value = []
  } catch (error) {
    message.error(t('common.installFailed', { error: (error as Error).message }))
  } finally {
    installing.value = false
  }
}

function handleUninstall(skillId: string) {
  dialog.warning({
    title: t('skills.confirmUninstallTitle'),
    content: t('skills.confirmUninstallContent'),
    positiveText: t('common.uninstall'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await skillStore.uninstallSkill(skillId)
        message.success(t('common.uninstallSuccess'))
      } catch (error) {
        message.error(t('common.uninstallFailed', { error: (error as Error).message }))
      }
    }
  })
}

onMounted(async () => {
  await Promise.all([
    skillStore.loadSkills(),
    agentStore.loadAgents()
  ])
  selectedAgents.value = agentStore.enabledAgents.map(a => a.id)
})
</script>

<style scoped>
.skills-view {
  padding: 20px;
}

.skills-hero {
  align-items: center;
  padding: 18px 22px;
}

.skills-hero .hero-copy {
  max-width: 600px;
}

.skills-hero .hero-kicker {
  margin-bottom: 8px;
}

.skills-hero .hero-subtitle {
  margin-top: 6px;
}

.skills-hero .hero-stats {
  min-width: min(380px, 100%);
  gap: 8px;
}

.skills-hero .hero-stat {
  padding: 12px 14px;
  border-radius: 20px;
}

.skills-hero .hero-stat-value {
  font-size: 22px;
}

.skills-toolbar-panel {
  padding: 16px 18px;
}

.skills-toolbar-left {
  flex: 1 1 320px;
  min-width: 0;
}

.search-input {
  width: min(320px, 100%);
}

.error-alert {
  border-radius: 24px;
  overflow: hidden;
}

.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 14px;
  overflow: auto;
  flex: 1;
  min-height: 0;
  align-content: start;
  padding-right: 4px;
}

@media (max-width: 768px) {
  .skills-view {
    padding: 16px;
  }

  .search-input,
  .skills-toolbar-left {
    width: 100%;
  }
}
</style>
