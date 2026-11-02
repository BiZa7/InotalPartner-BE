package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type PageBannerRepository struct {
	db *gorm.DB
}

func NewPageBannerRepository(db *gorm.DB) *PageBannerRepository {
	return &PageBannerRepository{db: db}
}

// FindByPageType mengambil data banner berdasarkan jenis halaman (page_type).
func (r *PageBannerRepository) FindByPageType(pageType string) (*model.PageBanner, error) {
	var banner model.PageBanner
	err := r.db.Where("page_type = ?", pageType).First(&banner).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// FindActiveByPageType mengambil banner yang aktif saja (is_active = true) berdasarkan page_type.
func (r *PageBannerRepository) FindActiveByPageType(pageType string) (*model.PageBanner, error) {
	var banner model.PageBanner
	err := r.db.Where("page_type = ? AND is_active = ?", pageType, true).First(&banner).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// Create menyimpan banner baru ke database.
func (r *PageBannerRepository) Create(banner *model.PageBanner) error {
	return r.db.Create(banner).Error
}

// Update memperbarui seluruh field banner menggunakan r.db.Save()
// agar nilai is_active (termasuk false) tersimpan dengan benar.
func (r *PageBannerRepository) Update(banner *model.PageBanner) error {
	return r.db.Save(banner).Error
}

// Delete menghapus banner berdasarkan ID.
func (r *PageBannerRepository) Delete(id uint) error {
	return r.db.Delete(&model.PageBanner{}, id).Error
}

