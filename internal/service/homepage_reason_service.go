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
	ErrHomepageReasonNotFound = errors.New("kartu alasan tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type SaveHomepageReasonSettingRequest struct {
	SectionLabel       string `json:"section_label"`
	SectionTitle       string `json:"section_title"`
	SectionDescription string `json:"section_description"`
	ImageOneURL        string `json:"image_one_url"`
	ImageTwoURL        string `json:"image_two_url"`
	ImageThreeURL      string `json:"image_three_url"`
	IsActive           *bool  `json:"is_active"`
}

type CreateHomepageReasonRequest struct {
	Number       string `json:"number"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateHomepageReasonRequest struct {
	Number       string `json:"number"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type PublicHomepageReasonsResponse struct {
	Setting *model.HomepageReasonSetting `json:"setting"`
	Reasons []model.HomepageReason       `json:"reasons"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type HomepageReasonService struct {
	repo *repository.HomepageReasonRepository
}

func NewHomepageReasonService(repo *repository.HomepageReasonRepository) *HomepageReasonService {
	return &HomepageReasonService{repo: repo}
}

// ─── Section Settings ─────────────────────────────────────────────────────────

// GetSetting mengambil konfigurasi section alasan dari database.
// Jika belum ada, mengembalikan default fallback tanpa error.
func (s *HomepageReasonService) GetSetting() (*model.HomepageReasonSetting, error) {
	setting, err := s.repo.FindFirstSetting()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return &model.HomepageReasonSetting{
			SectionLabel:       "MENGAPA BERGABUNG",
			SectionTitle:       "Saatnya Bersama Beraksi",
			SectionDescription: "Indonesia butuh pemuda generasi masa depan emas 2045.",
			ImageOneURL:        "https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=500&q=80",
			ImageTwoURL:        "https://images.unsplash.com/photo-1573497491208-6b1acb260507?w=500&q=80",
			ImageThreeURL:      "https://images.unsplash.com/photo-1600880292203-757bb62b4baf?w=500&q=80",
			IsActive:           true,
		}, nil
	}
	return setting, nil
}

// SaveSetting membuat atau memperbarui konfigurasi section alasan (upsert single-record).
func (s *HomepageReasonService) SaveSetting(req SaveHomepageReasonSettingRequest) (*model.HomepageReasonSetting, error) {
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

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	existing, err := s.repo.FindFirstSetting()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if existing == nil {
		newSetting := &model.HomepageReasonSetting{
			SectionLabel:       sectionLabel,
			SectionTitle:       sectionTitle,
			SectionDescription: strings.TrimSpace(req.SectionDescription),
			ImageOneURL:        strings.TrimSpace(req.ImageOneURL),
			ImageTwoURL:        strings.TrimSpace(req.ImageTwoURL),
			ImageThreeURL:      strings.TrimSpace(req.ImageThreeURL),
			IsActive:           isActive,
		}
		if err := s.repo.CreateSetting(newSetting); err != nil {
			return nil, fmt.Errorf("gagal membuat setting: %w", err)
		}
		return newSetting, nil
	}

	existing.SectionLabel = sectionLabel
	existing.SectionTitle = sectionTitle
	existing.SectionDescription = strings.TrimSpace(req.SectionDescription)
	existing.ImageOneURL = strings.TrimSpace(req.ImageOneURL)
	existing.ImageTwoURL = strings.TrimSpace(req.ImageTwoURL)
	existing.ImageThreeURL = strings.TrimSpace(req.ImageThreeURL)
	existing.IsActive = isActive

	if err := s.repo.UpdateSetting(existing); err != nil {
		return nil, fmt.Errorf("gagal memperbarui setting: %w", err)
	}
	return existing, nil
}

// ─── Reason Cards ─────────────────────────────────────────────────────────────

// GetAllReasons mengambil seluruh data kartu alasan untuk admin (termasuk nonaktif).
func (s *HomepageReasonService) GetAllReasons() ([]model.HomepageReason, error) {
	reasons, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if reasons == nil {
		reasons = []model.HomepageReason{}
	}
	return reasons, nil
}

// GetReasonByID mengambil detail kartu alasan berdasarkan ID.
func (s *HomepageReasonService) GetReasonByID(id uint) (*model.HomepageReason, error) {
	reason, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrHomepageReasonNotFound
	}
	return reason, nil
}

// CreateReason membuat kartu alasan baru.
func (s *HomepageReasonService) CreateReason(req CreateHomepageReasonRequest) (*model.HomepageReason, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("judul kartu alasan (title) wajib diisi")
	}
	if utf8.RuneCountInString(title) > 255 {
		return nil, newValidationError("judul kartu alasan (title) maksimal 255 karakter")
	}

	number := strings.TrimSpace(req.Number)
	if utf8.RuneCountInString(number) > 20 {
		return nil, newValidationError("nomor kartu (number) maksimal 20 karakter")
	}

	if req.DisplayOrder < 0 {
		return nil, newValidationError("display_order tidak boleh bernilai negatif")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	reason := &model.HomepageReason{
		Number:       number,
		Title:        title,
		Description:  strings.TrimSpace(req.Description),
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	if err := s.repo.Create(reason); err != nil {
		return nil, fmt.Errorf("gagal membuat kartu alasan: %w", err)
	}
	return reason, nil
}

// UpdateReason memperbarui kartu alasan berdasarkan ID.
func (s *HomepageReasonService) UpdateReason(id uint, req UpdateHomepageReasonRequest) (*model.HomepageReason, error) {
	reason, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrHomepageReasonNotFound
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("judul kartu alasan (title) wajib diisi")
	}
	if utf8.RuneCountInString(title) > 255 {
		return nil, newValidationError("judul kartu alasan (title) maksimal 255 karakter")
	}

	number := strings.TrimSpace(req.Number)
	if utf8.RuneCountInString(number) > 20 {
		return nil, newValidationError("nomor kartu (number) maksimal 20 karakter")
	}

	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, newValidationError("display_order tidak boleh bernilai negatif")
		}
		reason.DisplayOrder = *req.DisplayOrder
	}

	if req.IsActive != nil {
		reason.IsActive = *req.IsActive
	}

	reason.Number = number
	reason.Title = title
	reason.Description = strings.TrimSpace(req.Description)

	if err := s.repo.Update(reason); err != nil {
		return nil, fmt.Errorf("gagal memperbarui kartu alasan: %w", err)
	}
	return reason, nil
}

// DeleteReason menghapus kartu alasan berdasarkan ID.
func (s *HomepageReasonService) DeleteReason(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return ErrHomepageReasonNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus kartu alasan: %w", err)
	}
	return nil
}

// ─── Public ───────────────────────────────────────────────────────────────────

// GetPublicReasonsResponse mengembalikan konfigurasi aktif dan seluruh kartu alasan aktif.
func (s *HomepageReasonService) GetPublicReasonsResponse() (*PublicHomepageReasonsResponse, error) {
	setting, err := s.repo.FindActiveSetting()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Fallback jika setting belum ada di DB
	if setting == nil {
		setting = &model.HomepageReasonSetting{
			SectionLabel:       "MENGAPA BERGABUNG",
			SectionTitle:       "Saatnya Bersama Beraksi",
			SectionDescription: "Indonesia butuh pemuda generasi masa depan emas 2045.",
			ImageOneURL:        "https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=500&q=80",
			ImageTwoURL:        "https://images.unsplash.com/photo-1573497491208-6b1acb260507?w=500&q=80",
			ImageThreeURL:      "https://images.unsplash.com/photo-1600880292203-757bb62b4baf?w=500&q=80",
			IsActive:           true,
		}
	}

	reasons, err := s.repo.FindAllActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if reasons == nil {
		reasons = []model.HomepageReason{}
	}

	return &PublicHomepageReasonsResponse{
		Setting: setting,
		Reasons: reasons,
	}, nil
}

