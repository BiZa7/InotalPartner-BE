package service

import (
	"errors"
	"fmt"
	"strings"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type UpdateTagRequest struct {
	Name string `json:"name" binding:"omitempty,min=1,max=100"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}

// GetTags mengambil semua data tag
func (s *TagService) GetTags() ([]model.Tag, error) {
	tags, err := s.tagRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return tags, nil
}

// GetTagByID mengambil detail tag berdasarkan ID
func (s *TagService) GetTagByID(id uint) (*model.Tag, error) {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if tag == nil {
		return nil, errors.New("tag tidak ditemukan")
	}
	return tag, nil
}

// CreateTag membuat tag baru
func (s *TagService) CreateTag(req CreateTagRequest) (*model.Tag, error) {
	// Trim input
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("nama tag tidak boleh kosong")
	}

	// Validasi name unique
	existing, err := s.tagRepo.FindByName(name)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, errors.New("tag sudah ada")
	}

	// Generate slug dari name
	baseSlug := generateSlug(name)
	if baseSlug == "" {
		return nil, errors.New("nama tag tidak valid untuk dijadikan slug")
	}

	// Pastikan slug unique
	slug, err := generateUniqueSlug(baseSlug, func(candidate string) (bool, error) {
		found, err := s.tagRepo.FindBySlug(candidate)
		return found != nil, err
	})
	if err != nil {
		return nil, err
	}

	tag := &model.Tag{
		Name: name,
		Slug: slug,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, fmt.Errorf("gagal membuat tag: %w", err)
	}

	return tag, nil
}

// UpdateTag memperbarui data tag
func (s *TagService) UpdateTag(id uint, req UpdateTagRequest) (*model.Tag, error) {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if tag == nil {
		return nil, errors.New("tag tidak ditemukan")
	}

	// Jika name dikirim dan berbeda, lakukan validasi dan update
	name := strings.TrimSpace(req.Name)
	if name != "" && name != tag.Name {
		// Cek keunikan nama
		exists, err := s.tagRepo.NameExistsExcept(name, id)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		if exists {
			return nil, errors.New("tag sudah ada")
		}

		// Update slug sesuai name baru
		baseSlug := generateSlug(name)
		if baseSlug == "" {
			return nil, errors.New("nama tag tidak valid untuk dijadikan slug")
		}

		slug, err := generateUniqueSlug(baseSlug, func(candidate string) (bool, error) {
			exists, err := s.tagRepo.SlugExistsExcept(candidate, id)
			return exists, err
		})
		if err != nil {
			return nil, err
		}

		tag.Name = name
		tag.Slug = slug
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, fmt.Errorf("gagal memperbarui tag: %w", err)
	}

	return tag, nil
}

// DeleteTag menghapus tag berdasarkan ID
// Relasi many-to-many (article_tags) menggunakan OnDelete:CASCADE,
// sehingga asosiasi di article_tags akan terhapus otomatis oleh database,
// tanpa menghapus artikel itu sendiri.
func (s *TagService) DeleteTag(id uint) error {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if tag == nil {
		return errors.New("tag tidak ditemukan")
	}

	return s.tagRepo.Delete(id)
}
