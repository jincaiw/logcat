package repository

import (
	"syslog-alert/internal/models"
	"syslog-alert/pkg/constants"
)

// ---- 机器人/推送通道 ----

func CreateRobot(robot *models.Robot) error {
	return DB().Create(robot).Error
}

func GetRobots() []models.Robot {
	var robots []models.Robot
	DB().Find(&robots)
	return robots
}

func GetActiveRobotCount() int64 {
	var count int64
	DB().Model(&models.Robot{}).
		Where("is_active = ? AND (platform IN ? OR platform = '')", true, constants.SupportedNotificationPlatforms()).
		Count(&count)
	return count
}

func GetRobotByID(id uint) (*models.Robot, error) {
	var robot models.Robot
	err := DB().First(&robot, id).Error
	return &robot, err
}

func UpdateRobot(robot *models.Robot) error {
	return DB().Save(robot).Error
}

func DeleteRobot(id uint) error {
	return DB().Delete(&models.Robot{}, id).Error
}

// ---- 告警记录 ----

func CreateAlertRecord(record *models.AlertRecord) error {
	return DB().Create(record).Error
}

func GetAlertRecords(page, pageSize int) ([]models.AlertRecord, int64) {
	var records []models.AlertRecord
	var total int64

	DB().Model(&models.AlertRecord{}).Count(&total)
	DB().Order("sent_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records)

	return records, total
}

func GetAlertCount() int64 {
	var count int64
	DB().Model(&models.AlertRecord{}).Where("status = ?", "sent").Count(&count)
	return count
}
