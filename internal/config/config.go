package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Syslog   SyslogConfig   `yaml:"syslog"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
	Log      LogConfig      `yaml:"log"`
	Worker   WorkerConfig   `yaml:"worker"`
	Queue    QueueConfig    `yaml:"queue"`
	CORS     CORSConfig     `yaml:"cors"`
	Theme    string         `yaml:"theme"`
	Debug    bool           `yaml:"debug"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// SyslogConfig holds syslog receiver configuration
type SyslogConfig struct {
	Enabled bool `yaml:"enabled"`
	UDPPort int  `yaml:"udp_port"`
	TCPPort int  `yaml:"tcp_port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type        string       `yaml:"type"`
	AutoMigrate bool         `yaml:"auto_migrate"`
	SQLite      SQLiteConfig `yaml:"sqlite"`
	MySQL       MySQLConfig  `yaml:"mysql"`
}

// SQLiteConfig holds SQLite-specific configuration
type SQLiteConfig struct {
	Path string `yaml:"path"`
	WAL  bool   `yaml:"wal"`
}

// MySQLConfig holds MySQL-specific configuration
type MySQLConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Database     string `yaml:"database"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Charset      string `yaml:"charset"`
	Timezone     string `yaml:"timezone"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	SessionExpireHours  int `yaml:"session_expire_hours"`
	MaxFailedLogin      int `yaml:"max_failed_login"`
	LockDurationMinutes int `yaml:"lock_duration_minutes"`
}

// LogConfig holds log retention configuration
type LogConfig struct {
	RetentionDays          int `yaml:"retention_days"`
	MaxLogSize             int `yaml:"max_log_size"`
	UnmatchedRetentionDays int `yaml:"unmatched_retention_days"`
}

// WorkerConfig holds worker pool configuration
type WorkerConfig struct {
	ParseWorkers  int `yaml:"parse_workers"`
	FilterWorkers int `yaml:"filter_workers"`
	PushWorkers   int `yaml:"push_workers"`
}

// QueueConfig holds queue configuration
type QueueConfig struct {
	Capacity   int    `yaml:"capacity"`
	FullPolicy string `yaml:"full_policy"`
}

var (
	configPath string
	loaded     *Config
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
}

// Load loads configuration from the file and applies environment variable overrides.
// It should only be called once after flag.Parse().
func Load() (*Config, error) {
	if loaded != nil {
		return loaded, nil
	}

	cfg := &Config{}

	// Set defaults
	cfg.Server = ServerConfig{Host: "0.0.0.0", Port: 5080}
	cfg.Database = DatabaseConfig{
		Type:        "sqlite",
		AutoMigrate: true,
		SQLite:      SQLiteConfig{Path: "data/logcat.db", WAL: true},
		MySQL: MySQLConfig{
			Host:         "127.0.0.1",
			Port:         3306,
			Database:     "logcat",
			Username:     "logcat",
			Password:     "",
			Charset:      "utf8mb4",
			Timezone:     "Asia/Shanghai",
			MaxOpenConns: 50,
			MaxIdleConns: 10,
		},
	}
	cfg.Auth = AuthConfig{
		SessionExpireHours:  24,
		MaxFailedLogin:      5,
		LockDurationMinutes: 30,
	}
	cfg.Log = LogConfig{
		RetentionDays:          90,
		MaxLogSize:             10000,
		UnmatchedRetentionDays: 30,
	}
	cfg.Worker = WorkerConfig{
		ParseWorkers:  4,
		FilterWorkers: 4,
		PushWorkers:   4,
	}
	cfg.Queue = QueueConfig{
		Capacity:   10000,
		FullPolicy: "block_drop",
	}
	cfg.Theme = "light"

	// Load from config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
		}
		// Config file not found is OK, use defaults
	} else {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
		}
	}

	// Apply environment variable overrides
	applyEnvOverrides(cfg)

	loaded = cfg
	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("LOGCAT_SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("LOGCAT_SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("LOGCAT_DATABASE_TYPE"); v != "" {
		cfg.Database.Type = v
	}
	if v := os.Getenv("LOGCAT_DATABASE_AUTO_MIGRATE"); v != "" {
		cfg.Database.AutoMigrate = v == "true" || v == "1"
	}
	if v := os.Getenv("LOGCAT_SQLITE_PATH"); v != "" {
		cfg.Database.SQLite.Path = v
	}
	if v := os.Getenv("LOGCAT_MYSQL_HOST"); v != "" {
		cfg.Database.MySQL.Host = v
	}
	if v := os.Getenv("LOGCAT_MYSQL_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.MySQL.Port = port
		}
	}
	if v := os.Getenv("LOGCAT_MYSQL_DATABASE"); v != "" {
		cfg.Database.MySQL.Database = v
	}
	if v := os.Getenv("LOGCAT_MYSQL_USERNAME"); v != "" {
		cfg.Database.MySQL.Username = v
	}
	if v := os.Getenv("LOGCAT_MYSQL_PASSWORD"); v != "" {
		cfg.Database.MySQL.Password = v
	}
	if v := os.Getenv("LOGCAT_MYSQL_CHARSET"); v != "" {
		cfg.Database.MySQL.Charset = v
	}
	if v := os.Getenv("LOGCAT_MYSQL_TIMEZONE"); v != "" {
		cfg.Database.MySQL.Timezone = v
	}
	if v := os.Getenv("LOGCAT_MYSQL_MAX_OPEN_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Database.MySQL.MaxOpenConns = n
		}
	}
	if v := os.Getenv("LOGCAT_MYSQL_MAX_IDLE_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Database.MySQL.MaxIdleConns = n
		}
	}
	if v := os.Getenv("LOGCAT_SESSION_EXPIRE_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Auth.SessionExpireHours = n
		}
	}
	if v := os.Getenv("LOGCAT_MAX_FAILED_LOGIN"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Auth.MaxFailedLogin = n
		}
	}
	if v := os.Getenv("LOGCAT_LOCK_DURATION_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Auth.LockDurationMinutes = n
		}
	}
	if v := os.Getenv("LOGCAT_RETENTION_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Log.RetentionDays = n
		}
	}
	if v := os.Getenv("LOGCAT_UNMATCHED_RETENTION_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Log.UnmatchedRetentionDays = n
		}
	}
	if v := os.Getenv("LOGCAT_MAX_LOG_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Log.MaxLogSize = n
		}
	}
	if v := os.Getenv("LOGCAT_THEME"); v != "" {
		cfg.Theme = v
	}
	if v := os.Getenv("LOGCAT_SYSLOG_ENABLED"); v != "" {
		cfg.Syslog.Enabled = v == "true" || v == "1"
	}
	if v := os.Getenv("LOGCAT_SYSLOG_UDP_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Syslog.UDPPort = port
		}
	}
	if v := os.Getenv("LOGCAT_SYSLOG_TCP_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Syslog.TCPPort = port
		}
	}
	if v := os.Getenv("LOGCAT_PARSE_WORKERS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Worker.ParseWorkers = n
		}
	}
	if v := os.Getenv("LOGCAT_FILTER_WORKERS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Worker.FilterWorkers = n
		}
	}
	if v := os.Getenv("LOGCAT_PUSH_WORKERS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Worker.PushWorkers = n
		}
	}
	if v := os.Getenv("LOGCAT_QUEUE_CAPACITY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Queue.Capacity = n
		}
	}
	if v := os.Getenv("LOGCAT_QUEUE_FULL_POLICY"); v != "" {
		cfg.Queue.FullPolicy = v
	}
}

// Get returns the loaded config, or nil if not loaded yet
func Get() *Config {
	return loaded
}

// SetForTest replaces the loaded configuration and returns a restore function.
func SetForTest(cfg *Config) func() {
	previous := loaded
	loaded = cfg
	return func() {
		loaded = previous
	}
}

// GetAdminPassword returns the admin password from env or empty string
func GetAdminPassword() string {
	return os.Getenv("LOGCAT_ADMIN_PASSWORD")
}

// DSN returns the database connection string based on config
func (c *Config) DSN() string {
	if c.Database.Type == "mysql" {
		m := c.Database.MySQL
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=%s",
			m.Username, m.Password, m.Host, m.Port, m.Database, m.Charset, m.Timezone)
	}
	// SQLite
	return c.Database.SQLite.Path
}
