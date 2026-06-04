<script setup lang="ts">
import { h } from 'vue'
import { useRouter } from 'vue-router'
import { NDropdown, NButton, NIcon, NAvatar, NSpace, NTooltip } from 'naive-ui'
import type { DropdownOption } from 'naive-ui'
import {
  SunnyOutline, MoonOutline, LogOutOutline, KeyOutline, MenuOutline,
} from '@vicons/ionicons5'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useIsMobile } from '@/composables/useIsMobile'
import Breadcrumb from './Breadcrumb.vue'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const { isMobile } = useIsMobile()

const userOptions: DropdownOption[] = [
  {
    label: '修改密码',
    key: 'change-password',
    icon: () => h(NIcon, null, { default: () => h(KeyOutline) }),
  },
  {
    type: 'divider',
    key: 'd1',
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h(NIcon, null, { default: () => h(LogOutOutline) }),
  },
]

async function handleUserSelect(key: string) {
  if (key === 'logout') {
    await authStore.logout()
    router.push('/login')
  } else if (key === 'change-password') {
    router.push('/change-password')
  }
}
</script>

<template>
  <div class="top-header">
    <div class="top-header-left">
      <n-button text class="menu-toggle-btn" @click="appStore.toggleSidebar()">
        <n-icon :component="MenuOutline" size="20" />
      </n-button>
      <Breadcrumb />
    </div>

    <div class="top-header-right">
      <n-space align="center" :size="8">
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button text class="theme-toggle-btn" @click="appStore.toggleTheme()">
              <n-icon v-if="appStore.theme === 'light'" :component="MoonOutline" size="18" />
              <n-icon v-else :component="SunnyOutline" size="18" />
            </n-button>
          </template>
          {{ appStore.theme === 'light' ? '切换暗色模式' : '切换亮色模式' }}
        </n-tooltip>

        <n-dropdown :options="userOptions" @select="handleUserSelect" trigger="click">
          <div class="user-info">
            <n-avatar size="small" round class="user-avatar">
              {{ (authStore.user?.displayName || authStore.username || 'U').charAt(0).toUpperCase() }}
            </n-avatar>
            <span v-if="!isMobile" class="username-text">{{ authStore.user?.displayName || authStore.username }}</span>
          </div>
        </n-dropdown>
      </n-space>
    </div>
  </div>
</template>

<style scoped>
.top-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 20px;
  background: var(--bg-color-card);
  border-bottom: 1px solid var(--border-color);
  transition: background-color 0.3s ease, border-color 0.3s ease;
}

.top-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.top-header-right {
  display: flex;
  align-items: center;
}

.menu-toggle-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-color-secondary);
  transition: all 0.2s ease;
}

.menu-toggle-btn:hover {
  background: var(--bg-color-embedded);
  color: var(--text-color);
}

.theme-toggle-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-color-secondary);
  transition: all 0.2s ease;
}

.theme-toggle-btn:hover {
  background: var(--bg-color-embedded);
  color: var(--text-color);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.user-info:hover {
  background: var(--bg-color-embedded);
}

.user-avatar {
  background-color: var(--primary-color) !important;
  font-size: 13px;
}

.username-text {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  color: var(--text-color);
}
</style>
