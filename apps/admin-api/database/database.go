package database

import (
	"fmt"
	"os"
	"path/filepath"

	"blog-admin-api/models"

	"github.com/ncruces/go-sqlite3/gormlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type SeedConfig struct {
	AdminUser string
	AdminPass string
}

func Init(dbPath string, seedCfg SeedConfig) error {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	logLevel := logger.Info
	if os.Getenv("GIN_MODE") == "release" {
		logLevel = logger.Warn
	}

	var err error
	DB, err = gorm.Open(gormlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}

	if err := DB.AutoMigrate(
		&models.AdminUser{},
		&models.ArticleDraft{},
		&models.ArticlePublishRecord{},
		&models.UploadFile{},
		&models.SiteSetting{},
		&models.BuildRecord{},
		&models.WritingActivity{},
	); err != nil {
		return err
	}

	if err := seedAdminUser(seedCfg); err != nil {
		return err
	}

	if err := seedSiteSettings(); err != nil {
		return err
	}

	return nil
}

func seedAdminUser(cfg SeedConfig) error {
	var count int64
	DB.Model(&models.AdminUser{}).Count(&count)
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.AdminUser{
		Username:     cfg.AdminUser,
		PasswordHash: string(hash),
		Nickname:     "管理员",
	}
	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	fmt.Printf("✓ Admin account created: %s (change password after first login)\n", cfg.AdminUser)
	return nil
}

func seedSiteSettings() error {
	var count int64
	DB.Model(&models.SiteSetting{}).Count(&count)
	if count > 0 {
		return nil
	}

	defaults := []models.SiteSetting{
		{SettingKey: "title", SettingValue: "InkSpace"},
		{SettingKey: "description", SettingValue: "一个安静写作的地方"},
		{SettingKey: "author", SettingValue: "Blog Author"},
		{SettingKey: "url", SettingValue: "https://your-domain.com"},
		{SettingKey: "lang", SettingValue: "zh-CN"},
		{SettingKey: "posts_per_page", SettingValue: "10"},
	}

	for _, s := range defaults {
		if err := DB.Create(&s).Error; err != nil {
			return err
		}
	}
	return nil
}
