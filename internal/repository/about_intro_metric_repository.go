package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type AboutIntroMetricRepository struct {
	db *gorm.DB
}

func NewAboutIntroMetricRepository(db *gorm.DB) *AboutIntroMetricRepository {
	return &AboutIntroMetricRepository{db: db}
}

// FindAll mengambil seluruh data metric untuk dashboard admin.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *AboutIntroMetricRepository) FindAll() ([]model.AboutIntroMetric, error) {
	var metrics []model.AboutIntroMetric
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&metrics).Error
	if err != nil {
		return nil, err
	}
	if metrics == nil {
		metrics = []model.AboutIntroMetric{}
	}
	return metrics, nil
}

// FindByID mengambil satu data metric berdasarkan ID.
// Mengembalikan nil, nil jika tidak ditemukan.
func (r *AboutIntroMetricRepository) FindByID(id uint) (*model.AboutIntroMetric, error) {
	var metric model.AboutIntroMetric
	err := r.db.First(&metric, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

// FindBySettingID mengambil seluruh metric berdasarkan SettingID.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *AboutIntroMetricRepository) FindBySettingID(settingID uint) ([]model.AboutIntroMetric, error) {
	var metrics []model.AboutIntroMetric
	err := r.db.
		Where("setting_id = ?", settingID).
		Order("display_order ASC, id ASC").
		Find(&metrics).Error
	if err != nil {
		return nil, err
	}
	if metrics == nil {
		metrics = []model.AboutIntroMetric{}
	}
	return metrics, nil
}

// FindActiveBySettingID mengambil hanya metric yang aktif berdasarkan SettingID.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *AboutIntroMetricRepository) FindActiveBySettingID(settingID uint) ([]model.AboutIntroMetric, error) {
	var metrics []model.AboutIntroMetric
	err := r.db.
		Where("setting_id = ? AND is_active = ?", settingID, true).
		Order("display_order ASC, id ASC").
		Find(&metrics).Error
	if err != nil {
		return nil, err
	}
	if metrics == nil {
		metrics = []model.AboutIntroMetric{}
	}
	return metrics, nil
}

// Create menyimpan data metric baru ke database.
func (r *AboutIntroMetricRepository) Create(metric *model.AboutIntroMetric) error {
	return r.db.Create(metric).Error
}

// Update memperbarui data metric di database.
func (r *AboutIntroMetricRepository) Update(metric *model.AboutIntroMetric) error {
	return r.db.Save(metric).Error
}

// Delete menghapus data metric berdasarkan ID secara permanen.
func (r *AboutIntroMetricRepository) Delete(id uint) error {
	return r.db.Delete(&model.AboutIntroMetric{}, id).Error
}

