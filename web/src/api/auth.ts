import http, { extractData } from './index'
import type { ApiResponse, User } from '@/types'

export function login(username: string, password: string): Promise<ApiResponse<{ user: User }>> {
  return http.post('/auth/login', { username, password }).then(extractData)
}

export function logout(): Promise<ApiResponse<null>> {
  return http.post('/auth/logout').then(extractData)
}

export function getCurrentUser(): Promise<ApiResponse<User>> {
  return http.get('/auth/me').then(extractData)
}

export function changePassword(oldPwd: string, newPwd: string): Promise<ApiResponse<null>> {
  return http.post('/auth/change-password', { oldPassword: oldPwd, newPassword: newPwd }).then(extractData)
}

export function refreshToken(): Promise<ApiResponse<{ token: string }>> {
  return http.post('/auth/refresh').then(extractData)
}

export function initAdmin(password: string): Promise<ApiResponse<null>> {
  return http.post('/auth/init-admin', { username: 'admin', password }).then(extractData)
}

export function checkInitStatus(): Promise<ApiResponse<{ initialized: boolean }>> {
  return http.get('/auth/init-status').then(extractData)
}
