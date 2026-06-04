import http, { extractData } from './index'
import type { ApiResponse, User, PageResponse, PageParams } from '@/types'

export function getUsers(params: PageParams & { username?: string; status?: number }): Promise<ApiResponse<PageResponse<User>>> {
  return http.get('/users', { params }).then(extractData)
}

export function getUser(id: number): Promise<ApiResponse<User>> {
  return http.get(`/users/${id}`).then(extractData)
}

export function createUser(data: Partial<User> & { password?: string }): Promise<ApiResponse<User>> {
  return http.post('/users', data).then(extractData)
}

export function updateUser(id: number, data: Partial<User>): Promise<ApiResponse<User>> {
  return http.put(`/users/${id}`, data).then(extractData)
}

export function deleteUser(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/users/${id}`).then(extractData)
}

export function resetPassword(userId: number, password: string): Promise<ApiResponse<null>> {
  return http.post(`/users/${userId}/reset-password`, { password }).then(extractData)
}

export function unlockUser(userId: number): Promise<ApiResponse<null>> {
  return http.post(`/users/${userId}/unlock`).then(extractData)
}

export function forcePasswordChange(userId: number): Promise<ApiResponse<null>> {
  return http.post(`/users/${userId}/force-password-change`).then(extractData)
}

export function assignRoles(userId: number, roleIds: number[]): Promise<ApiResponse<null>> {
  return http.post(`/users/${userId}/roles`, { roleIds }).then(extractData)
}

export function getUserRoles(userId: number): Promise<ApiResponse<{ roles: import('@/types').Role[] }>> {
  return http.get(`/users/${userId}/roles`).then(extractData)
}

export function getCurrentUser(): Promise<ApiResponse<User & { permissions?: string[] }>> {
  return http.get('/auth/me').then(extractData)
}
