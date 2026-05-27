import http, { extractData } from './index'
import type { ApiResponse, AlertRule, PageResponse, PageParams } from '@/types'

export function getAlertRules(params: PageParams & { name?: string; severity?: string; status?: number }): Promise<ApiResponse<PageResponse<AlertRule>>> {
  return http.get('/alert-rules', { params }).then(extractData)
}

export function getAlertRule(id: number): Promise<ApiResponse<AlertRule>> {
  return http.get(`/alert-rules/${id}`).then(extractData)
}

export function createAlertRule(data: Partial<AlertRule>): Promise<ApiResponse<AlertRule>> {
  return http.post('/alert-rules', data).then(extractData)
}

export function updateAlertRule(id: number, data: Partial<AlertRule>): Promise<ApiResponse<AlertRule>> {
  return http.put(`/alert-rules/${id}`, data).then(extractData)
}

export function deleteAlertRule(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/alert-rules/${id}`).then(extractData)
}

export function toggleAlertRule(id: number, status: number): Promise<ApiResponse<null>> {
  return http.put(`/alert-rules/${id}/status`, { status }).then(extractData)
}