import http, { extractData } from './index'
import type { ApiResponse, Device, DeviceGroup, PageResponse, PageParams } from '@/types'

// ============ Devices ============

export function getDevices(params: PageParams & { name?: string; ipAddress?: string; enabled?: boolean; groupId?: number }): Promise<ApiResponse<PageResponse<Device>>> {
  return http.get('/devices', { params }).then(extractData)
}

export function getAllDevices(): Promise<ApiResponse<Device[]>> {
  return http.get('/devices/all').then(extractData)
}

export function getDevice(id: number): Promise<ApiResponse<Device>> {
  return http.get(`/devices/${id}`).then(extractData)
}

export function createDevice(data: Partial<Device>): Promise<ApiResponse<Device>> {
  return http.post('/devices', data).then(extractData)
}

export function updateDevice(id: number, data: Partial<Device>): Promise<ApiResponse<Device>> {
  return http.put(`/devices/${id}`, data).then(extractData)
}

export function deleteDevice(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/devices/${id}`).then(extractData)
}

// ============ Device Groups ============

export function getDeviceGroups(params: PageParams & { name?: string }): Promise<ApiResponse<PageResponse<DeviceGroup>>> {
  return http.get('/device-groups', { params }).then(extractData)
}

export function getAllDeviceGroups(): Promise<ApiResponse<DeviceGroup[]>> {
  return http.get('/device-groups/all').then(extractData)
}

export function getDeviceGroup(id: number): Promise<ApiResponse<DeviceGroup>> {
  return http.get(`/device-groups/${id}`).then(extractData)
}

export function createDeviceGroup(data: Partial<DeviceGroup>): Promise<ApiResponse<DeviceGroup>> {
  return http.post('/device-groups', data).then(extractData)
}

export function updateDeviceGroup(id: number, data: Partial<DeviceGroup>): Promise<ApiResponse<DeviceGroup>> {
  return http.put(`/device-groups/${id}`, data).then(extractData)
}

export function deleteDeviceGroup(id: number): Promise<ApiResponse<null>> {
  return http.delete(`/device-groups/${id}`).then(extractData)
}