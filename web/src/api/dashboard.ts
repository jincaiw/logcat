import http, { extractData } from './index'
import type { ApiResponse, DashboardStats } from '@/types'

export function getDashboardStats(): Promise<ApiResponse<DashboardStats>> {
  return http.get('/dashboard/stats').then(extractData)
}

export function getReceiveRateHistory(params?: { minutes?: number }): Promise<ApiResponse<{ times: string[]; rates: number[] }>> {
  return http.get('/dashboard/receive-rate', { params }).then(extractData)
}

export function getAlertTrend(params?: { days?: number }): Promise<ApiResponse<{ dates: string[]; counts: number[] }>> {
  return http.get('/dashboard/alert-trend', { params }).then(extractData)
}

export function getTopDevices(params?: { topN?: number }): Promise<ApiResponse<{ device: string; count: number }[]>> {
  return http.get('/dashboard/top-devices', { params }).then(extractData)
}