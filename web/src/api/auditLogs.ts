import http, { extractData } from './index'
import type { ApiResponse, AuditLog, PageResponse, PageParams } from '@/types'

export function getAuditLogs(params: PageParams & { username?: string; action?: string; resource?: string; result?: string; startTime?: string; endTime?: string }): Promise<ApiResponse<PageResponse<AuditLog>>> {
  return http.get('/audit-logs', { params }).then(extractData)
}

export function getAuditLog(id: number): Promise<ApiResponse<AuditLog>> {
  return http.get(`/audit-logs/${id}`).then(extractData)
}

export function exportAuditLogs(params: { startTime?: string; endTime?: string }): Promise<ApiResponse<{ url: string }>> {
  return http.get('/audit-logs/export', { params }).then(extractData)
}