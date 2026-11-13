package model

import (
	"time"
)

// PartnershipImage merepresentasikan gambar-gambar kolase pada section Partnership Beranda.
// Mendukung banyak gambar dengan pengaturan urutan tampilan dan status aktif/nonaktif.
type PartnershipImage struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"              json:"id"`
	ImageURL     string    `gorm:"type:varchar(500);not null;default:''" json:"image_url"`
	AltText      string    `gorm:"type:varchar(255);default:''"          json:"alt_text"`
	DisplayOrder int       `gorm:"default:0;not null"                    json:"display_order"`
	IsActive     bool      `gorm:"default:true;not null"                 json:"is_active"`
	CreatedAt    time.Time `                                             json:"created_at"`
	UpdatedAt    time.Time `                                             json:"updated_at"`
}

// TableName menentukan nama tabel database untuk PartnershipImage.
func (PartnershipImage) TableName() string {
	return "partnership_images"
}

