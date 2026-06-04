package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/middleware"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
)

// ParseTemplateHandler handles parse template endpoints
type ParseTemplateHandler struct {
	parseService *services.ParseService
}

var allowedParseTypes = map[string]struct{}{
	"json":         {},
	"delimiter":    {},
	"kv":           {},
	"regex":        {},
	"syslog_json":  {},
	"sub_template": {},
}

// NewParseTemplateHandler creates a new ParseTemplateHandler
func NewParseTemplateHandler(parseService *services.ParseService) *ParseTemplateHandler {
	return &ParseTemplateHandler{parseService: parseService}
}

func validateParseTemplatePayload(template *models.ParseTemplate) error {
	parseType := strings.TrimSpace(template.ParseType)
	if parseType == "" {
		return fmt.Errorf("parseType is required")
	}
	if _, ok := allowedParseTypes[parseType]; !ok {
		return fmt.Errorf("unsupported parse type: %s", parseType)
	}
	if template.FieldMapping != "" && !json.Valid([]byte(template.FieldMapping)) {
		return fmt.Errorf("fieldMapping must be valid JSON")
	}
	if template.ValueTransform != "" && !json.Valid([]byte(template.ValueTransform)) {
		return fmt.Errorf("valueTransform must be valid JSON")
	}
	if template.SubTemplates != "" && !json.Valid([]byte(template.SubTemplates)) {
		return fmt.Errorf("subTemplates must be valid JSON")
	}
	if parseType == "regex" && strings.TrimSpace(template.HeaderRegex) == "" {
		return fmt.Errorf("headerRegex is required for regex parse type")
	}
	if (parseType == "delimiter" || parseType == "kv") && strings.TrimSpace(template.Delimiter) == "" {
		return fmt.Errorf("delimiter is required for %s parse type", parseType)
	}
	if parseType == "sub_template" && strings.TrimSpace(template.SubTemplates) == "" {
		return fmt.Errorf("subTemplates is required for sub_template parse type")
	}
	return nil
}

func validateParseTemplateUpdates(updates map[string]interface{}) error {
	template := &models.ParseTemplate{}
	if v, ok := updates["parseType"].(string); ok {
		template.ParseType = v
	}
	if v, ok := updates["fieldMapping"].(string); ok {
		template.FieldMapping = v
	}
	if v, ok := updates["valueTransform"].(string); ok {
		template.ValueTransform = v
	}
	if v, ok := updates["subTemplates"].(string); ok {
		template.SubTemplates = v
	}
	if v, ok := updates["headerRegex"].(string); ok {
		template.HeaderRegex = v
	}
	if v, ok := updates["delimiter"].(string); ok {
		template.Delimiter = v
	}

	if template.ParseType == "" {
		if v, ok := updates["parseType"]; ok && v != nil {
			return fmt.Errorf("parseType must be a string")
		}
		if raw, ok := updates["fieldMapping"]; ok && raw != nil {
			if str, ok := raw.(string); !ok || !json.Valid([]byte(str)) {
				return fmt.Errorf("fieldMapping must be valid JSON")
			}
		}
		if raw, ok := updates["valueTransform"]; ok && raw != nil {
			if str, ok := raw.(string); !ok || !json.Valid([]byte(str)) {
				return fmt.Errorf("valueTransform must be valid JSON")
			}
		}
		if raw, ok := updates["subTemplates"]; ok && raw != nil {
			if str, ok := raw.(string); !ok || !json.Valid([]byte(str)) {
				return fmt.Errorf("subTemplates must be valid JSON")
			}
		}
		return nil
	}

	return validateParseTemplatePayload(template)
}

