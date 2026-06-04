<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NBreadcrumb, NBreadcrumbItem, NIcon } from 'naive-ui'
import { HomeOutline, ChevronForwardOutline } from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()

const breadcrumbs = computed(() => {
  const matched = route.matched.filter((r) => r.meta?.title)
  return matched.map((r) => ({
    title: r.meta?.title as string,
    path: r.path,
  }))
})

function navigateTo(path: string) {
  if (path !== route.path) {
    router.push(path)
  }
}
</script>

<template>
  <div class="breadcrumb-wrapper">
    <n-breadcrumb separator=">">
      <n-breadcrumb-item @click="navigateTo('/')">
        <n-icon :component="HomeOutline" :size="16" />
      </n-breadcrumb-item>
      <n-breadcrumb-item
        v-for="(item, index) in breadcrumbs"
        :key="index"
        @click="navigateTo(item.path)"
      >
        {{ item.title }}
      </n-breadcrumb-item>
    </n-breadcrumb>
  </div>
</template>

<style scoped>
.breadcrumb-wrapper {
  display: flex;
  align-items: center;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>