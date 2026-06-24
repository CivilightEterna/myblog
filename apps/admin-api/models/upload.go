package models

import "time"

type UploadFile struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Filename     string    `gorm:"not null" json:"filename"`
	OriginalName string    `json:"original_name"`
	Path         string    `gorm:"not null" json:"path"`
	URL          string    `gorm:"not null" json:"url"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	Hash         string    `json:"hash"`
	CreatedAt    time.Time `json:"created_at"`
}

func (UploadFile) TableName() string {
	return "upload_file"
}
