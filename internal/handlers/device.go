package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)

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

	devices, total, err := h.deviceService.ListDevices(page, pageSize, groupID, enabled, keyword)
	if err != nil {
		response.InternalError(c, "failed to list devices")
		return
	}

	response.SuccessWithPage(c, devices, total, page, pageSize)
}

// Create handles POST /api/devices
func (h *DeviceHandler) Create(c *gin.Context) {
	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		writeAuditLog(c, "create", "device", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.CreateDevice(&device); err != nil {
		writeAuditLog(c, "create", "device", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "create", "device", strconv.FormatUint(uint64(device.ID), 10), "success", "device created")
	response.Created(c, device)
}

// Update handles PUT /api/devices/:id
func (h *DeviceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "device", c.Param("id"), "failure", "invalid device id")
		response.BadRequest(c, "invalid device id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeAuditLog(c, "update", "device", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.UpdateDevice(uint(id), updates); err != nil {
		writeAuditLog(c, "update", "device", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "update", "device", c.Param("id"), "success", "device updated")
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
		writeAuditLog(c, "delete", "device", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to delete device")
		return
	}

	writeAuditLog(c, "delete", "device", c.Param("id"), "success", "device deleted")
	response.SuccessWithMessage(c, "device deleted", nil)
}

// --- DeviceGroup handlers ---

// ListGroups handles GET /api/device-groups
func (h *DeviceHandler) ListGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	groups, total, err := h.deviceService.ListDeviceGroups(page, pageSize, keyword)
	if err != nil {
		response.InternalError(c, "failed to list device groups")
		return
	}

	response.SuccessWithPage(c, groups, total, page, pageSize)
}

// CreateGroup handles POST /api/device-groups
func (h *DeviceHandler) CreateGroup(c *gin.Context) {
	var group models.DeviceGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		writeAuditLog(c, "create", "device-group", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.CreateDeviceGroup(&group); err != nil {
		writeAuditLog(c, "create", "device-group", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "create", "device-group", strconv.FormatUint(uint64(group.ID), 10), "success", "device group created")
	response.Created(c, group)
}

// UpdateGroup handles PUT /api/device-groups/:id
func (h *DeviceHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "device-group", c.Param("id"), "failure", "invalid group id")
		response.BadRequest(c, "invalid group id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeAuditLog(c, "update", "device-group", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.UpdateDeviceGroup(uint(id), updates); err != nil {
		writeAuditLog(c, "update", "device-group", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "update", "device-group", c.Param("id"), "success", "device group updated")
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
		writeAuditLog(c, "delete", "device-group", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to delete device group")
		return
	}

	writeAuditLog(c, "delete", "device-group", c.Param("id"), "success", "device group deleted")
	response.SuccessWithMessage(c, "device group deleted", nil)
}

func (h *DeviceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	var item models.Device
	if err := db.First(&item, id).Error; err != nil {
		response.NotFound(c, "record not found")
		return
	}
	response.Success(c, item)
}

func (h *DeviceHandler) GetAll(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	var items []models.Device
	if err := db.Find(&items).Error; err != nil {
		response.InternalError(c, "failed to fetch records")
		return
	}
	response.Success(c, items)
}

func (h *DeviceHandler) GetGroupByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	var item models.DeviceGroup
	if err := db.First(&item, id).Error; err != nil {
		response.NotFound(c, "record not found")
		return
	}
	response.Success(c, item)
}

func (h *DeviceHandler) GetAllGroups(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	var items []models.DeviceGroup
	if err := db.Find(&items).Error; err != nil {
		response.InternalError(c, "failed to fetch records")
		return
	}
	response.Success(c, items)
}

// RegisterDeviceRoutes registers device routes
func RegisterDeviceRoutes(router *gin.RouterGroup, deviceService *services.DeviceService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewDeviceHandler(deviceService)

	devices := router.Group("/devices")
	devices.Use(AuthRequired())
	{
		devices.GET("/all", requirePerm("devices:list"), handler.GetAll)
		devices.GET("/:id", requirePerm("devices:list"), handler.GetByID)
		devices.GET("", requirePerm("devices:list"), handler.List)
		devices.POST("", requirePerm("devices:create"), handler.Create)
		devices.PUT("/:id", requirePerm("devices:update"), handler.Update)
		devices.DELETE("/:id", requirePerm("devices:delete"), handler.Delete)
	}

	groups := router.Group("/device-groups")
	groups.Use(AuthRequired())
	{
		groups.GET("/all", requirePerm("device-groups:list"), handler.GetAllGroups)
		groups.GET("/:id", requirePerm("device-groups:list"), handler.GetGroupByID)
		groups.GET("", requirePerm("device-groups:list"), handler.ListGroups)
		groups.POST("", requirePerm("device-groups:create"), handler.CreateGroup)
		groups.PUT("/:id", requirePerm("device-groups:update"), handler.UpdateGroup)
		groups.DELETE("/:id", requirePerm("device-groups:delete"), handler.DeleteGroup)
	}
}
