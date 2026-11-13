package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type AboutIntroSettingRepository struct {
	db *gorm.DB
}

func NewAboutIntroSettingRepository(db *gorm.DB) *AboutIntroSettingRepository {
	return &AboutIntroSettingRepository{db: db}
}

// FindAll mengambil seluruh konfigurasi section Hero/Perkenalan untuk dashboard admin.
// Diurutkan berdasarkan id DESC.
func (r *AboutIntroSettingRepository) FindAll() ([]model.AboutIntroSetting, error) {
	var settings []model.AboutIntroSetting
	err := r.db.
		Order("id DESC").
		Find(&settings).Error
	if err != nil {
		return nil, err
	}
	if settings == nil {
		settings = []model.AboutIntroSetting{}
	}
	return settings, nil
}

// FindByID mengambil konfigurasi berdasarkan ID beserta relasi metrics dan points.
// Mengembalikan nil, nil jika data tidak ditemukan.
func (r *AboutIntroSettingRepository) FindByID(id uint) (*model.AboutIntroSetting, error) {
	var setting model.AboutIntroSetting
	err := r.db.
		Preload("Metrics", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC, id ASC")
		}).
		Preload("Points", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC, id ASC")
		}).
		First(&setting, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// FindActive mengambil satu konfigurasi aktif untuk halaman publik.
// Preload relasi metrics dan points yang berstatus aktif saja (is_active = true)
// dan diurutkan berdasarkan display_order ASC, id ASC.
func (r *AboutIntroSettingRepository) FindActive() (*model.AboutIntroSetting, error) {
	var setting model.AboutIntroSetting
	err := r.db.
		Where("is_active = ?", true).
		Preload("Metrics", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("display_order ASC, id ASC")
		}).
		Preload("Points", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("display_order ASC, id ASC")
		}).
		Order("id DESC").
		First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// Create menyimpan konfigurasi baru.
// Jika is_active = true, deaktifkan konfigurasi aktif lainnya dalam satu transaksi database.
func (r *AboutIntroSettingRepository) Create(setting *model.AboutIntroSetting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if setting.IsActive {
			if err := tx.Model(&model.AboutIntroSetting{}).
				Where("is_active = ?", true).
				Update("is_active", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(setting).Error
	})
}

// Update memperbarui data konfigurasi.
// Jika is_active diaktifkan (true), deaktifkan konfigurasi lain dalam satu transaksi database.
func (r *AboutIntroSettingRepository) Update(setting *model.AboutIntroSetting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if setting.IsActive {
			if err := tx.Model(&model.AboutIntroSetting{}).
				Where("is_active = ? AND id <> ?", true, setting.ID).
				Update("is_active", false).Error; err != nil {
				return err
			}
		}
		return tx.Save(setting).Error
	})
}

// Delete menghapus konfigurasi berdasarkan ID.
// Relasi metrics dan points akan terhapus via OnDelete:CASCADE di database level.
func (r *AboutIntroSettingRepository) Delete(id uint) error {
	return r.db.Delete(&model.AboutIntroSetting{}, id).Error
}

