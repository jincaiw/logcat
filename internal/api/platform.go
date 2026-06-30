package api

import "syslog-alert/pkg/constants"

func normalizeNotificationPlatform(platform string) (string, bool) {
	return constants.NormalizeNotificationPlatform(platform)
}
