package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type AboutVideoRepository struct {
	db *gorm.DB
}

func NewAboutVideoRepository(db *gorm.DB) *AboutVideoRepository {
	return &AboutVideoRepository{db: db}
}

// FindAll mengambil seluruh video profil untuk kebutuhan dashboard admin (termasuk yang nonaktif).
// Diurutkan berdasarkan id DESC agar video terbaru berada di atas.
func (r *AboutVideoRepository) FindAll() ([]model.AboutVideo, error) {
	var videos []model.AboutVideo
	err := r.db.Order("id DESC").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	if videos == nil {
		videos = []model.AboutVideo{}
	}
	return videos, nil
}

// FindByID mengambil satu video profil berdasarkan ID.
func (r *AboutVideoRepository) FindByID(id uint) (*model.AboutVideo, error) {
	var video model.AboutVideo
	err := r.db.First(&video, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// FindActivePublic mengambil seluruh video profil yang aktif untuk halaman publik.
// Diurutkan secara konsisten (id DESC) dan tidak menampilkan video nonaktif (is_active = true).
func (r *AboutVideoRepository) FindActivePublic() ([]model.AboutVideo, error) {
	var videos []model.AboutVideo
	err := r.db.
		Where("is_active = ?", true).
		Order("id DESC").
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	if videos == nil {
		videos = []model.AboutVideo{}
	}
	return videos, nil
}

// Create membuat record video profil baru.
func (r *AboutVideoRepository) Create(video *model.AboutVideo) error {
	return r.db.Create(video).Error
}

// Update memperbarui data video profil.
func (r *AboutVideoRepository) Update(video *model.AboutVideo) error {
	return r.db.Save(video).Error
}

// Delete menghapus data video profil berdasarkan ID.
func (r *AboutVideoRepository) Delete(id uint) error {
	return r.db.Delete(&model.AboutVideo{}, id).Error
}
