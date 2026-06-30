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
	"net/http"
	"strings"
	"sync"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"
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

var authSessions = &sessionStore{sessions: make(map[string]*session)}

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

// AuthMiddleware 认证中间件：校验 token，将用户 ID 注入 context，并自动续期会话。
// 未认证返回 401。
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
