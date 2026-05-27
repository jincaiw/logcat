import http, { extractData } from './index'
import type { ApiResponse, ParseTemplate, ParseTestResult, PageResponse, PageParams } from '@/types'

export function getParseTemplates(params: PageParams & { name?: string; type?: string }): Promise<ApiResponse<PageResponse<ParseTemplate>>> {
  return http.get('/parse-templates', { params }).then(extractData)
}

export function getParseTemplate(id: number): Promise<ApiResponse<ParseTemplate>> {
  return http.get(`/parse-templates/${id}`).then(extractData)
}

export function createParseTemplate(data: Partial<ParseTemplate>): Promise<ApiResponse<ParseTemplate>> {
  return http.post('/parse-templates', data).then(extractData)
}

export function updateParseTemplate(id: number, data: Partial<ParseTemplate>): Promise<ApiResponse<ParseTemplate>> {
  return http.put(`/parse-templates/${id}`, data).then(extractData)
}

export function deleteParseTemplate(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/parse-templates/${id}`).then(extractData)
}

export function testParseTemplate(id: number, sample: string): Promise<ApiResponse<ParseTestResult>> {
  return http.post(`/parse-templates/${id}/test`, { sample }).then(extractData)
}