package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/pkg/response"
)

var startTime = time.Now()

// Healthz handles GET /healthz
func Healthz(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Readyz handles GET /readyz
func Readyz(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		c.JSON(503, gin.H{
			"status": "not ready",
			"reason": "database not available",
		})
		return
	}

	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.JSON(503, gin.H{
			"status": "not ready",
			"reason": "database ping failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ready",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// RuntimeMetrics handles GET /metrics/runtime
func RuntimeMetrics(c *gin.Context) {
	response.Success(c, gin.H{
		"uptime_seconds": int64(time.Since(startTime).Seconds()),
		"status":         "running",
	})
}

// Metrics handles GET /metrics
func Metrics(c *gin.Context) {
	db := database.GetDB()
	dbStatus := "connected"
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "disconnected"
		}
	} else {
		dbStatus = "not initialized"
	}

	var dbStats gin.H
	if db != nil {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			dbStats = gin.H{
				"open_connections": sqlDB.Stats().OpenConnections,
				"in_use":           sqlDB.Stats().InUse,
				"idle":             sqlDB.Stats().Idle,
				"wait_count":       sqlDB.Stats().WaitCount,
				"max_open":         sqlDB.Stats().MaxOpenConnections,
			}
		}
	}

	c.JSON(200, gin.H{
		"status":    "ok",
		"database":  dbStatus,
		"db_stats":  dbStats,
		"uptime":    int64(time.Since(startTime).Seconds()),
		"started_at": startTime.Format(time.RFC3339),
	})
}