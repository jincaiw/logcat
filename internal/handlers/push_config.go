package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/crypto"
	"github.com/logcat/logcat/pkg/response"
)

// validatePushConfig validates push config fields based on type
func validatePushConfig(config *models.PushConfig) string {
	validTypes := map[string]bool{"http": true, "email": true, "syslog": true}
	if !validTypes[config.Type] {
		return "type must be one of: http, email, syslog"
	}
	switch config.Type {
	case "http":
		if config.URL == "" {
			return "url is required when type is http"
		}
	case "email":
		if config.SMTPHost == "" {
			return "smtpHost is required when type is email"
		}
		if config.SMTPPort == 0 {
			return "smtpPort is required when type is email"
		}
	case "syslog":
		if config.SyslogHost == "" {
			return "syslogHost is required when type is syslog"
		}
		if config.SyslogPort == 0 {
			return "syslogPort is required when type is syslog"
		}
	}
	return ""
}

// validatePushConfigUpdates validates push config update fields based on type
func validatePushConfigUpdates(updates map[string]interface{}) string {
	validTypes := map[string]bool{"http": true, "email": true, "syslog": true}
	if t, ok := updates["type"].(string); ok {
		if !validTypes[t] {
			return "type must be one of: http, email, syslog"
		}
	}
	// Determine the effective type: either from the update or we can't validate further
	t, _ := updates["type"].(string)
	if t == "" {
		// Type not being updated; cannot validate type-specific fields without DB lookup
		return ""
	}
	switch t {
	case "http":
		if url, ok := updates["url"].(string); ok && url == "" {
			return "url is required when type is http"
		}
	case "email":
		if host, ok := updates["smtpHost"].(string); ok && host == "" {
			return "smtpHost is required when type is email"
		}
		if port, ok := updates["smtpPort"]; ok {
			switch v := port.(type) {
			case float64:
				if v == 0 {
					return "smtpPort is required when type is email"
				}
			case int:
				if v == 0 {
					return "smtpPort is required when type is email"
				}
			}
		}
	case "syslog":
		if host, ok := updates["syslogHost"].(string); ok && host == "" {
			return "syslogHost is required when type is syslog"
		}
		if port, ok := updates["syslogPort"]; ok {
			switch v := port.(type) {
			case float64:
				if v == 0 {
					return "syslogPort is required when type is syslog"
				}
			case int:
				if v == 0 {
					return "syslogPort is required when type is syslog"
				}
			}
		}
	}
	return ""
}

// PushConfigHandler handles push config endpoints
type PushConfigHandler struct {
	pushService *services.PushService
}

// NewPushConfigHandler creates a new PushConfigHandler
func NewPushConfigHandler(pushService *services.PushService) *PushConfigHandler {
	return &PushConfigHandler{pushService: pushService}
}

// List handles GET /api/push-configs
func (h *PushConfigHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	query := db.Model(&models.PushConfig{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count push configs")
		return
	}

	var configs []models.PushConfig
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&configs).Error; err != nil {
		response.InternalError(c, "failed to list push configs")
		return
	}
	response.SuccessWithPage(c, configs, total, page, pageSize)
}

// GetByID handles GET /api/push-configs/:id
func (h *PushConfigHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid push config id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var config models.PushConfig
	if err := db.First(&config, id).Error; err != nil {
		response.NotFound(c, "push config not found")
		return
	}

	response.Success(c, config)
}

