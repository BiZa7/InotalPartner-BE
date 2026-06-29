package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"inotal-be/config"
	"inotal-be/internal/database"
	"inotal-be/internal/repository"
	"inotal-be/internal/service"
	"inotal-be/routes"
)

func main() {
	// 1. Load config dari .env
	config.Load()

	// 2. Koneksi ke PostgreSQL + AutoMigrate
	database.Connect()

	// 3. Dependency injection (manual, tanpa DI framework)
	userRepo := repository.NewUserRepository(database.DB)
	authSvc := service.NewAuthService(userRepo)

	// 4. Setup Gin
	if config.App.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 5. Daftarkan semua route
	routes.Setup(r, authSvc)

	// 6.  server
	addr := ":" + config.App.AppPort
	log.Printf("[server] Starting on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[server] Failed to start: %v", err)
	}
}