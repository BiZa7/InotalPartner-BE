package repository

import (
	"errors"

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

// CountAll menghitung total user (termasuk soft-deleted tidak dihitung)
func (r *UserRepository) CountAll(count *int64) error {
	return r.db.Model(&model.User{}).Count(count).Error
}

// CountByRoleID menghitung total user berdasarkan role ID tertentu
func (r *UserRepository) CountByRoleID(roleID uint, count *int64) error {
	return r.db.Model(&model.User{}).Where("role_id = ?", roleID).Count(count).Error
}

// FindByEmail mencari user berdasarkan email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // tidak ditemukan bukan error fatal
	}
	return &user, err
}

// FindByID mencari user berdasarkan ID
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// FindByGoogleID mencari user berdasarkan Google ID
func (r *UserRepository) FindByGoogleID(googleID string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").Where("google_id = ?", googleID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// UpdateLastLogin memperbarui waktu login terakhir
func (r *UserRepository) UpdateLastLogin(user *model.User) error {
	return r.db.Save(user).Error
}

// FindAll mengambil semua user dengan pagination
func (r *UserRepository) FindAll(page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Role").Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
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
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("password_hash", passwordHash).Error
}

// EmailExistsExcept cek apakah email sudah dipakai user lain (selain dirinya sendiri)
func (r *UserRepository) EmailExistsExcept(email string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).
		Where("email = ? AND id != ?", email, excludeID).
		Count(&count).Error
	return count > 0, err
}