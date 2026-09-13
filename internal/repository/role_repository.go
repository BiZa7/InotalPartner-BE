package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// FindByName mencari role berdasarkan nama (e.g. "guest", "super_admin")
func (r *RoleRepository) FindByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

// FindByID mencari role berdasarkan ID
func (r *RoleRepository) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

// FindAll mengambil seluruh data role diurutkan berdasarkan ID
func (r *RoleRepository) FindAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Order("id ASC").Find(&roles).Error
	return roles, err
}

// Create membuat data role baru
func (r *RoleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// Update memperbarui data role
func (r *RoleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

// Delete menghapus role berdasarkan ID
func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

// NameExistsExcept mengecek apakah nama role sudah digunakan oleh role lain
func (r *RoleRepository) NameExistsExcept(name string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Role{}).Where("name = ? AND id != ?", name, excludeID).Count(&count).Error
	return count > 0, err
}

// HasUsers mengecek apakah ada user aktif yang menggunakan role ID ini
func (r *RoleRepository) HasUsers(roleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("role_id = ?", roleID).Count(&count).Error
	return count > 0, err
}

// SeedDefaultRoles memastikan 4 role bawaan (super_admin, admin, operator, guest) ada di DB
func (r *RoleRepository) SeedDefaultRoles() error {
	defaultRoles := []model.Role{
		{Name: "super_admin", DisplayName: "Super Admin", IsSystem: true},
		{Name: "admin", DisplayName: "Admin", IsSystem: true},
		{Name: "operator", DisplayName: "Operator", IsSystem: true},
		{Name: "guest", DisplayName: "Guest", IsSystem: true},
	}

	for _, role := range defaultRoles {
		var existing model.Role
		err := r.db.Where("name = ?", role.Name).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.Create(&role).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

