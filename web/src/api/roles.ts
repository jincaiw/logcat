import http, { extractData } from './index'
import type { ApiResponse, Role, PageResponse, PageParams, Permission } from '@/types'

export function getRoles(params: PageParams & { name?: string }): Promise<ApiResponse<PageResponse<Role>>> {
  return http.get('/roles', { params }).then(extractData)
}

export function getAllRoles(): Promise<ApiResponse<Role[]>> {
  return http.get('/roles', { params: { pageSize: 1000 } }).then(extractData).then((res) => {
    const raw = res.data
    if (Array.isArray(raw)) return res
    if (raw?.list && Array.isArray(raw.list)) {
      return { ...res, data: raw.list }
    }
    if (raw?.items && Array.isArray(raw.items)) {
      return { ...res, data: raw.items }
    }
    return { ...res, data: [] }
  })
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
  return http.get('/permissions').then(extractData).then((res) => {
    const list = Array.isArray(res.data) ? res.data : []
    return {
      ...res,
      data: list.map((item: Permission) => ({
        ...item,
        group: item.resource || item.type || 'default',
        description: item.description || item.action || item.code,
      })),
    }
  })
}