// Create handles POST /api/push-configs
func (h *PushConfigHandler) Create(c *gin.Context) {
	var config models.PushConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		writeAuditLog(c, "create", "push-config", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if msg := validatePushConfig(&config); msg != "" {
		writeAuditLog(c, "create", "push-config", "", "failure", msg)
		response.BadRequest(c, msg)
		return
	}
	if config.Token != "" {
		enc, err := crypto.Encrypt(config.Token)
		if err != nil {
			response.InternalError(c, "failed to encrypt token")
			return
		}
		config.Token = enc
	}
	if config.SMTPPassword != "" {
		enc, err := crypto.Encrypt(config.SMTPPassword)
		if err != nil {
			response.InternalError(c, "failed to encrypt smtp password")
			return
		}
		config.SMTPPassword = enc
	}
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	if err := db.Create(&config).Error; err != nil {
		writeAuditLog(c, "create", "push-config", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}
	writeAuditLog(c, "create", "push-config", strconv.FormatUint(uint64(config.ID), 10), "success", "push config created")
	response.Created(c, config)
}

// Update handles PUT /api/push-configs/:id
func (h *PushConfigHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "push-config", c.Param("id"), "failure", "invalid push config id")
		response.BadRequest(c, "invalid push config id")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeAuditLog(c, "update", "push-config", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if msg := validatePushConfigUpdates(updates); msg != "" {
		writeAuditLog(c, "update", "push-config", c.Param("id"), "failure", msg)
		response.BadRequest(c, msg)
		return
	}

	allowedFields := map[string]bool{
		"name": true, "type": true, "enabled": true,
		"url": true, "method": true, "timeout": true, "retry_count": true, "retry_delay": true,
		"notes_ids": true, "headers": true, "body_template": true,
		"success_status_codes": true, "success_body_keyword": true,
		"auth_type": true, "token": true, "content_type": true,
		"retry_on_status_codes": true, "max_response_log_size": true,
		"smtp_host": true, "smtp_port": true, "smtp_username": true, "smtp_password": true,
		"from_address": true, "to_addresses": true,
		"subject_template": true, "email_body_template": true,
		"syslog_host": true, "syslog_port": true, "syslog_protocol": true,
		"syslog_format": true, "syslog_fields": true,
	}
	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			filtered[k] = v
		}
	}
	if tokenVal, ok := filtered["token"].(string); ok && tokenVal != "" {
		enc, err := crypto.Encrypt(tokenVal)
		if err != nil {
			response.InternalError(c, "failed to encrypt token")
			return
		}
		filtered["token"] = enc
	}
	if pwdVal, ok := filtered["smtp_password"].(string); ok && pwdVal != "" {
		enc, err := crypto.Encrypt(pwdVal)
		if err != nil {
			response.InternalError(c, "failed to encrypt smtp password")
			return
		}
		filtered["smtp_password"] = enc
	}
	if len(filtered) == 0 {
		response.BadRequest(c, "no valid fields to update")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	if err := db.Model(&models.PushConfig{}).Where("id = ?", id).Updates(filtered).Error; err != nil {
		writeAuditLog(c, "update", "push-config", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}
	writeAuditLog(c, "update", "push-config", c.Param("id"), "success", "push config updated")
	response.SuccessWithMessage(c, "push config updated", nil)
}

// Delete handles DELETE /api/push-configs/:id
func (h *PushConfigHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "delete", "push-config", c.Param("id"), "failure", "invalid push config id")
		response.BadRequest(c, "invalid push config id")
		return
	}
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	result := db.Delete(&models.PushConfig{}, id)
	if result.Error != nil {
		writeAuditLog(c, "delete", "push-config", c.Param("id"), "failure", result.Error.Error())
		response.BadRequest(c, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		writeAuditLog(c, "delete", "push-config", c.Param("id"), "failure", "push config not found")
		response.NotFound(c, "push config not found")
		return
	}
	writeAuditLog(c, "delete", "push-config", c.Param("id"), "success", "push config deleted")
	response.SuccessWithMessage(c, "push config deleted", nil)
}

// Test handles POST /api/push-configs/:id/test
func (h *PushConfigHandler) Test(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "test", "push-config", c.Param("id"), "failure", "invalid push config id")
		response.BadRequest(c, "invalid push config id")
		return
	}
	result, err := h.pushService.TestPush(uint(id))
	if err != nil {
		writeAuditLog(c, "test", "push-config", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}
	detail := result.Summary
	if detail == "" {
		detail = "push config tested"
	}
	status := "success"
	if !result.Success {
		status = "failure"
	}
	writeAuditLog(c, "test", "push-config", c.Param("id"), status, detail)
	response.Success(c, result)
}

// RegisterPushConfigRoutes registers push config routes
func RegisterPushConfigRoutes(router *gin.RouterGroup, pushService *services.PushService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewPushConfigHandler(pushService)
	pc := router.Group("/push-configs")
	pc.Use(AuthRequired())
	{
		pc.GET("", requirePerm("push-configs:list"), handler.List)
		pc.GET("/:id", requirePerm("push-configs:list"), handler.GetByID)
		pc.POST("", requirePerm("push-configs:create"), handler.Create)
		pc.PUT("/:id", requirePerm("push-configs:update"), handler.Update)
		pc.DELETE("/:id", requirePerm("push-configs:delete"), handler.Delete)
		pc.POST("/:id/test", requirePerm("push-configs:test"), handler.Test)
	}
}

// AlertRuleHandler handles alert rule endpoints
type AlertRuleHandler struct{}

// NewAlertRuleHandler creates a new AlertRuleHandler
func NewAlertRuleHandler() *AlertRuleHandler {
	return &AlertRuleHandler{}
}

