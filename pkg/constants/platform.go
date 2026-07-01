package constants

import "strings"

// SupportedNotificationPlatforms 返回当前支持的通知渠道类型。
func SupportedNotificationPlatforms() []string {
	return []string{PlatformFeishu, PlatformEmail, PlatformSyslog, PlatformHTTP}
}

// NormalizeNotificationPlatform 规范化通知渠道类型；空值按默认飞书处理。
func NormalizeNotificationPlatform(platform string) (string, bool) {
	platform = strings.TrimSpace(platform)
	if platform == "" {
		return PlatformFeishu, true
	}

	switch platform {
	case PlatformFeishu, PlatformEmail, PlatformSyslog, PlatformHTTP:
		return platform, true
	default:
		return "", false
	}
}

// IsSupportedNotificationPlatform 判断通知渠道是否仍被当前版本支持。
func IsSupportedNotificationPlatform(platform string) bool {
	_, ok := NormalizeNotificationPlatform(platform)
	return ok
}
