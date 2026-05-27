import http, { extractData } from './index'
import type { ApiResponse, StatsQueryParams, StatsResult } from '@/types'

export function getStats(params: StatsQueryParams): Promise<ApiResponse<StatsResult[]>> {
  return http.get('/stats', { params }).then(extractData)
}

export function exportStatsCsv(params: StatsQueryParams): Promise<ApiResponse<{ url: string }>> {
  return http.get('/stats/export-csv', { params }).then(extractData)
}

export function copyIpList(params: { startTime?: string; endTime?: string }): Promise<ApiResponse<{ ips: string[] }>> {
  return http.get('/stats/ip-list', { params }).then(extractData)
}