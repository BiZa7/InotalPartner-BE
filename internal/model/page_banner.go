package model

import (
	"time"
)

// PageBanner merepresentasikan data hero banner dinamis untuk halaman tertentu (about, event, news, contact).
type PageBanner struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"                 json:"id"`
	PageType    string    `gorm:"type:varchar(50);not null;uniqueIndex"    json:"page_type"`
	ImageURL    string    `gorm:"type:text;default:''"                     json:"image_url"`
	Subtitle    string    `gorm:"type:varchar(255);default:''"             json:"subtitle"`
	Title       string    `gorm:"type:varchar(255);not null"               json:"title"`
	Description string    `gorm:"type:text;default:''"                     json:"description"`
	IsActive    bool      `gorm:"default:true;not null"                    json:"is_active"`
	CreatedAt   time.Time `                                                json:"created_at"`
	UpdatedAt   time.Time `                                                json:"updated_at"`
}

// TableName menentukan nama tabel database untuk PageBanner.
func (PageBanner) TableName() string {
	return "page_banners"
}

