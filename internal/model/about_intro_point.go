package model

import (
	"time"
)

// AboutIntroPoint merepresentasikan daftar poin Misi dan Visi pada section Hero/Perkenalan Halaman Tentang.
// SectionType menerima 'mission' atau 'vision'.
type AboutIntroPoint struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"              json:"id"`
	SettingID    uint      `gorm:"not null;index"                        json:"setting_id"`
	SectionType  string    `gorm:"type:varchar(20);not null;default:''"  json:"section_type"` // 'mission' atau 'vision'
	Content      string    `gorm:"type:text;not null;default:''"         json:"content"`
	DisplayOrder int       `gorm:"default:0;not null"                    json:"display_order"`
	IsActive     bool      `gorm:"default:true;not null"                 json:"is_active"`
	CreatedAt    time.Time `                                             json:"created_at"`
	UpdatedAt    time.Time `                                             json:"updated_at"`
}

// TableName menentukan nama tabel database untuk AboutIntroPoint.
func (AboutIntroPoint) TableName() string {
	return "about_intro_points"
}

