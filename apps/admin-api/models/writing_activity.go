package models

import "time"

type WritingActivity struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ArticleID      *uint     `json:"article_id"`
	ArticleSlug    string    `json:"article_slug"`
	Action         string    `gorm:"not null" json:"action"`
	WordCount      int       `gorm:"default:0" json:"word_count"`
	DeltaWordCount int       `gorm:"default:0" json:"delta_word_count"`
	ActivityDate   string    `gorm:"not null" json:"activity_date"` // YYYY-MM-DD
	CreatedAt      time.Time `json:"created_at"`
}

func (WritingActivity) TableName() string {
	return "writing_activity"
}
