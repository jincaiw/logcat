import http, { extractData } from './index'
import type { ApiResponse, HighFreqIp, HighFreqIpConfig, PageResponse, PageParams } from '@/types'

export function getHighFreqIps(params: PageParams & { startTime?: string; endTime?: string }): Promise<ApiResponse<PageResponse<HighFreqIp>>> {
  return http.get('/high-freq-ips', { params }).then(extractData)
}

export function getHighFreqIpConfig(): Promise<ApiResponse<HighFreqIpConfig>> {
  return http.get('/high-freq-ips/config').then(extractData)
}

export function updateHighFreqIpConfig(data: HighFreqIpConfig): Promise<ApiResponse<null>> {
  return http.put('/high-freq-ips/config', data).then(extractData)
}

export function refreshHighFreqIps(): Promise<ApiResponse<null>> {
  return http.post('/high-freq-ips/refresh').then(extractData)
}