import axios from 'axios'
import http, { extractData } from './index'
import type {
  ApiResponse,
  RuntimeMetricsSnapshot,
  SystemHealthCheck,
  SystemMetricsSnapshot,
} from '@/types'

const monitorHttp = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true,
})

export function getHealthz(): Promise<SystemHealthCheck> {
  return monitorHttp.get('/healthz').then((res) => {
    const d = res.data
    return d.data ?? d
  })
}

export function getReadyz(): Promise<SystemHealthCheck> {
  return monitorHttp.get('/readyz').then((res) => {
    const d = res.data
    return d.data ?? d
  })
}

export function getMetricsSnapshot(): Promise<SystemMetricsSnapshot> {
  return monitorHttp.get('/metrics').then((res) => {
    const d = res.data
    return d.data ?? d
  })
}

export function getRuntimeMetrics(): Promise<ApiResponse<RuntimeMetricsSnapshot>> {
  return http.get('/metrics/runtime').then(extractData)
}
