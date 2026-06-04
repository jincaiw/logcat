import http, { extractData } from './index'
import type { ApiResponse, AlertRecord, AlertDisposition, PageResponse, PageParams } from '@/types'

function mapAlertRecord(item: any): AlertRecord {
  return {
    id: item.id,
    logId: item.logId || '',
    alertRuleId: item.alertRuleId || null,
    pushConfigId: item.pushConfigId || null,
    channelType: item.channelType || '',
    status: item.status || 'success',
    retryCount: item.retryCount || 0,
    requestSummary: item.requestSummary || '',
    responseStatusCode: item.responseStatusCode || 0,
    responseSummary: item.responseSummary || '',
    errorMessage: item.errorMessage || '',
    dispositionStatus: item.dispositionStatus || '',
    sentAt: item.sentAt || null,
    createdAt: item.createdAt,
    alertRule: item.alertRule || null,
    pushConfig: item.pushConfig || null,
  }
}

export function getAlertRecords(params: PageParams & { severity?: string; status?: string; deviceName?: string; startTime?: string; endTime?: string }): Promise<ApiResponse<PageResponse<AlertRecord>>> {
  return http.get('/alerts', { params }).then(extractData).then((res) => ({
    ...res,
    data: {
      ...(res.data || {}),
      list: ((res.data?.items || res.data?.list || []) as any[]).map(mapAlertRecord),
      total: Number(res.data?.total || 0),
      page: Number(res.data?.page || params.page || 1),
      pageSize: Number(res.data?.pageSize || params.pageSize || 20),
    },
  }))
}

export function getAlertRecord(id: number): Promise<ApiResponse<AlertRecord>> {
  return http.get(`/alerts/${id}`).then(extractData).then((res) => ({
    ...res,
    data: mapAlertRecord(res.data || {}),
  }))
}

export function disposeAlert(alertId: number, data: { action: string; note?: string }): Promise<ApiResponse<AlertDisposition>> {
  const statusMap: Record<string, string> = {
    confirm: 'confirmed',
    ignore: 'ignored',
    close: 'closed',
  }
  return http.post(`/alerts/${alertId}/dispositions`, {
    status: statusMap[data.action] || data.action,
    note: data.note,
  }).then(extractData)
}

export function getAlertDispositions(params: PageParams & { alertId?: number; action?: string; operatorId?: number }): Promise<ApiResponse<PageResponse<AlertDisposition>>> {
  const request = params.alertId
    ? http.get(`/alerts/${params.alertId}/dispositions`)
    : http.get('/alert-dispositions', { params })

  return request.then(extractData).then((res) => ({
    ...res,
    data: {
      list: Array.isArray(res.data) ? res.data : (res.data?.items || res.data?.list || []),
      total: Array.isArray(res.data) ? res.data.length : Number(res.data?.total || 0),
      page: Number((res.data as any)?.page || params.page || 1),
      pageSize: Number((res.data as any)?.pageSize || params.pageSize || 20),
    },
  }))
}

export function batchDisposeAlerts(alertIds: number[], action: string, note?: string): Promise<ApiResponse<null>> {
  return Promise.all(alertIds.map((id) => disposeAlert(id, { action, note }))).then(() => ({
    code: 0,
    message: 'success',
    data: null,
  } as ApiResponse<null>))
}
