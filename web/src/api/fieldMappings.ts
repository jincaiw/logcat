import http, { extractData } from './index'
import type { ApiResponse, FieldMappingDoc, PageResponse, PageParams } from '@/types'

export function getFieldMappings(params: PageParams & { name?: string }): Promise<ApiResponse<PageResponse<FieldMappingDoc>>> {
  return http.get('/field-mappings', { params }).then(extractData)
}

export function getFieldMapping(id: number): Promise<ApiResponse<FieldMappingDoc>> {
  return http.get(`/field-mappings/${id}`).then(extractData)
}

export function createFieldMapping(data: Partial<FieldMappingDoc>): Promise<ApiResponse<FieldMappingDoc>> {
  return http.post('/field-mappings', data).then(extractData)
}

export function updateFieldMapping(id: number, data: Partial<FieldMappingDoc>): Promise<ApiResponse<FieldMappingDoc>> {
  return http.put(`/field-mappings/${id}`, data).then(extractData)
}

export function deleteFieldMapping(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/field-mappings/${id}`).then(extractData)
}