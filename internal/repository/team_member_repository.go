package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type TeamMemberRepository struct {
	db *gorm.DB
}

func NewTeamMemberRepository(db *gorm.DB) *TeamMemberRepository {
	return &TeamMemberRepository{db: db}
}

// FindAll mengambil seluruh anggota tim, termasuk yang nonaktif, untuk kebutuhan pengelolaan admin.
// Diurutkan berdasarkan display_order ASC, id ASC.
func (r *TeamMemberRepository) FindAll() ([]model.TeamMember, error) {
	var members []model.TeamMember
	err := r.db.
		Order("display_order ASC, id ASC").
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	if members == nil {
		members = []model.TeamMember{}
	}
	return members, nil
}

// FindByID mengambil detail anggota tim berdasarkan ID.
func (r *TeamMemberRepository) FindByID(id uint) (*model.TeamMember, error) {
	var member model.TeamMember
	err := r.db.First(&member, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FindByDisplayOrder mencari anggota tim berstatus AKTIF dengan display_order tertentu.
// Jika excludeID > 0, record dengan ID tersebut akan diabaikan (untuk validasi edit).
// Hanya anggota berstatus Aktif (is_active = true) yang dianggap memblokir nomor urut tersebut.
func (r *TeamMemberRepository) FindByDisplayOrder(displayOrder int, excludeID uint) (*model.TeamMember, error) {
	var member model.TeamMember
	query := r.db.Where("display_order = ? AND is_active = ?", displayOrder, true)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	err := query.First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// Create membuat record anggota tim baru.
func (r *TeamMemberRepository) Create(member *model.TeamMember) error {
	return r.db.Create(member).Error
}

// Update memperbarui data anggota tim berdasarkan ID.
func (r *TeamMemberRepository) Update(member *model.TeamMember) error {
	return r.db.Save(member).Error
}

// Delete menghapus anggota tim berdasarkan ID.
func (r *TeamMemberRepository) Delete(id uint) error {
	return r.db.Delete(&model.TeamMember{}, id).Error
}

// FindActivePublic mengambil daftar anggota aktif yang memiliki nama dan foto tidak kosong (bukan hanya whitespace),
// diurutkan berdasarkan display_order ASC, id ASC, untuk kebutuhan API publik.
func (r *TeamMemberRepository) FindActivePublic() ([]model.TeamMember, error) {
	var members []model.TeamMember
	err := r.db.
		Where("is_active = ?", true).
		Where("name IS NOT NULL AND TRIM(name) <> ''").
		Where("photo_url IS NOT NULL AND TRIM(photo_url) <> ''").
		Order("display_order ASC, id ASC").
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	if members == nil {
		members = []model.TeamMember{}
	}
	return members, nil
}
