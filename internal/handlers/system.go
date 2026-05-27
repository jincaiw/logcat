package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
)

// SystemHandler handles system configuration and status endpoints
type SystemHandler struct {
	statsService *services.StatsService
}

// NewSystemHandler creates a new SystemHandler
func NewSystemHandler(statsService *services.StatsService) *SystemHandler {
	return &SystemHandler{statsService: statsService}
}

// GetConfig handles GET /api/system/config
func (h *SystemHandler) GetConfig(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var configs []models.SystemConfig
	if err := db.Find(&configs).Error; err != nil {
		response.InternalError(c, "failed to get system config")
		return
	}

	response.Success(c, configs)
}

// UpdateConfigRequest is the config update request body
type UpdateSystemConfigRequest struct {
	Configs []struct {
		Key   string `json:"configKey"`
		Value string `json:"configValue"`
	} `json:"configs"`
}

// UpdateConfig handles PUT /api/system/config
func (h *SystemHandler) UpdateConfig(c *gin.Context) {
	var req UpdateSystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	for _, cfg := range req.Configs {
		var existing models.SystemConfig
		result := db.Where("config_key = ?", cfg.Key).First(&existing)
		if result.Error != nil {
			db.Create(&models.SystemConfig{
				ConfigKey:   cfg.Key,
				ConfigValue: cfg.Value,
			})
		} else {
			existing.ConfigValue = cfg.Value
			db.Save(&existing)
		}
	}

	response.SuccessWithMessage(c, "system config updated", nil)
}

// Status handles GET /api/system/status
func (h *SystemHandler) Status(c *gin.Context) {
	status := h.statsService.GetSystemStatus()
	response.Success(c, status)
}

// SyslogStart handles POST /api/system/syslog/start
func (h *SystemHandler) SyslogStart(c *gin.Context) {
	// TODO: Implement syslog receiver start
	response.SuccessWithMessage(c, "syslog receiver started", nil)
}

// SyslogStop handles POST /api/system/syslog/stop
func (h *SystemHandler) SyslogStop(c *gin.Context) {
	// TODO: Implement syslog receiver stop
	response.SuccessWithMessage(c, "syslog receiver stopped", nil)
}

// RegisterSystemRoutes registers system routes
func RegisterSystemRoutes(router *gin.RouterGroup, statsService *services.StatsService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewSystemHandler(statsService)
	sys := router.Group("/system")
	sys.Use(AuthRequired())
	{
		sys.GET("/config", requirePerm("system:config:read"), handler.GetConfig)
		sys.PUT("/config", requirePerm("system:config:update"), handler.UpdateConfig)
		sys.GET("/status", requirePerm("system:status"), handler.Status)
		sys.POST("/syslog/start", requirePerm("system:syslog"), handler.SyslogStart)
		sys.POST("/syslog/stop", requirePerm("system:syslog"), handler.SyslogStop)
	}
}

// AuditLogHandler handles audit log endpoints
type AuditLogHandler struct {
	auditService *services.AuditService
}

// NewAuditLogHandler creates a new AuditLogHandler
func NewAuditLogHandler(auditService *services.AuditService) *AuditLogHandler {
	return &AuditLogHandler{auditService: auditService}
}

// List handles GET /api/audit-logs
func (h *AuditLogHandler) List(c *gin.Context) {
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		page, _ = parseInt(p)
	}
	if ps := c.Query("pageSize"); ps != "" {
		pageSize, _ = parseInt(ps)
	}

	logs, total, err := h.auditService.List(page, pageSize,
		c.Query("username"), c.Query("action"), c.Query("result"), c.Query("resourceType"))
	if err != nil {
		response.InternalError(c, "failed to list audit logs")
		return
	}

	response.SuccessWithPage(c, logs, total, page, pageSize)
}

// RegisterAuditLogRoutes registers audit log routes
func RegisterAuditLogRoutes(router *gin.RouterGroup, auditService *services.AuditService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewAuditLogHandler(auditService)
	al := router.Group("/audit-logs")
	al.Use(AuthRequired())
	{
		al.GET("", requirePerm("audit-logs:list"), handler.List)
	}
}

// parseInt is a helper to parse integer from string
func parseInt(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}