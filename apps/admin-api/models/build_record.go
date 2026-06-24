package models

import "time"

type BuildRecord struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Status      string     `gorm:"not null" json:"status"` // running, success, failed
	TriggerType string     `json:"trigger_type"`           // manual, auto
	Log         string     `json:"log"`
	ReleasePath string     `json:"release_path"`
	StartedAt   time.Time  `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at"`
}

func (BuildRecord) TableName() string {
	return "build_record"
}
