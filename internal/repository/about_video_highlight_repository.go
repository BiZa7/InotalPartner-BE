package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type AboutVideoHighlightRepository struct {
	db *gorm.DB
}

func NewAboutVideoHighlightRepository(db *gorm.DB) *AboutVideoHighlightRepository {
	return &AboutVideoHighlightRepository{db: db}
}

// FindAll mengambil seluruh highlight video untuk dashboard admin (termasuk yang nonaktif).
// Diurutkan berdasarkan display_order ASC, lalu id ASC sebagai urutan sekunder.
func (r *AboutVideoHighlightRepository) FindAll() ([]model.AboutVideoHighlight, error) {
	var highlights []model.AboutVideoHighlight
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&highlights).Error
	if err != nil {
		return nil, err
	}
	if highlights == nil {
		highlights = []model.AboutVideoHighlight{}
	}
	return highlights, nil
}

// FindByID mengambil highlight video berdasarkan ID.
func (r *AboutVideoHighlightRepository) FindByID(id uint) (*model.AboutVideoHighlight, error) {
	var highlight model.AboutVideoHighlight
	err := r.db.First(&highlight, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &highlight, nil
}

// FindActivePublic mengambil hanya highlight video aktif untuk halaman publik.
// Diurutkan berdasarkan display_order ASC, lalu id ASC sebagai urutan sekunder.
func (r *AboutVideoHighlightRepository) FindActivePublic() ([]model.AboutVideoHighlight, error) {
	var highlights []model.AboutVideoHighlight
	err := r.db.
		Where("is_active = ?", true).
		Order("display_order ASC, id ASC").
		Find(&highlights).Error
	if err != nil {
		return nil, err
	}
	if highlights == nil {
		highlights = []model.AboutVideoHighlight{}
	}
	return highlights, nil
}

// Create membuat highlight baru.
func (r *AboutVideoHighlightRepository) Create(highlight *model.AboutVideoHighlight) error {
	return r.db.Create(highlight).Error
}

// Update memperbarui data highlight.
func (r *AboutVideoHighlightRepository) Update(highlight *model.AboutVideoHighlight) error {
	return r.db.Save(highlight).Error
}

// Delete menghapus highlight berdasarkan ID.
func (r *AboutVideoHighlightRepository) Delete(id uint) error {
	return r.db.Delete(&model.AboutVideoHighlight{}, id).Error
}

