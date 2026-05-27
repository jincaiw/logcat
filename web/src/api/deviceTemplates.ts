import http, { extractData } from './index'
import type { ApiResponse, DeviceTemplate, PageResponse, PageParams } from '@/types'

export function getDeviceTemplates(params: PageParams & { name?: string }): Promise<ApiResponse<PageResponse<DeviceTemplate>>> {
  return http.get('/device-templates', { params }).then(extractData)
}

export function getAllDeviceTemplates(): Promise<ApiResponse<DeviceTemplate[]>> {
  return http.get('/device-templates/all').then(extractData)
}

export function getDeviceTemplate(id: number): Promise<ApiResponse<DeviceTemplate>> {
  return http.get(`/device-templates/${id}`).then(extractData)
}

export function createDeviceTemplate(data: Partial<DeviceTemplate>): Promise<ApiResponse<DeviceTemplate>> {
  return http.post('/device-templates', data).then(extractData)
}

export function updateDeviceTemplate(id: number, data: Partial<DeviceTemplate>): Promise<ApiResponse<DeviceTemplate>> {
  return http.put(`/device-templates/${id}`, data).then(extractData)
}

export function deleteDeviceTemplate(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/device-templates/${id}`).then(extractData)
}