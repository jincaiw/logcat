import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const theme = ref<'light' | 'dark'>((localStorage.getItem('theme') as 'light' | 'dark') || 'light')
  const sidebarCollapsed = ref(localStorage.getItem('sidebarCollapsed') === 'true')
  const locale = ref<'zhCN' | 'enUS'>('zhCN')

  function toggleTheme() {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem('theme', theme.value)
  }

  function setTheme(t: 'light' | 'dark') {
    theme.value = t
    localStorage.setItem('theme', t)
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
    localStorage.setItem('sidebarCollapsed', String(sidebarCollapsed.value))
  }

  function setSidebarCollapsed(collapsed: boolean) {
    sidebarCollapsed.value = collapsed
    localStorage.setItem('sidebarCollapsed', String(collapsed))
  }

  return {
    theme,
    sidebarCollapsed,
    locale,
    toggleTheme,
    setTheme,
    toggleSidebar,
    setSidebarCollapsed,
  }
})