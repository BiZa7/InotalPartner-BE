package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type AboutVideoSettingRepository struct {
	db *gorm.DB
}

func NewAboutVideoSettingRepository(db *gorm.DB) *AboutVideoSettingRepository {
	return &AboutVideoSettingRepository{db: db}
}

// FindFirst mengambil record pengaturan section video profil pertama.
func (r *AboutVideoSettingRepository) FindFirst() (*model.AboutVideoSetting, error) {
	var setting model.AboutVideoSetting
	err := r.db.Order("id ASC").First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// FindActive mengambil record pengaturan section video profil yang aktif (is_active = true).
func (r *AboutVideoSettingRepository) FindActive() (*model.AboutVideoSetting, error) {
	var setting model.AboutVideoSetting
	err := r.db.Where("is_active = ?", true).Order("id ASC").First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// Create membuat record pengaturan section baru.
func (r *AboutVideoSettingRepository) Create(setting *model.AboutVideoSetting) error {
	return r.db.Create(setting).Error
}

// Update memperbarui data pengaturan section.
func (r *AboutVideoSettingRepository) Update(setting *model.AboutVideoSetting) error {
	return r.db.Save(setting).Error
}

