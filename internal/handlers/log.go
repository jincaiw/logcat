package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/middleware"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
)

// LogHandler handles syslog log endpoints
type LogHandler struct {
	traceService   *services.TraceService
	cleanupService *services.CleanupService
}

// NewLogHandler creates a new LogHandler
func NewLogHandler(traceService *services.TraceService, cleanupService *services.CleanupService) *LogHandler {
	return &LogHandler{traceService: traceService, cleanupService: cleanupService}
}

// List handles GET /api/logs
func (h *LogHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	sourceIP := c.Query("sourceIp")
	eventType := c.Query("eventType")
	severity := c.Query("severity")
	filterStatus := c.Query("filterStatus")
	keyword := c.Query("keyword")
	deviceIDStr := c.Query("deviceId")
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")
	destinationIP := c.Query("destinationIp")
	alertStatus := c.Query("alertStatus")
	logID := c.Query("logId")

	var logs []models.SyslogLog
	var total int64

	query := db.Model(&models.SyslogLog{})
	if sourceIP != "" {
		query = query.Where("source_ip = ?", sourceIP)
	}
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if filterStatus != "" {
		query = query.Where("filter_status = ?", filterStatus)
	}
	if keyword != "" {
		query = query.Where("raw_message LIKE ?", "%"+response.EscapeLike(keyword)+"%")
	}
	if deviceIDStr != "" {
		if deviceID, err := strconv.ParseUint(deviceIDStr, 10, 64); err == nil {
			query = query.Where("device_id = ?", deviceID)
		}
	}
	if startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			query = query.Where("received_at >= ?", t)
		}
	}
	if endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			query = query.Where("received_at <= ?", t)
		}
	}
	if destinationIP != "" {
		query = query.Where("destination_ip = ?", destinationIP)
	}
	if alertStatus != "" {
		query = query.Where("alert_status = ?", alertStatus)
	}
	if logID != "" {
		query = query.Where("log_id = ?", logID)
	}

	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count logs")
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("received_at DESC").Find(&logs).Error; err != nil {
		response.InternalError(c, "failed to list logs")
		return
	}

	response.SuccessWithPage(c, logs, total, page, pageSize)
}

// GetByID handles GET /api/logs/:id
func (h *LogHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid log id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var log models.SyslogLog
	if err := db.First(&log, id).Error; err != nil {
		response.NotFound(c, "log not found")
		return
	}

	response.Success(c, log)
}

// Export handles GET /api/logs/export
func (h *LogHandler) Export(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var logs []models.SyslogLog
	query := db.Model(&models.SyslogLog{})

	if sourceIP := c.Query("sourceIp"); sourceIP != "" {
		query = query.Where("source_ip = ?", sourceIP)
	}
	if eventType := c.Query("eventType"); eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if severity := c.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if filterStatus := c.Query("filterStatus"); filterStatus != "" {
		query = query.Where("filter_status = ?", filterStatus)
	}
	if startTimeStr := c.Query("startTime"); startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			query = query.Where("received_at >= ?", t)
		}
	}
	if endTimeStr := c.Query("endTime"); endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			query = query.Where("received_at <= ?", t)
		}
	}

	if err := query.Order("received_at DESC").Limit(10000).Find(&logs).Error; err != nil {
		response.InternalError(c, "failed to export logs")
		return
	}

	response.Success(c, logs)
}

// UnmatchedCount handles GET /api/logs/unmatched-count
func (h *LogHandler) UnmatchedCount(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var count int64
	if err := db.Model(&models.SyslogLog{}).Where("filter_status = ?", "no_device").Count(&count).Error; err != nil {
		response.InternalError(c, "failed to count unmatched logs")
		return
	}

	response.Success(c, gin.H{"count": count})
}

// Cleanup handles DELETE /api/logs/cleanup
func (h *LogHandler) Cleanup(c *gin.Context) {
	if err := h.cleanupService.RunCleanupNow(); err != nil {
		writeAuditLog(c, "cleanup", "logs", "", "failure", err.Error())
		response.InternalError(c, "cleanup failed: "+err.Error())
		return
	}
	writeAuditLog(c, "cleanup", "logs", "", "success", "log cleanup completed")
	response.SuccessWithMessage(c, "cleanup completed", nil)
}

