<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <div class="app-frame">
            <div class="ambient ambient-a" />
            <div class="ambient ambient-b" />
            <div class="ambient ambient-c" />

            <n-layout has-sider class="app-layout">
              <n-layout-sider
                class="app-sider"
                bordered
                collapse-mode="width"
                :collapsed-width="72"
                :width="248"
                :collapsed="collapsed"
                show-trigger
                @collapse="collapsed = true"
                @expand="collapsed = false"
              >
                <div class="logo">
                  <div class="logo-mark">
                    <img class="logo-mark-image" :src="appLogo" alt="SkillManager logo">
                  </div>
                  <div v-if="!collapsed" class="logo-copy">
                    <strong>SkillManager</strong>
                    <span>{{ t('app.brandSubtitle') }}</span>
                  </div>
                </div>

                <div v-if="!collapsed" class="sider-section-label">
                  {{ t('app.workspaceLabel') }}
                </div>

                <n-menu
                  class="app-menu"
                  :collapsed="collapsed"
                  :collapsed-width="72"
                  :collapsed-icon-size="22"
                  :options="menuOptions"
                  :value="currentMenuKey"
                  @update:value="handleMenuSelect"
                />
              </n-layout-sider>

              <n-layout class="main-layout">
                <n-layout-content class="app-content">
                  <div class="page-container">
                    <router-view />
                  </div>
                </n-layout-content>
              </n-layout>
            </n-layout>
          </div>
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import appLogo from '@/assets/logo.png'
import {
  NConfigProvider,
  NLayout,
  NLayoutSider,
  NLayoutContent,
  NMenu,
  NMessageProvider,
  NDialogProvider,
  NNotificationProvider,
  NIcon,
  type MenuOption
} from 'naive-ui'
import {
  FolderOutline,
  CloudOutline,
  SettingsOutline,
  PersonOutline
} from '@vicons/ionicons5'
import { useI18n } from '@/i18n'

const router = useRouter()
const route = useRoute()
const collapsed = ref(false)
const { t } = useI18n()

const themeOverrides = {
  common: {
    fontFamily: '"SF Pro Display", "SF Pro Text", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    primaryColor: '#2a7fff',
    primaryColorHover: '#5094ff',
    primaryColorPressed: '#1868e5',
    primaryColorSuppl: '#2a7fff',
    successColor: '#18a46f',
    warningColor: '#ff9f0a',
    errorColor: '#ff5f57',
    infoColor: '#5c84ff',
    borderRadius: '24px',
    bodyColor: 'transparent',
    cardColor: 'rgba(255, 255, 255, 0.72)',
    modalColor: 'rgba(255, 255, 255, 0.88)',
    popoverColor: 'rgba(255, 255, 255, 0.9)',
    inputColor: 'rgba(255, 255, 255, 0.76)',
    tableColor: 'rgba(255, 255, 255, 0.7)',
    textColorBase: '#15304f',
    textColor1: '#15304f',
    textColor2: '#49617d',
    textColor3: '#70839f',
    borderColor: 'rgba(255, 255, 255, 0.62)'
  },
  Layout: {
    color: 'transparent',
    siderColor: 'transparent',
    headerColor: 'transparent'
  },
  Menu: {
    itemColorActive: 'rgba(255, 255, 255, 0.92)',
    itemColorActiveCollapsed: 'rgba(255, 255, 255, 0.92)',
    itemTextColorActive: '#15304f',
    itemTextColorActiveHover: '#15304f',
    itemTextColorHover: '#1b4068',
    itemIconColorActive: '#2a7fff',
    itemIconColorActiveHover: '#2a7fff',
    itemIconColorHover: '#1b4068',
    itemBorderRadius: '18px'
  },
  Button: {
    borderRadiusMedium: '18px',
    borderRadiusLarge: '22px',
    heightMedium: '42px',
    heightLarge: '48px',
    fontWeight: '600'
  },
  Input: {
    borderRadius: '18px'
  },
  Card: {
    borderRadius: '28px',
    color: 'rgba(255, 255, 255, 0.72)',
    borderColor: 'rgba(255, 255, 255, 0.64)'
  },
  Tag: {
    borderRadius: '999px'
  },
  Modal: {
    borderRadius: '30px'
  },
  Alert: {
    borderRadius: '24px'
  },
  InputNumber: {
    borderRadius: '18px'
  },
  Popover: {
    borderRadius: '22px'
  },
  Tooltip: {
    borderRadius: '16px'
  },
  DataTable: {
    borderRadius: '22px'
  },
  Select: {
    peers: {
      InternalSelection: {
        borderRadius: '18px'
      },
      InternalSelectMenu: {
        borderRadius: '22px'
      }
    }
  },
  Dropdown: {
    peers: {
      Popover: {
        borderRadius: '22px'
      }
    }
  },
  Popconfirm: {
    peers: {
      Popover: {
        borderRadius: '22px'
      }
    }
  }
}

