import http, { extractData } from './index'
import type { ApiResponse, PageResponse, AuditLog } from '@/types'

export function getAuditLogs(params: {
  page?: number
  pageSize?: number
  username?: string
  action?: string
  result?: string
  resourceType?: string
  startTime?: string
  endTime?: string
}): Promise<ApiResponse<PageResponse<AuditLog>>> {
  return http.get('/audit-logs', { params }).then(extractData)
}
