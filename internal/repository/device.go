// Package repository 数据访问层，按实体拆分，每文件负责一个实体的 CRUD。
// 所有函数通过 database.Get() 获取 DB 单例，保持与原代码一致的调用方式。
package repository

import (
	"syslog-alert/internal/database"
	"syslog-alert/internal/models"
	"syslog-alert/internal/service/cache"

	"gorm.io/gorm"
)

// DB 获取数据库实例
func DB() *gorm.DB {
	return database.Get()
}

// ---- 设备分组 ----

func CreateDeviceGroup(group *models.DeviceGroup) error {
	err := DB().Create(group).Error
	if err == nil {
		cache.InvalidateDevices()
	}
	return err
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
	err := DB().Save(group).Error
	if err == nil {
		cache.InvalidateDevices()
	}
	return err
}

func DeleteDeviceGroup(id uint) error {
	err := DB().Delete(&models.DeviceGroup{}, id).Error
	if err == nil {
		cache.InvalidateDevices()
	}
	return err
}

// ---- 设备 ----

func CreateDevice(device *models.Device) error {
	err := DB().Create(device).Error
	if err == nil {
		cache.InvalidateDevices()
		cache.InvalidateStatsCaches()
	}
	return err
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
	err := DB().Save(device).Error
	if err == nil {
		cache.InvalidateDevices()
		cache.InvalidateStatsCaches()
	}
	return err
}

func DeleteDevice(id uint) error {
	err := DB().Delete(&models.Device{}, id).Error
	if err == nil {
		cache.InvalidateDevices()
		cache.InvalidateStatsCaches()
	}
	return err
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
