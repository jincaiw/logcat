<script setup lang="ts">
import { computed, watch, onMounted } from 'vue'
import { zhCN, dateZhCN } from 'naive-ui'
import { darkTheme } from 'naive-ui'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

const naiveTheme = computed(() => appStore.theme === 'dark' ? darkTheme : undefined)
const naiveLocale = computed(() => zhCN)
const naiveDateLocale = computed(() => dateZhCN)

// Sync theme to html class
function syncTheme() {
  document.documentElement.classList.toggle('dark', appStore.theme === 'dark')
  document.documentElement.classList.toggle('light', appStore.theme === 'light')
}

onMounted(syncTheme)
watch(() => appStore.theme, syncTheme)
</script>

<template>
  <n-config-provider
    :locale="naiveLocale"
    :date-locale="naiveDateLocale"
    :theme="naiveTheme"
  >
    <n-notification-provider>
      <n-message-provider>
        <n-dialog-provider>
          <router-view />
        </n-dialog-provider>
      </n-message-provider>
    </n-notification-provider>
  </n-config-provider>
</template>