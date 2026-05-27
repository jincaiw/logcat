import { computed } from 'vue'
import { useAppStore } from '@/stores/app'

export function useTheme() {
  const appStore = useAppStore()

  const theme = computed(() => appStore.theme)
  const isDark = computed(() => appStore.theme === 'dark')

  function toggleTheme() {
    appStore.toggleTheme()
  }

  function setTheme(t: 'light' | 'dark') {
    appStore.setTheme(t)
  }

  return {
    theme,
    isDark,
    toggleTheme,
    setTheme,
  }
}