// Trace handles GET /api/logs/:id/trace
func (h *LogHandler) Trace(c *gin.Context) {
	logID := c.Param("id")
	if logID == "" {
		response.BadRequest(c, "log id is required")
		return
	}

	trace, err := h.traceService.GetTraceByLogID(logID)
	if err != nil {
		response.NotFound(c, "trace not found")
		return
	}

	response.Success(c, trace)
}

// RegisterLogRoutes registers log routes
func RegisterLogRoutes(router *gin.RouterGroup, traceService *services.TraceService, cleanupService *services.CleanupService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewLogHandler(traceService, cleanupService)
	logs := router.Group("/logs")
	logs.Use(AuthRequired())
	{
		logs.GET("", requirePerm("logs:list"), handler.List)
		logs.GET("/export", requirePerm("logs:list"), handler.Export)
		logs.GET("/unmatched-count", requirePerm("logs:list"), handler.UnmatchedCount)
		logs.GET("/:id", requirePerm("logs:list"), handler.GetByID)
		logs.DELETE("/cleanup", requirePerm("logs:cleanup"), handler.Cleanup)
		logs.GET("/:id/trace", requirePerm("logs:trace"), handler.Trace)
	}
}

// AlertRecordHandler handles alert record endpoints
type AlertRecordHandler struct {
	alertService *services.AlertService
}

// NewAlertRecordHandler creates a new AlertRecordHandler
func NewAlertRecordHandler(alertService *services.AlertService) *AlertRecordHandler {
	return &AlertRecordHandler{alertService: alertService}
}

// List handles GET /api/alerts
func (h *AlertRecordHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	status := c.Query("status")
	channelType := c.Query("channelType")
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")
	logID := c.Query("logId")

	var records []models.AlertRecord
	var total int64

	query := db.Model(&models.AlertRecord{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if channelType != "" {
		query = query.Where("channel_type = ?", channelType)
	}
	if startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			query = query.Where("sent_at >= ?", t)
		}
	}
	if endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			query = query.Where("sent_at <= ?", t)
		}
	}
	if logID != "" {
		query = query.Where("log_id = ?", logID)
	}

	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count alert records")
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").
		Preload("AlertRule").Preload("PushConfig").Find(&records).Error; err != nil {
		response.InternalError(c, "failed to list alert records")
		return
	}

	response.SuccessWithPage(c, records, total, page, pageSize)
}

// GetByID handles GET /api/alerts/:id
func (h *AlertRecordHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid alert id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var record models.AlertRecord
	if err := db.Preload("AlertRule").Preload("PushConfig").First(&record, id).Error; err != nil {
		response.NotFound(c, "alert record not found")
		return
	}

	response.Success(c, record)
}

// CreateDispositionRequest is the request body for creating a disposition
type CreateDispositionRequest struct {
	Status string `json:"status" binding:"required"`
	Note   string `json:"note"`
}

// CreateDisposition handles POST /api/alerts/:id/dispositions
func (h *AlertRecordHandler) CreateDisposition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "create_disposition", "alert", c.Param("id"), "failure", "invalid alert id")
		response.BadRequest(c, "invalid alert id")
		return
	}

	var req CreateDispositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAuditLog(c, "create_disposition", "alert", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	alertID := uint(id)
	operatorID := middleware.GetUserID(c)
	operatorName := middleware.GetUsername(c)

	disposition, createErr := h.alertService.CreateDisposition(&alertID, nil, req.Status, req.Note, &operatorID, operatorName)
	if createErr != nil {
		writeAuditLog(c, "create_disposition", "alert", c.Param("id"), "failure", createErr.Error())
		response.InternalError(c, "failed to create disposition")
		return
	}

	writeAuditLog(c, "create_disposition", "alert", c.Param("id"), "success", "alert disposition created with status "+req.Status)
	response.Created(c, disposition)
}

// ListDispositions handles GET /api/alerts/:id/dispositions
func (h *AlertRecordHandler) ListDispositions(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	dispositions, err := h.alertService.ListDispositions(uint(id))
	if err != nil {
		response.InternalError(c, "failed to list dispositions")
		return
	}

	response.Success(c, dispositions)
}

// ListAllDispositions handles GET /api/alert-dispositions
func (h *AlertRecordHandler) ListAllDispositions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	status := c.Query("status")

	dispositions, total, err := h.alertService.ListAllDispositions(page, pageSize, status)
	if err != nil {
		response.InternalError(c, "failed to list alert dispositions")
		return
	}

	response.SuccessWithPage(c, dispositions, total, page, pageSize)
}

