package service

import (
	"errors"
	"fmt"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateRoleRequest struct {
	Name        string `json:"name"         binding:"required,min=2,max=50"`
	DisplayName string `json:"display_name" binding:"omitempty,max=100"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"         binding:"omitempty,min=2,max=50"`
	DisplayName string `json:"display_name" binding:"omitempty,max=100"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

// GetRoles mengambil semua data role
func (s *RoleService) GetRoles() ([]model.Role, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return roles, nil
}

// GetRoleByID mengambil detail role berdasarkan ID
func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if role == nil {
		return nil, errors.New("role tidak ditemukan")
	}
	return role, nil
}

// CreateRole membuat role baru (selalu is_system = false)
func (s *RoleService) CreateRole(req CreateRoleRequest) (*model.Role, error) {
	// Validasi name unique
	existing, err := s.roleRepo.FindByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, errors.New("nama role sudah digunakan")
	}

	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Name
	}

	role := &model.Role{
		Name:        req.Name,
		DisplayName: displayName,
		IsSystem:    false, // Role buatan CRUD selalu is_system = false
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// UpdateRole memperbarui data role
func (s *RoleService) UpdateRole(id uint, req UpdateRoleRequest) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if role == nil {
		return nil, errors.New("role tidak ditemukan")
	}

	// Jika nama diubah, lakukan validasi
	if req.Name != "" && req.Name != role.Name {
		// System role tidak boleh diubah namanya
		if role.IsSystem {
			return nil, errors.New("nama role system tidak dapat diubah")
		}

		// Cek keunikan nama
		exists, err := s.roleRepo.NameExistsExcept(req.Name, id)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		if exists {
			return nil, errors.New("nama role sudah digunakan")
		}
		role.Name = req.Name
	}

	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return role, nil
}

// DeleteRole menghapus role berdasarkan ID
func (s *RoleService) DeleteRole(id uint) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if role == nil {
		return errors.New("role tidak ditemukan")
	}

	// Proteksi role system
	if role.IsSystem {
		return errors.New("role system tidak dapat dihapus")
	}

	// Cek apakah ada user yang masih menggunakan role ini
	hasUsers, err := s.roleRepo.HasUsers(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if hasUsers {
		return errors.New("role tidak dapat dihapus karena masih digunakan oleh user")
	}

	return s.roleRepo.Delete(id)
}
