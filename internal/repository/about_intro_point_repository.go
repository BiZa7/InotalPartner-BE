package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type AboutIntroPointRepository struct {
	db *gorm.DB
}

func NewAboutIntroPointRepository(db *gorm.DB) *AboutIntroPointRepository {
	return &AboutIntroPointRepository{db: db}
}

// FindAll mengambil seluruh data point untuk dashboard admin.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *AboutIntroPointRepository) FindAll() ([]model.AboutIntroPoint, error) {
	var points []model.AboutIntroPoint
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&points).Error
	if err != nil {
		return nil, err
	}
	if points == nil {
		points = []model.AboutIntroPoint{}
	}
	return points, nil
}

// FindByID mengambil satu data point berdasarkan ID.
// Mengembalikan nil, nil jika tidak ditemukan.
func (r *AboutIntroPointRepository) FindByID(id uint) (*model.AboutIntroPoint, error) {
	var point model.AboutIntroPoint
	err := r.db.First(&point, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &point, nil
}

// FindBySettingID mengambil seluruh point berdasarkan SettingID.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *AboutIntroPointRepository) FindBySettingID(settingID uint) ([]model.AboutIntroPoint, error) {
	var points []model.AboutIntroPoint
	err := r.db.
		Where("setting_id = ?", settingID).
		Order("display_order ASC, id ASC").
		Find(&points).Error
	if err != nil {
		return nil, err
	}
	if points == nil {
		points = []model.AboutIntroPoint{}
	}
	return points, nil
}

// FindBySettingIDAndType mengambil seluruh point berdasarkan SettingID dan SectionType ('mission' atau 'vision').
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *AboutIntroPointRepository) FindBySettingIDAndType(settingID uint, sectionType string) ([]model.AboutIntroPoint, error) {
	var points []model.AboutIntroPoint
	err := r.db.
		Where("setting_id = ? AND section_type = ?", settingID, sectionType).
		Order("display_order ASC, id ASC").
		Find(&points).Error
	if err != nil {
		return nil, err
	}
	if points == nil {
		points = []model.AboutIntroPoint{}
	}
	return points, nil
}

// FindActiveBySettingID mengambil hanya point yang aktif berdasarkan SettingID.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *AboutIntroPointRepository) FindActiveBySettingID(settingID uint) ([]model.AboutIntroPoint, error) {
	var points []model.AboutIntroPoint
	err := r.db.
		Where("setting_id = ? AND is_active = ?", settingID, true).
		Order("display_order ASC, id ASC").
		Find(&points).Error
	if err != nil {
		return nil, err
	}
	if points == nil {
		points = []model.AboutIntroPoint{}
	}
	return points, nil
}

// Create menyimpan data point baru ke database.
func (r *AboutIntroPointRepository) Create(point *model.AboutIntroPoint) error {
	return r.db.Create(point).Error
}

// Update memperbarui data point di database.
func (r *AboutIntroPointRepository) Update(point *model.AboutIntroPoint) error {
	return r.db.Save(point).Error
}

// Delete menghapus data point berdasarkan ID secara permanen.
func (r *AboutIntroPointRepository) Delete(id uint) error {
	return r.db.Delete(&model.AboutIntroPoint{}, id).Error
}

