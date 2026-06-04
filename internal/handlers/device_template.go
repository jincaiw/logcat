package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/pkg/response"
)

// DeviceTemplateHandler handles device template endpoints
type DeviceTemplateHandler struct{}

// NewDeviceTemplateHandler creates a new DeviceTemplateHandler
func NewDeviceTemplateHandler() *DeviceTemplateHandler {
	return &DeviceTemplateHandler{}
}

// List handles GET /api/device-templates
func (h *DeviceTemplateHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	query := db.Model(&models.DeviceTemplate{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count device templates")
		return
	}

	var templates []models.DeviceTemplate
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&templates).Error; err != nil {
		response.InternalError(c, "failed to list device templates")
		return
	}

	response.SuccessWithPage(c, templates, total, page, pageSize)
}

// GetByID handles GET /api/device-templates/:id
func (h *DeviceTemplateHandler) GetByID(c *gin.Context) {
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

	var template models.DeviceTemplate
	if err := db.First(&template, id).Error; err != nil {
		response.NotFound(c, "device template not found")
		return
	}

	response.Success(c, template)
}

// GetAll handles GET /api/device-templates/all
func (h *DeviceTemplateHandler) GetAll(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var templates []models.DeviceTemplate
	if err := db.Order("id DESC").Find(&templates).Error; err != nil {
		response.InternalError(c, "failed to list device templates")
		return
	}

	response.Success(c, templates)
}

