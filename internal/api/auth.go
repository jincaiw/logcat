// Package api 认证相关 handler：登录、登出、获取/更新个人信息、修改密码。
package api

import (
	"net/http"
	"strconv"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// loginRequest 登录请求
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// loginResponse 登录响应
type loginResponse struct {
	Token string          `json:"token"`
	User  models.UserView `json:"user"`
}

// Login 用户登录，校验用户名密码后颁发 token。
func (ws *WebServer) Login(w http.ResponseWriter, r *http.Request) {
	ensureAuthMaintenance()
	if retryAfter, blocked := isLoginBlocked(r); blocked {
		w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
		JSONError(w, "登录失败次数过多，请稍后再试", http.StatusTooManyRequests)
		return
	}

	var req loginRequest
	if !DecodeJSON(w, r, &req) {
		return
	}
	if req.Username == "" || req.Password == "" {
		recordLoginFailure(r)
		JSONError(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}

	user, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		applogger.Warn("登录失败：用户不存在 %s", req.Username)
		recordLoginFailure(r)
		JSONError(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		applogger.Warn("登录失败：密码错误 %s", req.Username)
		recordLoginFailure(r)
		JSONError(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	resetLoginFailures(r)
	token := createSession(user.ID, user.Username)
	applogger.Info("用户登录成功: %s", user.Username)

	JSONResponse(w, loginResponse{
		Token: token,
		User:  user.ToView(),
	})
}

// Logout 用户登出，销毁当前会话。
func (ws *WebServer) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token != "" {
		removeSession(token)
	}
	JSONResponse(w, map[string]string{"message": "已登出"})
}

// GetProfile 获取当前登录用户个人信息。
func (ws *WebServer) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := requireUser(w, r)
	if !ok {
		return
	}
	JSONResponse(w, user.ToView())
}

// updateProfileRequest 更新个人信息请求
type updateProfileRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// UpdateProfile 更新当前登录用户个人信息（昵称、邮箱、头像）。
func (ws *WebServer) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := requireUser(w, r)
	if !ok {
		return
	}
	var req updateProfileRequest
	if !DecodeJSON(w, r, &req) {
		return
	}
	if req.Nickname == "" {
		JSONError(w, "昵称不能为空", http.StatusBadRequest)
		return
	}
	if err := repository.UpdateUserProfile(user.ID, req.Nickname, req.Email, req.Avatar); err != nil {
		JSONError(w, "更新个人信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	updated, err := repository.GetUserByID(user.ID)
	if err != nil {
		JSONError(w, "获取用户信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, updated.ToView())
}

// changePasswordRequest 修改密码请求
type changePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// ChangePassword 修改当前登录用户密码。
func (ws *WebServer) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := requireUser(w, r)
	if !ok {
		return
	}
	var req changePasswordRequest
	if !DecodeJSON(w, r, &req) {
		return
	}
	if req.OldPassword == "" || req.NewPassword == "" {
		JSONError(w, "原密码和新密码不能为空", http.StatusBadRequest)
		return
	}
	if len(req.NewPassword) < 6 {
		JSONError(w, "新密码长度不能少于 6 位", http.StatusBadRequest)
		return
	}

	// 校验原密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		JSONError(w, "原密码错误", http.StatusBadRequest)
		return
	}

	// 生成新密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), constants.AuthBCryptCost)
	if err != nil {
		JSONError(w, "密码加密失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repository.UpdateUserPassword(user.ID, string(hash)); err != nil {
		JSONError(w, "修改密码失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	applogger.Info("用户修改密码成功: %s", user.Username)
	JSONResponse(w, map[string]string{"message": "密码修改成功"})
}
