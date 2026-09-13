package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"inotal-be/config"
	"inotal-be/internal/handler"
	"inotal-be/internal/middleware"
	"inotal-be/internal/service"
)

func Setup(r *gin.Engine, authSvc *service.AuthService, userSvc *service.UserService, roleSvc *service.RoleService) {
	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.App.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "inotal-be"})
	})

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	roleHandler := handler.NewRoleHandler(roleSvc)

	// API v1
	api := r.Group("/api")
	{
		// ── Auth ────────────────────────────────────────────────────────────
		auth := api.Group("/auth")
		{
			// Public routes
			auth.POST("/setup", authHandler.Setup)    // ← first-time setup (DB kosong)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/google", authHandler.GoogleLogin)
			auth.GET("/google/callback", authHandler.GoogleCallback)


			// Protected routes (hanya perlu login)
			protected := auth.Group("")
			protected.Use(middleware.AuthMiddleware(authSvc))
			{
				protected.GET("/me", authHandler.Me)
			}
		}

		// ── Users ───────────────────────────────────────────────────────────
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(authSvc))
		{
			// Ganti password sendiri — semua role yang sudah login
			users.PATCH("/me/password", userHandler.ChangePassword)

			// CRUD user — hanya admin & super_admin
			adminRoutes := users.Group("")
			adminRoutes.Use(middleware.AdminOrAbove())
			{
				adminRoutes.GET("", userHandler.GetUsers)          // GET  /api/users
				adminRoutes.POST("", userHandler.CreateUser)       // POST /api/users
				adminRoutes.GET("/:id", userHandler.GetUser)       // GET  /api/users/:id
				adminRoutes.PUT("/:id", userHandler.UpdateUser)    // PUT  /api/users/:id
				adminRoutes.DELETE("/:id", userHandler.DeleteUser) // DELETE /api/users/:id

				// Reset password user lain
				adminRoutes.PATCH("/:id/password", userHandler.AdminResetPassword) // PATCH /api/users/:id/password
			}
		}

		// ── Roles ───────────────────────────────────────────────────────────
		roles := api.Group("/roles")
		roles.Use(middleware.AuthMiddleware(authSvc))
		roles.Use(middleware.AdminOrAbove())
		{
			roles.GET("", roleHandler.GetRoles)          // GET    /api/roles
			roles.POST("", roleHandler.CreateRole)       // POST   /api/roles
			roles.GET("/:id", roleHandler.GetRole)       // GET    /api/roles/:id
			roles.PUT("/:id", roleHandler.UpdateRole)    // PUT    /api/roles/:id
			roles.DELETE("/:id", roleHandler.DeleteRole) // DELETE /api/roles/:id
		}
	}
}