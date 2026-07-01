package cache

import (
	"fmt"
	"sync"
	"time"

	"syslog-alert/internal/models"
)

type ttlValueCache[T any] struct {
	mu      sync.Mutex
	value   T
	expires time.Time
	loaded  bool
}

func (c *ttlValueCache[T]) get(ttl time.Duration, loader func() T) T {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.loaded && time.Now().Before(c.expires) {
		return c.value
	}
	c.value = loader()
	c.expires = time.Now().Add(ttl)
	c.loaded = true
	return c.value
}

func (c *ttlValueCache[T]) invalidate() {
	c.mu.Lock()
	c.loaded = false
	var zero T
	c.value = zero
	c.expires = time.Time{}
	c.mu.Unlock()
}

type ttlMapCache[T any] struct {
	mu      sync.Mutex
	entries map[string]ttlEntry[T]
}

type ttlEntry[T any] struct {
	value   T
	expires time.Time
}

func (c *ttlMapCache[T]) get(key string, ttl time.Duration, loader func() T) T {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		c.entries = make(map[string]ttlEntry[T])
	}
	if entry, ok := c.entries[key]; ok && time.Now().Before(entry.expires) {
		return entry.value
	}
	value := loader()
	c.entries[key] = ttlEntry[T]{value: value, expires: time.Now().Add(ttl)}
	return value
}

func (c *ttlMapCache[T]) invalidate() {
	c.mu.Lock()
	c.entries = make(map[string]ttlEntry[T])
	c.mu.Unlock()
}

var (
	systemStatsCache     ttlValueCache[models.SystemStats]
	fieldStatsCache      ttlMapCache[models.FieldStatsResult]
	availableFieldsCache ttlMapCache[[]models.StatsField]
)

// GetCachedSystemStats 按 TTL 缓存系统统计结果。
func GetCachedSystemStats(ttl time.Duration, loader func() models.SystemStats) models.SystemStats {
	return systemStatsCache.get(ttl, loader)
}

// GetCachedFieldStats 按请求 key 缓存字段统计结果。
func GetCachedFieldStats(key string, ttl time.Duration, loader func() models.FieldStatsResult) models.FieldStatsResult {
	return fieldStatsCache.get(key, ttl, loader)
}

// GetCachedAvailableStatsFields 按 policy key 缓存可统计字段列表。
func GetCachedAvailableStatsFields(key string, ttl time.Duration, loader func() []models.StatsField) []models.StatsField {
	fields := availableFieldsCache.get(key, ttl, loader)
	return cloneStatsFields(fields)
}

// InvalidateStatsCaches 清理统计缓存。
func InvalidateStatsCaches() {
	systemStatsCache.invalidate()
	fieldStatsCache.invalidate()
	availableFieldsCache.invalidate()
}

func cloneStatsFields(src []models.StatsField) []models.StatsField {
	if len(src) == 0 {
		return nil
	}
	out := make([]models.StatsField, len(src))
	copy(out, src)
	return out
}

func statsFieldStatsKey(req models.FieldStatsRequest) string {
	return fmt.Sprintf("device:%d|policy:%d|start:%s|end:%s|field:%s|top:%d", req.DeviceID, req.FilterPolicyID, req.StartTime, req.EndTime, req.Field, req.TopN)
}

func StatsFieldStatsKey(req models.FieldStatsRequest) string {
	return statsFieldStatsKey(req)
}

func StatsAvailableFieldsKey(policyID uint) string {
	return fmt.Sprintf("policy:%d", policyID)
}
