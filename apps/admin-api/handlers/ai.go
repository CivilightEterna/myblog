package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"blog-admin-api/database"
	"blog-admin-api/models"

	"github.com/gin-gonic/gin"
)

type AIHandler struct{}

func NewAIHandler() *AIHandler {
	return &AIHandler{}
}

func (h *AIHandler) GenerateSummary(c *gin.Context) {
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

	// Check AI config
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "未配置 AI API Key，请设置环境变量 DEEPSEEK_API_KEY"})
		return
	}

	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	model := os.Getenv("DEEPSEEK_MODEL")
	if model == "" {
		model = "deepseek-chat"
	}

	// Build tags string
	var tags []string
	if article.TagsJSON != "" {
		json.Unmarshal([]byte(article.TagsJSON), &tags)
	}

	// Build prompt
	prompt := buildPrompt(article.Title, article.Category, tags, article.Content)

	// Call DeepSeek
	summary, err := callDeepSeek(baseURL, apiKey, model, prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 生成失败: " + err.Error()})
		return
	}

	// Save to article
	now := time.Now()
	article.AiSummary = summary
	article.AiSummaryModel = model
	article.AiSummaryGeneratedAt = &now
	database.DB.Save(&article)

	c.JSON(http.StatusOK, gin.H{
		"ai_summary":             summary,
		"ai_summary_model":       model,
		"ai_summary_generated_at": now,
	})
}

func buildPrompt(title, category string, tags []string, content string) string {
	tagStr := ""
	if len(tags) > 0 {
		for i, t := range tags {
			if i > 0 {
				tagStr += "、"
			}
			tagStr += t
		}
	}

	return fmt.Sprintf(`你是一位技术博客编辑。请为以下文章生成一段中文摘要。

要求：
- 100 到 200 字
- 用词清晰简单
- 不要使用"本文主要介绍了"、"文章指出"这类模板化开头
- 不要编造文章中不存在的信息
- 适合放在技术博客文章顶部
- 保留关键术语和核心观点

文章信息：
- 标题：%s
- 分类：%s
- 标签：%s

正文：
%s`, title, category, tagStr, content)
}

func callDeepSeek(baseURL, apiKey, model, prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens":    600,
		"temperature":   0.3,
		"stream":        false,
	}

	jsonBody, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", baseURL+"/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求 AI 服务失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("AI 服务返回错误 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析 AI 响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("AI 未返回有效结果")
	}

	return result.Choices[0].Message.Content, nil
}
