package model

import (
	"time"
)

// AboutPageSetting merepresentasikan konfigurasi judul utama pada Halaman Tentang (About Us).
// Pengaturan ini dirancang sebagai konfigurasi tunggal aktif (single active configuration).
type AboutPageSetting struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"              json:"id"`
	Title     string    `gorm:"type:varchar(255);not null;default:''" json:"title"`
	TeamTitle string    `gorm:"type:varchar(255);not null;default:''" json:"team_title"`
	IsActive  bool      `gorm:"default:true;not null"                 json:"is_active"`
	CreatedAt time.Time `                                             json:"created_at"`
	UpdatedAt time.Time `                                             json:"updated_at"`
}

// TableName menentukan nama tabel database untuk AboutPageSetting.
func (AboutPageSetting) TableName() string {
	return "about_page_settings"
}

// TeamMember merepresentasikan anggota tim pada tabel khusus team_members.
// Nama dan foto anggota dapat dikosongkan saat berstatus draft atau dalam persiapan.
// Jabatan (position) bersifat opsional.
type TeamMember struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"        json:"id"`
	Name         string    `gorm:"type:varchar(150);default:''"    json:"name"`
	PhotoURL     string    `gorm:"type:text;default:''"            json:"photo_url"`
	Position     string    `gorm:"type:varchar(100);default:''"    json:"position"`
	DisplayOrder int       `gorm:"default:0;not null"              json:"display_order"`
	IsActive     bool      `gorm:"default:true;not null"           json:"is_active"`
	CreatedAt    time.Time `                                       json:"created_at"`
	UpdatedAt    time.Time `                                       json:"updated_at"`
}

// TableName menentukan nama tabel database untuk TeamMember.
func (TeamMember) TableName() string {
	return "team_members"
}

// AboutVideo merepresentasikan video profil perusahaan pada halaman Tentang Kami (About Us).
type AboutVideo struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"        json:"id"`
	Title        string    `gorm:"type:varchar(255);not null"      json:"title"`
	Description  string    `gorm:"type:text;default:''"            json:"description"`
	VideoURL     string    `gorm:"type:text;not null"              json:"video_url"`
	ThumbnailURL string    `gorm:"type:text;default:''"            json:"thumbnail_url"`
	IsActive     bool      `gorm:"default:true;not null"           json:"is_active"`
	CreatedAt    time.Time `                                       json:"created_at"`
	UpdatedAt    time.Time `                                       json:"updated_at"`
}

// TableName menentukan nama tabel database untuk AboutVideo.
func (AboutVideo) TableName() string {
	return "about_videos"
}

// AboutFAQ merepresentasikan tanya jawab FAQ pada Halaman Tentang (About Us).
type AboutFAQ struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"        json:"id"`
	Question     string    `gorm:"type:varchar(255);not null"      json:"question"`
	Answer       string    `gorm:"type:text;not null"              json:"answer"`
	DisplayOrder int       `gorm:"default:0;not null"              json:"display_order"`
	IsActive     bool      `gorm:"default:true;not null"           json:"is_active"`
	CreatedAt    time.Time `                                       json:"created_at"`
	UpdatedAt    time.Time `                                       json:"updated_at"`
}

// TableName menentukan nama tabel database untuk AboutFAQ.
func (AboutFAQ) TableName() string {
	return "about_faqs"
}

// AboutVideoSetting merepresentasikan pengaturan header/section video profil pada halaman Tentang Kami.
// Dirancang sebagai single-record configuration.
type AboutVideoSetting struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"              json:"id"`
	SectionLabel       string    `gorm:"type:varchar(100);default:''"          json:"section_label"`
	SectionTitle       string    `gorm:"type:varchar(255);not null"            json:"section_title"`
	SectionDescription string    `gorm:"type:text;default:''"                  json:"section_description"`
	IsActive           bool      `gorm:"default:true;not null"                 json:"is_active"`
	CreatedAt          time.Time `                                             json:"created_at"`
	UpdatedAt          time.Time `                                             json:"updated_at"`
}

// TableName menentukan nama tabel database untuk AboutVideoSetting.
func (AboutVideoSetting) TableName() string {
	return "about_video_settings"
}

// AboutVideoHighlight merepresentasikan poin-poin keunggulan/highlight video profil perusahaan.
type AboutVideoHighlight struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"        json:"id"`
	Icon         string    `gorm:"type:varchar(50);default:''"     json:"icon"`
	Title        string    `gorm:"type:varchar(150);not null"      json:"title"`
	Description  string    `gorm:"type:text;not null"              json:"description"`
	DisplayOrder int       `gorm:"default:0;not null"              json:"display_order"`
	IsActive     bool      `gorm:"default:true;not null"           json:"is_active"`
	CreatedAt    time.Time `                                       json:"created_at"`
	UpdatedAt    time.Time `                                       json:"updated_at"`
}

// TableName menentukan nama tabel database untuk AboutVideoHighlight.
func (AboutVideoHighlight) TableName() string {
	return "about_video_highlights"
}
