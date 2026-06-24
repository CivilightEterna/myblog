package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"blog-admin-api/config"
	"blog-admin-api/database"
	"blog-admin-api/models"
	"blog-admin-api/router"
)

func main() {
	exportPath := flag.String("export-settings", "", "Export all settings as JSON to the given file path")
	flag.Parse()

	cfg := config.Load()

	// Initialize database
	if err := database.Init(cfg.DBPath, database.SeedConfig{
		AdminUser: cfg.AdminUser,
		AdminPass: cfg.AdminPass,
	}); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// If -export-settings flag is set, export and exit
	if *exportPath != "" {
		if err := exportSettings(*exportPath); err != nil {
			log.Fatalf("Failed to export settings: %v", err)
		}
		log.Printf("Settings exported to %s", *exportPath)
		return
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

func exportSettings(path string) error {
	var settings []models.SiteSetting
	if err := database.DB.Order("id ASC").Find(&settings).Error; err != nil {
		return fmt.Errorf("query settings: %w", err)
	}

	data := make(map[string]string)
	for _, s := range settings {
		data[s.SettingKey] = s.SettingValue
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
