package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type PartnershipImageRepository struct {
	db *gorm.DB
}

func NewPartnershipImageRepository(db *gorm.DB) *PartnershipImageRepository {
	return &PartnershipImageRepository{db: db}
}

// FindAll mengambil seluruh data gambar kolase untuk dashboard admin.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *PartnershipImageRepository) FindAll() ([]model.PartnershipImage, error) {
	var images []model.PartnershipImage
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	if images == nil {
		images = []model.PartnershipImage{}
	}
	return images, nil
}

// FindByID mengambil satu data gambar berdasarkan ID.
// Mengembalikan nil, nil jika data tidak ditemukan.
func (r *PartnershipImageRepository) FindByID(id uint) (*model.PartnershipImage, error) {
	var image model.PartnershipImage
	err := r.db.First(&image, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// FindActive mengambil hanya gambar-gambar yang aktif untuk ditampilkan di publik.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *PartnershipImageRepository) FindActive() ([]model.PartnershipImage, error) {
	var images []model.PartnershipImage
	err := r.db.
		Where("is_active = ?", true).
		Order("display_order ASC, id ASC").
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	if images == nil {
		images = []model.PartnershipImage{}
	}
	return images, nil
}

// Create menyimpan data gambar baru ke database.
func (r *PartnershipImageRepository) Create(image *model.PartnershipImage) error {
	return r.db.Create(image).Error
}

// Update memperbarui data gambar di database.
func (r *PartnershipImageRepository) Update(image *model.PartnershipImage) error {
	return r.db.Save(image).Error
}

// Delete menghapus data gambar berdasarkan ID secara permanen.
func (r *PartnershipImageRepository) Delete(id uint) error {
	return r.db.Delete(&model.PartnershipImage{}, id).Error
}

