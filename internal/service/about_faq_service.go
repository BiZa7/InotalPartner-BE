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
	ErrAboutFAQNotFound = errors.New("FAQ tidak ditemukan")
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateAboutFAQRequest struct {
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateAboutFAQRequest struct {
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type AboutFAQService struct {
	faqRepo *repository.AboutFAQRepository
}

func NewAboutFAQService(faqRepo *repository.AboutFAQRepository) *AboutFAQService {
	return &AboutFAQService{faqRepo: faqRepo}
}

// GetAllFAQs mengambil seluruh FAQ untuk dashboard admin.
func (s *AboutFAQService) GetAllFAQs() ([]model.AboutFAQ, error) {
	faqs, err := s.faqRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if faqs == nil {
		faqs = []model.AboutFAQ{}
	}
	return faqs, nil
}

// GetPublicFAQs mengambil FAQ aktif saja untuk halaman publik.
func (s *AboutFAQService) GetPublicFAQs() ([]model.AboutFAQ, error) {
	faqs, err := s.faqRepo.FindActivePublic()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if faqs == nil {
		faqs = []model.AboutFAQ{}
	}
	return faqs, nil
}

// GetFAQByID mengambil detail FAQ berdasarkan ID.
func (s *AboutFAQService) GetFAQByID(id uint) (*model.AboutFAQ, error) {
	faq, err := s.faqRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if faq == nil {
		return nil, ErrAboutFAQNotFound
	}
	return faq, nil
}

// CreateFAQ membuat FAQ baru dengan validasi field question, answer, display_order, dan is_active.
func (s *AboutFAQService) CreateFAQ(req CreateAboutFAQRequest) (*model.AboutFAQ, error) {
	question := strings.TrimSpace(req.Question)
	if question == "" {
		return nil, newValidationError("pertanyaan wajib diisi")
	}
	if utf8.RuneCountInString(question) > 255 {
		return nil, newValidationError("pertanyaan maksimal 255 karakter")
	}

	answer := strings.TrimSpace(req.Answer)
	if answer == "" {
		return nil, newValidationError("jawaban wajib diisi")
	}

	if req.DisplayOrder < 0 {
		return nil, newValidationError("display_order tidak boleh bernilai negatif")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	faq := &model.AboutFAQ{
		Question:     question,
		Answer:       answer,
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	if err := s.faqRepo.Create(faq); err != nil {
		return nil, fmt.Errorf("gagal menambahkan FAQ: %w", err)
	}

	return faq, nil
}

// UpdateFAQ memperbarui FAQ dengan validasi.
func (s *AboutFAQService) UpdateFAQ(id uint, req UpdateAboutFAQRequest) (*model.AboutFAQ, error) {
	faq, err := s.faqRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if faq == nil {
		return nil, ErrAboutFAQNotFound
	}

	question := strings.TrimSpace(req.Question)
	if question == "" {
		return nil, newValidationError("pertanyaan wajib diisi")
	}
	if utf8.RuneCountInString(question) > 255 {
		return nil, newValidationError("pertanyaan maksimal 255 karakter")
	}

	answer := strings.TrimSpace(req.Answer)
	if answer == "" {
		return nil, newValidationError("jawaban wajib diisi")
	}

	targetDisplayOrder := faq.DisplayOrder
	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, newValidationError("display_order tidak boleh bernilai negatif")
		}
		targetDisplayOrder = *req.DisplayOrder
	}

	targetIsActive := faq.IsActive
	if req.IsActive != nil {
		targetIsActive = *req.IsActive
	}

	faq.Question = question
	faq.Answer = answer
	faq.DisplayOrder = targetDisplayOrder
	faq.IsActive = targetIsActive

	if err := s.faqRepo.Update(faq); err != nil {
		return nil, fmt.Errorf("gagal memperbarui FAQ: %w", err)
	}

	return faq, nil
}

// DeleteFAQ menghapus FAQ berdasarkan ID.
func (s *AboutFAQService) DeleteFAQ(id uint) error {
	faq, err := s.faqRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if faq == nil {
		return ErrAboutFAQNotFound
	}

	if err := s.faqRepo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus FAQ: %w", err)
	}

	return nil
}

