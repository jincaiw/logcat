<script setup lang="ts">
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent } from 'naive-ui'
import { useAppStore } from '@/stores/app'
import SideMenu from './SideMenu.vue'
import TopHeader from './TopHeader.vue'

const appStore = useAppStore()
</script>

<template>
  <n-layout has-sider position="absolute" style="top: 0; left: 0; right: 0; bottom: 0">
    <n-layout-sider
      bordered
      :collapsed="appStore.sidebarCollapsed"
      :collapsed-width="64"
      :width="240"
      :native-scrollbar="false"
      collapse-mode="width"
    >
      <div class="sidebar-logo">
        <span v-if="!appStore.sidebarCollapsed" class="logo-text">GoLog</span>
        <span v-else class="logo-text-collapsed">GL</span>
      </div>
      <SideMenu />
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered>
        <TopHeader />
      </n-layout-header>

      <n-layout-content
        :native-scrollbar="false"
        content-style="padding: 16px 24px; min-height: calc(100vh - 56px)"
      >
        <slot />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style scoped>
.sidebar-logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-color-card);
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: var(--primary-color);
  letter-spacing: 2px;
}

.logo-text-collapsed {
  font-size: 18px;
  font-weight: 700;
  color: var(--primary-color);
}
</style>