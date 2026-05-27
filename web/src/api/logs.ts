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

export function cleanupLogs(beforeTime?: string): Promise<ApiResponse<{ deleted: number }>> {
  return http.post('/logs/cleanup', { beforeTime }).then(extractData)
}

export function exportLogs(params: LogQueryParams): Promise<ApiResponse<{ url: string }>> {
  return http.get('/logs/export', { params }).then(extractData)
}