package model

import (
	"time"
)

// HomepageReasonSetting merepresentasikan pengaturan utama section "Mengapa Bergabung" di beranda.
// Dirancang sebagai single active configuration (pengaturan tunggal).
type HomepageReasonSetting struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"                 json:"id"`
	SectionLabel       string    `gorm:"type:varchar(100);default:''"             json:"section_label"`
	SectionTitle       string    `gorm:"type:varchar(255);not null;default:''"    json:"section_title"`
	SectionDescription string    `gorm:"type:text;default:''"                     json:"section_description"`
	ImageOneURL        string    `gorm:"type:text;default:''"                     json:"image_one_url"`
	ImageTwoURL        string    `gorm:"type:text;default:''"                     json:"image_two_url"`
	ImageThreeURL      string    `gorm:"type:text;default:''"                     json:"image_three_url"`
	IsActive           bool      `gorm:"default:true;not null"                    json:"is_active"`
	CreatedAt          time.Time `                                                json:"created_at"`
	UpdatedAt          time.Time `                                                json:"updated_at"`
}

// TableName menentukan nama tabel database untuk HomepageReasonSetting.
func (HomepageReasonSetting) TableName() string {
	return "homepage_reason_settings"
}

// HomepageReason merepresentasikan kartu alasan pada section "Mengapa Bergabung".
type HomepageReason struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"                 json:"id"`
	Number       string    `gorm:"type:varchar(20);default:''"              json:"number"`
	Title        string    `gorm:"type:varchar(255);not null;default:''"    json:"title"`
	Description  string    `gorm:"type:text;default:''"                     json:"description"`
	DisplayOrder int       `gorm:"default:0;not null"                       json:"display_order"`
	IsActive     bool      `gorm:"default:true;not null"                    json:"is_active"`
	CreatedAt    time.Time `                                                json:"created_at"`
	UpdatedAt    time.Time `                                                json:"updated_at"`
}

// TableName menentukan nama tabel database untuk HomepageReason.
func (HomepageReason) TableName() string {
	return "homepage_reasons"
}

