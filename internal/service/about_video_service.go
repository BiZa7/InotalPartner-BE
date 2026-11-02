package service

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateAboutVideoRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	VideoURL     string `json:"video_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateAboutVideoRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	VideoURL     string `json:"video_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	IsActive     *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type AboutVideoService struct {
	videoRepo *repository.AboutVideoRepository
}

func NewAboutVideoService(videoRepo *repository.AboutVideoRepository) *AboutVideoService {
	return &AboutVideoService{videoRepo: videoRepo}
}

// GetAllVideos mengambil seluruh video profil untuk kebutuhan dashboard admin (termasuk yang nonaktif).
func (s *AboutVideoService) GetAllVideos() ([]model.AboutVideo, error) {
	videos, err := s.videoRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if videos == nil {
		videos = []model.AboutVideo{}
	}
	return videos, nil
}

// GetPublicVideos mengambil daftar video profil yang aktif untuk halaman publik.
func (s *AboutVideoService) GetPublicVideos() ([]model.AboutVideo, error) {
	videos, err := s.videoRepo.FindActivePublic()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if videos == nil {
		videos = []model.AboutVideo{}
	}
	return videos, nil
}

// GetVideoByID mengambil detail satu video profil berdasarkan ID.
func (s *AboutVideoService) GetVideoByID(id uint) (*model.AboutVideo, error) {
	video, err := s.videoRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if video == nil {
		return nil, errors.New("video profil tidak ditemukan")
	}
	return video, nil
}

// CreateVideo membuat video profil baru dengan validasi judul dan URL YouTube.
func (s *AboutVideoService) CreateVideo(req CreateAboutVideoRequest) (*model.AboutVideo, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("judul video wajib diisi")
	}
	if len(title) > 255 {
		return nil, errors.New("judul video tidak boleh melebihi 255 karakter")
	}

	normalizedVideoURL, err := validateAndNormalizeYouTubeURL(req.VideoURL)
	if err != nil {
		return nil, err
	}

	thumbnailURL := strings.TrimSpace(req.ThumbnailURL)
	if thumbnailURL != "" {
		if err := validateURL(thumbnailURL); err != nil {
			return nil, fmt.Errorf("URL thumbnail tidak valid: %w", err)
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	video := &model.AboutVideo{
		Title:        title,
		Description:  strings.TrimSpace(req.Description),
		VideoURL:     normalizedVideoURL,
		ThumbnailURL: thumbnailURL,
		IsActive:     isActive,
	}

	if err := s.videoRepo.Create(video); err != nil {
		return nil, fmt.Errorf("gagal menambahkan video profil: %w", err)
	}

	return video, nil
}

// UpdateVideo memperbarui video profil berdasarkan ID.
func (s *AboutVideoService) UpdateVideo(id uint, req UpdateAboutVideoRequest) (*model.AboutVideo, error) {
	video, err := s.videoRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if video == nil {
		return nil, errors.New("video profil tidak ditemukan")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("judul video wajib diisi")
	}
	if len(title) > 255 {
		return nil, errors.New("judul video tidak boleh melebihi 255 karakter")
	}

	normalizedVideoURL, err := validateAndNormalizeYouTubeURL(req.VideoURL)
	if err != nil {
		return nil, err
	}

	thumbnailURL := strings.TrimSpace(req.ThumbnailURL)
	if thumbnailURL != "" {
		if err := validateURL(thumbnailURL); err != nil {
			return nil, fmt.Errorf("URL thumbnail tidak valid: %w", err)
		}
	}

	targetIsActive := video.IsActive
	if req.IsActive != nil {
		targetIsActive = *req.IsActive
	}

	video.Title = title
	video.Description = strings.TrimSpace(req.Description)
	video.VideoURL = normalizedVideoURL
	video.ThumbnailURL = thumbnailURL
	video.IsActive = targetIsActive

	if err := s.videoRepo.Update(video); err != nil {
		return nil, fmt.Errorf("gagal memperbarui video profil: %w", err)
	}

	return video, nil
}

// DeleteVideo menghapus video profil berdasarkan ID.
func (s *AboutVideoService) DeleteVideo(id uint) error {
	video, err := s.videoRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if video == nil {
		return errors.New("video profil tidak ditemukan")
	}

	if err := s.videoRepo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus video profil: %w", err)
	}

	return nil
}

// ─── Helpers Validasi URL ─────────────────────────────────────────────────────

var youtubeIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)

// validateAndNormalizeYouTubeURL memeriksa dan mengekstrak ID video YouTube yang aman,
// lalu mengembalikan URL standar yang sudah dinormalisasi.
func validateAndNormalizeYouTubeURL(rawURL string) (string, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", errors.New("URL video wajib diisi")
	}

	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return "", errors.New("format URL video tidak valid")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("skema URL video harus menggunakan http atau https")
	}

	host := strings.ToLower(parsed.Hostname())
	var videoID string

	switch {
	case host == "youtu.be" || strings.HasSuffix(host, ".youtu.be"):
		// Format: https://youtu.be/<video_id>
		pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(pathParts) > 0 && pathParts[0] != "" {
			videoID = pathParts[0]
		}
	case host == "youtube.com" || strings.HasSuffix(host, ".youtube.com"):
		path := parsed.Path
		switch {
		case strings.HasPrefix(path, "/watch"):
			// Format: https://www.youtube.com/watch?v=<video_id>
			videoID = parsed.Query().Get("v")
		case strings.HasPrefix(path, "/embed/"):
			// Format: https://www.youtube.com/embed/<video_id>
			parts := strings.Split(strings.TrimPrefix(path, "/embed/"), "/")
			if len(parts) > 0 {
				videoID = parts[0]
			}
		case strings.HasPrefix(path, "/shorts/"):
			// Format: https://www.youtube.com/shorts/<video_id>
			parts := strings.Split(strings.TrimPrefix(path, "/shorts/"), "/")
			if len(parts) > 0 {
				videoID = parts[0]
			}
		case strings.HasPrefix(path, "/live/"):
			// Format: https://www.youtube.com/live/<video_id>
			parts := strings.Split(strings.TrimPrefix(path, "/live/"), "/")
			if len(parts) > 0 {
				videoID = parts[0]
			}
		}
	default:
		return "", errors.New("URL video harus berasal dari YouTube (youtube.com atau youtu.be)")
	}

	if videoID == "" || !youtubeIDRegex.MatchString(videoID) {
		return "", errors.New("ID video YouTube tidak valid atau tidak ditemukan pada URL yang diberikan")
	}

	// Normalisasi menjadi URL YouTube standar yang aman
	normalized := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	return normalized, nil
}

// validateURL memvalidasi bahwa string merupakan URL yang valid dengan skema http atau https.
func validateURL(rawURL string) error {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return errors.New("format URL tidak valid")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("skema URL harus http atau https")
	}
	if parsed.Host == "" {
		return errors.New("host URL tidak boleh kosong")
	}
	return nil
}
