<script setup lang="ts">
import { h } from 'vue'
import { useRouter } from 'vue-router'
import { NDropdown, NButton, NIcon, NAvatar, NSpace } from 'naive-ui'
import type { DropdownOption } from 'naive-ui'
import {
  SunnyOutline, MoonOutline, PersonCircleOutline,
  LogOutOutline, KeyOutline, MenuOutline,
} from '@vicons/ionicons5'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import Breadcrumb from './Breadcrumb.vue'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

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

function handleUserSelect(key: string) {
  if (key === 'logout') {
    authStore.logout()
    router.push('/login')
  } else if (key === 'change-password') {
    router.push('/change-password')
  }
}

function renderUserLabel() {
  return h(
    'div',
    { style: { display: 'flex', alignItems: 'center', gap: '8px' } },
    [
      h(NAvatar, { size: 'small', round: true, style: { backgroundColor: 'var(--primary-color)' } }, {
        default: () => (authStore.user?.displayName || authStore.username || 'U').charAt(0).toUpperCase(),
      }),
      h('span', null, authStore.user?.displayName || authStore.username),
    ]
  )
}
</script>

<template>
  <div class="top-header">
    <div class="top-header-left">
      <n-button text style="font-size: 20px" @click="appStore.toggleSidebar()">
        <n-icon :component="MenuOutline" />
      </n-button>
      <Breadcrumb />
    </div>

    <div class="top-header-right">
      <n-space align="center">
        <n-button text style="font-size: 20px" @click="appStore.toggleTheme()">
          <n-icon v-if="appStore.theme === 'light'" :component="MoonOutline" />
          <n-icon v-else :component="SunnyOutline" />
        </n-button>

        <n-dropdown :options="userOptions" @select="handleUserSelect" trigger="click">
          <n-button text>
            <component :is="renderUserLabel()" />
          </n-button>
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
</style>