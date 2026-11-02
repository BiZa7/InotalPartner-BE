package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type AboutFAQRepository struct {
	db *gorm.DB
}

func NewAboutFAQRepository(db *gorm.DB) *AboutFAQRepository {
	return &AboutFAQRepository{db: db}
}

// FindAll mengambil seluruh FAQ untuk kebutuhan dashboard admin (termasuk yang nonaktif).
// Diurutkan berdasarkan display_order ASC, lalu id ASC sebagai urutan sekunder.
func (r *AboutFAQRepository) FindAll() ([]model.AboutFAQ, error) {
	var faqs []model.AboutFAQ
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	if faqs == nil {
		faqs = []model.AboutFAQ{}
	}
	return faqs, nil
}

// FindByID mengambil FAQ berdasarkan ID.
func (r *AboutFAQRepository) FindByID(id uint) (*model.AboutFAQ, error) {
	var faq model.AboutFAQ
	err := r.db.First(&faq, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &faq, nil
}

// FindActivePublic mengambil hanya FAQ aktif untuk halaman publik.
// Diurutkan berdasarkan display_order ASC, lalu id ASC sebagai urutan sekunder.
func (r *AboutFAQRepository) FindActivePublic() ([]model.AboutFAQ, error) {
	var faqs []model.AboutFAQ
	err := r.db.
		Where("is_active = ?", true).
		Order("display_order ASC, id ASC").
		Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	if faqs == nil {
		faqs = []model.AboutFAQ{}
	}
	return faqs, nil
}

// Create membuat FAQ baru.
func (r *AboutFAQRepository) Create(faq *model.AboutFAQ) error {
	return r.db.Create(faq).Error
}

// Update memperbarui data FAQ.
func (r *AboutFAQRepository) Update(faq *model.AboutFAQ) error {
	return r.db.Save(faq).Error
}

// Delete menghapus FAQ berdasarkan ID.
func (r *AboutFAQRepository) Delete(id uint) error {
	return r.db.Delete(&model.AboutFAQ{}, id).Error
}

