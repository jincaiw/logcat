import http, { extractData } from './index'
import type { ApiResponse, DesensitizeRule, PageResponse, PageParams } from '@/types'

export function getDesensitizeRules(params: PageParams & { name?: string; status?: number }): Promise<ApiResponse<PageResponse<DesensitizeRule>>> {
  return http.get('/desensitize-rules', { params }).then(extractData)
}

export function getDesensitizeRule(id: number): Promise<ApiResponse<DesensitizeRule>> {
  return http.get(`/desensitize-rules/${id}`).then(extractData)
}

export function createDesensitizeRule(data: Partial<DesensitizeRule>): Promise<ApiResponse<DesensitizeRule>> {
  return http.post('/desensitize-rules', data).then(extractData)
}

export function updateDesensitizeRule(id: number, data: Partial<DesensitizeRule>): Promise<ApiResponse<DesensitizeRule>> {
  return http.put(`/desensitize-rules/${id}`, data).then(extractData)
}

export function deleteDesensitizeRule(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/desensitize-rules/${id}`).then(extractData)
}

export function toggleDesensitizeRule(id: number, status: number): Promise<ApiResponse<null>> {
  return http.put(`/desensitize-rules/${id}/status`, { status }).then(extractData)
}