package service

import (
    "errors"
    "fmt"

    "golang.org/x/crypto/bcrypt"

    "inotal-be/internal/model"
    "inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

type CreateUserRequest struct {
    FullName string `json:"full_name" binding:"required,min=2,max=150"`
    Email    string `json:"email"     binding:"required,email"`
    Password string `json:"password"  binding:"required,min=8"`
    Company  string `json:"company"   binding:"omitempty,max=200"`
    Role     string `json:"role"      binding:"required,oneof=super_admin admin operator guest"`
}

type UpdateUserRequest struct {
    FullName string `json:"full_name" binding:"omitempty,min=2,max=150"`
    Email    string `json:"email"     binding:"omitempty,email"`
    Company  string `json:"company"   binding:"omitempty,max=200"`
    Role     string `json:"role"      binding:"omitempty,oneof=super_admin admin operator guest"`
}

type ChangePasswordRequest struct {
    OldPassword     string `json:"old_password" binding:"required"`
    NewPassword     string `json:"new_password" binding:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

type AdminChangePasswordRequest struct {
    NewPassword     string `json:"new_password" binding:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

type UserListResponse struct {
    Data       []PublicUser `json:"data"`
    Total      int64        `json:"total"`
    Page       int          `json:"page"`
    Limit      int          `json:"limit"`
    TotalPages int          `json:"total_pages"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type UserService struct {
    userRepo *repository.UserRepository
    roleRepo *repository.RoleRepository
}

func NewUserService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
    return &UserService{userRepo: userRepo, roleRepo: roleRepo}
}

// GetAll mengambil semua user dengan pagination
func (s *UserService) GetAll(page, limit int) (*UserListResponse, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }

    users, total, err := s.userRepo.FindAll(page, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to get users: %w", err)
    }

    publicUsers := make([]PublicUser, len(users))
    for i, u := range users {
        publicUsers[i] = toPublicUser(&u)
    }

    totalPages := int(total) / limit
    if int(total)%limit != 0 {
        totalPages++
    }

    return &UserListResponse{
        Data:       publicUsers,
        Total:      total,
        Page:       page,
        Limit:      limit,
        TotalPages: totalPages,
    }, nil
}

// GetByID mengambil detail user berdasarkan ID
func (s *UserService) GetByID(id uint) (*PublicUser, error) {
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    if user == nil {
        return nil, errors.New("user tidak ditemukan")
    }
    result := toPublicUser(user)
    return &result, nil
}

// CreateUser membuat user baru (oleh admin/super_admin)
func (s *UserService) CreateUser(req CreateUserRequest) (*PublicUser, error) {
    // Cek duplikat email
    existing, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    if existing != nil {
        return nil, errors.New("email sudah terdaftar")
    }

    // Cari role
    role, err := s.roleRepo.FindByName(req.Role)
    if err != nil {
        return nil, fmt.Errorf("failed to find role: %w", err)
    }
    if role == nil {
        return nil, errors.New("role tidak ditemukan")
    }

    // Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    hashStr := string(hash)

    user := &model.User{
        FullName:     req.FullName,
        Email:        req.Email,
        Company:      req.Company,
        Provider:     model.ProviderLocal,
        PasswordHash: &hashStr,
        RoleID:       role.ID,
        Role:         *role,
        LegacyRole:   role.Name,
        IsVerified:   true,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    result := toPublicUser(user)
    return &result, nil
}

// UpdateUser memperbarui data user
func (s *UserService) UpdateUser(id uint, req UpdateUserRequest, requestorRole string) (*PublicUser, error) {
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    if user == nil {
        return nil, errors.New("user tidak ditemukan")
    }

    targetRole := user.Role.Name
    
    // Guard 1: Non-super_admin dilarang mengubah super_admin
    if targetRole == "super_admin" && requestorRole != "super_admin" {
        return nil, errors.New("tidak punya izin mengubah akun super_admin")
    }
    // Guard 2: Admin dilarang mengubah admin lain
    if targetRole == "admin" && requestorRole == "admin" {
        return nil, errors.New("admin tidak memiliki izin mengubah akun admin lain")
    }
    // Guard 3: Hanya super_admin yang bisa mengangkat orang lain menjadi super_admin
    if req.Role == "super_admin" && requestorRole != "super_admin" {
        return nil, errors.New("hanya super_admin yang bisa mengangkat super_admin")
    }

    // Cek email unik jika email diubah
    if req.Email != "" && req.Email != user.Email {
        exists, err := s.userRepo.EmailExistsExcept(req.Email, id)
        if err != nil {
            return nil, fmt.Errorf("database error: %w", err)
        }
        if exists {
            return nil, errors.New("email sudah digunakan oleh user lain")
        }
        user.Email = req.Email
    }

    if req.FullName != "" {
        user.FullName = req.FullName
    }
    if req.Company != "" {
        user.Company = req.Company
    }
    if req.Role != "" {
        role, err := s.roleRepo.FindByName(req.Role)
        if err != nil {
            return nil, fmt.Errorf("database error: %w", err)
        }
        if role == nil {
            return nil, errors.New("role tidak ditemukan")
        }
        user.RoleID = role.ID
        user.Role = *role
        user.LegacyRole = role.Name
    }

    if err := s.userRepo.Update(user); err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err)
    }

    result := toPublicUser(user)
    return &result, nil
}

// DeleteUser soft-delete user berdasarkan ID
func (s *UserService) DeleteUser(id uint, requestorID uint, requestorRole string) error {
    if id == requestorID {
        return errors.New("tidak bisa menghapus akun sendiri")
    }

    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return fmt.Errorf("database error: %w", err)
    }
    if user == nil {
        return errors.New("user tidak ditemukan")
    }
    
    targetRole := user.Role.Name
    
    if targetRole == "super_admin" && requestorRole != "super_admin" {
        return errors.New("tidak punya izin menghapus akun super_admin")
    }
    if targetRole == "admin" && requestorRole == "admin" {
        return errors.New("admin tidak memiliki izin menghapus akun admin lain")
    }

    return s.userRepo.Delete(id)
}

// ChangePassword ganti password sendiri (user yang sedang login)
func (s *UserService) ChangePassword(userID uint, req ChangePasswordRequest) error {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return fmt.Errorf("database error: %w", err)
    }
    if user == nil {
        return errors.New("user tidak ditemukan")
    }

    if user.Provider != model.ProviderLocal || user.PasswordHash == nil {
        return errors.New("akun Google tidak bisa ganti password melalui endpoint ini")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.OldPassword)); err != nil {
        return errors.New("password lama tidak sesuai")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    return s.userRepo.UpdatePassword(userID, string(hash))
}

// AdminResetPassword reset password user lain (oleh admin/super_admin)
func (s *UserService) AdminResetPassword(targetID uint, req AdminChangePasswordRequest, requestorRole string) error {
    user, err := s.userRepo.FindByID(targetID)
    if err != nil {
        return fmt.Errorf("database error: %w", err)
    }
    if user == nil {
        return errors.New("user tidak ditemukan")
    }

    targetRole := user.Role.Name

    if targetRole == "super_admin" && requestorRole != "super_admin" {
        return errors.New("tidak punya izin mereset password akun super_admin")
    }
    if targetRole == "admin" && requestorRole == "admin" {
        return errors.New("admin tidak memiliki izin mereset password akun admin lain")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    return s.userRepo.UpdatePassword(targetID, string(hash))
}