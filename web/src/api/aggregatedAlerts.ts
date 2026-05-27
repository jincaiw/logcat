import http, { extractData } from './index'
import type { ApiResponse, AggregatedAlert, SyslogLog, PageResponse, PageParams } from '@/types'

export function getAggregatedAlerts(params: PageParams & { severity?: string; status?: string; startTime?: string; endTime?: string }): Promise<ApiResponse<PageResponse<AggregatedAlert>>> {
  return http.get('/aggregated-alerts', { params }).then(extractData)
}

export function getAggregatedAlert(id: number): Promise<ApiResponse<AggregatedAlert>> {
  return http.get(`/aggregated-alerts/${id}`).then(extractData)
}

export function getAggregatedAlertLogs(id: number, params: PageParams): Promise<ApiResponse<PageResponse<SyslogLog>>> {
  return http.get(`/aggregated-alerts/${id}/logs`, { params }).then(extractData)
}

export function acknowledgeAggregatedAlert(id: number): Promise<ApiResponse<null>> {
  return http.post(`/aggregated-alerts/${id}/acknowledge`).then(extractData)
}

export function resolveAggregatedAlert(id: number): Promise<ApiResponse<null>> {
  return http.post(`/aggregated-alerts/${id}/resolve`).then(extractData)
}