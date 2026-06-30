import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

const STORAGE_KEY = 'logcat-theme'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(loadSavedTheme())

  function loadSavedTheme(): boolean {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved === 'dark' || saved === 'light') return saved === 'dark'
    return false // default light
  }

  function toggleTheme() {
    isDark.value = !isDark.value
  }

  watch(isDark, (val) => {
    localStorage.setItem(STORAGE_KEY, val ? 'dark' : 'light')
    document.documentElement.classList.toggle('dark', val)
    document.documentElement.classList.toggle('light', !val)
  }, { immediate: true })

  return { isDark, toggleTheme }
})
