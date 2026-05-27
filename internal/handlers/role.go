package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
)

// RoleHandler handles role management endpoints
type RoleHandler struct {
	userService *services.UserService
}

// NewRoleHandler creates a new RoleHandler
func NewRoleHandler(userService *services.UserService) *RoleHandler {
	return &RoleHandler{userService: userService}
}

// List handles GET /api/roles
func (h *RoleHandler) List(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var roles []models.Role
	if err := db.Preload("Permissions").Find(&roles).Error; err != nil {
		response.InternalError(c, "failed to list roles")
		return
	}

	response.Success(c, roles)
}

// Create handles POST /api/roles
func (h *RoleHandler) Create(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&role).Error; err != nil {
		response.BadRequest(c, "failed to create role: "+err.Error())
		return
	}

	response.Created(c, role)
}

// Update handles PUT /api/roles/:id
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid role id")
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

	if err := db.Model(&models.Role{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.BadRequest(c, "failed to update role: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "role updated", nil)
}

// Delete handles DELETE /api/roles/:id
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid role id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	// Don't allow deleting built-in roles
	var role models.Role
	if err := db.First(&role, id).Error; err != nil {
		response.NotFound(c, "role not found")
		return
	}
	if role.BuiltIn {
		response.BadRequest(c, "cannot delete built-in role")
		return
	}

	// Delete role permissions first
	db.Where("role_id = ?", id).Delete(&models.RolePermission{})
	db.Delete(&models.Role{}, id)

	response.SuccessWithMessage(c, "role deleted", nil)
}

// GetPermissions handles GET /api/roles/:id/permissions
func (h *RoleHandler) GetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid role id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var role models.Role
	if err := db.Preload("Permissions").First(&role, id).Error; err != nil {
		response.NotFound(c, "role not found")
		return
	}

	response.Success(c, role.Permissions)
}

// SetPermissionsRequest is the request body for setting permissions
type SetPermissionsRequest struct {
	PermissionIDs []uint `json:"permissionIds" binding:"required"`
}

// SetPermissions handles POST /api/roles/:id/permissions
func (h *RoleHandler) SetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid role id")
		return
	}

	var req SetPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	// Remove existing permissions
	db.Where("role_id = ?", id).Delete(&models.RolePermission{})

	// Add new permissions
	for _, permID := range req.PermissionIDs {
		rp := models.RolePermission{RoleID: uint(id), PermissionID: permID}
		db.Create(&rp)
	}

	response.SuccessWithMessage(c, "permissions updated", nil)
}

// ListPermissions handles GET /api/permissions
func (h *RoleHandler) ListPermissions(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var perms []models.Permission
	if err := db.Find(&perms).Error; err != nil {
		response.InternalError(c, "failed to list permissions")
		return
	}

	response.Success(c, perms)
}

// RegisterRoleRoutes registers role and permission routes
func RegisterRoleRoutes(router *gin.RouterGroup, userService *services.UserService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewRoleHandler(userService)

	roles := router.Group("/roles")
	roles.Use(AuthRequired())
	{
		roles.GET("", requirePerm("roles:list"), handler.List)
		roles.POST("", requirePerm("roles:create"), handler.Create)
		roles.PUT("/:id", requirePerm("roles:update"), handler.Update)
		roles.DELETE("/:id", requirePerm("roles:delete"), handler.Delete)
		roles.GET("/:id/permissions", requirePerm("roles:permissions"), handler.GetPermissions)
		roles.POST("/:id/permissions", requirePerm("roles:permissions"), handler.SetPermissions)
	}

	// Permissions
	perms := router.Group("/permissions")
	perms.Use(AuthRequired())
	{
		perms.GET("", requirePerm("permissions:list"), handler.ListPermissions)
	}
}