package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
)

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
	var configs []models.PushConfig
	if err := db.Find(&configs).Error; err != nil {
		response.InternalError(c, "failed to list push configs")
		return
	}
	response.Success(c, configs)
}

// Create handles POST /api/push-configs
func (h *PushConfigHandler) Create(c *gin.Context) {
	var config models.PushConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	if err := db.Create(&config).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, config)
}

// Update handles PUT /api/push-configs/:id
func (h *PushConfigHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	db := database.GetDB()
	if err := db.Model(&models.PushConfig{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "push config updated", nil)
}

// Delete handles DELETE /api/push-configs/:id
func (h *PushConfigHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	db := database.GetDB()
	db.Delete(&models.PushConfig{}, id)
	response.SuccessWithMessage(c, "push config deleted", nil)
}

// Test handles POST /api/push-configs/:id/test
func (h *PushConfigHandler) Test(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	result, err := h.pushService.TestPush(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// RegisterPushConfigRoutes registers push config routes
func RegisterPushConfigRoutes(router *gin.RouterGroup, pushService *services.PushService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewPushConfigHandler(pushService)
	pc := router.Group("/push-configs")
	pc.Use(AuthRequired())
	{
		pc.GET("", requirePerm("push-configs:list"), handler.List)
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
	var rules []models.AlertRule
	if err := db.Preload("FilterPolicy").Preload("PushConfig").Preload("OutputTemplate").Find(&rules).Error; err != nil {
		response.InternalError(c, "failed to list alert rules")
		return
	}
	response.Success(c, rules)
}

// Create handles POST /api/alert-rules
func (h *AlertRuleHandler) Create(c *gin.Context) {
	var rule models.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}
	if err := db.Create(&rule).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, rule)
}

// Update handles PUT /api/alert-rules/:id
func (h *AlertRuleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	db := database.GetDB()
	if err := db.Model(&models.AlertRule{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "alert rule updated", nil)
}

// Delete handles DELETE /api/alert-rules/:id
func (h *AlertRuleHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	db := database.GetDB()
	db.Delete(&models.AlertRule{}, id)
	response.SuccessWithMessage(c, "alert rule deleted", nil)
}

// RegisterAlertRuleRoutes registers alert rule routes
func RegisterAlertRuleRoutes(router *gin.RouterGroup, requirePerm func(string) gin.HandlerFunc) {
	handler := NewAlertRuleHandler()
	ar := router.Group("/alert-rules")
	ar.Use(AuthRequired())
	{
		ar.GET("", requirePerm("alert-rules:list"), handler.List)
		ar.POST("", requirePerm("alert-rules:create"), handler.Create)
		ar.PUT("/:id", requirePerm("alert-rules:update"), handler.Update)
		ar.DELETE("/:id", requirePerm("alert-rules:delete"), handler.Delete)
	}
}