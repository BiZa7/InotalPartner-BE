package service

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type SaveAboutVideoSettingRequest struct {
	SectionLabel       string `json:"section_label"`
	SectionTitle       string `json:"section_title"`
	SectionDescription string `json:"section_description"`
	IsActive           *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type AboutVideoSettingService struct {
	settingRepo *repository.AboutVideoSettingRepository
}

func NewAboutVideoSettingService(settingRepo *repository.AboutVideoSettingRepository) *AboutVideoSettingService {
	return &AboutVideoSettingService{settingRepo: settingRepo}
}

// GetSetting mengambil konfigurasi section video profil yang ada di database.
// Jika belum pernah disimpan, mengembalikan default fallback tanpa error.
func (s *AboutVideoSettingService) GetSetting() (*model.AboutVideoSetting, error) {
	setting, err := s.settingRepo.FindFirst()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		// Nilai fallback default
		return &model.AboutVideoSetting{
			SectionLabel:       "VIDEO PROFIL",
			SectionTitle:       "Lihat sendiri bagaimana kolaborasi bekerja.",
			SectionDescription: "Keripitannya kami tampilkan apa adanya: ruang kerja nyata, interaksi harian, dan energi tim yang menghidupkan setiap ide.",
			IsActive:           true,
		}, nil
	}
	return setting, nil
}

// GetPublicSetting mengambil konfigurasi section video profil yang aktif untuk halaman publik.
// Jika tidak ada atau nonaktif, kembalikan fallback aman atau data default aktif.
func (s *AboutVideoSettingService) GetPublicSetting() (*model.AboutVideoSetting, error) {
	setting, err := s.settingRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		// Fallback aman untuk publik jika belum ada record yang dibuat atau disetting
		return &model.AboutVideoSetting{
			SectionLabel:       "VIDEO PROFIL",
			SectionTitle:       "Lihat sendiri bagaimana kolaborasi bekerja.",
			SectionDescription: "Keripitannya kami tampilkan apa adanya: ruang kerja nyata, interaksi harian, dan energi tim yang menghidupkan setiap ide.",
			IsActive:           true,
		}, nil
	}
	return setting, nil
}

// SaveSetting menyimpan atau memperbarui record pengaturan section video profil (upsert single-record).
func (s *AboutVideoSettingService) SaveSetting(req SaveAboutVideoSettingRequest) (*model.AboutVideoSetting, error) {
	sectionTitle := strings.TrimSpace(req.SectionTitle)
	if sectionTitle == "" {
		return nil, newValidationError("section_title wajib diisi")
	}
	if utf8.RuneCountInString(sectionTitle) > 255 {
		return nil, newValidationError("section_title maksimal 255 karakter")
	}

	sectionLabel := strings.TrimSpace(req.SectionLabel)
	if utf8.RuneCountInString(sectionLabel) > 100 {
		return nil, newValidationError("section_label maksimal 100 karakter")
	}

	sectionDesc := strings.TrimSpace(req.SectionDescription)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	existing, err := s.settingRepo.FindFirst()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if existing == nil {
		newSetting := &model.AboutVideoSetting{
			SectionLabel:       sectionLabel,
			SectionTitle:       sectionTitle,
			SectionDescription: sectionDesc,
			IsActive:           isActive,
		}
		if err := s.settingRepo.Create(newSetting); err != nil {
			return nil, fmt.Errorf("gagal membuat pengaturan section video: %w", err)
		}
		return newSetting, nil
	}

	existing.SectionLabel = sectionLabel
	existing.SectionTitle = sectionTitle
	existing.SectionDescription = sectionDesc
	existing.IsActive = isActive

	if err := s.settingRepo.Update(existing); err != nil {
		return nil, fmt.Errorf("gagal memperbarui pengaturan section video: %w", err)
	}

	return existing, nil
}

