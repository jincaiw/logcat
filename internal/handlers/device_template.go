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

	var templates []models.DeviceTemplate
	if err := db.Find(&templates).Error; err != nil {
		response.InternalError(c, "failed to list device templates")
		return
	}

	response.Success(c, templates)
}

// Create handles POST /api/device-templates
func (h *DeviceTemplateHandler) Create(c *gin.Context) {
	var template models.DeviceTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&template).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, template)
}

// Update handles PUT /api/device-templates/:id
func (h *DeviceTemplateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	if err := db.Model(&models.DeviceTemplate{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "template updated", nil)
}

// Delete handles DELETE /api/device-templates/:id
func (h *DeviceTemplateHandler) Delete(c *gin.Context) {
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
	db.Delete(&models.DeviceTemplate{}, id)

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
		response.NotFound(c, "template not found")
		return
	}

	// Apply to devices matching device_type
	if err := db.Model(&models.Device{}).Where("device_type = ?", template.DeviceType).
		Updates(map[string]interface{}{
			"template_id":       template.ID,
			"parse_template_id": template.ParseTemplateID,
		}).Error; err != nil {
		response.InternalError(c, "failed to apply template")
		return
	}

	response.SuccessWithMessage(c, "template applied successfully", nil)
}

// RegisterDeviceTemplateRoutes registers device template routes
func RegisterDeviceTemplateRoutes(router *gin.RouterGroup, requirePerm func(string) gin.HandlerFunc) {
	handler := NewDeviceTemplateHandler()
	templates := router.Group("/device-templates")
	templates.Use(AuthRequired())
	{
		templates.GET("", requirePerm("device-templates:list"), handler.List)
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

	var mappings []models.FieldMappingDoc
	if err := db.Find(&mappings).Error; err != nil {
		response.InternalError(c, "failed to list field mappings")
		return
	}

	response.Success(c, mappings)
}

// Create handles POST /api/field-mappings
func (h *FieldMappingHandler) Create(c *gin.Context) {
	var mapping models.FieldMappingDoc
	if err := c.ShouldBindJSON(&mapping); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&mapping).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, mapping)
}

// Update handles PUT /api/field-mappings/:id
func (h *FieldMappingHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid mapping id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if err := db.Model(&models.FieldMappingDoc{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}

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
	db.Delete(&models.FieldMappingDoc{}, id)

	response.SuccessWithMessage(c, "field mapping deleted", nil)
}

// RegisterFieldMappingRoutes registers field mapping routes
func RegisterFieldMappingRoutes(router *gin.RouterGroup, requirePerm func(string) gin.HandlerFunc) {
	handler := NewFieldMappingHandler()
	fm := router.Group("/field-mappings")
	fm.Use(AuthRequired())
	{
		fm.GET("", requirePerm("field-mappings:list"), handler.List)
		fm.POST("", requirePerm("field-mappings:create"), handler.Create)
		fm.PUT("/:id", requirePerm("field-mappings:update"), handler.Update)
		fm.DELETE("/:id", requirePerm("field-mappings:delete"), handler.Delete)
	}
}