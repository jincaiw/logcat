<script setup lang="ts">
import { computed, ref, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from '@/i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useThemeStore } from '@/stores/theme'
import { useRouter, RouterLink } from 'vue-router'
import { APP_VERSION } from '@/version'
import {
  LayoutDashboard,
  ScrollText,
  Bot,
  Router as RouterIcon,
  Filter,
  FileCode2,
  Settings,
  User,
  LogOut,
  BarChart3,
  ChevronDown,
  Play,
  Square,
  Menu,
  Sun,
  Moon,
} from '@lucide/vue'

const { t, locale, setLocale } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const themeStore = useThemeStore()
const router = useRouter()

const showUserMenu = ref(false)
const mobileMenuOpen = ref(false)

onMounted(() => {
  appStore.initApp()
  if (authStore.token && !authStore.user) {
    authStore.fetchProfile().catch(() => {})
  }
})

const menuItems = computed(() => [
  { key: 'dashboard', label: t('menu.dashboard'), icon: LayoutDashboard },
  { key: 'logs', label: t('menu.logs'), icon: ScrollText },
  { key: 'devices', label: t('menu.devices'), icon: RouterIcon },
  { key: 'log-parser', label: t('menu.parseTemplates'), icon: FileCode2 },
  { key: 'filter-policies', label: t('menu.filterPolicies'), icon: Filter },
  { key: 'robots', label: t('menu.robots'), icon: Bot },
  { key: 'stats', label: t('menu.stats'), icon: BarChart3 },
  { key: 'settings', label: t('menu.settings'), icon: Settings },
  { key: 'profile', label: t('menu.profile'), icon: User },
])

function toggleService() {
  if (appStore.serviceRunning) {
    appStore.stopService().catch(() => {})
  } else {
    appStore.startService(appStore.listenPort, appStore.protocol).catch(() => {})
  }
}

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}

// Close user menu on outside click
function handleClickOutside(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('.user-menu-wrap')) {
    showUserMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})
onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="app-layout">
    <!-- Sidebar -->
    <aside class="sidebar" :class="{ open: mobileMenuOpen }">
      <div class="sidebar-header">
        <div class="brand-logo">L</div>
        <span class="brand-text">logcat</span>
        <span class="brand-version">v{{ APP_VERSION }}</span>
      </div>

      <nav class="sidebar-nav">
        <RouterLink
          v-for="item in menuItems"
          :key="item.key"
          :to="`/${item.key}`"
          class="nav-link"
          active-class="active"
          @click="mobileMenuOpen = false"
        >
          <component :is="item.icon" :size="18" />
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>

      <div class="sidebar-footer">
        <div class="service-status">
          <span class="status-dot" :class="{ online: appStore.serviceRunning, offline: !appStore.serviceRunning }"></span>
          <span class="status-text">
            {{ appStore.serviceRunning ? t('common.running') : t('common.stopped') }}
            : {{ appStore.listenPort }}/{{ appStore.protocol.toUpperCase() }}
          </span>
        </div>
        <button class="btn-toggle-service"
                :class="{ danger: appStore.serviceRunning, success: !appStore.serviceRunning }"
                @click="toggleService"
                :disabled="appStore.loading">
          <component :is="appStore.serviceRunning ? Square : Play" :size="14" />
          {{ appStore.serviceRunning ? t('service.stop') : t('service.start') }}
        </button>
      </div>
    </aside>

    <!-- Mobile overlay -->
    <div v-if="mobileMenuOpen" class="mobile-overlay" @click="mobileMenuOpen = false"></div>

    <!-- Main content -->
    <div class="main-area">
      <header class="top-bar">
        <button class="btn-menu" @click="mobileMenuOpen = !mobileMenuOpen">
          <Menu :size="20" />
        </button>

        <div class="top-bar-right">
          <button class="btn-lang" @click="setLocale(locale === 'zh-CN' ? 'en-US' : 'zh-CN')">
            {{ locale === 'zh-CN' ? 'EN' : '中' }}
          </button>

          <button class="btn-theme-toggle" @click="themeStore.toggleTheme" :title="themeStore.isDark ? 'Switch to light mode' : 'Switch to dark mode'">
            <Moon v-if="!themeStore.isDark" :size="16" />
            <Sun v-else :size="16" />
          </button>

          <div class="user-menu-wrap">
            <button class="btn-user" @click="showUserMenu = !showUserMenu">
              <User :size="16" />
              <span>{{ authStore.user?.username || 'admin' }}</span>
              <ChevronDown :size="14" :class="{ rotate: showUserMenu }" />
            </button>
            <div v-if="showUserMenu" class="user-dropdown">
              <RouterLink to="/profile" class="dropdown-item" @click="showUserMenu = false">
                <User :size="16" />
                {{ t('menu.profile') }}
              </RouterLink>
              <RouterLink to="/settings" class="dropdown-item" @click="showUserMenu = false">
                <Settings :size="16" />
                {{ t('menu.settings') }}
              </RouterLink>
              <div class="dropdown-divider"></div>
              <button class="dropdown-item danger" @click="handleLogout">
                <LogOut :size="16" />
                {{ t('auth.logout') }}
              </button>
            </div>
          </div>
        </div>
      </header>

      <main class="content-area">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>