// Create handles POST /api/device-templates
func (h *DeviceTemplateHandler) Create(c *gin.Context) {
	var template models.DeviceTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		writeAuditLog(c, "create", "device-template", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&template).Error; err != nil {
		writeAuditLog(c, "create", "device-template", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "create", "device-template", strconv.FormatUint(uint64(template.ID), 10), "success", "device template created")
	response.Created(c, template)
}

// Update handles PUT /api/device-templates/:id
func (h *DeviceTemplateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "device-template", c.Param("id"), "failure", "invalid template id")
		response.BadRequest(c, "invalid template id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeAuditLog(c, "update", "device-template", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	allowedFields := map[string]bool{
		"name": true, "device_type": true, "parse_template_id": true,
		"field_mapping_doc_id": true, "recommended_policy": true, "enabled": true,
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
	if err := db.Model(&models.DeviceTemplate{}).Where("id = ?", id).Updates(filtered).Error; err != nil {
		writeAuditLog(c, "update", "device-template", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "update", "device-template", c.Param("id"), "success", "device template updated")
	response.SuccessWithMessage(c, "template updated", nil)
}

// Delete handles DELETE /api/device-templates/:id
func (h *DeviceTemplateHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "delete", "device-template", c.Param("id"), "failure", "invalid template id")
		response.BadRequest(c, "invalid template id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	result := db.Delete(&models.DeviceTemplate{}, id)
	if result.Error != nil {
		writeAuditLog(c, "delete", "device-template", c.Param("id"), "failure", result.Error.Error())
		response.BadRequest(c, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		writeAuditLog(c, "delete", "device-template", c.Param("id"), "failure", "device template not found")
		response.NotFound(c, "device template not found")
		return
	}

	writeAuditLog(c, "delete", "device-template", c.Param("id"), "success", "device template deleted")
	response.SuccessWithMessage(c, "template deleted", nil)
}

// Apply handles POST /api/device-templates/:id/apply
func (h *DeviceTemplateHandler) Apply(c *gin.Context) {
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

	// Get template
	var template models.DeviceTemplate
	if err := db.First(&template, id).Error; err != nil {
		writeAuditLog(c, "apply", "device-template", c.Param("id"), "failure", "template not found")
		response.NotFound(c, "template not found")
		return
	}

	// Apply to devices matching device_type
	applyUpdates := map[string]interface{}{
		"template_id":        template.ID,
		"parse_template_id":  template.ParseTemplateID,
		"field_mapping_doc_id": template.FieldMappingDocID,
	}
	result := db.Model(&models.Device{}).Where("device_type = ?", template.DeviceType).
		Updates(applyUpdates)
	if result.Error != nil {
		writeAuditLog(c, "apply", "device-template", c.Param("id"), "failure", result.Error.Error())
		response.InternalError(c, "failed to apply template")
		return
	}

	writeAuditLog(c, "apply", "device-template", c.Param("id"), "success", "device template applied")
	response.Success(c, gin.H{"affectedDevices": result.RowsAffected})
}

// RegisterDeviceTemplateRoutes registers device template routes
func RegisterDeviceTemplateRoutes(router *gin.RouterGroup, requirePerm func(string) gin.HandlerFunc) {
	handler := NewDeviceTemplateHandler()
	templates := router.Group("/device-templates")
	templates.Use(AuthRequired())
	{
		templates.GET("", requirePerm("device-templates:list"), handler.List)
		templates.GET("/all", requirePerm("device-templates:list"), handler.GetAll)
		templates.GET("/:id", requirePerm("device-templates:list"), handler.GetByID)
		templates.POST("", requirePerm("device-templates:create"), handler.Create)
		templates.PUT("/:id", requirePerm("device-templates:update"), handler.Update)
		templates.DELETE("/:id", requirePerm("device-templates:delete"), handler.Delete)
		templates.POST("/:id/apply", requirePerm("device-templates:apply"), handler.Apply)
	}
}

// FieldMappingHandler handles field mapping endpoints
type FieldMappingHandler struct{}

// NewFieldMappingHandler creates a new FieldMappingHandler
func NewFieldMappingHandler() *FieldMappingHandler {
	return &FieldMappingHandler{}
}

// List handles GET /api/field-mappings
func (h *FieldMappingHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	query := db.Model(&models.FieldMappingDoc{})
	if keyword != "" {
		query = query.Where("device_type LIKE ? OR standard_field LIKE ?", "%"+response.EscapeLike(keyword)+"%", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count field mappings")
		return
	}

	var mappings []models.FieldMappingDoc
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&mappings).Error; err != nil {
		response.InternalError(c, "failed to list field mappings")
		return
	}

	response.SuccessWithPage(c, mappings, total, page, pageSize)
}

// GetByID handles GET /api/field-mappings/:id
func (h *FieldMappingHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid mapping id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var mapping models.FieldMappingDoc
	if err := db.First(&mapping, id).Error; err != nil {
		response.NotFound(c, "field mapping not found")
		return
	}

	response.Success(c, mapping)
}

// Create handles POST /api/field-mappings
func (h *FieldMappingHandler) Create(c *gin.Context) {
	var mapping models.FieldMappingDoc
	if err := c.ShouldBindJSON(&mapping); err != nil {
		writeAuditLog(c, "create", "field-mapping", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&mapping).Error; err != nil {
		writeAuditLog(c, "create", "field-mapping", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "create", "field-mapping", strconv.FormatUint(uint64(mapping.ID), 10), "success", "field mapping created")
	response.Created(c, mapping)
}

// Update handles PUT /api/field-mappings/:id
func (h *FieldMappingHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "field-mapping", c.Param("id"), "failure", "invalid mapping id")
		response.BadRequest(c, "invalid mapping id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeAuditLog(c, "update", "field-mapping", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	allowedFields := map[string]bool{
		"device_type": true, "standard_field": true, "original_field": true,
		"field_type": true, "description": true,
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
	if err := db.Model(&models.FieldMappingDoc{}).Where("id = ?", id).Updates(filtered).Error; err != nil {
		writeAuditLog(c, "update", "field-mapping", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	writeAuditLog(c, "update", "field-mapping", c.Param("id"), "success", "field mapping updated")
	response.SuccessWithMessage(c, "field mapping updated", nil)
}

// Delete handles DELETE /api/field-mappings/:id
func (h *FieldMappingHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid mapping id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	result := db.Delete(&models.FieldMappingDoc{}, id)
	if result.Error != nil {
		writeAuditLog(c, "delete", "field-mapping", c.Param("id"), "failure", result.Error.Error())
		response.BadRequest(c, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		writeAuditLog(c, "delete", "field-mapping", c.Param("id"), "failure", "field mapping not found")
		response.NotFound(c, "field mapping not found")
		return
	}

	writeAuditLog(c, "delete", "field-mapping", c.Param("id"), "success", "field mapping deleted")
	response.SuccessWithMessage(c, "field mapping deleted", nil)
}

// RegisterFieldMappingRoutes registers field mapping routes
func RegisterFieldMappingRoutes(router *gin.RouterGroup, requirePerm func(string) gin.HandlerFunc) {
	handler := NewFieldMappingHandler()
	fm := router.Group("/field-mappings")
	fm.Use(AuthRequired())
	{
		fm.GET("", requirePerm("field-mappings:list"), handler.List)
		fm.GET("/:id", requirePerm("field-mappings:list"), handler.GetByID)
		fm.POST("", requirePerm("field-mappings:create"), handler.Create)
		fm.PUT("/:id", requirePerm("field-mappings:update"), handler.Update)
		fm.DELETE("/:id", requirePerm("field-mappings:delete"), handler.Delete)
	}
}
