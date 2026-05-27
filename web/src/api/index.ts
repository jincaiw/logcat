import axios from 'axios'
import type { ApiResponse } from '@/types'

const http = axios.create({
  baseURL: '/api',
  timeout: 30000,
  withCredentials: true,
})

// Request interceptor
http.interceptors.request.use(
  (config) => {
    config.withCredentials = true
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
http.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse
    if (res.code !== 0 && res.code !== 200) {
      return Promise.reject(new Error(res.message || 'Request failed'))
    }
    return response
  },
  (error) => {
    if (error.response) {
      const status = error.response.status
      if (status === 401) {
        // Redirect to login
        const currentPath = window.location.pathname
        if (currentPath !== '/login' && currentPath !== '/init') {
          window.location.href = '/login'
        }
      }
    }
    return Promise.reject(error)
  }
)

export default http

// Helper to extract data from response
export function extractData<T>(response: any): T {
  return response.data.data as T
}