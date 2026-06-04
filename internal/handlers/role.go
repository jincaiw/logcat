package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	keyword := c.Query("keyword")

	query := db.Model(&models.Role{}).Preload("Permissions")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.InternalError(c, "failed to count roles")
		return
	}

	var roles []models.Role
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&roles).Error; err != nil {
		response.InternalError(c, "failed to list roles")
		return
	}

	response.SuccessWithPage(c, roles, total, page, pageSize)
}

// GetByID handles GET /api/roles/:id
func (h *RoleHandler) GetByID(c *gin.Context) {
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

	response.Success(c, role)
}

// Create handles POST /api/roles
func (h *RoleHandler) Create(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		writeAuditLog(c, "create", "role", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Create(&role).Error; err != nil {
		writeAuditLog(c, "create", "role", "", "failure", err.Error())
		response.BadRequest(c, "failed to create role: "+err.Error())
		return
	}

	writeAuditLog(c, "create", "role", strconv.FormatUint(uint64(role.ID), 10), "success", "role created")
	response.Created(c, role)
}

// Update handles PUT /api/roles/:id
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "role", c.Param("id"), "failure", "invalid role id")
		response.BadRequest(c, "invalid role id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		writeAuditLog(c, "update", "role", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	allowedFields := map[string]bool{"name": true, "description": true}
	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			filtered[k] = v
		}
	}
	if len(filtered) == 0 {
		writeAuditLog(c, "update", "role", c.Param("id"), "failure", "no updatable fields")
		response.BadRequest(c, "no updatable fields provided")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	if err := db.Model(&models.Role{}).Where("id = ?", id).Updates(filtered).Error; err != nil {
		writeAuditLog(c, "update", "role", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, "failed to update role: "+err.Error())
		return
	}

	writeAuditLog(c, "update", "role", c.Param("id"), "success", "role updated")
	response.SuccessWithMessage(c, "role updated", nil)
}

// Delete handles DELETE /api/roles/:id
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "delete", "role", c.Param("id"), "failure", "invalid role id")
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
		writeAuditLog(c, "delete", "role", c.Param("id"), "failure", "role not found")
		response.NotFound(c, "role not found")
		return
	}
	if role.BuiltIn {
		writeAuditLog(c, "delete", "role", c.Param("id"), "failure", "cannot delete built-in role")
		response.BadRequest(c, "cannot delete built-in role")
		return
	}

	// Delete role permissions and role in a transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Role{}, id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		writeAuditLog(c, "delete", "role", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, "failed to delete role: "+err.Error())
		return
	}

	writeAuditLog(c, "delete", "role", c.Param("id"), "success", "role deleted")
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
		writeAuditLog(c, "set_permissions", "role", c.Param("id"), "failure", "invalid role id")
		response.BadRequest(c, "invalid role id")
		return
	}

	var req SetPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAuditLog(c, "set_permissions", "role", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	// Remove existing permissions and add new permissions in a transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		for _, permID := range req.PermissionIDs {
			rp := models.RolePermission{RoleID: uint(id), PermissionID: permID}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		writeAuditLog(c, "set_permissions", "role", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, "failed to set permissions: "+err.Error())
		return
	}

	writeAuditLog(c, "set_permissions", "role", c.Param("id"), "success", "role permissions updated")
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
		roles.GET("/:id", requirePerm("roles:list"), handler.GetByID)
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
