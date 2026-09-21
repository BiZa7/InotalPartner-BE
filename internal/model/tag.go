package model

import (
	"time"
)

type Tag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"           json:"id"`
	CreatedAt time.Time `                                          json:"created_at"`
	UpdatedAt time.Time `                                          json:"updated_at"`

	Name string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Slug string `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
}
