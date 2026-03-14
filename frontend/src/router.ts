import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/skills'
  },
  {
    path: '/skills',
    name: 'Skills',
    component: () => import('@/views/SkillsView.vue')
  },
  {
    path: '/skills/:id',
    name: 'SkillDetail',
    component: () => import('@/views/SkillDetailView.vue')
  },
  {
    path: '/registry',
    name: 'Registry',
    component: () => import('@/views/RegistryView.vue')
  },
  {
    path: '/agents',
    name: 'Agents',
    component: () => import('@/views/AgentsView.vue')
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/SettingsView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
