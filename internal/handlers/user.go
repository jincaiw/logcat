package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles user management endpoints
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUserRequest is the create user request body
type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	RoleIDs     []uint `json:"roleIds"`
}

// UpdateUserRequest is the update user request body
type UpdateUserRequest struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	RoleIDs     []uint `json:"roleIds"`
	RoleIDsSet  bool   `json:"-"` // tracks whether roleIds was explicitly provided
}

// List handles GET /api/users
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	page, pageSize = response.NormalizePagination(page, pageSize)
	status := c.Query("status")
	keyword := c.Query("keyword")

	users, total, err := h.userService.ListUsers(page, pageSize, status, keyword)
	if err != nil {
		response.InternalError(c, "failed to list users")
		return
	}

	response.SuccessWithPage(c, users, total, page, pageSize)
}

// GetByID handles GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	db := database.GetDB()
	if db == nil {
		response.InternalError(c, "database not available")
		return
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, user)
}

// Create handles POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAuditLog(c, "create", "user", "", "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := services.ValidatePasswordStrength(req.Password); err != nil {
		writeAuditLog(c, "create", "user", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeAuditLog(c, "create", "user", "", "failure", "failed to hash password")
		response.InternalError(c, "failed to hash password")
		return
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		Status:       models.UserStatusEnabled,
	}

	if err := h.userService.CreateUser(user); err != nil {
		writeAuditLog(c, "create", "user", "", "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	// Assign roles if provided
	if len(req.RoleIDs) > 0 {
		if err := h.userService.AssignRoles(user.ID, req.RoleIDs); err != nil {
			writeAuditLog(c, "assign_roles", "user", strconv.FormatUint(uint64(user.ID), 10), "failure", err.Error())
			response.InternalError(c, "failed to assign roles")
			return
		}
	}

	writeAuditLog(c, "create", "user", strconv.FormatUint(uint64(user.ID), 10), "success", "user created")
	response.Created(c, user)
}

// Update handles PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "update", "user", c.Param("id"), "failure", "invalid user id")
		response.BadRequest(c, "invalid user id")
		return
	}

	var req UpdateUserRequest
	rawBody := make(map[string]interface{})
	if err := c.ShouldBindJSON(&rawBody); err != nil {
		writeAuditLog(c, "update", "user", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	// Map raw fields to request
	if v, ok := rawBody["displayName"].(string); ok {
		req.DisplayName = v
	}
	if v, ok := rawBody["email"].(string); ok {
		req.Email = v
	}
	if v, ok := rawBody["status"].(string); ok {
		req.Status = v
	}
	if _, ok := rawBody["roleIds"]; ok {
		req.RoleIDsSet = true
		if arr, ok := rawBody["roleIds"].([]interface{}); ok {
			for _, item := range arr {
				if f, ok := item.(float64); ok {
					req.RoleIDs = append(req.RoleIDs, uint(f))
				}
			}
		}
	}

	updates := make(map[string]interface{})
	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := h.userService.UpdateUser(uint(id), updates); err != nil {
			writeAuditLog(c, "update", "user", c.Param("id"), "failure", err.Error())
			response.BadRequest(c, err.Error())
			return
		}
	}

	// Update roles if explicitly provided (including empty array to clear roles)
	if req.RoleIDsSet {
		if err := h.userService.AssignRoles(uint(id), req.RoleIDs); err != nil {
			writeAuditLog(c, "assign_roles", "user", c.Param("id"), "failure", err.Error())
			response.InternalError(c, "failed to assign roles")
			return
		}
	}

	writeAuditLog(c, "update", "user", c.Param("id"), "success", "user updated")
	response.SuccessWithMessage(c, "user updated", nil)
}

// Delete handles DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "delete", "user", c.Param("id"), "failure", "invalid user id")
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		writeAuditLog(c, "delete", "user", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to delete user")
		return
	}

	writeAuditLog(c, "delete", "user", c.Param("id"), "success", "user deleted")
	response.SuccessWithMessage(c, "user deleted", nil)
}

