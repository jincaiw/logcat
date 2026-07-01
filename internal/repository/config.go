package repository

import (
	configpkg "syslog-alert/internal/config"
	"syslog-alert/internal/models"
	"syslog-alert/internal/service/cache"
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
	err := DB().Save(&cfg).Error
	if err == nil {
		cache.InvalidateSystemConfig()
		cache.InvalidateStatsCaches()
	}
	return err
}
