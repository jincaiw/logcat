package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/middleware"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
	"gorm.io/gorm"
)

// StatsHandler handles statistics and dashboard endpoints
type StatsHandler struct {
	statsService *services.StatsService
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// Query handles GET /api/stats
func (h *StatsHandler) Query(c *gin.Context) {
	field := c.Query("field")
	if strings.TrimSpace(field) == "" {
		response.BadRequest(c, "field is required")
		return
	}

	topN := 20
	if v := c.Query("topN"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			topN = parsed
		}
	}

	items, err := h.statsService.QueryFieldStats(field, c.Query("startTime"), c.Query("endTime"), topN)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, items)
}

// FieldStats handles GET /api/stats/fields
func (h *StatsHandler) FieldStats(c *gin.Context) {
	stats, err := h.statsService.GetFieldStats()
	if err != nil {
		response.InternalError(c, "failed to get field stats")
		return
	}
	response.Success(c, stats)
}

// AvailableFields handles GET /api/stats/available-fields
func (h *StatsHandler) AvailableFields(c *gin.Context) {
	fields := h.statsService.AvailableFields()
	response.Success(c, fields)
}

// ExportCSV handles GET /api/stats/export-csv
func (h *StatsHandler) ExportCSV(c *gin.Context) {
	field := c.Query("field")
	if strings.TrimSpace(field) == "" {
		response.BadRequest(c, "field is required")
		return
	}

	topN := 20
	if v := c.Query("topN"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			topN = parsed
		}
	}

	items, err := h.statsService.QueryFieldStats(field, c.Query("startTime"), c.Query("endTime"), topN)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var builder strings.Builder
	builder.WriteString("field,value,count,percentage,lastSeenAt\n")
	for _, item := range items {
		builder.WriteString(fmt.Sprintf("%s,%s,%d,%.6f,%s\n",
			sanitizeCSVField(item.Field),
			sanitizeCSVField(item.Value),
			item.Count,
			item.Percentage,
			item.LastSeenAt,
		))
	}

	response.Success(c, gin.H{
		"url": "data:text/csv;charset=utf-8," + url.PathEscape(builder.String()),
	})
}

// IPList handles GET /api/stats/ip-list
func (h *StatsHandler) IPList(c *gin.Context) {
	ips, err := h.statsService.GetIPList(c.Query("startTime"), c.Query("endTime"))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"ips": ips})
}

// Dashboard handles GET /api/dashboard (aggregated dashboard data)
func (h *StatsHandler) Dashboard(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		response.InternalError(c, "failed to get dashboard stats")
		return
	}
	response.Success(c, stats)
}

// RegisterStatsRoutes registers stats routes
func RegisterStatsRoutes(router *gin.RouterGroup, statsService *services.StatsService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewStatsHandler(statsService)
	s := router.Group("/stats")
	s.Use(AuthRequired())
	{
		s.GET("", requirePerm("stats:fields"), handler.Query)
		s.GET("/fields", requirePerm("stats:fields"), handler.FieldStats)
		s.GET("/available-fields", requirePerm("stats:available-fields"), handler.AvailableFields)
		s.GET("/export-csv", requirePerm("stats:fields"), handler.ExportCSV)
		s.GET("/ip-list", requirePerm("stats:fields"), handler.IPList)
	}
}

// DashboardHandler handles dashboard endpoints
type DashboardHandler struct {
	statsService *services.StatsService
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(statsService *services.StatsService) *DashboardHandler {
	return &DashboardHandler{statsService: statsService}
}

// GetDashboard handles GET /api/dashboard (redirect from dashboard endpoint)
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		response.InternalError(c, "failed to get dashboard stats")
		return
	}
	response.Success(c, stats)
}

// RegisterDashboardRoutes registers dashboard routes
func RegisterDashboardRoutes(router *gin.RouterGroup, statsService *services.StatsService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewDashboardHandler(statsService)
	d := router.Group("/dashboard")
	d.Use(AuthRequired())
	{
		d.GET("", requirePerm("dashboard:view"), handler.GetDashboard)
		d.GET("/stats", requirePerm("dashboard:view"), handler.GetDashboard)
	}
}

