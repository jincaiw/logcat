import http, { extractData } from './index'
import type { ApiResponse, ImportResult } from '@/types'

export type ExportType = 'device-templates' | 'parse-templates' | 'filter-policies' | 'push-configs'

export interface ExportResult {
  url: string
  version: string
  resourceType: string
  count: number
}

export function exportData(type: ExportType): Promise<ApiResponse<ExportResult>> {
  return http.get(`/export/${type}`).then(extractData)
}

export function importData(type: ExportType, content: string): Promise<ApiResponse<ImportResult>> {
  return http.post(`/import/${type}`, content, {
    headers: { 'Content-Type': 'application/json' },
  }).then(extractData)
}

export function getExportHistory(params?: { type?: string; page?: number; pageSize?: number }): Promise<ApiResponse<any>> {
  return http.get('/export/history', { params }).then(extractData)
}

export function getImportHistory(params?: { type?: string; page?: number; pageSize?: number }): Promise<ApiResponse<any>> {
  return http.get('/import/history', { params }).then(extractData)
}
