<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { NButton, NCard, NGi, NGrid, NIcon, NSpin, NStatistic, NSpace, NTable, NTag } from 'naive-ui'
import { RefreshOutline } from '@vicons/ionicons5'
import { getHealthz, getMetricsSnapshot, getReadyz, getRuntimeMetrics } from '@/api/monitoring'
import type { RuntimeMetricsSnapshot, SystemHealthCheck, SystemMetricsSnapshot } from '@/types'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAppMessage } from '@/composables/useMessage'
import { formatDuration } from '@/api/system'

const message = useAppMessage()

const loading = ref(false)
const health = ref<SystemHealthCheck | null>(null)
const ready = ref<SystemHealthCheck | null>(null)
const metrics = ref<SystemMetricsSnapshot | null>(null)
const runtime = ref<RuntimeMetricsSnapshot | null>(null)

let timer: ReturnType<typeof setInterval> | null = null

const dbStatsRows = computed(() => {
  const stats = metrics.value?.db_stats || {}
  return [
    { label: '打开连接', value: stats.open_connections ?? 0 },
    { label: '使用中', value: stats.in_use ?? 0 },
    { label: '空闲', value: stats.idle ?? 0 },
    { label: '等待次数', value: stats.wait_count ?? 0 },
    { label: '最大连接', value: stats.max_open ?? 0 },
  ]
})

const receiverRows = computed(() => {
  const value = runtime.value?.receiver || {}
  return [
    { label: 'UDP 接收', value: value.udpReceived ?? 0 },
    { label: 'TCP 接收', value: value.tcpReceived ?? 0 },
    { label: 'UDP 错误', value: value.udpErrors ?? 0 },
    { label: 'TCP 错误', value: value.tcpErrors ?? 0 },
    { label: '解析错误', value: value.parseErrors ?? 0 },
    { label: '通道丢弃', value: value.channelDropped ?? 0 },
    { label: 'TCP 连接数', value: value.tcpConnections ?? 0 },
    { label: '最后接收时间', value: formatTimestamp(value.lastReceiveAt) },
  ]
})

const pipelineRows = computed(() => {
  const value = runtime.value?.pipeline || {}
  return [
    { label: '原始队列', value: value.rawQueueDepth ?? 0 },
    { label: '解析队列', value: value.parsedQueueDepth ?? 0 },
    { label: '入库队列', value: value.dbQueueDepth ?? 0 },
    { label: '推送队列', value: value.pushQueueDepth ?? 0 },
    { label: '解析处理数', value: value.parseProcessed ?? 0 },
    { label: '解析错误数', value: value.parseErrors ?? 0 },
    { label: '筛选处理数', value: value.filterProcessed ?? 0 },
    { label: '筛掉数', value: value.filterDropped ?? 0 },
    { label: '入库成功数', value: value.dbWritten ?? 0 },
    { label: '入库错误数', value: value.dbErrors ?? 0 },
    { label: '推送处理数', value: value.pushProcessed ?? 0 },
    { label: '推送错误数', value: value.pushErrors ?? 0 },
    { label: '原始丢弃数', value: value.rawDropped ?? 0 },
    { label: '入库丢弃数', value: value.dbDropped ?? 0 },
    { label: '推送丢弃数', value: value.pushDropped ?? 0 },
  ]
})

async function loadData() {
  loading.value = true
  try {
    const [healthValue, readyValue, metricsValue, runtimeValue] = await Promise.all([
      getHealthz(),
      getReadyz(),
      getMetricsSnapshot(),
      getRuntimeMetrics(),
    ])
    health.value = healthValue
    ready.value = readyValue
    metrics.value = metricsValue
    runtime.value = runtimeValue.data
  } catch (error: any) {
    message.error(error?.message || '获取指标失败')
  } finally {
    loading.value = false
  }
}

function formatNumber(value?: number) {
  if (typeof value !== 'number') return '--'
  return Number.isInteger(value) ? String(value) : value.toFixed(2)
}

