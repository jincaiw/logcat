import axios from 'axios'
import type { AxiosResponse } from 'axios'
import type { ApiResponse } from '@/types'

const http = axios.create({
  baseURL: '/api',
  timeout: 30000,
  withCredentials: true,
})

let onUnauthorized: (() => void) | null = null

export function setOnUnauthorized(fn: () => void) {
  onUnauthorized = fn
}

let isRefreshing = false
let pendingRequests: Array<(token: string | null) => void> = []

function onRefreshed(token: string | null) {
  pendingRequests.forEach(cb => cb(token))
  pendingRequests = []
}

http.interceptors.request.use(
  (config) => {
    config.withCredentials = true
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

http.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse
    if (res.code !== 0 && res.code !== 200 && res.code !== 201) {
      return Promise.reject(new Error(res.message || 'Request failed'))
    }
    return response
  },
  (error) => {
    const originalRequest = error.config
    if (!originalRequest) {
      return Promise.reject(error)
    }
    if (error.response) {
      const status = error.response.status
      if (status === 401 && !originalRequest._retry) {
        const requestUrl = String(originalRequest?.url || '')
        if (requestUrl.includes('/auth/refresh') || requestUrl.includes('/auth/login')) {
          handleUnauthorized()
          return Promise.reject(error)
        }
        originalRequest._retry = true
        if (!isRefreshing) {
          isRefreshing = true
          return http.post('/auth/refresh').then((refreshRes) => {
            const res = refreshRes.data as ApiResponse
            if (res.code === 0 || res.code === 200) {
              isRefreshing = false
              onRefreshed('refreshed')
              return http(originalRequest)
            } else {
              isRefreshing = false
              onRefreshed(null)
              handleUnauthorized()
              return Promise.reject(error)
            }
          }).catch(() => {
            isRefreshing = false
            onRefreshed(null)
            handleUnauthorized()
            return Promise.reject(error)
          })
        }
        return new Promise((resolve, reject) => {
          pendingRequests.push((token) => {
            if (token) {
              resolve(http(originalRequest))
            } else {
              reject(error)
            }
          })
        })
      }
    }
    return Promise.reject(error)
  }
)

function handleUnauthorized() {
  // Clear pending requests to prevent memory leak
  pendingRequests = []
  const currentPath = window.location.pathname
  if (currentPath !== '/login') {
    if (onUnauthorized) {
      onUnauthorized()
    } else {
      window.location.href = '/login'
    }
  }
}

export default http

export function extractData<T = any>(response: AxiosResponse<ApiResponse<T>>): ApiResponse<T> {
  return response.data
}
