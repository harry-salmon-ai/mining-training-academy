package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"go-backend-react-frontend/internal/auth"
	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/handlers"
	"go-backend-react-frontend/internal/middleware"
	"go-backend-react-frontend/internal/seed"

	"github.com/gin-gonic/gin"
)

//go:generate sh -c "cd frontend && npm install && npm run build"
//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	auth.Init()

	database := db.Init()
	db.Migrate()

	if os.Getenv("SEED") == "true" {
		seed.Run(database, true)
	}
	// Always seed if no users exist
	var userCount int64
	database.Table("users").Count(&userCount)
	if userCount == 0 {
		seed.Run(database, false)
	}

	r := gin.Default()

	// CORS for dev mode
	if os.Getenv("ENV") == "dev" {
		r.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})
	}

	api := r.Group("/api")
	{
		// Public auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", handlers.Register)
			authGroup.POST("/login", handlers.Login)
			authGroup.POST("/logout", handlers.Logout)
			authGroup.GET("/google", handlers.GoogleLogin)
			authGroup.GET("/google/callback", handlers.GoogleCallback)
			authGroup.GET("/me", auth.AuthMiddleware(), handlers.GetMe)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			// Modules (read accessible by all authenticated users)
			protected.GET("/modules", handlers.ListModules)
			protected.GET("/modules/:id", handlers.GetModule)
			protected.GET("/modules/:id/sections", handlers.ListSections)

			// Module write operations (requires permission)
			protectedModules := protected.Group("/modules")
			protectedModules.Use(middleware.RequirePermission("modules", "create"))
			{
				protectedModules.POST("", handlers.CreateModule)
			}
			protected.PATCH("/modules/:id", middleware.RequirePermission("modules", "update"), handlers.UpdateModule)
			protected.DELETE("/modules/:id", middleware.RequirePermission("modules", "delete"), handlers.DeleteModule)
			protected.POST("/modules/:id/publish", middleware.RequirePermission("modules", "publish"), handlers.PublishModule)
			protected.POST("/modules/:id/archive", middleware.RequirePermission("modules", "update"), handlers.ArchiveModule)
			protected.POST("/modules/:id/sections", middleware.RequirePermission("modules", "update"), handlers.CreateSection)

			// SME notes on slides (admin only)
			protected.GET("/slides/:id/notes", middleware.RequireAdmin(), handlers.ListSlideNotes)
			protected.POST("/slides/:id/notes", middleware.RequireAdmin(), handlers.CreateSlideNote)
			protected.PATCH("/notes/:id", middleware.RequireAdmin(), handlers.UpdateNote)
			protected.DELETE("/notes/:id", middleware.RequireAdmin(), handlers.DeleteNote)

			// Sections
			protected.PATCH("/sections/:id", middleware.RequirePermission("modules", "update"), handlers.UpdateSection)
			protected.DELETE("/sections/:id", middleware.RequirePermission("modules", "delete"), handlers.DeleteSection)
			protected.POST("/sections/:id/slides", middleware.RequirePermission("modules", "update"), handlers.CreateSlide)

			// Slides
			protected.GET("/slides/:id", handlers.GetSlide)
			protected.PATCH("/slides/:id", middleware.RequirePermission("modules", "update"), handlers.UpdateSlide)
			protected.DELETE("/slides/:id", middleware.RequirePermission("modules", "delete"), handlers.DeleteSlide)

			// Categories
			protected.GET("/categories", handlers.ListCategories)
			protected.POST("/categories", middleware.RequirePermission("categories", "create"), handlers.CreateCategory)
			protected.PATCH("/categories/:id", middleware.RequirePermission("categories", "update"), handlers.UpdateCategory)
			protected.DELETE("/categories/:id", middleware.RequirePermission("categories", "delete"), handlers.DeleteCategory)

			// Quizzes
			protected.GET("/quizzes/:id", handlers.GetQuiz)
			protected.POST("/quizzes", middleware.RequirePermission("modules", "create"), handlers.CreateQuiz)
			protected.PATCH("/quizzes/:id", middleware.RequirePermission("modules", "update"), handlers.UpdateQuiz)
			protected.DELETE("/quizzes/:id", middleware.RequirePermission("modules", "delete"), handlers.DeleteQuiz)
			protected.POST("/quizzes/:id/questions", middleware.RequirePermission("modules", "update"), handlers.AddQuestion)
			protected.POST("/quizzes/:id/attempt", handlers.StartQuizAttempt)
			protected.PATCH("/quizzes/:id/attempt/:attemptId", handlers.SubmitQuizAttempt)
			protected.GET("/quizzes/:id/results", handlers.GetQuizResults)

			// Enrollments
			protected.GET("/enrollments", handlers.ListEnrollments)
			protected.POST("/enrollments", middleware.RequirePermission("enrollments", "create"), handlers.CreateEnrollment)
			protected.POST("/enrollments/bulk", middleware.RequirePermission("enrollments", "assign"), handlers.BulkEnroll)
			protected.POST("/enrollments/self", handlers.SelfEnroll)
			protected.PATCH("/enrollments/:id", middleware.RequirePermission("enrollments", "update"), handlers.UpdateEnrollment)
			protected.DELETE("/enrollments/:id", middleware.RequirePermission("enrollments", "delete"), handlers.DeleteEnrollment)

			// Progress
			protected.POST("/progress/slide", handlers.RecordSlideProgress)
			protected.GET("/progress/module/:id", handlers.GetModuleProgress)
			protected.GET("/progress/overview", handlers.GetProgressOverview)

			// Users (admin)
			protected.GET("/users", middleware.RequirePermission("users", "read"), handlers.ListUsers)
			protected.POST("/users", middleware.RequirePermission("users", "create"), handlers.CreateUser)
			protected.GET("/users/:id", middleware.RequirePermission("users", "read"), handlers.GetUser)
			protected.PATCH("/users/:id", middleware.RequirePermission("users", "update"), handlers.UpdateUser)
			protected.DELETE("/users/:id", middleware.RequirePermission("users", "delete"), handlers.DeactivateUser)
			protected.PATCH("/profile", handlers.UpdateProfile)

			// Certificates
			protected.GET("/certificates", handlers.ListCertificates)
			protected.POST("/certificates/generate", middleware.RequireAdmin(), handlers.GenerateCertificate)

			// Reports
			protected.GET("/reports/dashboard", handlers.GetDashboardStats)

			// Notifications
			protected.GET("/notifications", handlers.ListNotifications)
			protected.PATCH("/notifications/:id/read", handlers.MarkNotificationRead)
			protected.POST("/notifications/mark-all-read", handlers.MarkAllNotificationsRead)

			// Audit
			protected.GET("/audit", middleware.RequirePermission("audit", "read"), handlers.ListAuditLogs)

			// Settings
			protected.GET("/settings", handlers.GetSettings)
			protected.PATCH("/settings", middleware.RequirePermission("settings", "update"), handlers.UpdateSettings)

			// Media upload (admin only)
			protected.POST("/media/upload", middleware.RequireAdmin(), handlers.UploadMedia)
		}
	}

	// Media serve — UUID-named files, no auth required
	r.GET("/api/media/:filename", handlers.ServeMedia)

	// Serve embedded frontend in production
	if os.Getenv("ENV") == "dev" {
		log.Println("Running in dev mode - frontend should be served by Vite on :3000")
	} else {
		distFS, err := fs.Sub(frontendFS, "frontend/dist")
		if err != nil {
			log.Fatal(err)
		}
		// Pre-read index.html to avoid http.FileServer's index.html→"./" redirect loop
		indexHTML, err := frontendFS.ReadFile("frontend/dist/index.html")
		if err != nil {
			log.Fatal(err)
		}
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// Try to serve an actual asset file (JS, CSS, images, etc.)
			f, err := http.FS(distFS).Open(path)
			if err != nil {
				// SPA fallback — serve index.html bytes directly (no http.FileServer redirect)
				c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
				return
			}
			stat, statErr := f.Stat()
			f.Close()
			if statErr != nil || stat.IsDir() {
				c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
				return
			}
			c.FileFromFS(path, http.FS(distFS))
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Server starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
