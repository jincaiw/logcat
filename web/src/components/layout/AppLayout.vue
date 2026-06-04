<script setup lang="ts">
import { watch, onMounted, computed } from 'vue'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent } from 'naive-ui'
import { useAppStore } from '@/stores/app'
import { useIsMobile } from '@/composables/useIsMobile'
import { useTimeFormat } from '@/composables/useTimeFormat'
import SideMenu from './SideMenu.vue'
import TopHeader from './TopHeader.vue'

const appStore = useAppStore()
const { isMobile } = useIsMobile()
const { loadTimezoneConfig } = useTimeFormat()

// True when the sidebar is fully hidden and should not capture any clicks.
const siderInactive = computed(() => isMobile.value && appStore.sidebarCollapsed)

watch(isMobile, (val) => {
  if (val) {
    appStore.sidebarCollapsed = true
  }
}, { immediate: true })

onMounted(() => {
  if (isMobile.value) {
    appStore.sidebarCollapsed = true
  }
  loadTimezoneConfig()
})

function closeSidebar() {
  appStore.sidebarCollapsed = true
}
</script>

<template>
  <n-layout has-sider position="absolute" style="top: 0; left: 0; right: 0; bottom: 0">
    <n-layout-sider
      bordered
      :collapsed="appStore.sidebarCollapsed"
      :collapsed-width="isMobile ? 0 : 64"
      :width="240"
      :native-scrollbar="false"
      :collapse-mode="isMobile ? 'transform' : 'width'"
      :position="isMobile ? 'absolute' : 'static'"
      :show-trigger="false"
      :show-content="!isMobile || !appStore.sidebarCollapsed"
      class="app-sider"
      :class="{ 'sider-inactive': siderInactive }"
    >
      <div class="sidebar-logo">
        <div class="sidebar-logo-icon">
          <svg width="24" height="24" viewBox="0 0 36 36" fill="none">
            <rect width="36" height="36" rx="10" fill="var(--primary-color)" />
            <path d="M10 12h4l4 8 4-8h4v14h-4v-8l-4 8-4-8v8h-4V12z" fill="white" />
          </svg>
        </div>
        <span v-if="!appStore.sidebarCollapsed" class="logo-text">logcat</span>
        <span v-else class="logo-text-collapsed">GL</span>
      </div>
      <SideMenu />
    </n-layout-sider>

    <div
      v-if="isMobile && !appStore.sidebarCollapsed"
      class="sidebar-overlay"
      @click="closeSidebar"
    />

    <n-layout>
      <n-layout-header bordered>
        <TopHeader />
      </n-layout-header>

      <n-layout-content
        :native-scrollbar="false"
        :content-style="{ padding: isMobile ? '8px' : '16px 24px', minHeight: 'calc(100vh - 56px)' }"
      >
        <slot />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style scoped>
.app-sider {
  transition: all 0.2s ease;
  z-index: 10;
}

/* When the sider is hidden on mobile (collapsed with width=0) we must disable
   pointer events so it doesn't intercept clicks meant for the page content.
   A CSS class driven by Vue is more reliable than mutating inline styles from
   a watch because the class is applied synchronously with the collapse state. */
.app-sider.sider-inactive {
  pointer-events: none;
}

.sidebar-logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid var(--border-color);
  background: var(--sidebar-bg);
  padding: 0 16px;
  transition: background-color 0.3s ease, border-color 0.3s ease;
}

.sidebar-logo-icon {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-color);
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.logo-text-collapsed {
  font-size: 16px;
  font-weight: 700;
  color: var(--primary-color);
}

.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 9;
  cursor: pointer;
  backdrop-filter: blur(2px);
}
</style>
