import http, { extractData } from './index'
import type { ApiResponse, AggregatedAlert, SyslogLog, PageResponse, PageParams } from '@/types'

function mapAggregatedAlert(item: any): AggregatedAlert {
  return {
    id: item.id,
    aggregateKey: item.aggregateKey || item.aggregateType || '',
    aggregateType: item.aggregateType || '',
    sourceIp: item.sourceIp || '',
    destinationIp: item.destinationIp || '',
    eventType: item.eventType || '',
    deviceId: item.deviceId || null,
    severity: item.severity || 'info',
    count: item.count || 0,
    firstSeenAt: item.firstSeenAt || item.firstAt || '',
    lastSeenAt: item.lastSeenAt || item.lastAt || '',
    status: item.status || 'active',
    createdAt: item.createdAt,
    updatedAt: item.updatedAt,
  }
}

export function getAggregatedAlerts(params: PageParams & { severity?: string; status?: string; startTime?: string; endTime?: string }): Promise<ApiResponse<PageResponse<AggregatedAlert>>> {
  return http.get('/aggregated-alerts', { params }).then(extractData).then((res) => ({
    ...res,
    data: {
      ...(res.data || {}),
      list: ((res.data?.items || res.data?.list || []) as any[]).map(mapAggregatedAlert),
      total: Number(res.data?.total || 0),
      page: Number(res.data?.page || params.page || 1),
      pageSize: Number(res.data?.pageSize || params.pageSize || 20),
    },
  }))
}

export function getAggregatedAlert(id: number): Promise<ApiResponse<AggregatedAlert>> {
  return http.get(`/aggregated-alerts/${id}`).then(extractData).then((res) => ({
    ...res,
    data: mapAggregatedAlert(res.data || {}),
  }))
}

export function getAggregatedAlertLogs(id: number, params: PageParams): Promise<ApiResponse<PageResponse<SyslogLog>>> {
  return http.get(`/aggregated-alerts/${id}/logs`, { params }).then(extractData).then((res) => ({
    ...res,
    data: {
      list: Array.isArray(res.data) ? res.data : (res.data?.list || res.data?.items || []),
      total: Array.isArray(res.data) ? res.data.length : Number(res.data?.total || 0),
      page: params.page || 1,
      pageSize: params.pageSize || 20,
    },
  }))
}

export function acknowledgeAggregatedAlert(id: number): Promise<ApiResponse<null>> {
  return http.post(`/aggregated-alerts/${id}/acknowledge`).then(extractData)
}

export function resolveAggregatedAlert(id: number): Promise<ApiResponse<null>> {
  return http.post(`/aggregated-alerts/${id}/resolve`).then(extractData)
}
