// Package repository 数据访问层，按实体拆分，每文件负责一个实体的 CRUD。
// 所有函数通过 database.Get() 获取 DB 单例，保持与原代码一致的调用方式。
package repository

import (
	"syslog-alert/internal/database"
	"syslog-alert/internal/models"

	"gorm.io/gorm"
)

// DB 获取数据库实例
func DB() *gorm.DB {
	return database.Get()
}

// ---- 设备分组 ----

func CreateDeviceGroup(group *models.DeviceGroup) error {
	return DB().Create(group).Error
}

func GetDeviceGroups() []models.DeviceGroup {
	var groups []models.DeviceGroup
	DB().Order("sort_order ASC").Find(&groups)
	return groups
}

func GetDeviceGroupByID(id uint) (*models.DeviceGroup, error) {
	var group models.DeviceGroup
	err := DB().First(&group, id).Error
	return &group, err
}

func UpdateDeviceGroup(group *models.DeviceGroup) error {
	return DB().Save(group).Error
}

func DeleteDeviceGroup(id uint) error {
	return DB().Delete(&models.DeviceGroup{}, id).Error
}

// ---- 设备 ----

func CreateDevice(device *models.Device) error {
	return DB().Create(device).Error
}

func GetDevices() []models.Device {
	var devices []models.Device
	DB().Find(&devices)
	return devices
}

func GetDeviceByID(id uint) (*models.Device, error) {
	var device models.Device
	err := DB().First(&device, id).Error
	return &device, err
}

func GetDeviceByIP(ip string) (*models.Device, error) {
	var device models.Device
	err := DB().Where("ip_address = ?", ip).First(&device).Error
	return &device, err
}

func UpdateDevice(device *models.Device) error {
	return DB().Save(device).Error
}

func DeleteDevice(id uint) error {
	return DB().Delete(&models.Device{}, id).Error
}

func GetDeviceCount() int64 {
	var count int64
	DB().Model(&models.Device{}).Count(&count)
	return count
}

func GetActiveDeviceCount() int64 {
	var count int64
	DB().Model(&models.Device{}).Where("is_active = ?", true).Count(&count)
	return count
}
