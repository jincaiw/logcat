package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
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
	sourceIP := c.Query("sourceIp")
	eventType := c.Query("eventType")
	severity := c.Query("severity")
	filterStatus := c.Query("filterStatus")
	keyword := c.Query("keyword")

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
		query = query.Where("raw_message LIKE ?", "%"+keyword+"%")
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

// Cleanup handles DELETE /api/logs/cleanup
func (h *LogHandler) Cleanup(c *gin.Context) {
	if err := h.cleanupService.RunCleanupNow(); err != nil {
		response.InternalError(c, "cleanup failed: "+err.Error())
		return
	}
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	channelType := c.Query("channelType")

	records, total, err := h.alertService.ListAlertRecords(page, pageSize, status, channelType)
	if err != nil {
		response.InternalError(c, "failed to list alert records")
		return
	}

	response.SuccessWithPage(c, records, total, page, pageSize)
}

// CreateDispositionRequest is the request body for creating a disposition
type CreateDispositionRequest struct {
	Status string `json:"status" binding:"required"`
	Note   string `json:"note"`
}

// CreateDisposition handles POST /api/alerts/:id/dispositions
func (h *AlertRecordHandler) CreateDisposition(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req CreateDispositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	alertID := uint(id)
	operatorID := uint(1) // TODO: get from context
	operatorName := "admin"

	disposition, err := h.alertService.CreateDisposition(&alertID, nil, req.Status, req.Note, &operatorID, operatorName)
	if err != nil {
		response.InternalError(c, "failed to create disposition")
		return
	}

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

// RegisterAlertRecordRoutes registers alert record routes
func RegisterAlertRecordRoutes(router *gin.RouterGroup, alertService *services.AlertService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewAlertRecordHandler(alertService)
	alerts := router.Group("/alerts")
	alerts.Use(AuthRequired())
	{
		alerts.GET("", requirePerm("alerts:list"), handler.List)
		alerts.POST("/:id/dispositions", requirePerm("alerts:disposition:create"), handler.CreateDisposition)
		alerts.GET("/:id/dispositions", requirePerm("alerts:disposition:list"), handler.ListDispositions)
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

// RegisterAggregatedAlertRoutes registers aggregated alert routes
func RegisterAggregatedAlertRoutes(router *gin.RouterGroup, aggregateService *services.AggregateService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewAggregatedAlertHandler(aggregateService)
	aa := router.Group("/aggregated-alerts")
	aa.Use(AuthRequired())
	{
		aa.GET("", requirePerm("aggregated-alerts:list"), handler.List)
		aa.GET("/:id/logs", requirePerm("aggregated-alerts:logs"), handler.GetLogs)
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
	ips := h.highFreqService.GetHighFreqIPs()
	response.Success(c, ips)
}

// UpdateConfigRequest is the config update request body
type HighFreqConfigRequest struct {
	Threshold     int `json:"threshold"`
	WindowSeconds int `json:"windowSeconds"`
}

// UpdateConfig handles PUT /api/high-frequency-ips/config
func (h *HighFreqIPHandler) UpdateConfig(c *gin.Context) {
	var req HighFreqConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	h.highFreqService.UpdateConfig(req.Threshold, req.WindowSeconds)
	response.Success(c, h.highFreqService.GetConfig())
}

// RegisterHighFreqIPRoutes registers high frequency IP routes
func RegisterHighFreqIPRoutes(router *gin.RouterGroup, highFreqService *services.HighFreqService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewHighFreqIPHandler(highFreqService)
	hf := router.Group("/high-frequency-ips")
	hf.Use(AuthRequired())
	{
		hf.GET("", requirePerm("high-freq-ips:list"), handler.List)
		hf.PUT("/config", requirePerm("high-freq-ips:config"), handler.UpdateConfig)
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
	rules, err := h.desensitizeService.ListRules()
	if err != nil {
		response.InternalError(c, "failed to list rules")
		return
	}
	response.Success(c, rules)
}

// Create handles POST /api/desensitize-rules
func (h *DesensitizeRuleHandler) Create(c *gin.Context) {
	var rule models.DesensitizeRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if err := h.desensitizeService.CreateRule(&rule); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, rule)
}

// Update handles PUT /api/desensitize-rules/:id
func (h *DesensitizeRuleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if err := h.desensitizeService.UpdateRule(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "rule updated", nil)
}

// Delete handles DELETE /api/desensitize-rules/:id
func (h *DesensitizeRuleHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.desensitizeService.DeleteRule(uint(id)); err != nil {
		response.InternalError(c, "failed to delete rule")
		return
	}
	response.SuccessWithMessage(c, "rule deleted", nil)
}

// RegisterDesensitizeRoutes registers desensitize rule routes
func RegisterDesensitizeRoutes(router *gin.RouterGroup, desensitizeService *services.DesensitizeService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewDesensitizeRuleHandler(desensitizeService)
	dr := router.Group("/desensitize-rules")
	dr.Use(AuthRequired())
	{
		dr.GET("", requirePerm("desensitize-rules:list"), handler.List)
		dr.POST("", requirePerm("desensitize-rules:create"), handler.Create)
		dr.PUT("/:id", requirePerm("desensitize-rules:update"), handler.Update)
		dr.DELETE("/:id", requirePerm("desensitize-rules:delete"), handler.Delete)
	}
}