func validateFilterPolicyPayload(policy *models.FilterPolicy) error {
	if strings.TrimSpace(policy.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if policy.Conditions != "" && !json.Valid([]byte(policy.Conditions)) {
		return fmt.Errorf("conditions must be valid JSON")
	}
	if logic := strings.ToUpper(strings.TrimSpace(policy.ConditionLogic)); logic != "" && logic != "AND" && logic != "OR" {
		return fmt.Errorf("conditionLogic must be AND or OR")
	}
	action := strings.ToLower(strings.TrimSpace(policy.Action))
	if action == "" {
		return fmt.Errorf("action is required")
	}
	if action != "keep" && action != "drop" {
		return fmt.Errorf("action must be keep or drop")
	}
	if policy.WhitelistEnabled && strings.TrimSpace(policy.WhitelistField) == "" {
		return fmt.Errorf("whitelistField is required when whitelist is enabled")
	}
	return nil
}

func validateFilterPolicyUpdates(updates map[string]interface{}) error {
	policy := &models.FilterPolicy{}
	if v, ok := updates["name"].(string); ok {
		policy.Name = v
	}
	if v, ok := updates["conditions"].(string); ok {
		policy.Conditions = v
	}
	if v, ok := updates["conditionLogic"].(string); ok {
		policy.ConditionLogic = v
	}
	if v, ok := updates["action"].(string); ok {
		policy.Action = v
	}
	if v, ok := updates["whitelistField"].(string); ok {
		policy.WhitelistField = v
	}
	if v, ok := updates["whitelistEnabled"].(bool); ok {
		policy.WhitelistEnabled = v
	}
	if raw, ok := updates["conditions"]; ok && raw != nil {
		if str, ok := raw.(string); !ok || !json.Valid([]byte(str)) {
			return fmt.Errorf("conditions must be valid JSON")
		}
	}
	if raw, ok := updates["conditionLogic"]; ok && raw != nil {
		if str, ok := raw.(string); !ok || (strings.ToUpper(strings.TrimSpace(str)) != "AND" && strings.ToUpper(strings.TrimSpace(str)) != "OR") {
			return fmt.Errorf("conditionLogic must be AND or OR")
		}
	}
	if raw, ok := updates["action"]; ok && raw != nil {
		if str, ok := raw.(string); !ok || (strings.ToLower(strings.TrimSpace(str)) != "keep" && strings.ToLower(strings.TrimSpace(str)) != "drop") {
			return fmt.Errorf("action must be keep or drop")
		}
	}
	if enabledRaw, ok := updates["whitelistEnabled"]; ok {
		if enabled, ok := enabledRaw.(bool); ok && enabled {
			if fieldRaw, ok := updates["whitelistField"].(string); !ok || strings.TrimSpace(fieldRaw) == "" {
				return fmt.Errorf("whitelistField is required when whitelist is enabled")
			}
		}
	}
	return nil
}

func writeTemplateAudit(c *gin.Context, action, resourceType, resourceID, result, detail string) {
	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	var uid *uint
	if userID > 0 {
		uid = &userID
	}
	_ = middleware.AuditLogWriter(uid, username, action, resourceType, resourceID, c.ClientIP(), c.GetHeader("User-Agent"), result, detail)
}

// List handles GET /api/parse-templates
func (h *ParseTemplateHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	query := db.Model(&models.ParseTemplate{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count parse templates")
		return
	}

	var templates []models.ParseTemplate
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&templates).Error; err != nil {
		response.InternalError(c, "failed to list parse templates")
		return
	}

	response.SuccessWithPage(c, templates, total, page, pageSize)
}

// GetByID handles GET /api/parse-templates/:id
func (h *ParseTemplateHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var template models.ParseTemplate
	if err := db.First(&template, id).Error; err != nil {
		response.NotFound(c, "parse template not found")
		return
	}

	response.Success(c, template)
}

// Create handles POST /api/parse-templates
func (h *ParseTemplateHandler) Create(c *gin.Context) {
	var template models.ParseTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		writeTemplateAudit(c, "create", "parse-template", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if err := validateParseTemplatePayload(&template); err != nil {
		writeTemplateAudit(c, "create", "parse-template", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&template).Error; err != nil {
		writeTemplateAudit(c, "create", "parse-template", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeTemplateAudit(c, "create", "parse-template", strconv.FormatUint(uint64(template.ID), 10), "success", "parse template created")
	response.Created(c, template)
}

// Update handles PUT /api/parse-templates/:id
func (h *ParseTemplateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeTemplateAudit(c, "update", "parse-template", c.Param("id"), "failure", "invalid template id")
		response.BadRequest(c, "invalid template id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeTemplateAudit(c, "update", "parse-template", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if err := validateParseTemplateUpdates(updates); err != nil {
		writeTemplateAudit(c, "update", "parse-template", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	allowedParseFields := map[string]bool{
		"name": true, "device_type": true, "parse_type": true,
		"header_regex": true, "delimiter": true, "field_mapping": true,
		"value_transform": true, "sample_log": true, "sub_templates": true, "enabled": true,
	}
	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedParseFields[k] {
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
	if err := db.Model(&models.ParseTemplate{}).Where("id = ?", id).Updates(filtered).Error; err != nil {
		writeTemplateAudit(c, "update", "parse-template", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeTemplateAudit(c, "update", "parse-template", c.Param("id"), "success", "parse template updated")
	response.SuccessWithMessage(c, "parse template updated", nil)
}

// Delete handles DELETE /api/parse-templates/:id
func (h *ParseTemplateHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	result := db.Delete(&models.ParseTemplate{}, id)
	if result.Error != nil {
		writeTemplateAudit(c, "delete", "parse-template", c.Param("id"), "failure", result.Error.Error())
		response.BadRequest(c, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		writeTemplateAudit(c, "delete", "parse-template", c.Param("id"), "failure", "template not found")
		response.NotFound(c, "parse template not found")
		return
	}
	writeTemplateAudit(c, "delete", "parse-template", c.Param("id"), "success", "parse template deleted")

	response.SuccessWithMessage(c, "parse template deleted", nil)
}

// TestRequest is the request body for test parsing
type TestParseRequest struct {
	TemplateID uint   `json:"templateId" binding:"required"`
	RawLog     string `json:"rawLog" binding:"required"`
}

// Test handles POST /api/parse-templates/test
func (h *ParseTemplateHandler) Test(c *gin.Context) {
	var req TestParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeTemplateAudit(c, "test", "parse-template", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if strings.TrimSpace(req.RawLog) == "" {
		writeTemplateAudit(c, "test", "parse-template", strconv.FormatUint(uint64(req.TemplateID), 10), "failure", "rawLog is required")
		response.BadRequest(c, "rawLog is required")
		return
	}

	result, err := h.parseService.TestParse(req.TemplateID, req.RawLog)
	if err != nil {
		writeTemplateAudit(c, "test", "parse-template", strconv.FormatUint(uint64(req.TemplateID), 10), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeTemplateAudit(c, "test", "parse-template", strconv.FormatUint(uint64(req.TemplateID), 10), "success", "parse template tested")
	response.Success(c, result)
}

// RegisterParseTemplateRoutes registers parse template routes
func RegisterParseTemplateRoutes(router *gin.RouterGroup, parseService *services.ParseService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewParseTemplateHandler(parseService)
	pt := router.Group("/parse-templates")
	pt.Use(AuthRequired())
	{
		pt.GET("", requirePerm("parse-templates:list"), handler.List)
		pt.GET("/:id", requirePerm("parse-templates:list"), handler.GetByID)
		pt.POST("", requirePerm("parse-templates:create"), handler.Create)
		pt.PUT("/:id", requirePerm("parse-templates:update"), handler.Update)
		pt.DELETE("/:id", requirePerm("parse-templates:delete"), handler.Delete)
		pt.POST("/test", requirePerm("parse-templates:test"), handler.Test)
	}
}

// FilterPolicyHandler handles filter policy endpoints
type FilterPolicyHandler struct {
	filterService *services.FilterService
}

// NewFilterPolicyHandler creates a new FilterPolicyHandler
func NewFilterPolicyHandler(filterService *services.FilterService) *FilterPolicyHandler {
	return &FilterPolicyHandler{filterService: filterService}
}

// List handles GET /api/filter-policies
func (h *FilterPolicyHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	query := db.Model(&models.FilterPolicy{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count filter policies")
		return
	}

	var policies []models.FilterPolicy
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&policies).Error; err != nil {
		response.InternalError(c, "failed to list filter policies")
		return
	}

	response.SuccessWithPage(c, policies, total, page, pageSize)
}

// GetByID handles GET /api/filter-policies/:id
func (h *FilterPolicyHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid policy id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var policy models.FilterPolicy
	if err := db.First(&policy, id).Error; err != nil {
		response.NotFound(c, "filter policy not found")
		return
	}

	response.Success(c, policy)
}

// Create handles POST /api/filter-policies
func (h *FilterPolicyHandler) Create(c *gin.Context) {
	var policy models.FilterPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		writeTemplateAudit(c, "create", "filter-policy", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if err := validateFilterPolicyPayload(&policy); err != nil {
		writeTemplateAudit(c, "create", "filter-policy", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&policy).Error; err != nil {
		writeTemplateAudit(c, "create", "filter-policy", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeTemplateAudit(c, "create", "filter-policy", strconv.FormatUint(uint64(policy.ID), 10), "success", "filter policy created")
	response.Created(c, policy)
}

// Update handles PUT /api/filter-policies/:id
func (h *FilterPolicyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeTemplateAudit(c, "update", "filter-policy", c.Param("id"), "failure", "invalid policy id")
		response.BadRequest(c, "invalid policy id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeTemplateAudit(c, "update", "filter-policy", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if err := validateFilterPolicyUpdates(updates); err != nil {
		writeTemplateAudit(c, "update", "filter-policy", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	allowedFilterFields := map[string]bool{
		"name": true, "device_id": true, "device_group_id": true,
		"parse_template_id": true, "conditions": true, "condition_logic": true,
		"whitelist_enabled": true, "whitelist_field": true, "whitelist_values": true,
		"action": true, "priority": true, "dedup_enabled": true, "dedup_window": true,
		"enabled": true,
	}
	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedFilterFields[k] {
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
	if err := db.Model(&models.FilterPolicy{}).Where("id = ?", id).Updates(filtered).Error; err != nil {
		writeTemplateAudit(c, "update", "filter-policy", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeTemplateAudit(c, "update", "filter-policy", c.Param("id"), "success", "filter policy updated")
	response.SuccessWithMessage(c, "filter policy updated", nil)
}

// Delete handles DELETE /api/filter-policies/:id
func (h *FilterPolicyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid policy id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	result := db.Delete(&models.FilterPolicy{}, id)
	if result.Error != nil {
		writeTemplateAudit(c, "delete", "filter-policy", c.Param("id"), "failure", result.Error.Error())
		response.BadRequest(c, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		writeTemplateAudit(c, "delete", "filter-policy", c.Param("id"), "failure", "policy not found")
		response.NotFound(c, "filter policy not found")
		return
	}
	writeTemplateAudit(c, "delete", "filter-policy", c.Param("id"), "success", "filter policy deleted")

	response.SuccessWithMessage(c, "filter policy deleted", nil)
}

// TestFilterRequest is the request body for test filtering
type TestFilterRequest struct {
	PolicyID   uint                   `json:"policyId" binding:"required"`
	ParsedData map[string]interface{} `json:"parsedData" binding:"required"`
}

// TestFilter handles POST /api/filter-policies/test
func (h *FilterPolicyHandler) TestFilter(c *gin.Context) {
	var req TestFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeTemplateAudit(c, "test", "filter-policy", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	if len(req.ParsedData) == 0 {
		writeTemplateAudit(c, "test", "filter-policy", strconv.FormatUint(uint64(req.PolicyID), 10), "failure", "parsedData is required")
		response.BadRequest(c, "parsedData is required")
		return
	}

	result, err := h.filterService.TestFilter(req.PolicyID, req.ParsedData)
	if err != nil {
		writeTemplateAudit(c, "test", "filter-policy", strconv.FormatUint(uint64(req.PolicyID), 10), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeTemplateAudit(c, "test", "filter-policy", strconv.FormatUint(uint64(req.PolicyID), 10), "success", "filter policy tested")
	response.Success(c, result)
}

// RegisterFilterPolicyRoutes registers filter policy routes
func RegisterFilterPolicyRoutes(router *gin.RouterGroup, filterService *services.FilterService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewFilterPolicyHandler(filterService)
	fp := router.Group("/filter-policies")
	fp.Use(AuthRequired())
	{
		fp.GET("", requirePerm("filter-policies:list"), handler.List)
		fp.GET("/:id", requirePerm("filter-policies:list"), handler.GetByID)
		fp.POST("", requirePerm("filter-policies:create"), handler.Create)
		fp.PUT("/:id", requirePerm("filter-policies:update"), handler.Update)
		fp.DELETE("/:id", requirePerm("filter-policies:delete"), handler.Delete)
		fp.POST("/test", requirePerm("filter-policies:test"), handler.TestFilter)
	}
}

// OutputTemplateHandler handles output template endpoints
type OutputTemplateHandler struct{}

// NewOutputTemplateHandler creates a new OutputTemplateHandler
func NewOutputTemplateHandler() *OutputTemplateHandler {
	return &OutputTemplateHandler{}
}

// List handles GET /api/output-templates
func (h *OutputTemplateHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	query := db.Model(&models.OutputTemplate{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count output templates")
		return
	}

	var templates []models.OutputTemplate
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&templates).Error; err != nil {
		response.InternalError(c, "failed to list output templates")
		return
	}

	response.SuccessWithPage(c, templates, total, page, pageSize)
}

// GetByID handles GET /api/output-templates/:id
func (h *OutputTemplateHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var template models.OutputTemplate
	if err := db.First(&template, id).Error; err != nil {
		response.NotFound(c, "output template not found")
		return
	}

	response.Success(c, template)
}

// Create handles POST /api/output-templates
func (h *OutputTemplateHandler) Create(c *gin.Context) {
	var template models.OutputTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		writeAuditLog(c, "create", "output-template", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&template).Error; err != nil {
		writeAuditLog(c, "create", "output-template", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "create", "output-template", strconv.FormatUint(uint64(template.ID), 10), "success", "output template created")
	response.Created(c, template)
}

// Update handles PUT /api/output-templates/:id
func (h *OutputTemplateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "output-template", c.Param("id"), "failure", "invalid template id")
		response.BadRequest(c, "invalid template id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeAuditLog(c, "update", "output-template", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	allowedOutputFields := map[string]bool{
		"name": true, "channel_type": true, "content": true,
		"fields": true, "device_type": true, "enabled": true,
	}
	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedOutputFields[k] {
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
	if err := db.Model(&models.OutputTemplate{}).Where("id = ?", id).Updates(filtered).Error; err != nil {
		writeAuditLog(c, "update", "output-template", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "update", "output-template", c.Param("id"), "success", "output template updated")
	response.SuccessWithMessage(c, "output template updated", nil)
}

// Delete handles DELETE /api/output-templates/:id
func (h *OutputTemplateHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	result := db.Delete(&models.OutputTemplate{}, id)
	if result.Error != nil {
		writeAuditLog(c, "delete", "output-template", c.Param("id"), "failure", result.Error.Error())
		response.BadRequest(c, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		writeAuditLog(c, "delete", "output-template", c.Param("id"), "failure", "output template not found")
		response.NotFound(c, "output template not found")
		return
	}

	writeAuditLog(c, "delete", "output-template", c.Param("id"), "success", "output template deleted")
	response.SuccessWithMessage(c, "output template deleted", nil)
}

// RegisterOutputTemplateRoutes registers output template routes
func RegisterOutputTemplateRoutes(router *gin.RouterGroup, requirePerm func(string) gin.HandlerFunc) {
	handler := NewOutputTemplateHandler()
	ot := router.Group("/output-templates")
	ot.Use(AuthRequired())
	{
		ot.GET("", requirePerm("output-templates:list"), handler.List)
		ot.GET("/:id", requirePerm("output-templates:list"), handler.GetByID)
		ot.POST("", requirePerm("output-templates:create"), handler.Create)
		ot.PUT("/:id", requirePerm("output-templates:update"), handler.Update)
		ot.DELETE("/:id", requirePerm("output-templates:delete"), handler.Delete)
	}
}
