package repository

import (
	"strings"

	configpkg "syslog-alert/internal/config"
	"syslog-alert/internal/models"
	"syslog-alert/internal/service/cache"
	"syslog-alert/pkg/constants"
)

// GetSystemConfig 获取系统配置
func GetSystemConfig() models.SystemConfig {
	var cfg models.SystemConfig
	DB().First(&cfg)
	if normalizeSystemConfigProtocol(&cfg) {
		_ = DB().Model(&models.SystemConfig{}).Where("id = ?", cfg.ID).Update("protocol", cfg.Protocol).Error
	}
	paths := configpkg.Get()
	cfg.DataDir = paths.DataDir
	cfg.ConfigDir = paths.ConfigDir
	return cfg
}

// UpdateSystemConfig 更新系统配置
func UpdateSystemConfig(cfg models.SystemConfig) error {
	normalizeSystemConfigProtocol(&cfg)
	err := DB().Save(&cfg).Error
	if err == nil {
		cache.InvalidateSystemConfig()
		cache.InvalidateStatsCaches()
	}
	return err
}

func normalizeSystemConfigProtocol(cfg *models.SystemConfig) bool {
	if cfg == nil {
		return false
	}
	protocol := strings.ToLower(strings.TrimSpace(cfg.Protocol))
	if protocol != constants.ProtocolBoth {
		cfg.Protocol = constants.ProtocolBoth
		return true
	}
	cfg.Protocol = constants.ProtocolBoth
	return false
}