// RegisterAlertRecordRoutes registers alert record routes
func RegisterAlertRecordRoutes(router *gin.RouterGroup, alertService *services.AlertService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewAlertRecordHandler(alertService)
	alerts := router.Group("/alerts")
	alerts.Use(AuthRequired())
	{
		alerts.GET("", requirePerm("alerts:list"), handler.List)
		alerts.GET("/:id", requirePerm("alerts:list"), handler.GetByID)
		alerts.POST("/:id/dispositions", requirePerm("alerts:disposition:create"), handler.CreateDisposition)
		alerts.GET("/:id/dispositions", requirePerm("alerts:disposition:list"), handler.ListDispositions)
	}

	dispositions := router.Group("/alert-dispositions")
	dispositions.Use(AuthRequired())
	{
		dispositions.GET("", requirePerm("alerts:disposition:list"), handler.ListAllDispositions)
	}
}

// AggregatedAlertHandler handles aggregated alert endpoints
type AggregatedAlertHandler struct {
	aggregateService *services.AggregateService
}

// NewAggregatedAlertHandler creates a new AggregatedAlertHandler
func NewAggregatedAlertHandler(aggregateService *services.AggregateService) *AggregatedAlertHandler {
	return &AggregatedAlertHandler{aggregateService: aggregateService}
}

// List handles GET /api/aggregated-alerts
func (h *AggregatedAlertHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)

	alerts, total, err := h.aggregateService.ListAggregatedAlerts(page, pageSize,
		c.Query("status"), c.Query("severity"), c.Query("eventType"))
	if err != nil {
		response.InternalError(c, "failed to list aggregated alerts")
		return
	}

	response.SuccessWithPage(c, alerts, total, page, pageSize)
}

// GetLogs handles GET /api/aggregated-alerts/:id/logs
func (h *AggregatedAlertHandler) GetLogs(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var logs []models.SyslogLog
	if err := db.Where("aggregated_alert_id = ?", id).Order("received_at DESC").Limit(100).Find(&logs).Error; err != nil {
		response.InternalError(c, "failed to get logs")
		return
	}

	response.Success(c, logs)
}

// Get handles GET /api/aggregated-alerts/:id
func (h *AggregatedAlertHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	alert, err := h.aggregateService.GetAggregatedAlert(uint(id))
	if err != nil {
		response.NotFound(c, "aggregated alert not found")
		return
	}

	response.Success(c, alert)
}

// Acknowledge handles POST /api/aggregated-alerts/:id/acknowledge
func (h *AggregatedAlertHandler) Acknowledge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "acknowledge", "aggregated-alert", c.Param("id"), "failure", "invalid aggregated alert id")
		response.BadRequest(c, "invalid aggregated alert id")
		return
	}
	if err := h.aggregateService.UpdateStatus(uint(id), "acknowledged"); err != nil {
		writeAuditLog(c, "acknowledge", "aggregated-alert", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to acknowledge aggregated alert")
		return
	}
	writeAuditLog(c, "acknowledge", "aggregated-alert", c.Param("id"), "success", "aggregated alert acknowledged")
	response.SuccessWithMessage(c, "aggregated alert acknowledged", nil)
}

// Resolve handles POST /api/aggregated-alerts/:id/resolve
func (h *AggregatedAlertHandler) Resolve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "resolve", "aggregated-alert", c.Param("id"), "failure", "invalid aggregated alert id")
		response.BadRequest(c, "invalid aggregated alert id")
		return
	}
	if err := h.aggregateService.UpdateStatus(uint(id), "resolved"); err != nil {
		writeAuditLog(c, "resolve", "aggregated-alert", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to resolve aggregated alert")
		return
	}
	writeAuditLog(c, "resolve", "aggregated-alert", c.Param("id"), "success", "aggregated alert resolved")
	response.SuccessWithMessage(c, "aggregated alert resolved", nil)
}

