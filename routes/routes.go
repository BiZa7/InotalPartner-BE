package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"inotal-be/config"
	"inotal-be/internal/handler"
	"inotal-be/internal/middleware"
	"inotal-be/internal/service"
)

func Setup(
	r *gin.Engine,
	authSvc *service.AuthService,
	userSvc *service.UserService,
	roleSvc *service.RoleService,
	categorySvc *service.CategoryService,
	tagSvc *service.TagService,
	articleSvc *service.ArticleService,
) {
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
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	tagHandler := handler.NewTagHandler(tagSvc)
	articleHandler := handler.NewArticleHandler(articleSvc)

	// API v1
	api := r.Group("/api")
	{
		// ── Auth ────────────────────────────────────────────────────────────
		auth := api.Group("/auth")
		{
			// Public routes
			auth.POST("/setup", authHandler.Setup) // ← first-time setup (DB kosong)
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

		// ── Categories ──────────────────────────────────────────────────────
		categories := api.Group("/categories")
		categories.Use(middleware.AuthMiddleware(authSvc))
		{
			operatorCategories := categories.Group("")
			operatorCategories.Use(middleware.OperatorOnly())
			{
				operatorCategories.GET("", categoryHandler.GetCategories)      // GET    /api/categories
				operatorCategories.POST("", categoryHandler.CreateCategory)    // POST   /api/categories
				operatorCategories.GET("/:id", categoryHandler.GetCategory)    // GET    /api/categories/:id
				operatorCategories.PUT("/:id", categoryHandler.UpdateCategory) // PUT    /api/categories/:id
			}

			// DELETE category ❌ (Operator, Admin, SuperAdmin, Guest tidak punya akses delete)
			categories.DELETE("/:id", middleware.DenyAll(), categoryHandler.DeleteCategory)
		}

		// ── Tags ────────────────────────────────────────────────────────────
		tags := api.Group("/tags")
		tags.Use(middleware.AuthMiddleware(authSvc))
		{
			operatorTags := tags.Group("")
			operatorTags.Use(middleware.OperatorOnly())
			{
				operatorTags.GET("", tagHandler.GetTags)       // GET    /api/tags
				operatorTags.POST("", tagHandler.CreateTag)    // POST   /api/tags
				operatorTags.GET("/:id", tagHandler.GetTag)    // GET    /api/tags/:id
				operatorTags.PUT("/:id", tagHandler.UpdateTag) // PUT    /api/tags/:id
			}

			// DELETE tag ❌ (Operator, Admin, SuperAdmin, Guest tidak punya akses delete)
			tags.DELETE("/:id", middleware.DenyAll(), tagHandler.DeleteTag)
		}

		// ── Articles ────────────────────────────────────────────────────────
		articles := api.Group("/articles")
		articles.Use(middleware.AuthMiddleware(authSvc))
		{
			articles.GET("", middleware.RequireRole("operator", "admin", "super_admin"), articleHandler.GetArticles)    // GET  /api/articles
			articles.POST("", middleware.OperatorOnly(), articleHandler.CreateArticle)                                  // POST /api/articles (only Operator)
			articles.GET("/:id", middleware.RequireRole("operator", "admin", "super_admin"), articleHandler.GetArticle) // GET  /api/articles/:id
			articles.PUT("/:id", middleware.OperatorOnly(), articleHandler.UpdateArticle)                               // PUT  /api/articles/:id (only Operator)
			articles.POST("/:id/publish", middleware.AdminOrAbove(), articleHandler.PublishArticle)                     // POST /api/articles/:id/publish (only Admin & Super Admin)
			articles.POST("/:id/takedown", middleware.AdminOrAbove(), articleHandler.TakedownArticle)                   // POST /api/articles/:id/takedown (only Admin & Super Admin)
		}

		// ── Public Content (No Auth Required) ──────────────────────────────
		publicArticleHandler := handler.NewPublicArticleHandler(articleSvc)
		public := api.Group("/public")
		{
			// News
			public.GET("/news", publicArticleHandler.GetPublicNews)            // GET /api/public/news
			public.GET("/news/:slug", publicArticleHandler.GetPublicNewsBySlug) // GET /api/public/news/:slug

			// Events
			public.GET("/events", publicArticleHandler.GetPublicEvents)            // GET /api/public/events
			public.GET("/events/:slug", publicArticleHandler.GetPublicEventsBySlug) // GET /api/public/events/:slug

			// Articles
			public.GET("/articles", publicArticleHandler.GetPublicArticles)            // GET /api/public/articles
			public.GET("/articles/:slug", publicArticleHandler.GetPublicArticleBySlug) // GET /api/public/articles/:slug
		}
	}
}
