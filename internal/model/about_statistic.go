package model

import (
	"time"
)

// AboutStatistic merepresentasikan kartu statistik pencapaian pada Halaman Tentang (About Us).
type AboutStatistic struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"                 json:"id"`
	Title           string    `gorm:"type:varchar(150);not null"               json:"title"`
	Icon            string    `gorm:"type:varchar(100);not null"               json:"icon"`
	StatisticValue  string    `gorm:"type:varchar(30);not null"                json:"statistic_value"`
	StatisticSuffix string    `gorm:"type:varchar(20);default:''"              json:"statistic_suffix"`
	Description     string    `gorm:"type:text;default:''"                     json:"description"`
	DisplayOrder    int       `gorm:"default:0;not null"                       json:"display_order"`
	IsActive        bool      `gorm:"default:true;not null"                    json:"is_active"`
	CreatedAt       time.Time `                                                json:"created_at"`
	UpdatedAt       time.Time `                                                json:"updated_at"`
}

// TableName menentukan nama tabel database untuk AboutStatistic.
func (AboutStatistic) TableName() string {
	return "about_statistics"
}