// ImportExportHandler handles import/export endpoints
type ImportExportHandler struct{}

// NewImportExportHandler creates a new ImportExportHandler
func NewImportExportHandler() *ImportExportHandler {
	return &ImportExportHandler{}
}

const importExportSchemaVersion = "1.0"

type importExportEnvelope struct {
	Version      string          `json:"version"`
	ResourceType string          `json:"resourceType"`
	ExportedAt   string          `json:"exportedAt"`
	Items        json.RawMessage `json:"items"`
}

type importExportResult struct {
	ResourceType string   `json:"resourceType"`
	Version      string   `json:"version"`
	Created      int      `json:"created"`
	Updated      int      `json:"updated"`
	Failed       int      `json:"failed"`
	Errors       []string `json:"errors"`
}

func decodeImportItems[T any](c *gin.Context, expectedResourceType string) ([]T, string, error) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read request body: %w", err)
	}
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return nil, "", fmt.Errorf("empty request body")
	}

	if raw[0] == '[' {
		var items []T
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil, "", err
		}
		return items, "legacy", nil
	}

	var envelope importExportEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, "", err
	}
	if envelope.ResourceType != "" && envelope.ResourceType != expectedResourceType {
		return nil, "", fmt.Errorf("resourceType mismatch: expected %s", expectedResourceType)
	}
	if envelope.Version != "" && envelope.Version != importExportSchemaVersion {
		return nil, envelope.Version, fmt.Errorf("unsupported version: %s", envelope.Version)
	}

	var items []T
	if err := json.Unmarshal(envelope.Items, &items); err != nil {
		return nil, envelope.Version, err
	}
	return items, envelope.Version, nil
}

func buildExportPayload(resourceType string, items interface{}) gin.H {
	payload, _ := json.MarshalIndent(importExportEnvelope{
		Version:      importExportSchemaVersion,
		ResourceType: resourceType,
		ExportedAt:   time.Now().Format(time.RFC3339),
		Items:        mustJSON(items),
	}, "", "  ")

	count := 0
	value := reflect.ValueOf(items)
	if value.IsValid() && value.Kind() == reflect.Slice {
		count = value.Len()
	}

	return gin.H{
		"url":          "data:application/json;charset=utf-8," + url.PathEscape(string(payload)),
		"version":      importExportSchemaVersion,
		"resourceType": resourceType,
		"count":        count,
	}
}

func mustJSON(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}

// sanitizeCSVField neutralizes CSV injection (formula injection) by prefixing
// dangerous leading characters with a single quote, then escaping commas and
// double quotes. It prevents attacks where a value starting with =, +, -, @,
// TAB, or CR is interpreted as a formula by spreadsheet applications.
func sanitizeCSVField(v string) string {
	v = strings.ReplaceAll(v, ",", " ")
	v = strings.ReplaceAll(v, "\"", "\"\"")
	if len(v) > 0 {
		switch v[0] {
		case '=', '+', '-', '@', '\t', '\r':
			v = "'" + v
		}
	}
	return v
}

func auditImportExport(c *gin.Context, action, resourceType, result, detail string) {
	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	var uid *uint
	if userID > 0 {
		uid = &userID
	}
	_ = middleware.AuditLogWriter(
		uid,
		username,
		action,
		resourceType,
		"",
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		result,
		detail,
	)
}

// ImportParseTemplates handles POST /api/import/parse-templates
func (h *ImportExportHandler) ImportParseTemplates(c *gin.Context) {
	items, version, err := decodeImportItems[models.ParseTemplate](c, "parse-templates")
	if err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	result := importExportResult{ResourceType: "parse-templates", Version: version, Errors: []string{}}
	for _, tmpl := range items {
		tmpl.ID = 0
		var existing models.ParseTemplate
		err := db.Where("name = ?", tmpl.Name).First(&existing).Error
		if err == nil {
			tmpl.ID = existing.ID
			tmpl.CreatedAt = existing.CreatedAt
			if saveErr := db.Save(&tmpl).Error; saveErr != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", tmpl.Name, saveErr))
				continue
			}
			result.Updated++
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", tmpl.Name, err))
			continue
		}
		if createErr := db.Create(&tmpl).Error; createErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", tmpl.Name, createErr))
			continue
		}
		result.Created++
	}
	auditImportExport(c, "import", "parse-templates", "success", fmt.Sprintf("created=%d updated=%d failed=%d", result.Created, result.Updated, result.Failed))
	response.Success(c, result)
}

