package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create menyimpan kategori baru ke database
func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// FindAll mengambil seluruh data kategori diurutkan berdasarkan ID
func (r *CategoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("id ASC").Find(&categories).Error
	return categories, err
}

// FindByID mencari kategori berdasarkan ID
func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

// FindByName mencari kategori berdasarkan nama
func (r *CategoryRepository) FindByName(name string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

// FindBySlug mencari kategori berdasarkan slug
func (r *CategoryRepository) FindBySlug(slug string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

// Update memperbarui data kategori
func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete menghapus kategori berdasarkan ID
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// NameExistsExcept mengecek apakah nama kategori sudah digunakan oleh kategori lain
func (r *CategoryRepository) NameExistsExcept(name string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Where("name = ? AND id != ?", name, excludeID).Count(&count).Error
	return count > 0, err
}

// SlugExistsExcept mengecek apakah slug sudah digunakan oleh kategori lain
func (r *CategoryRepository) SlugExistsExcept(slug string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Where("slug = ? AND id != ?", slug, excludeID).Count(&count).Error
	return count > 0, err
}

// HasArticles mengecek apakah ada artikel yang menggunakan kategori ini
func (r *CategoryRepository) HasArticles(categoryID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Article{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count > 0, err
}
