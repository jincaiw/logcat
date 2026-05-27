package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// AuditLogWriter writes an audit log entry
func AuditLogWriter(userID *uint, username, action, resourceType, resourceID, clientIP, userAgent, result, detail string) error {
	db := database.GetDB()
	if db == nil {
		return nil
	}

	audit := models.AuditLog{
		UserID:       userID,
		Username:     username,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ClientIP:     clientIP,
		UserAgent:    userAgent,
		Result:       result,
		Detail:       detail,
	}

	return db.Create(&audit).Error
}

// AuditAction is a middleware that writes audit logs for API actions
func AuditAction(action, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only log successful operations
		status := c.Writer.Status()
		if status >= 200 && status < 300 {
			userID := GetUserID(c)
			username := GetUsername(c)
			resourceID := c.Param("id")

			var uid *uint
			if userID > 0 {
				uid = &userID
			}

			_ = AuditLogWriter(
				uid,
				username,
				action,
				resourceType,
				resourceID,
				c.ClientIP(),
				c.GetHeader("User-Agent"),
				"success",
				c.Request.Method+" "+c.Request.URL.Path,
			)
		}
	}
}