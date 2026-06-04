package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	appconfig "github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/middleware"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	logsyslog "github.com/logcat/logcat/internal/syslog"
	"github.com/logcat/logcat/pkg/response"
	"gorm.io/gorm"
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

	merged := mergeSystemConfigs(configs)
	response.Success(c, merged)
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

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, cfg := range req.Configs {
			if cfg.Key == "" {
				continue
			}
			var existing models.SystemConfig
			result := tx.Where("config_key = ?", cfg.Key).First(&existing)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					if err := tx.Create(&models.SystemConfig{
						ConfigKey:   cfg.Key,
						ConfigValue: cfg.Value,
						Description: "",
					}).Error; err != nil {
						return err
					}
				} else {
					return result.Error
				}
			} else {
				existing.ConfigValue = cfg.Value
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		writeAuditLog(c, "update_config", "system", "", "failure", err.Error())
		response.InternalError(c, "failed to update system config")
		return
	}

	updatedKeys := make([]string, 0, len(req.Configs))
	for _, cfg := range req.Configs {
		if cfg.Key != "" {
			updatedKeys = append(updatedKeys, cfg.Key)
		}
	}

	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	var uid *uint
	if userID > 0 {
		uid = &userID
	}
	middleware.AuditLogWriter(
		uid,
		username,
		"update_config",
		"system",
		"",
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		"success",
		"system config updated: "+strings.Join(updatedKeys, ", "),
	)

	response.SuccessWithMessage(c, "system config updated", nil)
}

func mergeSystemConfigs(stored []models.SystemConfig) []models.SystemConfig {
	defaults := defaultSystemConfigs()
	index := make(map[string]int, len(defaults))
	for i, cfg := range defaults {
		index[cfg.ConfigKey] = i
	}

	for _, cfg := range stored {
		if i, ok := index[cfg.ConfigKey]; ok {
			defaults[i].ConfigValue = cfg.ConfigValue
			if cfg.Description != "" {
				defaults[i].Description = cfg.Description
			}
			continue
		}
		defaults = append(defaults, cfg)
		index[cfg.ConfigKey] = len(defaults) - 1
	}

	return defaults
}

func defaultSystemConfigs() []models.SystemConfig {
	cfg := appconfig.Get()
	if cfg == nil {
		return nil
	}

	return []models.SystemConfig{
		{ConfigKey: "serverHost", ConfigValue: cfg.Server.Host, Description: "Web/API 监听地址，保存后需重启服务"},
		{ConfigKey: "serverPort", ConfigValue: intToString(cfg.Server.Port), Description: "Web/API 监听端口，默认 5080，保存后需重启服务"},
		{ConfigKey: "syslogEnabled", ConfigValue: boolToString(cfg.Syslog.Enabled), Description: "是否启用 Syslog 接收器，保存后需重启服务"},
		{ConfigKey: "syslogUdpPort", ConfigValue: intToString(cfg.Syslog.UDPPort), Description: "Syslog UDP 监听端口，保存后需重启服务"},
		{ConfigKey: "syslogTcpPort", ConfigValue: intToString(cfg.Syslog.TCPPort), Description: "Syslog TCP 监听端口，保存后需重启服务"},
		{ConfigKey: "databaseType", ConfigValue: cfg.Database.Type, Description: "数据库类型，sqlite 或 mysql，保存后需重启服务"},
		{ConfigKey: "sqlitePath", ConfigValue: cfg.Database.SQLite.Path, Description: "SQLite 数据文件路径，保存后需重启服务"},
		{ConfigKey: "sessionExpireHours", ConfigValue: intToString(cfg.Auth.SessionExpireHours), Description: "会话过期时间（小时），保存后需重启服务"},
		{ConfigKey: "maxFailedLogin", ConfigValue: intToString(cfg.Auth.MaxFailedLogin), Description: "连续登录失败阈值，保存后需重启服务"},
		{ConfigKey: "lockDurationMinutes", ConfigValue: intToString(cfg.Auth.LockDurationMinutes), Description: "账号锁定时长（分钟），保存后需重启服务"},
		{ConfigKey: "retentionDays", ConfigValue: intToString(cfg.Log.RetentionDays), Description: "日志保留天数，保存后需重启服务"},
		{ConfigKey: "unmatchedRetentionDays", ConfigValue: intToString(cfg.Log.UnmatchedRetentionDays), Description: "未匹配日志保留天数，保存后需重启服务"},
		{ConfigKey: "maxLogSize", ConfigValue: intToString(cfg.Log.MaxLogSize), Description: "最大日志大小（MB），保存后需重启服务"},
		{ConfigKey: "parseWorkers", ConfigValue: intToString(cfg.Worker.ParseWorkers), Description: "解析 worker 数量，保存后需重启服务"},
		{ConfigKey: "filterWorkers", ConfigValue: intToString(cfg.Worker.FilterWorkers), Description: "筛选 worker 数量，保存后需重启服务"},
		{ConfigKey: "pushWorkers", ConfigValue: intToString(cfg.Worker.PushWorkers), Description: "推送 worker 数量，保存后需重启服务"},
		{ConfigKey: "queueCapacity", ConfigValue: intToString(cfg.Queue.Capacity), Description: "队列容量，保存后需重启服务"},
		{ConfigKey: "queueFullPolicy", ConfigValue: cfg.Queue.FullPolicy, Description: "队列满时策略，保存后需重启服务"},
		{ConfigKey: "timezone", ConfigValue: "Asia/Shanghai", Description: "系统显示时区，即时生效"},
		{ConfigKey: "timeFormat", ConfigValue: "YYYY-MM-DD HH:mm:ss", Description: "时间显示格式，即时生效"},
	}
}

