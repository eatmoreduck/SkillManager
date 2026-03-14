<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-layout has-sider class="app-layout">
            <!-- 侧边栏 -->
            <n-layout-sider
              bordered
              collapse-mode="width"
              :collapsed-width="64"
              :width="200"
              :collapsed="collapsed"
              show-trigger
              @collapse="collapsed = true"
              @expand="collapsed = false"
            >
              <div class="logo">
                <span v-if="!collapsed">SkillManager</span>
                <span v-else>SM</span>
              </div>

              <n-menu
                :collapsed="collapsed"
                :collapsed-width="64"
                :collapsed-icon-size="22"
                :options="menuOptions"
                :value="currentMenuKey"
                @update:value="handleMenuSelect"
              />
            </n-layout-sider>

            <!-- 主内容区 -->
            <n-layout>
              <n-layout-header bordered class="app-header">
                <span class="page-title">{{ pageTitle }}</span>
              </n-layout-header>

              <n-layout-content class="app-content">
                <router-view />
              </n-layout-content>
            </n-layout>
          </n-layout>
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NConfigProvider,
  NLayout,
  NLayoutSider,
  NLayoutHeader,
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

const router = useRouter()
const route = useRoute()
const collapsed = ref(false)

const themeOverrides = {
  common: {
    primaryColor: '#18a058',
    primaryColorHover: '#36ad6a',
    primaryColorPressed: '#0c7a43'
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

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    skills: 'My Skills',
    registry: 'Registry',
    agents: 'Agents',
    settings: 'Settings'
  }
  return titles[currentMenuKey.value] || 'SkillManager'
})

function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions: MenuOption[] = [
  { label: 'My Skills', key: 'skills', icon: renderIcon(FolderOutline) },
  { label: 'Registry', key: 'registry', icon: renderIcon(CloudOutline) },
  { label: 'Agents', key: 'agents', icon: renderIcon(PersonOutline) },
  { label: 'Settings', key: 'settings', icon: renderIcon(SettingsOutline) }
]

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
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
html, body, #app {
  height: 100%;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
}
</style>

<style scoped>
.app-layout {
  height: 100%;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: var(--n-text-color);
  border-bottom: 1px solid var(--n-border-color);
}

.app-header {
  height: 48px;
  padding: 0 20px;
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 16px;
  font-weight: 600;
}

.app-content {
  height: calc(100% - 48px);
  overflow: hidden;
}
</style>
