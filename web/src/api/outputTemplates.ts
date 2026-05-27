import http, { extractData } from './index'
import type { ApiResponse, OutputTemplate, PageResponse, PageParams } from '@/types'

export function getOutputTemplates(params: PageParams & { name?: string }): Promise<ApiResponse<PageResponse<OutputTemplate>>> {
  return http.get('/output-templates', { params }).then(extractData)
}

export function getOutputTemplate(id: number): Promise<ApiResponse<OutputTemplate>> {
  return http.get(`/output-templates/${id}`).then(extractData)
}

export function createOutputTemplate(data: Partial<OutputTemplate>): Promise<ApiResponse<OutputTemplate>> {
  return http.post('/output-templates', data).then(extractData)
}

export function updateOutputTemplate(id: number, data: Partial<OutputTemplate>): Promise<ApiResponse<OutputTemplate>> {
  return http.put(`/output-templates/${id}`, data).then(extractData)
}

export function deleteOutputTemplate(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/output-templates/${id}`).then(extractData)
}