func intToString(v int) string {
	return strconv.Itoa(v)
}

func boolToString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// Status handles GET /api/system/status
func (h *SystemHandler) Status(c *gin.Context) {
	status := h.statsService.GetSystemStatus()
	response.Success(c, status)
}

// SyslogStart handles POST /api/system/syslog/start
func (h *SystemHandler) SyslogStart(c *gin.Context) {
	receiver := logsyslog.GetGlobalReceiver()
	if receiver == nil {
		response.InternalError(c, "syslog receiver not initialized")
		return
	}
	if err := receiver.Start(); err != nil {
		userID := middleware.GetUserID(c)
		username := middleware.GetUsername(c)
		var uid *uint
		if userID > 0 {
			uid = &userID
		}
		middleware.AuditLogWriter(uid, username, "syslog_start", "system", "", c.ClientIP(), c.GetHeader("User-Agent"), "failure", err.Error())
		response.InternalError(c, "failed to start syslog receiver: "+err.Error())
		return
	}
	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	var uid *uint
	if userID > 0 {
		uid = &userID
	}
	middleware.AuditLogWriter(uid, username, "syslog_start", "system", "", c.ClientIP(), c.GetHeader("User-Agent"), "success", "syslog receiver started")
	response.Success(c, gin.H{
		"running": receiver.IsRunning(),
		"metrics": receiver.Metrics(),
	})
}

// SyslogStop handles POST /api/system/syslog/stop
func (h *SystemHandler) SyslogStop(c *gin.Context) {
	receiver := logsyslog.GetGlobalReceiver()
	if receiver == nil {
		response.InternalError(c, "syslog receiver not initialized")
		return
	}
	receiver.Stop()
	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	var uid *uint
	if userID > 0 {
		uid = &userID
	}
	middleware.AuditLogWriter(uid, username, "syslog_stop", "system", "", c.ClientIP(), c.GetHeader("User-Agent"), "success", "syslog receiver stopped")
	response.Success(c, gin.H{
		"running": receiver.IsRunning(),
		"metrics": receiver.Metrics(),
	})
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
	page, pageSize := 1, 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			pageSize = v
		}
	}
	page, pageSize = response.NormalizePagination(page, pageSize)

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

