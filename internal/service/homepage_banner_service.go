package service

import (
	"errors"
	"strings"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

type BannerButtonRequest struct {
	Label       string `json:"label"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type CreateBannerRequest struct {
	Title           string                `json:"title"`
	Subtitle        string                `json:"subtitle"`
	BackgroundImage string                `json:"background_image"`
	Order           int                   `json:"order"`
	IsActive        *bool                 `json:"is_active"`
	Buttons         []BannerButtonRequest `json:"buttons"`
}

type UpdateBannerRequest struct {
	Title           string                `json:"title"`
	Subtitle        string                `json:"subtitle"`
	BackgroundImage string                `json:"background_image"`
	Order           int                   `json:"order"`
	IsActive        *bool                 `json:"is_active"`
	Buttons         []BannerButtonRequest `json:"buttons"`
}

type HomepageBannerService struct {
	bannerRepo *repository.HomepageBannerRepository
}

func NewHomepageBannerService(bannerRepo *repository.HomepageBannerRepository) *HomepageBannerService {
	return &HomepageBannerService{bannerRepo: bannerRepo}
}

// GetAllBanners mengambil seluruh banner (untuk Admin)
func (s *HomepageBannerService) GetAllBanners() ([]model.HomepageBanner, error) {
	banners, err := s.bannerRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if banners == nil {
		banners = []model.HomepageBanner{}
	}
	return banners, nil
}

// GetActiveBanners mengambil seluruh banner aktif (untuk Public Carousel)
func (s *HomepageBannerService) GetActiveBanners() ([]model.HomepageBanner, error) {
	banners, err := s.bannerRepo.FindActive()
	if err != nil {
		return nil, err
	}
	if banners == nil {
		banners = []model.HomepageBanner{}
	}
	return banners, nil
}

// GetBannerByID mengambil detail satu banner berdasarkan ID
func (s *HomepageBannerService) GetBannerByID(id uint) (*model.HomepageBanner, error) {
	banner, err := s.bannerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if banner == nil {
		return nil, errors.New("banner tidak ditemukan")
	}
	return banner, nil
}

// CreateBanner membuat record HomepageBanner baru beserta relasi buttons
func (s *HomepageBannerService) CreateBanner(req CreateBannerRequest) (*model.HomepageBanner, error) {
	if req.Order < 1 {
		return nil, errors.New("nomor urut slider harus minimal 1")
	}

	orderExists, err := s.bannerRepo.OrderExists(req.Order, 0)
	if err != nil {
		return nil, err
	}
	if orderExists {
		return nil, errors.New("nomor urut slider sudah digunakan. Silakan pilih nomor urut lain")
	}

	title := strings.TrimSpace(req.Title)
	bgImage := strings.TrimSpace(req.BackgroundImage)

	// Banner baru tidak langsung tampil di Landing Page.
	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	buttons := make([]model.HomepageBannerButton, 0, len(req.Buttons))
	for _, btnReq := range req.Buttons {
		label := strings.TrimSpace(btnReq.Label)
		url := strings.TrimSpace(btnReq.URL)

		// Tombol hanya dibuat jika label DAN URL diisi.
		if label == "" || url == "" {
			continue
		}

		buttons = append(buttons, model.HomepageBannerButton{
			Label:       label,
			URL:         url,
			Description: strings.TrimSpace(btnReq.Description),
			Order:       btnReq.Order,
		})
	}

	banner := &model.HomepageBanner{
		Title:           title,
		Subtitle:        "",
		BackgroundImage: bgImage,
		Order:           req.Order,
		IsActive:        isActive,
	}

	if err := s.bannerRepo.Create(banner, buttons); err != nil {
		return nil, err
	}

	return s.bannerRepo.FindByID(banner.ID)
}

// UpdateBanner memperbarui data HomepageBanner dan mereplace list button miliknya
func (s *HomepageBannerService) UpdateBanner(id uint, req UpdateBannerRequest) (*model.HomepageBanner, error) {
	if req.Order < 1 {
		return nil, errors.New("nomor urut slider harus minimal 1")
	}

	banner, err := s.bannerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if banner == nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	orderExists, err := s.bannerRepo.OrderExists(req.Order, id)
	if err != nil {
		return nil, err
	}
	if orderExists {
		return nil, errors.New("nomor urut slider sudah digunakan. Silakan pilih nomor urut lain")
	}

	banner.Title = strings.TrimSpace(req.Title)
	banner.Subtitle = ""
	banner.BackgroundImage = strings.TrimSpace(req.BackgroundImage)
	banner.Order = req.Order

	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}

	buttons := make([]model.HomepageBannerButton, 0, len(req.Buttons))
	for _, btnReq := range req.Buttons {
		label := strings.TrimSpace(btnReq.Label)
		url := strings.TrimSpace(btnReq.URL)

		if label == "" || url == "" {
			continue
		}

		buttons = append(buttons, model.HomepageBannerButton{
			Label:       label,
			URL:         url,
			Description: strings.TrimSpace(btnReq.Description),
			Order:       btnReq.Order,
		})
	}

	if err := s.bannerRepo.Update(banner, buttons); err != nil {
		return nil, err
	}

	return s.bannerRepo.FindByID(id)
}

// DeleteBanner menghapus HomepageBanner berdasarkan ID
func (s *HomepageBannerService) DeleteBanner(id uint) error {
	banner, err := s.bannerRepo.FindByID(id)
	if err != nil {
		return err
	}
	if banner == nil {
		return errors.New("banner tidak ditemukan")
	}
	return s.bannerRepo.Delete(id)
}