// RegisterAggregatedAlertRoutes registers aggregated alert routes
func RegisterAggregatedAlertRoutes(router *gin.RouterGroup, aggregateService *services.AggregateService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewAggregatedAlertHandler(aggregateService)
	aa := router.Group("/aggregated-alerts")
	aa.Use(AuthRequired())
	{
		aa.GET("", requirePerm("aggregated-alerts:list"), handler.List)
		aa.GET("/:id", requirePerm("aggregated-alerts:list"), handler.Get)
		aa.GET("/:id/logs", requirePerm("aggregated-alerts:logs"), handler.GetLogs)
		aa.POST("/:id/acknowledge", requirePerm("aggregated-alerts:update"), handler.Acknowledge)
		aa.POST("/:id/resolve", requirePerm("aggregated-alerts:update"), handler.Resolve)
	}
}

// HighFreqIPHandler handles high frequency IP endpoints
type HighFreqIPHandler struct {
	highFreqService *services.HighFreqService
}

// NewHighFreqIPHandler creates a new HighFreqIPHandler
func NewHighFreqIPHandler(highFreqService *services.HighFreqService) *HighFreqIPHandler {
	return &HighFreqIPHandler{highFreqService: highFreqService}
}

// List handles GET /api/high-frequency-ips
func (h *HighFreqIPHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	rawIPs := h.highFreqService.GetHighFreqIPs()

	// Filter by keyword if provided
	filtered := rawIPs
	if keyword != "" {
		filtered = make([]services.HighFreqIPEntry, 0)
		for _, item := range rawIPs {
			if strings.Contains(item.IP, keyword) {
				filtered = append(filtered, item)
			}
		}
	}

	total := int64(len(filtered))
	offset := (page - 1) * pageSize
	if offset > len(filtered) {
		offset = len(filtered)
	}
	end := offset + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	paged := filtered[offset:end]
	items := make([]gin.H, 0, len(paged))
	for _, item := range paged {
		items = append(items, gin.H{
			"ip":          item.IP,
			"count":       item.Count,
			"firstSeen":   item.FirstSeen,
			"lastSeen":    item.LastSeen,
			"deviceNames": []string{},
		})
	}
	response.SuccessWithPage(c, items, total, page, pageSize)
}

// UpdateConfigRequest is the config update request body
type HighFreqConfigRequest struct {
	Threshold         int `json:"threshold"`
	WindowSeconds     int `json:"windowSeconds"`
	TimeWindowSeconds int `json:"timeWindowSeconds"`
}

// UpdateConfig handles PUT /api/high-frequency-ips/config
func (h *HighFreqIPHandler) UpdateConfig(c *gin.Context) {
	var req HighFreqConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAuditLog(c, "update_config", "high-frequency-ip", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	windowSeconds := req.WindowSeconds
	if req.TimeWindowSeconds > 0 {
		windowSeconds = req.TimeWindowSeconds
	}
	if windowSeconds <= 0 {
		windowSeconds = 60
	}
	if req.Threshold <= 0 {
		req.Threshold = 100
	}

	h.highFreqService.UpdateConfig(req.Threshold, windowSeconds)
	writeAuditLog(c, "update_config", "high-frequency-ip", "", "success", "high frequency IP config updated")
	response.Success(c, gin.H{
		"threshold":         req.Threshold,
		"timeWindowSeconds": windowSeconds,
	})
}

// GetConfig handles GET /api/high-frequency-ips/config
func (h *HighFreqIPHandler) GetConfig(c *gin.Context) {
	cfg := h.highFreqService.GetConfig()
	response.Success(c, gin.H{
		"threshold":         cfg["threshold"],
		"timeWindowSeconds": cfg["window_seconds"],
	})
}

// Refresh handles POST /api/high-frequency-ips/refresh
func (h *HighFreqIPHandler) Refresh(c *gin.Context) {
	response.SuccessWithMessage(c, "refreshed", gin.H{
		"count": len(h.highFreqService.GetHighFreqIPs()),
	})
}

// RegisterHighFreqIPRoutes registers high frequency IP routes
func RegisterHighFreqIPRoutes(router *gin.RouterGroup, highFreqService *services.HighFreqService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewHighFreqIPHandler(highFreqService)
	hf := router.Group("/high-frequency-ips")
	hf.Use(AuthRequired())
	{
		hf.GET("", requirePerm("high-freq-ips:list"), handler.List)
		hf.GET("/config", requirePerm("high-freq-ips:config"), handler.GetConfig)
		hf.PUT("/config", requirePerm("high-freq-ips:config"), handler.UpdateConfig)
		hf.POST("/refresh", requirePerm("high-freq-ips:list"), handler.Refresh)
	}
}

// DesensitizeRuleHandler handles desensitize rule endpoints
type DesensitizeRuleHandler struct {
	desensitizeService *services.DesensitizeService
}

// NewDesensitizeRuleHandler creates a new DesensitizeRuleHandler
func NewDesensitizeRuleHandler(desensitizeService *services.DesensitizeService) *DesensitizeRuleHandler {
	return &DesensitizeRuleHandler{desensitizeService: desensitizeService}
}

// List handles GET /api/desensitize-rules
func (h *DesensitizeRuleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	rules, total, err := h.desensitizeService.ListRules(page, pageSize, keyword)
	if err != nil {
		response.InternalError(c, "failed to list rules")
		return
	}
	response.SuccessWithPage(c, rules, total, page, pageSize)
}

// GetByID handles GET /api/desensitize-rules/:id
func (h *DesensitizeRuleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid desensitize rule id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var rule models.DesensitizeRule
	if err := db.First(&rule, id).Error; err != nil {
		response.NotFound(c, "desensitize rule not found")
		return
	}

	response.Success(c, rule)
}

// ToggleStatus handles PUT /api/desensitize-rules/:id/status
func (h *DesensitizeRuleHandler) ToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid desensitize rule id")
		return
	}
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	if err := db.Model(&models.DesensitizeRule{}).Where("id = ?", id).Update("enabled", req.Enabled).Error; err != nil {
		response.InternalError(c, "failed to update status")
		return
	}
	response.Success(c, nil)
}

