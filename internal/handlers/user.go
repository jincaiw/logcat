package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
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
	DisplayName string          `json:"displayName"`
	Email       string          `json:"email"`
	Status      string          `json:"status"`
	RoleIDs     []uint          `json:"roleIds"`
}

// List handles GET /api/users
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	users, total, err := h.userService.ListUsers(page, pageSize, status, keyword)
	if err != nil {
		response.InternalError(c, "failed to list users")
		return
	}

	response.SuccessWithPage(c, users, total, page, pageSize)
}

// Create handles POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
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
		response.BadRequest(c, err.Error())
		return
	}

	// Assign roles if provided
	if len(req.RoleIDs) > 0 {
		h.userService.AssignRoles(user.ID, req.RoleIDs)
	}

	response.Created(c, user)
}

// Update handles PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
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
			response.BadRequest(c, err.Error())
			return
		}
	}

	// Update roles if provided
	if len(req.RoleIDs) > 0 {
		if err := h.userService.AssignRoles(uint(id), req.RoleIDs); err != nil {
			response.InternalError(c, "failed to assign roles")
			return
		}
	}

	response.SuccessWithMessage(c, "user updated", nil)
}

// Delete handles DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		response.InternalError(c, "failed to delete user")
		return
	}

	response.SuccessWithMessage(c, "user deleted", nil)
}

// ResetPassword handles POST /api/users/:id/reset-password
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "failed to hash password")
		return
	}

	if err := h.userService.ResetPassword(uint(id), string(hashedPassword)); err != nil {
		response.InternalError(c, "failed to reset password")
		return
	}

	response.SuccessWithMessage(c, "password reset successful", nil)
}

// Unlock handles POST /api/users/:id/unlock
func (h *UserHandler) Unlock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userService.UnlockUser(uint(id)); err != nil {
		response.InternalError(c, "failed to unlock user")
		return
	}

	response.SuccessWithMessage(c, "user unlocked", nil)
}

// RegisterUserRoutes registers user routes
func RegisterUserRoutes(router *gin.RouterGroup, userService *services.UserService, requirePerm func(string) gin.HandlerFunc) {
	handler := NewUserHandler(userService)
	users := router.Group("/users")
	users.Use(AuthRequired())
	{
		users.GET("", requirePerm("users:list"), handler.List)
		users.POST("", requirePerm("users:create"), handler.Create)
		users.PUT("/:id", requirePerm("users:update"), handler.Update)
		users.DELETE("/:id", requirePerm("users:delete"), handler.Delete)
		users.POST("/:id/reset-password", requirePerm("users:reset-password"), handler.ResetPassword)
		users.POST("/:id/unlock", requirePerm("users:unlock"), handler.Unlock)
	}
}