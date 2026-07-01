// Package api 认证中间件与会话管理。
//
// 设计要点：
//   - 单用户模式：系统仅维护一个用户，登录后颁发随机 token
//   - 会话存储在内存中（sync.RWMutex 保护），进程重启后需重新登录
//   - 所有 /api/ 路由（除 /api/auth/login）均需携带有效 token
//   - token 通过 Authorization: Bearer <token> 头部传递
package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"
)

const (
	loginFailureWindow    = 5 * time.Minute
	loginFailureThreshold = 5
	loginBlockDuration    = 60 * time.Second
	authCleanupInterval   = 10 * time.Minute
)

// session 会话条目
type session struct {
	username  string
	userID    uint
	expiresAt time.Time
}

// sessionStore 内存会话存储
type sessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*session // token -> session
}

// loginAttempt 登录失败记录
type loginAttempt struct {
	firstFailure time.Time
	failureCount int
	blockedUntil time.Time
}

type loginAttemptStore struct {
	mu      sync.RWMutex
	entries map[string]*loginAttempt
}

var (
	authSessions    = &sessionStore{sessions: make(map[string]*session)}
	loginAttempts   = &loginAttemptStore{entries: make(map[string]*loginAttempt)}
	authCleanupOnce sync.Once
)

func ensureAuthMaintenance() {
	authCleanupOnce.Do(func() {
		go authMaintenanceLoop()
	})
}

func authMaintenanceLoop() {
	ticker := time.NewTicker(authCleanupInterval)
	defer ticker.Stop()

	cleanupAuthState := func() {
		now := time.Now()

		authSessions.mu.Lock()
		for token, s := range authSessions.sessions {
			if now.After(s.expiresAt) {
				delete(authSessions.sessions, token)
			}
		}
		authSessions.mu.Unlock()

		loginAttempts.mu.Lock()
		for key, attempt := range loginAttempts.entries {
			if !attempt.blockedUntil.IsZero() && now.After(attempt.blockedUntil) {
				delete(loginAttempts.entries, key)
				continue
			}
			if attempt.blockedUntil.IsZero() && !attempt.firstFailure.IsZero() && now.Sub(attempt.firstFailure) > loginFailureWindow {
				delete(loginAttempts.entries, key)
			}
		}
		loginAttempts.mu.Unlock()
	}

	cleanupAuthState()
	for range ticker.C {
		cleanupAuthState()
	}
}