// ImportFilterPolicies handles POST /api/import/filter-policies
func (h *ImportExportHandler) ImportFilterPolicies(c *gin.Context) {
	items, version, err := decodeImportItems[models.FilterPolicy](c, "filter-policies")
	if err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	result := importExportResult{ResourceType: "filter-policies", Version: version, Errors: []string{}}
	for _, policy := range items {
		policy.ID = 0
		var existing models.FilterPolicy
		err := db.Where("name = ?", policy.Name).First(&existing).Error
		if err == nil {
			policy.ID = existing.ID
			policy.CreatedAt = existing.CreatedAt
			if saveErr := db.Save(&policy).Error; saveErr != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", policy.Name, saveErr))
				continue
			}
			result.Updated++
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", policy.Name, err))
			continue
		}
		if createErr := db.Create(&policy).Error; createErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", policy.Name, createErr))
			continue
		}
		result.Created++
	}
	auditImportExport(c, "import", "filter-policies", "success", fmt.Sprintf("created=%d updated=%d failed=%d", result.Created, result.Updated, result.Failed))
	response.Success(c, result)
}

// ImportPushConfigs handles POST /api/import/push-configs
func (h *ImportExportHandler) ImportPushConfigs(c *gin.Context) {
	items, version, err := decodeImportItems[models.PushConfig](c, "push-configs")
	if err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	result := importExportResult{ResourceType: "push-configs", Version: version, Errors: []string{}}
	for _, cfg := range items {
		cfg.ID = 0
		var existing models.PushConfig
		err := db.Where("name = ?", cfg.Name).First(&existing).Error
		if err == nil {
			cfg.ID = existing.ID
			cfg.CreatedAt = existing.CreatedAt
			if saveErr := db.Save(&cfg).Error; saveErr != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", cfg.Name, saveErr))
				continue
			}
			result.Updated++
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", cfg.Name, err))
			continue
		}
		if createErr := db.Create(&cfg).Error; createErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", cfg.Name, createErr))
			continue
		}
		result.Created++
	}
	auditImportExport(c, "import", "push-configs", "success", fmt.Sprintf("created=%d updated=%d failed=%d", result.Created, result.Updated, result.Failed))
	response.Success(c, result)
}

// ImportDeviceTemplates handles POST /api/import/device-templates
func (h *ImportExportHandler) ImportDeviceTemplates(c *gin.Context) {
	items, version, err := decodeImportItems[models.DeviceTemplate](c, "device-templates")
	if err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	result := importExportResult{ResourceType: "device-templates", Version: version, Errors: []string{}}
	for _, tmpl := range items {
		tmpl.ID = 0
		var existing models.DeviceTemplate
		err := db.Where("name = ?", tmpl.Name).First(&existing).Error
		if err == nil {
			tmpl.ID = existing.ID
			tmpl.CreatedAt = existing.CreatedAt
			if saveErr := db.Save(&tmpl).Error; saveErr != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", tmpl.Name, saveErr))
				continue
			}
			result.Updated++
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", tmpl.Name, err))
			continue
		}
		if createErr := db.Create(&tmpl).Error; createErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", tmpl.Name, createErr))
			continue
		}
		result.Created++
	}
	auditImportExport(c, "import", "device-templates", "success", fmt.Sprintf("created=%d updated=%d failed=%d", result.Created, result.Updated, result.Failed))
	response.Success(c, result)
}

