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

	// Migrate existing data from articles.category_id to article_categories pivot table
	migrationQuery := `
		INSERT INTO article_categories (article_id, category_id)
		SELECT id, category_id FROM articles
		WHERE category_id IS NOT NULL AND category_id > 0
		ON CONFLICT (article_id, category_id) DO NOTHING;
	`
	if err := db.Exec(migrationQuery).Error; err != nil {
		log.Printf("[database] Warning: Failed to migrate existing category_id data: %v", err)
	}

	// Migrate event_time column type to TIME WITHOUT TIME ZONE if currently timestamp/timestamptz
	alterEventTimeQuery := `
		DO $$ 
		BEGIN 
			IF EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'articles' AND column_name = 'event_time' AND data_type LIKE '%timestamp%'
			) THEN 
				ALTER TABLE articles ALTER COLUMN event_time TYPE time without time zone USING event_time::time without time zone;
			END IF;
		END $$;
	`
	if err := db.Exec(alterEventTimeQuery).Error; err != nil {
		log.Printf("[database] Warning: Failed to alter event_time column type: %v", err)
	}

	log.Println("[database] Connected and migrated successfully")
	DB = db
}
