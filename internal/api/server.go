package api

import (
	"os"
	"runtime"
	"slices"
	"sync"
	"time"

	"syslog-alert/internal/config"
	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/internal/service/cache"
	"syslog-alert/internal/service/syslog"
	"syslog-alert/pkg/constants"
)

// WebServer Web 应用服务器，持有 Syslog 服务与运行时统计。
type WebServer struct {
	syslogServer *syslog.Server
	startTime    time.Time
	stats        models.SystemStats
	statsMutex   sync.RWMutex
}

// NewWebServer 创建 Web 服务器实例。
func NewWebServer() *WebServer {
	ws := &WebServer{
		startTime: time.Now(),
		stats: models.SystemStats{
			ListenPort: constants.DefaultListenPort,
		},
	}
	ws.syslogServer = syslog.NewServer(ws)
	return ws
}

// UpdateStats 实现 syslog.StatsUpdater 接口，更新运行时统计。
func (ws *WebServer) UpdateStats(logs int64, devices int, running bool) {
	ws.statsMutex.Lock()
	defer ws.statsMutex.Unlock()
	ws.stats.TotalLogs = logs
	ws.stats.DeviceCount = devices
	ws.stats.ServiceRunning = running
}

// GetSystemStats 获取系统统计信息。
func (ws *WebServer) GetSystemStats() models.SystemStats {
	return cache.GetCachedSystemStats(3*time.Second, ws.loadSystemStats)
}

func (ws *WebServer) loadSystemStats() models.SystemStats {
	ws.statsMutex.RLock()
	defer ws.statsMutex.RUnlock()

	stats := ws.stats
	stats.StartTime = ws.startTime.Format("2006-01-02 15:04:05")
	stats.ListenPort = ws.syslogServer.GetPort()
	stats.ServiceRunning = ws.syslogServer.IsRunning()
	stats.ReceiveRate = ws.syslogServer.GetReceiveRate()
	stats.Connections = ws.syslogServer.GetConnections()
	stats.DroppedLogs = ws.syslogServer.GetDroppedCount()
	stats.QueueLength = ws.syslogServer.GetQueueLength()
	stats.TraceCacheSize = ws.syslogServer.GetTraceCacheSize()
	stats.Protocol = cache.GetSystemConfig(repository.GetSystemConfig).Protocol

	// 日志统计
	stats.TotalLogs = repository.GetLogCount()
	stats.MatchedLogs = repository.GetMatchedLogCount()
	stats.UnmatchedLogs = repository.GetUnmatchedLogsCount()
	stats.AlertCount = repository.GetAlertCount()

	// 设备与配置统计
	stats.DeviceCount = int(repository.GetDeviceCount())
	stats.ActiveDevices = int(repository.GetActiveDeviceCount())
	stats.ParseTemplateCount = int64(len(cache.GetParseTemplates(repository.GetParseTemplates)))
	stats.ActiveFilterPolicies = func() int {
		count := 0
		for _, policy := range cache.GetFilterPolicies(repository.GetFilterPolicies) {
			if policy.IsActive {
				count++
			}
		}
		return count
	}()
	stats.ActiveAlertPolicies = len(repository.GetActiveAlertPolicies())
	stats.ActiveRobots = func() int {
		count := 0
		for _, robot := range cache.GetRobots(repository.GetRobots) {
			if robot.IsActive && (robot.Platform == "" || slices.Contains(constants.SupportedNotificationPlatforms(), robot.Platform)) {
				count++
			}
		}
		return count
	}()

	// 运行时统计
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	stats.MemoryUsage = memStats.Alloc / 1024 / 1024 // MB
	stats.GoroutineCount = runtime.NumGoroutine()

	// 数据库文件大小
	if dbPath := config.Get().DatabasePath; dbPath != "" {
		if fi, err := os.Stat(dbPath); err == nil {
			stats.DatabaseSize = fi.Size()
		}
	}

	return stats
}

// SyslogServer 返回 Syslog 服务器实例（供 handler 使用）。
func (ws *WebServer) SyslogServer() *syslog.Server {
	return ws.syslogServer
}

// StartTime 返回服务启动时间。
func (ws *WebServer) StartTime() time.Time {
	return ws.startTime
}
