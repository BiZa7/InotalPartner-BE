package service

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs Sesuai Kontrak Frontend docs/api-intro-tentang.md ───────────────────

type AboutIntroProgressDTO struct {
	Label   string `json:"label"`
	Percent int    `json:"percent"`
}

type AboutIntroCardDTO struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Items       []string `json:"items"`
}

type AboutIntroDocument struct {
	Eyebrow        string                  `json:"eyebrow"`
	Headline       string                  `json:"headline"`
	HeadlineAccent string                  `json:"headline_accent"`
	PhotoURL       string                  `json:"photo_url"`
	PhotoAlt       string                  `json:"photo_alt"`
	ProgressTitle  string                  `json:"progress_title"`
	Progress       []AboutIntroProgressDTO `json:"progress"`
	Mission        AboutIntroCardDTO       `json:"mission"`
	Vision         AboutIntroCardDTO       `json:"vision"`
	CTALabel       string                  `json:"cta_label"`
	CTATarget      string                  `json:"cta_target"`
}

type AboutIntroUpsertService struct {
	db          *gorm.DB
	settingRepo *repository.AboutIntroSettingRepository
	metricRepo  *repository.AboutIntroMetricRepository
	pointRepo   *repository.AboutIntroPointRepository
}

func NewAboutIntroUpsertService(
	db *gorm.DB,
	settingRepo *repository.AboutIntroSettingRepository,
	metricRepo *repository.AboutIntroMetricRepository,
	pointRepo *repository.AboutIntroPointRepository,
) *AboutIntroUpsertService {
	return &AboutIntroUpsertService{
		db:          db,
		settingRepo: settingRepo,
		metricRepo:  metricRepo,
		pointRepo:   pointRepo,
	}
}

// GetDocument mengambil dokumen terpadu Intro & Visi Misi sesuai kontrak frontend.
func (s *AboutIntroUpsertService) GetDocument() (*AboutIntroDocument, error) {
	setting, err := s.settingRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if setting == nil {
		return nil, nil
	}

	// Parsing Headline dan HeadlineAccent dari setting.Title
	headline := setting.Title
	headlineAccent := ""
	if parts := strings.Split(setting.Title, " | "); len(parts) == 2 {
		headline = parts[0]
		headlineAccent = parts[1]
	}

	progressList := make([]AboutIntroProgressDTO, 0)
	for _, m := range setting.Metrics {
		progressList = append(progressList, AboutIntroProgressDTO{
			Label:   m.Label,
			Percent: m.Value,
		})
	}

	missionItems := make([]string, 0)
	visionItems := make([]string, 0)

	for _, p := range setting.Points {
		if strings.EqualFold(p.SectionType, "mission") {
			missionItems = append(missionItems, p.Content)
		} else if strings.EqualFold(p.SectionType, "vision") {
			visionItems = append(visionItems, p.Content)
		}
	}

	doc := &AboutIntroDocument{
		Eyebrow:        setting.Eyebrow,
		Headline:       headline,
		HeadlineAccent: headlineAccent,
		PhotoURL:       setting.ImageURL,
		PhotoAlt:       setting.ImageAlt,
		ProgressTitle:  setting.ProgressTitle,
		Progress:       progressList,
		Mission: AboutIntroCardDTO{
			Title:       setting.MissionTitle,
			Description: setting.MissionDescription,
			Items:       missionItems,
		},
		Vision: AboutIntroCardDTO{
			Title:       setting.VisionTitle,
			Description: setting.VisionDescription,
			Items:       visionItems,
		},
		CTALabel:  setting.CTAButtonText,
		CTATarget: setting.CTAButtonURL,
	}

	return doc, nil
}

