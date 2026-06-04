package services

import (
	"errors"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/pkg/response"
)

// UserService handles user management logic
type UserService struct{}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{}
}

// ListUsers returns paginated list of users
func (s *UserService) ListUsers(page, pageSize int, status, keyword string) ([]models.User, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not available")
	}

	var users []models.User
	var total int64

	query := db.Model(&models.User{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR display_name LIKE ? OR email LIKE ?",
			"%"+response.EscapeLike(keyword)+"%", "%"+response.EscapeLike(keyword)+"%", "%"+response.EscapeLike(keyword)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUser returns a single user by ID
func (s *UserService) GetUser(id uint) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername returns a single user by username
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	// Check for duplicate username
	var count int64
	db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return errors.New("username already exists")
	}

	return db.Create(user).Error
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	// Check if user exists
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return err
	}

	// Whitelist allowed fields to prevent arbitrary field modification
	allowedFields := map[string]bool{
		"display_name": true,
		"email":        true,
		"status":       true,
		"phone":        true,
		"department":   true,
	}

	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			filtered[k] = v
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	return db.Model(&user).Updates(filtered).Error
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uint) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	// Delete user roles first
	db.Where("user_id = ?", id).Delete(&models.UserRole{})

	return db.Delete(&models.User{}, id).Error
}

// ResetPassword resets a user's password
func (s *UserService) ResetPassword(userID uint, newPasswordHash string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_hash":        newPasswordHash,
		"must_change_password": true,
	}).Error
}

// UnlockUser unlocks a user account
func (s *UserService) UnlockUser(userID uint) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"status":             models.UserStatusEnabled,
		"failed_login_count": 0,
		"locked_until":       nil,
	}).Error
}

// AssignRoles assigns roles to a user
func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	// Remove existing roles
	if err := db.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		return err
	}

	// Add new roles
	for _, roleID := range roleIDs {
		ur := models.UserRole{UserID: userID, RoleID: roleID}
		if err := db.Create(&ur).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetUserRoles returns roles assigned to a user.
func (s *UserService) GetUserRoles(userID uint) ([]models.Role, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var user models.User
	if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user.Roles, nil
}

// ForcePasswordChange toggles the must-change-password flag for a user.
func (s *UserService) ForcePasswordChange(userID uint, required bool) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Model(&models.User{}).Where("id = ?", userID).Update("must_change_password", required).Error
}
