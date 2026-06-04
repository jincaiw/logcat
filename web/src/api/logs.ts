import http, { extractData } from './index'
import type { ApiResponse, SyslogLog, PageResponse, PageParams } from '@/types'

export interface LogQueryParams extends PageParams {
  deviceName?: string
  deviceHost?: string
  sourceIp?: string
  destIp?: string
  eventType?: string
  severity?: string
  pushStatus?: string
  keyword?: string
  logId?: string
  startTime?: string
  endTime?: string
  parsedFieldKey?: string
  parsedFieldValue?: string
  filterStatus?: string
}

export function queryLogs(params: LogQueryParams): Promise<ApiResponse<PageResponse<SyslogLog>>> {
  return http.get('/logs', { params }).then(extractData)
}

export function getLogById(id: string): Promise<ApiResponse<SyslogLog>> {
  return http.get(`/logs/${id}`).then(extractData)
}

export function getLogTrace(id: string): Promise<ApiResponse<{ log: SyslogLog; trace: any[] }>> {
  return http.get(`/logs/${id}/trace`).then(extractData)
}

export function getUnmatchedLogCount(): Promise<ApiResponse<{ count: number }>> {
  return http.get('/logs/unmatched-count').then(extractData)
}

export function cleanupLogs(beforeTime?: string, days?: number): Promise<ApiResponse<{ deleted: number }>> {
  const data: Record<string, any> = {}
  if (days !== undefined) data.days = days
  else if (beforeTime) data.beforeTime = beforeTime
  return http.delete('/logs/cleanup', { data }).then(extractData)
}

export function exportLogs(params: LogQueryParams): Promise<ApiResponse<{ url: string }>> {
  return http.get('/logs/export', { params }).then(extractData)
}
