package model

import (
	"time"
)

// PartnershipService merepresentasikan kartu layanan pada section Partnership Beranda.
// Mendukung penambahan kartu dinamis dengan nama ikon Lucide, urutan tampilan, dan status aktif/nonaktif.
type PartnershipService struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"              json:"id"`
	Title        string    `gorm:"type:varchar(150);not null;default:''" json:"title"`
	Description  string    `gorm:"type:text;default:''"                  json:"description"`
	Icon         string    `gorm:"type:varchar(100);default:''"          json:"icon"`
	DisplayOrder int       `gorm:"default:0;not null"                    json:"display_order"`
	IsActive     bool      `gorm:"default:true;not null"                 json:"is_active"`
	CreatedAt    time.Time `                                             json:"created_at"`
	UpdatedAt    time.Time `                                             json:"updated_at"`
}

// TableName menentukan nama tabel database untuk PartnershipService.
func (PartnershipService) TableName() string {
	return "partnership_services"
}

