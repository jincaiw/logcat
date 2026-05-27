import http, { extractData } from './index'
import type { ApiResponse, FilterPolicy, FilterTestResult, PageResponse, PageParams } from '@/types'

export function getFilterPolicies(params: PageParams & { name?: string; status?: number }): Promise<ApiResponse<PageResponse<FilterPolicy>>> {
  return http.get('/filter-policies', { params }).then(extractData)
}

export function getFilterPolicy(id: number): Promise<ApiResponse<FilterPolicy>> {
  return http.get(`/filter-policies/${id}`).then(extractData)
}

export function createFilterPolicy(data: Partial<FilterPolicy>): Promise<ApiResponse<FilterPolicy>> {
  return http.post('/filter-policies', data).then(extractData)
}

export function updateFilterPolicy(id: number, data: Partial<FilterPolicy>): Promise<ApiResponse<FilterPolicy>> {
  return http.put(`/filter-policies/${id}`, data).then(extractData)
}

export function deleteFilterPolicy(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/filter-policies/${id}`).then(extractData)
}

export function testFilterPolicy(id: number, logData: Record<string, any>): Promise<ApiResponse<FilterTestResult>> {
  return http.post(`/filter-policies/${id}/test`, logData).then(extractData)
}