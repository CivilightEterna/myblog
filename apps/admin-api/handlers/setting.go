package handlers

import (
	"net/http"

	"blog-admin-api/database"
	"blog-admin-api/models"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct{}

func NewSettingHandler() *SettingHandler {
	return &SettingHandler{}
}

func (h *SettingHandler) GetSettings(c *gin.Context) {
	var settings []models.SiteSetting
	if err := database.DB.Order("id ASC").Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询设置失败"})
		return
	}

	// Return as key-value map
	result := make(map[string]string)
	for _, s := range settings {
		result[s.SettingKey] = s.SettingValue
	}

	c.JSON(http.StatusOK, result)
}

func (h *SettingHandler) UpdateSettings(c *gin.Context) {
	var updates map[string]string
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的设置数据"})
		return
	}

	for key, value := range updates {
		var setting models.SiteSetting
		result := database.DB.Where("setting_key = ?", key).First(&setting)
		if result.Error != nil {
			// Create new setting
			setting = models.SiteSetting{
				SettingKey:   key,
				SettingValue: value,
			}
			database.DB.Create(&setting)
		} else {
			setting.SettingValue = value
			database.DB.Save(&setting)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置已更新"})
}

// GetPublicConfig returns whitelisted settings for the public frontend.
// NO auth required — only safe, non-sensitive keys are exposed.
func (h *SettingHandler) GetPublicConfig(c *gin.Context) {
	// Whitelist: only these keys are safe to expose publicly
	whitelist := map[string]bool{
		"site_name":                true,
		"site_description":         true,
		"site_author":              true,
		"site_background_enabled":  true,
		"site_background_image":    true,
		"site_background_opacity":  true,
		"site_background_mode":     true,
		"site_background_position": true,
		"comment_enabled":           true,
		"comment_provider":          true,
		"twikoo_env_id":             true,
		"twikoo_server_url":           true,
		"landing_enabled":             true,
		"landing_cover_image":         true,
		"landing_title":               true,
		"landing_subtitle":            true,
		"landing_description":         true,
		"landing_button_text":         true,
		"landing_button_link":         true,
		"landing_overlay_opacity":     true,
		"landing_text_position":       true,
		"landing_quote_enabled":       true,
		"landing_quote_text":          true,
		"landing_quote_author":        true,
		"landing_quote_typing_effect": true,
		"landing_animation_enabled":   true,
		"landing_animation_style":     true,
		"landing_animation_intensity": true,
		"landing_anime_style_enabled": true,
		"landing_blur_card_enabled":   true,
		"landing_card_opacity":        true,
		"landing_accent_color":        true,
		"title":                       true,
		"site_avatar":                 true,
		"site_logo":                   true,
		"description":                 true,
		"author":                      true,
		"lang":                        true,
		"url":                         true,
	}

	var settings []models.SiteSetting
	if err := database.DB.Order("id ASC").Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询配置失败"})
		return
	}

	result := make(map[string]string)
	for _, s := range settings {
		if whitelist[s.SettingKey] {
			result[s.SettingKey] = s.SettingValue
		}
	}

	// Apply defaults for missing keys
	defaults := map[string]string{
		"site_background_enabled":  "false",
		"site_background_image":    "",
		"site_background_opacity":  "0.18",
		"site_background_mode":     "cover",
		"site_background_position": "center",
		"comment_enabled":          "false",
		"comment_provider":         "none",
		"twikoo_env_id":            "",
		"twikoo_server_url":             "",
		"landing_enabled":               "true",
		"landing_cover_image":           "",
		"landing_title":                 "InkSpace",
		"landing_subtitle":              "写作、技术与生活的私人空间",
		"landing_description":           "记录后端、数据库、AI 与那些偶然闪光的想法。",
		"landing_button_text":           "进入博客",
		"landing_button_link":           "/blog",
		"landing_overlay_opacity":       "0.35",
		"landing_text_position":         "center",
		"landing_quote_enabled":         "true",
		"landing_quote_text":            "愿你在代码与星光之间，找到自己的小宇宙。",
		"landing_quote_author":          "",
		"landing_quote_typing_effect":   "true",
		"landing_animation_enabled":     "true",
		"landing_animation_style":       "sakura",
		"landing_animation_intensity":   "medium",
		"landing_anime_style_enabled":   "true",
		"landing_blur_card_enabled":     "true",
		"landing_card_opacity":          "0.35",
		"landing_accent_color":          "#8b5cf6",
		"site_name":                     "InkSpace",
		"site_description":              "一个安静写作的地方",
		"site_author":                   "",
		"site_avatar":                   "",
		"site_logo":                     "",
	}
	for k, v := range defaults {
		if _, ok := result[k]; !ok {
			result[k] = v
		}
	}

	c.JSON(http.StatusOK, result)
}