// SaveDocument melakukan upsert dokumen terpadu Intro & Visi Misi dalam satu transaksi database atomik.
func (s *AboutIntroUpsertService) SaveDocument(req AboutIntroDocument) error {
	eyebrow := strings.TrimSpace(req.Eyebrow)
	if eyebrow == "" {
		return newValidationError("eyebrow wajib diisi")
	}
	if utf8.RuneCountInString(eyebrow) > 60 {
		return newValidationError("eyebrow maksimal 60 karakter")
	}

	headline := strings.TrimSpace(req.Headline)
	if headline == "" {
		return newValidationError("headline wajib diisi")
	}
	if utf8.RuneCountInString(headline) > 160 {
		return newValidationError("headline maksimal 160 karakter")
	}

	headlineAccent := strings.TrimSpace(req.HeadlineAccent)
	if utf8.RuneCountInString(headlineAccent) > 80 {
		return newValidationError("headline_accent maksimal 80 karakter")
	}

	combinedTitle := headline
	if headlineAccent != "" {
		combinedTitle = headline + " | " + headlineAccent
	}

	photoURL := strings.TrimSpace(req.PhotoURL)
	if photoURL == "" {
		return newValidationError("photo_url wajib diisi")
	}

	photoAlt := strings.TrimSpace(req.PhotoAlt)
	if photoAlt == "" {
		return newValidationError("photo_alt wajib diisi")
	}
	if utf8.RuneCountInString(photoAlt) > 160 {
		return newValidationError("photo_alt maksimal 160 karakter")
	}

	progressTitle := strings.TrimSpace(req.ProgressTitle)
	if progressTitle == "" {
		return newValidationError("progress_title wajib diisi")
	}
	if utf8.RuneCountInString(progressTitle) > 60 {
		return newValidationError("progress_title maksimal 60 karakter")
	}

	if len(req.Progress) != 2 {
		return newValidationError("progress harus berisi tepat 2 baris")
	}
	for i, p := range req.Progress {
		if strings.TrimSpace(p.Label) == "" {
			return newValidationError(fmt.Sprintf("label progress %d wajib diisi", i+1))
		}
		if p.Percent < 0 || p.Percent > 100 {
			return newValidationError(fmt.Sprintf("persentase progress %d harus antara 0-100", i+1))
		}
	}

	// Validasi Mission
	missionTitle := strings.TrimSpace(req.Mission.Title)
	if missionTitle == "" {
		return newValidationError("mission.title wajib diisi")
	}
	if utf8.RuneCountInString(missionTitle) > 60 {
		return newValidationError("mission.title maksimal 60 karakter")
	}
	if utf8.RuneCountInString(req.Mission.Description) > 300 {
		return newValidationError("mission.description maksimal 300 karakter")
	}
	if len(req.Mission.Items) != 3 {
		return newValidationError("poin misi harus berisi tepat 3 item")
	}
	for i, item := range req.Mission.Items {
		if strings.TrimSpace(item) == "" {
			return newValidationError(fmt.Sprintf("poin misi %d wajib diisi", i+1))
		}
	}

	// Validasi Vision
	visionTitle := strings.TrimSpace(req.Vision.Title)
	if visionTitle == "" {
		return newValidationError("vision.title wajib diisi")
	}
	if utf8.RuneCountInString(visionTitle) > 60 {
		return newValidationError("vision.title maksimal 60 karakter")
	}
	if utf8.RuneCountInString(req.Vision.Description) > 300 {
		return newValidationError("vision.description maksimal 300 karakter")
	}
	if len(req.Vision.Items) != 3 {
		return newValidationError("poin visi harus berisi tepat 3 item")
	}
	for i, item := range req.Vision.Items {
		if strings.TrimSpace(item) == "" {
			return newValidationError(fmt.Sprintf("poin visi %d wajib diisi", i+1))
		}
	}

	ctaLabel := strings.TrimSpace(req.CTALabel)
	if ctaLabel == "" {
		return newValidationError("cta_label wajib diisi")
	}
	if utf8.RuneCountInString(ctaLabel) > 60 {
		return newValidationError("cta_label maksimal 60 karakter")
	}

	ctaTarget := strings.TrimSpace(req.CTATarget)
	if ctaTarget == "" {
		return newValidationError("cta_target wajib diisi")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var setting model.AboutIntroSetting
		err := tx.Where("is_active = ?", true).Order("id DESC").First(&setting).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			setting = model.AboutIntroSetting{
				Eyebrow:            eyebrow,
				Title:              combinedTitle,
				ImageURL:           photoURL,
				ImageAlt:           photoAlt,
				ProgressTitle:      progressTitle,
				MissionTitle:       missionTitle,
				MissionDescription: strings.TrimSpace(req.Mission.Description),
				VisionTitle:        visionTitle,
				VisionDescription:  strings.TrimSpace(req.Vision.Description),
				CTAButtonText:      ctaLabel,
				CTAButtonURL:       ctaTarget,
				IsActive:           true,
			}
			if err := tx.Create(&setting).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			setting.Eyebrow = eyebrow
			setting.Title = combinedTitle
			setting.ImageURL = photoURL
			setting.ImageAlt = photoAlt
			setting.ProgressTitle = progressTitle
			setting.MissionTitle = missionTitle
			setting.MissionDescription = strings.TrimSpace(req.Mission.Description)
			setting.VisionTitle = visionTitle
			setting.VisionDescription = strings.TrimSpace(req.Vision.Description)
			setting.CTAButtonText = ctaLabel
			setting.CTAButtonURL = ctaTarget
			if err := tx.Save(&setting).Error; err != nil {
				return err
			}
		}

		// Bersihkan metric lama dan simpan metric baru
		if err := tx.Where("setting_id = ?", setting.ID).Delete(&model.AboutIntroMetric{}).Error; err != nil {
			return err
		}

		for idx, p := range req.Progress {
			metric := model.AboutIntroMetric{
				SettingID:    setting.ID,
				Label:        strings.TrimSpace(p.Label),
				Value:        p.Percent,
				DisplayOrder: idx + 1,
				IsActive:     true,
			}
			if err := tx.Create(&metric).Error; err != nil {
				return err
			}
		}

		// Bersihkan points lama dan simpan points baru (misi & visi)
		if err := tx.Where("setting_id = ?", setting.ID).Delete(&model.AboutIntroPoint{}).Error; err != nil {
			return err
		}

		for idx, item := range req.Mission.Items {
			pt := model.AboutIntroPoint{
				SettingID:    setting.ID,
				SectionType:  "mission",
				Content:      strings.TrimSpace(item),
				DisplayOrder: idx + 1,
				IsActive:     true,
			}
			if err := tx.Create(&pt).Error; err != nil {
				return err
			}
		}

		for idx, item := range req.Vision.Items {
			pt := model.AboutIntroPoint{
				SettingID:    setting.ID,
				SectionType:  "vision",
				Content:      strings.TrimSpace(item),
				DisplayOrder: idx + 1,
				IsActive:     true,
			}
			if err := tx.Create(&pt).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