const currentMenuKey = computed(() => {
  const path = route.path
  if (path.startsWith('/skills')) return 'skills'
  if (path.startsWith('/registry')) return 'registry'
  if (path.startsWith('/agents')) return 'agents'
  if (path.startsWith('/settings')) return 'settings'
  return 'skills'
})

function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions = computed<MenuOption[]>(() => [
  { label: t('nav.skills'), key: 'skills', icon: renderIcon(FolderOutline) },
  { label: t('nav.registry'), key: 'registry', icon: renderIcon(CloudOutline) },
  { label: t('nav.agents'), key: 'agents', icon: renderIcon(PersonOutline) },
  { label: t('nav.settings'), key: 'settings', icon: renderIcon(SettingsOutline) }
])

function handleMenuSelect(key: string) {
  const routes: Record<string, string> = {
    skills: '/skills',
    registry: '/registry',
    agents: '/agents',
    settings: '/settings'
  }
  router.push(routes[key])
}
</script>

<style>
:root {
  --app-bg: linear-gradient(145deg, #eef6ff 0%, #f7fbff 45%, #eef7f4 100%);
  --glass-strong: rgba(255, 255, 255, 0.82);
  --glass-medium: rgba(255, 255, 255, 0.68);
  --glass-soft: rgba(255, 255, 255, 0.56);
  --glass-border: rgba(255, 255, 255, 0.7);
  --glass-shadow: 0 24px 80px rgba(90, 118, 154, 0.16);
  --glass-shadow-soft: 0 16px 44px rgba(70, 102, 136, 0.12);
  --text-primary: #14314f;
  --text-secondary: #55708d;
  --text-tertiary: #7890a8;
  --accent-blue: #2a7fff;
  --accent-cyan: #72d6ff;
  --accent-mint: #87e3b4;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body,
#app {
  height: 100%;
  overflow: hidden;
  font-family: "SF Pro Display", "SF Pro Text", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  color: var(--text-primary);
  background: var(--app-bg);
}

body {
  background-attachment: fixed;
}

button,
input,
textarea,
select {
  font: inherit;
}

#app > .n-config-provider {
  height: 100%;
}

.page-shell {
  width: 100%;
  height: 100%;
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 18px;
  overflow: hidden;
}

.page-hero,
.glass-panel,
.state-surface {
  position: relative;
  overflow: hidden;
  border: 1px solid var(--glass-border);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.82), rgba(255, 255, 255, 0.58));
  box-shadow: var(--glass-shadow-soft);
  backdrop-filter: blur(24px) saturate(160%);
  -webkit-backdrop-filter: blur(24px) saturate(160%);
}

.page-hero,
.glass-panel,
.state-surface {
  width: 100%;
  border-radius: 32px;
}

.apple-card {
  position: relative;
  overflow: hidden;
  border-radius: 30px;
  border: 1px solid rgba(255, 255, 255, 0.78);
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 18px 44px rgba(90, 118, 154, 0.14);
  isolation: isolate;
}

.page-hero::before,
.glass-panel::before,
.state-surface::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.74), rgba(255, 255, 255, 0) 44%);
}

.page-hero > *,
.glass-panel > *,
.state-surface > *,
.apple-card > * {
  position: relative;
  z-index: 1;
}

.page-hero {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
  padding: 26px 28px;
}

.hero-copy {
  position: relative;
  z-index: 1;
  max-width: 720px;
}

.hero-kicker {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--accent-blue);
}

.hero-title {
  font-size: clamp(28px, 4vw, 40px);
  line-height: 1.05;
  font-weight: 700;
  letter-spacing: -0.04em;
  color: var(--text-primary);
}

