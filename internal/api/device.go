package api

import (
	"net/http"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	applogger "syslog-alert/pkg/logger"
)

// ---- 设备 ----

// ListDevices 列出全部设备。
func (ws *WebServer) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices := repository.GetDevices()
	JSONResponse(w, devices)
}

// CreateDevice 创建设备。
func (ws *WebServer) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var device models.Device
	if !DecodeJSON(w, r, &device) {
		return
	}
	if err := repository.CreateDevice(&device); err != nil {
		applogger.Error("创建设备失败: %v", err)
		JSONError(w, "创建设备失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, device)
}

// GetDevice 获取单个设备。
func (ws *WebServer) GetDevice(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	device, err := repository.GetDeviceByID(id)
	if err != nil {
		applogger.Error("获取设备失败: %v", err)
		JSONError(w, "设备不存在", http.StatusNotFound)
		return
	}
	JSONResponse(w, device)
}

// UpdateDevice 更新设备。
func (ws *WebServer) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	var device models.Device
	if !DecodeJSON(w, r, &device) {
		return
	}
	device.ID = id
	if err := repository.UpdateDevice(&device); err != nil {
		applogger.Error("更新设备失败: %v", err)
		JSONError(w, "更新设备失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, device)
}

// DeleteDevice 删除设备。
func (ws *WebServer) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	if err := repository.DeleteDevice(id); err != nil {
		applogger.Error("删除设备失败: %v", err)
		JSONError(w, "删除设备失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}

// ---- 设备分组 ----

// ListDeviceGroups 列出全部设备分组。
func (ws *WebServer) ListDeviceGroups(w http.ResponseWriter, r *http.Request) {
	groups := repository.GetDeviceGroups()
	JSONResponse(w, groups)
}

// CreateDeviceGroup 创建设备分组。
func (ws *WebServer) CreateDeviceGroup(w http.ResponseWriter, r *http.Request) {
	var group models.DeviceGroup
	if !DecodeJSON(w, r, &group) {
		return
	}
	if err := repository.CreateDeviceGroup(&group); err != nil {
		applogger.Error("创建设备分组失败: %v", err)
		JSONError(w, "创建设备分组失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, group)
}

// GetDeviceGroup 获取单个设备分组。
func (ws *WebServer) GetDeviceGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	group, err := repository.GetDeviceGroupByID(id)
	if err != nil {
		applogger.Error("获取设备分组失败: %v", err)
		JSONError(w, "设备分组不存在", http.StatusNotFound)
		return
	}
	JSONResponse(w, group)
}

// UpdateDeviceGroup 更新设备分组。
func (ws *WebServer) UpdateDeviceGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	var group models.DeviceGroup
	if !DecodeJSON(w, r, &group) {
		return
	}
	group.ID = id
	if err := repository.UpdateDeviceGroup(&group); err != nil {
		applogger.Error("更新设备分组失败: %v", err)
		JSONError(w, "更新设备分组失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, group)
}

// DeleteDeviceGroup 删除设备分组。
func (ws *WebServer) DeleteDeviceGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	if err := repository.DeleteDeviceGroup(id); err != nil {
		applogger.Error("删除设备分组失败: %v", err)
		JSONError(w, "删除设备分组失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}
