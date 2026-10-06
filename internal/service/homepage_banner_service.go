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
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("title wajib diisi")
	}

	subtitle := strings.TrimSpace(req.Subtitle)
	if subtitle == "" {
		return nil, errors.New("subtitle wajib diisi")
	}

	bgImage := strings.TrimSpace(req.BackgroundImage)
	if bgImage == "" {
		return nil, errors.New("background_image wajib diisi")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	buttons := make([]model.HomepageBannerButton, 0, len(req.Buttons))
	for _, btnReq := range req.Buttons {
		btnLabel := strings.TrimSpace(btnReq.Label)
		if btnLabel == "" {
			return nil, errors.New("label button wajib diisi")
		}
		btnURL := strings.TrimSpace(btnReq.URL)
		if btnURL == "" {
			return nil, errors.New("url button wajib diisi")
		}

		buttons = append(buttons, model.HomepageBannerButton{
			Label:       btnLabel,
			URL:         btnURL,
			Description: strings.TrimSpace(btnReq.Description),
			Order:       btnReq.Order,
		})
	}

	banner := &model.HomepageBanner{
		Title:           title,
		Subtitle:        subtitle,
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
	banner, err := s.bannerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if banner == nil {
		return nil, errors.New("banner tidak ditemukan")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("title wajib diisi")
	}

	subtitle := strings.TrimSpace(req.Subtitle)
	if subtitle == "" {
		return nil, errors.New("subtitle wajib diisi")
	}

	bgImage := strings.TrimSpace(req.BackgroundImage)
	if bgImage == "" {
		return nil, errors.New("background_image wajib diisi")
	}

	banner.Title = title
	banner.Subtitle = subtitle
	banner.BackgroundImage = bgImage
	banner.Order = req.Order

	if req.IsActive != nil {
		banner.IsActive = *req.IsActive
	}

	buttons := make([]model.HomepageBannerButton, 0, len(req.Buttons))
	for _, btnReq := range req.Buttons {
		btnLabel := strings.TrimSpace(btnReq.Label)
		if btnLabel == "" {
			return nil, errors.New("label button wajib diisi")
		}
		btnURL := strings.TrimSpace(btnReq.URL)
		if btnURL == "" {
			return nil, errors.New("url button wajib diisi")
		}

		buttons = append(buttons, model.HomepageBannerButton{
			Label:       btnLabel,
			URL:         btnURL,
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
