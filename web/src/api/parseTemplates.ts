import http, { extractData } from './index'
import type { ApiResponse, ParseTemplate, ParseTestResult, PageResponse, PageParams } from '@/types'

export function getParseTemplates(params: PageParams & { name?: string; type?: string }): Promise<ApiResponse<PageResponse<ParseTemplate>>> {
  return http.get('/parse-templates', { params }).then(extractData).then((res) => ({
    ...res,
    data: {
      list: Array.isArray(res.data) ? res.data : (res.data?.items || res.data?.list || []),
      total: Array.isArray(res.data) ? res.data.length : Number(res.data?.total || 0),
      page: params.page || 1,
      pageSize: params.pageSize || 20,
    },
  }))
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
  return http.post('/parse-templates/test', { templateId: id, rawLog: sample }).then(extractData)
}
