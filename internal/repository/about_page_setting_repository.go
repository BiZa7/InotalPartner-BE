package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type AboutPageSettingRepository struct {
	db *gorm.DB
}

func NewAboutPageSettingRepository(db *gorm.DB) *AboutPageSettingRepository {
	return &AboutPageSettingRepository{db: db}
}

// FindActive mengambil konfigurasi judul halaman Tentang yang aktif (is_active = true).
func (r *AboutPageSettingRepository) FindActive() (*model.AboutPageSetting, error) {
	var setting model.AboutPageSetting
	err := r.db.Where("is_active = ?", true).Order("id DESC").First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// FindByID mengambil konfigurasi berdasarkan ID.
func (r *AboutPageSettingRepository) FindByID(id uint) (*model.AboutPageSetting, error) {
	var setting model.AboutPageSetting
	err := r.db.First(&setting, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// FindAll mengambil seluruh konfigurasi judul halaman Tentang untuk kebutuhan pengelolaan admin.
func (r *AboutPageSettingRepository) FindAll() ([]model.AboutPageSetting, error) {
	var settings []model.AboutPageSetting
	err := r.db.Order("id ASC").Find(&settings).Error
	if err != nil {
		return nil, err
	}
	if settings == nil {
		settings = []model.AboutPageSetting{}
	}
	return settings, nil
}

// Create membuat konfigurasi baru.
func (r *AboutPageSettingRepository) Create(setting *model.AboutPageSetting) error {
	return r.db.Create(setting).Error
}

// Update memperbarui data konfigurasi berdasarkan ID.
func (r *AboutPageSettingRepository) Update(setting *model.AboutPageSetting) error {
	return r.db.Save(setting).Error
}

// DeactivateAll menonaktifkan seluruh konfigurasi aktif yang ada di database.
func (r *AboutPageSettingRepository) DeactivateAll() error {
	return r.db.Model(&model.AboutPageSetting{}).Where("is_active = ?", true).Update("is_active", false).Error
}

// DeactivateAllExcept menonaktifkan seluruh konfigurasi aktif selain ID yang ditentukan.
func (r *AboutPageSettingRepository) DeactivateAllExcept(excludeID uint) error {
	query := r.db.Model(&model.AboutPageSetting{}).Where("is_active = ?", true)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	return query.Update("is_active", false).Error
}

// Count menghitung total record konfigurasi judul di database.
func (r *AboutPageSettingRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.AboutPageSetting{}).Count(&count).Error
	return count, err
}
