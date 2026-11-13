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
	ErrAboutIntroMetricNotFound = errors.New("metric perkenalan tentang kami tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateAboutIntroMetricRequest struct {
	SettingID    uint   `json:"setting_id"`
	Label        string `json:"label"`
	Value        *int   `json:"value"` // Persentase 0–100
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateAboutIntroMetricRequest struct {
	SettingID    *uint  `json:"setting_id"`
	Label        string `json:"label"`
	Value        *int   `json:"value"` // Persentase 0–100
	DisplayOrder *int   `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type AboutIntroMetricService struct {
	repo        *repository.AboutIntroMetricRepository
	settingRepo *repository.AboutIntroSettingRepository
}

func NewAboutIntroMetricService(
	repo *repository.AboutIntroMetricRepository,
	settingRepo *repository.AboutIntroSettingRepository,
) *AboutIntroMetricService {
	return &AboutIntroMetricService{
		repo:        repo,
		settingRepo: settingRepo,
	}
}

// GetAll mengambil seluruh metric untuk dashboard admin.
func (s *AboutIntroMetricService) GetAll() ([]model.AboutIntroMetric, error) {
	metrics, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if metrics == nil {
		metrics = []model.AboutIntroMetric{}
	}
	return metrics, nil
}

// GetByID mengambil metric berdasarkan ID.
func (s *AboutIntroMetricService) GetByID(id uint) (*model.AboutIntroMetric, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	metric, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if metric == nil {
		return nil, ErrAboutIntroMetricNotFound
	}
	return metric, nil
}

// GetBySettingID mengambil seluruh metric berdasarkan SettingID.
func (s *AboutIntroMetricService) GetBySettingID(settingID uint) ([]model.AboutIntroMetric, error) {
	if settingID == 0 {
		return nil, newValidationError("setting_id tidak valid")
	}

	metrics, err := s.repo.FindBySettingID(settingID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if metrics == nil {
		metrics = []model.AboutIntroMetric{}
	}
	return metrics, nil
}

// GetActiveBySettingID mengambil seluruh metric aktif berdasarkan SettingID.
func (s *AboutIntroMetricService) GetActiveBySettingID(settingID uint) ([]model.AboutIntroMetric, error) {
	if settingID == 0 {
		return nil, newValidationError("setting_id tidak valid")
	}

	metrics, err := s.repo.FindActiveBySettingID(settingID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if metrics == nil {
		metrics = []model.AboutIntroMetric{}
	}
	return metrics, nil
}

// Create membuat data metric baru.
func (s *AboutIntroMetricService) Create(req CreateAboutIntroMetricRequest) (*model.AboutIntroMetric, error) {
	if req.SettingID == 0 {
		return nil, newValidationError("setting_id wajib diisi dan harus valid")
	}

	setting, err := s.settingRepo.FindByID(req.SettingID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, newValidationError("setting_id tidak ditemukan")
	}

	label := strings.TrimSpace(req.Label)
	if label == "" {
		return nil, newValidationError("label wajib diisi")
	}
	if utf8.RuneCountInString(label) > 150 {
		return nil, newValidationError("label maksimal 150 karakter")
	}

	if req.Value == nil {
		return nil, newValidationError("value wajib diisi")
	}
	if *req.Value < 0 || *req.Value > 100 {
		return nil, newValidationError("value harus berada pada rentang 0 sampai 100")
	}

	if req.DisplayOrder < 0 {
		return nil, newValidationError("display_order tidak boleh bernilai negatif")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	metric := &model.AboutIntroMetric{
		SettingID:    req.SettingID,
		Label:        label,
		Value:        *req.Value,
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	if err := s.repo.Create(metric); err != nil {
		return nil, fmt.Errorf("gagal membuat metric perkenalan tentang kami: %w", err)
	}

	return metric, nil
}

// Update memperbarui data metric yang sudah ada.
func (s *AboutIntroMetricService) Update(id uint, req UpdateAboutIntroMetricRequest) (*model.AboutIntroMetric, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	metric, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if metric == nil {
		return nil, ErrAboutIntroMetricNotFound
	}

	if req.SettingID != nil {
		if *req.SettingID == 0 {
			return nil, newValidationError("setting_id tidak valid")
		}
		setting, err := s.settingRepo.FindByID(*req.SettingID)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		if setting == nil {
			return nil, newValidationError("setting_id tidak ditemukan")
		}
		metric.SettingID = *req.SettingID
	}

	label := strings.TrimSpace(req.Label)
	if label == "" {
		return nil, newValidationError("label wajib diisi")
	}
	if utf8.RuneCountInString(label) > 150 {
		return nil, newValidationError("label maksimal 150 karakter")
	}

	if req.Value != nil {
		if *req.Value < 0 || *req.Value > 100 {
			return nil, newValidationError("value harus berada pada rentang 0 sampai 100")
		}
		metric.Value = *req.Value
	}

	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, newValidationError("display_order tidak boleh bernilai negatif")
		}
		metric.DisplayOrder = *req.DisplayOrder
	}

	if req.IsActive != nil {
		metric.IsActive = *req.IsActive
	}

	metric.Label = label

	if err := s.repo.Update(metric); err != nil {
		return nil, fmt.Errorf("gagal memperbarui metric perkenalan tentang kami: %w", err)
	}

	return metric, nil
}

// Delete menghapus data metric berdasarkan ID.
func (s *AboutIntroMetricService) Delete(id uint) error {
	if id == 0 {
		return newValidationError("ID tidak valid")
	}

	metric, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if metric == nil {
		return ErrAboutIntroMetricNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus metric perkenalan tentang kami: %w", err)
	}

	return nil
}

