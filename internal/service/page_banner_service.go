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
	ErrPageBannerNotFound = errors.New("banner halaman tidak ditemukan")
)

// ─── Supported Page Types ─────────────────────────────────────────────────────

var validPageTypes = map[string]bool{
	"about":   true,
	"event":   true,
	"news":    true,
	"article": true,
	"contact": true,
}

// NormalizeAndValidatePageType memvalidasi dan mengembalikan format lowercase dari page_type.
// Mendukung format kanonikal dan alias dwibahasa:
// - tentang <-> about
// - berita <-> news
// - artikel <-> article
// - kontak <-> contact
// - event
func NormalizeAndValidatePageType(raw string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))

	// Mapping alias bahasa Indonesia / alternatif ke canonical page_type backend
	aliasMap := map[string]string{
		"tentang": "about",
		"about":   "about",
		"berita":  "news",
		"news":    "news",
		"artikel": "article",
		"article": "article",
		"kontak":  "contact",
		"contact": "contact",
		"event":   "event",
	}

	if canonical, isAlias := aliasMap[normalized]; isAlias {
		normalized = canonical
	}

	if !validPageTypes[normalized] {
		return "", newValidationError("page_type tidak valid. Pilihan yang didukung: about/tentang, event, news/berita, article/artikel, contact/kontak")
	}
	return normalized, nil
}

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreatePageBannerRequest struct {
	ImageURL    string `json:"image_url"`
	Subtitle    string `json:"subtitle"`
	Label       string `json:"label"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type UpdatePageBannerRequest struct {
	ImageURL    string `json:"image_url"`
	Subtitle    string `json:"subtitle"`
	Label       string `json:"label"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type PageBannerService struct {
	bannerRepo *repository.PageBannerRepository
}

func NewPageBannerService(bannerRepo *repository.PageBannerRepository) *PageBannerService {
	return &PageBannerService{bannerRepo: bannerRepo}
}

// GetBannerByPageType mengambil data banner untuk admin berdasarkan page_type.
func (s *PageBannerService) GetBannerByPageType(rawPageType string) (*model.PageBanner, error) {
	pageType, err := NormalizeAndValidatePageType(rawPageType)
	if err != nil {
		return nil, err
	}

	banner, err := s.bannerRepo.FindByPageType(pageType)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if banner == nil {
		return nil, ErrPageBannerNotFound
	}
	return banner, nil
}

// GetPublicBannerByPageType mengambil data banner aktif saja untuk pengunjung halaman publik.
func (s *PageBannerService) GetPublicBannerByPageType(rawPageType string) (*model.PageBanner, error) {
	pageType, err := NormalizeAndValidatePageType(rawPageType)
	if err != nil {
		return nil, err
	}

	banner, err := s.bannerRepo.FindActiveByPageType(pageType)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if banner == nil {
		return nil, ErrPageBannerNotFound
	}
	return banner, nil
}

// CreateBanner membuat data banner baru untuk page_type tertentu.
func (s *PageBannerService) CreateBanner(rawPageType string, req CreatePageBannerRequest) (*model.PageBanner, error) {
	pageType, err := NormalizeAndValidatePageType(rawPageType)
	if err != nil {
		return nil, err
	}

	// Periksa apakah banner untuk halaman tersebut sudah ada (karena 1 halaman = 1 banner)
	existing, err := s.bannerRepo.FindByPageType(pageType)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, newValidationError(fmt.Sprintf("banner untuk halaman '%s' sudah ada. Silakan gunakan update (PUT)", pageType))
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("judul banner (title) wajib diisi")
	}
	if utf8.RuneCountInString(title) > 255 {
		return nil, newValidationError("judul banner (title) maksimal 255 karakter")
	}

	subtitle := strings.TrimSpace(req.Subtitle)
	if subtitle == "" && strings.TrimSpace(req.Label) != "" {
		subtitle = strings.TrimSpace(req.Label)
	}
	if utf8.RuneCountInString(subtitle) > 255 {
		return nil, newValidationError("subtitle banner maksimal 255 karakter")
	}

	imageURL := strings.TrimSpace(req.ImageURL)
	description := strings.TrimSpace(req.Description)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	banner := &model.PageBanner{
		PageType:    pageType,
		ImageURL:    imageURL,
		Subtitle:    subtitle,
		Title:       title,
		Description: description,
		IsActive:    isActive,
	}

	if err := s.bannerRepo.Create(banner); err != nil {
		return nil, fmt.Errorf("gagal membuat banner: %w", err)
	}

	return banner, nil
}

// UpdateBanner memperbarui data banner untuk page_type tertentu.
// Jika belum ada, akan mengembalikan error not found agar admin membuat melalui POST terlebih dahulu.
func (s *PageBannerService) UpdateBanner(rawPageType string, req UpdatePageBannerRequest) (*model.PageBanner, error) {
	pageType, err := NormalizeAndValidatePageType(rawPageType)
	if err != nil {
		return nil, err
	}

	banner, err := s.bannerRepo.FindByPageType(pageType)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if banner == nil {
		return nil, ErrPageBannerNotFound
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, newValidationError("judul banner (title) wajib diisi")
	}
	if utf8.RuneCountInString(title) > 255 {
		return nil, newValidationError("judul banner (title) maksimal 255 karakter")
	}

	subtitle := strings.TrimSpace(req.Subtitle)
	if subtitle == "" && strings.TrimSpace(req.Label) != "" {
		subtitle = strings.TrimSpace(req.Label)
	}
	if utf8.RuneCountInString(subtitle) > 255 {
		return nil, newValidationError("subtitle banner maksimal 255 karakter")
	}

	imageURL := strings.TrimSpace(req.ImageURL)
	description := strings.TrimSpace(req.Description)

	isActive := banner.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	banner.ImageURL = imageURL
	banner.Subtitle = subtitle
	banner.Title = title
	banner.Description = description
	banner.IsActive = isActive

	if err := s.bannerRepo.Update(banner); err != nil {
		return nil, fmt.Errorf("gagal memperbarui banner: %w", err)
	}

	return banner, nil
}

// DeleteBanner menghapus banner untuk page_type tertentu.
func (s *PageBannerService) DeleteBanner(rawPageType string) error {
	pageType, err := NormalizeAndValidatePageType(rawPageType)
	if err != nil {
		return err
	}

	banner, err := s.bannerRepo.FindByPageType(pageType)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if banner == nil {
		return ErrPageBannerNotFound
	}

	if err := s.bannerRepo.Delete(banner.ID); err != nil {
		return fmt.Errorf("gagal menghapus banner: %w", err)
	}

	return nil
}

