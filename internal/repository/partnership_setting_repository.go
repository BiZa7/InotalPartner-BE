package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type PartnershipSettingRepository struct {
	db *gorm.DB
}

func NewPartnershipSettingRepository(db *gorm.DB) *PartnershipSettingRepository {
	return &PartnershipSettingRepository{db: db}
}

// FindAll mengambil seluruh konfigurasi section Partnership untuk dashboard admin.
// Diurutkan berdasarkan id DESC agar record terbaru berada di atas.
func (r *PartnershipSettingRepository) FindAll() ([]model.PartnershipSetting, error) {
	var settings []model.PartnershipSetting
	err := r.db.
		Order("id DESC").
		Find(&settings).Error
	if err != nil {
		return nil, err
	}
	if settings == nil {
		settings = []model.PartnershipSetting{}
	}
	return settings, nil
}

// FindByID mengambil konfigurasi berdasarkan ID.
// Mengembalikan nil, nil jika data tidak ditemukan.
func (r *PartnershipSettingRepository) FindByID(id uint) (*model.PartnershipSetting, error) {
	var setting model.PartnershipSetting
	err := r.db.First(&setting, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// FindActive mengambil satu konfigurasi aktif untuk kebutuhan publik di halaman Beranda.
// Mengutamakan ID terbaru jika ada record aktif.
func (r *PartnershipSettingRepository) FindActive() (*model.PartnershipSetting, error) {
	var setting model.PartnershipSetting
	err := r.db.
		Where("is_active = ?", true).
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
// Jika is_active = true, deaktifkan konfigurasi aktif sebelumnya dalam satu transaksi database.
func (r *PartnershipSettingRepository) Create(setting *model.PartnershipSetting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if setting.IsActive {
			if err := tx.Model(&model.PartnershipSetting{}).
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
func (r *PartnershipSettingRepository) Update(setting *model.PartnershipSetting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if setting.IsActive {
			if err := tx.Model(&model.PartnershipSetting{}).
				Where("is_active = ? AND id <> ?", true, setting.ID).
				Update("is_active", false).Error; err != nil {
				return err
			}
		}
		return tx.Save(setting).Error
	})
}

// Delete menghapus konfigurasi berdasarkan ID secara permanen.
func (r *PartnershipSettingRepository) Delete(id uint) error {
	return r.db.Delete(&model.PartnershipSetting{}, id).Error
}

