// Package config 负责配置管理：数据目录、模板目录、数据库路径等路径解析。
package config

import (
	"os"
	"path/filepath"

	"syslog-alert/pkg/constants"
	"syslog-alert/pkg/logger"
)

// Paths 应用路径配置
type Paths struct {
	DataDir      string
	ConfigDir    string
	TemplatesDir string
	DatabasePath string
}

var paths *Paths

// Init 初始化路径配置，应在程序启动时调用一次。
func Init() *Paths {
	paths = &Paths{
		DataDir:      resolveDataDir(),
		TemplatesDir: resolveTemplatesDir(),
		ConfigDir:    resolveConfigDir(),
	}
	paths.DatabasePath = filepath.Join(paths.DataDir, constants.DatabaseFile)
	logger.Info("数据目录: %s", paths.DataDir)
	logger.Info("模板目录: %s", paths.TemplatesDir)
	return paths
}

// Get 获取已初始化的路径配置
func Get() *Paths {
	if paths == nil {
		return Init()
	}
	return paths
}

// resolveDataDir 解析数据目录
// 优先级：环境变量 > 可执行文件同级 data/ > 当前目录 data/
func resolveDataDir() string {
	if envDir := os.Getenv(constants.EnvDataDir); envDir != "" {
		return envDir
	}

	exePath, err := os.Executable()
	if err != nil {
		logger.Warn("获取可执行文件路径失败，使用 ./data: %v", err)
		return constants.DataDirName
	}
	dataPath := filepath.Join(filepath.Dir(exePath), constants.DataDirName)
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		logger.Warn("创建数据目录失败: %v", err)
	}
	return dataPath
}

// resolveTemplatesDir 解析模板目录
// 优先级：环境变量 > 可执行文件同级 templates/
func resolveTemplatesDir() string {
	if envDir := os.Getenv(constants.EnvTemplatesDir); envDir != "" {
		return envDir
	}

	exePath, err := os.Executable()
	if err != nil {
		logger.Warn("获取可执行文件路径失败，使用 ./templates: %v", err)
		return constants.TemplatesDirName
	}
	return filepath.Join(filepath.Dir(exePath), constants.TemplatesDirName)
}

// resolveConfigDir 解析配置目录
func resolveConfigDir() string {
	if envDir := os.Getenv(constants.EnvConfigDir); envDir != "" {
		return envDir
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		// 优先使用可执行文件同级 templates/
		configPath := filepath.Join(exeDir, constants.TemplatesDirName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}
	return constants.TemplatesDirName
}
