package model

import (
	"time"
)

// HomepageBanner merepresentasikan satu slide pada Banner section carousel Beranda.
type HomepageBanner struct {
	ID              uint                   `gorm:"primaryKey;autoIncrement"                             json:"id"`
	Title           string                 `gorm:"type:varchar(255);not null"                           json:"title"`
	Subtitle        string                 `gorm:"type:text;not null"                                  json:"subtitle"`
	BackgroundImage string                 `gorm:"type:text;not null"                                  json:"background_image"`
	Order           int                    `gorm:"default:0"                                            json:"order"`
	IsActive        bool                   `gorm:"default:true"                                         json:"is_active"`
	Buttons         []HomepageBannerButton `gorm:"foreignKey:BannerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"buttons"`
	CreatedAt       time.Time              `                                                            json:"created_at"`
	UpdatedAt       time.Time              `                                                            json:"updated_at"`
}

// HomepageBannerButton merepresentasikan tombol action yang terikat pada satu slide HomepageBanner.
type HomepageBannerButton struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BannerID    uint      `gorm:"not null;index"            json:"banner_id"`
	Label       string    `gorm:"type:varchar(100);not null" json:"label"`
	URL         string    `gorm:"type:varchar(255);not null" json:"url"`
	Description string    `gorm:"type:text"                  json:"description,omitempty"`
	Order       int       `gorm:"default:0"                  json:"order"`
	CreatedAt   time.Time `                                  json:"created_at"`
	UpdatedAt   time.Time `                                  json:"updated_at"`
}
