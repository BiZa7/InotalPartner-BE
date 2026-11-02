package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type TestimonialRepository struct {
	db *gorm.DB
}

func NewTestimonialRepository(db *gorm.DB) *TestimonialRepository {
	return &TestimonialRepository{db: db}
}

// Create menyimpan testimoni baru.
func (r *TestimonialRepository) Create(testimonial *model.Testimonial) error {
	return r.db.Create(testimonial).Error
}

// FindAll mengambil seluruh testimoni, termasuk yang nonaktif, diurutkan berdasarkan display_order ASC, id ASC.
func (r *TestimonialRepository) FindAll() ([]model.Testimonial, error) {
	var testimonials []model.Testimonial
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&testimonials).Error
	if err != nil {
		return nil, err
	}
	if testimonials == nil {
		testimonials = []model.Testimonial{}
	}
	return testimonials, nil
}

// FindActive mengambil hanya testimoni dengan is_active = true, diurutkan berdasarkan display_order ASC, id ASC.
func (r *TestimonialRepository) FindActive() ([]model.Testimonial, error) {
	var testimonials []model.Testimonial
	err := r.db.
		Where("is_active = ?", true).
		Order("display_order ASC, id ASC").
		Find(&testimonials).Error
	if err != nil {
		return nil, err
	}
	if testimonials == nil {
		testimonials = []model.Testimonial{}
	}
	return testimonials, nil
}

// FindByID mengambil testimoni berdasarkan ID.
func (r *TestimonialRepository) FindByID(id uint) (*model.Testimonial, error) {
	var testimonial model.Testimonial
	err := r.db.First(&testimonial, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &testimonial, nil
}

// Update memperbarui data testimoni.
func (r *TestimonialRepository) Update(testimonial *model.Testimonial) error {
	return r.db.Save(testimonial).Error
}

// Delete menghapus testimoni berdasarkan ID.
func (r *TestimonialRepository) Delete(id uint) error {
	return r.db.Delete(&model.Testimonial{}, id).Error
}
