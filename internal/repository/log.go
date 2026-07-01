package repository

import (
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/service/cache"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"
)

// ---- 日志 ----

func CreateLog(log *models.SyslogLog) error {
	err := DB().Create(log).Error
	if err == nil {
		cache.InvalidateStatsCaches()
	}
	return err
}

func GetLogs(page, pageSize int, deviceID *int, startTime, endTime, keyword string) ([]models.SyslogLog, int64) {
	var logs []models.SyslogLog
	var total int64

	query := DB().Model(&models.SyslogLog{})

	if deviceID != nil && *deviceID > 0 {
		query = query.Where("device_id = ?", *deviceID)
	}
	if startTime != "" {
		query = query.Where("received_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("received_at <= ?", endTime)
	}
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("raw_message LIKE ? OR parsed_fields LIKE ?", searchPattern, searchPattern)
	}

	query.Count(&total)
	query.Order("received_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	return logs, total
}

func GetLogCount() int64 {
	var count int64
	DB().Model(&models.SyslogLog{}).Count(&count)
	return count
}

func GetMatchedLogCount() int64 {
	var count int64
	DB().Model(&models.SyslogLog{}).Where("filter_status = ?", constants.FilterStatusMatched).Count(&count)
	return count
}

func GetUnmatchedLogsCount() int64 {
	var count int64
	DB().Model(&models.SyslogLog{}).Where("filter_status = ?", constants.FilterStatusUnmatched).Count(&count)
	return count
}

func UpdateLogFilterStatus(logID uint, status string, policyID uint) error {
	err := DB().Model(&models.SyslogLog{}).Where("id = ?", logID).Updates(map[string]interface{}{
		"filter_status":     status,
		"matched_policy_id": policyID,
	}).Error
	if err == nil {
		cache.InvalidateStatsCaches()
	}
	return err
}

func UpdateLogAlertStatus(logID uint, status string, policyID uint) error {
	err := DB().Model(&models.SyslogLog{}).Where("id = ?", logID).Updates(map[string]interface{}{
		"alert_status":    status,
		"alert_policy_id": policyID,
	}).Error
	if err == nil {
		cache.InvalidateStatsCaches()
	}
	return err
}

func UpdateLogParsedFields(logID uint, parsedData, parsedFields string) error {
	err := DB().Model(&models.SyslogLog{}).Where("id = ?", logID).Updates(map[string]interface{}{
		"parsed_data":   parsedData,
		"parsed_fields": parsedFields,
	}).Error
	if err == nil {
		cache.InvalidateStatsCaches()
	}
	return err
}

func DeleteLog(logID uint) error {
	err := DB().Delete(&models.SyslogLog{}, logID).Error
	if err == nil {
		cache.InvalidateStatsCaches()
	}
	return err
}

// CleanupOldLogs 清理指定天数前的日志
func CleanupOldLogs(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := DB().Where("received_at < ?", cutoff).Delete(&models.SyslogLog{})
	if result.Error != nil {
		return result.Error
	}

	applogger.Info("已清理 %d 条旧日志（截止: %s）", result.RowsAffected, cutoff.Format("2006-01-02 15:04:05"))

	cache.InvalidateStatsCaches()
	if result.RowsAffected > 1000 {
		go func() {
			database := DB()
			sqlDB, _ := database.DB()
			if sqlDB != nil {
				sqlDB.Exec("PRAGMA wal_checkpoint(FULL)")
				sqlDB.Exec("VACUUM")
				applogger.Info("数据库压缩完成（删除 %d 行）", result.RowsAffected)
			}
		}()
	}
	return nil
}

// CleanupUnmatchedLogs 清理指定天数前的未匹配日志
func CleanupUnmatchedLogs(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := DB().Where("filter_status = ? AND received_at < ?", constants.FilterStatusUnmatched, cutoff).Delete(&models.SyslogLog{})
	if result.Error != nil {
		return result.Error
	}

	applogger.Info("已清理 %d 条未匹配日志", result.RowsAffected)

	if result.RowsAffected > 1000 {
		go func() {
			sqlDB, _ := DB().DB()
			if sqlDB != nil {
				sqlDB.Exec("PRAGMA wal_checkpoint(FULL)")
				sqlDB.Exec("VACUUM")
			}
		}()
	}
	return nil
}

// CleanupAllLogs 清空所有日志
func CleanupAllLogs() error {
	result := DB().Exec("DELETE FROM syslog_logs")
	if result.Error != nil {
		return result.Error
	}
	applogger.Info("已清空所有日志")
	DB().Exec("PRAGMA wal_checkpoint(PASSIVE)")
	DB().Exec("VACUUM")
	return nil
}

// GetLogByID 根据 ID 获取日志
func GetLogByID(id uint) (*models.SyslogLog, error) {
	var log models.SyslogLog
	err := DB().First(&log, id).Error
	return &log, err
}
