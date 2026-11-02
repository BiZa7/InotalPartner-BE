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
	ErrTestimonialNotFound = errors.New("testimoni tidak ditemukan")
)

// ValidationError merepresentasikan kesalahan validasi input data.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func newValidationError(message string) error {
	return &ValidationError{Message: message}
}

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateTestimonialRequest struct {
	Name         string `json:"name"`
	Position     string `json:"position"`
	Company      string `json:"company"`
	Content      string `json:"content"`
	PhotoURL     string `json:"photo_url"`
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateTestimonialRequest struct {
	Name         string  `json:"name"`
	Position     *string `json:"position"`
	Company      *string `json:"company"`
	Content      string  `json:"content"`
	PhotoURL     *string `json:"photo_url"`
	DisplayOrder *int    `json:"display_order"`
	IsActive     *bool   `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type TestimonialService struct {
	repo *repository.TestimonialRepository
}

func NewTestimonialService(repo *repository.TestimonialRepository) *TestimonialService {
	return &TestimonialService{repo: repo}
}

// GetAllTestimonials mengambil seluruh testimoni (termasuk yang nonaktif).
func (s *TestimonialService) GetAllTestimonials() ([]model.Testimonial, error) {
	testimonials, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if testimonials == nil {
		testimonials = []model.Testimonial{}
	}
	return testimonials, nil
}

// GetActiveTestimonials mengambil testimoni yang aktif saja untuk kebutuhan publik.
func (s *TestimonialService) GetActiveTestimonials() ([]model.Testimonial, error) {
	testimonials, err := s.repo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if testimonials == nil {
		testimonials = []model.Testimonial{}
	}
	return testimonials, nil
}

// GetTestimonialByID mengambil detail testimoni berdasarkan ID.
func (s *TestimonialService) GetTestimonialByID(id uint) (*model.Testimonial, error) {
	testimonial, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if testimonial == nil {
		return nil, ErrTestimonialNotFound
	}
	return testimonial, nil
}

// CreateTestimonial membuat testimoni baru dengan validasi field wajib dan batas panjang.
func (s *TestimonialService) CreateTestimonial(req CreateTestimonialRequest) (*model.Testimonial, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, newValidationError("nama wajib diisi")
	}
	if utf8.RuneCountInString(name) > 150 {
		return nil, newValidationError("nama maksimal 150 karakter")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, newValidationError("isi testimoni wajib diisi")
	}

	position := strings.TrimSpace(req.Position)
	if utf8.RuneCountInString(position) > 100 {
		return nil, newValidationError("jabatan/posisi maksimal 100 karakter")
	}

	company := strings.TrimSpace(req.Company)
	if utf8.RuneCountInString(company) > 150 {
		return nil, newValidationError("nama perusahaan/instansi maksimal 150 karakter")
	}

	if req.DisplayOrder < 0 {
		return nil, newValidationError("display_order tidak boleh bernilai negatif")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	testimonial := &model.Testimonial{
		Name:         name,
		Position:     position,
		Company:      company,
		Content:      content,
		PhotoURL:     strings.TrimSpace(req.PhotoURL),
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	if err := s.repo.Create(testimonial); err != nil {
		return nil, fmt.Errorf("gagal menambahkan testimoni: %w", err)
	}

	return testimonial, nil
}

// UpdateTestimonial memperbarui data testimoni berdasarkan ID.
func (s *TestimonialService) UpdateTestimonial(id uint, req UpdateTestimonialRequest) (*model.Testimonial, error) {
	testimonial, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if testimonial == nil {
		return nil, ErrTestimonialNotFound
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, newValidationError("nama wajib diisi")
	}
	if utf8.RuneCountInString(name) > 150 {
		return nil, newValidationError("nama maksimal 150 karakter")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, newValidationError("isi testimoni wajib diisi")
	}

	targetDisplayOrder := testimonial.DisplayOrder
	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, newValidationError("display_order tidak boleh bernilai negatif")
		}
		targetDisplayOrder = *req.DisplayOrder
	}

	targetIsActive := testimonial.IsActive
	if req.IsActive != nil {
		targetIsActive = *req.IsActive
	}

	targetPosition := testimonial.Position
	if req.Position != nil {
		trimmedPos := strings.TrimSpace(*req.Position)
		if utf8.RuneCountInString(trimmedPos) > 100 {
			return nil, newValidationError("jabatan/posisi maksimal 100 karakter")
		}
		targetPosition = trimmedPos
	}

	targetCompany := testimonial.Company
	if req.Company != nil {
		trimmedComp := strings.TrimSpace(*req.Company)
		if utf8.RuneCountInString(trimmedComp) > 150 {
			return nil, newValidationError("nama perusahaan/instansi maksimal 150 karakter")
		}
		targetCompany = trimmedComp
	}

	targetPhotoURL := testimonial.PhotoURL
	if req.PhotoURL != nil {
		targetPhotoURL = strings.TrimSpace(*req.PhotoURL)
	}

	testimonial.Name = name
	testimonial.Position = targetPosition
	testimonial.Company = targetCompany
	testimonial.Content = content
	testimonial.PhotoURL = targetPhotoURL
	testimonial.DisplayOrder = targetDisplayOrder
	testimonial.IsActive = targetIsActive

	if err := s.repo.Update(testimonial); err != nil {
		return nil, fmt.Errorf("gagal memperbarui testimoni: %w", err)
	}

	return testimonial, nil
}

// DeleteTestimonial menghapus testimoni berdasarkan ID.
func (s *TestimonialService) DeleteTestimonial(id uint) error {
	testimonial, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if testimonial == nil {
		return ErrTestimonialNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus testimoni: %w", err)
	}

	return nil
}
