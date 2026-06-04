package handlers

import (
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/engine"
	logsyslog "github.com/logcat/logcat/internal/syslog"
	"github.com/logcat/logcat/pkg/response"
)

var startTime = time.Now()

// Healthz handles GET /healthz and /api/healthz
func Healthz(c *gin.Context) {
	payload := gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}
	if strings.HasPrefix(c.FullPath(), "/api/") {
		response.Success(c, payload)
		return
	}
	c.JSON(200, payload)
}

// Readyz handles GET /readyz and /api/readyz
func Readyz(c *gin.Context) {
	db := database.GetDB()
	notReady := func(reason string) {
		payload := gin.H{
			"status": "not ready",
			"reason": reason,
		}
		if strings.HasPrefix(c.FullPath(), "/api/") {
			c.JSON(503, gin.H{
				"code":    503,
				"message": reason,
				"data":    payload,
			})
			return
		}
		c.JSON(503, payload)
	}

	if db == nil {
		notReady("database not available")
		return
	}

	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		notReady("database ping failed")
		return
	}

	payload := gin.H{
		"status": "ready",
		"time":   time.Now().Format(time.RFC3339),
	}
	if strings.HasPrefix(c.FullPath(), "/api/") {
		response.Success(c, payload)
		return
	}
	c.JSON(200, payload)
}

// RuntimeMetrics handles GET /metrics/runtime
func RuntimeMetrics(c *gin.Context) {
	response.Success(c, collectRuntimeMetrics())
}

// Metrics handles GET /metrics and /api/metrics
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

	receiverMetrics := gin.H{}
	if receiver := logsyslog.GetGlobalReceiver(); receiver != nil {
		metrics := receiver.Metrics()
		receiverMetrics = gin.H{
			"udp_received":     metrics.UDPReceived,
			"tcp_received":     metrics.TCPReceived,
			"udp_errors":       metrics.UDPErrors,
			"tcp_errors":       metrics.TCPErrors,
			"parse_errors":     metrics.ParseErrors,
			"channel_dropped":  metrics.ChannelDropped,
			"tcp_connections":  metrics.TCPConnections,
			"last_receive_at":  metrics.LastReceiveAt,
		}
	}

	pipelineMetrics := gin.H{}
	if pipeline := engine.GetGlobalPipeline(); pipeline != nil {
		metrics := pipeline.Metrics()
		pipelineMetrics = gin.H{
			"raw_queue_depth":    metrics.RawQueueDepth,
			"parsed_queue_depth": metrics.ParsedQueueDepth,
			"db_queue_depth":     metrics.DBQueueDepth,
			"push_queue_depth":   metrics.PushQueueDepth,
			"parse_processed":    metrics.ParseProcessed,
			"parse_errors":       metrics.ParseErrors,
			"filter_processed":   metrics.FilterProcessed,
			"filter_dropped":     metrics.FilterDropped,
			"db_written":         metrics.DBWritten,
			"db_errors":          metrics.DBErrors,
			"push_processed":     metrics.PushProcessed,
			"push_errors":        metrics.PushErrors,
			"raw_dropped":        metrics.RawDropped,
			"db_dropped":         metrics.DBDropped,
			"push_dropped":       metrics.PushDropped,
		}
	}

	payload := gin.H{
		"status":           "ok",
		"database":         dbStatus,
		"db_stats":         dbStats,
		"uptime":           int64(time.Since(startTime).Seconds()),
		"started_at":       startTime.Format(time.RFC3339),
		"receiver_metrics": receiverMetrics,
		"pipeline_metrics": pipelineMetrics,
	}

	// For /metrics (unauthenticated, e.g. Prometheus), return raw JSON.
	// For /api/metrics (authenticated, used by the UI), wrap in the standard
	// {code,message,data} envelope so the front-end interceptor works.
	if strings.HasPrefix(c.FullPath(), "/api/") {
		response.Success(c, payload)
		return
	}
	c.JSON(200, payload)
}

func collectRuntimeMetrics() gin.H {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	receiverMetrics := gin.H{}
	if receiver := logsyslog.GetGlobalReceiver(); receiver != nil {
		metrics := receiver.Metrics()
		receiverMetrics = gin.H{
			"udpReceived":    metrics.UDPReceived,
			"tcpReceived":    metrics.TCPReceived,
			"udpErrors":      metrics.UDPErrors,
			"tcpErrors":      metrics.TCPErrors,
			"parseErrors":    metrics.ParseErrors,
			"channelDropped": metrics.ChannelDropped,
			"tcpConnections": metrics.TCPConnections,
			"lastReceiveAt":  metrics.LastReceiveAt,
		}
	}

	pipelineMetrics := gin.H{}
	if pipeline := engine.GetGlobalPipeline(); pipeline != nil {
		metrics := pipeline.Metrics()
		pipelineMetrics = gin.H{
			"rawQueueDepth":    metrics.RawQueueDepth,
			"parsedQueueDepth": metrics.ParsedQueueDepth,
			"dbQueueDepth":     metrics.DBQueueDepth,
			"pushQueueDepth":   metrics.PushQueueDepth,
			"parseProcessed":   metrics.ParseProcessed,
			"parseErrors":      metrics.ParseErrors,
			"filterProcessed":  metrics.FilterProcessed,
			"filterDropped":    metrics.FilterDropped,
			"dbWritten":        metrics.DBWritten,
			"dbErrors":         metrics.DBErrors,
			"pushProcessed":    metrics.PushProcessed,
			"pushErrors":       metrics.PushErrors,
			"rawDropped":       metrics.RawDropped,
			"dbDropped":        metrics.DBDropped,
			"pushDropped":      metrics.PushDropped,
		}
	}

	return gin.H{
		"uptime_seconds": int64(time.Since(startTime).Seconds()),
		"status":         "running",
		"goroutines":     runtime.NumGoroutine(),
		"heap_alloc_mb":  float64(memStats.HeapAlloc) / 1024 / 1024,
		"heap_sys_mb":    float64(memStats.HeapSys) / 1024 / 1024,
		"heap_total_mb":  float64(memStats.TotalAlloc) / 1024 / 1024,
		"num_gc":         memStats.NumGC,
		"receiver":       receiverMetrics,
		"pipeline":       pipelineMetrics,
	}
}
