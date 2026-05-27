import http, { extractData } from './index'
import type { ApiResponse, ImportResult } from '@/types'

export type ExportType = 'devices' | 'deviceGroups' | 'deviceTemplates' | 'parseTemplates' | 'filterPolicies' | 'outputTemplates' | 'pushConfigs' | 'alertRules' | 'desensitizeRules' | 'roles'

export function exportData(type: ExportType): Promise<ApiResponse<{ url: string }>> {
  return http.get(`/export/${type}`).then(extractData)
}

export function importData(type: ExportType, file: File): Promise<ApiResponse<ImportResult>> {
  const formData = new FormData()
  formData.append('file', file)
  return http.post(`/import/${type}`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }).then(extractData)
}

export function getExportHistory(params?: { type?: string; page?: number; pageSize?: number }): Promise<ApiResponse<any>> {
  return http.get('/export/history', { params }).then(extractData)
}

export function getImportHistory(params?: { type?: string; page?: number; pageSize?: number }): Promise<ApiResponse<any>> {
  return http.get('/import/history', { params }).then(extractData)
}