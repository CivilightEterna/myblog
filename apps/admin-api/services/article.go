package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog-admin-api/models"
)

type ArticleService struct {
	BlogRoot string
}

func NewArticleService(blogRoot string) *ArticleService {
	return &ArticleService{BlogRoot: blogRoot}
}

// PublishMarkdown generates a Markdown file with frontmatter from an article draft
func (s *ArticleService) PublishMarkdown(article *models.ArticleDraft) (string, error) {
	// Determine target directory: apps/blog-web/src/content/posts/YYYY/MM/
	now := time.Now()
	if article.PublishedAt != nil {
		now = *article.PublishedAt
	}
	yearMonth := now.Format("2006/01")
	postsDir := filepath.Join(s.BlogRoot, "apps", "blog-web", "src", "content", "posts", yearMonth)

	if err := os.MkdirAll(postsDir, 0755); err != nil {
		return "", fmt.Errorf("创建文章目录失败: %w", err)
	}

	// Parse tags
	var tags []string
	if article.TagsJSON != "" {
		json.Unmarshal([]byte(article.TagsJSON), &tags)
	}

	// Generate frontmatter
	frontmatter := s.buildFrontmatter(article, tags, now)

	// Generate full Markdown file
	mdContent := frontmatter + "\n" + article.Content

	// Write to file
	filePath := filepath.Join(postsDir, article.Slug+".md")
	if err := os.WriteFile(filePath, []byte(mdContent), 0644); err != nil {
		return "", fmt.Errorf("写入Markdown文件失败: %w", err)
	}

	return filePath, nil
}

func (s *ArticleService) buildFrontmatter(article *models.ArticleDraft, tags []string, now time.Time) string {
	dateStr := now.Format("2006-01-02")
	updatedStr := article.UpdatedAt.Format("2006-01-02")
	if updatedStr == dateStr {
		updatedStr = dateStr
	}

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("contentVersion: %d\n", article.ContentVersion))
	sb.WriteString(fmt.Sprintf("title: \"%s\"\n", escapeYAML(article.Title)))
	sb.WriteString(fmt.Sprintf("slug: \"%s\"\n", escapeYAML(article.Slug)))
	sb.WriteString(fmt.Sprintf("description: \"%s\"\n", escapeYAML(article.Description)))
	sb.WriteString(fmt.Sprintf("date: \"%s\"\n", dateStr))
	sb.WriteString(fmt.Sprintf("updated: \"%s\"\n", updatedStr))
	sb.WriteString(fmt.Sprintf("category: \"%s\"\n", escapeYAML(article.Category)))

	if len(tags) > 0 {
		sb.WriteString("tags:\n")
		for _, tag := range tags {
			sb.WriteString(fmt.Sprintf("  - %s\n", tag))
		}
	} else {
		sb.WriteString("tags: []\n")
	}

	if article.Cover != "" {
		sb.WriteString(fmt.Sprintf("cover: \"%s\"\n", escapeYAML(article.Cover)))
	} else {
		sb.WriteString("cover: \"\"\n")
	}

	sb.WriteString(fmt.Sprintf("draft: %t\n", article.Draft))
	sb.WriteString(fmt.Sprintf("pinned: %t\n", article.Pinned))
	sb.WriteString(fmt.Sprintf("comment: %t\n", article.CommentEnabled))
	sb.WriteString(fmt.Sprintf("toc: %t\n", article.TocEnabled))
	if article.AiSummary != "" {
		sb.WriteString(fmt.Sprintf("aiSummary: \"%s\"\n", escapeYAML(article.AiSummary)))
	} else {
		sb.WriteString("aiSummary: \"\"\n")
	}
	sb.WriteString("---\n")

	return sb.String()
}

func escapeYAML(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
