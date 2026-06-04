import http, { extractData } from './index'
import type { ApiResponse, HighFreqIp, HighFreqIpConfig, PageResponse, PageParams } from '@/types'

export function getHighFreqIps(params: PageParams & { startTime?: string; endTime?: string }): Promise<ApiResponse<PageResponse<HighFreqIp>>> {
  return http.get('/high-frequency-ips', { params }).then(extractData)
}

export function getHighFreqIpConfig(): Promise<ApiResponse<HighFreqIpConfig>> {
  return http.get('/high-frequency-ips/config').then(extractData)
}

export function updateHighFreqIpConfig(data: HighFreqIpConfig): Promise<ApiResponse<null>> {
  return http.put('/high-frequency-ips/config', data).then(extractData)
}

export function refreshHighFreqIps(): Promise<ApiResponse<null>> {
  return http.post('/high-frequency-ips/refresh').then(extractData)
}
