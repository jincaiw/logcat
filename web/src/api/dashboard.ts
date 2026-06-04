import http, { extractData } from './index'
import type { ApiResponse, DashboardStats } from '@/types'

export function getDashboardStats(): Promise<ApiResponse<DashboardStats>> {
  return http.get('/dashboard/stats').then(extractData)
}
