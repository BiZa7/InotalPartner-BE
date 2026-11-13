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
	ErrAboutIntroSettingNotFound = errors.New("pengaturan perkenalan tentang kami tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateAboutIntroSettingRequest struct {
	Eyebrow            string `json:"eyebrow"`
	Title              string `json:"title"`
	ImageURL           string `json:"image_url"`
	ImageAlt           string `json:"image_alt"`
	ProgressTitle      string `json:"progress_title"`
	MissionTitle       string `json:"mission_title"`
	MissionDescription string `json:"mission_description"`
	VisionTitle        string `json:"vision_title"`
	VisionDescription  string `json:"vision_description"`
	CTAButtonText      string `json:"cta_button_text"`
	CTAButtonURL       string `json:"cta_button_url"`
	IsActive           *bool  `json:"is_active"`
}

type UpdateAboutIntroSettingRequest struct {
	Eyebrow            string `json:"eyebrow"`
	Title              string `json:"title"`
	ImageURL           string `json:"image_url"`
	ImageAlt           string `json:"image_alt"`
	ProgressTitle      string `json:"progress_title"`
	MissionTitle       string `json:"mission_title"`
	MissionDescription string `json:"mission_description"`
	VisionTitle        string `json:"vision_title"`
	VisionDescription  string `json:"vision_description"`
	CTAButtonText      string `json:"cta_button_text"`
	CTAButtonURL       string `json:"cta_button_url"`
	IsActive           *bool  `json:"is_active"`
}