// ExportParseTemplates handles GET /api/export/parse-templates
func (h *ImportExportHandler) ExportParseTemplates(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var templates []models.ParseTemplate
	if err := db.Find(&templates).Error; err != nil {
		response.InternalError(c, "failed to export parse templates")
		return
	}
	auditImportExport(c, "export", "parse-templates", "success", fmt.Sprintf("count=%d", len(templates)))
	response.Success(c, buildExportPayload("parse-templates", templates))
}

// ExportFilterPolicies handles GET /api/export/filter-policies
func (h *ImportExportHandler) ExportFilterPolicies(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var policies []models.FilterPolicy
	if err := db.Find(&policies).Error; err != nil {
		response.InternalError(c, "failed to export filter policies")
		return
	}
	auditImportExport(c, "export", "filter-policies", "success", fmt.Sprintf("count=%d", len(policies)))
	response.Success(c, buildExportPayload("filter-policies", policies))
}

// ExportPushConfigs handles GET /api/export/push-configs
func (h *ImportExportHandler) ExportPushConfigs(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var configs []models.PushConfig
	if err := db.Find(&configs).Error; err != nil {
		response.InternalError(c, "failed to export push configs")
		return
	}
	auditImportExport(c, "export", "push-configs", "success", fmt.Sprintf("count=%d", len(configs)))
	response.Success(c, buildExportPayload("push-configs", configs))
}

// ExportDeviceTemplates handles GET /api/export/device-templates
func (h *ImportExportHandler) ExportDeviceTemplates(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var templates []models.DeviceTemplate
	if err := db.Find(&templates).Error; err != nil {
		response.InternalError(c, "failed to export device templates")
		return
	}
	auditImportExport(c, "export", "device-templates", "success", fmt.Sprintf("count=%d", len(templates)))
	response.Success(c, buildExportPayload("device-templates", templates))
}

// ExportHistory handles GET /api/export/history
func (h *ImportExportHandler) ExportHistory(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)

	var total int64
	query := db.Model(&models.ExportHistory{})
	if resourceType := c.Query("resourceType"); resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count export history")
		return
	}

	var histories []models.ExportHistory
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&histories).Error; err != nil {
		response.InternalError(c, "failed to list export history")
		return
	}

	response.SuccessWithPage(c, histories, total, page, pageSize)
}

// ImportHistory handles GET /api/import/history
func (h *ImportExportHandler) ImportHistory(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)

	var total int64
	query := db.Model(&models.ImportHistory{})
	if resourceType := c.Query("resourceType"); resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count import history")
		return
	}

	var histories []models.ImportHistory
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&histories).Error; err != nil {
		response.InternalError(c, "failed to list import history")
		return
	}

	response.SuccessWithPage(c, histories, total, page, pageSize)
}

// RegisterImportExportRoutes registers import/export routes
func RegisterImportExportRoutes(router *gin.RouterGroup, requirePerm func(string) gin.HandlerFunc) {
	handler := NewImportExportHandler()

	imp := router.Group("/import")
	imp.Use(AuthRequired())
	{
		imp.POST("/parse-templates", requirePerm("import:parse-templates"), handler.ImportParseTemplates)
		imp.POST("/filter-policies", requirePerm("import:filter-policies"), handler.ImportFilterPolicies)
		imp.POST("/push-configs", requirePerm("import:push-configs"), handler.ImportPushConfigs)
		imp.POST("/device-templates", requirePerm("import:device-templates"), handler.ImportDeviceTemplates)
		imp.GET("/history", requirePerm("import:parse-templates"), handler.ImportHistory)
	}

	exp := router.Group("/export")
	exp.Use(AuthRequired())
	{
		exp.GET("/parse-templates", requirePerm("export:config"), handler.ExportParseTemplates)
		exp.GET("/filter-policies", requirePerm("export:config"), handler.ExportFilterPolicies)
		exp.GET("/push-configs", requirePerm("export:config"), handler.ExportPushConfigs)
		exp.GET("/device-templates", requirePerm("export:config"), handler.ExportDeviceTemplates)
		exp.GET("/history", requirePerm("export:config"), handler.ExportHistory)
	}
}
