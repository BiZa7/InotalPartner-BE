package model

import (
	"time"
)

// HomeCTA merepresentasikan data Call to Action (CTA) pada section Beranda sebelum footer.
type HomeCTA struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"                 json:"id"`
	SectionLabel string    `gorm:"type:varchar(100);not null"               json:"section_label"`
	Title        string    `gorm:"type:text;not null"                       json:"title"`
	ButtonText   string    `gorm:"type:varchar(100);not null"               json:"button_text"`
	ButtonURL    string    `gorm:"type:varchar(255);not null"               json:"button_url"`
	IsActive     bool      `gorm:"default:true;not null"                    json:"is_active"`
	CreatedAt    time.Time `                                                json:"created_at"`
	UpdatedAt    time.Time `                                                json:"updated_at"`
}

// TableName menentukan nama tabel database untuk HomeCTA.
func (HomeCTA) TableName() string {
	return "home_ctas"
}

