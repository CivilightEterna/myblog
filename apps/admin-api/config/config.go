package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ServerPort   string
	DBPath       string
	JWTSecret    string
	BlogRoot     string
	CORSOrigin   string
	AdminUser    string
	AdminPass    string
}

func Load() *Config {
	blogRoot := os.Getenv("BLOG_ROOT")
	if blogRoot == "" {
		blogRoot = filepath.Join("..", "..")
	}
	absRoot, err := filepath.Abs(blogRoot)
	if err != nil {
		absRoot = blogRoot
	}
	blogRoot = absRoot

	// Database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(blogRoot, "data", "blog_admin.db")
	}

	// JWT Secret: MUST be set in production, auto-generate in dev
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = generateRandomSecret(32)
		fmt.Println("⚠ JWT_SECRET not set, generated random secret (will change on restart)")
		fmt.Println("  Set JWT_SECRET environment variable for production use")
	}

	// CORS origin: restrict in production
	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		// Default: allow local dev servers
		corsOrigin = "http://localhost:5173,http://localhost:4321"
	}

	// Admin credentials: read from env, never use hardcoded defaults in production
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminUser == "" || adminPass == "" {
		fmt.Println("⚠ ADMIN_USERNAME or ADMIN_PASSWORD not set")
		fmt.Println("  Default admin account will be created: admin / admin123")
		fmt.Println("  Change immediately after first login in production!")
		adminUser = "admin"
		adminPass = "admin123"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	cfg := &Config{
		ServerPort: serverPort,
		DBPath:     dbPath,
		JWTSecret:  jwtSecret,
		BlogRoot:   blogRoot,
		CORSOrigin: corsOrigin,
		AdminUser:  adminUser,
		AdminPass:  adminPass,
	}

	cfg.printStartupInfo()
	return cfg
}

func (c *Config) printStartupInfo() {
	isProd := os.Getenv("GIN_MODE") == "release" || os.Getenv("PRODUCTION") == "true"

	if isProd {
		fmt.Println("========================================")
		fmt.Println("  Blog Admin API — Production Mode")
		fmt.Println("========================================")
	} else {
		fmt.Println("========================================")
		fmt.Println("  Blog Admin API — Development Mode")
		fmt.Println("========================================")
	}

	if isProd {
		// In production, verify critical settings
		warnings := []string{}
		if os.Getenv("JWT_SECRET") == "" {
			warnings = append(warnings, "JWT_SECRET is not set — using auto-generated key")
		}
		if os.Getenv("ADMIN_PASSWORD") == "" {
			warnings = append(warnings, "ADMIN_PASSWORD is not set — using default password")
		}
		for _, w := range warnings {
			fmt.Printf("⚠ %s\n", w)
		}
	}
}

func generateRandomSecret(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback — shouldn't happen
		return "fallback-secret-do-not-use-in-production"
	}
	return hex.EncodeToString(bytes)
}
