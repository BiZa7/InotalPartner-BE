package service

import (
	"errors"
	"fmt"
	"strings"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── Sentinel Errors ──────────────────────────────────────────────────────────

var (
	ErrAboutIntroPointNotFound = errors.New("point perkenalan tentang kami tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateAboutIntroPointRequest struct {
	SettingID    uint   `json:"setting_id"`
	SectionType  string `json:"section_type"` // 'mission' atau 'vision'
	Content      string `json:"content"`
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateAboutIntroPointRequest struct {
	SettingID    *uint  `json:"setting_id"`
	SectionType  string `json:"section_type"` // 'mission' atau 'vision'
	Content      string `json:"content"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type AboutIntroPointService struct {
	repo        *repository.AboutIntroPointRepository
	settingRepo *repository.AboutIntroSettingRepository
}

func NewAboutIntroPointService(
	repo *repository.AboutIntroPointRepository,
	settingRepo *repository.AboutIntroSettingRepository,
) *AboutIntroPointService {
	return &AboutIntroPointService{
		repo:        repo,
		settingRepo: settingRepo,
	}
}

// GetAll mengambil seluruh point untuk dashboard admin.
func (s *AboutIntroPointService) GetAll() ([]model.AboutIntroPoint, error) {
	points, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if points == nil {
		points = []model.AboutIntroPoint{}
	}
	return points, nil
}

// GetByID mengambil point berdasarkan ID.
func (s *AboutIntroPointService) GetByID(id uint) (*model.AboutIntroPoint, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	point, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if point == nil {
		return nil, ErrAboutIntroPointNotFound
	}
	return point, nil
}

// GetBySettingID mengambil seluruh point berdasarkan SettingID.
func (s *AboutIntroPointService) GetBySettingID(settingID uint) ([]model.AboutIntroPoint, error) {
	if settingID == 0 {
		return nil, newValidationError("setting_id tidak valid")
	}

	points, err := s.repo.FindBySettingID(settingID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if points == nil {
		points = []model.AboutIntroPoint{}
	}
	return points, nil
}

// GetBySettingIDAndType mengambil seluruh point berdasarkan SettingID dan SectionType ('mission' atau 'vision').
func (s *AboutIntroPointService) GetBySettingIDAndType(settingID uint, sectionType string) ([]model.AboutIntroPoint, error) {
	if settingID == 0 {
		return nil, newValidationError("setting_id tidak valid")
	}

	st := strings.ToLower(strings.TrimSpace(sectionType))
	if st != "mission" && st != "vision" {
		return nil, newValidationError("section_type harus bernilai 'mission' atau 'vision'")
	}

	points, err := s.repo.FindBySettingIDAndType(settingID, st)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if points == nil {
		points = []model.AboutIntroPoint{}
	}
	return points, nil
}

// GetActiveBySettingID mengambil seluruh point aktif berdasarkan SettingID.
func (s *AboutIntroPointService) GetActiveBySettingID(settingID uint) ([]model.AboutIntroPoint, error) {
	if settingID == 0 {
		return nil, newValidationError("setting_id tidak valid")
	}

	points, err := s.repo.FindActiveBySettingID(settingID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if points == nil {
		points = []model.AboutIntroPoint{}
	}
	return points, nil
}

// Create membuat data point baru.
func (s *AboutIntroPointService) Create(req CreateAboutIntroPointRequest) (*model.AboutIntroPoint, error) {
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

	sectionType := strings.ToLower(strings.TrimSpace(req.SectionType))
	if sectionType != "mission" && sectionType != "vision" {
		return nil, newValidationError("section_type wajib bernilai 'mission' atau 'vision'")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, newValidationError("content wajib diisi")
	}

	if req.DisplayOrder < 0 {
		return nil, newValidationError("display_order tidak boleh bernilai negatif")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	point := &model.AboutIntroPoint{
		SettingID:    req.SettingID,
		SectionType:  sectionType,
		Content:      content,
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	if err := s.repo.Create(point); err != nil {
		return nil, fmt.Errorf("gagal membuat point perkenalan tentang kami: %w", err)
	}

	return point, nil
}

// Update memperbarui data point yang sudah ada.
func (s *AboutIntroPointService) Update(id uint, req UpdateAboutIntroPointRequest) (*model.AboutIntroPoint, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	point, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if point == nil {
		return nil, ErrAboutIntroPointNotFound
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
		point.SettingID = *req.SettingID
	}

	sectionType := strings.ToLower(strings.TrimSpace(req.SectionType))
	if sectionType != "mission" && sectionType != "vision" {
		return nil, newValidationError("section_type wajib bernilai 'mission' atau 'vision'")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, newValidationError("content wajib diisi")
	}

	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, newValidationError("display_order tidak boleh bernilai negatif")
		}
		point.DisplayOrder = *req.DisplayOrder
	}

	if req.IsActive != nil {
		point.IsActive = *req.IsActive
	}

	point.SectionType = sectionType
	point.Content = content

	if err := s.repo.Update(point); err != nil {
		return nil, fmt.Errorf("gagal memperbarui point perkenalan tentang kami: %w", err)
	}

	return point, nil
}

// Delete menghapus data point berdasarkan ID.
func (s *AboutIntroPointService) Delete(id uint) error {
	if id == 0 {
		return newValidationError("ID tidak valid")
	}

	point, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if point == nil {
		return ErrAboutIntroPointNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus point perkenalan tentang kami: %w", err)
	}

	return nil
}

