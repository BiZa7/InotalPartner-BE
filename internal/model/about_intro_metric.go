package model

import (
	"time"
)

// AboutIntroMetric merepresentasikan indikator progres kemitraan pada section Hero/Perkenalan Halaman Tentang.
// Mendukung jumlah metric fleksibel dengan persentase 0-100 dan urutan tampilan.
type AboutIntroMetric struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"              json:"id"`
	SettingID    uint      `gorm:"not null;index"                        json:"setting_id"`
	Label        string    `gorm:"type:varchar(150);not null;default:''" json:"label"`
	Value        int       `gorm:"not null;default:0"                    json:"value"` // Persentase 0 - 100
	DisplayOrder int       `gorm:"default:0;not null"                    json:"display_order"`
	IsActive     bool      `gorm:"default:true;not null"                 json:"is_active"`
	CreatedAt    time.Time `                                             json:"created_at"`
	UpdatedAt    time.Time `                                             json:"updated_at"`
}

// TableName menentukan nama tabel database untuk AboutIntroMetric.
func (AboutIntroMetric) TableName() string {
	return "about_intro_metrics"
}

