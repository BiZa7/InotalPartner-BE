package model

import (
	"time"
)

// PartnershipSetting merepresentasikan pengaturan utama section Partnership pada halaman Beranda.
// Dirancang sebagai single active configuration (pengaturan tunggal).
type PartnershipSetting struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"              json:"id"`
	Eyebrow     string    `gorm:"type:varchar(100);not null;default:''" json:"eyebrow"`
	Title       string    `gorm:"type:varchar(150);not null;default:''" json:"title"`
	Description string    `gorm:"type:text;default:''"                  json:"description"`
	BadgeValue  string    `gorm:"type:varchar(50);default:''"           json:"badge_value"`
	BadgeLabel  string    `gorm:"type:varchar(100);default:''"          json:"badge_label"`
	IsActive    bool      `gorm:"default:true;not null"                 json:"is_active"`
	CreatedAt   time.Time `                                             json:"created_at"`
	UpdatedAt   time.Time `                                             json:"updated_at"`
}

// TableName menentukan nama tabel database untuk PartnershipSetting.
func (PartnershipSetting) TableName() string {
	return "partnership_settings"
}

