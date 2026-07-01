package repository

import (
	"syslog-alert/internal/models"
	"syslog-alert/internal/service/cache"
)

// ---- 解析模板 ----

func CreateParseTemplate(template *models.ParseTemplate) error {
	err := DB().Create(template).Error
	if err == nil {
		cache.InvalidateParseTemplates()
		cache.InvalidateStatsCaches()
	}
	return err
}

func GetParseTemplates() []models.ParseTemplate {
	var templates []models.ParseTemplate
	DB().Find(&templates)
	return templates
}

func GetParseTemplateCount() int64 {
	var count int64
	DB().Model(&models.ParseTemplate{}).Count(&count)
	return count
}

func GetParseTemplateByID(id uint) (*models.ParseTemplate, error) {
	var template models.ParseTemplate
	err := DB().First(&template, id).Error
	return &template, err
}

func UpdateParseTemplate(template *models.ParseTemplate) error {
	err := DB().Save(template).Error
	if err == nil {
		cache.InvalidateParseTemplates()
		cache.InvalidateStatsCaches()
	}
	return err
}

func DeleteParseTemplate(id uint) error {
	err := DB().Delete(&models.ParseTemplate{}, id).Error
	if err == nil {
		cache.InvalidateParseTemplates()
		cache.InvalidateStatsCaches()
	}
	return err
}

// ---- 输出模板 ----

func CreateOutputTemplate(template *models.OutputTemplate) error {
	err := DB().Create(template).Error
	if err == nil {
		cache.InvalidateOutputTemplates()
		cache.InvalidateStatsCaches()
	}
	return err
}

func GetOutputTemplates() []models.OutputTemplate {
	var templates []models.OutputTemplate
	DB().Find(&templates)
	return templates
}

func GetOutputTemplateByID(id uint) (*models.OutputTemplate, error) {
	var template models.OutputTemplate
	err := DB().First(&template, id).Error
	return &template, err
}

func GetOutputTemplateByPlatform(platform string) (*models.OutputTemplate, error) {
	var template models.OutputTemplate
	err := DB().Where("platform = ? AND is_active = ?", platform, true).First(&template).Error
	return &template, err
}

func UpdateOutputTemplate(template *models.OutputTemplate) error {
	err := DB().Save(template).Error
	if err == nil {
		cache.InvalidateOutputTemplates()
		cache.InvalidateStatsCaches()
	}
	return err
}

func DeleteOutputTemplate(id uint) error {
	err := DB().Delete(&models.OutputTemplate{}, id).Error
	if err == nil {
		cache.InvalidateOutputTemplates()
		cache.InvalidateStatsCaches()
	}
	return err
}

// ---- 字段映射文档 ----

func CreateFieldMappingDoc(doc *models.FieldMappingDoc) error {
	return DB().Create(doc).Error
}

func GetFieldMappingDocs() []models.FieldMappingDoc {
	var docs []models.FieldMappingDoc
	DB().Order("device_type ASC").Find(&docs)
	return docs
}

func GetFieldMappingDocByID(id uint) (*models.FieldMappingDoc, error) {
	var doc models.FieldMappingDoc
	err := DB().First(&doc, id).Error
	return &doc, err
}

func GetFieldMappingDocByDeviceType(deviceType string) (*models.FieldMappingDoc, error) {
	var doc models.FieldMappingDoc
	err := DB().Where("device_type = ?", deviceType).First(&doc).Error
	return &doc, err
}

func GetFieldMappingDocByName(name string) (*models.FieldMappingDoc, error) {
	var doc models.FieldMappingDoc
	err := DB().Where("name = ?", name).First(&doc).Error
	return &doc, err
}

func UpdateFieldMappingDoc(doc *models.FieldMappingDoc) error {
	return DB().Save(doc).Error
}

func DeleteFieldMappingDoc(id uint) error {
	return DB().Delete(&models.FieldMappingDoc{}, id).Error
}
