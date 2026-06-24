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
