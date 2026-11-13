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
	ErrHomeCTANotFound = errors.New("data cta tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateHomeCTARequest struct {
	SectionLabel string `json:"section_label"`
	Title        string `json:"title"`
	ButtonText   string `json:"button_text"`
	ButtonURL    string `json:"button_url"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateHomeCTARequest struct {
	SectionLabel string `json:"section_label"`
	Title        string `json:"title"`
	ButtonText   string `json:"button_text"`
	ButtonURL    string `json:"button_url"`
	IsActive     *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type HomeCTAService struct {
	ctaRepo *repository.HomeCTARepository
}

func NewHomeCTAService(ctaRepo *repository.HomeCTARepository) *HomeCTAService {
	return &HomeCTAService{ctaRepo: ctaRepo}
}

// GetAll mengambil seluruh data CTA untuk dashboard admin.
func (s *HomeCTAService) GetAll() ([]model.HomeCTA, error) {
	ctas, err := s.ctaRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if ctas == nil {
		ctas = []model.HomeCTA{}
	}
	return ctas, nil
}

// GetByID mengambil detail data CTA berdasarkan ID.
func (s *HomeCTAService) GetByID(id uint) (*model.HomeCTA, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	cta, err := s.ctaRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if cta == nil {
		return nil, ErrHomeCTANotFound
	}
	return cta, nil
}

// GetActivePublic mengambil satu data CTA aktif untuk publik.
// Jika belum ada data aktif di DB, kembalikan default CTA agar tampilan tidak kosong.
func (s *HomeCTAService) GetActivePublic() (*model.HomeCTA, error) {
	cta, err := s.ctaRepo.FindActivePublic()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if cta == nil {
		return &model.HomeCTA{
			SectionLabel: "SIAP BERKOLABORASI?",
			Title:        "Mari Bergabung dan Bergerak Bersama INOTAL.",
			ButtonText:   "Hubungi Kami",
			ButtonURL:    "/kontak",
			IsActive:     true,
		}, nil
	}
	return cta, nil
}

// Create membuat record CTA baru.
func (s *HomeCTAService) Create(req CreateHomeCTARequest) (*model.HomeCTA, error) {
	sectionLabel := strings.TrimSpace(req.SectionLabel)
	if sectionLabel == "" {
		return nil, newValidationError("section_label wajib diisi")
	}
	if utf8.RuneCountInString(sectionLabel) > 100 {
		return nil, newValidationError("section_label maksimal 100 karakter")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("title wajib diisi")
	}

	buttonText := strings.TrimSpace(req.ButtonText)
	if buttonText == "" {
		return nil, newValidationError("button_text wajib diisi")
	}
	if utf8.RuneCountInString(buttonText) > 100 {
		return nil, newValidationError("button_text maksimal 100 karakter")
	}

	buttonURL := strings.TrimSpace(req.ButtonURL)
	if buttonURL == "" {
		return nil, newValidationError("button_url wajib diisi")
	}
	if utf8.RuneCountInString(buttonURL) > 255 {
		return nil, newValidationError("button_url maksimal 255 karakter")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	cta := &model.HomeCTA{
		SectionLabel: sectionLabel,
		Title:        title,
		ButtonText:   buttonText,
		ButtonURL:    buttonURL,
		IsActive:     isActive,
	}

	if err := s.ctaRepo.Create(cta); err != nil {
		return nil, fmt.Errorf("gagal membuat cta: %w", err)
	}
	return cta, nil
}

// Update memperbarui data CTA berdasarkan ID.
func (s *HomeCTAService) Update(id uint, req UpdateHomeCTARequest) (*model.HomeCTA, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	cta, err := s.ctaRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if cta == nil {
		return nil, ErrHomeCTANotFound
	}

	sectionLabel := strings.TrimSpace(req.SectionLabel)
	if sectionLabel == "" {
		return nil, newValidationError("section_label wajib diisi")
	}
	if utf8.RuneCountInString(sectionLabel) > 100 {
		return nil, newValidationError("section_label maksimal 100 karakter")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("title wajib diisi")
	}

	buttonText := strings.TrimSpace(req.ButtonText)
	if buttonText == "" {
		return nil, newValidationError("button_text wajib diisi")
	}
	if utf8.RuneCountInString(buttonText) > 100 {
		return nil, newValidationError("button_text maksimal 100 karakter")
	}

	buttonURL := strings.TrimSpace(req.ButtonURL)
	if buttonURL == "" {
		return nil, newValidationError("button_url wajib diisi")
	}
	if utf8.RuneCountInString(buttonURL) > 255 {
		return nil, newValidationError("button_url maksimal 255 karakter")
	}

	if req.IsActive != nil {
		cta.IsActive = *req.IsActive
	}

	cta.SectionLabel = sectionLabel
	cta.Title = title
	cta.ButtonText = buttonText
	cta.ButtonURL = buttonURL

	if err := s.ctaRepo.Update(cta); err != nil {
		return nil, fmt.Errorf("gagal memperbarui cta: %w", err)
	}
	return cta, nil
}

// Delete menghapus data CTA berdasarkan ID.
func (s *HomeCTAService) Delete(id uint) error {
	if id == 0 {
		return newValidationError("ID tidak valid")
	}

	cta, err := s.ctaRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if cta == nil {
		return ErrHomeCTANotFound
	}

	if err := s.ctaRepo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus cta: %w", err)
	}
	return nil
}

