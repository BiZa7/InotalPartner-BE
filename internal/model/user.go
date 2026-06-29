package model

import (
	"time"

	"gorm.io/gorm"
)


type AuthProvider string

const (
	ProviderLocal  AuthProvider = "local"  //email/password
	ProviderGoogle AuthProvider = "google" //google login
)


type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RolePartner UserRole = "partner"
)


type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"           json:"id"`
	CreatedAt time.Time      `                                          json:"created_at"`
	UpdatedAt time.Time      `                                          json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                              json:"-"`

	// Identitas
	FullName string `gorm:"type:varchar(150);not null"         json:"full_name"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`

	// Auth
	Provider AuthProvider `gorm:"type:varchar(20);not null;default:'local'" json:"provider"`
	GoogleID string       `gorm:"type:varchar(100);index"            json:"google_id,omitempty"`

	// Password — nullable karena Google login tidak butuh password
	PasswordHash *string `gorm:"type:varchar(255)"                  json:"-"`

	// Profil
	AvatarURL string   `gorm:"type:text"                          json:"avatar_url"`
	Company   string   `gorm:"type:varchar(200)"                  json:"company"`
	Role      UserRole `gorm:"type:varchar(20);not null;default:'partner'" json:"role"`

	// Status
	IsVerified bool       `gorm:"default:false"                      json:"is_verified"`
	LastLoginAt *time.Time `                                          json:"last_login_at"`
}