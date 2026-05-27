package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// RequirePermission checks if the current user has the specified permission code
func RequirePermission(permCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":      401,
				"message":   "unauthorized",
				"requestId": GetRequestID(c),
			})
			return
		}

		db := database.GetDB()
		if db == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":      500,
				"message":   "database not available",
				"requestId": GetRequestID(c),
			})
			return
		}

		// Check if user has the required permission through any of their roles
		var count int64
		err := db.Table("user_roles").
			Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Where("user_roles.user_id = ? AND permissions.code = ?", userID, permCode).
			Count(&count).Error

		if err != nil || count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":      403,
				"message":   "forbidden: insufficient permissions",
				"requestId": GetRequestID(c),
			})
			return
		}

		c.Next()
	}
}

// GetUserRoles returns all role codes for the current user
func GetUserRoles(userID uint) ([]string, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}

	var roles []string
	err := db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("roles.code", &roles).Error

	return roles, err
}

// GetUserPermissions returns all permission codes for the current user
func GetUserPermissions(userID uint) ([]string, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}

	var perms []string
	err := db.Table("user_roles").
		Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("user_roles.user_id = ?", userID).
		Distinct("permissions.code").
		Pluck("permissions.code", &perms).Error

	return perms, err
}

// GetAllPermModels returns all permissions as models
func GetAllPermModels() ([]models.Permission, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}

	var perms []models.Permission
	if err := db.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}