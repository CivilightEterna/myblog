package handlers

import (
	"net/http"
	"strconv"

	"blog-admin-api/database"
	"blog-admin-api/models"
	"blog-admin-api/services"

	"github.com/gin-gonic/gin"
)

type BuildHandler struct {
	BuildService *services.BuildService
}

func NewBuildHandler(buildService *services.BuildService) *BuildHandler {
	return &BuildHandler{BuildService: buildService}
}

func (h *BuildHandler) TriggerBuild(c *gin.Context) {
	var req struct {
		TriggerType string `json:"trigger_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.TriggerType = "manual"
	}
	if req.TriggerType == "" {
		req.TriggerType = "manual"
	}

	record, err := h.BuildService.ExecuteBuild(req.TriggerType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "构建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (h *BuildHandler) ListBuilds(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.BuildRecord{}).Count(&total)

	var builds []models.BuildRecord
	if err := database.DB.Order("started_at DESC").Offset(offset).Limit(pageSize).Find(&builds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询构建记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     builds,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *BuildHandler) GetBuild(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的构建ID"})
		return
	}

	var record models.BuildRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "构建记录不存在"})
		return
	}

	c.JSON(http.StatusOK, record)
}
