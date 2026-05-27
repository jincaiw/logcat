package services

import (
	"errors"
	"sync"
	"time"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// HighFreqIPEntry tracks high-frequency IP access
type HighFreqIPEntry struct {
	IP        string
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
}

// HighFreqService handles high-frequency IP detection
type HighFreqService struct {
	mu          sync.RWMutex
	ipCounters  map[string]*HighFreqIPEntry
	threshold   int           // default threshold per window
	window      time.Duration  // time window
	stopCh      chan struct{}
}

// NewHighFreqService creates a new HighFreqService
func NewHighFreqService() *HighFreqService {
	svc := &HighFreqService{
		ipCounters: make(map[string]*HighFreqIPEntry),
		threshold:  100,
		window:     60 * time.Second,
		stopCh:     make(chan struct{}),
	}
	go svc.cleanupLoop()
	return svc
}

// Detect records an access from an IP and checks if it exceeds threshold
func (s *HighFreqService) Detect(ip string) (bool, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	entry, exists := s.ipCounters[ip]

	if !exists {
		s.ipCounters[ip] = &HighFreqIPEntry{
			IP:        ip,
			Count:     1,
			FirstSeen: now,
			LastSeen:  now,
		}
		return false, 1
	}

	// Check if within window
	if now.Sub(entry.FirstSeen) > s.window {
		entry.FirstSeen = now
		entry.LastSeen = now
		entry.Count = 1
		return false, 1
	}

	entry.LastSeen = now
	entry.Count++

	return entry.Count >= s.threshold, entry.Count
}

// GetHighFreqIPs returns all IPs that exceed the threshold
func (s *HighFreqService) GetHighFreqIPs() []HighFreqIPEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []HighFreqIPEntry
	now := time.Now()
	for _, entry := range s.ipCounters {
		if entry.Count >= s.threshold && now.Sub(entry.FirstSeen) <= s.window {
			result = append(result, *entry)
		}
	}
	return result
}

// UpdateConfig updates the detection threshold and window
func (s *HighFreqService) UpdateConfig(threshold int, windowSeconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.threshold = threshold
	s.window = time.Duration(windowSeconds) * time.Second
}

// GetConfig returns current configuration
func (s *HighFreqService) GetConfig() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]interface{}{
		"threshold":      s.threshold,
		"window_seconds": int(s.window.Seconds()),
	}
}

// cleanupLoop periodically cleans up expired entries
func (s *HighFreqService) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			now := time.Now()
			for key, entry := range s.ipCounters {
				if now.Sub(entry.FirstSeen) > s.window*2 {
					delete(s.ipCounters, key)
				}
			}
			s.mu.Unlock()
		case <-s.stopCh:
			return
		}
	}
}

// PersistHighFreqIP persists a high-frequency IP as a system config for future reference
func (s *HighFreqService) PersistHighFreqIP(ip string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	var config models.SystemConfig
	result := db.Where("config_key = ?", "high_freq_ip:"+ip).First(&config)
	if result.Error != nil {
		config = models.SystemConfig{
			ConfigKey:   "high_freq_ip:" + ip,
			ConfigValue: ip,
			Description: "High frequency IP detected",
		}
		return db.Create(&config).Error
	}

	config.UpdatedAt = time.Now()
	return db.Save(&config).Error
}

// defaultHighFreqService is the global instance
var defaultHighFreqService = NewHighFreqService()

// GetHighFreqService returns the global HighFreqService
func GetHighFreqService() *HighFreqService {
	return defaultHighFreqService
}