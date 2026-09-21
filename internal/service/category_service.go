package service

import (
	"errors"
	"fmt"
	"strings"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"omitempty,min=1,max=100"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// GetCategories mengambil semua data kategori
func (s *CategoryService) GetCategories() ([]model.Category, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return categories, nil
}

// GetCategoryByID mengambil detail kategori berdasarkan ID
func (s *CategoryService) GetCategoryByID(id uint) (*model.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if category == nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return category, nil
}

// CreateCategory membuat kategori baru
func (s *CategoryService) CreateCategory(req CreateCategoryRequest) (*model.Category, error) {
	// Trim input
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("nama kategori tidak boleh kosong")
	}

	// Validasi name unique
	existing, err := s.categoryRepo.FindByName(name)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, errors.New("kategori sudah ada")
	}

	// Generate slug dari name
	baseSlug := generateSlug(name)
	if baseSlug == "" {
		return nil, errors.New("nama kategori tidak valid untuk dijadikan slug")
	}

	// Pastikan slug unique
	slug, err := generateUniqueSlug(baseSlug, func(candidate string) (bool, error) {
		found, err := s.categoryRepo.FindBySlug(candidate)
		return found != nil, err
	})
	if err != nil {
		return nil, err
	}

	category := &model.Category{
		Name: name,
		Slug: slug,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, fmt.Errorf("gagal membuat kategori: %w", err)
	}

	return category, nil
}

// UpdateCategory memperbarui data kategori
func (s *CategoryService) UpdateCategory(id uint, req UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if category == nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	// Jika name dikirim dan berbeda, lakukan validasi dan update
	name := strings.TrimSpace(req.Name)
	if name != "" && name != category.Name {
		// Cek keunikan nama
		exists, err := s.categoryRepo.NameExistsExcept(name, id)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		if exists {
			return nil, errors.New("kategori sudah ada")
		}

		// Update slug sesuai name baru
		baseSlug := generateSlug(name)
		if baseSlug == "" {
			return nil, errors.New("nama kategori tidak valid untuk dijadikan slug")
		}

		slug, err := generateUniqueSlug(baseSlug, func(candidate string) (bool, error) {
			exists, err := s.categoryRepo.SlugExistsExcept(candidate, id)
			return exists, err
		})
		if err != nil {
			return nil, err
		}

		category.Name = name
		category.Slug = slug
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, fmt.Errorf("gagal memperbarui kategori: %w", err)
	}

	return category, nil
}

// DeleteCategory menghapus kategori berdasarkan ID
func (s *CategoryService) DeleteCategory(id uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if category == nil {
		return errors.New("kategori tidak ditemukan")
	}

	// Cek apakah kategori masih digunakan oleh artikel
	hasArticles, err := s.categoryRepo.HasArticles(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if hasArticles {
		return errors.New("kategori masih digunakan oleh artikel")
	}

	return s.categoryRepo.Delete(id)
}