.hero-subtitle {
  margin-top: 10px;
  max-width: 680px;
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-secondary);
}

.hero-stats {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  min-width: min(440px, 100%);
}

.hero-stat {
  padding: 16px 18px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.74);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.78), rgba(255, 255, 255, 0.52));
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.65);
}

.hero-stat-value {
  display: block;
  font-size: 26px;
  line-height: 1;
  font-weight: 700;
  letter-spacing: -0.04em;
  color: var(--text-primary);
}

.hero-stat-label {
  display: block;
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-tertiary);
}

.section-bar {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.section-title {
  font-size: 18px;
  line-height: 1.2;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.section-subtitle {
  margin-top: 4px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-secondary);
}

.section-meta {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.72);
  color: var(--text-secondary);
  font-size: 12px;
  font-weight: 600;
}

.glass-panel {
  padding: 20px 22px;
}

.toolbar-panel {
  padding-bottom: 18px;
}

.content-surface {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 18px;
  overflow-y: auto !important;
}

.state-surface {
  flex: 1;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 28px;
}

.filter-toolbar {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.filter-toolbar-left,
.filter-toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.badge-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.74);
  color: var(--text-secondary);
  font-size: 12px;
  font-weight: 600;
}

.badge-chip strong {
  color: var(--text-primary);
}

.apple-card {
  transition: transform 180ms ease, box-shadow 180ms ease, border-color 180ms ease;
}

.apple-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--glass-shadow);
  border-color: rgba(255, 255, 255, 0.84);
}

.n-dialog,
.n-modal .n-card,
.n-modal-body,
.n-base-select-menu,
.n-base-selection-menu {
  backdrop-filter: blur(26px) saturate(180%);
  -webkit-backdrop-filter: blur(26px) saturate(180%);
}

.n-modal .n-card,
.n-dialog {
  border-radius: 30px;
  border: 1px solid rgba(255, 255, 255, 0.7);
  box-shadow: var(--glass-shadow);
}

/* Naive UI global border-radius fixes */
.n-alert {
  border-radius: 24px !important;
  overflow: hidden;
}

.n-descriptions {
  border-radius: 20px;
  overflow: hidden;
}

.n-descriptions .n-descriptions-table-wrapper {
  border-radius: 20px;
  overflow: hidden;
}

.n-descriptions table {
  border-radius: 20px;
  overflow: hidden;
}

.n-descriptions table tr:first-child td:first-child,
.n-descriptions table tr:first-child th:first-child {
  border-top-left-radius: 20px;
}

.n-descriptions table tr:first-child td:last-child,
.n-descriptions table tr:first-child th:last-child {
  border-top-right-radius: 20px;
}

.n-descriptions table tr:last-child td:first-child,
.n-descriptions table tr:last-child th:first-child {
  border-bottom-left-radius: 20px;
}

.n-descriptions table tr:last-child td:last-child,
.n-descriptions table tr:last-child th:last-child {
  border-bottom-right-radius: 20px;
}

.n-base-select-menu,
.n-base-selection-menu {
  border-radius: 22px !important;
  overflow: hidden;
}

.n-popover {
  border-radius: 22px !important;
}

.n-popconfirm .n-popover {
  border-radius: 22px !important;
}

.n-dropdown-menu {
  border-radius: 22px;
  overflow: hidden;
}

.n-input-number {
  border-radius: 18px;
}

.n-input-number .n-input {
  border-radius: 18px;
}

.n-input-number .n-input .n-input__border,
.n-input-number .n-input .n-input__state-border {
  border-radius: 18px;
}

.n-data-table {
  border-radius: 22px;
  overflow: hidden;
}

.n-divider {
  border-radius: 1px;
}

::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  border: 2px solid transparent;
  border-radius: 999px;
  background: rgba(79, 118, 163, 0.26);
  background-clip: content-box;
}

@media (max-width: 1100px) {
  .page-hero {
    flex-direction: column;
    align-items: stretch;
  }

  .hero-stats {
    min-width: 0;
  }
}

