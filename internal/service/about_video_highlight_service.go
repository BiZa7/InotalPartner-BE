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
	ErrAboutVideoHighlightNotFound = errors.New("highlight video tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateAboutVideoHighlightRequest struct {
	Icon         string `json:"icon"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateAboutVideoHighlightRequest struct {
	Icon         string `json:"icon"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type AboutVideoHighlightService struct {
	highlightRepo *repository.AboutVideoHighlightRepository
}

func NewAboutVideoHighlightService(highlightRepo *repository.AboutVideoHighlightRepository) *AboutVideoHighlightService {
	return &AboutVideoHighlightService{highlightRepo: highlightRepo}
}

// GetAllHighlights mengambil seluruh highlight video untuk dashboard admin.
func (s *AboutVideoHighlightService) GetAllHighlights() ([]model.AboutVideoHighlight, error) {
	highlights, err := s.highlightRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if highlights == nil {
		highlights = []model.AboutVideoHighlight{}
	}
	return highlights, nil
}

// GetPublicHighlights mengambil highlight video aktif saja untuk halaman publik.
func (s *AboutVideoHighlightService) GetPublicHighlights() ([]model.AboutVideoHighlight, error) {
	highlights, err := s.highlightRepo.FindActivePublic()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if highlights == nil {
		highlights = []model.AboutVideoHighlight{}
	}
	return highlights, nil
}

// GetHighlightByID mengambil detail highlight video berdasarkan ID.
func (s *AboutVideoHighlightService) GetHighlightByID(id uint) (*model.AboutVideoHighlight, error) {
	highlight, err := s.highlightRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if highlight == nil {
		return nil, ErrAboutVideoHighlightNotFound
	}
	return highlight, nil
}

// CreateHighlight membuat highlight baru dengan validasi.
func (s *AboutVideoHighlightService) CreateHighlight(req CreateAboutVideoHighlightRequest) (*model.AboutVideoHighlight, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("judul highlight wajib diisi")
	}
	if utf8.RuneCountInString(title) > 150 {
		return nil, newValidationError("judul highlight maksimal 150 karakter")
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		return nil, newValidationError("deskripsi highlight wajib diisi")
	}

	icon := strings.TrimSpace(req.Icon)
	if utf8.RuneCountInString(icon) > 50 {
		return nil, newValidationError("ikon maksimal 50 karakter")
	}

	if req.DisplayOrder < 0 {
		return nil, newValidationError("display_order tidak boleh bernilai negatif")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	highlight := &model.AboutVideoHighlight{
		Icon:         icon,
		Title:        title,
		Description:  description,
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	if err := s.highlightRepo.Create(highlight); err != nil {
		return nil, fmt.Errorf("gagal membuat highlight video: %w", err)
	}

	return highlight, nil
}

// UpdateHighlight memperbarui highlight video yang ada.
func (s *AboutVideoHighlightService) UpdateHighlight(id uint, req UpdateAboutVideoHighlightRequest) (*model.AboutVideoHighlight, error) {
	highlight, err := s.highlightRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if highlight == nil {
		return nil, ErrAboutVideoHighlightNotFound
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("judul highlight wajib diisi")
	}
	if utf8.RuneCountInString(title) > 150 {
		return nil, newValidationError("judul highlight maksimal 150 karakter")
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		return nil, newValidationError("deskripsi highlight wajib diisi")
	}

	icon := strings.TrimSpace(req.Icon)
	if utf8.RuneCountInString(icon) > 50 {
		return nil, newValidationError("ikon maksimal 50 karakter")
	}

	displayOrder := highlight.DisplayOrder
	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, newValidationError("display_order tidak boleh bernilai negatif")
		}
		displayOrder = *req.DisplayOrder
	}

	isActive := highlight.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	highlight.Icon = icon
	highlight.Title = title
	highlight.Description = description
	highlight.DisplayOrder = displayOrder
	highlight.IsActive = isActive

	if err := s.highlightRepo.Update(highlight); err != nil {
		return nil, fmt.Errorf("gagal memperbarui highlight video: %w", err)
	}

	return highlight, nil
}

// DeleteHighlight menghapus highlight video berdasarkan ID.
func (s *AboutVideoHighlightService) DeleteHighlight(id uint) error {
	highlight, err := s.highlightRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if highlight == nil {
		return ErrAboutVideoHighlightNotFound
	}

	if err := s.highlightRepo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus highlight video: %w", err)
	}

	return nil
}

