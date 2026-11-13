package service

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── Sentinel Errors ──────────────────────────────────────────────────────────

var (
	ErrPartnershipImageNotFound = errors.New("gambar partnership tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreatePartnershipImageRequest struct {
	ImageURL     string `json:"image_url"`
	AltText      string `json:"alt_text"`
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdatePartnershipImageRequest struct {
	ImageURL     string `json:"image_url"`
	AltText      string `json:"alt_text"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type PartnershipImageService struct {
	repo *repository.PartnershipImageRepository
}

func NewPartnershipImageService(repo *repository.PartnershipImageRepository) *PartnershipImageService {
	return &PartnershipImageService{repo: repo}
}

// GetAll mengambil seluruh data gambar untuk dashboard admin.
func (s *PartnershipImageService) GetAll() ([]model.PartnershipImage, error) {
	images, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if images == nil {
		images = []model.PartnershipImage{}
	}
	return images, nil
}

// GetByID mengambil gambar berdasarkan ID.
func (s *PartnershipImageService) GetByID(id uint) (*model.PartnershipImage, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	image, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if image == nil {
		return nil, ErrPartnershipImageNotFound
	}
	return image, nil
}

// GetActive mengambil gambar-gambar aktif untuk publik.
// Jika di database belum ada gambar aktif, kembalikan gambar default agar tampilan Beranda tidak kosong.
func (s *PartnershipImageService) GetActive() ([]model.PartnershipImage, error) {
	images, err := s.repo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if len(images) == 0 {
		return []model.PartnershipImage{
			{
				ImageURL:     "https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=800&q=80",
				AltText:      "Partnership Meeting",
				DisplayOrder: 1,
				IsActive:     true,
			},
			{
				ImageURL:     "https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?w=600&q=80",
				AltText:      "Collaboration Discussion",
				DisplayOrder: 2,
				IsActive:     true,
			},
		}, nil
	}
	return images, nil
}

// Create membuat gambar baru.
func (s *PartnershipImageService) Create(req CreatePartnershipImageRequest) (*model.PartnershipImage, error) {
	imageURL := strings.TrimSpace(req.ImageURL)
	if imageURL == "" {
		return nil, newValidationError("image_url wajib diisi")
	}
	if utf8.RuneCountInString(imageURL) > 500 {
		return nil, newValidationError("image_url maksimal 500 karakter")
	}

	altText := strings.TrimSpace(req.AltText)
	if utf8.RuneCountInString(altText) > 255 {
		return nil, newValidationError("alt_text maksimal 255 karakter")
	}

	if req.DisplayOrder < 0 {
		return nil, newValidationError("display_order tidak boleh bernilai negatif")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	image := &model.PartnershipImage{
		ImageURL:     imageURL,
		AltText:      altText,
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	if err := s.repo.Create(image); err != nil {
		return nil, fmt.Errorf("gagal membuat gambar partnership: %w", err)
	}

	return image, nil
}

// Update memperbarui data gambar yang sudah ada.
func (s *PartnershipImageService) Update(id uint, req UpdatePartnershipImageRequest) (*model.PartnershipImage, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	image, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if image == nil {
		return nil, ErrPartnershipImageNotFound
	}

	imageURL := strings.TrimSpace(req.ImageURL)
	if imageURL == "" {
		return nil, newValidationError("image_url wajib diisi")
	}
	if utf8.RuneCountInString(imageURL) > 500 {
		return nil, newValidationError("image_url maksimal 500 karakter")
	}

	altText := strings.TrimSpace(req.AltText)
	if utf8.RuneCountInString(altText) > 255 {
		return nil, newValidationError("alt_text maksimal 255 karakter")
	}

	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, newValidationError("display_order tidak boleh bernilai negatif")
		}
		image.DisplayOrder = *req.DisplayOrder
	}

	if req.IsActive != nil {
		image.IsActive = *req.IsActive
	}

	image.ImageURL = imageURL
	image.AltText = altText

	if err := s.repo.Update(image); err != nil {
		return nil, fmt.Errorf("gagal memperbarui gambar partnership: %w", err)
	}

	return image, nil
}

// Delete menghapus gambar berdasarkan ID.
func (s *PartnershipImageService) Delete(id uint) error {
	if id == 0 {
		return newValidationError("ID tidak valid")
	}

	image, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if image == nil {
		return ErrPartnershipImageNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus gambar partnership: %w", err)
	}

	return nil
}

