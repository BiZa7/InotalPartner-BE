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
	ErrPartnershipServiceNotFound = errors.New("layanan partnership tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreatePartnershipServiceRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdatePartnershipServiceRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type PartnershipServiceService struct {
	repo *repository.PartnershipServiceRepository
}

func NewPartnershipServiceService(repo *repository.PartnershipServiceRepository) *PartnershipServiceService {
	return &PartnershipServiceService{repo: repo}
}

// GetAll mengambil seluruh kartu layanan untuk dashboard admin.
func (s *PartnershipServiceService) GetAll() ([]model.PartnershipService, error) {
	services, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if services == nil {
		services = []model.PartnershipService{}
	}
	return services, nil
}

// GetByID mengambil kartu layanan berdasarkan ID.
func (s *PartnershipServiceService) GetByID(id uint) (*model.PartnershipService, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	service, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if service == nil {
		return nil, ErrPartnershipServiceNotFound
	}
	return service, nil
}

// GetActive mengambil kartu-kartu layanan aktif untuk publik.
// Jika di database belum ada data aktif, kembalikan 4 kartu default agar tampilan Beranda tidak kosong.
func (s *PartnershipServiceService) GetActive() ([]model.PartnershipService, error) {
	services, err := s.repo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if len(services) == 0 {
		return []model.PartnershipService{
			{
				Title:        "IT CONSULTANT",
				Description:  "Layanan konsultasi IT profesional untuk bisnis yang ingin memanfaatkan teknologi secara strategis demi pertumbuhan dan efisiensi operasional.",
				Icon:         "monitor-cog",
				DisplayOrder: 1,
				IsActive:     true,
			},
			{
				Title:        "SYSTEM INTEGRATION",
				Description:  "Memastikan komponen perangkat lunak dan keras dalam ekosistem IT bisnis Anda bekerja secara harmonis dan terpadu.",
				Icon:         "layers",
				DisplayOrder: 2,
				IsActive:     true,
			},
			{
				Title:        "RESOURCE MANAGEMENT",
				Description:  "Layanan konsultasi dan implementasi strategis untuk mengalokasikan sumber daya secara efektif demi efisiensi biaya maksimal.",
				Icon:         "settings-2",
				DisplayOrder: 3,
				IsActive:     true,
			},
			{
				Title:        "SOFTWARE DEVELOPMENT",
				Description:  "Pengembangan solusi perangkat lunak kustom yang disesuaikan dengan kebutuhan spesifik dan visi bisnis klien.",
				Icon:         "code-2",
				DisplayOrder: 4,
				IsActive:     true,
			},
		}, nil
	}
	return services, nil
}

// Create membuat kartu layanan baru.
func (s *PartnershipServiceService) Create(req CreatePartnershipServiceRequest) (*model.PartnershipService, error) {
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

	icon := strings.TrimSpace(req.Icon)
	if icon == "" {
		return nil, newValidationError("icon wajib diisi")
	}
	if utf8.RuneCountInString(icon) > 100 {
		return nil, newValidationError("icon maksimal 100 karakter")
	}

	if req.DisplayOrder < 0 {
		return nil, newValidationError("display_order tidak boleh bernilai negatif")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	service := &model.PartnershipService{
		Title:        title,
		Description:  description,
		Icon:         icon,
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	if err := s.repo.Create(service); err != nil {
		return nil, fmt.Errorf("gagal membuat layanan partnership: %w", err)
	}

	return service, nil
}

// Update memperbarui kartu layanan yang sudah ada.
func (s *PartnershipServiceService) Update(id uint, req UpdatePartnershipServiceRequest) (*model.PartnershipService, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	service, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if service == nil {
		return nil, ErrPartnershipServiceNotFound
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

	icon := strings.TrimSpace(req.Icon)
	if icon == "" {
		return nil, newValidationError("icon wajib diisi")
	}
	if utf8.RuneCountInString(icon) > 100 {
		return nil, newValidationError("icon maksimal 100 karakter")
	}

	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, newValidationError("display_order tidak boleh bernilai negatif")
		}
		service.DisplayOrder = *req.DisplayOrder
	}

	if req.IsActive != nil {
		service.IsActive = *req.IsActive
	}

	service.Title = title
	service.Description = description
	service.Icon = icon

	if err := s.repo.Update(service); err != nil {
		return nil, fmt.Errorf("gagal memperbarui layanan partnership: %w", err)
	}

	return service, nil
}

// Delete menghapus kartu layanan berdasarkan ID.
func (s *PartnershipServiceService) Delete(id uint) error {
	if id == 0 {
		return newValidationError("ID tidak valid")
	}

	service, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if service == nil {
		return ErrPartnershipServiceNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus layanan partnership: %w", err)
	}

	return nil
}

