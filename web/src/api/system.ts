import http, { extractData } from './index'
import type { ApiResponse, SystemConfig, SystemStatus } from '@/types'

function formatDuration(seconds?: number): string {
  if (!seconds) return '0s'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  if (h > 0) return `${h}h ${m}m ${s}s`
  if (m > 0) return `${m}m ${s}s`
  return `${s}s`
}

export function getSystemStatus(): Promise<ApiResponse<SystemStatus>> {
  return http.get('/system/status').then(extractData)
}

export function startSyslogService(): Promise<ApiResponse<null>> {
  return http.post('/system/syslog/start').then(extractData)
}

export function stopSyslogService(): Promise<ApiResponse<null>> {
  return http.post('/system/syslog/stop').then(extractData)
}

export function restartSyslogService(): Promise<ApiResponse<null>> {
  return stopSyslogService().then(() => startSyslogService())
}

export function getSystemConfigs(): Promise<ApiResponse<SystemConfig[]>> {
  return http.get('/system/config').then(extractData).then((res) => {
    const raw = res.data
    if (Array.isArray(raw)) return res
    if (raw?.list && Array.isArray(raw.list)) {
      return { ...res, data: raw.list }
    }
    if (raw?.items && Array.isArray(raw.items)) {
      return { ...res, data: raw.items }
    }
    return { ...res, data: [] }
  })
}

export function updateSystemConfigs(configs: Record<string, string>): Promise<ApiResponse<null>> {
  return http.put('/system/config', {
    configs: Object.entries(configs).map(([configKey, configValue]) => ({ configKey, configValue })),
  }).then(extractData)
}

export function getSystemConfig(key: string): Promise<ApiResponse<SystemConfig>> {
  return getSystemConfigs().then((res) => ({
    ...res,
    data: (res.data || []).find((item) => item.configKey === key) as SystemConfig,
  }))
}

export { formatDuration }
