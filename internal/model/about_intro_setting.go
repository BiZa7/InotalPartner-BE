package model

import (
	"time"
)

// AboutIntroSetting merepresentasikan konfigurasi utama section Hero/Perkenalan pada Halaman Tentang (About Us).
// Dirancang sebagai single active configuration (pengaturan tunggal aktif).
type AboutIntroSetting struct {
	ID                 uint                `gorm:"primaryKey;autoIncrement"              json:"id"`
	Eyebrow            string              `gorm:"type:varchar(150);not null;default:''" json:"eyebrow"`
	Title              string              `gorm:"type:varchar(255);not null;default:''" json:"title"`
	ImageURL           string              `gorm:"type:varchar(500);not null;default:''" json:"image_url"`
	ImageAlt           string              `gorm:"type:varchar(255);default:''"          json:"image_alt"`
	ProgressTitle      string              `gorm:"type:varchar(150);not null;default:''" json:"progress_title"`
	MissionTitle       string              `gorm:"type:varchar(150);not null;default:''" json:"mission_title"`
	MissionDescription string              `gorm:"type:text;default:''"                  json:"mission_description"`
	VisionTitle        string              `gorm:"type:varchar(150);not null;default:''" json:"vision_title"`
	VisionDescription  string              `gorm:"type:text;default:''"                  json:"vision_description"`
	CTAButtonText      string              `gorm:"type:varchar(100);not null;default:''" json:"cta_button_text"`
	CTAButtonURL       string              `gorm:"type:varchar(500);not null;default:''" json:"cta_button_url"`
	IsActive           bool                `gorm:"default:true;not null"                 json:"is_active"`
	CreatedAt          time.Time           `                                             json:"created_at"`
	UpdatedAt          time.Time           `                                             json:"updated_at"`
	Metrics            []AboutIntroMetric  `gorm:"foreignKey:SettingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"metrics,omitempty"`
	Points             []AboutIntroPoint   `gorm:"foreignKey:SettingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"points,omitempty"`
}

// TableName menentukan nama tabel database untuk AboutIntroSetting.
func (AboutIntroSetting) TableName() string {
	return "about_intro_settings"
}

