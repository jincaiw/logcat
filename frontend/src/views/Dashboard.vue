<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { NButton, NProgress, NSpace, useMessage } from 'naive-ui'
import { useAppStore } from '@/stores/app'
import { API } from '@/api'
import { useI18n } from '@/i18n'

const appStore = useAppStore()
const message = useMessage()
const { t } = useI18n()

const localIP = ref('')
const port = ref(5140)
const displayProtocol = computed(() => (appStore.protocol === 'both' ? 'TCP+UDP' : appStore.protocol.toUpperCase()))
let refreshTimer: ReturnType<typeof setInterval> | null = null

const systemStats = ref({
  totalLogs: 0,
  deviceCount: 0,
  matchedLogs: 0,
  alertCount: 0,
  parseTemplateCount: 0,
  activeFilterPolicies: 0,
  activeAlertPolicies: 0,
  activeRobots: 0,
  memoryUsage: 0,
  cpuUsage: 0,
  goroutineCount: 0,
  connections: 0,
  receiveRate: 0,
  databaseSize: 0,
  activeDevices: 0,
})

onMounted(async () => {
  await appStore.refreshStats()
  try {
    const ips = await API.GetLocalIPs()
    if (Array.isArray(ips)) {
      const preferredIP = ips.find((ip: string) => ip.startsWith('10.')) ||
                         ips.find((ip: string) => ip.startsWith('192.168.')) ||
                         ips.find((ip: string) => ip.startsWith('172.')) ||
                         ips[0] || '127.0.0.1'
      localIP.value = preferredIP
    } else {
      localIP.value = ips
    }
  } catch {
    localIP.value = '127.0.0.1'
  }
  port.value = appStore.listenPort
  await loadSystemStats()
  refreshTimer = setInterval(() => {
    appStore.refreshStats()
    loadSystemStats()
  }, 5000)
})

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})

async function loadSystemStats() {
  try {
    const stats = await API.GetSystemStats()
    systemStats.value = {
      totalLogs: stats.totalLogs || 0,
      deviceCount: stats.deviceCount || 0,
      matchedLogs: stats.matchedLogs || 0,
      alertCount: stats.alertCount || 0,
      parseTemplateCount: stats.parseTemplateCount || 0,
      activeFilterPolicies: stats.activeFilterPolicies || 0,
      activeAlertPolicies: stats.activeAlertPolicies || 0,
      activeRobots: stats.activeRobots || 0,
      memoryUsage: stats.memoryUsage || 0,
      cpuUsage: stats.cpuUsage || 0,
      goroutineCount: stats.goroutineCount || 0,
      connections: stats.connections || 0,
      receiveRate: stats.receiveRate || 0,
      databaseSize: stats.databaseSize || 0,
      activeDevices: stats.activeDevices || 0,
    }
  } catch (e) {
    console.error(e)
  }
}

async function handleStart() {
  try {
    await appStore.startService(port.value, appStore.protocol || 'both')
    message.success(t('message.serviceStarted'))
  } catch (error: any) {
    message.error(t('message.startFailed') + (error.message || error))
  }
}

