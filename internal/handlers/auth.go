package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/middleware"
	"github.com/logcat/logcat/internal/services"
	"github.com/logcat/logcat/pkg/response"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// LoginRequest is the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest is the change password request body
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	expireHours := 24
	if cfg := config.Get(); cfg != nil {
		expireHours = cfg.Auth.SessionExpireHours
	}

	result := h.authService.Login(req.Username, req.Password, expireHours)
	if !result.Success {
		middleware.AuditLogWriter(
			nil,
			req.Username,
			"login",
			"auth",
			"",
			c.ClientIP(),
			c.GetHeader("User-Agent"),
			"failure",
			result.Message,
		)
		response.Error(c, http.StatusUnauthorized, http.StatusUnauthorized, result.Message)
		return
	}

	setSessionCookie(c, result.Token, expireHours*3600)

	// Log audit
	middleware.AuditLogWriter(
		&result.User.ID,
		result.User.Username,
		"login",
		"auth",
		"",
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		"success",
		"user login",
	)

	response.Success(c, gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	token, _ := c.Cookie("session_token")
	h.authService.Logout(token)
	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	var uid *uint
	if userID > 0 {
		uid = &userID
	}

	clearSessionCookie(c)

	middleware.AuditLogWriter(
		uid,
		username,
		"logout",
		"auth",
		"",
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		"success",
		"user logout",
	)

	response.SuccessWithMessage(c, "logout successful", nil)
}

// Me handles GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		response.InternalError(c, "failed to get user info")
		return
	}

	response.Success(c, user)
}

// ChangePassword handles POST /api/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		username := middleware.GetUsername(c)
		middleware.AuditLogWriter(
			&userID,
			username,
			"change_password",
			"auth",
			"",
			c.ClientIP(),
			c.GetHeader("User-Agent"),
			"failure",
			err.Error(),
		)
		response.BadRequest(c, err.Error())
		return
	}

	username := middleware.GetUsername(c)
	middleware.AuditLogWriter(
		&userID,
		username,
		"change_password",
		"auth",
		"",
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		"success",
		"password changed",
	)

	response.SuccessWithMessage(c, "password changed successfully", nil)
}

// InitAdmin handles POST /api/auth/init-admin
func (h *AuthHandler) InitAdmin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	user, err := h.authService.InitAdmin(req.Username, req.Password)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "admin initialized successfully", user)
}

// InitStatus handles GET /api/auth/init-status
func (h *AuthHandler) InitStatus(c *gin.Context) {
	initialized := h.authService.IsInitialized()
	response.Success(c, gin.H{"initialized": initialized})
}

// Refresh handles POST /api/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	username := middleware.GetUsername(c)

	oldToken, _ := c.Cookie("session_token")
	if oldToken != "" {
		middleware.DefaultSessionStore.Delete(oldToken)
	}

	expireHours := 24
	if cfg := config.Get(); cfg != nil {
		expireHours = cfg.Auth.SessionExpireHours
	}

	token := middleware.DefaultSessionStore.Create(userID, username, expireHours)
	setSessionCookie(c, token, expireHours*3600)

	response.Success(c, gin.H{
		"token": token,
	})
}

func setSessionCookie(c *gin.Context, token string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", token, maxAge, "/", "", requestUsesHTTPS(c), true)
}

func clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", "", -1, "/", "", requestUsesHTTPS(c), true)
}

func requestUsesHTTPS(c *gin.Context) bool {
	if c.Request != nil && c.Request.TLS != nil {
		return true
	}
	proto := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Forwarded-Proto")))
	if proto == "https" {
		return true
	}
	forwarded := strings.ToLower(c.GetHeader("Forwarded"))
	return strings.Contains(forwarded, "proto=https")
}

// RegisterRoutes registers auth routes
func RegisterRoutes(router *gin.RouterGroup, authService *services.AuthService, middlewareFn func(string) gin.HandlerFunc) {
	handler := NewAuthHandler(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/login", middleware.LoginRateLimit(), handler.Login)
		auth.POST("/init-admin", middleware.InitRateLimit(), handler.InitAdmin)
		auth.GET("/init-status", handler.InitStatus)

		// Protected routes
		protected := auth.Group("")
		protected.Use(middleware.AuthRequired())
		{
			protected.POST("/logout", handler.Logout)
			protected.GET("/me", handler.Me)
			protected.POST("/change-password", handler.ChangePassword)
			protected.POST("/refresh", handler.Refresh)
		}
	}
}

// Ensure AuthRequired is available from middleware package
var AuthRequired = middleware.AuthRequired
