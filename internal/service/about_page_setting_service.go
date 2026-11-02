package service

import (
	"errors"
	"fmt"
	"strings"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateAboutPageSettingRequest struct {
	Title     string `json:"title"`
	TeamTitle string `json:"team_title"`
	IsActive  *bool  `json:"is_active"`
}

type UpdateAboutPageSettingRequest struct {
	Title     string `json:"title"`
	TeamTitle string `json:"team_title"`
	IsActive  *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type AboutPageSettingService struct {
	settingRepo *repository.AboutPageSettingRepository
}

func NewAboutPageSettingService(settingRepo *repository.AboutPageSettingRepository) *AboutPageSettingService {
	return &AboutPageSettingService{settingRepo: settingRepo}
}

// GetActiveSetting mengambil konfigurasi judul halaman Tentang yang sedang aktif.
func (s *AboutPageSettingService) GetActiveSetting() (*model.AboutPageSetting, error) {
	setting, err := s.settingRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, errors.New("pengaturan judul aktif tidak ditemukan")
	}
	return setting, nil
}

// GetSettingByID mengambil konfigurasi judul berdasarkan ID.
func (s *AboutPageSettingService) GetSettingByID(id uint) (*model.AboutPageSetting, error) {
	setting, err := s.settingRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, errors.New("pengaturan judul tidak ditemukan")
	}
	return setting, nil
}

// GetAllSettings mengambil seluruh riwayat konfigurasi judul untuk kebutuhan admin.
func (s *AboutPageSettingService) GetAllSettings() ([]model.AboutPageSetting, error) {
	settings, err := s.settingRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if settings == nil {
		settings = []model.AboutPageSetting{}
	}
	return settings, nil
}

// CreateSetting membuat konfigurasi judul baru.
// Jika status aktif diset true, seluruh konfigurasi aktif sebelumnya akan dinonaktifkan
// untuk memastikan hanya ada satu konfigurasi aktif.
func (s *AboutPageSettingService) CreateSetting(req CreateAboutPageSettingRequest) (*model.AboutPageSetting, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = "Tentang Kami"
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Pastikan hanya satu konfigurasi yang aktif
	if isActive {
		if err := s.settingRepo.DeactivateAll(); err != nil {
			return nil, fmt.Errorf("gagal menonaktifkan konfigurasi aktif sebelumnya: %w", err)
		}
	}

	setting := &model.AboutPageSetting{
		Title:     title,
		TeamTitle: strings.TrimSpace(req.TeamTitle),
		IsActive:  isActive,
	}

	if err := s.settingRepo.Create(setting); err != nil {
		return nil, fmt.Errorf("gagal membuat konfigurasi judul: %w", err)
	}

	return setting, nil
}

// UpdateSetting memperbarui data konfigurasi judul berdasarkan ID.
// Jika status aktif diset true, konfigurasi lain akan otomatis dinonaktifkan.
func (s *AboutPageSettingService) UpdateSetting(id uint, req UpdateAboutPageSettingRequest) (*model.AboutPageSetting, error) {
	setting, err := s.settingRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, errors.New("pengaturan judul tidak ditemukan")
	}

	// Jangan menimpa title yang sudah ada jika req.Title kosong
	if trimmedTitle := strings.TrimSpace(req.Title); trimmedTitle != "" {
		setting.Title = trimmedTitle
	}
	setting.TeamTitle = strings.TrimSpace(req.TeamTitle)

	if req.IsActive != nil {
		if *req.IsActive {
			// Nonaktifkan konfigurasi aktif selain ID ini
			if err := s.settingRepo.DeactivateAllExcept(id); err != nil {
				return nil, fmt.Errorf("gagal menonaktifkan konfigurasi aktif lainnya: %w", err)
			}
			setting.IsActive = true
		} else {
			setting.IsActive = false
		}
	}

	if err := s.settingRepo.Update(setting); err != nil {
		return nil, fmt.Errorf("gagal memperbarui konfigurasi judul: %w", err)
	}

	return setting, nil
}

// SaveActiveTitle memperbarui judul konfigurasi aktif saat ini, atau membuat konfigurasi baru
// jika belum pernah ada konfigurasi di database. Tetap menjamin hanya satu konfigurasi aktif
// tanpa mengisi teks default secara otomatis.
func (s *AboutPageSettingService) SaveActiveTitle(title string, teamTitle string, isActive *bool) (*model.AboutPageSetting, error) {
	activeSetting, err := s.settingRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	activeStatus := true
	if isActive != nil {
		activeStatus = *isActive
	}

	if activeSetting != nil {
		finalTitle := title
		if strings.TrimSpace(finalTitle) == "" {
			finalTitle = activeSetting.Title
		}
		req := UpdateAboutPageSettingRequest{
			Title:     finalTitle,
			TeamTitle: teamTitle,
			IsActive:  &activeStatus,
		}
		return s.UpdateSetting(activeSetting.ID, req)
	}

	// Jika belum ada konfigurasi aktif, buat baru
	req := CreateAboutPageSettingRequest{
		Title:     title,
		TeamTitle: teamTitle,
		IsActive:  &activeStatus,
	}
	return s.CreateSetting(req)
}
