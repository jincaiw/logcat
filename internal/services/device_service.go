package services

import (
	"errors"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// DeviceService handles device management logic
type DeviceService struct{}

// NewDeviceService creates a new DeviceService
func NewDeviceService() *DeviceService {
	return &DeviceService{}
}

// ListDevices returns all devices with optional filters
func (s *DeviceService) ListDevices(groupID *uint, enabled *bool, keyword string) ([]models.Device, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	query := db.Model(&models.Device{}).Preload("Group").Preload("Template").Preload("ParseTemplate")

	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	}
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR ip_address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var devices []models.Device
	if err := query.Order("id DESC").Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// GetDevice returns a single device by ID
func (s *DeviceService) GetDevice(id uint) (*models.Device, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var device models.Device
	if err := db.Preload("Group").Preload("Template").Preload("ParseTemplate").First(&device, id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// GetDeviceByIP returns a device by IP address
func (s *DeviceService) GetDeviceByIP(ip string) (*models.Device, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var device models.Device
	if err := db.Where("ip_address = ?", ip).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// CreateDevice creates a new device
func (s *DeviceService) CreateDevice(device *models.Device) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	// Check for duplicate IP
	var count int64
	db.Model(&models.Device{}).Where("ip_address = ?", device.IPAddress).Count(&count)
	if count > 0 {
		return errors.New("device with this IP address already exists")
	}

	return db.Create(device).Error
}

// UpdateDevice updates an existing device
func (s *DeviceService) UpdateDevice(id uint, updates map[string]interface{}) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	var device models.Device
	if err := db.First(&device, id).Error; err != nil {
		return err
	}

	// If updating IP, check for duplicates
	if ip, ok := updates["ip_address"].(string); ok && ip != device.IPAddress {
		var count int64
		db.Model(&models.Device{}).Where("ip_address = ? AND id != ?", ip, id).Count(&count)
		if count > 0 {
			return errors.New("device with this IP address already exists")
		}
	}

	return db.Model(&device).Updates(updates).Error
}

// DeleteDevice deletes a device
func (s *DeviceService) DeleteDevice(id uint) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Delete(&models.Device{}, id).Error
}

// --- DeviceGroup operations ---

// ListDeviceGroups returns all device groups
func (s *DeviceService) ListDeviceGroups() ([]models.DeviceGroup, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var groups []models.DeviceGroup
	if err := db.Order("sort_order ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetDeviceGroup returns a single device group by ID
func (s *DeviceService) GetDeviceGroup(id uint) (*models.DeviceGroup, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var group models.DeviceGroup
	if err := db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// CreateDeviceGroup creates a new device group
func (s *DeviceService) CreateDeviceGroup(group *models.DeviceGroup) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Create(group).Error
}

// UpdateDeviceGroup updates an existing device group
func (s *DeviceService) UpdateDeviceGroup(id uint, updates map[string]interface{}) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	var group models.DeviceGroup
	if err := db.First(&group, id).Error; err != nil {
		return err
	}

	return db.Model(&group).Updates(updates).Error
}

// DeleteDeviceGroup deletes a device group
func (s *DeviceService) DeleteDeviceGroup(id uint) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Delete(&models.DeviceGroup{}, id).Error
}