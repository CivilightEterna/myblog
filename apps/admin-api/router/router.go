package router

import (
	"net/http"
	"strings"

	"blog-admin-api/config"
	"blog-admin-api/handlers"
	"blog-admin-api/middleware"
	"blog-admin-api/services"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS with configurable origins — never use * in production
	r.Use(corsMiddleware(cfg.CORSOrigin))

	// Security headers
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})

	// Initialize services
	articleService := services.NewArticleService(cfg.BlogRoot)
	buildService := services.NewBuildService(cfg.BlogRoot)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg.JWTSecret)
	articleHandler := handlers.NewArticleHandler(articleService)
	uploadHandler := handlers.NewUploadHandler(cfg.BlogRoot)
	settingHandler := handlers.NewSettingHandler()
	buildHandler := handlers.NewBuildHandler(buildService)

	api := r.Group("/api/admin")
	{
		// Rate limit login attempts (simple — in production use a proper rate limiter)
		api.POST("/login", authHandler.Login)

		// All other routes require JWT
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Auth
			protected.GET("/profile", authHandler.Profile)

			// Articles
			protected.GET("/articles", articleHandler.ListArticles)
			protected.GET("/articles/:id", articleHandler.GetArticle)
			protected.POST("/articles", articleHandler.CreateArticle)
			protected.PUT("/articles/:id", articleHandler.UpdateArticle)
			protected.DELETE("/articles/:id", articleHandler.DeleteArticle)
			protected.POST("/articles/:id/publish", articleHandler.PublishArticle)
			protected.POST("/articles/:id/unpublish", articleHandler.UnpublishArticle)

			// Uploads
			protected.POST("/uploads", uploadHandler.Upload)
			protected.GET("/uploads", uploadHandler.ListUploads)
			protected.DELETE("/uploads/:id", uploadHandler.DeleteUpload)

			// Settings
			protected.GET("/settings", settingHandler.GetSettings)
			protected.PUT("/settings", settingHandler.UpdateSettings)

			// Builds
			protected.POST("/build", buildHandler.TriggerBuild)
			protected.GET("/builds", buildHandler.ListBuilds)
			protected.GET("/builds/:id", buildHandler.GetBuild)
		}
	}

	return r
}

// corsMiddleware validates the origin against the allowed list
func corsMiddleware(allowedOrigins string) gin.HandlerFunc {
	allowed := strings.Split(allowedOrigins, ",")
	for i := range allowed {
		allowed[i] = strings.TrimSpace(allowed[i])
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed_ := false
		for _, o := range allowed {
			if o == origin || o == "*" {
				allowed_ = true
				break
			}
		}

		if allowed_ {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if origin != "" {
			// Don't set Allow-Origin for disallowed origins — browser will block
			c.Header("Access-Control-Allow-Origin", allowed[0])
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
