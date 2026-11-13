package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type PartnershipServiceRepository struct {
	db *gorm.DB
}

func NewPartnershipServiceRepository(db *gorm.DB) *PartnershipServiceRepository {
	return &PartnershipServiceRepository{db: db}
}

// FindAll mengambil seluruh data kartu layanan untuk dashboard admin (termasuk yang nonaktif).
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *PartnershipServiceRepository) FindAll() ([]model.PartnershipService, error) {
	var services []model.PartnershipService
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&services).Error
	if err != nil {
		return nil, err
	}
	if services == nil {
		services = []model.PartnershipService{}
	}
	return services, nil
}

// FindByID mengambil satu data kartu layanan berdasarkan ID.
// Mengembalikan nil, nil jika data tidak ditemukan.
func (r *PartnershipServiceRepository) FindByID(id uint) (*model.PartnershipService, error) {
	var service model.PartnershipService
	err := r.db.First(&service, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// FindActive mengambil hanya kartu layanan yang aktif untuk ditampilkan pada halaman publik.
// Diurutkan berdasarkan display_order ASC, lalu id ASC.
func (r *PartnershipServiceRepository) FindActive() ([]model.PartnershipService, error) {
	var services []model.PartnershipService
	err := r.db.
		Where("is_active = ?", true).
		Order("display_order ASC, id ASC").
		Find(&services).Error
	if err != nil {
		return nil, err
	}
	if services == nil {
		services = []model.PartnershipService{}
	}
	return services, nil
}

// Create menyimpan data kartu layanan baru ke database.
func (r *PartnershipServiceRepository) Create(service *model.PartnershipService) error {
	return r.db.Create(service).Error
}

// Update memperbarui data kartu layanan di database.
func (r *PartnershipServiceRepository) Update(service *model.PartnershipService) error {
	return r.db.Save(service).Error
}

// Delete menghapus data kartu layanan berdasarkan ID secara permanen.
func (r *PartnershipServiceRepository) Delete(id uint) error {
	return r.db.Delete(&model.PartnershipService{}, id).Error
}

