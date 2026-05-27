package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SessionInfo holds user session data
type SessionInfo struct {
	UserID   uint
	Username string
	ExpireAt time.Time
}

// SessionStore is an in-memory session store
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*SessionInfo
}

// NewSessionStore creates a new in-memory session store
func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*SessionInfo),
	}
	// Start cleanup goroutine
	go store.cleanupExpired()
	return store
}

// Create creates a new session and returns the session token
func (s *SessionStore) Create(userID uint, username string, expireHours int) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	token := uuid.New().String()
	s.sessions[token] = &SessionInfo{
		UserID:   userID,
		Username: username,
		ExpireAt: time.Now().Add(time.Duration(expireHours) * time.Hour),
	}
	return token
}

// Get retrieves session info by token
func (s *SessionStore) Get(token string) *SessionInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[token]
	if !ok {
		return nil
	}
	if time.Now().After(session.ExpireAt) {
		return nil
	}
	return session
}

// Delete removes a session
func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

// DeleteByUserID removes all sessions for a user
func (s *SessionStore) DeleteByUserID(userID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for token, session := range s.sessions {
		if session.UserID == userID {
			delete(s.sessions, token)
		}
	}
}

// cleanupExpired periodically removes expired sessions
func (s *SessionStore) cleanupExpired() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for token, session := range s.sessions {
			if now.After(session.ExpireAt) {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}

// DefaultSessionStore is the global session store
var DefaultSessionStore = NewSessionStore()

// Context keys for storing user info
const (
	ContextKeyUserID   = "UserID"
	ContextKeyUsername = "Username"
	ContextKeyUserRoles = "UserRoles"
	ContextKeyRequestID = "RequestID"
)

// RequestID generates a request ID for each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Set(ContextKeyRequestID, rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}

// AuthRequired checks for a valid session cookie
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":      401,
				"message":   "unauthorized: no session token",
				"requestId": c.GetString(ContextKeyRequestID),
			})
			return
		}

		session := DefaultSessionStore.Get(token)
		if session == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":      401,
				"message":   "unauthorized: invalid or expired session",
				"requestId": c.GetString(ContextKeyRequestID),
			})
			return
		}

		// Set user info in context
		c.Set(ContextKeyUserID, session.UserID)
		c.Set(ContextKeyUsername, session.Username)
		c.Next()
	}
}

// GetUserID extracts the user ID from the context
func GetUserID(c *gin.Context) uint {
	if id, exists := c.Get(ContextKeyUserID); exists {
		if uid, ok := id.(uint); ok {
			return uid
		}
	}
	return 0
}

// GetUsername extracts the username from the context
func GetUsername(c *gin.Context) string {
	if name, exists := c.Get(ContextKeyUsername); exists {
		if s, ok := name.(string); ok {
			return s
		}
	}
	return ""
}

// GetRequestID extracts the request ID from the context
func GetRequestID(c *gin.Context) string {
	if rid, exists := c.Get(ContextKeyRequestID); exists {
		if s, ok := rid.(string); ok {
			return s
		}
	}
	return ""
}