// PublicAboutIntroResponse merepresentasikan DTO response khusus untuk kebutuhan publik frontend.
// Memisahkan poin Misi dan Visi menjadi mission_points dan vision_points.
type PublicAboutIntroResponse struct {
	ID                 uint                     `json:"id"`
	Eyebrow            string                   `json:"eyebrow"`
	Title              string                   `json:"title"`
	ImageURL           string                   `json:"image_url"`
	ImageAlt           string                   `json:"image_alt"`
	ProgressTitle      string                   `json:"progress_title"`
	Metrics            []model.AboutIntroMetric `json:"metrics"`
	MissionTitle       string                   `json:"mission_title"`
	MissionDescription string                   `json:"mission_description"`
	MissionPoints      []model.AboutIntroPoint  `json:"mission_points"`
	VisionTitle        string                   `json:"vision_title"`
	VisionDescription  string                   `json:"vision_description"`
	VisionPoints       []model.AboutIntroPoint  `json:"vision_points"`
	CTAButtonText      string                   `json:"cta_button_text"`
	CTAButtonURL       string                   `json:"cta_button_url"`
	IsActive           bool                     `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type AboutIntroSettingService struct {
	repo *repository.AboutIntroSettingRepository
}

func NewAboutIntroSettingService(repo *repository.AboutIntroSettingRepository) *AboutIntroSettingService {
	return &AboutIntroSettingService{repo: repo}
}

// GetAll mengambil seluruh konfigurasi untuk dashboard admin.
func (s *AboutIntroSettingService) GetAll() ([]model.AboutIntroSetting, error) {
	settings, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if settings == nil {
		settings = []model.AboutIntroSetting{}
	}
	return settings, nil
}

// GetByID mengambil konfigurasi berdasarkan ID beserta metrics dan points.
func (s *AboutIntroSettingService) GetByID(id uint) (*model.AboutIntroSetting, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	setting, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, ErrAboutIntroSettingNotFound
	}
	return setting, nil
}

// GetActive mengambil satu konfigurasi aktif untuk kebutuhan publik di halaman Tentang.
// Mengembalikan nil, nil jika belum ada pengaturan aktif di database (tanpa fallback hardcoded).
func (s *AboutIntroSettingService) GetActive() (*model.AboutIntroSetting, error) {
	setting, err := s.repo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return setting, nil
}

// GetActivePublicResponse mengambil setting aktif dan memformatnya menjadi PublicAboutIntroResponse
// dengan mission_points dan vision_points yang terpisah secara terstruktur untuk frontend.
func (s *AboutIntroSettingService) GetActivePublicResponse() (*PublicAboutIntroResponse, error) {
	setting, err := s.repo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, nil
	}

	metrics := setting.Metrics
	if metrics == nil {
		metrics = []model.AboutIntroMetric{}
	}

	missionPoints := make([]model.AboutIntroPoint, 0)
	visionPoints := make([]model.AboutIntroPoint, 0)

	for _, pt := range setting.Points {
		if strings.EqualFold(pt.SectionType, "mission") {
			missionPoints = append(missionPoints, pt)
		} else if strings.EqualFold(pt.SectionType, "vision") {
			visionPoints = append(visionPoints, pt)
		}
	}

	resp := &PublicAboutIntroResponse{
		ID:                 setting.ID,
		Eyebrow:            setting.Eyebrow,
		Title:              setting.Title,
		ImageURL:           setting.ImageURL,
		ImageAlt:           setting.ImageAlt,
		ProgressTitle:      setting.ProgressTitle,
		Metrics:            metrics,
		MissionTitle:       setting.MissionTitle,
		MissionDescription: setting.MissionDescription,
		MissionPoints:      missionPoints,
		VisionTitle:        setting.VisionTitle,
		VisionDescription:  setting.VisionDescription,
		VisionPoints:       visionPoints,
		CTAButtonText:      setting.CTAButtonText,
		CTAButtonURL:       setting.CTAButtonURL,
		IsActive:           setting.IsActive,
	}

	return resp, nil
}

// Create membuat konfigurasi baru.
func (s *AboutIntroSettingService) Create(req CreateAboutIntroSettingRequest) (*model.AboutIntroSetting, error) {
	eyebrow := strings.TrimSpace(req.Eyebrow)
	if eyebrow == "" {
		return nil, newValidationError("eyebrow wajib diisi")
	}
	if utf8.RuneCountInString(eyebrow) > 150 {
		return nil, newValidationError("eyebrow maksimal 150 karakter")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("title wajib diisi")
	}
	if utf8.RuneCountInString(title) > 255 {
		return nil, newValidationError("title maksimal 255 karakter")
	}

	imageURL := strings.TrimSpace(req.ImageURL)
	if imageURL == "" {
		return nil, newValidationError("image_url wajib diisi")
	}
	if utf8.RuneCountInString(imageURL) > 500 {
		return nil, newValidationError("image_url maksimal 500 karakter")
	}

	imageAlt := strings.TrimSpace(req.ImageAlt)
	if utf8.RuneCountInString(imageAlt) > 255 {
		return nil, newValidationError("image_alt maksimal 255 karakter")
	}

	progressTitle := strings.TrimSpace(req.ProgressTitle)
	if progressTitle == "" {
		return nil, newValidationError("progress_title wajib diisi")
	}
	if utf8.RuneCountInString(progressTitle) > 150 {
		return nil, newValidationError("progress_title maksimal 150 karakter")
	}

	missionTitle := strings.TrimSpace(req.MissionTitle)
	if missionTitle == "" {
		return nil, newValidationError("mission_title wajib diisi")
	}
	if utf8.RuneCountInString(missionTitle) > 150 {
		return nil, newValidationError("mission_title maksimal 150 karakter")
	}

	missionDescription := strings.TrimSpace(req.MissionDescription)

	visionTitle := strings.TrimSpace(req.VisionTitle)
	if visionTitle == "" {
		return nil, newValidationError("vision_title wajib diisi")
	}
	if utf8.RuneCountInString(visionTitle) > 150 {
		return nil, newValidationError("vision_title maksimal 150 karakter")
	}

	visionDescription := strings.TrimSpace(req.VisionDescription)

	ctaButtonText := strings.TrimSpace(req.CTAButtonText)
	if ctaButtonText == "" {
		return nil, newValidationError("cta_button_text wajib diisi")
	}
	if utf8.RuneCountInString(ctaButtonText) > 100 {
		return nil, newValidationError("cta_button_text maksimal 100 karakter")
	}

	ctaButtonURL := strings.TrimSpace(req.CTAButtonURL)
	if ctaButtonURL == "" {
		return nil, newValidationError("cta_button_url wajib diisi")
	}
	if utf8.RuneCountInString(ctaButtonURL) > 500 {
		return nil, newValidationError("cta_button_url maksimal 500 karakter")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	setting := &model.AboutIntroSetting{
		Eyebrow:            eyebrow,
		Title:              title,
		ImageURL:           imageURL,
		ImageAlt:           imageAlt,
		ProgressTitle:      progressTitle,
		MissionTitle:       missionTitle,
		MissionDescription: missionDescription,
		VisionTitle:        visionTitle,
		VisionDescription:  visionDescription,
		CTAButtonText:      ctaButtonText,
		CTAButtonURL:       ctaButtonURL,
		IsActive:           isActive,
	}

	if err := s.repo.Create(setting); err != nil {
		return nil, fmt.Errorf("gagal membuat konfigurasi perkenalan tentang kami: %w", err)
	}

	return setting, nil
}

// Update memperbarui konfigurasi yang sudah ada.
func (s *AboutIntroSettingService) Update(id uint, req UpdateAboutIntroSettingRequest) (*model.AboutIntroSetting, error) {
	if id == 0 {
		return nil, newValidationError("ID tidak valid")
	}

	setting, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, ErrAboutIntroSettingNotFound
	}

	eyebrow := strings.TrimSpace(req.Eyebrow)
	if eyebrow == "" {
		return nil, newValidationError("eyebrow wajib diisi")
	}
	if utf8.RuneCountInString(eyebrow) > 150 {
		return nil, newValidationError("eyebrow maksimal 150 karakter")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("title wajib diisi")
	}
	if utf8.RuneCountInString(title) > 255 {
		return nil, newValidationError("title maksimal 255 karakter")
	}

	imageURL := strings.TrimSpace(req.ImageURL)
	if imageURL == "" {
		return nil, newValidationError("image_url wajib diisi")
	}
	if utf8.RuneCountInString(imageURL) > 500 {
		return nil, newValidationError("image_url maksimal 500 karakter")
	}

	imageAlt := strings.TrimSpace(req.ImageAlt)
	if utf8.RuneCountInString(imageAlt) > 255 {
		return nil, newValidationError("image_alt maksimal 255 karakter")
	}

	progressTitle := strings.TrimSpace(req.ProgressTitle)
	if progressTitle == "" {
		return nil, newValidationError("progress_title wajib diisi")
	}
	if utf8.RuneCountInString(progressTitle) > 150 {
		return nil, newValidationError("progress_title maksimal 150 karakter")
	}

	missionTitle := strings.TrimSpace(req.MissionTitle)
	if missionTitle == "" {
		return nil, newValidationError("mission_title wajib diisi")
	}
	if utf8.RuneCountInString(missionTitle) > 150 {
		return nil, newValidationError("mission_title maksimal 150 karakter")
	}

	missionDescription := strings.TrimSpace(req.MissionDescription)

	visionTitle := strings.TrimSpace(req.VisionTitle)
	if visionTitle == "" {
		return nil, newValidationError("vision_title wajib diisi")
	}
	if utf8.RuneCountInString(visionTitle) > 150 {
		return nil, newValidationError("vision_title maksimal 150 karakter")
	}

	visionDescription := strings.TrimSpace(req.VisionDescription)

	ctaButtonText := strings.TrimSpace(req.CTAButtonText)
	if ctaButtonText == "" {
		return nil, newValidationError("cta_button_text wajib diisi")
	}
	if utf8.RuneCountInString(ctaButtonText) > 100 {
		return nil, newValidationError("cta_button_text maksimal 100 karakter")
	}

	ctaButtonURL := strings.TrimSpace(req.CTAButtonURL)
	if ctaButtonURL == "" {
		return nil, newValidationError("cta_button_url wajib diisi")
	}
	if utf8.RuneCountInString(ctaButtonURL) > 500 {
		return nil, newValidationError("cta_button_url maksimal 500 karakter")
	}

	setting.Eyebrow = eyebrow
	setting.Title = title
	setting.ImageURL = imageURL
	setting.ImageAlt = imageAlt
	setting.ProgressTitle = progressTitle
	setting.MissionTitle = missionTitle
	setting.MissionDescription = missionDescription
	setting.VisionTitle = visionTitle
	setting.VisionDescription = visionDescription
	setting.CTAButtonText = ctaButtonText
	setting.CTAButtonURL = ctaButtonURL

	if req.IsActive != nil {
		setting.IsActive = *req.IsActive
	}

	if err := s.repo.Update(setting); err != nil {
		return nil, fmt.Errorf("gagal memperbarui konfigurasi perkenalan tentang kami: %w", err)
	}

	return setting, nil
}

// Delete menghapus konfigurasi berdasarkan ID.
func (s *AboutIntroSettingService) Delete(id uint) error {
	if id == 0 {
		return newValidationError("ID tidak valid")
	}

	setting, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return ErrAboutIntroSettingNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus konfigurasi perkenalan tentang kami: %w", err)
	}

	return nil
}

