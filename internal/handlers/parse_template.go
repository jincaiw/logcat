package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
)

// ParseTemplateHandler handles parse template endpoints
type ParseTemplateHandler struct {
	parseService *services.ParseService
}

// NewParseTemplateHandler creates a new ParseTemplateHandler
func NewParseTemplateHandler(parseService *services.ParseService) *ParseTemplateHandler {
	return &ParseTemplateHandler{parseService: parseService}
}

// List handles GET /api/parse-templates
func (h *ParseTemplateHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var templates []models.ParseTemplate
	if err := db.Find(&templates).Error; err != nil {
		response.InternalError(c, "failed to list parse templates")
		return
	}

	response.Success(c, templates)
}

// Create handles POST /api/parse-templates
func (h *ParseTemplateHandler) Create(c *gin.Context) {
	var template models.ParseTemplate
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

// Update handles PUT /api/parse-templates/:id
func (h *ParseTemplateHandler) Update(c *gin.Context) {
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
	if err := db.Model(&models.ParseTemplate{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}

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
	db.Delete(&models.ParseTemplate{}, id)

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
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	result, err := h.parseService.TestParse(req.TemplateID, req.RawLog)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// RegisterParseTemplateRoutes registers parse template routes
func RegisterParseTemplateRoutes(router *gin.RouterGroup, parseService *services.ParseService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewParseTemplateHandler(parseService)
	pt := router.Group("/parse-templates")
	pt.Use(AuthRequired())
	{
		pt.GET("", requirePerm("parse-templates:list"), handler.List)
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

	var policies []models.FilterPolicy
	if err := db.Find(&policies).Error; err != nil {
		response.InternalError(c, "failed to list filter policies")
		return
	}

	response.Success(c, policies)
}

// Create handles POST /api/filter-policies
func (h *FilterPolicyHandler) Create(c *gin.Context) {
	var policy models.FilterPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&policy).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, policy)
}

// Update handles PUT /api/filter-policies/:id
func (h *FilterPolicyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid policy id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if err := db.Model(&models.FilterPolicy{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}

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
	db.Delete(&models.FilterPolicy{}, id)

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
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	result, err := h.filterService.TestFilter(req.PolicyID, req.ParsedData)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// RegisterFilterPolicyRoutes registers filter policy routes
func RegisterFilterPolicyRoutes(router *gin.RouterGroup, filterService *services.FilterService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewFilterPolicyHandler(filterService)
	fp := router.Group("/filter-policies")
	fp.Use(AuthRequired())
	{
		fp.GET("", requirePerm("filter-policies:list"), handler.List)
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

	var templates []models.OutputTemplate
	if err := db.Find(&templates).Error; err != nil {
		response.InternalError(c, "failed to list output templates")
		return
	}

	response.Success(c, templates)
}

// Create handles POST /api/output-templates
func (h *OutputTemplateHandler) Create(c *gin.Context) {
	var template models.OutputTemplate
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

// Update handles PUT /api/output-templates/:id
func (h *OutputTemplateHandler) Update(c *gin.Context) {
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
	if err := db.Model(&models.OutputTemplate{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}

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
	db.Delete(&models.OutputTemplate{}, id)

	response.SuccessWithMessage(c, "output template deleted", nil)
}

// RegisterOutputTemplateRoutes registers output template routes
func RegisterOutputTemplateRoutes(router *gin.RouterGroup, requirePerm func(string) gin.HandlerFunc) {
	handler := NewOutputTemplateHandler()
	ot := router.Group("/output-templates")
	ot.Use(AuthRequired())
	{
		ot.GET("", requirePerm("output-templates:list"), handler.List)
		ot.POST("", requirePerm("output-templates:create"), handler.Create)
		ot.PUT("/:id", requirePerm("output-templates:update"), handler.Update)
		ot.DELETE("/:id", requirePerm("output-templates:delete"), handler.Delete)
	}
}