// List handles GET /api/alert-rules
func (h *AlertRuleHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	query := db.Model(&models.AlertRule{}).Preload("FilterPolicy").Preload("PushConfig").Preload("OutputTemplate")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count alert rules")
		return
	}

	var rules []models.AlertRule
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&rules).Error; err != nil {
		response.InternalError(c, "failed to list alert rules")
		return
	}
	response.SuccessWithPage(c, rules, total, page, pageSize)
}

// GetByID handles GET /api/alert-rules/:id
func (h *AlertRuleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid alert rule id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var rule models.AlertRule
	if err := db.Preload("FilterPolicy").Preload("PushConfig").Preload("OutputTemplate").First(&rule, id).Error; err != nil {
		response.NotFound(c, "alert rule not found")
		return
	}

	response.Success(c, rule)
}

// ToggleStatus handles PUT /api/alert-rules/:id/status
func (h *AlertRuleHandler) ToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid alert rule id")
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
	if err := db.Model(&models.AlertRule{}).Where("id = ?", id).Update("enabled", req.Enabled).Error; err != nil {
		response.InternalError(c, "failed to update status")
		return
	}
	response.Success(c, nil)
}

// Create handles POST /api/alert-rules
func (h *AlertRuleHandler) Create(c *gin.Context) {
	var rule models.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		writeAuditLog(c, "create", "alert-rule", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if rule.FilterPolicyID == nil || *rule.FilterPolicyID == 0 {
		writeAuditLog(c, "create", "alert-rule", "", "failure", "filterPolicyId is required")
		response.BadRequest(c, "filterPolicyId is required")
		return
	}
	if rule.PushConfigID == nil || *rule.PushConfigID == 0 {
		writeAuditLog(c, "create", "alert-rule", "", "failure", "pushConfigId is required")
		response.BadRequest(c, "pushConfigId is required")
		return
	}
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	if err := db.Create(&rule).Error; err != nil {
		writeAuditLog(c, "create", "alert-rule", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}
	writeAuditLog(c, "create", "alert-rule", strconv.FormatUint(uint64(rule.ID), 10), "success", "alert rule created")
	response.Created(c, rule)
}

// Update handles PUT /api/alert-rules/:id
func (h *AlertRuleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "alert-rule", c.Param("id"), "failure", "invalid alert rule id")
		response.BadRequest(c, "invalid alert rule id")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeAuditLog(c, "update", "alert-rule", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	allowedFields := map[string]bool{
		"name": true, "enabled": true, "severity": true,
		"filter_policy_id": true, "push_config_id": true, "output_template_id": true,
		"throttle_minutes": true, "cooldown_minutes": true,
		"description": true, "condition": true,
	}
	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			filtered[k] = v
		}
	}
	if len(filtered) == 0 {
		response.BadRequest(c, "no valid fields to update")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	if err := db.Model(&models.AlertRule{}).Where("id = ?", id).Updates(filtered).Error; err != nil {
		writeAuditLog(c, "update", "alert-rule", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}
	writeAuditLog(c, "update", "alert-rule", c.Param("id"), "success", "alert rule updated")
	response.SuccessWithMessage(c, "alert rule updated", nil)
}

// Delete handles DELETE /api/alert-rules/:id
func (h *AlertRuleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "delete", "alert-rule", c.Param("id"), "failure", "invalid alert rule id")
		response.BadRequest(c, "invalid alert rule id")
		return
	}
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	result := db.Delete(&models.AlertRule{}, id)
	if result.Error != nil {
		writeAuditLog(c, "delete", "alert-rule", c.Param("id"), "failure", result.Error.Error())
		response.BadRequest(c, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		writeAuditLog(c, "delete", "alert-rule", c.Param("id"), "failure", "alert rule not found")
		response.NotFound(c, "alert rule not found")
		return
	}
	writeAuditLog(c, "delete", "alert-rule", c.Param("id"), "success", "alert rule deleted")
	response.SuccessWithMessage(c, "alert rule deleted", nil)
}

// RegisterAlertRuleRoutes registers alert rule routes
func RegisterAlertRuleRoutes(router *gin.RouterGroup, requirePerm func(string) gin.HandlerFunc) {
	handler := NewAlertRuleHandler()
	ar := router.Group("/alert-rules")
	ar.Use(AuthRequired())
	{
		ar.GET("", requirePerm("alert-rules:list"), handler.List)
		ar.GET("/:id", requirePerm("alert-rules:list"), handler.GetByID)
		ar.POST("", requirePerm("alert-rules:create"), handler.Create)
		ar.PUT("/:id", requirePerm("alert-rules:update"), handler.Update)
		ar.PUT("/:id/status", requirePerm("alert-rules:update"), handler.ToggleStatus)
		ar.DELETE("/:id", requirePerm("alert-rules:delete"), handler.Delete)
	}
}
