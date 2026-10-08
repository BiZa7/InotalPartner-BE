package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create menyimpan user baru ke database
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// CountAll menghitung total user
func (r *UserRepository) CountAll(count *int64) error {
	return r.db.Model(&model.User{}).Count(count).Error
}

// CountByRoleID menghitung total user berdasarkan role ID tertentu
func (r *UserRepository) CountByRoleID(roleID uint, count *int64) error {
	return r.db.Model(&model.User{}).
		Where("role_id = ?", roleID).
		Count(count).Error
}

// FindByEmail mencari user berdasarkan email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

// FindByID mencari user berdasarkan ID
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User

	err := r.db.
		Preload("Role").
		First(&user, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

// FindByGoogleID mencari user berdasarkan Google ID
func (r *UserRepository) FindByGoogleID(googleID string) (*model.User, error) {
	var user model.User

	err := r.db.
		Preload("Role").
		Where("google_id = ?", googleID).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

// UpdateLastLogin memperbarui waktu login terakhir
func (r *UserRepository) UpdateLastLogin(user *model.User) error {
	return r.db.Save(user).Error
}

// FindAll mengambil user dengan pagination dan pencarian berdasarkan nama
func (r *UserRepository) FindAll(page int, limit int, search string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&model.User{})

	search = strings.TrimSpace(search)

	if search != "" {
		search = strings.ToLower(search)

		query = query.Where(
			"LOWER(full_name) LIKE ?",
			"%"+search+"%",
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Role").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update memperbarui data user
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete soft-delete user berdasarkan ID
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// UpdatePassword memperbarui password hash user
func (r *UserRepository) UpdatePassword(id uint, passwordHash string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).Error
}

// EmailExistsExcept cek apakah email sudah dipakai user lain
func (r *UserRepository) EmailExistsExcept(email string, excludeID uint) (bool, error) {
	var count int64

	err := r.db.
		Model(&model.User{}).
		Where("email = ? AND id != ?", email, excludeID).
		Count(&count).Error

	return count > 0, err
}