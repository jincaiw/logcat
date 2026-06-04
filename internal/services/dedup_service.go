package services

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// DedupCacheEntry represents a deduplication cache entry
type DedupCacheEntry struct {
	Hash      string
	FirstSeen time.Time
	LastSeen  time.Time
	Count     int
}

// DedupService handles log deduplication
type DedupService struct {
	mu        sync.RWMutex
	cache     map[string]*DedupCacheEntry
	stopCh    chan struct{}
	stopOnce  sync.Once
}

// NewDedupService creates a new DedupService
func NewDedupService() *DedupService {
	svc := &DedupService{
		cache:  make(map[string]*DedupCacheEntry),
		stopCh: make(chan struct{}),
	}
	// Start cleanup goroutine
	go svc.cleanupLoop()
	return svc
}

// Stop stops the cleanup loop. Safe to call multiple times.
func (s *DedupService) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
	})
}

// IsDuplicate checks if a log message is a duplicate within the time window
func (s *DedupService) IsDuplicate(rawMessage string, windowSeconds int) (bool, int) {
	hash := s.hash(rawMessage)

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	entry, exists := s.cache[hash]

	if !exists {
		s.cache[hash] = &DedupCacheEntry{
			Hash:      hash,
			FirstSeen: now,
			LastSeen:  now,
			Count:     1,
		}
		return false, 1
	}

	// Check if within the time window
	if now.Sub(entry.FirstSeen).Seconds() > float64(windowSeconds) {
		// Reset the entry
		entry.FirstSeen = now
		entry.LastSeen = now
		entry.Count = 1
		return false, 1
	}

	// It's a duplicate
	entry.LastSeen = now
	entry.Count++
	return true, entry.Count
}

// hash generates a hash for deduplication
func (s *DedupService) hash(message string) string {
	h := sha256.Sum256([]byte(message))
	return fmt.Sprintf("%x", h)
}

// cleanupLoop periodically cleans up expired entries
func (s *DedupService) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, entry := range s.cache {
			// Remove entries older than 1 hour
			if now.Sub(entry.FirstSeen) > time.Hour {
				delete(s.cache, key)
			}
		}
		s.mu.Unlock()
	}
}

// GetStats returns dedup statistics
func (s *DedupService) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"total_entries": len(s.cache),
	}
}

// defaultDedupService is the global dedup service instance
var defaultDedupService = NewDedupService()

// GetDedupService returns the global dedup service
func GetDedupService() *DedupService {
	return defaultDedupService
}