async function handleStop() {
  try {
    await appStore.stopService()
    message.success(t('message.serviceStopped'))
  } catch (error: any) {
    message.error(t('message.stopFailed') + (error.message || error))
  }
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function getMetricColor(percent: number): string {
  if (percent > 90) return 'var(--danger)'
  if (percent > 70) return 'var(--warning)'
  return 'var(--accent)'
}

const statCards = computed(() => [
  { key: 'total', label: t('service.totalLogs'), value: systemStats.value.totalLogs, tone: 'blue' },
  { key: 'matched', label: t('dashboard.matchedLogsShort'), value: systemStats.value.matchedLogs, tone: 'green' },
  { key: 'alert', label: t('dashboard.alertTimes'), value: systemStats.value.alertCount, tone: 'red' },
  { key: 'device', label: t('service.deviceCount'), value: systemStats.value.deviceCount, tone: 'violet' },
])

const metrics = ref([
  { key: 'memory', label: t('service.memoryUsage'), display: '', percent: 0 },
  { key: 'cpu', label: t('service.cpuUsage'), display: '', percent: 0 },
  { key: 'goroutine', label: t('dashboard.goroutines'), display: '', percent: 0 },
  { key: 'rate', label: t('dashboard.processRate'), display: '', percent: 0 },
  { key: 'db', label: t('service.databaseSize'), display: '', percent: 0 },
  { key: 'devices', label: t('dashboard.activeServers'), display: '', percent: 0 },
])

watch(systemStats, () => {
  const s = systemStats.value
  metrics.value = [
    { key: 'memory', label: t('service.memoryUsage'), display: `${s.memoryUsage} MB`, percent: Math.min(s.memoryUsage / 500 * 100, 100) },
    { key: 'cpu', label: t('service.cpuUsage'), display: `${s.cpuUsage.toFixed(1)}%`, percent: Math.min(s.cpuUsage, 100) },
    { key: 'goroutine', label: t('dashboard.goroutines'), display: `${s.goroutineCount}`, percent: Math.min(s.goroutineCount / 100 * 100, 100) },
    { key: 'rate', label: t('dashboard.processRate'), display: `${s.receiveRate.toFixed(1)}${t('common.perSecond')}`, percent: Math.min(s.receiveRate, 100) },
    { key: 'db', label: t('service.databaseSize'), display: formatBytes(s.databaseSize), percent: Math.min(s.databaseSize / 524288000 * 100, 100) },
    { key: 'devices', label: t('dashboard.activeServers'), display: `${s.activeDevices} ${t('common.deviceUnit')}`, percent: Math.min(s.activeDevices / 50 * 100, 100) },
  ]
}, { deep: true, immediate: true })
</script>

<template>
  <div class="page dashboard-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('menu.dashboard') }}</h1>
        <p class="page-subtitle text-muted">{{ t('dashboard.subtitle') }}</p>
      </div>
      <NSpace class="page-actions">
        <NButton :disabled="appStore.serviceRunning" type="primary" @click="handleStart">
          {{ t('service.start') }}
        </NButton>
        <NButton :disabled="!appStore.serviceRunning" @click="handleStop">
          {{ t('service.stop') }}
        </NButton>
      </NSpace>
    </div>

    <!-- Service Hero -->
    <div class="service-hero panel-card" :class="{ active: appStore.serviceRunning }">
      <div class="hero-status">
        <span class="status-dot" :class="{ running: appStore.serviceRunning }"></span>
        <div class="status-text-wrap">
          <span class="status-text">{{ appStore.serviceRunning ? t('service.running') : t('service.stopped') }}</span>
          <span class="status-sub text-muted">{{ appStore.serviceRunning ? t('dashboard.serviceListening') : t('dashboard.serviceIdle') }}</span>
        </div>
      </div>
      <div class="hero-info">
        <div class="info-item">
          <span class="info-label text-muted">{{ t('service.listenAddress') }}</span>
          <span class="info-value mono">{{ localIP || '127.0.0.1' }}</span>
        </div>
        <div class="info-item">
          <span class="info-label text-muted">{{ t('service.listenPort') }}</span>
          <span class="info-value mono">{{ port }}</span>
        </div>
        <div class="info-item">
          <span class="info-label text-muted">{{ t('service.protocol') }}</span>
          <span class="info-value mono">{{ displayProtocol }}</span>
        </div>
      </div>
    </div>

    <!-- Stat Cards -->
    <div class="stats-grid mt-4">
      <div v-for="card in statCards" :key="card.key" class="stat-card">
        <div class="stat-label text-muted">{{ card.label }}</div>
        <div class="stat-value" :class="`value-${card.tone}`">{{ card.value.toLocaleString() }}</div>
      </div>
    </div>

    <!-- Resource Metrics -->
    <div class="mt-6">
      <h3 class="card-title text-secondary">{{ t('dashboard.resourceUsage') }}</h3>
    </div>
    <div class="metrics-grid mt-4">
      <div v-for="m in metrics" :key="m.key" class="metric-bar">
        <div class="metric-bar-header">
          <span class="metric-bar-label text-muted">{{ m.label }}</span>
          <span class="metric-bar-value">{{ m.display }}</span>
        </div>
        <NProgress
          type="line"
          :percentage="m.percent"
          :color="getMetricColor(m.percent)"
          :show-indicator="false"
          :height="6"
          :border-radius="3"
        />
      </div>
    </div>
  </div>
</template>