// ResetPassword handles POST /api/users/:id/reset-password
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "reset_password", "user", c.Param("id"), "failure", "invalid user id")
		response.BadRequest(c, "invalid user id")
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAuditLog(c, "reset_password", "user", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := services.ValidatePasswordStrength(req.Password); err != nil {
		writeAuditLog(c, "reset_password", "user", c.Param("id"), "failure", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeAuditLog(c, "reset_password", "user", c.Param("id"), "failure", "failed to hash password")
		response.InternalError(c, "failed to hash password")
		return
	}

	if err := h.userService.ResetPassword(uint(id), string(hashedPassword)); err != nil {
		writeAuditLog(c, "reset_password", "user", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to reset password")
		return
	}

	writeAuditLog(c, "reset_password", "user", c.Param("id"), "success", "password reset successful")
	response.SuccessWithMessage(c, "password reset successful", nil)
}

// Unlock handles POST /api/users/:id/unlock
func (h *UserHandler) Unlock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "unlock", "user", c.Param("id"), "failure", "invalid user id")
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userService.UnlockUser(uint(id)); err != nil {
		writeAuditLog(c, "unlock", "user", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to unlock user")
		return
	}

	writeAuditLog(c, "unlock", "user", c.Param("id"), "success", "user unlocked")
	response.SuccessWithMessage(c, "user unlocked", nil)
}

// GetRoles handles GET /api/users/:id/roles
func (h *UserHandler) GetRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	roles, err := h.userService.GetUserRoles(uint(id))
	if err != nil {
		response.InternalError(c, "failed to get user roles")
		return
	}

	response.Success(c, gin.H{"roles": roles})
}

// AssignRolesRequest is the request body for assigning roles.
type AssignRolesRequest struct {
	RoleIDs []uint `json:"roleIds" binding:"required"`
}

// AssignRoles handles POST /api/users/:id/roles
func (h *UserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "assign_roles", "user", c.Param("id"), "failure", "invalid user id")
		response.BadRequest(c, "invalid user id")
		return
	}

	var req AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAuditLog(c, "assign_roles", "user", c.Param("id"), "failure", "invalid request: "+err.Error())
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.userService.AssignRoles(uint(id), req.RoleIDs); err != nil {
		writeAuditLog(c, "assign_roles", "user", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to assign roles")
		return
	}

	writeAuditLog(c, "assign_roles", "user", c.Param("id"), "success", "roles assigned")
	response.SuccessWithMessage(c, "roles assigned", nil)
}

// ForcePasswordChange handles POST /api/users/:id/force-password-change
func (h *UserHandler) ForcePasswordChange(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeAuditLog(c, "force_password_change", "user", c.Param("id"), "failure", "invalid user id")
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userService.ForcePasswordChange(uint(id), true); err != nil {
		writeAuditLog(c, "force_password_change", "user", c.Param("id"), "failure", err.Error())
		response.InternalError(c, "failed to update password change flag")
		return
	}

	writeAuditLog(c, "force_password_change", "user", c.Param("id"), "success", "password change required")
	response.SuccessWithMessage(c, "password change required", nil)
}

// RegisterUserRoutes registers user routes
func RegisterUserRoutes(router *gin.RouterGroup, userService *services.UserService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewUserHandler(userService)
	users := router.Group("/users")
	users.Use(AuthRequired())
	{
		users.GET("", requirePerm("users:list"), handler.List)
		users.GET("/:id", requirePerm("users:list"), handler.GetByID)
		users.POST("", requirePerm("users:create"), handler.Create)
		users.PUT("/:id", requirePerm("users:update"), handler.Update)
		users.DELETE("/:id", requirePerm("users:delete"), handler.Delete)
		users.POST("/:id/reset-password", requirePerm("users:reset-password"), handler.ResetPassword)
		users.POST("/:id/unlock", requirePerm("users:unlock"), handler.Unlock)
		users.GET("/:id/roles", requirePerm("users:list"), handler.GetRoles)
		users.POST("/:id/roles", requirePerm("users:update"), handler.AssignRoles)
		users.POST("/:id/force-password-change", requirePerm("users:update"), handler.ForcePasswordChange)
	}
}
