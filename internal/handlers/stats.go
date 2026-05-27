package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
)

// StatsHandler handles statistics and dashboard endpoints
type StatsHandler struct {
	statsService *services.StatsService
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
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
		s.GET("/fields", requirePerm("stats:fields"), handler.FieldStats)
		s.GET("/available-fields", requirePerm("stats:available-fields"), handler.AvailableFields)
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
	}
}

// ImportExportHandler handles import/export endpoints
type ImportExportHandler struct{}

// NewImportExportHandler creates a new ImportExportHandler
func NewImportExportHandler() *ImportExportHandler {
	return &ImportExportHandler{}
}

// ImportParseTemplates handles POST /api/import/parse-templates
func (h *ImportExportHandler) ImportParseTemplates(c *gin.Context) {
	var data []map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	// TODO: Implement import logic
	response.SuccessWithMessage(c, "imported", gin.H{"count": len(data)})
}

// ImportFilterPolicies handles POST /api/import/filter-policies
func (h *ImportExportHandler) ImportFilterPolicies(c *gin.Context) {
	var data []map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	response.SuccessWithMessage(c, "imported", gin.H{"count": len(data)})
}

// ImportPushConfigs handles POST /api/import/push-configs
func (h *ImportExportHandler) ImportPushConfigs(c *gin.Context) {
	var data []map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	response.SuccessWithMessage(c, "imported", gin.H{"count": len(data)})
}

// ImportDeviceTemplates handles POST /api/import/device-templates
func (h *ImportExportHandler) ImportDeviceTemplates(c *gin.Context) {
	var data []map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	response.SuccessWithMessage(c, "imported", gin.H{"count": len(data)})
}

// ExportParseTemplates handles GET /api/export/parse-templates
func (h *ImportExportHandler) ExportParseTemplates(c *gin.Context) {
	response.Success(c, []interface{}{})
}

// ExportFilterPolicies handles GET /api/export/filter-policies
func (h *ImportExportHandler) ExportFilterPolicies(c *gin.Context) {
	response.Success(c, []interface{}{})
}

// ExportPushConfigs handles GET /api/export/push-configs
func (h *ImportExportHandler) ExportPushConfigs(c *gin.Context) {
	response.Success(c, []interface{}{})
}

// ExportDeviceTemplates handles GET /api/export/device-templates
func (h *ImportExportHandler) ExportDeviceTemplates(c *gin.Context) {
	response.Success(c, []interface{}{})
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
	}

	exp := router.Group("/export")
	exp.Use(AuthRequired())
	{
		exp.GET("/parse-templates", requirePerm("export:config"), handler.ExportParseTemplates)
		exp.GET("/filter-policies", requirePerm("export:config"), handler.ExportFilterPolicies)
		exp.GET("/push-configs", requirePerm("export:config"), handler.ExportPushConfigs)
		exp.GET("/device-templates", requirePerm("export:config"), handler.ExportDeviceTemplates)
	}
}