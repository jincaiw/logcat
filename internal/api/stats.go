package api

import (
	"net/http"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/internal/service/cache"
)

// GetStats 返回系统统计信息。
func (ws *WebServer) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := ws.GetSystemStats()
	JSONResponse(w, stats)
}

// GetFieldStats 按字段聚合统计。
func (ws *WebServer) GetFieldStats(w http.ResponseWriter, r *http.Request) {
	var req models.FieldStatsRequest
	if !DecodeJSON(w, r, &req) {
		return
	}
	key := cache.StatsFieldStatsKey(req)
	result := cache.GetCachedFieldStats(key, 5*time.Second, func() models.FieldStatsResult {
		return repository.GetFieldStats(req)
	})
	JSONResponse(w, result)
}

// GetAvailableStatsFields 获取指定筛选策略下可统计的字段列表。
func (ws *WebServer) GetAvailableStatsFields(w http.ResponseWriter, r *http.Request) {
	policyID, ok := ParseUintID(w, r.PathValue("policyId"))
	if !ok {
		return
	}
	key := cache.StatsAvailableFieldsKey(policyID)
	fields := cache.GetCachedAvailableStatsFields(key, 30*time.Second, func() []models.StatsField {
		return repository.GetAvailableStatsFields(policyID)
	})
	if fields == nil {
		fields = []models.StatsField{}
	}
	JSONResponse(w, fields)
}
