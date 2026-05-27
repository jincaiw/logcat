import http, { extractData } from './index'
import type { ApiResponse, PushConfig, PageResponse, PageParams } from '@/types'

export function getPushConfigs(params: PageParams & { name?: string; type?: string; status?: number }): Promise<ApiResponse<PageResponse<PushConfig>>> {
  return http.get('/push-configs', { params }).then(extractData)
}

export function getPushConfig(id: number): Promise<ApiResponse<PushConfig>> {
  return http.get(`/push-configs/${id}`).then(extractData)
}

export function createPushConfig(data: Partial<PushConfig>): Promise<ApiResponse<PushConfig>> {
  return http.post('/push-configs', data).then(extractData)
}

export function updatePushConfig(id: number, data: Partial<PushConfig>): Promise<ApiResponse<PushConfig>> {
  return http.put(`/push-configs/${id}`, data).then(extractData)
}

export function deletePushConfig(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/push-configs/${id}`).then(extractData)
}

export function testPushConfig(id: number): Promise<ApiResponse<{ success: boolean; message: string }>> {
  return http.post(`/push-configs/${id}/test`).then(extractData)
}