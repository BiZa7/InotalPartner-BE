package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create menyimpan tag baru ke database
func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

// FindAll mengambil seluruh data tag diurutkan berdasarkan ID
func (r *TagRepository) FindAll() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Order("id ASC").Find(&tags).Error
	return tags, err
}

// FindByID mencari tag berdasarkan ID
func (r *TagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tag, err
}

// FindByIDs mencari list tag berdasarkan slice ID
func (r *TagRepository) FindByIDs(ids []uint) ([]model.Tag, error) {
	if len(ids) == 0 {
		return []model.Tag{}, nil
	}
	var tags []model.Tag
	err := r.db.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

// FindByName mencari tag berdasarkan nama
func (r *TagRepository) FindByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tag, err
}

// FindBySlug mencari tag berdasarkan slug
func (r *TagRepository) FindBySlug(slug string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tag, err
}

// Update memperbarui data tag
func (r *TagRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

// Delete menghapus tag berdasarkan ID
func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

// NameExistsExcept mengecek apakah nama tag sudah digunakan oleh tag lain
func (r *TagRepository) NameExistsExcept(name string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Tag{}).Where("name = ? AND id != ?", name, excludeID).Count(&count).Error
	return count > 0, err
}

// SlugExistsExcept mengecek apakah slug sudah digunakan oleh tag lain
func (r *TagRepository) SlugExistsExcept(slug string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Tag{}).Where("slug = ? AND id != ?", slug, excludeID).Count(&count).Error
	return count > 0, err
}
