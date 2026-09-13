package model

import (
	"time"
)

type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	DisplayName string    `gorm:"type:varchar(100);not null" json:"display_name"`
	IsSystem    bool      `gorm:"default:true" json:"is_system"`
}
