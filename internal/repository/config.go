package repository

import (
	configpkg "syslog-alert/internal/config"
	"syslog-alert/internal/models"
)

// GetSystemConfig 获取系统配置
func GetSystemConfig() models.SystemConfig {
	var cfg models.SystemConfig
	DB().First(&cfg)
	paths := configpkg.Get()
	cfg.DataDir = paths.DataDir
	cfg.ConfigDir = paths.ConfigDir
	return cfg
}

// UpdateSystemConfig 更新系统配置
func UpdateSystemConfig(cfg models.SystemConfig) error {
	return DB().Save(&cfg).Error
}
