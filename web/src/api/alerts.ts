import http, { extractData } from './index'
import type { ApiResponse, AlertRecord, AlertDisposition, PageResponse, PageParams } from '@/types'

export function getAlertRecords(params: PageParams & { severity?: string; status?: string; deviceName?: string; startTime?: string; endTime?: string }): Promise<ApiResponse<PageResponse<AlertRecord>>> {
  return http.get('/alerts', { params }).then(extractData)
}

export function getAlertRecord(id: number): Promise<ApiResponse<AlertRecord>> {
  return http.get(`/alerts/${id}`).then(extractData)
}

export function disposeAlert(alertId: number, data: { action: string; note?: string }): Promise<ApiResponse<AlertDisposition>> {
  return http.post(`/alerts/${alertId}/dispose`, data).then(extractData)
}

export function getAlertDispositions(params: PageParams & { alertId?: number; action?: string; operatorId?: number }): Promise<ApiResponse<PageResponse<AlertDisposition>>> {
  return http.get('/alerts/dispositions', { params }).then(extractData)
}

export function batchDisposeAlerts(alertIds: number[], action: string, note?: string): Promise<ApiResponse<null>> {
  return http.post('/alerts/batch-dispose', { alertIds, action, note }).then(extractData)
}