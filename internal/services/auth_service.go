package services

import (
	"errors"
	"regexp"
	"time"

	"github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/middleware"
	"github.com/logcat/logcat/internal/models"

	"golang.org/x/crypto/bcrypt"
)

var passwordMinLength = 8

func ValidatePasswordStrength(password string) error {
	if len(password) < passwordMinLength {
		return errors.New("password must be at least 8 characters long")
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("password must contain uppercase, lowercase letters and digits")
	}
	return nil
}

// AuthService handles authentication logic
type AuthService struct{}

// NewAuthService creates a new AuthService
func NewAuthService() *AuthService {
	return &AuthService{}
}

// LoginResult holds the result of a login attempt
type LoginResult struct {
	Success  bool
	Message  string
	Token    string
	User     *UserResponse
}

// UserResponse is the public user data returned to clients
type UserResponse struct {
	ID                 uint              `json:"id"`
	Username           string            `json:"username"`
	DisplayName        string            `json:"displayName"`
	Email              string            `json:"email"`
	Status             string            `json:"status"`
	LastLoginAt        *time.Time        `json:"lastLoginAt"`
	MustChangePassword bool              `json:"mustChangePassword"`
	Roles              []string          `json:"roles"`
	Permissions        []string          `json:"permissions"`
	CreatedAt          time.Time         `json:"createdAt"`
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(username, password string, expireHours int) *LoginResult {
	db := database.GetDB()
	if db == nil {
		return &LoginResult{Success: false, Message: "database not available"}
	}

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return &LoginResult{Success: false, Message: "invalid username or password"}
	}

	// Check account status
	if user.Status == models.UserStatusDisabled {
		return &LoginResult{Success: false, Message: "account is disabled"}
	}

	// Check if account is locked
	if user.Status == models.UserStatusLocked {
		if user.LockedUntil != nil {
			if time.Now().Before(*user.LockedUntil) {
				return &LoginResult{Success: false, Message: "account is locked, please try again later"}
			}
			// Unlock if lock duration has passed
			user.Status = models.UserStatusEnabled
			user.FailedLoginCount = 0
			user.LockedUntil = nil
			db.Model(&user).Updates(map[string]interface{}{
				"status":             models.UserStatusEnabled,
				"failed_login_count": 0,
				"locked_until":       nil,
			})
		} else {
			// LockedUntil is nil but status is locked - treat as locked
			return &LoginResult{Success: false, Message: "account is locked, please try again later"}
		}
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		user.FailedLoginCount++
		cfg := config.Get()
		maxFailed := 5
		lockDuration := 30
		if cfg != nil {
			maxFailed = cfg.Auth.MaxFailedLogin
			lockDuration = cfg.Auth.LockDurationMinutes
		}

		if user.FailedLoginCount >= maxFailed {
			user.Status = models.UserStatusLocked
			lockUntil := time.Now().Add(time.Duration(lockDuration) * time.Minute)
			user.LockedUntil = &lockUntil
		}
		db.Model(&user).Updates(map[string]interface{}{
			"failed_login_count": user.FailedLoginCount,
			"status":             user.Status,
			"locked_until":       user.LockedUntil,
		})
		return &LoginResult{Success: false, Message: "invalid username or password"}
	}

	// Reset failed count on success
	user.FailedLoginCount = 0
	user.LockedUntil = nil
	now := time.Now()
	user.LastLoginAt = &now
	db.Model(&user).Updates(map[string]interface{}{
		"failed_login_count": 0,
		"locked_until":       nil,
		"last_login_at":      &now,
	})

	// Create session
	token := middleware.DefaultSessionStore.Create(user.ID, user.Username, expireHours)

	// Build user response
	userResp := buildUserResponse(&user)

	return &LoginResult{
		Success: true,
		Message: "login successful",
		Token:   token,
		User:    userResp,
	}
}

// Logout removes the session
func (s *AuthService) Logout(token string) {
	middleware.DefaultSessionStore.Delete(token)
}

// GetCurrentUser returns the current user's information
func (s *AuthService) GetCurrentUser(userID uint) (*UserResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var user models.User
	if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return buildUserResponse(&user), nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	if err := ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.MustChangePassword = false
	return db.Save(&user).Error
}

// InitAdmin initializes the admin user if no users exist
func (s *AuthService) InitAdmin(username, password string) (*UserResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	// Check if any user already exists
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return nil, errors.New("system already initialized")
	}

	if err := ValidatePasswordStrength(password); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := models.User{
		Username:           username,
		PasswordHash:       string(hashedPassword),
		DisplayName:        "系统管理员",
		Status:             models.UserStatusEnabled,
		MustChangePassword: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return nil, err
	}

	// Assign admin role
	var adminRole models.Role
	if err := db.Where("code = ?", "admin").First(&adminRole).Error; err == nil {
		ur := models.UserRole{UserID: admin.ID, RoleID: adminRole.ID}
		db.Create(&ur)
	}

	return buildUserResponse(&admin), nil
}

// IsInitialized checks if the system has been initialized (any user exists)
func (s *AuthService) IsInitialized() bool {
	db := database.GetDB()
	if db == nil {
		return false
	}
	var count int64
	db.Model(&models.User{}).Count(&count)
	return count > 0
}

// buildUserResponse builds a UserResponse from a User model
func buildUserResponse(user *models.User) *UserResponse {
	resp := &UserResponse{
		ID:                 user.ID,
		Username:           user.Username,
		DisplayName:        user.DisplayName,
		Email:              user.Email,
		Status:             string(user.Status),
		LastLoginAt:        user.LastLoginAt,
		MustChangePassword: user.MustChangePassword,
		CreatedAt:          user.CreatedAt,
	}

	// Get roles
	roles, _ := middleware.GetUserRoles(user.ID)
	resp.Roles = roles

	// Get permissions
	perms, _ := middleware.GetUserPermissions(user.ID)
	resp.Permissions = perms

	return resp
}