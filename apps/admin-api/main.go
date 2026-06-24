package main

import (
	"fmt"
	"log"

	"blog-admin-api/config"
	"blog-admin-api/database"
	"blog-admin-api/router"
)

func main() {
	cfg := config.Load()

	// Initialize database
	if err := database.Init(cfg.DBPath, database.SeedConfig{
		AdminUser: cfg.AdminUser,
		AdminPass: cfg.AdminPass,
	}); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized")

	// Setup router
	r := router.Setup(cfg)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Admin API server starting on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
