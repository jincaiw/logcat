// Package logger 提供结构化日志能力，替代散落的 stdlog.Printf("[DEBUG] ...")。
// 按级别输出，支持调试开关，避免生产环境输出过多调试信息。
package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Level 日志级别
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	currentLevel = LevelInfo // 默认 Info 级别，可通过 SetLevel 调整
	prefixes     = map[Level]string{
		LevelDebug: "[DEBUG]",
		LevelInfo:  "[INFO] ",
		LevelWarn:  "[WARN] ",
		LevelError: "[ERROR]",
	}
)

// SetLevel 设置全局日志级别
func SetLevel(l Level) {
	currentLevel = l
}

// SetDebug 快捷开启/关闭调试日志
func SetDebug(enable bool) {
	if enable {
		currentLevel = LevelDebug
	} else {
		currentLevel = LevelInfo
	}
}

func logf(level Level, format string, args ...interface{}) {
	if level < currentLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	prefix := prefixes[level]
	// 去除多余空格，保持对齐
	log.Printf("%s %s", prefix, strings.TrimSpace(msg))
}

// Debug 调试日志
func Debug(format string, args ...interface{}) { logf(LevelDebug, format, args...) }

// Info 信息日志
func Info(format string, args ...interface{}) { logf(LevelInfo, format, args...) }

// Warn 警告日志
func Warn(format string, args ...interface{}) { logf(LevelWarn, format, args...) }

// Error 错误日志
func Error(format string, args ...interface{}) { logf(LevelError, format, args...) }

// Fatal 输出错误后退出
func Fatal(format string, args ...interface{}) {
	logf(LevelError, format, args...)
	os.Exit(1)
}

func init() {
	// 保持与标准 log 一致的标志：日期 + 时间
	log.SetFlags(log.LstdFlags)
}
