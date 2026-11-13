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

// ─── DTOs Sesuai Kontrak Frontend docs/api-kartu-layanan-beranda.md ────────────

type HomepageServiceItemDTO struct {
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type HomepageServicesDocument struct {
	Eyebrow            string                   `json:"eyebrow"`
	Title              string                   `json:"title"`
	Description        string                   `json:"description"`
	BadgeValue         string                   `json:"badge_value"`
	BadgeLabel         string                   `json:"badge_label"`
	PhotoPrimaryURL    string                   `json:"photo_primary_url"`
	PhotoPrimaryAlt    string                   `json:"photo_primary_alt"`
	PhotoSecondaryURL  string                   `json:"photo_secondary_url"`
	PhotoSecondaryAlt  string                   `json:"photo_secondary_alt"`
	Services           []HomepageServiceItemDTO `json:"services"`
}

type HomepageServicesUpsertService struct {
	db          *gorm.DB
	settingRepo *repository.PartnershipSettingRepository
	imageRepo   *repository.PartnershipImageRepository
	serviceRepo *repository.PartnershipServiceRepository
}

func NewHomepageServicesUpsertService(
	db *gorm.DB,
	settingRepo *repository.PartnershipSettingRepository,
	imageRepo *repository.PartnershipImageRepository,
	serviceRepo *repository.PartnershipServiceRepository,
) *HomepageServicesUpsertService {
	return &HomepageServicesUpsertService{
		db:          db,
		settingRepo: settingRepo,
		imageRepo:   imageRepo,
		serviceRepo: serviceRepo,
	}
}

// GetDocument mengambil dokumen gabungan kartu layanan Beranda untuk admin maupun publik.
func (s *HomepageServicesUpsertService) GetDocument() (*HomepageServicesDocument, error) {
	setting, err := s.settingRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	images, err := s.imageRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	services, err := s.serviceRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Jika belum pernah ada setting yang dibuat, kembalikan nil
	if setting == nil && len(images) == 0 && len(services) == 0 {
		return nil, nil
	}

	doc := &HomepageServicesDocument{
		Eyebrow:     "Layanan Kami",
		Title:       "Partnership",
		Description: "",
		BadgeValue:  "7+",
		BadgeLabel:  "Partnership Aktif",
		Services:    make([]HomepageServiceItemDTO, 0),
	}

	if setting != nil {
		if setting.Eyebrow != "" {
			doc.Eyebrow = setting.Eyebrow
		}
		doc.Title = setting.Title
		doc.Description = setting.Description
		doc.BadgeValue = setting.BadgeValue
		doc.BadgeLabel = setting.BadgeLabel
	}

	if len(images) > 0 {
		doc.PhotoPrimaryURL = images[0].ImageURL
		doc.PhotoPrimaryAlt = images[0].AltText
	}
	if len(images) > 1 {
		doc.PhotoSecondaryURL = images[1].ImageURL
		doc.PhotoSecondaryAlt = images[1].AltText
	}

	for _, srv := range services {
		doc.Services = append(doc.Services, HomepageServiceItemDTO{
			Icon:        srv.Icon,
			Title:       srv.Title,
			Description: srv.Description,
		})
	}

	return doc, nil
}

// SaveDocument melakukan upsert menyeluruh (setting, 2 foto, 4 kartu layanan) dalam satu transaksi database.
func (s *HomepageServicesUpsertService) SaveDocument(req HomepageServicesDocument) error {
	eyebrow := strings.TrimSpace(req.Eyebrow)
	if eyebrow == "" {
		return newValidationError("eyebrow wajib diisi")
	}
	if utf8.RuneCountInString(eyebrow) > 60 {
		return newValidationError("eyebrow maksimal 60 karakter")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return newValidationError("title wajib diisi")
	}
	if utf8.RuneCountInString(title) > 80 {
		return newValidationError("title maksimal 80 karakter")
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		return newValidationError("description wajib diisi")
	}
	if utf8.RuneCountInString(description) > 500 {
		return newValidationError("description maksimal 500 karakter")
	}

	badgeValue := strings.TrimSpace(req.BadgeValue)
	if badgeValue == "" {
		return newValidationError("badge_value wajib diisi")
	}
	if utf8.RuneCountInString(badgeValue) > 8 {
		return newValidationError("badge_value maksimal 8 karakter")
	}

	badgeLabel := strings.TrimSpace(req.BadgeLabel)
	if badgeLabel == "" {
		return newValidationError("badge_label wajib diisi")
	}
	if utf8.RuneCountInString(badgeLabel) > 40 {
		return newValidationError("badge_label maksimal 40 karakter")
	}

	photoPrimaryURL := strings.TrimSpace(req.PhotoPrimaryURL)
	if photoPrimaryURL == "" {
		return newValidationError("photo_primary_url wajib diisi")
	}

	photoSecondaryURL := strings.TrimSpace(req.PhotoSecondaryURL)
	if photoSecondaryURL == "" {
		return newValidationError("photo_secondary_url wajib diisi")
	}

	if len(req.Services) != 4 {
		return newValidationError("services harus berjumlah tepat 4 kartu layanan")
	}

	for i, srv := range req.Services {
		t := strings.TrimSpace(srv.Title)
		d := strings.TrimSpace(srv.Description)
		ic := strings.TrimSpace(srv.Icon)
		if t == "" {
			return newValidationError(fmt.Sprintf("judul kartu layanan %d wajib diisi", i+1))
		}
		if d == "" {
			return newValidationError(fmt.Sprintf("deskripsi kartu layanan %d wajib diisi", i+1))
		}
		if ic == "" {
			return newValidationError(fmt.Sprintf("ikon kartu layanan %d wajib dipilih", i+1))
		}
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Setting: ambil setting aktif atau buat baru
		var setting model.PartnershipSetting
		err := tx.Where("is_active = ?", true).Order("id DESC").First(&setting).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Buat baru
			setting = model.PartnershipSetting{
				Eyebrow:     eyebrow,
				Title:       title,
				Description: description,
				BadgeValue:  badgeValue,
				BadgeLabel:  badgeLabel,
				IsActive:    true,
			}
			if err := tx.Create(&setting).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// Update setting aktif
			setting.Eyebrow = eyebrow
			setting.Title = title
			setting.Description = description
			setting.BadgeValue = badgeValue
			setting.BadgeLabel = badgeLabel
			if err := tx.Save(&setting).Error; err != nil {
				return err
			}
		}

		// 2. Images: bersihkan dan perbarui tepat 2 gambar
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.PartnershipImage{}).Error; err != nil {
			return err
		}

		img1 := model.PartnershipImage{
			ImageURL:     photoPrimaryURL,
			AltText:      strings.TrimSpace(req.PhotoPrimaryAlt),
			DisplayOrder: 1,
			IsActive:     true,
		}
		if err := tx.Create(&img1).Error; err != nil {
			return err
		}

		img2 := model.PartnershipImage{
			ImageURL:     photoSecondaryURL,
			AltText:      strings.TrimSpace(req.PhotoSecondaryAlt),
			DisplayOrder: 2,
			IsActive:     true,
		}
		if err := tx.Create(&img2).Error; err != nil {
			return err
		}

		// 3. Services: bersihkan dan perbarui 4 kartu layanan
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.PartnershipService{}).Error; err != nil {
			return err
		}

		for idx, srv := range req.Services {
			sModel := model.PartnershipService{
				Title:        strings.TrimSpace(srv.Title),
				Description:  strings.TrimSpace(srv.Description),
				Icon:         strings.TrimSpace(srv.Icon),
				DisplayOrder: idx + 1,
				IsActive:     true,
			}
			if err := tx.Create(&sModel).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

