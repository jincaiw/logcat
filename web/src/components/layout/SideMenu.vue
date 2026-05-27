<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NMenu, NIcon } from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import {
  GridOutline, SettingsOutline, PeopleOutline, HardwareChipOutline,
  DocumentTextOutline, NotificationsOutline, ReaderOutline,
  WarningOutline, ShieldCheckmarkOutline, StatsChartOutline,
  ClipboardOutline,
} from '@vicons/ionicons5'
import { usePermissionStore, type MenuItem } from '@/stores/permission'
import { useAppStore } from '@/stores/app'

const router = useRouter()
const route = useRoute()
const permissionStore = usePermissionStore()
const appStore = useAppStore()

const iconMap: Record<string, any> = {
  'grid-outline': GridOutline,
  'settings-outline': SettingsOutline,
  'people-outline': PeopleOutline,
  'hardware-chip-outline': HardwareChipOutline,
  'document-text-outline': DocumentTextOutline,
  'notifications-outline': NotificationsOutline,
  'reader-outline': ReaderOutline,
  'warning-outline': WarningOutline,
  'shield-checkmark-outline': ShieldCheckmarkOutline,
  'stats-chart-outline': StatsChartOutline,
  'clipboard-outline': ClipboardOutline,
}

function renderIcon(iconName: string) {
  const IconComp = iconMap[iconName]
  if (IconComp) {
    return () => h(NIcon, null, { default: () => h(IconComp) })
  }
  return undefined
}

function convertMenu(menus: MenuItem[]): MenuOption[] {
  return menus.map((item) => ({
    label: item.label,
    key: item.key,
    icon: item.icon ? renderIcon(item.icon) : undefined,
    children: item.children ? convertMenu(item.children) : undefined,
  }))
}

const menuOptions = computed(() => convertMenu(permissionStore.visibleMenus))

const selectedKey = ref<string | null>(null)

function getSelectedKey(): string | null {
  const path = route.path
  // Match most specific path first
  const allKeys = getAllKeys(permissionStore.visibleMenus)
  let bestMatch: string | null = null
  for (const key of allKeys) {
    if (path.startsWith(key) && (!bestMatch || key.length > bestMatch.length)) {
      bestMatch = key
    }
  }
  return bestMatch || '/'
}

function getAllKeys(menus: MenuItem[]): string[] {
  const keys: string[] = []
  for (const menu of menus) {
    keys.push(menu.key)
    if (menu.children) {
      keys.push(...getAllKeys(menu.children))
    }
  }
  return keys
}

selectedKey.value = getSelectedKey()

function handleMenuClick(key: string) {
  router.push(key)
}
</script>

<template>
  <n-menu
    :value="selectedKey"
    :collapsed="appStore.sidebarCollapsed"
    :collapsed-width="64"
    :collapsed-icon-size="22"
    :options="menuOptions"
    :indent="24"
    @update:value="handleMenuClick"
  />
</template>

<style scoped>
.n-menu {
  height: 100%;
  padding-top: 8px;
}
</style>