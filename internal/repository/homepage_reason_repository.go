package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type HomepageReasonRepository struct {
	db *gorm.DB
}

func NewHomepageReasonRepository(db *gorm.DB) *HomepageReasonRepository {
	return &HomepageReasonRepository{db: db}
}

// ─── Settings ─────────────────────────────────────────────────────────────────

// FindFirstSetting mengambil record pengaturan section alasan pertama.
func (r *HomepageReasonRepository) FindFirstSetting() (*model.HomepageReasonSetting, error) {
	var setting model.HomepageReasonSetting
	err := r.db.Order("id ASC").First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// FindActiveSetting mengambil record pengaturan section alasan yang aktif (is_active = true).
func (r *HomepageReasonRepository) FindActiveSetting() (*model.HomepageReasonSetting, error) {
	var setting model.HomepageReasonSetting
	err := r.db.Where("is_active = ?", true).Order("id ASC").First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// CreateSetting membuat record baru untuk pengaturan section alasan.
func (r *HomepageReasonRepository) CreateSetting(setting *model.HomepageReasonSetting) error {
	return r.db.Create(setting).Error
}

// UpdateSetting memperbarui data pengaturan section alasan yang sudah ada.
func (r *HomepageReasonRepository) UpdateSetting(setting *model.HomepageReasonSetting) error {
	return r.db.Save(setting).Error
}

// ─── Reason Cards ─────────────────────────────────────────────────────────────

// FindAll mengambil seluruh kartu alasan untuk dashboard admin (termasuk yang nonaktif).
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *HomepageReasonRepository) FindAll() ([]model.HomepageReason, error) {
	var reasons []model.HomepageReason
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&reasons).Error
	if err != nil {
		return nil, err
	}
	if reasons == nil {
		reasons = []model.HomepageReason{}
	}
	return reasons, nil
}

// FindByID mengambil satu data kartu alasan berdasarkan ID.
func (r *HomepageReasonRepository) FindByID(id uint) (*model.HomepageReason, error) {
	var reason model.HomepageReason
	err := r.db.First(&reason, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("kartu alasan tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	return &reason, nil
}

// FindAllActive mengambil hanya kartu alasan aktif (is_active = true) untuk halaman publik.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *HomepageReasonRepository) FindAllActive() ([]model.HomepageReason, error) {
	var reasons []model.HomepageReason
	err := r.db.
		Where("is_active = ?", true).
		Order("display_order ASC, id ASC").
		Find(&reasons).Error
	if err != nil {
		return nil, err
	}
	if reasons == nil {
		reasons = []model.HomepageReason{}
	}
	return reasons, nil
}

// Create membuat kartu alasan baru.
func (r *HomepageReasonRepository) Create(reason *model.HomepageReason) error {
	return r.db.Create(reason).Error
}

// Update memperbarui data kartu alasan.
func (r *HomepageReasonRepository) Update(reason *model.HomepageReason) error {
	return r.db.Save(reason).Error
}

// Delete menghapus kartu alasan secara permanen berdasarkan ID.
func (r *HomepageReasonRepository) Delete(id uint) error {
	result := r.db.Delete(&model.HomepageReason{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("kartu alasan tidak ditemukan")
	}
	return nil
}

