package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"inotal-be/config"
	"inotal-be/internal/model"
)

var DB *gorm.DB

func Connect() {
	cfg := config.App

	gormCfg := &gorm.Config{}
	if cfg.AppEnv == "development" {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
	if err != nil {
		log.Fatalf("[database] Failed to connect to PostgreSQL: %v", err)
	}

	// AutoMigrate: buat/update tabel otomatis sesuai model
	if err := db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
	); err != nil {
		log.Fatalf("[database] AutoMigrate failed: %v", err)
	}

	log.Println("[database] Connected and migrated successfully")
	DB = db
}
