package service

import (
	"errors"
	"fmt"
	"strings"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateTeamMemberRequest struct {
	Name         string `json:"name"`
	PhotoURL     string `json:"photo_url"`
	Position     string `json:"position"`
	DisplayOrder int    `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateTeamMemberRequest struct {
	Name         string `json:"name"`
	PhotoURL     string `json:"photo_url"`
	Position     string `json:"position"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     *bool  `json:"is_active"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type TeamMemberService struct {
	memberRepo *repository.TeamMemberRepository
}

func NewTeamMemberService(memberRepo *repository.TeamMemberRepository) *TeamMemberService {
	return &TeamMemberService{memberRepo: memberRepo}
}

// GetAllTeamMembers mengambil seluruh anggota tim untuk kebutuhan pengelolaan admin
// (termasuk nonaktif maupun data yang belum lengkap).
func (s *TeamMemberService) GetAllTeamMembers() ([]model.TeamMember, error) {
	members, err := s.memberRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if members == nil {
		members = []model.TeamMember{}
	}
	return members, nil
}

// GetPublicTeamMembers mengambil anggota tim aktif dengan nama dan foto yang terisi lengkap
// untuk kebutuhan publik, diurutkan berdasarkan display_order ASC, id ASC.
func (s *TeamMemberService) GetPublicTeamMembers() ([]model.TeamMember, error) {
	members, err := s.memberRepo.FindActivePublic()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if members == nil {
		members = []model.TeamMember{}
	}
	return members, nil
}

// GetTeamMemberByID mengambil detail satu anggota tim berdasarkan ID.
func (s *TeamMemberService) GetTeamMemberByID(id uint) (*model.TeamMember, error) {
	member, err := s.memberRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if member == nil {
		return nil, errors.New("anggota tim tidak ditemukan")
	}
	return member, nil
}

// CreateTeamMember membuat record anggota tim baru.
// Nama dan foto boleh kosong saat persiapan (draft), position bersifat opsional.
// Jumlah anggota tidak dibatasi.
func (s *TeamMemberService) CreateTeamMember(req CreateTeamMemberRequest) (*model.TeamMember, error) {
	if req.DisplayOrder < 0 {
		return nil, errors.New("display_order tidak boleh bernilai negatif")
	}

	// Periksa apakah nomor urut tampil sudah digunakan oleh anggota lain
	existing, err := s.memberRepo.FindByDisplayOrder(req.DisplayOrder, 0)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, errors.New("display_order sudah digunakan")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	member := &model.TeamMember{
		Name:         strings.TrimSpace(req.Name),
		PhotoURL:     strings.TrimSpace(req.PhotoURL),
		Position:     strings.TrimSpace(req.Position),
		DisplayOrder: req.DisplayOrder,
		IsActive:     isActive,
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, fmt.Errorf("gagal menambahkan anggota tim: %w", err)
	}

	return member, nil
}

// UpdateTeamMember memperbarui data anggota tim berdasarkan ID.
func (s *TeamMemberService) UpdateTeamMember(id uint, req UpdateTeamMemberRequest) (*model.TeamMember, error) {
	member, err := s.memberRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if member == nil {
		return nil, errors.New("anggota tim tidak ditemukan")
	}

	targetDisplayOrder := member.DisplayOrder
	if req.DisplayOrder != nil {
		if *req.DisplayOrder < 0 {
			return nil, errors.New("display_order tidak boleh bernilai negatif")
		}
		targetDisplayOrder = *req.DisplayOrder
	}

	targetIsActive := member.IsActive
	if req.IsActive != nil {
		targetIsActive = *req.IsActive
	}

	// Periksa keunikan display_order terhadap anggota aktif lain dengan mengecualikan ID anggota yang sedang diedit
	existing, err := s.memberRepo.FindByDisplayOrder(targetDisplayOrder, id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, errors.New("display_order sudah digunakan")
	}

	member.DisplayOrder = targetDisplayOrder
	member.Name = strings.TrimSpace(req.Name)
	member.PhotoURL = strings.TrimSpace(req.PhotoURL)
	member.Position = strings.TrimSpace(req.Position)
	member.IsActive = targetIsActive

	if err := s.memberRepo.Update(member); err != nil {
		return nil, fmt.Errorf("gagal memperbarui anggota tim: %w", err)
	}

	return member, nil
}

// DeleteTeamMember menghapus anggota tim berdasarkan ID.
func (s *TeamMemberService) DeleteTeamMember(id uint) error {
	member, err := s.memberRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if member == nil {
		return errors.New("anggota tim tidak ditemukan")
	}

	if err := s.memberRepo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus anggota tim: %w", err)
	}

	return nil
}
