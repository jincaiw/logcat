import http, { extractData } from './index'
import type { ApiResponse, SystemConfig, SystemStatus } from '@/types'

export function getSystemStatus(): Promise<ApiResponse<SystemStatus>> {
  return http.get('/system/status').then(extractData)
}

export function startSyslogService(): Promise<ApiResponse<null>> {
  return http.post('/system/start').then(extractData)
}

export function stopSyslogService(): Promise<ApiResponse<null>> {
  return http.post('/system/stop').then(extractData)
}

export function restartSyslogService(): Promise<ApiResponse<null>> {
  return http.post('/system/restart').then(extractData)
}

export function getSystemConfigs(): Promise<ApiResponse<SystemConfig[]>> {
  return http.get('/system/configs').then(extractData)
}

export function updateSystemConfigs(configs: Record<string, string>): Promise<ApiResponse<null>> {
  return http.put('/system/configs', configs).then(extractData)
}

export function getSystemConfig(key: string): Promise<ApiResponse<SystemConfig>> {
  return http.get(`/system/configs/${key}`).then(extractData)
}