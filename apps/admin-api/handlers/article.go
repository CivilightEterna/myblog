package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"blog-admin-api/database"
	"blog-admin-api/models"
	"blog-admin-api/services"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	ArticleService *services.ArticleService
}

func NewArticleHandler(articleService *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{ArticleService: articleService}
}

type CreateArticleRequest struct {
	Title          string   `json:"title" binding:"required"`
	Slug           string   `json:"slug" binding:"required"`
	Description    string   `json:"description"`
	Content        string   `json:"content" binding:"required"`
	Cover          string   `json:"cover"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	Draft          bool     `json:"draft"`
	Pinned         bool     `json:"pinned"`
	CommentEnabled bool     `json:"comment_enabled"`
	TocEnabled     bool     `json:"toc_enabled"`
}

// ListArticles returns all articles with optional filters
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	query := database.DB.Model(&models.ArticleDraft{})

	// Filter by draft status
	if draftFilter := c.Query("draft"); draftFilter != "" {
		isDraft := draftFilter == "true"
		query = query.Where("draft = ?", isDraft)
	}

	// Search by title
	if search := c.Query("search"); search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	var articles []models.ArticleDraft
	if err := query.Order("pinned DESC, updated_at DESC").
		Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     articles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetArticle returns a single article by ID
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var article models.ArticleDraft
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// CreateArticle creates a new article draft
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写必要字段"})
		return
	}

	tagsJSON, _ := json.Marshal(req.Tags)

	article := models.ArticleDraft{
		Title:          req.Title,
		Slug:           req.Slug,
		Description:    req.Description,
		Content:        req.Content,
		Cover:          req.Cover,
		Category:       req.Category,
		TagsJSON:       string(tagsJSON),
		Draft:          req.Draft,
		Pinned:         req.Pinned,
		CommentEnabled: req.CommentEnabled,
		TocEnabled:     req.TocEnabled,
		ContentVersion: 1,
		WordCount:      countWords(req.Content),
	}

	if err := database.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败: " + err.Error()})
		return
	}

	// Record writing activity
	h.recordActivity(&article, "create_draft", 0)

	c.JSON(http.StatusCreated, article)
}

// UpdateArticle updates an existing article draft
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var article models.ArticleDraft
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写必要字段"})
		return
	}

	oldWordCount := article.WordCount
	tagsJSON, _ := json.Marshal(req.Tags)
	newWordCount := countWords(req.Content)

	article.Title = req.Title
	article.Slug = req.Slug
	article.Description = req.Description
	article.Content = req.Content
	article.Cover = req.Cover
	article.Category = req.Category
	article.TagsJSON = string(tagsJSON)
	article.Draft = req.Draft
	article.Pinned = req.Pinned
	article.CommentEnabled = req.CommentEnabled
	article.TocEnabled = req.TocEnabled
	article.WordCount = newWordCount

	if err := database.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	// Record writing activity
	delta := newWordCount - oldWordCount
	action := "update_draft"
	if !article.Draft {
		action = "update_published"
	}
	h.recordActivity(&article, action, delta)

	c.JSON(http.StatusOK, article)
}

// DeleteArticle deletes an article draft
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var article models.ArticleDraft
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	h.recordActivity(&article, "delete_article", -article.WordCount)

	if err := database.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章已删除"})
}

// PublishArticle publishes an article draft as a Markdown file
func (h *ArticleHandler) PublishArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var article models.ArticleDraft
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// Generate Markdown file
	mdPath, err := h.ArticleService.PublishMarkdown(&article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败: " + err.Error()})
		return
	}

	// Update draft status
	now := timeNow()
	article.Draft = false
	article.PublishedAt = &now
	database.DB.Save(&article)

	// Record publish record
	record := models.ArticlePublishRecord{
		DraftID:      article.ID,
		Title:        article.Title,
		Slug:         article.Slug,
		MarkdownPath: mdPath,
	}
	database.DB.Create(&record)

	// Record writing activity
	h.recordActivity(&article, "publish_article", 0)

	c.JSON(http.StatusOK, gin.H{
		"message":       "发布成功",
		"markdown_path": mdPath,
	})
}

// UnpublishArticle sets an article back to draft
func (h *ArticleHandler) UnpublishArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var article models.ArticleDraft
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	article.Draft = true
	article.PublishedAt = nil
	database.DB.Save(&article)

	h.recordActivity(&article, "unpublish_article", 0)

	c.JSON(http.StatusOK, gin.H{"message": "已取消发布"})
}

func (h *ArticleHandler) recordActivity(article *models.ArticleDraft, action string, deltaWords int) {
	activity := models.WritingActivity{
		ArticleID:      &article.ID,
		ArticleSlug:    article.Slug,
		Action:         action,
		WordCount:      article.WordCount,
		DeltaWordCount: deltaWords,
		ActivityDate:   timeNow().Format("2006-01-02"),
	}
	database.DB.Create(&activity)
}

func timeNow() time.Time {
	return time.Now()
}

func countWords(content string) int {
	// Count Chinese characters and English words
	count := 0
	inWord := false
	for _, r := range content {
		if r >= 0x4e00 && r <= 0x9fff {
			count++
			inWord = false
		} else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			if !inWord {
				count++
				inWord = true
			}
		} else {
			inWord = false
		}
	}
	return count
}
