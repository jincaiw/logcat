import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { serviceApi, configApi, statsApi } from '@/api'
import type { Protocol, SystemStats } from '@/types'

export const useAppStore = defineStore('app', () => {
  const currentPageTitle = ref('')
  const serviceRunning = ref(false)
  const listenPort = ref(5140)
  const protocol = ref<Protocol>('both')
  const loading = ref(false)

  const stats = ref<SystemStats>({
    totalLogs: 0,
    deviceCount: 0,
    serviceRunning: false,
    listenPort: 5140,
    startTime: '',
    memoryUsage: 0,
    cpuUsage: 0,
    connections: 0,
    receiveRate: 0,
    protocol: 'both',
    databaseSize: 0,
  })

  const formattedStats = computed(() => stats.value)

  /** 初始化应用：加载配置与统计 */
  async function initApp() {
    try {
      const config = await configApi.get()
      listenPort.value = config.listenPort || 5140
      protocol.value = 'both'
      await refreshStats()
    } catch (error) {
      console.error('Failed to init app:', error)
    }
  }

  /** 刷新服务状态与系统统计 */
  async function refreshStats() {
    try {
      const [serviceStatus, systemStats] = await Promise.all([
        serviceApi.getStatus(),
        statsApi.getSystemStats(),
      ])

      serviceRunning.value = serviceStatus.serviceRunning

      stats.value = {
        ...systemStats,
        serviceRunning: serviceStatus.serviceRunning,
        listenPort: serviceStatus.listenPort,
        connections: serviceStatus.connections,
        receiveRate: serviceStatus.receiveRate,
      }
    } catch (error) {
      console.error('Failed to refresh stats:', error)
    }
  }

  /** 启动 Syslog 服务 */
  async function startService(port: number, proto: Protocol) {
    loading.value = true
    try {
      await serviceApi.start(port, proto)
      serviceRunning.value = true
      listenPort.value = port
      protocol.value = proto
      await refreshStats()
    } catch (error) {
      console.error('Failed to start service:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  /** 停止 Syslog 服务 */
  async function stopService() {
    loading.value = true
    try {
      await serviceApi.stop()
      serviceRunning.value = false
      await refreshStats()
    } catch (error) {
      console.error('Failed to stop service:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  function setPageTitle(title: string) {
    currentPageTitle.value = title
  }

  return {
    currentPageTitle,
    serviceRunning,
    stats,
    listenPort,
    protocol,
    loading,
    formattedStats,
    initApp,
    refreshStats,
    startService,
    stopService,
    setPageTitle,
  }
})
