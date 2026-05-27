package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
)

// DeviceHandler handles device management endpoints
type DeviceHandler struct {
	deviceService *services.DeviceService
}

// NewDeviceHandler creates a new DeviceHandler
func NewDeviceHandler(deviceService *services.DeviceService) *DeviceHandler {
	return &DeviceHandler{deviceService: deviceService}
}

// List handles GET /api/devices
func (h *DeviceHandler) List(c *gin.Context) {
	var groupID *uint
	if gid := c.Query("groupId"); gid != "" {
		if id, err := strconv.ParseUint(gid, 10, 64); err == nil {
			uid := uint(id)
			groupID = &uid
		}
	}

	var enabled *bool
	if e := c.Query("enabled"); e != "" {
		b := e == "true" || e == "1"
		enabled = &b
	}

	keyword := c.Query("keyword")

	devices, err := h.deviceService.ListDevices(groupID, enabled, keyword)
	if err != nil {
		response.InternalError(c, "failed to list devices")
		return
	}

	response.Success(c, devices)
}

// Create handles POST /api/devices
func (h *DeviceHandler) Create(c *gin.Context) {
	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.CreateDevice(&device); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, device)
}

// Update handles PUT /api/devices/:id
func (h *DeviceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid device id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.UpdateDevice(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "device updated", nil)
}

// Delete handles DELETE /api/devices/:id
func (h *DeviceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid device id")
		return
	}

	if err := h.deviceService.DeleteDevice(uint(id)); err != nil {
		response.InternalError(c, "failed to delete device")
		return
	}

	response.SuccessWithMessage(c, "device deleted", nil)
}

// --- DeviceGroup handlers ---

// ListGroups handles GET /api/device-groups
func (h *DeviceHandler) ListGroups(c *gin.Context) {
	groups, err := h.deviceService.ListDeviceGroups()
	if err != nil {
		response.InternalError(c, "failed to list device groups")
		return
	}

	response.Success(c, groups)
}

// CreateGroup handles POST /api/device-groups
func (h *DeviceHandler) CreateGroup(c *gin.Context) {
	var group models.DeviceGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.CreateDeviceGroup(&group); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, group)
}

// UpdateGroup handles PUT /api/device-groups/:id
func (h *DeviceHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.UpdateDeviceGroup(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "device group updated", nil)
}

// DeleteGroup handles DELETE /api/device-groups/:id
func (h *DeviceHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	if err := h.deviceService.DeleteDeviceGroup(uint(id)); err != nil {
		response.InternalError(c, "failed to delete device group")
		return
	}

	response.SuccessWithMessage(c, "device group deleted", nil)
}

// RegisterDeviceRoutes registers device routes
func RegisterDeviceRoutes(router *gin.RouterGroup, deviceService *services.DeviceService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewDeviceHandler(deviceService)

	devices := router.Group("/devices")
	devices.Use(AuthRequired())
	{
		devices.GET("", requirePerm("devices:list"), handler.List)
		devices.POST("", requirePerm("devices:create"), handler.Create)
		devices.PUT("/:id", requirePerm("devices:update"), handler.Update)
		devices.DELETE("/:id", requirePerm("devices:delete"), handler.Delete)
	}

	groups := router.Group("/device-groups")
	groups.Use(AuthRequired())
	{
		groups.GET("", requirePerm("device-groups:list"), handler.ListGroups)
		groups.POST("", requirePerm("device-groups:create"), handler.CreateGroup)
		groups.PUT("/:id", requirePerm("device-groups:update"), handler.UpdateGroup)
		groups.DELETE("/:id", requirePerm("device-groups:delete"), handler.DeleteGroup)
	}
}