package services

import (
	"errors"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/pkg/response"
)

// DeviceService handles device management logic
type DeviceService struct{}

// NewDeviceService creates a new DeviceService
func NewDeviceService() *DeviceService {
	return &DeviceService{}
}

// ListDevices returns devices with optional filters and pagination
func (s *DeviceService) ListDevices(page, pageSize int, groupID *uint, enabled *bool, keyword string) ([]models.Device, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not available")
	}

	query := db.Model(&models.Device{}).Preload("Group").Preload("Template").Preload("ParseTemplate")

	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	}
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR ip_address LIKE ?", "%"+response.EscapeLike(keyword)+"%", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var devices []models.Device
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&devices).Error; err != nil {
		return nil, 0, err
	}
	return devices, total, nil
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

	allowedFields := map[string]bool{
		"name": true, "ip_address": true, "group_id": true,
		"template_id": true, "parse_template_id": true,
		"device_type": true, "description": true, "enabled": true,
	}
	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			filtered[k] = v
		}
	}
	if len(filtered) == 0 {
		return errors.New("no valid fields to update")
	}

	var device models.Device
	if err := db.First(&device, id).Error; err != nil {
		return err
	}

	// If updating IP, check for duplicates
	if ip, ok := filtered["ip_address"].(string); ok && ip != device.IPAddress {
		var count int64
		db.Model(&models.Device{}).Where("ip_address = ? AND id != ?", ip, id).Count(&count)
		if count > 0 {
			return errors.New("device with this IP address already exists")
		}
	}

	return db.Model(&device).Updates(filtered).Error
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

// DeviceGroupListItem represents a device group with device count for list response
type DeviceGroupListItem struct {
	models.DeviceGroup
	DeviceCount int64 `json:"deviceCount"`
}

// ListDeviceGroups returns device groups with pagination and device count
func (s *DeviceService) ListDeviceGroups(page, pageSize int, keyword string) ([]DeviceGroupListItem, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not available")
	}

	query := db.Model(&models.DeviceGroup{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+response.EscapeLike(keyword)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var groups []models.DeviceGroup
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sort_order ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	items := make([]DeviceGroupListItem, len(groups))
	for i, g := range groups {
		var count int64
		if err := db.Model(&models.Device{}).Where("group_id = ?", g.ID).Count(&count).Error; err == nil {
			items[i] = DeviceGroupListItem{DeviceGroup: g, DeviceCount: count}
		} else {
			items[i] = DeviceGroupListItem{DeviceGroup: g, DeviceCount: 0}
		}
	}
	return items, total, nil
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

	allowedFields := map[string]bool{
		"name": true, "description": true, "color": true, "sort_order": true,
	}
	filtered := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			filtered[k] = v
		}
	}
	if len(filtered) == 0 {
		return errors.New("no valid fields to update")
	}

	var group models.DeviceGroup
	if err := db.First(&group, id).Error; err != nil {
		return err
	}

	return db.Model(&group).Updates(filtered).Error
}

// DeleteDeviceGroup deletes a device group
func (s *DeviceService) DeleteDeviceGroup(id uint) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Delete(&models.DeviceGroup{}, id).Error
}