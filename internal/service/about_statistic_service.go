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
	ErrAboutStatisticNotFound = errors.New("statistik tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateAboutStatisticRequest struct {
	Title           string `json:"title"`
	Icon            string `json:"icon"`
	StatisticValue  string `json:"statistic_value"`
	StatisticSuffix string `json:"statistic_suffix"`
	Description     string `json:"description"`
	DisplayOrder    int    `json:"display_order"`
	IsActive        *bool  `json:"is_active"`
}

type UpdateAboutStatisticRequest struct {
	Title           string `json:"title"`
	Icon            string `json:"icon"`
	StatisticValue  string `json:"statistic_value"`
	StatisticSuffix string `json:"statistic_suffix"`
	Description     string `json:"description"`
	DisplayOrder    *int   `json:"display_order"`
	IsActive        *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type AboutStatisticService struct {
	statisticRepo *repository.AboutStatisticRepository
}

func NewAboutStatisticService(statisticRepo *repository.AboutStatisticRepository) *AboutStatisticService {
	return &AboutStatisticService{statisticRepo: statisticRepo}
}

// GetAllStatistics mengambil seluruh data statistik untuk kebutuhan dashboard admin (termasuk nonaktif).
func (s *AboutStatisticService) GetAllStatistics() ([]model.AboutStatistic, error) {
	stats, err := s.statisticRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if stats == nil {
		stats = []model.AboutStatistic{}
	}
	return stats, nil
}

// GetPublicStatistics mengambil data statistik yang aktif saja (is_active = true) untuk halaman publik.
func (s *AboutStatisticService) GetPublicStatistics() ([]model.AboutStatistic, error) {
	stats, err := s.statisticRepo.FindAllActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if stats == nil {
		stats = []model.AboutStatistic{}
	}
	return stats, nil
}

// GetStatisticByID mengambil detail data statistik berdasarkan ID.
func (s *AboutStatisticService) GetStatisticByID(id uint) (*model.AboutStatistic, error) {
	stat, err := s.statisticRepo.FindByID(id)
	if err != nil {
		return nil, ErrAboutStatisticNotFound
	}
	return stat, nil
}

// CreateStatistic membuat data statistik baru dengan validasi lengkap.
func (s *AboutStatisticService) CreateStatistic(req CreateAboutStatisticRequest) (*model.AboutStatistic, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("judul statistik wajib diisi")
	}
	if utf8.RuneCountInString(title) > 150 {
		return nil, newValidationError("judul statistik maksimal 150 karakter")
	}

	icon := strings.TrimSpace(req.Icon)
	if icon == "" {
		return nil, newValidationError("ikon wajib diisi")
	}
	if utf8.RuneCountInString(icon) > 100 {
		return nil, newValidationError("ikon maksimal 100 karakter")
	}

	statisticValue := strings.TrimSpace(req.StatisticValue)
	if statisticValue == "" {
		return nil, newValidationError("nilai statistik wajib diisi")
	}
	if utf8.RuneCountInString(statisticValue) > 30 {
		return nil, newValidationError("nilai statistik maksimal 30 karakter")
	}

	statisticSuffix := strings.TrimSpace(req.StatisticSuffix)
	if utf8.RuneCountInString(statisticSuffix) > 20 {
		return nil, newValidationError("suffix statistik maksimal 20 karakter")
	}

	description := strings.TrimSpace(req.Description)

	if req.DisplayOrder < 0 {
		return nil, newValidationError("display_order tidak boleh bernilai negatif")
	}

	// Penanganan IsActive: default true jika nil / tidak dikirim,
	// namun menghormati nilai false jika dikirim secara eksplisit.
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	statistic := &model.AboutStatistic{
		Title:           title,
		Icon:            icon,
		StatisticValue:  statisticValue,
		StatisticSuffix: statisticSuffix,
		Description:     description,
		DisplayOrder:    req.DisplayOrder,
		IsActive:        isActive,
	}

	if err := s.statisticRepo.Create(statistic); err != nil {
		return nil, fmt.Errorf("gagal membuat statistik: %w", err)
	}

	return statistic, nil
}

// UpdateStatistic memperbarui data statistik berdasarkan ID dengan validasi lengkap.
func (s *AboutStatisticService) UpdateStatistic(id uint, req UpdateAboutStatisticRequest) (*model.AboutStatistic, error) {
	statistic, err := s.statisticRepo.FindByID(id)
	if err != nil {
		return nil, ErrAboutStatisticNotFound
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("judul statistik wajib diisi")
	}
	if utf8.RuneCountInString(title) > 150 {
		return nil, newValidationError("judul statistik maksimal 150 karakter")
	}

	icon := strings.TrimSpace(req.Icon)
	if icon == "" {
		return nil, newValidationError("ikon wajib diisi")
	}
	if utf8.RuneCountInString(icon) > 100 {
		return nil, newValidationError("ikon maksimal 100 karakter")
	}

	statisticValue := strings.TrimSpace(req.StatisticValue)
	if statisticValue == "" {
		return nil, newValidationError("nilai statistik wajib diisi")
	}
	if utf8.RuneCountInString(statisticValue) > 30 {
		return nil, newValidationError("nilai statistik maksimal 30 karakter")
	}

	statisticSuffix := strings.TrimSpace(req.StatisticSuffix)
	if utf8.RuneCountInString(statisticSuffix) > 20 {
		return nil, newValidationError("suffix statistik maksimal 20 karakter")
	}

	description := strings.TrimSpace(req.Description)

	displayOrder := statistic.DisplayOrder
	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, newValidationError("display_order tidak boleh bernilai negatif")
		}
		displayOrder = *req.DisplayOrder
	}

	// Penanganan IsActive: jika dikirim, update sesuai nilai pointer (*req.IsActive)
	// baik bernilai true maupun false secara eksplisit.
	isActive := statistic.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	statistic.Title = title
	statistic.Icon = icon
	statistic.StatisticValue = statisticValue
	statistic.StatisticSuffix = statisticSuffix
	statistic.Description = description
	statistic.DisplayOrder = displayOrder
	statistic.IsActive = isActive

	if err := s.statisticRepo.Update(statistic); err != nil {
		return nil, fmt.Errorf("gagal memperbarui statistik: %w", err)
	}

	return statistic, nil
}

// DeleteStatistic menghapus data statistik berdasarkan ID.
func (s *AboutStatisticService) DeleteStatistic(id uint) error {
	_, err := s.statisticRepo.FindByID(id)
	if err != nil {
		return ErrAboutStatisticNotFound
	}

	if err := s.statisticRepo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus statistik: %w", err)
	}

	return nil
}

