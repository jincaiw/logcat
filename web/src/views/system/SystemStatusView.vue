<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  NCard, NButton, NSpace, NGrid, NGi, NTag, NIcon, NSpin, NTable, useMessage,
} from 'naive-ui'
import {
  PlayOutline, StopOutline, RefreshOutline,
  ServerOutline,
} from '@vicons/ionicons5'
import { getSystemStatus, startSyslogService, stopSyslogService } from '@/api/system'
import type { SystemStatus } from '@/types'
import PageHeader from '@/components/common/PageHeader.vue'

const message = useMessage()

const status = ref<SystemStatus | null>(null)
const loading = ref(false)
const actionLoading = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

async function loadStatus() {
  loading.value = true
  try {
    const res = await getSystemStatus()
    status.value = res.data
  } catch {
    message.error('获取状态失败')
  } finally {
    loading.value = false
  }
}

async function handleStart() {
  actionLoading.value = true
  try {
    await startSyslogService()
    message.success('服务已启动')
    await loadStatus()
  } catch {
    message.error('启动失败')
  } finally {
    actionLoading.value = false
  }
}

async function handleStop() {
  actionLoading.value = true
  try {
    await stopSyslogService()
    message.success('服务已停止')
    await loadStatus()
  } catch {
    message.error('停止失败')
  } finally {
    actionLoading.value = false
  }
}

onMounted(() => {
  loadStatus()
  timer = setInterval(loadStatus, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="page-container">
    <PageHeader title="服务状态" description="Syslog 服务运行状态与控制">
      <n-space>
        <n-button
          v-if="!status?.serviceRunning"
          type="success"
          :loading="actionLoading"
          @click="handleStart"
        >
          <template #icon><n-icon :component="PlayOutline" /></template>
          启动服务
        </n-button>
        <n-button
          v-else
          type="error"
          :loading="actionLoading"
          @click="handleStop"
        >
          <template #icon><n-icon :component="StopOutline" /></template>
          停止服务
        </n-button>
        <n-button quaternary @click="loadStatus">
          <template #icon><n-icon :component="RefreshOutline" /></template>
        </n-button>
      </n-space>
    </PageHeader>

    <n-spin :show="loading">
      <!-- Service Status -->
      <n-grid cols="1 s:2 m:4" x-gap="12" y-gap="12" style="margin-bottom: 16px">
        <n-gi>
          <n-card size="small">
            <n-statistic label="服务状态">
              <n-tag :type="status?.serviceRunning ? 'success' : 'default'">
                {{ status?.serviceRunning ? '运行中' : '已停止' }}
              </n-tag>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="运行时长" :value="status?.uptime || '--'" />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="启动时间" :value="status?.startedAt || '--'" />
          </n-card>
        </n-gi>
      </n-grid>

      <!-- Listeners -->
      <n-card title="监听信息" size="small" style="margin-bottom: 16px">
        <n-table v-if="status?.listeners?.length" :single-line="false" size="small">
          <thead>
            <tr>
              <th>协议</th>
              <th>地址</th>
              <th>端口</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(listener, i) in status.listeners" :key="i">
              <td>{{ listener.protocol.toUpperCase() }}</td>
              <td>{{ listener.address }}</td>
              <td>{{ listener.port }}</td>
              <td><n-tag size="small" type="success">{{ listener.status }}</n-tag></td>
            </tr>
          </tbody>
        </n-table>
        <div v-else style="color: var(--text-color-secondary)">暂无监听信息</div>
      </n-card>

      <!-- Connection Status -->
      <n-card title="连接状态" size="small" style="margin-bottom: 16px">
        <n-grid v-if="status?.connections" cols="4" x-gap="12">
          <n-gi>
            <n-statistic label="总数" :value="status.connections.total" />
          </n-gi>
          <n-gi>
            <n-statistic label="活跃" :value="status.connections.active" />
          </n-gi>
          <n-gi>
            <n-statistic label="空闲" :value="status.connections.idle" />
          </n-gi>
          <n-gi>
            <n-statistic label="关闭" :value="status.connections.closed" />
          </n-gi>
        </n-grid>
        <div v-else style="color: var(--text-color-secondary)">暂无连接信息</div>
      </n-card>

      <!-- Queue Status -->
      <n-card title="队列状态" size="small" style="margin-bottom: 16px">
        <n-table v-if="status?.queue?.name" :single-line="false" size="small">
          <thead>
            <tr>
              <th>队列名称</th>
              <th>大小</th>
              <th>容量</th>
              <th>入队速率</th>
              <th>出队速率</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>{{ status.queue.name }}</td>
              <td>{{ status.queue.size }}</td>
              <td>{{ status.queue.capacity }}</td>
              <td>{{ status.queue.enqueueRate?.toFixed(1) }}/s</td>
              <td>{{ status.queue.dequeueRate?.toFixed(1) }}/s</td>
            </tr>
          </tbody>
        </n-table>
        <div v-else style="color: var(--text-color-secondary)">暂无队列信息</div>
      </n-card>

      <!-- Worker Status -->
      <n-card title="工作器状态" size="small">
        <n-table v-if="status?.workers?.length" :single-line="false" size="small">
          <thead>
            <tr>
              <th>ID</th>
              <th>状态</th>
              <th>已处理</th>
              <th>错误数</th>
              <th>最后活跃</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="worker in status.workers" :key="worker.id">
              <td>{{ worker.id }}</td>
              <td><n-tag size="small" :type="worker.status === 'running' ? 'success' : 'default'">{{ worker.status }}</n-tag></td>
              <td>{{ worker.processedCount }}</td>
              <td>{{ worker.errorCount }}</td>
              <td>{{ worker.lastActiveAt }}</td>
            </tr>
          </tbody>
        </n-table>
        <div v-else style="color: var(--text-color-secondary)">暂无工作器信息</div>
      </n-card>
    </n-spin>
  </div>
</template>

<script lang="ts">
import { NStatistic } from 'naive-ui'
</script>