function formatTimestamp(value?: number | string) {
  if (!value) return '--'
  if (typeof value === 'string') return value
  return new Date(value / 1_000_000).toLocaleString('zh-CN', { hour12: false })
}

function statusTagType(value?: string) {
  return value === 'ok' || value === 'ready' || value === 'running' ? 'success' : 'warning'
}

onMounted(() => {
  loadData()
  timer = setInterval(loadData, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="page-container">
    <PageHeader title="指标监控" description="健康检查、数据库连接与运行时指标">
      <n-space>
        <n-button quaternary @click="loadData">
          <template #icon><n-icon :component="RefreshOutline" /></template>
          刷新
        </n-button>
      </n-space>
    </PageHeader>

    <n-spin :show="loading">
      <n-grid cols="1 s:2 m:3 l:6" x-gap="12" y-gap="12" style="margin-bottom: 16px">
        <n-gi>
          <n-card size="small">
            <n-statistic label="健康检查">
              <n-tag :type="statusTagType(health?.status)">{{ health?.status || '--' }}</n-tag>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="就绪检查">
              <n-tag :type="statusTagType(ready?.status)">{{ ready?.status || '--' }}</n-tag>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="数据库状态">
              <n-tag :type="statusTagType(metrics?.database)">{{ metrics?.database || '--' }}</n-tag>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="运行时长" :value="formatDuration(runtime?.uptime_seconds || 0)" />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="Goroutines" :value="runtime?.goroutines || 0" />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card size="small">
            <n-statistic label="Heap(MB)" :value="formatNumber(runtime?.heap_alloc_mb)" />
          </n-card>
        </n-gi>
      </n-grid>

      <n-grid cols="1 l:2" x-gap="12" y-gap="12" style="margin-bottom: 16px">
        <n-gi>
          <n-card title="数据库连接" size="small">
            <div style="overflow-x: auto">
            <n-table :single-line="false" size="small">
              <tbody>
                <tr v-for="item in dbStatsRows" :key="item.label">
                  <td style="width: 40%">{{ item.label }}</td>
                  <td>{{ item.value }}</td>
                </tr>
              </tbody>
            </n-table>
            </div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card title="运行时概览" size="small">
            <div style="overflow-x: auto">
            <n-table :single-line="false" size="small">
              <tbody>
                <tr>
                  <td style="width: 40%">启动时间</td>
                  <td>{{ metrics?.started_at || '--' }}</td>
                </tr>
                <tr>
                  <td>累计运行秒数</td>
                  <td>{{ metrics?.uptime || 0 }}</td>
                </tr>
                <tr>
                  <td>GC 次数</td>
                  <td>{{ runtime?.num_gc || 0 }}</td>
                </tr>
                <tr>
                  <td>Heap Sys(MB)</td>
                  <td>{{ formatNumber(runtime?.heap_sys_mb) }}</td>
                </tr>
                <tr>
                  <td>Heap Total(MB)</td>
                  <td>{{ formatNumber(runtime?.heap_total_mb) }}</td>
                </tr>
                <tr>
                  <td>运行状态</td>
                  <td>{{ runtime?.status || '--' }}</td>
                </tr>
              </tbody>
            </n-table>
            </div>
          </n-card>
        </n-gi>
      </n-grid>

      <n-grid cols="1 l:2" x-gap="12" y-gap="12">
        <n-gi>
          <n-card title="接收器指标" size="small">
            <div style="overflow-x: auto">
            <n-table :single-line="false" size="small">
              <tbody>
                <tr v-for="item in receiverRows" :key="item.label">
                  <td style="width: 40%">{{ item.label }}</td>
                  <td>{{ item.value }}</td>
                </tr>
              </tbody>
            </n-table>
            </div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card title="流水线指标" size="small">
            <div style="overflow-x: auto">
            <n-table :single-line="false" size="small">
              <tbody>
                <tr v-for="item in pipelineRows" :key="item.label">
                  <td style="width: 40%">{{ item.label }}</td>
                  <td>{{ item.value }}</td>
                </tr>
              </tbody>
            </n-table>
            </div>
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>
