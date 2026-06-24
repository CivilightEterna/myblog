package models

import "time"

type ArticleDraft struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Title           string     `gorm:"not null" json:"title"`
	Slug            string     `gorm:"not null" json:"slug"`
	Description     string     `json:"description"`
	Content         string     `gorm:"not null" json:"content"`
	Cover           string     `json:"cover"`
	Category        string     `json:"category"`
	TagsJSON        string     `json:"tags_json"` // JSON array stored as string
	Draft           bool       `gorm:"default:true" json:"draft"`
	Pinned          bool       `gorm:"default:false" json:"pinned"`
	CommentEnabled  bool       `gorm:"default:true" json:"comment_enabled"`
	TocEnabled      bool       `gorm:"default:true" json:"toc_enabled"`
	ContentVersion     int        `gorm:"default:1" json:"content_version"`
	WordCount          int        `gorm:"default:0" json:"word_count"`
	AiSummary          string     `json:"ai_summary"`
	AiSummaryGeneratedAt *time.Time `json:"ai_summary_generated_at"`
	AiSummaryModel     string     `json:"ai_summary_model"`
	PublishedAt        *time.Time `json:"published_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (ArticleDraft) TableName() string {
	return "article_draft"
}

type ArticlePublishRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DraftID       uint      `json:"draft_id"`
	Title         string    `gorm:"not null" json:"title"`
	Slug          string    `gorm:"not null" json:"slug"`
	MarkdownPath  string    `gorm:"not null" json:"markdown_path"`
	PublishedAt   time.Time `json:"published_at"`
	BuildRecordID *uint     `json:"build_record_id"`
}

func (ArticlePublishRecord) TableName() string {
	return "article_publish_record"
}
