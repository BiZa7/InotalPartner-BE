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
	ErrPartnershipSettingNotFound = errors.New("konfigurasi partnership tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreatePartnershipSettingRequest struct {
	Eyebrow     string `json:"eyebrow"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BadgeValue  string `json:"badge_value"`
	BadgeLabel  string `json:"badge_label"`
	IsActive    *bool  `json:"is_active"`
}

type UpdatePartnershipSettingRequest struct {
	Eyebrow     string `json:"eyebrow"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BadgeValue  string `json:"badge_value"`
	BadgeLabel  string `json:"badge_label"`
	IsActive    *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type PartnershipSettingService struct {
	repo *repository.PartnershipSettingRepository
}

func NewPartnershipSettingService(repo *repository.PartnershipSettingRepository) *PartnershipSettingService {
	return &PartnershipSettingService{repo: repo}
}

// GetAll mengambil seluruh konfigurasi untuk dashboard admin.
func (s *PartnershipSettingService) GetAll() ([]model.PartnershipSetting, error) {
	settings, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if settings == nil {
		settings = []model.PartnershipSetting{}
	}
	return settings, nil
}

// GetByID mengambil konfigurasi berdasarkan ID.
func (s *PartnershipSettingService) GetByID(id uint) (*model.PartnershipSetting, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	setting, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, ErrPartnershipSettingNotFound
	}
	return setting, nil
}

// GetActive mengambil satu konfigurasi aktif untuk publik.
// Jika belum ada data aktif di DB, kembalikan default fallback agar tampilan Beranda tidak kosong.
func (s *PartnershipSettingService) GetActive() (*model.PartnershipSetting, error) {
	setting, err := s.repo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return &model.PartnershipSetting{
			Eyebrow:     "Layanan Kami",
			Title:       "Partnership",
			Description: "Kolaborasi strategis bersama mitra untuk mengakselerasi transformasi digital dan inovasi teknologi.",
			BadgeValue:  "7+",
			BadgeLabel:  "Partnership Aktif",
			IsActive:    true,
		}, nil
	}
	return setting, nil
}

// Create membuat konfigurasi baru.
func (s *PartnershipSettingService) Create(req CreatePartnershipSettingRequest) (*model.PartnershipSetting, error) {
	eyebrow := strings.TrimSpace(req.Eyebrow)
	if utf8.RuneCountInString(eyebrow) > 100 {
		return nil, newValidationError("eyebrow maksimal 100 karakter")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("title wajib diisi")
	}
	if utf8.RuneCountInString(title) > 150 {
		return nil, newValidationError("title maksimal 150 karakter")
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		return nil, newValidationError("description wajib diisi")
	}

	badgeValue := strings.TrimSpace(req.BadgeValue)
	if utf8.RuneCountInString(badgeValue) > 50 {
		return nil, newValidationError("badge_value maksimal 50 karakter")
	}

	badgeLabel := strings.TrimSpace(req.BadgeLabel)
	if utf8.RuneCountInString(badgeLabel) > 100 {
		return nil, newValidationError("badge_label maksimal 100 karakter")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	setting := &model.PartnershipSetting{
		Eyebrow:     eyebrow,
		Title:       title,
		Description: description,
		BadgeValue:  badgeValue,
		BadgeLabel:  badgeLabel,
		IsActive:    isActive,
	}

	if err := s.repo.Create(setting); err != nil {
		return nil, fmt.Errorf("gagal membuat konfigurasi partnership: %w", err)
	}

	return setting, nil
}

// Update memperbarui konfigurasi yang sudah ada.
func (s *PartnershipSettingService) Update(id uint, req UpdatePartnershipSettingRequest) (*model.PartnershipSetting, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	setting, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, ErrPartnershipSettingNotFound
	}

	eyebrow := strings.TrimSpace(req.Eyebrow)
	if utf8.RuneCountInString(eyebrow) > 100 {
		return nil, newValidationError("eyebrow maksimal 100 karakter")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("title wajib diisi")
	}
	if utf8.RuneCountInString(title) > 150 {
		return nil, newValidationError("title maksimal 150 karakter")
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		return nil, newValidationError("description wajib diisi")
	}

	badgeValue := strings.TrimSpace(req.BadgeValue)
	if utf8.RuneCountInString(badgeValue) > 50 {
		return nil, newValidationError("badge_value maksimal 50 karakter")
	}

	badgeLabel := strings.TrimSpace(req.BadgeLabel)
	if utf8.RuneCountInString(badgeLabel) > 100 {
		return nil, newValidationError("badge_label maksimal 100 karakter")
	}

	setting.Eyebrow = eyebrow
	setting.Title = title
	setting.Description = description
	setting.BadgeValue = badgeValue
	setting.BadgeLabel = badgeLabel

	if req.IsActive != nil {
		setting.IsActive = *req.IsActive
	}

	if err := s.repo.Update(setting); err != nil {
		return nil, fmt.Errorf("gagal memperbarui konfigurasi partnership: %w", err)
	}

	return setting, nil
}

// Delete menghapus konfigurasi berdasarkan ID.
func (s *PartnershipSettingService) Delete(id uint) error {
	if id == 0 {
		return newValidationError("ID tidak valid")
	}

	setting, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return ErrPartnershipSettingNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus konfigurasi partnership: %w", err)
	}

	return nil
}

