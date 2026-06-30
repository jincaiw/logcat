package repository

import (
	"syslog-alert/internal/models"

	"gorm.io/gorm"
)

// ---- 用户 ----

// CreateUser 创建用户
func CreateUser(user *models.User) error {
	return DB().Create(user).Error
}

// GetUserByUsername 按用户名查询用户
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB().Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByID 按 ID 查询用户
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := DB().First(&user, id).Error
	return &user, err
}

// GetFirstUser 获取第一个用户（单用户模式）
func GetFirstUser() (*models.User, error) {
	var user models.User
	err := DB().First(&user).Error
	return &user, err
}

// GetUserCount 返回用户总数
func GetUserCount() int64 {
	var count int64
	DB().Model(&models.User{}).Count(&count)
	return count
}

// UpdateUser 更新用户信息
func UpdateUser(user *models.User) error {
	return DB().Save(user).Error
}

// UpdateUserPassword 仅更新密码字段
func UpdateUserPassword(id uint, passwordHash string) error {
	return DB().Model(&models.User{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
}

// UpdateUserProfile 仅更新个人信息字段
func UpdateUserProfile(id uint, nickname, email, avatar string) error {
	updates := map[string]interface{}{
		"nickname": nickname,
		"email":    email,
		"avatar":   avatar,
	}
	return DB().Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// UserExists 判断是否存在任意用户（用于单用户模式判断是否已初始化）
func UserExists() bool {
	var count int64
	DB().Model(&models.User{}).Count(&count)
	return count > 0
}

// EnsureUserExists 确保至少存在一个用户记录，返回该用户。
// 用于单用户模式：若不存在则使用默认凭证创建。
func EnsureUserExists(defaultUsername, defaultPasswordHash, defaultNickname string) (*models.User, error) {
	var user models.User
	err := DB().First(&user).Error
	if err == nil {
		return &user, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	user = models.User{
		Username:     defaultUsername,
		PasswordHash: defaultPasswordHash,
		Nickname:     defaultNickname,
	}
	if err := DB().Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
