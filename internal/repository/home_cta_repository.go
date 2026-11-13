package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type HomeCTARepository struct {
	db *gorm.DB
}

func NewHomeCTARepository(db *gorm.DB) *HomeCTARepository {
	return &HomeCTARepository{db: db}
}

// FindAll mengambil seluruh data CTA untuk dashboard admin (termasuk yang nonaktif).
// Diurutkan berdasarkan id DESC agar record terbaru berada di atas.
func (r *HomeCTARepository) FindAll() ([]model.HomeCTA, error) {
	var ctas []model.HomeCTA
	err := r.db.
		Order("id DESC").
		Find(&ctas).Error
	if err != nil {
		return nil, err
	}
	if ctas == nil {
		ctas = []model.HomeCTA{}
	}
	return ctas, nil
}

// FindByID mengambil data CTA berdasarkan ID.
// Mengembalikan nil, nil jika data tidak ditemukan.
func (r *HomeCTARepository) FindByID(id uint) (*model.HomeCTA, error) {
	var cta model.HomeCTA
	err := r.db.First(&cta, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cta, nil
}

// FindActivePublic mengambil satu CTA yang aktif untuk ditampilkan pada halaman publik.
// Mengutamakan ID terbaru jika ada record aktif.
func (r *HomeCTARepository) FindActivePublic() (*model.HomeCTA, error) {
	var cta model.HomeCTA
	err := r.db.
		Where("is_active = ?", true).
		Order("id DESC").
		First(&cta).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cta, nil
}

// DeactivateAllExcept menonaktifkan seluruh CTA selain ID tertentu (atau semua jika exceptID == 0).
func (r *HomeCTARepository) DeactivateAllExcept(exceptID uint) error {
	query := r.db.Model(&model.HomeCTA{}).Where("is_active = ?", true)
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}
	return query.Update("is_active", false).Error
}

// Create menyimpan data CTA baru. Jika is_active = true, deaktifkan yang lain dalam transaksi.
func (r *HomeCTARepository) Create(cta *model.HomeCTA) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if cta.IsActive {
			if err := tx.Model(&model.HomeCTA{}).Where("is_active = ?", true).Update("is_active", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(cta).Error
	})
}

// Update memperbarui data CTA. Jika is_active diubah menjadi true, deaktifkan yang lain dalam transaksi.
func (r *HomeCTARepository) Update(cta *model.HomeCTA) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if cta.IsActive {
			if err := tx.Model(&model.HomeCTA{}).Where("is_active = ? AND id <> ?", true, cta.ID).Update("is_active", false).Error; err != nil {
				return err
			}
		}
		return tx.Save(cta).Error
	})
}

// Delete menghapus data CTA berdasarkan ID secara permanen.
func (r *HomeCTARepository) Delete(id uint) error {
	result := r.db.Delete(&model.HomeCTA{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("data cta tidak ditemukan")
	}
	return nil
}

