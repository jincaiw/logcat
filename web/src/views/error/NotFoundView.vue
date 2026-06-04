<script setup lang="ts">
import { useRouter } from 'vue-router'
import { NButton, NResult } from 'naive-ui'

const router = useRouter()
</script>

<template>
  <div class="not-found-container">
    <n-result
      status="404"
      title="404 页面未找到"
      description="您访问的页面不存在或已被移除"
    >
      <template #footer>
        <n-button type="primary" @click="router.push('/')">
          返回首页
        </n-button>
      </template>
    </n-result>
  </div>
</template>

<style scoped>
.not-found-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 24px;
}
</style>
