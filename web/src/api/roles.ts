import http, { extractData } from './index'
import type { ApiResponse, Role, PageResponse, PageParams, Permission } from '@/types'

export function getRoles(params: PageParams & { name?: string }): Promise<ApiResponse<PageResponse<Role>>> {
  return http.get('/roles', { params }).then(extractData)
}

export function getAllRoles(): Promise<ApiResponse<Role[]>> {
  return http.get('/roles/all').then(extractData)
}

export function getRole(id: number): Promise<ApiResponse<Role>> {
  return http.get(`/roles/${id}`).then(extractData)
}

export function createRole(data: Partial<Role>): Promise<ApiResponse<Role>> {
  return http.post('/roles', data).then(extractData)
}

export function updateRole(id: number, data: Partial<Role>): Promise<ApiResponse<Role>> {
  return http.put(`/roles/${id}`, data).then(extractData)
}

export function deleteRole(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/roles/${id}`).then(extractData)
}

export function getRolePermissions(roleId: number): Promise<ApiResponse<Permission[]>> {
  return http.get(`/roles/${roleId}/permissions`).then(extractData)
}

export function assignPermissions(roleId: number, permissionIds: number[]): Promise<ApiResponse<null>> {
  return http.post(`/roles/${roleId}/permissions`, { permissionIds }).then(extractData)
}

export function getAllPermissions(): Promise<ApiResponse<Permission[]>> {
  return http.get('/permissions').then(extractData)
}