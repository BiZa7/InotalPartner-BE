package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type AboutStatisticRepository struct {
	db *gorm.DB
}

func NewAboutStatisticRepository(db *gorm.DB) *AboutStatisticRepository {
	return &AboutStatisticRepository{db: db}
}

// FindAll mengambil seluruh data statistik untuk dashboard admin (termasuk yang nonaktif).
// Diurutkan berdasarkan display_order ASC, lalu id ASC sebagai urutan sekunder.
func (r *AboutStatisticRepository) FindAll() ([]model.AboutStatistic, error) {
	var statistics []model.AboutStatistic
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&statistics).Error
	if err != nil {
		return nil, err
	}
	if statistics == nil {
		statistics = []model.AboutStatistic{}
	}
	return statistics, nil
}

// FindByID mengambil data statistik berdasarkan ID.
// Mengembalikan error jika data tidak ditemukan.
func (r *AboutStatisticRepository) FindByID(id uint) (*model.AboutStatistic, error) {
	var statistic model.AboutStatistic
	err := r.db.First(&statistic, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("statistik tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	return &statistic, nil
}

// FindAllActive mengambil hanya data statistik aktif (is_active = true) untuk halaman publik.
// Diurutkan berdasarkan display_order ASC, lalu id ASC sebagai urutan sekunder.
func (r *AboutStatisticRepository) FindAllActive() ([]model.AboutStatistic, error) {
	var statistics []model.AboutStatistic
	err := r.db.
		Where("is_active = ?", true).
		Order("display_order ASC, id ASC").
		Find(&statistics).Error
	if err != nil {
		return nil, err
	}
	if statistics == nil {
		statistics = []model.AboutStatistic{}
	}
	return statistics, nil
}

// Create membuat data statistik baru.
func (r *AboutStatisticRepository) Create(statistic *model.AboutStatistic) error {
	return r.db.Create(statistic).Error
}

// Update memperbarui data statistik berdasarkan ID yang valid.
// Menggunakan r.db.Save(statistic) untuk memastikan seluruh field termasuk is_active (true/false) tersimpan.
func (r *AboutStatisticRepository) Update(statistic *model.AboutStatistic) error {
	return r.db.Save(statistic).Error
}

// Delete menghapus data statistik berdasarkan ID yang valid.
func (r *AboutStatisticRepository) Delete(id uint) error {
	return r.db.Delete(&model.AboutStatistic{}, id).Error
}