// generateToken 生成随机 token（hex 编码）。
func generateToken() string {
	b := make([]byte, constants.AuthTokenBytes)
	if _, err := rand.Read(b); err != nil {
		applogger.Error("生成 token 失败: %v", err)
		// 降级：使用时间戳（极端情况，不应发生）
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(b)
}

// createSession 创建新会话并返回 token。
func createSession(userID uint, username string) string {
	ensureAuthMaintenance()
	token := generateToken()
	authSessions.mu.Lock()
	defer authSessions.mu.Unlock()
	authSessions.sessions[token] = &session{
		username:  username,
		userID:    userID,
		expiresAt: time.Now().Add(time.Duration(constants.AuthSessionTTL) * time.Second),
	}
	return token
}

// getSession 查询会话，若已过期则删除并返回 nil。
func getSession(token string) *session {
	authSessions.mu.RLock()
	s, ok := authSessions.sessions[token]
	authSessions.mu.RUnlock()
	if !ok {
		return nil
	}
	if time.Now().After(s.expiresAt) {
		authSessions.mu.Lock()
		delete(authSessions.sessions, token)
		authSessions.mu.Unlock()
		return nil
	}
	return s
}

// renewSession 续期会话（滑动窗口：每次请求后重置过期时间）。
func renewSession(token string) {
	authSessions.mu.Lock()
	defer authSessions.mu.Unlock()
	if s, ok := authSessions.sessions[token]; ok {
		s.expiresAt = time.Now().Add(time.Duration(constants.AuthSessionTTL) * time.Second)
	}
}

// removeSession 注销会话。
func removeSession(token string) {
	authSessions.mu.Lock()
	defer authSessions.mu.Unlock()
	delete(authSessions.sessions, token)
}

// extractToken 从请求头解析 Bearer token。
func extractToken(r *http.Request) string {
	auth := r.Header.Get(constants.AuthHeaderName)
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != constants.AuthTokenType {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func loginAttemptKey(r *http.Request) string {
	return clientIP(r)
}

func isLoginBlocked(r *http.Request) (time.Duration, bool) {
	key := loginAttemptKey(r)
	loginAttempts.mu.Lock()
	defer loginAttempts.mu.Unlock()

	attempt, ok := loginAttempts.entries[key]
	if !ok {
		return 0, false
	}
	if attempt.blockedUntil.IsZero() {
		return 0, false
	}
	if time.Now().After(attempt.blockedUntil) {
		delete(loginAttempts.entries, key)
		return 0, false
	}
	return time.Until(attempt.blockedUntil), true
}

func recordLoginFailure(r *http.Request) {
	key := loginAttemptKey(r)
	now := time.Now()

	loginAttempts.mu.Lock()
	defer loginAttempts.mu.Unlock()

	attempt, ok := loginAttempts.entries[key]
	if !ok {
		attempt = &loginAttempt{firstFailure: now, failureCount: 1}
		loginAttempts.entries[key] = attempt
		return
	}

	if !attempt.blockedUntil.IsZero() {
		if now.After(attempt.blockedUntil) {
			attempt.blockedUntil = time.Time{}
			attempt.firstFailure = now
			attempt.failureCount = 1
		}
		return
	}

	if attempt.firstFailure.IsZero() || now.Sub(attempt.firstFailure) > loginFailureWindow {
		attempt.firstFailure = now
		attempt.failureCount = 1
		return
	}

	attempt.failureCount++
	if attempt.failureCount >= loginFailureThreshold {
		attempt.blockedUntil = now.Add(loginBlockDuration)
		attempt.firstFailure = time.Time{}
		attempt.failureCount = 0
	}
}

func resetLoginFailures(r *http.Request) {
	key := loginAttemptKey(r)
	loginAttempts.mu.Lock()
	defer loginAttempts.mu.Unlock()
	delete(loginAttempts.entries, key)
}

// currentUserID 从 context 中获取已认证用户 ID。
type contextKey string

const userContextKey contextKey = "user_id"

// CurrentUserID 从请求 context 中获取当前用户 ID。
func CurrentUserID(r *http.Request) uint {
	if v, ok := r.Context().Value(userContextKey).(uint); ok {
		return v
	}
	return 0
}

// SecurityHeadersMiddleware 为所有响应添加基础安全头。
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("X-Content-Type-Options", "nosniff")
		headers.Set("X-Frame-Options", "DENY")
		headers.Set("Referrer-Policy", "no-referrer")
		headers.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware 认证中间件：校验 token，将用户 ID 注入 context，并自动续期会话。
// 未认证返回 401。
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ensureAuthMaintenance()
		token := extractToken(r)
		if token == "" {
			JSONError(w, "未认证，请先登录", http.StatusUnauthorized)
			return
		}
		s := getSession(token)
		if s == nil {
			JSONError(w, "会话已过期，请重新登录", http.StatusUnauthorized)
			return
		}
		renewSession(token)
		ctx := context.WithValue(r.Context(), userContextKey, s.userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireUser 校验当前请求存在已认证用户，并返回用户视图。
// 供 handler 调用，失败时自动写入 401 响应。
func requireUser(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	uid := CurrentUserID(r)
	if uid == 0 {
		JSONError(w, "未认证", http.StatusUnauthorized)
		return nil, false
	}
	user, err := repository.GetUserByID(uid)
	if err != nil {
		JSONError(w, "用户不存在", http.StatusUnauthorized)
		return nil, false
	}
	return user, true
}
