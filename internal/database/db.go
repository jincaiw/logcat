// Package database 负责数据库连接、迁移与种子数据初始化。
package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"syslog-alert/internal/config"
	"syslog-alert/internal/models"
	applogger "syslog-alert/pkg/logger"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Get 返回数据库单例
func Get() *gorm.DB {
	once.Do(func() {
		if err := initDB(); err != nil {
			applogger.Fatal("初始化数据库失败: %v", err)
		}
	})
	return db
}

func initDB() error {
	paths := config.Get()
	if err := os.MkdirAll(paths.DataDir, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %w", err)
	}

	dbPath := paths.DatabasePath
	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=10000&_sync=NORMAL&_cache_size=-64000", dbPath)

	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	db = gormDB

	// 连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(0)

	// SQLite PRAGMA 优化
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = -64000",
		"PRAGMA temp_store = MEMORY",
		"PRAGMA mmap_size = 268435456",
	}
	for _, p := range pragmas {
		if err := db.Exec(p).Error; err != nil {
			applogger.Warn("执行 PRAGMA 失败 (%s): %v", p, err)
		}
	}

	if err := autoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	ensureIndexes()
	seedDefaults()
	loadTemplatesFromDir()

	applogger.Info("数据库初始化完成: %s", dbPath)
	return nil
}

func autoMigrate() error {
	// 清理旧索引（历史遗留）
	db.Exec("DROP INDEX IF EXISTS idx_field_mapping_docs_device_type")

	return db.AutoMigrate(
		&models.User{},
		&models.DeviceGroup{},
		&models.Device{},
		&models.ParseTemplate{},
		&models.OutputTemplate{},
		&models.FilterPolicy{},
		&models.AlertPolicy{},
		&models.SyslogLog{},
		&models.Robot{},
		&models.AlertRule{},
		&models.AlertRecord{},
		&models.SystemConfig{},
		&models.FieldMappingDoc{},
	)
}

func ensureIndexes() {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_syslog_logs_received_at ON syslog_logs(received_at)",
		"CREATE INDEX IF NOT EXISTS idx_syslog_logs_device_id ON syslog_logs(device_id)",
		"CREATE INDEX IF NOT EXISTS idx_syslog_logs_filter_status ON syslog_logs(filter_status)",
		"CREATE INDEX IF NOT EXISTS idx_syslog_logs_alert_status ON syslog_logs(alert_status)",
		"CREATE INDEX IF NOT EXISTS idx_syslog_logs_device_received_at ON syslog_logs(device_id, received_at)",
		"CREATE INDEX IF NOT EXISTS idx_syslog_logs_device_policy_received_at ON syslog_logs(device_id, matched_policy_id, received_at)",
		"CREATE INDEX IF NOT EXISTS idx_syslog_logs_policy_received_at ON syslog_logs(matched_policy_id, received_at)",
		"CREATE INDEX IF NOT EXISTS idx_syslog_logs_filter_received_at ON syslog_logs(filter_status, received_at)",
		"CREATE INDEX IF NOT EXISTS idx_syslog_logs_alert_received_at ON syslog_logs(alert_status, received_at)",
		"CREATE INDEX IF NOT EXISTS idx_alert_records_created_at ON alert_records(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_alert_records_status_sent_at ON alert_records(status, sent_at)",
		"CREATE INDEX IF NOT EXISTS idx_devices_enabled ON devices(is_active)",
	}
	for _, idx := range indexes {
		db.Exec(idx)
	}
	applogger.Debug("数据库索引已确保")
}

// Vacuum 执行数据库压缩
func Vacuum() {
	db.Exec("PRAGMA wal_checkpoint(PASSIVE)")
	db.Exec("VACUUM")
}

// VacuumFull 执行完整压缩（阻塞）
func VacuumFull() {
	sqlDB, err := db.DB()
	if err == nil && sqlDB != nil {
		sqlDB.Exec("PRAGMA wal_checkpoint(FULL)")
		sqlDB.Exec("VACUUM")
	}
}

// loadTemplatesFromDir 从 templates/ 目录加载 JSON 配置文件
// 兼容旧逻辑：仅在记录不存在时创建
func loadTemplatesFromDir() {
	templatesDir := config.Get().TemplatesDir
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		applogger.Debug("模板目录不存在: %s", templatesDir)
		return
	}

	parseTemplatesFile := filepath.Join(templatesDir, "parse_templates.json")
	if _, err := os.Stat(parseTemplatesFile); err == nil {
		loadParseTemplatesFromFile(parseTemplatesFile)
	}

	filterPoliciesFile := filepath.Join(templatesDir, "filter_policies.json")
	if _, err := os.Stat(filterPoliciesFile); err == nil {
		loadFilterPoliciesFromFile(filterPoliciesFile)
	}
}
