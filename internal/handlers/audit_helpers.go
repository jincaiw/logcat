package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/middleware"
)

func writeAuditLog(c *gin.Context, action, resourceType, resourceID, result, detail string) {
	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)
	var uid *uint
	if userID > 0 {
		uid = &userID
	}
	if err := middleware.AuditLogWriter(uid, username, action, resourceType, resourceID, c.ClientIP(), c.GetHeader("User-Agent"), result, detail); err != nil {
		log.Printf("audit log write failed: %v, action=%s resource=%s/%s", err, action, resourceType, resourceID)
	}
}