@media (max-width: 768px) {
  .page-shell {
    gap: 14px;
  }

  .page-hero,
  .glass-panel,
  .state-surface {
    border-radius: 26px;
  }

  .page-hero {
    padding: 22px 20px;
  }

  .glass-panel {
    padding: 18px;
  }

  .hero-title {
    font-size: 28px;
  }

  .hero-stats {
    grid-template-columns: 1fr;
  }
}
</style>

<style scoped>
.app-frame {
  position: relative;
  height: 100%;
  padding: 18px;
  overflow: hidden;
}

.ambient {
  position: absolute;
  border-radius: 999px;
  filter: blur(12px);
  pointer-events: none;
}

.ambient-a {
  top: -120px;
  right: -40px;
  width: 360px;
  height: 360px;
  background: radial-gradient(circle, rgba(110, 186, 255, 0.36), rgba(110, 186, 255, 0));
}

.ambient-b {
  bottom: -140px;
  left: -60px;
  width: 420px;
  height: 420px;
  background: radial-gradient(circle, rgba(135, 227, 180, 0.28), rgba(135, 227, 180, 0));
}

.ambient-c {
  top: 34%;
  left: 30%;
  width: 240px;
  height: 240px;
  background: radial-gradient(circle, rgba(255, 216, 168, 0.24), rgba(255, 216, 168, 0));
}

.app-layout {
  position: relative;
  z-index: 1;
  height: 100%;
  min-height: 0;
  background: transparent;
}

.app-sider,
.main-layout {
  height: 100%;
  min-height: 0;
}

.app-layout :deep(> .n-layout-scroll-container) {
  height: 100%;
  display: flex;
  gap: 18px;
  background: transparent;
}

.app-sider {
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: 34px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.78), rgba(255, 255, 255, 0.58));
  box-shadow: var(--glass-shadow);
  backdrop-filter: blur(26px) saturate(170%);
  -webkit-backdrop-filter: blur(26px) saturate(170%);
}

.app-sider :deep(> .n-layout-sider-scroll-container) {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 18px 16px 20px;
}

.main-layout {
  flex: 1 1 auto;
  width: 0;
  min-width: 0;
  background: transparent;
}

.main-layout :deep(> .n-layout-scroll-container) {
  height: 100%;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
  background: transparent;
}

.logo {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 8px 8px 18px;
}

.logo-mark {
  width: 44px;
  height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.logo-mark-image {
  display: block;
  width: 100%;
  height: 100%;
}

.logo-copy {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.logo-copy strong {
  font-size: 17px;
  line-height: 1.2;
  letter-spacing: -0.02em;
  color: var(--text-primary);
}

.logo-copy span,
.sider-section-label {
  color: var(--text-tertiary);
}

.logo-copy span {
  margin-top: 4px;
  font-size: 12px;
}

.sider-section-label {
  margin: 6px 10px 12px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.app-menu {
  flex: 1 1 auto;
  min-height: 0;
}

.app-sider :deep(.n-menu) {
  background: transparent;
}

.app-sider :deep(.n-menu-item-content) {
  margin-bottom: 8px;
  border-radius: 18px;
  height: 50px;
  font-weight: 600;
}

.app-sider :deep(.n-menu-item-content::before) {
  border-radius: 18px;
  left: 0 !important;
  right: 0 !important;
}

.app-sider :deep(.n-menu-item-content--selected) {
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.72), 0 10px 26px rgba(62, 106, 164, 0.12);
}

.app-content {
  flex: 1 1 auto;
  width: 100%;
  min-width: 0;
  min-height: 0;
  display: flex;
  overflow: hidden;
  background: transparent;
}

.app-content :deep(> .n-layout-scroll-container) {
  width: 100%;
  min-width: 0;
  height: 100%;
  display: flex;
}

.page-container {
  flex: 1 1 auto;
  width: 100%;
  min-width: 0;
  min-height: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.page-container > * {
  flex: 1 1 auto;
  width: 100%;
  min-width: 0;
}

@media (max-width: 900px) {
  .app-frame {
    padding: 14px;
  }

  .app-layout :deep(> .n-layout-scroll-container) {
    gap: 14px;
  }
}

@media (max-width: 768px) {
  .app-frame {
    padding: 10px;
  }

  .app-layout :deep(> .n-layout-scroll-container) {
    gap: 10px;
  }

  .app-sider {
    border-radius: 26px;
  }
}
</style>