// Create handles POST /api/desensitize-rules
func (h *DesensitizeRuleHandler) Create(c *gin.Context) {
	var rule models.DesensitizeRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		writeAuditLog(c, "create", "desensitize-rule", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if err := h.desensitizeService.CreateRule(&rule); err != nil {
		writeAuditLog(c, "create", "desensitize-rule", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}
	writeAuditLog(c, "create", "desensitize-rule", strconv.FormatUint(uint64(rule.ID), 10), "success", "desensitize rule created")
	response.Created(c, rule)
}

// Update handles PUT /api/desensitize-rules/:id
func (h *DesensitizeRuleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "desensitize-rule", c.Param("id"), "failure", "invalid desensitize rule id")
		response.BadRequest(c, "invalid desensitize rule id")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeAuditLog(c, "update", "desensitize-rule", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if err := h.desensitizeService.UpdateRule(uint(id), updates); err != nil {
		writeAuditLog(c, "update", "desensitize-rule", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}
	writeAuditLog(c, "update", "desensitize-rule", c.Param("id"), "success", "desensitize rule updated")
	response.SuccessWithMessage(c, "rule updated", nil)
}

// Delete handles DELETE /api/desensitize-rules/:id
func (h *DesensitizeRuleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "delete", "desensitize-rule", c.Param("id"), "failure", "invalid desensitize rule id")
		response.BadRequest(c, "invalid desensitize rule id")
		return
	}
	if err := h.desensitizeService.DeleteRule(uint(id)); err != nil {
		writeAuditLog(c, "delete", "desensitize-rule", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to delete rule")
		return
	}
	writeAuditLog(c, "delete", "desensitize-rule", c.Param("id"), "success", "desensitize rule deleted")
	response.SuccessWithMessage(c, "rule deleted", nil)
}

// RegisterDesensitizeRoutes registers desensitize rule routes
func RegisterDesensitizeRoutes(router *gin.RouterGroup, desensitizeService *services.DesensitizeService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewDesensitizeRuleHandler(desensitizeService)
	dr := router.Group("/desensitize-rules")
	dr.Use(AuthRequired())
	{
		dr.GET("", requirePerm("desensitize-rules:list"), handler.List)
		dr.GET("/:id", requirePerm("desensitize-rules:list"), handler.GetByID)
		dr.POST("", requirePerm("desensitize-rules:create"), handler.Create)
		dr.PUT("/:id", requirePerm("desensitize-rules:update"), handler.Update)
		dr.PUT("/:id/status", requirePerm("desensitize-rules:update"), handler.ToggleStatus)
		dr.DELETE("/:id", requirePerm("desensitize-rules:delete"), handler.Delete)
	}
}
