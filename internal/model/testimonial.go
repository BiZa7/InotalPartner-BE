package model

import (
	"time"
)

// Testimonial merepresentasikan data testimoni untuk section testimoni pada halaman Beranda.
type Testimonial struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"     json:"id"`
	Name         string    `gorm:"type:varchar(150);not null"   json:"name"`
	Position     string    `gorm:"type:varchar(100);default:''" json:"position"`
	Company      string    `gorm:"type:varchar(150);default:''" json:"company"`
	Content      string    `gorm:"type:text;not null"           json:"content"`
	PhotoURL     string    `gorm:"type:text;default:''"         json:"photo_url"`
	DisplayOrder int       `gorm:"default:0;not null"           json:"display_order"`
	IsActive     bool      `gorm:"default:true;not null"        json:"is_active"`
	CreatedAt    time.Time `                                    json:"created_at"`
	UpdatedAt    time.Time `                                    json:"updated_at"`
}

// TableName menentukan nama tabel database untuk Testimonial.
func (Testimonial) TableName() string {
	return "testimonials"
}
